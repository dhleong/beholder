package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/rivo/tview"

	beholder "github.com/dhleong/beholder/src"
	"github.com/dhleong/beholder/src/ui"
)

const beholderVersion = "0.2.0"

type options struct {
	// right now we only have the --version and --help options
}

func parseOptions() *options {
	usage := `beholder: A CLI tool for D&D players

Usage: beholder

That's it! Just run the command and start typing to search for whatever
you're curious about. Use <ctrl-k> and <ctrl-k> to scroll through
results. If the information doesn't fit entirely in the info pane, you
can focus on it by pressing <enter>.

When focused on the info pane, you can use the arrow keys, page up/down,
or j/k to scroll; press <esc> or <backspace> to return focus to the
search bar.

Options:
  -h, --help  Show this screen.
  --version   Show version.`
	args, _ := docopt.ParseArgs(
		usage,
		os.Args[1:],
		fmt.Sprintf("beholder version %s", beholderVersion),
	)

	opts := &options{}
	args.Bind(opts)
	return opts
}

func main() {

	parseOptions()

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
