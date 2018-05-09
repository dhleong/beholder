package tui

import (
	"fmt"
	"strings"

	beholder "github.com/dhleong/beholder/src"
)

var valueSeparator = "  "

// ItemRenderer can render an Item
var ItemRenderer = &EntityRenderer{
	replacements: func(e beholder.Entity) []string {
		i := e.(beholder.Item)

		rs := []string{
			"{rarity}", i.Rarity,
		}

		value := i.Value
		if value == "" && i.Magic > 0 && i.Rarity != "" {
			// suggested magic item price ranges
			switch i.Rarity {
			case "Common":
				value = "50—100 gp"
			case "Uncommon":
				value = "101-500 gp"
			case "Rare":
				value = "501-5000 gp"
			case "Very Rare":
				value = "5001-50000 gp"
			case "Legendary":
				value = "50000+ gp"
			}
		}

		actualSeparator := ""
		if value != "" {
			rs = append(rs, "{value}", fmt.Sprintf(
				"[::d]%s[-:-:-]", value,
			))
		} else {
			rs = append(rs, "{value}", "")
		}

		if i.Magic > 0 {
			rs = append(rs, "{magic}", " (magic)")
			if value != "" {
				actualSeparator = valueSeparator
			}
		} else {
			rs = append(rs, "{magic}", "")
		}

		if i.Rarity != "" {
			if value != "" {
				actualSeparator = valueSeparator
			}

			rs = append(rs, "{text}", strings.Replace(
				strings.Join(i.GetText(), "\n"),
				fmt.Sprintf("Rarity: %s\n", i.Rarity),
				"",
				1,
			))
		}

		rs = append(rs, "{value-sep}", actualSeparator)

		return rs
	},

	template: `
[::bu]{name}[-:-:-]
[::d]{rarity}{magic}[-:-:-]{value-sep}{value}

{text}
`,
}
