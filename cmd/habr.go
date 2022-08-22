package cmd

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/toqueteos/webbrowser"
	"log"
	"net/http"
)

const ROOT = "https://habr.com"
const ArticlesUrl = ROOT + "/ru/all/"
const LinksSelector = ".tm-article-snippet__title-link"

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

	form := tview.NewForm().
		SetItemPadding(1)

	checked := buildMenu(form, items)
	form.AddButton("Open", func() {
		openChecked(checked)
		exit(app, main)
	}).AddButton("Cancel", func() {
		exit(app, main)
	}).SetCancelFunc(func() {
		exit(app, main)
	}).SetInputCapture(bindKeys(app, main, checked)).
		SetBorder(true).
		SetTitle("Check articles to open an press 'o'")

	app.SetRoot(form, true).SetFocus(form)
}

func bindKeys(app *tview.Application, main tview.Primitive, checked map[linkItem]bool) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q':
			exit(app, main)
			return nil
		case 'o':
			openChecked(checked)
			exit(app, main)
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

func exit(app *tview.Application, main tview.Primitive) *tview.Application {
	return app.SetRoot(main, true).SetFocus(main)
}

func buildMenu(form *tview.Form, items []linkItem) map[linkItem]bool {
	var checked = make(map[linkItem]bool, 0)
	for _, item := range items {
		checked[item] = false
		checkBox := tview.NewCheckbox().
			SetChecked(false).
			SetChangedFunc(addCheckedCallback(item, checked))
		checkBox.
			SetLabel(item.title).
			SetBackgroundColor(tcell.ColorGrey)
		form.AddFormItem(checkBox)
	}

	return checked
}

func addCheckedCallback(item linkItem, checkedMap map[linkItem]bool) func(bool2 bool) {
	return func(checked bool) {
		checkedMap[item] = checked
	}
}

func openChecked(checked map[linkItem]bool) {
	for item, selected := range checked {
		if selected {
			openArticle(item.link)
		}
	}
}

func openArticle(link string) {
	webbrowser.Open(link)
}

func getAndParseLinkItems() (result []linkItem, found bool) {
	res, err := http.Get(ArticlesUrl)
	if err != nil {
		log.Fatal("Can't fetch " + ArticlesUrl + err.Error())
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

	counter := 0
	doc.Find(LinksSelector).Each(func(i int, selection *goquery.Selection) {
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
