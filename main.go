package main

import (
	"fmt"
	"github.com/alekseinovikov/tool/cmd"
	"github.com/rivo/tview"
)

var (
	app      *tview.Application // The tview application.
	mainView tview.Primitive
)

func main() {
	app = tview.NewApplication()

	showMainMenu()

	if err := app.Run(); err != nil {
		fmt.Printf("Error running application: %s\n", err)
	}
}

func showMainMenu() {
	list := tview.NewList().
		AddItem("Habr", "Read some habr articles", 'h', habrCommand).
		AddItem("Hello", "Print Hello World", 'w', helloWorld).
		AddItem("Quit", "Exit from the application", 'q', func() {
			app.Stop()
		})

	list.SetBorder(true)
	mainView = list
	app.SetRoot(mainView, true).SetFocus(mainView)
}

func helloWorld() {
	cmd.NewHello().Run(app, mainView)
}

func habrCommand() {
	cmd.NewHabr().Run(app, mainView)
}
