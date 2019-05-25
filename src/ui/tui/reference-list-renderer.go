package tui

import (
	"bytes"

	beholder "github.com/dhleong/beholder/src"
)

// ReferenceListRenderer can render a ReferenceListEntity
var ReferenceListRenderer = &EntityRenderer{
	replacements: func(e beholder.Entity) []string {
		r := e.(*beholder.ReferenceList)

		var text bytes.Buffer
		lastCategory := ""
		hasAnyCategory := false

		for _, ref := range r.References {
			if categorized, ok := ref.(beholder.CategorizedEntity); ok {
				hasAnyCategory = true
				category := categorized.GetCategory()
				if category != lastCategory {
					lastCategory = category

					if text.Len() > 0 {
						text.WriteString("\n")
					}

					text.WriteString(`[::bu]`)
					text.WriteString(category)
					text.WriteString(`[-:-:-]`)
					text.WriteString(":\n\n")
				}
			}

			if hasAnyCategory {
				// indent the contents
				text.WriteString("  ")
			}

			text.WriteString(formatString(ref.GetName()))
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
