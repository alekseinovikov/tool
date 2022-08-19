package cmd

import (
	"github.com/rivo/tview"
)

func NewHello() Command {
	return hello{}
}

type hello struct {
}

func (hello) Run(app *tview.Application, main tview.Primitive) {
	modal := tview.NewModal().
		SetText("Hello! World!").
		AddButtons([]string{"Back", "Quit"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				app.Stop()
			}

			if buttonLabel == "Back" {
				app.SetRoot(main, true).SetFocus(main)
			}
		})

	app.SetRoot(modal, true).SetFocus(modal)
}
