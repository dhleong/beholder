package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	beholder "github.com/dhleong/beholder/src"
)

// NewMainUI .
func NewMainUI(app *beholder.App, tapp *tview.Application) tview.Primitive {

	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault

	pages := tview.NewPages()
	mainGrid := tview.NewGrid() // contains topPane and input
	topPane := tview.NewPages() // alterntes between topGrid and home
	topGrid := tview.NewGrid()  // flexible container for choices + entity

	help := NewHelpUI()
	home := NewHomeUI()
	input := NewInputUI(app)
	choices := NewChoicesUI(app)
	entity := NewEntityUI()

	topGrid.SetRows(0)
	topGrid.SetColumns(15, -1, -4) // the first column is the min width for `choices`

	mainGrid.SetRows(0, 1)
	mainGrid.SetColumns(0)
	mainGrid.AddItem(topPane, 0, 0, 1, 1, 0, 0, false)
	mainGrid.AddItem(input.UI, 1, 0, 1, 1, 0, 0, true)

	// choices are above input; they have an exact width so as not to be too small
	// UNLESS the window is wide enough to allow for something more precise
	topGrid.AddItem(choices.UI,
		0, 0,
		1, 1,
		0, 0,
		false,
	)
	topGrid.AddItem(choices.UI,
		0, 0,
		1, 2, // if there's enough room for it to grow, do so
		0, 60,
		false,
	)

	topGrid.AddItem(entity.UI,
		0, 1,
		1, 2,
		0, 0,
		false,
	)
	topGrid.AddItem(entity.UI,
		0, 2, // when there's enough room for `choices` to grow, this doesn't need to use its column
		1, 1,
		0, 60,
		false,
	)

	topPane.AddAndSwitchToPage("empty", home.UI, true)
	topPane.AddPage("main", topGrid, true, false)

	pages.AddAndSwitchToPage("main", mainGrid, true)
	pages.AddPage("help", help.UI, true, false)

	var lastHelp HelpPage
	var oldFocus tview.Primitive
	showHelp := func(page HelpPage) {
		lastHelp = page
		oldFocus = tapp.GetFocus()
		help.SetPage(page)
		pages.SwitchToPage("help")
		tapp.SetFocus(help.UI)
	}

	input.KeyHandler = func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {

		case tcell.KeyF1:
			fallthrough
		case tcell.KeyRune:
			if ev.Rune() == '?' || ev.Key() == tcell.KeyF1 {
				showHelp(HelpPageHome)
				return nil
			}

		case tcell.KeyESC:
			app.Quit()
			return nil
		case tcell.KeyEnter:
			if choices.GetSelectedEntity() != nil {
				tapp.SetFocus(entity.UI)
				entity.SetFocused(true)
			}
			return nil

		case tcell.KeyDown:
			fallthrough
		case tcell.KeyCtrlP:
			fallthrough
		case tcell.KeyCtrlJ:
			choices.Scroll(-1)
			return nil

		case tcell.KeyUp:
			fallthrough
		case tcell.KeyCtrlN:
			fallthrough
		case tcell.KeyCtrlK:
			choices.Scroll(1)
			return nil

			// forward page up/down events to the entity view
		case tcell.KeyPgUp:
			fallthrough
		case tcell.KeyPgDn:
			handler := entity.text.InputHandler()
			handler(ev, nil)
		}
		return ev
	}

	entity.KeyHandler = func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {
		case tcell.KeyF1:
			fallthrough
		case tcell.KeyRune:
			if ev.Rune() == '?' || ev.Key() == tcell.KeyF1 {
				showHelp(HelpPageEntity)
				return nil
			}

		case tcell.KeyBackspace:
			fallthrough
		case tcell.KeyBackspace2:
			fallthrough
		case tcell.KeyESC:
			tapp.SetFocus(input.UI)
			entity.SetFocused(false)
			return nil
		}
		return ev
	}

	help.KeyHandler = func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {
		case tcell.KeyF1:
			fallthrough
		case tcell.KeyRune:
			if ev.Rune() == '?' || ev.Key() == tcell.KeyF1 {
				help.SetPage(HelpPageHelp)
				return nil
			}

		case tcell.KeyBackspace:
			fallthrough
		case tcell.KeyBackspace2:
			fallthrough
		case tcell.KeyESC:
			if help.CurrentPage == HelpPageHelp {
				help.SetPage(lastHelp)
			} else {
				tapp.SetFocus(oldFocus)
				pages.SwitchToPage("main")
			}
			return nil
		}

		return ev
	}

	choices.SetChangedFunc(entity.Set)

	app.OnResults = func(results []beholder.Entity) {
		// it'd be nice to have HomeUI handle this,
		// but even having an empty view in the grid
		// messes with drawing the other elements...
		if len(results) > 0 {
			topPane.SwitchToPage("main")
		} else if len(results) == 0 {
			topPane.SwitchToPage("empty")
		}

		choices.Set(results)
	}

	return pages
}
