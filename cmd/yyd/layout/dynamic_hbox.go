package customlayout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
)

type YAlignment uint8

const (
	CenterYAlignment YAlignment = iota
	TopYAlignment
	BottomYAlignment
)

type DynamicHBoxLayout struct {
	Gap      float32
	XPadding float32
	YPadding float32
	YAlign   YAlignment
}

func NewDynamicHBoxLayout(gap, xPadding, yPadding float32, yalign YAlignment) *DynamicHBoxLayout {
	return &DynamicHBoxLayout{
		gap,
		xPadding,
		yPadding,
		yalign,
	}
}

func (l *DynamicHBoxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	visibleObjects := 0
	w, h := float32(0), float32(0)
	for _, o := range objects {
		if !o.Visible() || isHorizontalSpacer(o) {
			continue
		}

		visibleObjects++

		var childSize fyne.Size

		if o.Size().Width > 0 { //took the resize value if set
			childSize = o.Size()
		} else {
			childSize = o.MinSize()
		}

		w += childSize.Width
		h = max(h, childSize.Height)
	}

	if visibleObjects > 1 {
		w += float32(visibleObjects-1) * l.Gap
	}

	w += 2 * l.XPadding
	h += 2 * l.YPadding

	return fyne.NewSize(w, h)
}

func (l *DynamicHBoxLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	len := len(objects)

	maxHeightObject := float32(0)
	spacers := 0
	visibleObjects := 0
	total := 2 * (l.XPadding)

	for _, o := range objects {
		if !o.Visible() {
			continue
		}

		if isHorizontalSpacer(o) {
			spacers++
			continue
		}

		visibleObjects++

		if o.Size().Width > 0 {
			total += o.Size().Width
		} else {
			total += o.MinSize().Width
		}

		if o.Size().Height > 0 { //took the maximum object height
			maxHeightObject = fyne.Max(maxHeightObject, o.Size().Height)
		} else {
			maxHeightObject = fyne.Max(maxHeightObject, o.MinSize().Height)
		}
	}

	//adding padding between only visible objects
	if visibleObjects > 1 {
		total += float32((visibleObjects - 1)) * l.Gap
	}

	extra := size.Width - total

	spacerSize := float32(0)
	if spacers > 0 {
		spacerSize = extra / float32(spacers)
	}

	x := float32(l.XPadding)
	for i, o := range objects {
		if !o.Visible() {
			continue
		}

		if isHorizontalSpacer(o) {
			x += spacerSize
			continue
		}

		var childSize fyne.Size

		if o.Size().Width > 0 { //took the resize value
			childSize = o.Size()
		} else {
			childSize = o.MinSize()
		}

		// default y-axis alignement is center
		y := (maxHeightObject - childSize.Height) / 2
		switch l.YAlign {
		case TopYAlignment:
			y = l.YPadding
		case BottomYAlignment:
			y = maxHeightObject - childSize.Height
		}

		o.Move(fyne.NewPos(x, y+l.YPadding))
		o.Resize(childSize)

		if i < len-1 {
			x += l.Gap
		}

		x += childSize.Width
	}
}

func isHorizontalSpacer(obj fyne.CanvasObject) bool {
	spacer, ok := obj.(layout.SpacerObject)
	return ok && spacer.ExpandHorizontal()
}
