package cmd

import (
	cmd "github.com/alekseinovikov/tool/tool"
	"github.com/rivo/tview"
	"strconv"
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

	var insertedHost string
	scanFunc := func() {
		results, cancelSender := cmd.NewPortScan(insertedHost).Scan()
		defer func() {
			close(cancelSender)
			cancelSender = nil
		}()

		for result := range results {
			if result.Open {
				addFoundPortOnForm(result.Port, form)
				app.Draw()
			}
		}
	}

	form.AddInputField("Host to scan", "", 20, nil, func(inserted string) {
		insertedHost = inserted
	}).
		AddButton("Scan", scanFunc).
		AddButton("Quit", func() {
			exit(app, main, closeChan)
		})
	form.SetBorder(true).SetTitle("Enter some data")

	app.SetRoot(form, true).SetFocus(form)
}

func NewPortScanner() Command {
	return &portScanner{}
}

func addFoundPortOnForm(port uint16, form *tview.Form) {
	portString := strconv.Itoa(int(port))
	form.AddInputField("Found open port:", portString, 20, nil, nil)
}
