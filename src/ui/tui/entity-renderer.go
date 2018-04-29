package tui

import (
	"strings"

	beholder "github.com/dhleong/beholder/src"
)

// EntityRenderer can render an Entity type
type EntityRenderer struct {
	replacements func(beholder.Entity) []string
	template     string
}

// Render the given entity to a TUI string
func (r *EntityRenderer) Render(entity beholder.Entity) string {
	replacer := r.replacer(entity)
	return replacer.Replace(r.template)
}

func (r *EntityRenderer) replacer(entity beholder.Entity) *strings.Replacer {
	replacements := []string{"{name}", entity.GetName()}
	replacements = append(replacements, r.replacements(entity)...)

	if t, ok := entity.(beholder.Textual); ok {
		if !contains(replacements, "{text}") {
			text := t.GetText()
			if text == nil {
				replacements = append(replacements, "{text}", "")
			} else {
				replacements = append(replacements, "{text}", strings.Join(t.GetText(), "\n"))
			}
		}
	}
	return strings.NewReplacer(replacements...)
}

func contains(haystack []string, needle string) bool {
	for _, straw := range haystack {
		if straw == needle {
			return true
		}
	}

	return false
}
