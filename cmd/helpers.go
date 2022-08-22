package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func exit(app *tview.Application, main tview.Primitive, quitChains ...chan<- interface{}) *tview.Application {
	for _, quitChain := range quitChains {
		quitChain <- true
	}
	return app.SetRoot(main, true).SetFocus(main)
}

func bindExitKeys(app *tview.Application, main tview.Primitive, quitChains ...chan<- interface{}) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q':
			exit(app, main, quitChains...)
			return nil
		case 'o':
			exit(app, main, quitChains...)
			return nil
		}

		switch {
		case event.Key() == tcell.KeyDown:
			return tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
		case event.Key() == tcell.KeyUp:
			return tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
		}

		return event
	}
}
