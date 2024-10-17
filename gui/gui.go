package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type YTAD struct {
	App           fyne.App
	Window        fyne.Window
	Preference    fyne.Preferences
	ParentContent *fyne.Container
}

func NewYTAD() *YTAD {
	a := app.NewWithID(AppID)
	w := a.NewWindow(AppTitle)
	w.SetMaster()
	w.Resize(fyne.NewSize(600, 350))

	return &YTAD{
		App:        a,
		Window:     w,
		Preference: a.Preferences(),
	}
}

func (y *YTAD) setParentContainer(getStartedCont, mainCont *fyne.Container) {
	isSetup := y.Preference.Bool(PreferenceIsSetup)

	if isSetup {
		y.ParentContent = mainCont
	} else {
		y.ParentContent = getStartedCont
	}
}

func (y *YTAD) setContent() {
	y.Window.SetContent(
		container.NewBorder(
			SetYpadding(20), SetYpadding(20),
			SetXpadding(20), SetXpadding(20),
			y.ParentContent,
		),
	)
}

func (y *YTAD) run() {
	y.Window.ShowAndRun()
}

func (y *YTAD) InitGUI() {

	entryWidget := widget.NewEntry()
	entryWidget.PlaceHolder = "your youtube link here"
	getStartedCont := container.NewStack(
		
	)
	mainCont := container.NewCenter(widget.NewLabel("Mian container"))
	y.setParentContainer(getStartedCont, mainCont)
	y.setContent()
	y.run()
}
