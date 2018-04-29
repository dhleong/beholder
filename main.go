package main

import (
	"github.com/rivo/tview"

	beholder "github.com/dhleong/beholder/src"
	"github.com/dhleong/beholder/src/ui"
)

func main() {

	dataSource, err := beholder.NewDataSource()
	if err != nil {
		panic(err)
	}

	app, err := beholder.NewApp(dataSource)
	if err != nil {
		panic(err)
	}

	tapp := tview.NewApplication()
	root := ui.NewMainUI(app, tapp)
	tapp.SetRoot(root, true)

	app.Quit = tapp.Stop

	if err := tapp.Run(); err != nil {
		panic(err)
	}
}
