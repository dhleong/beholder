package tui

import (
	"bytes"
	"fmt"
	"strconv"

	beholder "github.com/dhleong/beholder/src"
)

var sizes = map[string]string{
	"T": "Tiny",
	"S": "Small",
	"M": "Medium",
	"L": "Large",
	"H": "Huge",
	"G": "Gargantuan",
}

// MonsterRenderer can render an Monster
var MonsterRenderer = &EntityRenderer{
	replacements: func(e beholder.Entity) []string {
		m := e.(beholder.Monster)

		var immunities bytes.Buffer
		if m.DamageResistances != "" {
			immunities.WriteString("\n[::b]Damage Resistances[::-]: ")
			immunities.WriteString(m.DamageResistances)
		}
		if m.DamageVulnerabilities != "" {
			immunities.WriteString("\n[::b]Damage Vulnerabilities[::-]: ")
			immunities.WriteString(m.DamageVulnerabilities)
		}
		if m.DamageImmunities != "" {
			immunities.WriteString("\n[::b]Damage Immunities[::-]: ")
			immunities.WriteString(m.DamageImmunities)
		}
		if m.ConditionImmunities != "" {
			immunities.WriteString("\n[::b]Condition Immunities[::-]: ")
			immunities.WriteString(m.ConditionImmunities)
		}

		var savesRow = ""
		if m.SavingThrows != "" {
			savesRow = fmt.Sprintf("\n[::b]Saves[::-]: %s", m.SavingThrows)
		}

		var sensesRow = ""
		if m.Senses != "" {
			sensesRow = fmt.Sprintf("\n[::b]Senses[::-]: %s", m.Senses)
		}

		var skillsRow = ""
		if m.SkillModifiers != "" {
			skillsRow = fmt.Sprintf("\n[::b]Skills[::-]: %s", m.SkillModifiers)
		}

		var languagesRow = ""
		if m.Languages != "" {
			languagesRow = fmt.Sprintf("\n[::b]Languages[::-]: %s", m.Languages)
		}

		var allActions bytes.Buffer
		BuildTraits(&allActions, m.Actions)

		if m.Legendary != nil {
			if allActions.Len() > 0 {
				allActions.WriteString("\n")
			}

			allActions.WriteString("\n[::bu]Legendary Actions\n")
			BuildTraits(&allActions, m.Legendary)
		}

		statFormat := "%2d[::bd]%s[::-]"
		return []string{
			"{size}", sizes[m.Size],
			"{type}", m.Type,
			"{cr}", m.Challenge,
			"{immunities}", immunities.String(),
			"{actions}", allActions.String(),

			// stat block:
			"{ac}", m.ArmorClass,
			"{hp}", m.HP,
			"{speed}", m.Speed,
			"{str}", fmt.Sprintf(statFormat, m.Str, formatModifier(m.Str)),
			"{dex}", fmt.Sprintf(statFormat, m.Dex, formatModifier(m.Dex)),
			"{con}", fmt.Sprintf(statFormat, m.Con, formatModifier(m.Con)),
			"{int}", fmt.Sprintf(statFormat, m.Int, formatModifier(m.Int)),
			"{wis}", fmt.Sprintf(statFormat, m.Wis, formatModifier(m.Wis)),
			"{cha}", fmt.Sprintf(statFormat, m.Cha, formatModifier(m.Cha)),
			"{passive}", strconv.Itoa(m.PassivePerception),
			"{saves-row}", savesRow,
			"{skills-row}", skillsRow,
			"{languages-row}", languagesRow,
			"{senses-row}", sensesRow,
		}
	},

	template: `
[::bu]{name}[-:-:-]
[::d]{size} {type}  CR {cr}[-:-:-]

[::b]AC[::-]: {ac}
[::b]HP[::-]: {hp}
[::b]Speed[::-]: {speed}
[::b]Passive Perception[::-]: {passive}

[::b] STR    DEX    CON    INT    WIS    CHA[::-]
[::-]{str}  {dex}  {con}  {int}  {wis}  {cha}[::-]
{immunities}{saves-row}{senses-row}{skills-row}{languages-row}

{traits}

{actions}
`,
}

func formatModifier(stat int) string {
	modifier := (stat - 10) / 2
	if modifier == 0 {
		return " +0"
	} else if modifier < 0 {
		return fmt.Sprintf("%-3d", modifier)
	} else {
		modifierStr := fmt.Sprintf("+%d", modifier)
		return fmt.Sprintf("%3s", modifierStr)
	}
}
