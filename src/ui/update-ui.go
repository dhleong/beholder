package ui

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// UpdateUI .
type UpdateUI struct {
	UI tview.Primitive
}

// NewUpdateUI .
func NewUpdateUI(newVersion string) *UpdateUI {

	brewInstruction := ""
	if runtime.GOOS == "darwin" {
		brewInstruction = "\nOr, run: brew upgrade beholder"
	}

	text := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorDefault).
		SetDynamicColors(true).
		SetText(strings.Trim(fmt.Sprintf(`
Welcome to Beholder!

[::b]Version [::bu]%s[::b] is now available![::-]

Press [::b]F5[::-] to get it!%s
`, newVersion, brewInstruction), " \n"))

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow)
	flex.SetBorderPadding(2, 1, 2, 2)
	flex.AddItem(text, 0, 1, false)

	return &UpdateUI{
		UI: flex,
	}
}
