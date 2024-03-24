package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
	// "gorm.io/gorm"
	// "gorm.io/driver/postgres"
)

type Game struct {
	// gorm.Model
	Id          uint `gorm:"primarykey"`
	Name        string
	Price       string
	CoverUrl    string
	Description string
	Publisher   string
}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	host, _ := os.LookupEnv("HOST")
	user, _ := os.LookupEnv("USER")
	dbname, _ := os.LookupEnv("NAME")
	password, _ := os.LookupEnv("PASSWORD")
	port, _ := os.LookupEnv("PORT")

	fmt.Println(host, user, dbname, password, port)
	// dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s", host, user, dbname, password, port)
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// if err != nil {
	// 	panic("failed to connect database")
	// }

	// fmt.Println(dsn)
	// db.AutoMigrate(&Game{})

	// for _, game := range getAllNames() {
		// data := getOneItem(game)
		// getOneItem(game)
	// 	db.Create(&Game{
	// 		Name: data.Name,
	// 		Price: data.Price,
	// 		CoverUrl: data.CoverUrl,
	// 		Description: data.Description,
	// 		Publisher: data.Publisher,
	// 	})

	// fmt.Println(idx, data.Name)
	// fmt.Println(idx, data.CoverUrl)
	// 	fmt.Println(idx, game)
	// }
	fmt.Println(getOneItem("202006"))
	// fmt.Println("\n")
	fmt.Println(getOneItem("10000886"))
	// fmt.Println("\n")
	fmt.Println(getOneItem("10000649"))
}

func getOneItem(link string) Game {
// func getOneItem(link string) string {
	// func getOneItem(link string) *gabs.Container {
	fullLink := fmt.Sprintf("https://store.playstation.com/en-us/concept%s", link)
	res, err := http.Get(fullLink)

	if err != nil {
		log.Fatal("Failed to parse document", err)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Fatal("Failed to parse document", err)
	}

	var game Game

	game.Name = doc.Find("h1[data-qa='mfe-game-title#name']").First().Text()
	game.Price = doc.Find("span[data-qa='mfeCtaMain#offer0#finalPrice']").First().Text()

	// CoverUrl := doc.Find("div.pdp-background-image").Children().First().Text()

	// fmt.Println(CoverUrl)

	game.Description = doc.Find("p[data-qa='mfe-game-overview#description']").First().Text()

	if game.Description == "" {
		game.Description = doc.Find("div[data-ol-order-start='1']").First().Text()
	}

	game.Publisher = doc.Find("div[data-qa='mfe-game-title#publisher']").First().Text()
	if game.Publisher == "" {
		game.Publisher = doc.Find("div.publisher").First().Text()
	}
	// return game.CoverUrl
	return ""
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
			link = strings.Split(link, "concept/")[1]
			nameList = append(nameList, link)
		}
	})
	return nameList
}
