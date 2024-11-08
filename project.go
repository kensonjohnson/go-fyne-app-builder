package main

import (
	"fmt"
	"io"

	"fyne.io/fyne/v2"
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

	list, _ := storage.ListerForURI(dir)

	return list, err
}
