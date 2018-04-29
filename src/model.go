package beholder

// EntityKind is the kind of entity
type EntityKind int

// EntityKinds
const (
	ItemEntity EntityKind = iota + 1
	MonsterEntity
	SpellEntity
	TraitEntity
)

// Entity is some renderable datum
type Entity interface {
	GetName() string
	GetKind() EntityKind
}
