package tui

import (
	"fmt"
	"strings"

	beholder "github.com/dhleong/beholder/src"
)

// ItemRenderer can render an Item
var ItemRenderer = &EntityRenderer{
	replacements: func(e beholder.Entity) []string {
		i := e.(beholder.Item)

		rs := []string{
			"{rarity}", i.Rarity,
		}

		if i.Rarity != "" {
			rs = append(rs, "{text}", strings.Replace(
				strings.Join(i.GetText(), "\n"),
				fmt.Sprintf("Rarity: %s\n", i.Rarity),
				"",
				1,
			))
		}

		return rs
	},

	template: `
[::bu]{name}[-:-:-]
[::d]{rarity}[-:-:-]

{text}
`,
}
