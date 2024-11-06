package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	// We start by creating a Fyne App
	a := app.New()

	// We built a custom theme in theme.go
	a.Settings().SetTheme(newAppBuilderTheme())

	// Create a window and give it a title
	window := a.NewWindow("App Builder")

	// Request the window be resized. This is only a request, and the OS
	// may or may not enforce the change.
	window.Resize(fyne.NewSize(1024, 768))

	// We built a helper function to create a basic UI in gui.go
	window.SetContent(makeGui())

	// Finally, create the window(s) and run the app.
	// This is blocking until the app returns or errors
	window.ShowAndRun()
}
