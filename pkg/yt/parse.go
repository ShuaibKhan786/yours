package yt

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (
	validYTHost1         = "www.youtube.com"
	validYTHost2         = "youtu.be"
	validYTHost3         = "youtube.com"
	videoIDQuery         = "v"
	videoPlaylistIDQuery = "list"
)

var errInvalidYTLink = errors.New("invalid youtube link")

type ytDetails struct {
	isValidYTLink bool
	isPlaylist    bool
	videoID       string
	playlistID    string
}

func newYTDetails() *ytDetails {
	return &ytDetails{}
}

func (ytd *ytDetails) parsed(link string) error {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return err
	}

	switch parsedURL.Host {
	case validYTHost1, validYTHost2, validYTHost3:
		ytd.isValidYTLink = true
	default:
		return fmt.Errorf("%w: %s", errInvalidYTLink, link)
	}

	// for host = youtu.be ID is path
	if parsedURL.Host == validYTHost2 {
		splitedURL := strings.Split(parsedURL.Path, "/")
		ytd.videoID = splitedURL[1]
	} else {
		// for host = www.youtube.com ID is query 'v'
		if vID := parsedURL.Query().Get(videoIDQuery); vID != "" {
			ytd.videoID = vID
		}
	}

	if vpID := parsedURL.Query().Get(videoPlaylistIDQuery); vpID != "" {
		ytd.isPlaylist = true
		ytd.playlistID = vpID
	}

	return nil
}

func (ytd *ytDetails) validPlaylist() bool {
	return ytd.isPlaylist
}

func (ytd *ytDetails) getVideoID() string {
	return ytd.videoID
}

func (ytd *ytDetails) getPlaylistID() string {
	return ytd.playlistID
}
