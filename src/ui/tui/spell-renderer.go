package tui

import (
	"strconv"

	beholder "github.com/dhleong/beholder/src"
)

var schools = map[string]string{
	"A":  "Abjuration",
	"C":  "Conjuration",
	"D":  "Divination",
	"EN": "Enchantment",
	"EV": "Evocation",
	"I":  "Illusion",
	"N":  "Necromancy",
	"T":  "Transmutation",
}

// SpellRenderer can render a Spell
var SpellRenderer = &EntityRenderer{
	replacements: func(e beholder.Entity) []string {
		s := e.(beholder.Spell)

		var level string
		if s.Level >= 1 {
			level = strconv.Itoa(s.Level)
		} else {
			level = "Cantrip"
		}

		return []string{
			"{level}", level,
			"{school}", schools[s.School],
			"{cast-time}", s.Time,
			"{range}", s.Range,
			"{components}", s.Components,
			"{duration}", s.Duration,
			"{classes}", s.Classes,
		}
	},

	template: `
[::bu]{name}[-:-:-]
[::d]{level} {school}[-:-:-]

[::b]Casting Time[::-]: {cast-time}
[::b]Range[::-]: {range}
[::b]Components[::-]: {components}
[::b]Duration[::-]: {duration}
[::b]Classes[::-]: {classes}

{text}
`,
}
