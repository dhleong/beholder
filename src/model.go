package beholder

// EntityKind is the kind of entity
type EntityKind int

// EntityKinds
const (
	ActionEntity EntityKind = iota + 1
	ConditionEntity
	FeatEntity
	FeatureEntity
	ItemEntity
	MonsterEntity
	ReferenceListEntity
	RuleEntity
	SpellEntity
	TraitEntity
)

// Entity is some renderable datum
type Entity interface {
	GetName() string
	GetKind() EntityKind
}

// SearchResult .
type SearchResult interface {
	GetEntity() Entity
	GetSequences() []*MatchedSequence
}

// ReferenceList is a dynamic entity that points to other Entities
type ReferenceList struct {
	Named
	References []Entity
}

// GetKind from Entity interface
func (s ReferenceList) GetKind() EntityKind {
	return ReferenceListEntity
}

// A CategorizedEntity belongs to a category. This is primarily used
// in conjunction with ReferenceList
type CategorizedEntity interface {
	GetCategory() string
}
