package yt

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kkdai/youtube/v2"
)

const (
	videoExt = "mp4"
	audioExt = "m4a"
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
			return nil, err
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
		defer writer.Close()

		_, err = io.Copy(writer, reader)
		if err != nil {
			// defer os.Remove(outputFullFilepath)
			return err
		}

		return nil
	case 18: //audio + video
		outputVideoFilename := fmt.Sprintf("%s.%s", SanitizeFilename(video.Title), videoExt)
		outputFullFilepath := filepath.Join(dir, outputVideoFilename)

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
		defer writer.Close()

		_, err = io.Copy(writer, reader)
		if err != nil {
			return err
		}

		return nil
	default: //audio and video must merge
		if !isFFmpegInstalled() {
			return nil
		}

		outputVideoFilename := fmt.Sprintf("%s.%s", SanitizeFilename(video.Title), videoExt)
		outputFullFilepath := filepath.Join(dir, outputVideoFilename)

		tempVideoFilename := fmt.Sprintf("%s.%s", video.ID, videoExt)
		tempVideoFullFilepath := filepath.Join(dir, tempVideoFilename)

		tempAudioFilename := fmt.Sprintf("%s.%s", video.ID, audioExt)
		tempAudioFullFilepath := filepath.Join(dir, tempAudioFilename)

		// download the audio temporarily
		audioFormat := video.Formats.Itag(140)
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

		_, err = io.Copy(audioWriter, audioReader)
		if err != nil {
			return err
		}

		// download the video temporarily
		videoFormat := video.Formats.Itag(itagNo)
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

		_, err = io.Copy(videoWriter, videoReader)
		if err != nil {
			return err
		}

		//merge the audio and video in one video file
		ffmpegCmd := exec.Command("ffmpeg",
			"-i", tempVideoFullFilepath,
			"-i", tempAudioFullFilepath,
			"-c:v", "copy",
			"-c:a", "copy",
			outputFullFilepath,
			"-loglevel", "warning",
		)

		if err := ffmpegCmd.Start(); err != nil {
			return err
		}
		if err := ffmpegCmd.Wait(); err != nil {
			return err
		}
	}
	return nil
}

func isFFmpegInstalled() bool {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return false
	}
	return true
}
