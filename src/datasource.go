package beholder

// DataSource abstracts fetching and loading Entities
type DataSource interface {
	GetEntities() []Entity
}

// NewDataSource creates a new DataSource
func NewDataSource() DataSource {
	return nil
}

type networkDataSource struct {
	compendiumURL string
}
