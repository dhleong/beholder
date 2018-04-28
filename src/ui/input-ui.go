package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	beholder "github.com/dhleong/beholder/src"
)

// InputUI .
type InputUI struct {
	UI         tview.Primitive
	KeyHandler func(*tcell.EventKey) *tcell.EventKey
}

// NewInputUI .
func NewInputUI(app *beholder.App) *InputUI {

	ui := &InputUI{}

	prompt := tview.NewTextView()
	prompt.SetBackgroundColor(tcell.ColorDefault)
	prompt.SetTextColor(tcell.ColorDarkCyan)
	prompt.SetText("> ")

	input := tview.NewInputField()
	input.SetChangedFunc(app.OnQueryChanged)
	input.SetFieldBackgroundColor(tcell.ColorDefault)
	input.SetFieldTextColor(tcell.ColorDefault)

	input.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {
		case tcell.KeyESC:
			if input.GetText() != "" {
				input.SetText("")
				return nil
			}
		}
		return ui.KeyHandler(ev)
	})

	ui.UI = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(prompt, 2, 0, false).
		AddItem(input, 0, 1, true)
	return ui
}
