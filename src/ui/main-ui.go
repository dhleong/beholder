package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	beholder "github.com/dhleong/beholder/src"
)

// NewMainUI .
func NewMainUI(app *beholder.App) tview.Primitive {
	// app := beholder.NewApp()

	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault

	grid := tview.NewGrid()

	input := NewInputUI(app)
	choices := NewChoicesUI(app)
	entity := NewEntityUI()

	grid.SetRows(0, 1)
	grid.SetColumns(-1, -4)

	// input spans the bottom row
	grid.AddItem(input.UI,
		1, 0,
		1, 2,
		0, 0,
		true,
	)

	// choices are above input
	grid.AddItem(choices.UI,
		0, 0,
		1, 1,
		0, 0,
		false,
	)

	grid.AddItem(entity.UI,
		0, 1,
		1, 1,
		0, 0,
		false,
	)

	input.KeyHandler = func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {
		case tcell.KeyESC:
			app.Quit()
			return nil
		case tcell.KeyCtrlJ:
			choices.Scroll(-1)
			return nil
		case tcell.KeyCtrlK:
			choices.Scroll(1)
			return nil
		}
		return ev
	}

	choices.SetChangedFunc(entity.Set)

	app.OnResults = func(results []beholder.Entity) {
		choices.Set(results)
	}

	return grid
}
