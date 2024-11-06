package customlayout_test

import (
	"reflect"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	customlayout "github.com/ShuaibKhan786/yours/cmd/ytad/layout"
)

func TestDynamicHBoxLayout(t *testing.T) {
	// Initialize a Fyne test app
	test.NewApp()

	// Create the layout and a button
	cont := customlayout.NewDynamicHBoxLayout(100, 2, 2, customlayout.BottomYAlignment)
	btn := widget.NewButton("Click Me", func() {})
	btn2 := widget.NewButton("Click Me", func() {})

	// Ensure the button is initialized with a test window
	window := test.NewWindow(btn)
	defer window.Close() // Clean up

	window2 := test.NewWindow(btn2)
	defer window2.Close() // Clean up

	// Get the minimum size of the layout with the button
	got := cont.MinSize([]fyne.CanvasObject{btn, btn2})
	expected := fyne.NewSize(
		(2*btn.MinSize().Width+4)+100,
		btn.MinSize().Height+4,
	)

	// Compare expected and got sizes
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected: %v but got: %v", expected, got)
	}
}
