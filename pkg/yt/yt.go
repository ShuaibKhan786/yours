package yt

import (
	"context"

	"github.com/kkdai/youtube/v2"
)

type Metadata interface{}

type PlaylistMD struct {
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

	video.Formats = sanitizedMP3MP4FormatsOnly(video.Formats)
	return video, nil
}

func getPlaylistMetadata(ctx context.Context, playlistID string) (*PlaylistMD, error) {
	client := &youtube.Client{}
	playlist, err := client.GetPlaylistContext(ctx, playlistID)
	if err != nil {
		return nil, err
	}

	md := &PlaylistMD{
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

		video.Formats = sanitizedMP3MP4FormatsOnly(video.Formats)

		videos = append(videos, video)
	}

	md.Videos = videos

	return md, err
}
