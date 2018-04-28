package tui

import "github.com/gdamore/tcell"

// ColorScheme is a struct containing... well, you can guess
type ColorScheme struct {
	SelectedBackground tcell.Color
	SelectedFg         tcell.Color
}

// Colors is the current color scheme
var Colors = ColorScheme{
	SelectedBackground: tcell.GetColor("#333333"),
	SelectedFg:         tcell.GetColor("#f0f0f0"),
}
