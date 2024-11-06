package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type yydTheme struct{}

func (m *yydTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if variant == theme.VariantLight {
		switch name {
		case theme.ColorNameForeground:
			return color.RGBA{R: 28, G: 27, B: 29, A: 255}
		case theme.ColorNamePrimary,
			theme.ColorNameFocus:
			return color.RGBA{R: 220, G: 218, B: 245, A: 255}
		case theme.ColorNameSelection:
			return color.RGBA{R: 85, G: 63, B: 93, A: 255}
		default:
			return theme.DefaultTheme().Color(name, variant)
		}
	} else {
		switch name {
		case theme.ColorNameForeground:
			return color.RGBA{R: 230, G: 225, B: 227, A: 255}
		case theme.ColorNamePrimary,
			theme.ColorNameFocus:
			return color.RGBA{R: 69, G: 69, B: 90, A: 255}
		case theme.ColorNameSelection:
			return color.RGBA{R: 85, G: 63, B: 93, A: 255}
		default:
			return theme.DefaultTheme().Color(name, variant)
		}
	}
}

func (m *yydTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m *yydTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func (m *yydTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}
