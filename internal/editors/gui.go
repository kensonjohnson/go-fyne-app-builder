package editors

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/fyne-io/defyne/pkg/gui"
)

func makeGUI(u fyne.URI) (fyne.CanvasObject, error) {
	r, err := storage.Reader(u)
	if err != nil {
		return nil, err
	}

	obj, _, err := gui.DecodeObject(r)
	if err != nil {
		return nil, err
	}

	// TODO: Get project title from project type
	name := "Preview"
	window := container.NewInnerWindow(name, obj)
	window.CloseIntercept = func() {}

	picker := widget.NewSelect([]string{"Desktop App", "iPhone 15 Max"}, func(s string) {})
	picker.Selected = "Desktop App"

	preview := container.NewBorder(
		container.NewHBox(picker),
		nil, nil, nil,
		container.NewCenter(window),
	)

	content := container.NewStack(
		canvas.NewRectangle(color.Gray{Y: 0xee}),
		container.NewPadded(preview),
	)

	return content, nil
}
