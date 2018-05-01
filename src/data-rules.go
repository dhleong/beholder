package beholder

import "fmt"

type ruleEntity struct {
	Named
	textual
}

// GetKind from Entity interface
func (r ruleEntity) GetKind() EntityKind {
	return RuleEntity
}

// ruleParts is a recursive structure representing a Rule and
// its sub-rules. One Entity will be generated for top-most Rule
// (IE the one which is not part of any other rule) including all
// of its sub-rules as "sections," and one for each sub-rule as well.
type ruleParts struct {
	name  string
	parts []interface{}
}

func rule(name string, parts ...interface{}) ruleParts {
	return ruleParts{
		name,
		parts,
	}
}

// section is a semantic alias for rule
var section = rule

// Generate entities given ruleParts, returning the first
// ruleEntity and a slice of *all* entities generated (*including*
// the first one)
func generateEntities(part ruleParts) (ruleEntity, []Entity) {
	text := make([]string, 0, len(part.parts))
	parent := &ruleEntity{
		Named: Named{part.name},
	}
	generated := []Entity{parent}

	for _, child := range part.parts {
		switch v := child.(type) {
		case string:
			text = append(text, v)

		case ruleParts:
			// NOTE: if we wanted to get fancy, we could replace all the
			// headers fromt he child entity with, eg: h2 -> h3
			childEntity, newDescendants := generateEntities(v)
			generated = append(generated, newDescendants...)

			text = append(text,
				"", // start with a line break
				fmt.Sprintf("<h2>%s</h2>", childEntity.Name),
			)
			text = append(text,
				childEntity.GetText()...,
			)
		}
	}

	// install now, through the pointer; otherwise it gets lost
	parent.textual = textual{text}

	return *parent, generated
}

func rulesDataSource(rules ...ruleParts) DataSource {
	entities := make([]Entity, 0, len(rules))

	for _, part := range rules {
		_, generated := generateEntities(part)
		entities = append(entities, generated...)
	}

	return NewStaticDataSource(entities)
}

// RuleDataSource .
var RuleDataSource = rulesDataSource(
	rule("Death, Dying, Dropping to 0 Hitpoints",
		"When you drop to 0 hit points, you either die outright or fall unconscious, as explained in the following sections.",

		section("Instant Death",
			"	Massive damage can kill you instantly. When damage reduces you to 0 hit points and there is damage remaining, you die if the remaining damage equals or exceeds your hit point maximum.",
			"	For example, a cleric with a maximum of 12 hit points currently has 6 hit points. If she takes 18 damage from an attack, she is reduced to 0 hit points, but 12 damage remains. Because the remaining damage equals her hit point maximum, the cleric dies.",
		),

		section("Falling Unconscious",
			"	If damage reduces you to 0 hit points and fails to kill you, you fall unconscious. This unconsciousness ends if you regain any hit points.",
		),

		section("Death Saving Throws",
			"	Whenever you start your turn with 0 hit points, you must make a special saving throw, called a death saving throw, to determine whether you creep closer to death or hang onto life. Unlike other saving throws, this one isn't tied to any ability score. You are in the hands of fate now, aided only by spells and features that improve your chances of succeeding on a saving throw.",
			"	Roll a d20. If the roll is 10 or higher, you succeed. Otherwise, you fail. A success or failure has no effect by itself. On your third success, you become stable (see below). On your third failure, you die. The successes and failures don't need to be consecutive; keep track of both until you collect three of a kind. The number of both is reset to zero when you regain any hit points or become stable.",
			"	<b>Rolling 1 or 20.</b> When you make a death saving throw and roll a 1 on the d20, it counts as two failures. If you roll a 20 on the d20, you regain 1 hit point.",
			"	<b>Damage at 0 Hit Points.</b> If you take any damage while you have 0 hit points, you suffer a death saving throw failure. If the damage is from a critical hit, you suffer two failures instead. If the damage equals or exceeds your hit point maximum, you suffer instant death.",
		),

		section("Stabilizing a Creature",
			"	The best way to save a creature with 0 hit points is to heal it. If healing is unavailable, the creature can at least be stabilized so that it isn't killed by a failed death saving throw.",
			"	You can use your action to administer first aid to an unconscious creature and attempt to stabilize it, which requires a successful DC 10 Wisdom (Medicine) check.",
			"	A <b>stable</b> creature doesn't make death saving throws, even though it has 0 hit points, but it does remain unconscious. The creature stops being stable, and must start making death saving throws again, if it takes any damage. A stable creature that isn't healed regains 1 hit point after <b>1d4</b> hours.",
		),

		section("Monsters and Death",
			"	Most GMs have a monster die the instant it drops to 0 hit points, rather than having it fall unconscious and make death saving throws.",
			"	Mighty villains and special nonplayer characters are common exceptions; the GM might have them fall unconscious and follow the same rules as player characters.",
		),

		section("Knocking a creature out",
			"	Sometimes an attacker wants to incapacitate a foe, rather than deal a killing blow. When an attacker reduces a creature to 0 hit points with a melee attack, the attacker can knock the creature out. The attacker can make this choice the instant the damage is dealt. The creature falls unconscious and is stable.",
		),

		section("Temporary Hit Points",
			"	Some spells and special abilities confer temporary hit points to a creature. Temporary hit points aren't actual hit points; they are a buffer against damage, a pool of hit points that protect you from injury.",
			"	When you have temporary hit points and take damage, the temporary hit points are lost first, and any leftover damage carries over to your normal hit points. For example, if you have 5 temporary hit points and take 7 damage, you lose the temporary hit points and then take 2 damage.",
			"	Because temporary hit points are separate from your actual hit points, they can exceed your hit point maximum. A character can, therefore, be at full hit points and receive temporary hit points.",
			"	Healing can't restore temporary hit points, and they can't be added together. If you have temporary hit points and receive more of them, you decide whether to keep the ones you have or to gain the new ones. For example, if a spell grants you 12 temporary hit points when you already have 10, you can have 12 or 10, not 22.",
			"	If you have 0 hit points, receiving temporary hit points doesn't restore you to consciousness or stabilize you. They can still absorb damage directed at you while you're in that state, but only true healing can save you.",
			"	Unless a feature that grants you temporary hit points has a duration, they last until they're depleted or you finish a long rest.",
		),
	),
)
