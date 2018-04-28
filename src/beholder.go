package beholder

import (
	"bufio"
	"os"
	"sort"

	"github.com/mitchellh/go-homedir"
)

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
	OnResults ResultsListener
	Quit      func()

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

// NewApp creates a new App
func NewApp() (*App, error) {

	app := &App{}

	// loading is pretty fast, but by loading
	// asynchronously the app will feel even snappier
	go func() {
		entities, err := loadEntities()
		if err != nil && app.OnError != nil {
			app.OnError(err)
			return
		} else if err != nil {
			panic(err)
		}

		app.entities = entities
		app.loaded = true
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
	return &App{
		entities: entities,
	}
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

func loadEntities() ([]Entity, error) {
	// FIXME TODO: download using DataSource, don't reuse this dir, etc.
	compendiumPath, err := homedir.Expand("~/.config/lacona-dnd/character.xml")
	if err != nil {
		return nil, err
	}

	f, err := os.Open(compendiumPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	return ParseXML(reader)
}
