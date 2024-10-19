package yt_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/ShuaibKhan786/yours/pkg/yt"
	"github.com/kkdai/youtube/v2"
)

func TestYT(t *testing.T) {
	expectedTitle := "Diamond Eyes - Worship | DnB | NCS - Copyright Free Music"
	expectedAuthor := "NoCopyrightSounds"
	t.Run("video metadata test", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		mdInterface, err := yt.GetYTMetadata(ctx, "https://youtu.be/gH9L98XWmiQ?si=YU6fR_4s5ohIVHAh")
		if err != nil {
			t.Errorf("Expected no error but Got this error: %v", err.Error())
		}

		switch md := mdInterface.(type) {
		case *youtube.Video:
			assertEqual(t, expectedTitle, md.Title)
			assertEqual(t, expectedAuthor, md.Author)

			if len(md.Formats) > 9 {
				t.Errorf("Expected to be filter only the suppported mimetypes")
			}

			assertFormats(t, md.Formats)
		default:
			t.Errorf("Expected the type to be *youtube.Video")
		}
	})

	t.Run("playlist metadata testing", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		mdInterface, err := yt.GetYTMetadata(ctx, "https://www.youtube.com/watch?v=GfCqMv--ncA&list=PLbaPlkHgQC08gpJJ1zkucJNHFATs1FwCu")
		if err != nil {
			t.Errorf("Expected no error but Got this error: %v", err.Error())
		}

		switch md := mdInterface.(type) {
		case *yt.PlaylistMD:
			assertEqual(t, "track I enjoy", md.Title)
			assertEqual(t, "Shuaib Khan", md.Author)
			assertEqual(t, expectedTitle, md.Videos[5].Title)
			assertEqual(t, expectedAuthor, md.Videos[5].Author)
			if len(md.Videos) != 6 {
				t.Errorf("Expected number of videos in a playlist is 6 but got %d", len(md.Videos))
			}
			assertFormats(t, md.Videos[5].Formats)
		default:
			t.Errorf("Expected the type to be *youtube.Playlist")
		}
	})
}

func assertEqual(t *testing.T, expected, got string) {
	if expected != got {
		t.Errorf("Expected: %s\n but Got: %s", expected, got)
	}
}

func assertFormats(t *testing.T, formats youtube.FormatList) {
	isExpectedItagAvilable := false
	for _, format := range formats {
		if strings.Contains(format.MimeType, "video/webm") {
			t.Errorf("Expected no webm video format after sanitized")
		}

		if format.ItagNo == 140 { // codec: ACC, 128kb bitrate
			isExpectedItagAvilable = true
		}
	}

	if !isExpectedItagAvilable {
		t.Errorf("Expected ACC codec of 128kb bitrate")
	}

}
