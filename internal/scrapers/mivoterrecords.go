package scrapers

import (
	"encoding/csv"
	"github.com/elastic/go-elasticsearch/v7"
	"io"
	"log"
	"os"
	"peoplefind/internal/pkg/models"
)

// Available at https://michiganvoters.info/download/20200302/EntireStateVoter.zip
const miVoterRecordsCsvPath string = datadirUrl + "michigan_entire_state_voter.csv"
const miVoterRecordsCsvFieldCount int = 48

func StartMiVoterIngest(client *elasticsearch.Client) {
	var i int64 = 0
	var people []models.Person

	csvfile, err := os.Open(miVoterRecordsCsvPath)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// log.Fatal(err)
			log.Print(err, "skipping")
			continue
		}
		if len(record) < miVoterRecordsCsvFieldCount {
			continue
		}

		person := models.Person{}
		person.Deceased = false
		person.Name = models.TitleFormat(record[1], record[2], record[0])
		person.Addresses = append(person.Addresses, models.PersonAddress{
			AddressString: models.TitleFormat(record[8], record[9], record[10], record[11], record[12], record[13], record[14], record[15], record[16], record[17], record[18]),
		})
		people = append(people, person)

		if (i % 10000) == 0 {
			log.Print(i)
			elasticSearchBulkPeopleAdd(client, people)
			people = []models.Person{}
		}

		i += 1
	}
}
