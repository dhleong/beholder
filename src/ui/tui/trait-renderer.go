package tui

import (
	"strings"

	beholder "github.com/dhleong/beholder/src"
)

// TraitRenderer can render a RaceTrait
var TraitRenderer = &EntityRenderer{
	replacements: func(e beholder.Entity) []string {
		t := e.(*beholder.RaceTrait)

		return []string{
			"{races}", strings.Join(t.Races, ", "),
		}
	},

	template: `
[::bu]{name}[-:-:-]
[::d]{races}[-:-:-]

{text}
`,
}
