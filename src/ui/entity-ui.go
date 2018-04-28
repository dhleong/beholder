package ui

import (
	"github.com/rivo/tview"

	beholder "github.com/dhleong/beholder/src"
)

// EntityUI .
type EntityUI struct {
	UI tview.Primitive

	text *tview.TextView
}

// NewEntityUI .
func NewEntityUI() *EntityUI {
	text := tview.NewTextView()
	text.SetBorderPadding(1, 1, 1, 1)
	return &EntityUI{
		UI:   text,
		text: text,
	}
}

// Set the current Entity to be displayed
func (e *EntityUI) Set(entity beholder.Entity) {
	if entity == nil {
		e.text.SetText("")
		return
	}
	e.text.SetText(entity.GetName())
}
