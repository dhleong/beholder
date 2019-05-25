package beholder

import (
	"fmt"
	"sort"
)

type spellListsSource struct {
	delegate DataSource
}

// ClassSpell collects all the Variants of a given Class that
// can use the provided Spell
type ClassSpell struct {
	Spell        *Spell
	Variants     []string
	VariantsOnly bool
}

// GetKind from Entity interface
func (c ClassSpell) GetKind() EntityKind {
	return SpellEntity
}

// GetName from Entity interface
func (c ClassSpell) GetName() string {
	return c.Spell.GetName()
}

// GetCategory from CategorizedEntity interface
func (c ClassSpell) GetCategory() string {
	if c.Spell.Level == 0 {
		return "Cantrips"
	}

	return fmt.Sprintf("Level %d Spells", c.Spell.Level)
}

func (s *spellListsSource) GetEntities() ([]Entity, error) {
	entities, err := s.delegate.GetEntities()
	if err != nil {
		return nil, err
	}

	classToSpells := map[string][]*ClassSpell{}

	for _, e := range entities {
		if e.GetKind() != SpellEntity {
			continue
		}

		spell := e.(Spell)
		extractClassSpells(classToSpells, spell)
	}

	for class, spells := range classToSpells {
		// convert to a slice of Entity interfaces
		spellEntities := make([]Entity, 0, len(spells))
		for _, s := range spells {
			spellEntities = append(spellEntities, s)
		}

		sort.Slice(spellEntities, func(i, j int) bool {
			a := spellEntities[i].(*ClassSpell)
			b := spellEntities[j].(*ClassSpell)
			return compareSpells(a.Spell, b.Spell)
		})

		entities = append(entities, &ReferenceList{
			Named:      Named{class + " Spells"},
			References: spellEntities,
		})
	}

	return entities, nil
}

// NewSpellListsSource creates a dynamic DataSource that
// augments the one provided with extra entities containing
// the list of spells available for each known class
func NewSpellListsSource(delegate DataSource) DataSource {
	return &spellListsSource{
		delegate: delegate,
	}
}

func extractClassSpells(dest map[string][]*ClassSpell, spell Spell) {
	for _, spellUser := range spell.SpellUsers {
		destList := dest[spellUser.Name]
		var lastClassSpell *ClassSpell
		if len(destList) > 0 {
			lastClassSpell = destList[len(destList)-1]
		}

		var targetClassSpell *ClassSpell

		if lastClassSpell == nil || lastClassSpell.Spell.Name != spell.Name {
			// new spell for the class
			s := &ClassSpell{
				Spell:        &spell,
				VariantsOnly: true, // assume true by default
			}
			dest[spellUser.Name] = append(destList, s)

			targetClassSpell = s
		} else {
			// existing spell
			targetClassSpell = lastClassSpell
		}

		if spellUser.Variant != "" {
			targetClassSpell.Variants = append(targetClassSpell.Variants, spellUser.Variant)
		} else {
			targetClassSpell.VariantsOnly = false
		}
	}
}

func compareSpells(a, b *Spell) bool {
	if a.Level < b.Level {
		return true
	} else if a.Level > b.Level {
		return false
	}

	// at the same level, sort alphabetically

	return a.GetName() < b.GetName()
}
