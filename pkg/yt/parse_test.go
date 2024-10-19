package yt

import (
	"errors"
	"testing"
)

func TestParse(t *testing.T) {
	t.Run("with invalid yt link", func(t *testing.T) {
		ytd := newYTDetails()
		err := ytd.parsed("https://go.dev/play/")

		if err == nil {
			t.Errorf("Expected an error when parsed an invalid yt link")
		}

		if !errors.Is(err, errInvalidYTLink) {
			t.Errorf("Expected error: %s but got: %s", errInvalidYTLink.Error(), err.Error())
		}
	})

	t.Run("with valid yt link but no playlist", func(t *testing.T) {
		ytd := newYTDetails()
		assertDifferentYTLinks(t, ytd, "https://youtu.be/XoiOOiuH8iI?si=bmgpImeUlHZDB4rs", "XoiOOiuH8iI")
		assertDifferentYTLinks(t, ytd, "https://www.youtube.com/watch?v=XoiOOiuH8iI", "XoiOOiuH8iI")
	})

	t.Run("with valid yt link with playlist", func(t *testing.T) {
		ytd := newYTDetails()

		if err := ytd.parsed("https://youtube.com/playlist?list=PL8dPuuaLjXtNlUrzyH5r6jN9ulIgZBpdo&si=IS-ZNCrlmWyM9vSw"); err != nil {
			t.Errorf("Expected no error on this test: %s", err.Error())
		}

		expected := "PL8dPuuaLjXtNlUrzyH5r6jN9ulIgZBpdo"
		if !ytd.validPlaylist() && ytd.getPlaylistID() != expected {
			t.Errorf("Expected playlistID: %s but Got: %s", expected, ytd.getPlaylistID())
		}

	})
}

func assertDifferentYTLinks(t *testing.T, ytd *ytDetails, link, expected string) {
	err := ytd.parsed(link)

	if err != nil {
		t.Errorf("Expected no error on valid yt link")
	}

	if ytd.getVideoID() != expected {
		t.Errorf("Expected videoID: %s but Got: %s", expected, ytd.videoID)
	}
}
