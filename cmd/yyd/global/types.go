package global

import (
	"image"
	"sync"
)

type MediaDetails struct {
	Thumbnail image.Image
	Formats   []string
}

type PlaylistMap struct {
	Details map[string]*MediaDetails
	Mu      sync.RWMutex
}

