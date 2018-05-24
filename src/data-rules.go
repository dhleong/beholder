package beholder

import "fmt"

type ruleEntity struct {
	Named
	textual
	kind EntityKind
}

// GetKind from Entity interface
func (r ruleEntity) GetKind() EntityKind {
	return r.kind
}

func (r ruleEntity) String() string {
	return fmt.Sprintf("Rule{%s}", r.Name)
}

// ruleParts is a recursive structure representing a Rule and
// its sub-rules. One Entity will be generated for top-most Rule
// (IE the one which is not part of any other rule) including all
// of its sub-rules as "sections," and one for each sub-rule as well.
type ruleParts struct {
	name  string
	parts []interface{}
}

func rule(name string, parts ...interface{}) *ruleParts {
	return &ruleParts{
		name,
		parts,
	}
}

// section is a semantic alias for rule
var section = rule

// Generate entities given ruleParts, returning the first
// ruleEntity and a slice of *all* entities generated (*including*
// the first one)
func generateEntities(
	kind EntityKind,
	part *ruleParts,
	isTop bool,
	ignoredSections map[string]bool,
) (ruleEntity, []Entity) {
	// NOTE: the top-level entity is always a RuleEntity
	var actualKind EntityKind
	if isTop {
		actualKind = RuleEntity
	} else {
		actualKind = kind
	}

	text := make([]string, 0, len(part.parts))
	parent := &ruleEntity{
		Named: Named{part.name},
		kind:  actualKind,
	}
	generated := []Entity{parent}

	for _, child := range part.parts {
		switch v := child.(type) {
		case string:
			text = append(text, v)

		case *ruleParts:
			if _, ok := ignoreSections[v.name]; ok {
				continue
			}

			// NOTE: if we wanted to get fancy, we could replace all the
			// headers fromt he child entity with, eg: h2 -> h3
			childEntity, newDescendants := generateEntities(kind, v, false, ignoreSections)
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

// rulesDataSource generates a DataSource containing a bunch of parts.
// The top-level Entity will always be RuleEntity, but any child entities
// will use the provided `kind`
func rulesDataSource(kind EntityKind, rules ...*ruleParts) DataSource {
	entities := make([]Entity, 0, len(rules))
	ignored := map[string]bool{}

	for _, part := range rules {
		_, generated := generateEntities(kind, part, true, ignored)
		entities = append(entities, generated...)
	}

	return NewStaticDataSource(entities)
}
