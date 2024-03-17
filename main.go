package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Game struct {
	name string
	price string
	coverUrl string
	description string
	publisher string
}

func main(){
	for _, game := range getAllNames() {
		fmt.Println(getOneItem(game))
	}
	// fmt.Println(getOneItem("https://store.playstation.com/en-us/concept/10000886"))
	// fmt.Println(getOneItem("https://store.playstation.com/en-us/concept/232654"))
	// fmt.Println(getOneItem("https://store.playstation.com/en-us/concept/231096"))
}


func getOneItem(link string) Game {
	res, err := http.Get(link)

	if err != nil {
		log.Fatal("Failed to parse document", err)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Fatal("Failed to parse document", err)
	}
	var game Game

	game.name = doc.Find("h1[data-qa='mfe-game-title#name']").First().Text()
	game.description = doc.Find("p[data-qa='mfe-game-overview#description']").First().Text()

	return game
}


func getAllNames() []string {
	res, err := http.Get("https://www.playstation.com/en-us/ps-plus/games")

	if err != nil {
		log.Fatal("Failed to parse document", err)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Fatal("Failed to parse document", err)
	}

	var nameList []string

	doc.Find("a[module-name='PS Plus Games List']").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")

		if exists {
			nameList = append(nameList, link)
		}
	})
	return nameList
}