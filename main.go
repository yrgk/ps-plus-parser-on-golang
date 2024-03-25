package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Game struct {
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

// begin: 13.24

func main() {
	host, _ := os.LookupEnv("HOST")
	user, _ := os.LookupEnv("USER")
	dbname, _ := os.LookupEnv("NAME")
	password, _ := os.LookupEnv("PASSWORD")
	port, _ := os.LookupEnv("PORT")

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s", host, user, dbname, password, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Game{})

	for idx, game := range getAllNames() {
		data := getOneItem(game)
		db.Create(&Game{
			Name: data.Name,
			Price: data.Price,
			CoverUrl: data.CoverUrl,
			Description: data.Description,
			Publisher: data.Publisher,
		})

		fmt.Println(idx, data.Name)
	}
}

func getOneItem(link string) Game {
	fullLink := fmt.Sprintf("https://store.playstation.com/en-us/concept/%s", link)
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

	if game.Name == "" {
		game.Name = doc.Find("h1.game-title").Text()
	}

	game.Price = doc.Find("span[data-qa='mfeCtaMain#offer0#finalPrice']").First().Text()

	if game.Price == "Free" {
		game.Price = doc.Find("span[data-qa='mfeCtaMain#offer1#finalPrice']").First().Text()

		if game.Price == "" {
			game.Price = "Free"
		}
	}

	if game.Price == "" {
		game.Price = "Not available for purchase"
	}

	game.Description = doc.Find("p[data-qa='mfe-game-overview#description']").First().Text()

	if game.Description == "" {
		game.Description = doc.Find("div.text-block").Find("p").First().Text()
	}

	game.Publisher = doc.Find("div[data-qa='mfe-game-title#publisher']").First().Text()
	if game.Publisher == "" {
		game.Publisher = doc.Find("div.publisher").First().Text()
	}

	script := doc.Find("div.pdp-background-image").Children().First().Text()
	jsonParsed, err := gabs.ParseJSON([]byte(script))
	if err != nil {
		game.CoverUrl = ""
		return game
	} else {
		concept := fmt.Sprintf("Concept:%s", link)
		path := fmt.Sprintf("cache.%s.media", concept)
		images, _ := jsonParsed.Path(path).Children()
		game.CoverUrl = images[len(images)-1].Path("url").String()
		return game
	}

}

func getAllNames() []string {
	res, err := http.Get("https://www.playstation.com/en-us/ps-plus/games")

	if err != nil {
		log.Fatal("Failed to load page", err)
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