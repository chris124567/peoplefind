package main

import (
	"github.com/elastic/go-elasticsearch/v7"
	"log"
	"peoplefind/internal/scrapers"
)

func main() {
	log.Print("Starting scraper...")

	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal(err)
	}

	scrapers.StartTelephoneDirectoriesScrape(client)

}
