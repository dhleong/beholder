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

// Textual .
type Textual interface {
	GetText() []string
}

// textual is a mixin for anything with a sequences of <text>
type textual struct {
	Text []string `xml:"text"`
}

// GetText implements the Textual interface
func (t textual) GetText() []string {
	return t.Text
}

// Stats is a stat block for a creature
type Stats struct {
	ArmorClass            string `xml:"ac"`
	HP                    string `xml:"hp"`
	Speed                 string `xml:"speed"`
	Str                   int    `xml:"str"`
	Dex                   int    `xml:"dex"`
	Con                   int    `xml:"con"`
	Int                   int    `xml:"int"`
	Wis                   int    `xml:"wis"`
	Cha                   int    `xml:"cha"`
	PassivePerception     int    `xml:"passive"`
	SavingThrows          string `xml:"saving"`
	SkillModifiers        string `xml:"skill"`
	Senses                string `xml:"senses"`
	DamageImmunities      string `xml:"immune"`
	DamageResistances     string `xml:"resist"`
	DamageVulnerabilities string `xml:"vulnerable"`
	ConditionImmunities   string `xml:"conditionImmune"`
}

// Trait is a simple container for a Name and Text
type Trait struct {
	Named
	textual
}

// Traitor has Traits (bad pun, I know)
type Traitor interface {
	GetTraits() []*Trait
}

type traitor struct {
	Traits []*Trait `xml:"trait"`
}

func (t traitor) GetTraits() []*Trait {
	return t.Traits
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
	traitor
	Size      string   `xml:"size"`
	Type      string   `xml:"type"`
	Alignment string   `xml:"alignment"`
	Challenge string   `xml:"cr"`
	Languages string   `xml:"languages"`
	Actions   []*Trait `xml:"action"`
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
