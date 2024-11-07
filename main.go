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
	ui := &gui{window: window}
	window.SetContent(ui.makeGui())

	// Create a menu
	window.SetMainMenu(ui.makeMenu())

	// Finally, create the window(s) and run the app.
	// This is blocking until the app returns or errors
	window.ShowAndRun()
}

func (ui *gui) makeMenu() *fyne.MainMenu {
	file := fyne.NewMenu(
		"File",
		fyne.NewMenuItem("Open Project", ui.openProject),
	)

	return fyne.NewMainMenu(file)
}
