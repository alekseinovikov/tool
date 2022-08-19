/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"strconv"
)

const ROOT = "https://habr.com"

type linkItem struct {
	index int
	title string
	link  string
}

func (item *linkItem) printSelf() {
	blue := color.New(color.FgBlue)
	yellow := color.New(color.FgYellow)

	blue.Printf("[%d] ", item.index)
	yellow.Printf("%s\n", item.title)
}

// habrCmd represents the habr command
var habrCmd = &cobra.Command{
	Use:   "habr",
	Short: "Fetches latest posts from habr.com",
	Long:  `Fetches latest posts from habr.com and prints the list with titles. Use parameter of the index to open the article in browser.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			article_num, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatal("Wrong argument! Must be a number!")
				return
			}
			openArticle(article_num)
		} else {
			getAllArticles()
		}
	},
}

func init() {
	rootCmd.AddCommand(habrCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// habrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// habrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getAllArticles() {
	items, found := getAndParseLinkItems()
	if !found {
		return
	}

	for _, item := range items {
		item.printSelf()
	}
}

func openArticle(num int) {
	items, found := getAndParseLinkItems()
	if !found {
		return
	}

	var selectedItem *linkItem = nil
	for _, item := range items {
		if item.index == num {
			selectedItem = &item
			break
		}
	}

	if selectedItem == nil {
		log.Fatal("Wrong article number!")
		return
	}

	browser.OpenURL(selectedItem.link)
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
