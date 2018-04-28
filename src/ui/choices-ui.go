package ui

import (
	"fmt"

	"github.com/rivo/tview"

	beholder "github.com/dhleong/beholder/src"
	"github.com/dhleong/beholder/src/ui/tui"
)

// ChoicesUI .
type ChoicesUI struct {
	UI tview.Primitive

	list *tui.EntityList
}

// Set the choices
func (c *ChoicesUI) Set(choices []beholder.Entity) {
	c.list.Clear()
	for i := 0; i < 20; i++ {
		c.list.AddItem(&beholder.SpellEntity{fmt.Sprintf("%d", i)})
	}
	c.list.SetCurrentItem(c.list.GetItemCount() - 1)
}

// Scroll by the given number of items
func (c *ChoicesUI) Scroll(items int) {
	c.list.SetCurrentItem(
		c.list.GetCurrentItem() + items,
	)
}

// NewChoicesUI .
func NewChoicesUI(app *beholder.App) *ChoicesUI {
	list := tui.NewList()

	return &ChoicesUI{
		UI:   list,
		list: list,
	}
}
