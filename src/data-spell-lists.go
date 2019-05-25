package beholder

import (
	"fmt"
	"sort"
	"strings"
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

	// now that we have all the class names, remove
	// any existing spell lists
	resultEntities := make([]Entity, 0, len(entities))
	for _, e := range entities {
		name := e.GetName()
		if strings.HasSuffix(name, " Spells") {
			className := name[0:strings.Index(name, " Spells")]
			if _, ok := classToSpells[className]; ok {
				// we have class spells for this class; don't include it!
				continue
			}
		}

		// normal case; copy over the entity
		resultEntities = append(resultEntities, e)
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

		resultEntities = append(resultEntities, &ReferenceList{
			Named:      Named{class + " Spells"},
			References: spellEntities,
		})
	}

	return resultEntities, nil
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
