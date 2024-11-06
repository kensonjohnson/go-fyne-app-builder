package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func makeBanner() fyne.CanvasObject {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.HomeIcon(), func() {}),
	)

	logo := canvas.NewImageFromResource(resourceLogoPng)
	logo.FillMode = canvas.ImageFillContain

	return container.NewStack(toolbar, container.NewPadded(logo))
}

func makeGui() fyne.CanvasObject {

	left := widget.NewLabel("left")
	right := widget.NewLabel("right")

	content := canvas.NewRectangle(color.Gray{Y: 0xee})

	top := makeBanner()

	// return container.NewBorder(makeBanner(), nil, left, right, content)
	objects := []fyne.CanvasObject{content, top, left, right}
	return container.New(newAppBuilderLayout(top, left, right, content), objects...)
}
