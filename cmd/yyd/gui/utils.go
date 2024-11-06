package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func SetXpadding(size float32) *canvas.Rectangle {
	padding := canvas.NewRectangle(color.Transparent)
	padding.SetMinSize(fyne.NewSize(size, 0))
	return padding
}

func SetYpadding(size float32) *canvas.Rectangle {
	padding := canvas.NewRectangle(color.Transparent)
	padding.SetMinSize(fyne.NewSize(0, size))
	return padding
}
