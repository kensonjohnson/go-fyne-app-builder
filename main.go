package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/storage"
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
	ui := &gui{window: window, title: binding.NewString()}
	window.SetContent(ui.makeGui())

	// Create a menu
	window.SetMainMenu(ui.makeMenu())
	ui.title.AddListener(binding.NewDataListener(func() {
		name, _ := ui.title.Get()
		window.SetTitle("App Builder: " + name)
	}))

	// Take in possible path from the command line
	flag.Usage = func() {
		fmt.Println("Usage: builder [project directory]")
	}
	flag.Parse()

	if len(flag.Args()) > 0 {
		// A possible path argument was passed in
		dirPath := flag.Arg(0)
		dirPath, err := filepath.Abs(dirPath)
		if err != nil {
			fmt.Println("Error resolving project path", err)
			return
		}

		dirURI := storage.NewFileURI(dirPath)
		dir, err := storage.ListerForURI(dirURI)
		if err != nil {
			fmt.Println("Error opening project", err)
			return
		}

		ui.openProject(dir)
	} else {
		// No args passed in
		a.Lifecycle().SetOnStarted(ui.openProjectDialog)
	}

	// Finally, create the window(s) and run the app.
	// This is blocking until the app returns or errors
	window.ShowAndRun()
}

func (ui *gui) makeMenu() *fyne.MainMenu {
	file := fyne.NewMenu(
		"File",
		fyne.NewMenuItem("Open Project", ui.openProjectDialog),
	)

	return fyne.NewMainMenu(file)
}
