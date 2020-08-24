package scrapers

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/spaolacci/murmur3"
	"peoplefind/internal/pkg/models"
	"strconv"
	"strings"
)

const ELASTICSEARCH_INDEX_NAME = "peoplefind"

func elasticSearchBulkPeopleAdd(client *elasticsearch.Client, people []models.Person) error {
	var hashString string = ""

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  ELASTICSEARCH_INDEX_NAME,
		Client: client,
	})

	if err != nil {
		return err
	}

	backgroundContext := context.Background()
	for _, person := range people {
		personJson, err := json.Marshal(person)
		if err != nil {
			return err
		}

		if person.Name == "" {
			continue
		}

		if len(person.Addresses) > 0 {
			hashString = person.Addresses[0].AddressString + person.Name
		} else if len(person.PhoneNumbers) > 0 {
			hashString = person.PhoneNumbers[0].Number + person.Name
		} else {
			hashString = person.Name
		}

		err = bi.Add(backgroundContext, esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: strconv.FormatUint(murmur3.Sum64([]byte(hashString)), 10), // hash address and name
			Body:       strings.NewReader(string(personJson)),
		})
		if err != nil {
			return err
		}
	}

	err = bi.Close(backgroundContext)
	if err != nil {
		return err
	}

	return nil
}
