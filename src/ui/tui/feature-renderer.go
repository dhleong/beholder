package tui

import (
	"strings"

	beholder "github.com/dhleong/beholder/src"
)

// FeatureRenderer can render a ClassFeature
var FeatureRenderer = &EntityRenderer{
	replacements: func(e beholder.Entity) []string {
		f := e.(*beholder.ClassFeature)

		return []string{
			"{classes}", strings.Join(f.Classes, ", "),
		}
	},

	template: `
[::bu]{name}[-:-:-]
[::d]{classes}[-:-:-]

{text}
`,
}
