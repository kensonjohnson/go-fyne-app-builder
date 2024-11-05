package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	window := a.NewWindow("App Builder")

	window.SetContent(widget.NewLabel("Hello Builder!"))
	window.ShowAndRun()
}
