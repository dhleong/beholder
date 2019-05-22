package tui

import (
	"strings"

	beholder "github.com/dhleong/beholder/src"
)

// ReferenceListRenderer can render a ReferenceListEntity
var ReferenceListRenderer = &EntityRenderer{
	replacements: func(e beholder.Entity) []string {
		r := e.(*beholder.ReferenceList)

		var text strings.Builder

		for _, ref := range r.References {
			text.WriteString(ref.GetName())
			text.WriteString("\n")
		}

		return []string{
			"{text}", text.String(),
		}
	},

	template: `
[::bu]{name}[-:-:-]

{text}
`,
}
