package main

import (
	"fmt"
	"io"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/storage"
)

func createProject(name string, parent fyne.ListableURI) (fyne.ListableURI, error) {
	dir, err := storage.Child(parent, name)
	if err != nil {
		return nil, err
	}

	err = storage.CreateListable(dir)
	if err != nil {
		return nil, err
	}

	mod, err := storage.Child(dir, "go.mod")
	if err != nil {
		return nil, err
	}

	writer, err := storage.Writer(mod)
	if err != nil {
		return nil, err
	}
	defer writer.Close()

	_, err = io.WriteString(writer, fmt.Sprintf(`module %s

go 1.17

require fyne.io/fyne/v2 v2.5.0
`, name))
	if err != nil {
		return nil, err
	}

	json, err := storage.Child(dir, "main.gui.json")
	if err != nil {
		return nil, err
	}

	writer, err = storage.Writer(json)
	if err != nil {
		return nil, err
	}
	defer writer.Close()

	_, err = io.WriteString(writer, fmt.Sprintf(`{
  "Type": "*widget.Label",
  "Struct": {
    "Hidden": false,
    "Text": "Welcome %s!",
    "Alignment": 0,
    "Wrapping": 0,
    "TextStyle": {
      "Bold": false,
      "Italic": false,
      "Monospace": false,
      "Symbol": false,
      "TabWidth": 0,
      "Underline": false
    },
    "Truncation": 0,
    "Importance": 0
  }
}`, name))

	list, _ := storage.ListerForURI(dir)

	return list, err
}

func (g *gui) openProject(dir fyne.ListableURI) {
	name := dir.Name()

	g.title.Set(name)

	// Empty out data before binding new project directory
	g.fileTree.Set(map[string][]string{}, map[string]fyne.URI{})

	addFilesToTree(dir, g.fileTree, binding.DataTreeRootID)
}

func addFilesToTree(dir fyne.ListableURI, tree binding.URITree, root string) {
	items, _ := dir.List()
	for _, uri := range items {
		nodeId := uri.String()
		tree.Append(root, nodeId, uri)

		isDir, err := storage.CanList(uri)
		if err != nil {
			log.Println("Failed to check for listing")
		}

		if isDir {
			child, _ := storage.ListerForURI(uri)
			addFilesToTree(child, tree, nodeId)
		}
	}
}
