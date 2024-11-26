package editors

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func setPreviewTheme(c *container.ThemeOverride, th fyne.Theme) {
	c.Theme = th
	c.Refresh()
}
