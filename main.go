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
		AddItem("Weather", "Show current weather", 'w', weatherCommand).
		AddItem("Port Scanner", "Scans all ports by host", 's', portScanCommand).
		AddItem("Quit", "Exit from the application", 'q', func() {
			app.Stop()
		})

	list.SetBorder(true)
	mainView = list
	app.SetRoot(mainView, true).SetFocus(mainView)
}

func habrCommand() {
	cmd.NewHabr().Run(app, mainView)
}

func weatherCommand() {
	cmd.NewWeatherCommand().Run(app, mainView)
}

func portScanCommand() {
	cmd.NewPortScanner().Run(app, mainView)
}
