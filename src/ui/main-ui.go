package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	beholder "github.com/dhleong/beholder/src"
)

// NewMainUI .
func NewMainUI(app *beholder.App, tapp *tview.Application) tview.Primitive {
	// app := beholder.NewApp()

	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault

	grid := tview.NewGrid()

	input := NewInputUI(app)
	choices := NewChoicesUI(app)
	entity := NewEntityUI()

	grid.SetRows(0, 1)
	grid.SetColumns(15, -1, -4) // the first column is the min width for `choices`

	// input spans the bottom row
	grid.AddItem(input.UI,
		1, 0,
		1, 3,
		0, 0,
		true,
	)

	// choices are above input; they have an exact width so as not to be too small
	// UNLESS the window is wide enough to allow for something more precise
	grid.AddItem(choices.UI,
		0, 0,
		1, 1,
		0, 0,
		false,
	)
	grid.AddItem(choices.UI,
		0, 0,
		1, 2, // if there's enough room for it to grow, do so
		0, 60,
		false,
	)

	grid.AddItem(entity.UI,
		0, 1,
		1, 2,
		0, 0,
		false,
	)
	grid.AddItem(entity.UI,
		0, 2, // when there's enough room for `choices` to grow, this doesn't need to use its column
		1, 1,
		0, 60,
		false,
	)

	input.KeyHandler = func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {
		case tcell.KeyESC:
			app.Quit()
			return nil
		case tcell.KeyEnter:
			if choices.GetSelectedEntity() != nil {
				tapp.SetFocus(entity.UI)
			}
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

	entity.KeyHandler = func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {
		case tcell.KeyBackspace:
			fallthrough
		case tcell.KeyBackspace2:
			fallthrough
		case tcell.KeyESC:
			tapp.SetFocus(input.UI)
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
