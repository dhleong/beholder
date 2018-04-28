package beholder

// EntityKind is the kind of entity
type EntityKind int

// EntityKinds
const (
	SpellEntityKind EntityKind = iota << 1
)

// Entity is some renderable datum
type Entity interface {
	GetName() string
	GetKind() EntityKind
}

// SpellEntity .
type SpellEntity struct {
	Name string
}

// GetName .
func (s SpellEntity) GetName() string {
	return s.Name
}

// GetKind .
func (s SpellEntity) GetKind() EntityKind {
	return SpellEntityKind
}
