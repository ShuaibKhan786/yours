package navigation

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type StackNavigation struct {
	mainContainer *fyne.Container
	contents      []*fyne.Container
	currentPos    int
}

func NewStackNavigation() *StackNavigation {
	return &StackNavigation{
		mainContainer: container.NewStack(),
		contents:      []*fyne.Container{},
		currentPos:    0,
	}
}

func (sv *StackNavigation) InitAllContents(contents []*fyne.Container) {
	sv.contents = append(sv.contents, contents...)
	sv.mainContainer.Add(sv.contents[0])
}

func (sv *StackNavigation) MainContentStackNavigation() *fyne.Container {
	return sv.mainContainer
}

func (sv *StackNavigation) Push() {
	if sv.currentPos < len(sv.contents)-1 {
		sv.mainContainer.RemoveAll()
		sv.currentPos++
		sv.mainContainer.Add(sv.contents[sv.currentPos])
		sv.mainContainer.Refresh()
	}
}

func (sv *StackNavigation) Pop() {
	if sv.currentPos > 0 {
		sv.mainContainer.RemoveAll()
		sv.currentPos--
		sv.mainContainer.Add(sv.contents[sv.currentPos])
		sv.mainContainer.Refresh()
	}
}
