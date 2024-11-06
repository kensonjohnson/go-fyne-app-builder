//go:generate fyne bundle -o bundled.go assets

package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type appBuilderTheme struct {
	fyne.Theme // Inherit properties from the default theme
}

func newAppBuilderTheme() fyne.Theme {
	return &appBuilderTheme{theme.DefaultTheme()}
}

func (t *appBuilderTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return t.Theme.Color(name, theme.VariantLight)
}

func (t appBuilderTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return 12
	}

	return t.Theme.Size(name)
}
