package beholder

// QueryListener .
type QueryListener func(query string)

// ResultsListener .
type ResultsListener func([]Entity)

// App represents the main app state
type App struct {
	OnQueryChanged QueryListener

	OnResults ResultsListener
}

// NewApp creates a new App
func NewApp() *App {

	app := &App{}

	onQuery := func(query string) {
		if len(query) > 0 {
			app.OnResults([]Entity{})
		}
	}

	app.OnQueryChanged = onQuery
	return app
}
