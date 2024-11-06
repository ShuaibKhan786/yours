package backend

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ShuaibKhan786/yours/cmd/yyd/channel"
	"github.com/ShuaibKhan786/yours/cmd/yyd/global"
	customimage "github.com/ShuaibKhan786/yours/pkg/images"
	"github.com/ShuaibKhan786/yours/pkg/workerpool"
	"github.com/ShuaibKhan786/yours/pkg/yt"
	"github.com/kkdai/youtube/v2"
)

// TODO: refator to sperate file
type DownloadTask struct {
	IsPlaylist bool
	RootCtx    context.Context
	Ctx        context.Context
	CtxManager map[string]context.CancelFunc
	Video      *youtube.Video
	Dir        string
	ItagNo     int
	Chan       *channel.Channel
}

func (dt *DownloadTask) Execute() error {
	// priortizing context cancellation first
	select {
	case <-dt.RootCtx.Done():
		return nil
	default:
	}

	err := yt.Download(dt.Ctx, dt.Video, dt.Dir, dt.ItagNo)
	//check the rootCTX being cancelled so that
	//it prevents from sending to a closed channel which cause panic
	if dt.RootCtx.Err() != nil {
		return nil
	}
	if err != nil {
		if dt.IsPlaylist {
			dt.Chan.GUIChannel <- channel.DownloadCancelDone(dt.Video.ID)
		} else {
			dt.Chan.GUIChannel <- channel.DownloadCancelDone("")
		}
	} else {
		if dt.IsPlaylist {
			dt.Chan.GUIChannel <- channel.DownloadDone(dt.Video.ID)
		} else {
			dt.Chan.GUIChannel <- channel.DownloadDone("")
		}
	}

	if cancel, ok := dt.CtxManager[dt.Video.ID]; ok {
		cancel()                           // Cancel the context
		delete(dt.CtxManager, dt.Video.ID) // Remove from the map after cancellation
	}
	return nil
}

// TODO: refator to sperate file
type DownloadImage struct {
	RootCtx     context.Context
	PlaylistMap *global.PlaylistMap
	Video       *youtube.Video
}

func (dt *DownloadImage) Execute() error {
	// priortizing context cancellation first
	select {
	case <-dt.RootCtx.Done():
		return nil
	default:
	}

	ctx, cancel := context.WithTimeout(dt.RootCtx, 10*time.Second)
	defer cancel()

	img, err := customimage.DecodeImageFromURI(ctx, dt.Video.Thumbnails[0].URL)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	arrFormat := global.GetArrayOfFormatWithSize(dt.Video)

	dt.PlaylistMap.Mu.Lock()
	dt.PlaylistMap.Details[dt.Video.ID] = &global.MediaDetails{
		Thumbnail: img,
		Formats:   arrFormat,
	}
	dt.PlaylistMap.Mu.Unlock()

	return nil
}

type DownloadAllPlaylistTask struct {
	RootCtx    context.Context
	Ctx        context.Context
	CtxManager map[string]context.CancelFunc
	Video      *youtube.Video
	Dir        string
	ItagNo     int
	Chan       *channel.Channel
}

func (dt *DownloadAllPlaylistTask) Execute() error {
	// priortizing context cancellation first
	select {
	case <-dt.RootCtx.Done():
		return nil
	default:
	}

	if dt.Ctx.Err() != nil {
		return nil
	}

	err := yt.Download(dt.Ctx, dt.Video, dt.Dir, dt.ItagNo)
	//check the rootCTX being cancelled so that
	//it prevents from sending to a closed channel which cause panic
	if dt.RootCtx.Err() != nil {
		return nil
	}
	if err != nil {
		fmt.Println(err.Error())
		if err == context.Canceled {
			dt.Chan.GUIChannel <- channel.PlaylistDownloadAllCancelDone{}
		} else {
			dt.Chan.GUIChannel <- channel.PlaylistDownloadAllError{}
		}
	} else {
		dt.Chan.GUIChannel <- channel.PlaylistDownloadAllDone{}
	}

	if cancel, ok := dt.CtxManager[dt.Video.ID]; ok {
		cancel()                           // Cancel the context
		delete(dt.CtxManager, dt.Video.ID) // Remove from the map after cancellation
	}
	return nil
}

func InitBackend(rootCtx context.Context, c *channel.Channel) {
	ctxManager := make(map[string]context.CancelFunc, 0)

	wp := workerpool.NewWorkerPool(rootCtx, 16)

	//****ONE GO ROUTINE RUNNING HERE****
	go func() {
		for {
			select {
			case i := <-c.BackendChannel:
				switch ct := i.(type) {
				case channel.LinkChannel:
					c.GUIChannel <- channel.FetchMDStarted{}
					ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
					defer cancel()
					mdInterface, err := yt.GetYTMetadata(ctx, string(ct))
					if err != nil {
						if errors.Is(err, yt.ErrInvalidYTLink) {
							c.GUIChannel <- channel.ErrorChannel(err.Error())
						} else {
							c.GUIChannel <- channel.ErrorChannel("Sorry, YYD cant proceed for now try later")
						}
					} else {
						c.GUIChannel <- mdInterface
					}

				case *channel.DownloadChannel:
					ctx, cancel := context.WithCancel(rootCtx)
					ctxManager[ct.Video.ID] = cancel
					dt := &DownloadTask{
						IsPlaylist: ct.IsPlaylist,
						RootCtx:    rootCtx,
						Ctx:        ctx,
						CtxManager: ctxManager,
						Video:      ct.Video,
						Dir:        ct.Dir,
						ItagNo:     ct.ItagNo,
						Chan:       c,
					}
					wp.Add(dt)

				case channel.CancelDownload:
					fmt.Println("Will need to canceled")
					if cancel, ok := ctxManager[string(ct)]; ok {
						cancel()                       // Cancel the context
						delete(ctxManager, string(ct)) // Remove from the map after cancellation
					}
				case *channel.DownloadPlaylistImages:
					ctx, cancel := context.WithCancel(rootCtx)
					ctxManager["playlistImageDownload"] = cancel
					for _, video := range ct.MD.Videos {
						di := &DownloadImage{
							RootCtx:     ctx,
							PlaylistMap: ct.PlaylistMap,
							Video:       video,
						}
						wp.Add(di)
					}
				case channel.CancelPlaylistImageDownload:
					if cancel, ok := ctxManager["playlistImageDownload"]; ok {
						cancel()                                    // Cancel the context
						delete(ctxManager, "playlistImageDownload") // Remove from the map after cancellation
					}
				case *channel.DownloadAllPlaylist:
					ctx, cancel := context.WithCancel(rootCtx)
					ctxManager[ct.PlaylistMD.ID] = cancel

					playlistDirName := filepath.Join(
						ct.Dir,
						yt.SanitizeFilename(ct.PlaylistMD.Title),
					)
					err := os.Mkdir(playlistDirName, 0755)
					if err != nil {
						if errors.Is(err, os.ErrExist) {
							for _, video := range ct.PlaylistMD.Videos {
								dt := &DownloadAllPlaylistTask{
									RootCtx:    rootCtx,
									Ctx:        ctx,
									CtxManager: ctxManager,
									Video:      video,
									Dir:        playlistDirName,
									ItagNo:     ct.ItagNo,
									Chan:       c,
								}
								wp.Add(dt)
							}
						} else {
							c.GUIChannel <- channel.ErrorChannel(err.Error())
						}
					} else {
						for _, video := range ct.PlaylistMD.Videos {
							dt := &DownloadAllPlaylistTask{
								RootCtx:    rootCtx,
								Ctx:        ctx,
								CtxManager: ctxManager,
								Video:      video,
								Dir:        playlistDirName,
								ItagNo:     ct.ItagNo,
								Chan:       c,
							}
							wp.Add(dt)
						}
					}
				}
			case <-rootCtx.Done():
				wp.Close()
				return
			}
		}
	}()
}
