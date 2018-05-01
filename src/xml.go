package beholder

import (
	"encoding/xml"
	"io"
	"strings"
)

type compendium struct {
	Classes  []Class   `xml:"class"`
	Items    []Item    `xml:"item"`
	Monsters []Monster `xml:"monster"`
	Races    []Race    `xml:"race"`
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

// Class .
type Class struct {
	Named

	HitDice      int    `xml:"hd"`
	Proficiency  string `xml:"proficiency"`
	SpellAbility string `xml:"spellAbility"`

	Levels []Level `xml:"autolevel"`
}

// Level represents features, etc. granted at a given level
type Level struct {
	Level    int      `xml:"level,attr"`
	Features []*Trait `xml:"feature"`
}

// ClassFeature is basically a Trait for one or more classes
type ClassFeature struct {
	*Trait
	Classes []string
}

// GetKind from Entity interface
func (t ClassFeature) GetKind() EntityKind {
	return FeatureEntity
}

// RaceTrait is a Trait for one or more races
type RaceTrait struct {
	*Trait
	Races []string
}

// GetKind from Entity interface
func (t RaceTrait) GetKind() EntityKind {
	return TraitEntity
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

// Race .
type Race struct {
	Named
	traitor
}

// Spell .
type Spell struct {
	Named
	textual
	Level      int     `xml:"level"`
	School     string  `xml:"school"`
	Time       string  `xml:"time"`
	Range      string  `xml:"range"`
	Ritual     *string `xml:"ritual"`
	Components string  `xml:"components"`
	Duration   string  `xml:"duration"`
	Classes    string  `xml:"classes"`
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

	raceTraits := map[string]*RaceTrait{}
	for _, race := range compendium.Races {
		for _, trait := range race.Traits {
			if t := raceTraits[trait.Name]; t == nil {
				t = &RaceTrait{trait, []string{race.Name}}
				raceTraits[trait.Name] = t
			} else {
				t.Races = append(t.Races, race.Name)
			}
		}
	}
	for _, entity := range raceTraits {
		result = append(result, entity)
	}

	// yikes:
	classFeatures := map[string]*ClassFeature{}
	for _, class := range compendium.Classes {
		for _, level := range class.Levels {
			for _, feature := range level.Features {
				if f := classFeatures[feature.Name]; f == nil {
					f = &ClassFeature{feature, []string{class.Name}}
					classFeatures[feature.Name] = f
				} else if !strings.HasSuffix(class.Name, "!Base") {
					f.Classes = append(f.Classes, class.Name)
				}
			}
		}
	}
	for _, entity := range classFeatures {
		result = append(result, entity)
	}

	return result, nil
}
