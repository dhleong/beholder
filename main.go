package main

import (
	"github.com/rivo/tview"

	beholder "github.com/dhleong/beholder/src"
	"github.com/dhleong/beholder/src/ui"
)

func main() {

	app, err := beholder.NewApp()
	if err != nil {
		panic(err)
	}

	root := ui.NewMainUI(app)
	tapp := tview.NewApplication().
		SetRoot(root, true)

	app.Quit = tapp.Stop

	if err := tapp.Run(); err != nil {
		panic(err)
	}
}
