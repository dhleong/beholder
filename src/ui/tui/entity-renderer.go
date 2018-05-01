package tui

import (
	"bytes"
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
				formatted := formatText(t.GetText())
				replacements = append(replacements, "{text}", formatted)
			}
		}
	}

	if t, ok := entity.(beholder.Traitor); ok {
		traits := t.GetTraits()
		replacements = append(replacements, "{traits}", BuildTraitsString(traits))
	} else {
		replacements = append(replacements, "{traits}", "")
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

func formatText(text []string) string {
	// doesn't allow for nesting, but that's okay for now
	return strings.NewReplacer(
		"<h1>", "[::b]",
		"</h1>", "[::-]",
		"<h2>", "[::b]",
		"</h2>", "[::-]",
		"<b>", "[::b]",
		"</b>", "[::-]",
	).Replace(strings.Join(text, "\n"))
}

// BuildTraitsString .
func BuildTraitsString(traits []*beholder.Trait) string {
	if traits == nil {
		return ""
	}

	var traitsBuilder bytes.Buffer

	for i, trait := range traits {
		if i > 0 {
			traitsBuilder.WriteString("\n")
		}
		traitsBuilder.WriteString("[::b]")
		traitsBuilder.WriteString(trait.Name)
		traitsBuilder.WriteString("\n")
		traitsBuilder.WriteString(strings.Join(trait.GetText(), "\n"))

		if i < len(traits)-1 {
			traitsBuilder.WriteString("\n")
		}
	}

	return traitsBuilder.String()
}
