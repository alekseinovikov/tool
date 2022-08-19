package cmd

import "github.com/rivo/tview"

type Command interface {
	Run(app *tview.Application, main tview.Primitive)
}
