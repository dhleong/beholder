package tui

import "github.com/gdamore/tcell/v2"

// ColorScheme is a struct containing... well, you can guess
type ColorScheme struct {
	Border             tcell.Color
	Highlight          tcell.Color
	SelectedBackground tcell.Color
	SelectedFg         tcell.Color
	SelectedHighlight  tcell.Color
}

// Colors is the current color scheme
var Colors = ColorScheme{
	Border:             tcell.GetColor("#333333"),
	Highlight:          tcell.GetColor("#3f7f3f"),
	SelectedBackground: tcell.GetColor("#333333"),
	SelectedFg:         tcell.GetColor("#f0f0f0"),
	SelectedHighlight:  tcell.GetColor("#3f9f3f"),
}
