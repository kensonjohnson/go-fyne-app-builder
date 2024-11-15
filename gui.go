package main

import (
	"app-builder/internal/dialogs"
	"app-builder/internal/editors"
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
	window fyne.Window
	title  binding.String

	fileTree binding.URITree
	content  *container.DocTabs
	openTabs map[fyne.URI]*container.TabItem
}

// Creates a stack with the toolbar on top and logo centered underneath
func (g *gui) makeBanner() fyne.CanvasObject {
	title := canvas.NewText("App Creator", theme.Color(theme.ColorNameForeground))
	title.TextSize = 14
	title.TextStyle = fyne.TextStyle{Bold: true}

	g.title.AddListener(binding.NewDataListener(func() {
		name, _ := g.title.Get()
		if name == "" {
			name = "App Creator"
		}
		title.Text = name
		title.Refresh()
	}))

	home := widget.NewButtonWithIcon("", theme.HomeIcon(), func() {})
	left := container.NewHBox(home, title)

	logo := canvas.NewImageFromResource(resourceLogoPng)
	logo.FillMode = canvas.ImageFillContain

	return container.NewStack(container.NewPadded(left), container.NewPadded(logo))
}

// Creates a layout that has fixed sized left and right columns, and a center
// column that can grow and shrink with the window.
func (g *gui) makeGui() fyne.CanvasObject {

	top := g.makeBanner()

	g.fileTree = binding.NewURITree()
	files := widget.NewTreeWithData(
		g.fileTree,
		func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Filename.jpg")
		},
		func(data binding.DataItem, branch bool, object fyne.CanvasObject) {
			l := object.(*widget.Label)
			u, _ := data.(binding.URI).Get()

			l.SetText(u.Name())
		},
	)
	files.OnSelected = func(uid widget.TreeNodeID) {
		u, err := g.fileTree.GetValue(uid)
		if err != nil {
			dialog.ShowError(err, g.window)
			return
		}

		g.openFile(u)
	}

	left := widget.NewAccordion(
		widget.NewAccordionItem("Files", files),
		widget.NewAccordionItem("Sceens", widget.NewLabel("TODO Screens")),
	)
	left.Open(0)
	left.MultiOpen = true

	right := widget.NewRichTextFromMarkdown("## Settings")

	name, _ := g.title.Get()
	window := container.NewInnerWindow(
		name,
		widget.NewLabel("App preview here"),
	)
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

	g.content = container.NewDocTabs(
		container.NewTabItem("Preview", content),
	)
	g.content.CloseIntercept = func(ti *container.TabItem) {
		var uri fyne.URI

		for index, child := range g.openTabs {
			if child == ti {
				uri = index
			}
		}

		if uri != nil {
			delete(g.openTabs, uri)
		}

		g.content.Remove(ti)
	}

	objects := []fyne.CanvasObject{g.content, top, left, right}

	dividers := [3]fyne.CanvasObject{
		widget.NewSeparator(),
		widget.NewSeparator(),
		widget.NewSeparator(),
	}

	for i := 0; i < len(dividers); i++ {
		objects = append(objects, dividers[i])
	}

	return container.New(newAppBuilderLayout(top, left, right, g.content, dividers), objects...)
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

func (g *gui) openFile(u fyne.URI) {
	listable, err := storage.CanList(u)
	if listable || err != nil {
		// TODO: should we unselect this item
		return
	}

	if item, ok := g.openTabs[u]; ok {
		g.content.Select(item)
		return
	}

	edit := editors.ForURI(u)

	item := container.NewTabItem(u.Name(), edit)

	if g.openTabs == nil {
		g.openTabs = make(map[fyne.URI]*container.TabItem)
	}
	g.openTabs[u] = item

	g.content.Append(item)
	g.content.Select(item)

}
