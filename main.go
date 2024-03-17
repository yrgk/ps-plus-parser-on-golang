package main

import (
	"log"
	"net/http"
	"github.com/PuerkitoBio/goquery"
)

func main(){

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