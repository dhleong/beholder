package ui

import (
	"strings"

	"github.com/dhleong/beholder/src/ui/tui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// HelpPage enum
type HelpPage int

// HelpPages
const (
	HelpPageHome = iota
	HelpPageEntity
	HelpPageHelp
)

var helpPages = map[HelpPage]string{
	HelpPageHome: trim(`
This is the [::b]home screen[::-].

Just start typing to search for an item, monster, spell, or class or race feature. Searching is pretty smart, you don't have to spell it out completely: "mcw," for example, can find the spell "Mass Cure Wounds."

Results will show up in the [::b]info window[::-] on the right.

You can scroll through the results with the arrow keys up and down, or <ctrl-k> and <ctrl-j>, or <ctrl-n> and <ctrl-p>.

If you can't see everything in the [::b]info window[::-], you can scroll it with the page up/page down keys, or press Enter to focus on it and get more precise scrolling.
`),

	HelpPageEntity: trim(`
You've focused on the [::b]info window[::-].

You can scroll with the arrow keys, j and k, and page up and page down.

Press Backspace/Delete or the Escape key to return focus to the search bar.
`),

	HelpPageHelp: trim(`
This is the [::b]help page[::-] for a [::b]help page[::-].

You got here by pressing F1 or ? on another [::b]help page[::-].

Good job!

Press the Escape key to leave.
`),
}

// HelpUI .
type HelpUI struct {
	UI          tview.Primitive
	KeyHandler  func(ev *tcell.EventKey) *tcell.EventKey
	CurrentPage HelpPage

	text *tview.TextView
}

// NewHelpUI .
func NewHelpUI() *HelpUI {

	text := tview.NewTextView().
		SetScrollable(true).
		SetWordWrap(true).
		SetDynamicColors(true).
		SetTextColor(tcell.ColorDefault)

	text.SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(tui.Colors.Border).
		SetTitleColor(tcell.ColorDefault).
		SetTitle("Beholder Help")

	text.SetDrawFunc(func(
		screen tcell.Screen,
		x, y, w, h int,
	) (int, int, int, int) {

		tview.Print(screen, "Use ↑/↓ or j/k to scroll",
			x,
			y+h-1,
			w,
			tview.AlignCenter,
			tcell.ColorDefault,
		)

		// subtract out for border and padding
		x += 2
		y += 2
		w -= 4
		h -= 4
		return x, y, w, h
	})

	ui := &HelpUI{
		UI: text,

		text: text,
	}

	text.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		return ui.KeyHandler(ev)
	})

	return ui
}

// SetPage .
func (ui *HelpUI) SetPage(page HelpPage) {
	ui.CurrentPage = page
	ui.text.SetText(helpPages[page])
}

func trim(text string) string {
	return strings.Trim(text, " \n")
}
