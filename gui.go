package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type gui struct {
	window    fyne.Window
	title     binding.String
	directory *widget.Label
}

// Creates a stack with the toolbar on top and logo centered underneath
func makeBanner() fyne.CanvasObject {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.HomeIcon(), func() {}),
	)

	logo := canvas.NewImageFromResource(resourceLogoPng)
	logo.FillMode = canvas.ImageFillContain

	return container.NewStack(toolbar, container.NewPadded(logo))
}

func (g *gui) makeGui() fyne.CanvasObject {

	left := widget.NewLabel("left")
	right := widget.NewLabel("right")

	directory := widget.NewLabelWithData(g.title)
	content := container.NewStack(canvas.NewRectangle(color.Gray{Y: 0xee}), directory)

	top := makeBanner()

	objects := []fyne.CanvasObject{content, top, left, right}

	dividers := [3]fyne.CanvasObject{
		widget.NewSeparator(),
		widget.NewSeparator(),
		widget.NewSeparator(),
	}

	for i := 0; i < len(dividers); i++ {
		objects = append(objects, dividers[i])
	}

	return container.New(newAppBuilderLayout(top, left, right, content, dividers), objects...)
}

func (g *gui) openProjectDialog() {
	dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, g.window)
			return
		}

		if dir == nil {
			return
		}

		g.openProject(dir)
	}, g.window)
}

func (g *gui) openProject(dir fyne.ListableURI) {
	name := dir.Name()

	g.title.Set(name)

}
