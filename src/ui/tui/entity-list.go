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
	if item < 0 {
		l.currentItem = 0
		l.offset = 0
		return
	} else if item >= len(l.entities) {
		item = len(l.entities) - 1
	}

	l.currentItem = item

	// scroll the view to make currentItem visible
	_, _, _, height := l.Box.GetInnerRect()
	if item-l.offset < 0 {
		l.offset = item
	} else if item-l.offset >= height {
		l.offset = item - height + 1
	}
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

	var bgStyle tcell.Style
	bgStyle = bgStyle.Background(Colors.SelectedBackground)

	rowY := bottomLimit
	for i := l.offset; rowY >= y && i < len(l.entities); i++ {
		e := l.entities[i]

		rowY--
		itemX := x + 2
		itemWidth := width - 2
		screen.SetContent(x, rowY, ' ', nil,
			bgStyle,
		)
		tview.Print(
			screen, e.GetName(), itemX, rowY, itemWidth,
			tview.AlignLeft, tcell.ColorDefault,
		)

		// "selected" background
		if i == l.currentItem {
			textWidth := tview.StringWidth(e.GetName())
			textEnd := itemX + textWidth

			for bx := 1; bx < textEnd && bx < width; bx++ {
				m, c, style, _ := screen.GetContent(x+bx, rowY)

				style = style.
					Background(Colors.SelectedBackground).
					Foreground(Colors.SelectedFg)
				screen.SetContent(x+bx, rowY, m, c, style)
			}
		}
	}

}
