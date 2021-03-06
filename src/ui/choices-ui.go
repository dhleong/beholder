package ui

import (
	"github.com/rivo/tview"

	beholder "github.com/dhleong/beholder/src"
	"github.com/dhleong/beholder/src/ui/tui"
)

// ChoicesUI .
type ChoicesUI struct {
	UI tview.Primitive

	changedFunc func(entity beholder.Entity)
	list        *tui.EntityList
}

// Set the choices
func (c *ChoicesUI) Set(choices []beholder.SearchResult) {
	oldSelected := c.list.GetCurrentItem()

	// remove changed func briefly to avoid empty callback
	c.list.SetChangedFunc(nil)

	c.list.Clear()
	for _, entity := range choices {
		c.list.AddItem(entity)
	}

	c.list.SetChangedFunc(c.changedFunc)

	// persist selected position as best as possible
	if oldSelected < c.list.GetItemCount() {
		c.list.SetCurrentItem(oldSelected)
	} else {
		c.list.SetCurrentItem(c.list.GetItemCount() - 1)
	}
}

// Scroll by the given number of items
func (c *ChoicesUI) Scroll(items int) {
	c.list.SetCurrentItem(
		c.list.GetCurrentItem() + items,
	)
}

// GetSelectedEntity .
func (c *ChoicesUI) GetSelectedEntity() beholder.Entity {
	return c.list.GetCurrentEntity()
}

// SetChangedFunc .
func (c *ChoicesUI) SetChangedFunc(changed func(entity beholder.Entity)) {
	c.changedFunc = changed
	c.list.SetChangedFunc(changed)
}

// NewChoicesUI .
func NewChoicesUI(app *beholder.App) *ChoicesUI {
	list := tui.NewList()
	list.SetBorderPadding(0, 0, 0, 1)

	return &ChoicesUI{
		UI:   list,
		list: list,
	}
}
