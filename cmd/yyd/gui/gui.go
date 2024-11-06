//go:generate fyne bundle -o icons.go ../assets/icons

package gui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/ShuaibKhan786/yours/cmd/yyd/channel"
	"github.com/ShuaibKhan786/yours/cmd/yyd/global"
	"github.com/ShuaibKhan786/yours/cmd/yyd/navigation"
)

type YYD struct {
	App               fyne.App
	Window            fyne.Window
	Preference        fyne.Preferences
	ParentContent     *fyne.Container
	FileSavedLocation string
	Channel           *channel.Channel
	RootCtx           context.Context
	RootCancel        context.CancelFunc
	PlaylistMap       *global.PlaylistMap
}

func NewYYD(rootCtx context.Context, rootCancel context.CancelFunc, c *channel.Channel) *YYD {
	a := app.NewWithID(AppID)
	w := a.NewWindow(AppTitle)
	w.SetMaster()
	w.Resize(fyne.NewSize(MainWindowWidthSize, MainWindowHeightSize))

	playlistMap := &global.PlaylistMap{
		Details: make(map[string]*global.MediaDetails),
	}
	return &YYD{
		App:               a,
		Window:            w,
		Preference:        a.Preferences(),
		Channel:           c,
		RootCtx:           rootCtx,
		RootCancel:        rootCancel,
		PlaylistMap:       playlistMap,
		FileSavedLocation: a.Preferences().String(PreferenceFileSavedLocation),
	}
}

func (y *YYD) setParentContainer(pContainer *fyne.Container) {
	y.ParentContent = container.NewBorder(
		SetYpadding(2*DefaultPaddingSize), SetYpadding(2*DefaultPaddingSize),
		SetXpadding(2*DefaultPaddingSize), SetXpadding(2*DefaultPaddingSize),
		pContainer,
	)
}

func (y *YYD) setContent() {
	y.Window.CenterOnScreen()
	y.Window.SetContent(
		y.ParentContent,
	)
}

func (y *YYD) setTheme() {
	y.App.Settings().SetTheme(&yydTheme{})
}

func (y *YYD) run() {
	y.Window.ShowAndRun()
}

func (y *YYD) setOnClose() {
	y.Window.SetOnClosed(
		func() {
			y.RootCancel()
			close(y.Channel.GUIChannel)
			close(y.Channel.BackendChannel)
		},
	)
}

func (y *YYD) setFileSavedLocation(path string) {
	y.FileSavedLocation = path
	y.Preference.SetString(PreferenceFileSavedLocation, path)
}

func (y *YYD) setIsSetup() {
	y.Preference.SetBool(PreferenceIsSetup, true)
}

func (y *YYD) InitGUI() {
	isSetup := y.Preference.Bool(PreferenceIsSetup)

	if !isSetup {
		sn := navigation.NewStackNavigation()
		sn.InitAllContents([]*fyne.Container{
			y.getStartedPage(sn),
			y.mainPage(),
		})
		y.setParentContainer(sn.MainContentStackNavigation())
	} else {
		y.setParentContainer(y.mainPage())
	}

	y.setTheme()
	y.setContent()
	y.setOnClose()
	y.run()
}
