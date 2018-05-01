package tui

import beholder "github.com/dhleong/beholder/src"

// NewSimpleRenderer will create a Renderer that can render
// anything that's just a name and text. The label will be
// rendered immediately after the name.
func NewSimpleRenderer(label string) *EntityRenderer {
	return &EntityRenderer{
		replacements: func(e beholder.Entity) []string {
			return []string{
				"{label}", label,
			}
		},

		template: `
[::bu]{name}[-:-:-]{label}

{text}
`,
	}
}
