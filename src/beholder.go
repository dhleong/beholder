package beholder

import (
	"sort"
)

// Version is the current version of the app
const Version = "1.6.0"

const queryLimit int = 1024

// QueryListener .
type QueryListener func(query string)

// ResultsListener .
type ResultsListener func([]SearchResult)

// App represents the main app state
type App struct {
	// Callbacks provided by NewApp()
	OnQueryChanged QueryListener

	// required callbacks:
	OnUpdateAvailable func(newVersion string)
	OnResults         ResultsListener
	Quit              func()

	// optional callbacks:
	OnError func(err error)

	// internal state:
	entities []Entity
	loaded   bool
}

type scoredEntity struct {
	Entity
	sequences []*MatchedSequence
	score     float32
}

func (e *scoredEntity) GetEntity() Entity {
	return e.Entity
}

func (e *scoredEntity) GetSequences() []*MatchedSequence {
	return e.sequences
}

// NewApp creates a new App using the given DataSource
func NewApp(dataSource DataSource) (*App, error) {

	app := &App{}

	// loading is pretty fast, but by loading
	// asynchronously the app will feel even snappier
	go func() {
		entities, err := dataSource.GetEntities()
		if err != nil && app.OnError != nil {
			app.OnError(err)
			return
		} else if err != nil {
			panic(err)
		}
		app.entities = entities
		app.loaded = true
	}()

	go func() {
		if newVersion := CheckForUpdates(); newVersion != "" {
			// new version!
			app.OnUpdateAvailable(newVersion)
		}
	}()

	onQuery := func(query string) {
		if !app.loaded {
			// Could we defer until loaded?
			// maybe wait on a chan in a go block?
			return
		}

		if len(query) > 0 {
			app.OnResults(app.Query(query))
		} else {
			app.OnResults([]SearchResult{})
		}
	}

	app.OnQueryChanged = onQuery
	return app, nil
}

// NewAppWithEntities is a Factory that's convenient for testing
func NewAppWithEntities(entities []Entity) *App {
	app, _ := NewApp(NewStaticDataSource(entities))
	return app
}

// Query attempts to find Entities that match the query
func (a *App) Query(query string) []SearchResult {
	qm := NewQueryMatcher(query)
	results := make([]SearchResult, 0, queryLimit)
	for _, e := range a.entities {
		m := qm.Match(e.GetName())
		if m.Matched {
			results = append(results, &scoredEntity{
				e,
				m.Sequences,
				m.Score,
			})

			if len(results) >= queryLimit {
				break
			}
		}
	}

	// score results
	sort.Slice(results, func(i, j int) bool {
		iScore := results[i].(*scoredEntity).score
		jScore := results[j].(*scoredEntity).score
		return iScore > jScore
	})

	return results
}
