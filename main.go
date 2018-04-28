package main

import (
	"github.com/rivo/tview"

	beholder "github.com/dhleong/beholder/src"
	"github.com/dhleong/beholder/src/ui"
)

func main() {
	app := beholder.NewApp()

	root := ui.NewMainUI(app)
	tapp := tview.NewApplication().
		SetRoot(root, true)

	// app.OnResults = func(results []beholder.Entity) {
	// 	tapp.Stop()
	// }

	if err := tapp.Run(); err != nil {
		panic(err)
	}
}
