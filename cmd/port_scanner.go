package cmd

import (
	"fmt"
	cmd "github.com/alekseinovikov/tool/tool"
	"github.com/rivo/tview"
	"time"
)

type portScanner struct {
}

func (p *portScanner) Run(app *tview.Application, main tview.Primitive) {
	var cancelSender chan<- interface{}
	closeChan := make(chan interface{}, 1)
	go func() {
		if _, ok := <-closeChan; ok && cancelSender != nil {
			cancelSender <- true
		}
	}()

	form := tview.NewForm()

	textView := tview.NewTextView()
	textView.SetBorder(true).
		SetTitle("Results")

	flex := tview.NewFlex().
		AddItem(form, 0, 1, true).
		AddItem(textView, 0, 1, false)

	var insertedHost string
	scanFunc := func() {
		go func() {
			textView.Clear()
			results, cancelSender := cmd.NewPortScan(insertedHost).Scan()
			defer func() {
				close(cancelSender)
				cancelSender = nil
			}()

			addHelloOnTextView(textView, app)
			var found bool
			for result := range results {
				if result.Open {
					addFoundPortOnTextView(result.Port, textView, app)
					found = true
				}
			}

			if !found {
				addNotFoundOnTextView(textView, app)
			}

			addByeOnTextView(textView, app)
		}()
	}

	form.AddInputField("Host to scan", "", 20, nil, func(inserted string) {
		insertedHost = inserted
	}).
		AddButton("Scan", scanFunc).
		AddButton("Quit", func() {
			exit(app, main, closeChan)
		})
	form.SetBorder(true).SetTitle("Enter host to scan and press 'Scan'")

	app.SetRoot(flex, true).SetFocus(form)
}

func NewPortScanner() Command {
	return &portScanner{}
}

func addFoundPortOnTextView(port uint16, view *tview.TextView, app *tview.Application) {
	fmt.Fprintf(view, "Found open port: %d\n", port)
	app.Draw()
	time.Sleep(time.Millisecond * 100)
}

func addHelloOnTextView(view *tview.TextView, app *tview.Application) {
	fmt.Fprintf(view, "Started scanning...\n")
	app.Draw()
	time.Sleep(time.Millisecond * 100)
}

func addByeOnTextView(view *tview.TextView, app *tview.Application) {
	fmt.Fprintf(view, "Finished!\n")
	app.Draw()
	time.Sleep(time.Millisecond * 100)
}

func addNotFoundOnTextView(view *tview.TextView, app *tview.Application) {
	fmt.Fprintf(view, "Open ports not found :(\n")
	app.Draw()
}
