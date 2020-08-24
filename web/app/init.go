package web

import (
	"github.com/elastic/go-elasticsearch/v7"
	"log"
)

var elasticSearchClient *elasticsearch.Client

func init() {
	var err error

	elasticSearchClient, err = elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal(err)
	}
}
