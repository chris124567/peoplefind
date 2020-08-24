package web

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	// "log"
	"peoplefind/internal/pkg/models"
)

const INDEX string = "peoplefind"
const REQUEST_SIZE int = 20

func ElasticSearchNameQuery(es *elasticsearch.Client, offset int, name string) models.ElasticQueryResult {
	var results models.ElasticQueryResult

	var buf bytes.Buffer
	query := map[string]interface{}{
		"size": REQUEST_SIZE,
		"from": offset,
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"Name": map[string]interface{}{
					"query":    name,
					"operator": "and",
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		// log.Fatalf("Error encoding query: %s", err)
		return models.ElasticQueryResult{}
	}

	// Perform the search request.
	response, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(INDEX),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		// log.Fatalf("Error getting response: %s", err)
		return models.ElasticQueryResult{}
	}
	defer response.Body.Close()

	if response.IsError() {
		// var e map[string]interface{}
		// err := json.NewDecoder(response.Body).Decode(&e)
		// if err != nil {
		// 	log.Fatalf("Error parsing the error response body: %s", err)
		// } else {
		// 	// Print the response status and error information.
		// 	log.Fatalf("[%s] %s: %s",
		// 		response.Status(),
		// 		e["error"].(map[string]interface{})["type"],
		// 		e["error"].(map[string]interface{})["reason"],
		// 	)
		// }
		return models.ElasticQueryResult{}

	}

	err = json.NewDecoder(response.Body).Decode(&results)

	if err != nil {
		// log.Fatalf("Error parsing the response body: %s", err)
		return models.ElasticQueryResult{}
	}

	return results

}
