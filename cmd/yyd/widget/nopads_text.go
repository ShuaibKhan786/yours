package customwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TextWithoutPadding struct {
	widget.BaseWidget
	text      string
	alignment fyne.TextAlign
	style     fyne.TextStyle
}

func NewTextWithoutPadding(text string) *TextWithoutPadding {
	t := &TextWithoutPadding{
		text:      text,
		alignment: fyne.TextAlignLeading, // Default alignment
		style:     fyne.TextStyle{},      // Default text style
	}
	t.ExtendBaseWidget(t)
	return t
}

func (t *TextWithoutPadding) SetText(text string) {
	t.text = text
	t.Refresh()
}

func (t *TextWithoutPadding) SetTextStyle(style fyne.TextStyle) {
	t.style = style
	t.Refresh()
}

func (t *TextWithoutPadding) SetAlignment(alignment fyne.TextAlign) {
	t.alignment = alignment
	t.Refresh()
}

func (t *TextWithoutPadding) CreateRenderer() fyne.WidgetRenderer {
	txt := canvas.NewText(t.text, theme.Color(theme.ColorNameForeground)) // Always use the default foreground color
	txt.TextSize = theme.TextSize()
	txt.Alignment = t.alignment
	txt.TextStyle = t.style

	return &textRenderer{
		text: txt,
	}
}

type textRenderer struct {
	text *canvas.Text
}

func (r *textRenderer) Layout(size fyne.Size) {
	r.text.Color = theme.Color(theme.ColorNameForeground)
	r.text.Resize(size)
}

func (r *textRenderer) MinSize() fyne.Size {
	return r.text.MinSize()
}

func (r *textRenderer) Refresh() {
	r.text.Refresh()
}

func (r *textRenderer) Destroy() {}

func (r *textRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.text}
}
