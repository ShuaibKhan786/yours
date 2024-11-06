package channel

import (
	"github.com/ShuaibKhan786/yours/cmd/yyd/global"
	"github.com/ShuaibKhan786/yours/pkg/yt"
	"github.com/kkdai/youtube/v2"
)

type Channel struct {
	GUIChannel     chan interface{}
	BackendChannel chan interface{}
}

func InitChannel() *Channel {
	return &Channel{
		GUIChannel:     make(chan interface{}, 0),
		BackendChannel: make(chan interface{}, 0),
	}
}

func (c *Channel) Close() {
	close(c.GUIChannel)
	close(c.BackendChannel)
}

// backend channel types
type LinkChannel string

type DownloadChannel struct {
	IsPlaylist bool
	ItagNo     int
	Video      *youtube.Video
	Dir        string
}

type CancelDownload string

type CancelPlaylistImageDownload struct{}

type DownloadPlaylistImages struct {
	PlaylistMap *global.PlaylistMap
	MD          *yt.PlaylistMetadata
}

type DownloadAllPlaylist struct {
	PlaylistMD *yt.PlaylistMetadata
	ItagNo     int
	Dir        string
}

// gui channel types
type DownloadDone string

type DownloadCancelDone string

type ErrorChannel string

type FetchMDStarted struct{}

type PlaylistDownloadAllCancelDone struct{}
type PlaylistDownloadAllError struct{}
type PlaylistDownloadAllDone struct{}
