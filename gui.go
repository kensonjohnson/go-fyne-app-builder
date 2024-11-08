package main

import (
	"app-builder/internal/dialogs"
	"errors"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
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

// Creates a layout that has fixed sized left and right columns, and a center
// column that can grow and shrink with the window.
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

// Adds options to the built in file menu for the native OS
func (ui *gui) makeMenu() *fyne.MainMenu {
	file := fyne.NewMenu(
		"File",
		fyne.NewMenuItem("Open Project", ui.openProjectDialog),
	)

	return fyne.NewMainMenu(file)
}

// Creates a new dialog window with a file picker
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

func (g *gui) showCreate(w fyne.Window) {
	var wizard *dialogs.Wizard
	intro := widget.NewLabel(`Here you can create a new project!

Or open an existing on that you created earlier`)

	open := widget.NewButton("Open Project", func() {
		wizard.Hide()
		g.openProjectDialog()
	})

	create := widget.NewButton("Create Project", func() {
		wizard.Push("Project Details", g.makeCreateDetail(wizard))
	})
	create.Importance = widget.HighImportance

	buttons := container.NewGridWithColumns(2, open, create)
	home := container.NewVBox(intro, buttons)

	wizard = dialogs.NewWizard("Create Project", home)
	wizard.Show(w)

	wizard.Resize(home.MinSize().AddWidthHeight(40, 80))
}

func (g *gui) makeCreateDetail(wizard *dialogs.Wizard) fyne.CanvasObject {
	homeDir, _ := os.UserHomeDir()
	parent := storage.NewFileURI(homeDir)
	chosen, _ := storage.ListerForURI(parent)

	name := widget.NewEntry()
	name.Validator = func(s string) error {
		if s == "" {
			return errors.New("Project name is required")
		}

		return nil
	}

	var dir *widget.Button
	dir = widget.NewButton(chosen.Name(), func() {
		d := dialog.NewFolderOpen(func(lu fyne.ListableURI, err error) {
			if err != nil || lu == nil {
				return
			}

			chosen = lu
			dir.SetText(lu.Name())
		}, g.window)

		d.SetLocation(chosen)
		d.Show()
	})

	form := widget.NewForm(
		widget.NewFormItem("Name", name),
		widget.NewFormItem("Parent Directory", dir),
	)

	form.OnSubmit = func() {
		project, err := createProject(name.Text, chosen)
		if err != nil {
			dialog.ShowError(err, g.window)
			return
		}

		wizard.Hide()
		g.openProject(project)
	}

	return form
}
