package cmd

import (
	"fmt"
	"github.com/rivo/tview"
	"io"
	"log"
	"net/http"
	"time"
)

type weather struct {
}

func (w weather) Run(app *tview.Application, main tview.Primitive) {
	flex := tview.NewFlex()
	quit := make(chan interface{})
	flex.SetBorder(true).
		SetTitle("Weather (press 'q' to exit)").
		SetInputCapture(bindExitKeys(app, main, quit))

	textView := tview.NewTextView().
		SetDynamicColors(true)

	flex.AddItem(textView, 0, 1, false)

	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				updateWeather(app, textView)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	go updateWeather(app, textView)
	app.SetRoot(flex, true).SetFocus(flex)
}

func updateWeather(app *tview.Application, textView *tview.TextView) {
	textView.Clear()
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://wttr.in", nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "curl/7.79.1")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	if err != nil {
		textView.Write([]byte(fmt.Sprintf("Can't reach wttr.in! %v", err)))
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)
	str := string(bytes)
	translated := tview.TranslateANSI(str)

	textView.Write([]byte(translated))

	app.Draw()
}

func NewWeatherCommand() Command {
	return weather{}
}
