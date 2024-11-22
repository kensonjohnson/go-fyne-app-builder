package editors

import (
	"bytes"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/fyne-io/defyne/pkg/gui"
	"github.com/stretchr/testify/assert"
)

const labelJSON = `{
  "Type": "*widget.Label",
  "Struct": {
    "Hidden": false,
    "Text": "Testing",
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
}
`

func TestDecode(t *testing.T) {
	test.NewApp()
	obj, _, err := gui.DecodeObject(strings.NewReader(labelJSON))

	assert.Nil(t, err)
	assert.NotNil(t, obj)

	l, ok := obj.(*widget.Label)
	assert.True(t, ok)
	assert.Equal(t, "Testing", l.Text)

	test.AssertObjectRendersToImage(t, "label.png", l)
	test.AssertObjectRendersToMarkup(t, "label.xml", l)

}

func TestEncode(t *testing.T) {
	test.NewApp()
	l := widget.NewLabel("Testing")
	w := bytes.NewBuffer(nil)

	meta := make(map[fyne.CanvasObject]map[string]string)
	gui.EncodeObject(l, meta, w)
	json := w.String()
	assert.NotEmpty(t, json)
	assert.Equal(t, labelJSON, json)
}
