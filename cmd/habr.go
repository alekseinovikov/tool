package cmd

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/rivo/tview"
	"github.com/toqueteos/webbrowser"
	"log"
	"net/http"
)

const ROOT = "https://habr.com"

var SHORTCUTS = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

type linkItem struct {
	index int
	title string
	link  string
}

func NewHabr() Command {
	return habr{}
}

type habr struct {
}

func (habr) Run(app *tview.Application, main tview.Primitive) {
	items, found := getAndParseLinkItems()
	if !found {
		return
	}

	list := tview.NewList().ShowSecondaryText(false)

	buildMenu(list, items)
	list.AddItem("Quit", "", 'q', func() {
		app.SetRoot(main, true).SetFocus(main)
	})

	app.SetRoot(list, true).SetFocus(list)
}

func buildMenu(list *tview.List, items []linkItem) {
	for _, item := range items {
		list.AddItem(item.title, "", SHORTCUTS[item.index], openLinkCallback(item))
	}
}

func openLinkCallback(item linkItem) func() {
	return func() {
		openArticle(item.link)
	}
}

func openArticle(link string) {
	webbrowser.Open(link)
}

func getAndParseLinkItems() (result []linkItem, found bool) {
	res, err := http.Get(ROOT + "/ru/all/")
	if err != nil {
		log.Fatal("Can't fetch habr.com/ru/all/" + err.Error())
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal("Can't parse the response: " + err.Error())
		return
	}

	counter := 1
	doc.Find(".tm-article-snippet__title-link").Each(func(i int, selection *goquery.Selection) {
		title := selection.Find("span").Text()
		link, exists := selection.Attr("href")
		if exists {
			result = append(result, linkItem{
				index: counter,
				title: title,
				link:  ROOT + link,
			})

			counter = counter + 1
		}
	})

	return result, true
}
