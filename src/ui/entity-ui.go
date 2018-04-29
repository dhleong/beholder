package ui

import (
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	beholder "github.com/dhleong/beholder/src"
	"github.com/dhleong/beholder/src/ui/tui"
)

var entityRenderers = map[beholder.EntityKind]*tui.EntityRenderer{
	beholder.ItemEntity:    tui.ItemRenderer,
	beholder.MonsterEntity: tui.MonsterRenderer,
	beholder.SpellEntity:   tui.SpellRenderer,
	beholder.TraitEntity:   tui.TraitRenderer,
}

// EntityUI .
type EntityUI struct {
	UI         tview.Primitive
	KeyHandler func(*tcell.EventKey) *tcell.EventKey

	text *tview.TextView
}

// NewEntityUI .
func NewEntityUI() *EntityUI {
	text := tview.NewTextView()
	text.SetBorderPadding(2, 2, 4, 4)
	text.SetDynamicColors(true)
	text.SetTextColor(tcell.ColorDefault)

	ui := &EntityUI{
		UI:   text,
		text: text,
	}

	text.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		return ui.KeyHandler(ev)
	})

	return ui
}

// Set the current Entity to be displayed
func (e *EntityUI) Set(entity beholder.Entity) {
	if entity == nil {
		e.text.SetText("")
		return
	}

	r := entityRenderers[entity.GetKind()]
	if r != nil {
		e.text.SetText(strings.Trim(r.Render(entity), " \n\r"))
	} else {
		e.text.SetText(entity.GetName())
	}
	e.text.ScrollToBeginning()
}
