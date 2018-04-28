package beholder

import (
	"encoding/xml"
	"io"
)

type compendium struct {
	Items    []Item    `xml:"item"`
	Monsters []Monster `xml:"monster"`
	Spells   []Spell   `xml:"spell"`
}

// Named is a mixin for anything with a name.
type Named struct {
	Name string `xml:"name"`
}

// GetName helps implement the Entity interface
func (n Named) GetName() string {
	return n.Name
}

type textual struct {
	Text []string `xml:"text"`
}

// Stats is a stat block for a creature
type Stats struct {
	ArmorClass          string `xml:"AC"`
	HP                  string `xml:"hp"`
	Speed               string `xml:"speed"`
	Str                 int    `xml:"str"`
	Dex                 int    `xml:"dex"`
	Con                 int    `xml:"con"`
	Int                 int    `xml:"int"`
	Wis                 int    `xml:"wis"`
	Cha                 int    `xml:"cha"`
	PassivePerception   int    `xml:"passive"`
	SavingThrows        string `xml:"saving"`
	SkillModifiers      string `xml:"skill"`
	Senses              string `xml:"senses"`
	DamageImmunities    string `xml:"immune"`
	ConditionImmunities string `xml:"conditionImmune"`
}

// Item .
type Item struct {
	Named
	textual
	Type     string `xml:"type"`
	Magic    int    `xml:"magic"`
	Value    string `xml:"value"`
	Weight   string `xml:"weight"`
	Property string `xml:"property"`
	Rarity   string `xml:"rarity"`
}

// GetKind from Entity interface
func (i Item) GetKind() EntityKind {
	return ItemEntity
}

// Monster .
type Monster struct {
	Named
	Stats
	Size      string `xml:"size"`
	Type      string `xml:"type"`
	Alignment string `xml:"alignment"`
	Challenge string `xml:"cr"`
	Languages string `xml:"languages"`
}

// GetKind from Entity interface
func (m Monster) GetKind() EntityKind {
	return MonsterEntity
}

// Spell .
type Spell struct {
	Named
	textual
	Level      int    `xml:"level"`
	School     string `xml:"school"`
	Time       string `xml:"time"`
	Range      string `xml:"range"`
	Components string `xml:"components"`
	Duration   string `xml:"duration"`
	Classes    string `xml:"classes"`
}

// GetKind from Entity interface
func (s Spell) GetKind() EntityKind {
	return SpellEntity
}

// ParseXML extracts Entity instances from the Reader
func ParseXML(reader io.Reader) ([]Entity, error) {
	result := make([]Entity, 0, 4096)

	compendium := &compendium{}
	decoder := xml.NewDecoder(reader)
	err := decoder.Decode(compendium)
	if err != nil {
		return nil, err
	}

	// it'd be nice if I could just `append(result, Spells...)` :\
	for _, entity := range compendium.Spells {
		result = append(result, entity)
	}
	for _, entity := range compendium.Items {
		result = append(result, entity)
	}
	for _, entity := range compendium.Monsters {
		result = append(result, entity)
	}

	return result, nil
}
