package yt

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/kkdai/youtube/v2"
)

const (
	videoExt   = "mp4"
	audioExt   = "m4a"
	bufferSize = 32 * 1024
)

type Metadata interface{}

type PlaylistMetadata struct {
	ID          string
	Title       string
	Description string
	Author      string
	Videos      []*youtube.Video
}

func GetYTMetadata(ctx context.Context, link string) (Metadata, error) {
	ytd := newYTDetails()
	err := ytd.parsed(link)
	if err != nil {
		return nil, err
	}

	if ytd.validPlaylist() {
		return getPlaylistMetadata(ctx, link)
	} else {
		return getVideoMetadata(ctx, link)
	}
}

func getVideoMetadata(ctx context.Context, videoID string) (*youtube.Video, error) {
	client := &youtube.Client{}

	video, err := client.GetVideoContext(ctx, videoID)
	if err != nil {
		return nil, err
	}

	video.Formats = sanitizedMP4FormatsOnly(video.Formats)
	return video, nil
}

func getPlaylistMetadata(ctx context.Context, playlistID string) (*PlaylistMetadata, error) {
	client := &youtube.Client{}
	playlist, err := client.GetPlaylistContext(ctx, playlistID)
	if err != nil {
		return nil, err
	}

	md := &PlaylistMetadata{
		ID:          playlist.ID,
		Title:       playlist.Title,
		Description: playlist.Description,
		Author:      playlist.Author,
	}

	videos := make([]*youtube.Video, 0, len(playlist.Videos))

	for _, entry := range playlist.Videos {
		video, err := client.VideoFromPlaylistEntryContext(ctx, entry)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil, err
			}
			continue
		}

		video.Formats = sanitizedMP4FormatsOnly(video.Formats)

		videos = append(videos, video)
	}

	md.Videos = videos

	return md, err
}

func Download(ctx context.Context, video *youtube.Video, dir string, itagNo int) error {
	client := &youtube.Client{}

	switch itagNo {
	case 140: //audio only
		outputAudioFilename := fmt.Sprintf("%s.%s", SanitizeFilename(video.Title), audioExt)
		outputFullFilepath := filepath.Join(dir, outputAudioFilename)

		formats := video.Formats.Itag(itagNo)

		reader, _, err := client.GetStreamContext(ctx, video, &formats[0])
		if err != nil {
			return err
		}
		defer reader.Close()

		writer, err := os.Create(outputFullFilepath)
		if err != nil {
			return err
		}

		err = CopyContext(ctx, writer, reader)
		if err != nil {
			writer.Close()
			os.Remove(outputFullFilepath)
			return err
		}

		writer.Close()
		return nil
	case 18: //audio + video
		formats := video.Formats.Itag(itagNo)

		outputVideoFilename := fmt.Sprintf("%s(%s).%s", SanitizeFilename(video.Title), formats[0].QualityLabel, videoExt)
		outputFullFilepath := filepath.Join(dir, outputVideoFilename)

		reader, _, err := client.GetStreamContext(ctx, video, &formats[0])
		if err != nil {
			return err
		}
		defer reader.Close()

		writer, err := os.Create(outputFullFilepath)
		if err != nil {
			return err
		}

		err = CopyContext(ctx, writer, reader)
		if err != nil {
			writer.Close()
			os.Remove(outputFullFilepath)
			return err
		}

		writer.Close()
		return nil
	default: //audio and video must merge
		if !IsFFmpegInstalled() {
			return nil
		}

		audioFormat := video.Formats.Itag(140)
		videoFormat := video.Formats.Itag(itagNo)

		outputVideoFilename := fmt.Sprintf("%s(%s).%s", SanitizeFilename(video.Title), videoFormat[0].QualityLabel, videoExt)
		outputFullFilepath := filepath.Join(dir, outputVideoFilename)

		tempVideoFilename := fmt.Sprintf("%s(%s).%s", video.ID, videoFormat[0].QualityLabel, videoExt)
		tempVideoFullFilepath := filepath.Join(dir, tempVideoFilename)

		tempAudioFilename := fmt.Sprintf("%s.%s", video.ID, audioExt)
		tempAudioFullFilepath := filepath.Join(dir, tempAudioFilename)

		// download the audio temporarily
		audioReader, _, err := client.GetStreamContext(ctx, video, &audioFormat[0])
		if err != nil {
			return err
		}
		defer audioReader.Close()

		audioWriter, err := os.Create(tempAudioFullFilepath)
		if err != nil {
			return err
		}
		defer func() {
			audioWriter.Close()
			os.Remove(tempAudioFullFilepath)
		}()

		err = CopyContext(ctx, audioWriter, audioReader)
		if err != nil {
			return err
		}

		// download the video temporarily
		videoReader, _, err := client.GetStreamContext(ctx, video, &videoFormat[0])
		if err != nil {
			return err
		}
		defer videoReader.Close()

		videoWriter, err := os.Create(tempVideoFullFilepath)
		if err != nil {
			return err
		}
		defer func() {
			videoWriter.Close()
			os.Remove(tempVideoFullFilepath)
		}()

		err = CopyContext(ctx, videoWriter, videoReader)
		if err != nil {
			return err
		}

		//merge the audio and video in one video file
		ffmpegCmd := exec.CommandContext(
			ctx,
			"ffmpeg",
			"-i", tempVideoFullFilepath,
			"-i", tempAudioFullFilepath,
			"-c:v", "copy",
			"-c:a", "copy",
			outputFullFilepath,
			"-loglevel", "warning",
		)

		ffmpegCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

		if err := ffmpegCmd.Start(); err != nil {
			os.Remove(outputFullFilepath)
			return err
		}
		if err := ffmpegCmd.Wait(); err != nil {
			return err
		}
	}
	return nil
}

func CopyContext(ctx context.Context, dst io.Writer, src io.Reader) error {
	buffer := make([]byte, bufferSize)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			n, err := src.Read(buffer)
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}

			_, err = dst.Write(buffer[:n])
			if err != nil {
				return err
			}
		}
	}
}

func IsFFmpegInstalled() bool {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return false
	}
	return true
}
