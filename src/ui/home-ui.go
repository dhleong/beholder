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

type verticalCenterLayout struct {
	*tview.Box
	child       *tview.TextView
	childHeight int
}

func (l *verticalCenterLayout) Draw(screen tcell.Screen) {
	x, y, w, h := l.GetRect()

	childTop := (h / 2) - (l.childHeight / 2)

	l.child.SetRect(x, y+childTop, w, l.childHeight)
	l.child.Draw(screen)
}

// NewHomeUI .
func NewHomeUI() *HomeUI {

	text := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorDefault).
		SetText(strings.Trim(`
[::b]Welcome to Beholder![::-]

Just start typing to search for an item,
monster, spell, feat, rule, condition,
or class or race feature.

Press [::b]F1[::-] or [::b]?[::-] for detailed help anywhere.
`, " \n"))

	// flex := tview.NewFlex().
	// 	SetDirection(tview.FlexRow)
	// flex.SetBorderPadding(2, 1, 2, 2)
	// flex.AddItem(text, 0, 1, false)

	return &HomeUI{
		// UI: flex,
		UI: &verticalCenterLayout{
			Box:         tview.NewBox(),
			child:       text,
			childHeight: 7,
		},
	}
}
