package main

import (
	"fyne.io/fyne/v2"
)

const sideWidth = 220

type appBuilderLayout struct {
	top, left, right, content fyne.CanvasObject
}

func newAppBuilderLayout(top, left, right, content fyne.CanvasObject) fyne.Layout {
	return &appBuilderLayout{top: top, left: left, right: right, content: content}
}

func (l *appBuilderLayout) Layout(_ []fyne.CanvasObject, size fyne.Size) {
	topHeight := l.top.MinSize().Height
	l.top.Resize(fyne.NewSize(size.Width, topHeight))

	l.left.Move(fyne.NewPos(0, topHeight))
	l.left.Resize(fyne.NewSize(sideWidth, size.Height-topHeight))

	l.right.Move(fyne.NewPos(size.Width-sideWidth, topHeight))
	l.right.Resize(fyne.NewSize(sideWidth, size.Height-topHeight))

	l.content.Move(fyne.NewPos(sideWidth, topHeight))
	l.content.Resize(fyne.NewSize(size.Width-sideWidth*2, size.Height-topHeight))
}

func (l *appBuilderLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	borders := fyne.NewSize(
		sideWidth*2,
		l.top.MinSize().Height,
	)

	return borders.AddWidthHeight(100, 100)
}
