package editors

import (
	"errors"
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

var extentions = map[string]func(fyne.URI) (fyne.CanvasObject, fyne.CanvasObject, error){
	".go":       makeGo,
	".gui.json": makeGUI,
	".png":      makeImg,
	".txt":      makeTxt,
	".md":       makeTxt,
}

var mimes = map[string]func(fyne.URI) (fyne.CanvasObject, fyne.CanvasObject, error){
	"text/plain": makeTxt,
}

func ForURI(u fyne.URI) (fyne.CanvasObject, fyne.CanvasObject, error) {
	name := strings.ToLower(u.Name())
	var matched func(fyne.URI) (fyne.CanvasObject, fyne.CanvasObject, error)
	for ext, edit := range extentions {
		pos := strings.LastIndex(name, ext)
		if pos == -1 || pos != len(name)-len(ext) {
			continue
		}

		matched = edit
		break
	}

	if matched == nil {
		edit, ok := mimes[u.MimeType()]
		if !ok {
			return nil, nil, errors.New("unable to find editor for file: " + u.Name() + ", mime:" + u.MimeType())
		}

		return edit(u)
	}

	return matched(u)
}

func makeGo(u fyne.URI) (fyne.CanvasObject, fyne.CanvasObject, error) {
	// TODO code editor
	code, _, err := makeTxt(u)
	if err != nil {
		return nil, nil, err
	}
	code.(*widget.Entry).TextStyle = fyne.TextStyle{Monospace: true}

	return code, nil, nil
}

func makeImg(u fyne.URI) (fyne.CanvasObject, fyne.CanvasObject, error) {
	img := canvas.NewImageFromURI(u)
	img.FillMode = canvas.ImageFillContain

	return img, nil, nil
}

func makeTxt(u fyne.URI) (fyne.CanvasObject, fyne.CanvasObject, error) {
	code := widget.NewEntry()

	r, err := storage.Reader(u)
	if err != nil {
		return nil, nil, err
	}

	defer r.Close()

	data, _ := io.ReadAll(r)
	code.SetText(string(data))

	return code, nil, nil
}
