package ui

import (
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// HomeUI .
type HomeUI struct {
	UI tview.Primitive
}

// NewHomeUI .
func NewHomeUI() *HomeUI {

	text := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorDefault).
		SetText(strings.Trim(`
Welcome to Beholder!

Just start typing to search for
an item, monster, spell, or
class or race feature.

Press F1 or ? for detailed help anywhere.
`, " \n"))

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow)
	flex.SetBorderPadding(2, 1, 2, 2)
	flex.AddItem(text, 0, 1, false)

	return &HomeUI{
		UI: flex,
	}
}
