package tui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	beholder "github.com/dhleong/beholder/src"
)

// EntityList .
type EntityList struct {
	*tview.Box

	entities []beholder.Entity

	currentItem int
	offset      int
}

// NewList returns a new form.
func NewList() *EntityList {
	return &EntityList{
		Box: tview.NewBox(),
	}
}

// AddItem .
func (l *EntityList) AddItem(entity beholder.Entity) {
	l.entities = append(l.entities, entity)
}

// GetCurrentItem returns the index of the currently selected list item.
func (l *EntityList) GetCurrentItem() int {
	return l.currentItem
}

// SetCurrentItem sets the index of the currently selected list item.
func (l *EntityList) SetCurrentItem(item int) {
	l.currentItem = item
}

// GetCurrentEntity returns the currently selected entity.
func (l *EntityList) GetCurrentEntity() beholder.Entity {
	return l.entities[l.currentItem]
}

// GetItemCount .
func (l *EntityList) GetItemCount() int {
	return len(l.entities)
}

// Clear all entities from this List
func (l *EntityList) Clear() {
	l.entities = nil
	l.currentItem = 0
	l.offset = 0
}

// SetEntities sets the entities
func (l *EntityList) SetEntities(entities []beholder.Entity) {
	l.entities = entities
	if l.currentItem >= len(entities) {
		l.currentItem = len(entities) - 1
	}
}

// Draw .
func (l *EntityList) Draw(screen tcell.Screen) {

	// Determine the dimensions.
	x, y, width, height := l.GetInnerRect()
	bottomLimit := y + height

	for i := l.offset; i < height && i < len(l.entities); i++ {
		e := l.entities[i]
		tview.Print(
			screen, e.GetName(), x, bottomLimit-i, width,
			tview.AlignLeft, tcell.ColorDefault,
		)
	}

}
