package beholder

import (
	"sort"
)

// Version is the current version of the app
const Version = "0.2.0"

const queryLimit int = 128

// QueryListener .
type QueryListener func(query string)

// ResultsListener .
type ResultsListener func([]Entity)

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
	score float32
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
			app.OnResults([]Entity{})
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
func (a *App) Query(query string) []Entity {
	qm := NewQueryMatcher(query)
	results := make([]*scoredEntity, 0, queryLimit)
	for _, e := range a.entities {
		m := qm.Match(e.GetName())
		if m.Matched {
			results = append(results, &scoredEntity{e, m.Score})

			if len(results) >= queryLimit {
				break
			}
		}
	}

	// score results
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	// copy the sorted Entities
	scored := make([]Entity, 0, len(results))
	for _, se := range results {
		scored = append(scored, se.Entity)
	}

	return scored
}
