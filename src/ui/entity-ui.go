package ui

import (
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	beholder "github.com/dhleong/beholder/src"
	"github.com/dhleong/beholder/src/ui/tui"
)

var entityRenderers = map[beholder.EntityKind]*tui.EntityRenderer{
	beholder.ActionEntity:    tui.NewSimpleRenderer(" Action"),
	beholder.ConditionEntity: tui.NewSimpleRenderer(" Condition"),
	beholder.FeatEntity:      tui.NewSimpleRenderer(" Feat"),
	beholder.FeatureEntity:   tui.FeatureRenderer,
	beholder.ItemEntity:      tui.ItemRenderer,
	beholder.MonsterEntity:   tui.MonsterRenderer,
	beholder.RuleEntity:      tui.NewSimpleRenderer(" Rule"),
	beholder.SpellEntity:     tui.SpellRenderer,
	beholder.TraitEntity:     tui.TraitRenderer,
}

// EntityUI .
type EntityUI struct {
	UI         tview.Primitive
	KeyHandler func(*tcell.EventKey) *tcell.EventKey

	text    *tview.TextView
	current beholder.Entity
}

// NewEntityUI .
func NewEntityUI() *EntityUI {
	text := tview.NewTextView()
	text.SetDynamicColors(true)
	text.SetTextColor(tcell.ColorDefault)
	text.SetBorderColor(tui.Colors.Border)
	text.SetWordWrap(true)

	ui := &EntityUI{
		UI:   text,
		text: text,
	}

	ui.SetFocused(false)

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

	if e.current != nil &&
		e.current.GetKind() == entity.GetKind() &&
		e.current.GetName() == entity.GetName() {
		return
	}
	e.current = entity

	r := entityRenderers[entity.GetKind()]
	if r != nil {
		e.text.SetText(strings.Trim(r.Render(entity), " \n\r"))
	} else {
		e.text.SetText(entity.GetName())
	}
	e.text.ScrollToBeginning()
}

// SetFocused .
func (e *EntityUI) SetFocused(isFocused bool) {
	if isFocused {
		e.text.SetBorderPadding(1, 1, 3, 3)
	} else {
		e.text.SetBorderPadding(2, 2, 4, 4)
	}

	e.text.SetBorder(isFocused)
}
