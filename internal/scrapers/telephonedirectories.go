package scrapers

import (
	"github.com/antchfx/htmlquery"
	"github.com/dongri/phonenumber"
	"github.com/elastic/go-elasticsearch/v7"
	"log"
	"math/rand"
	"peoplefind/internal/pkg/htmlhelp"
	"peoplefind/internal/pkg/httphelp"
	"peoplefind/internal/pkg/models"
	"strings"
	"time"
)

const TD_RAND_WAIT_TIME int = 5
const telephoneDirectoriesRootUrl string = "https://www.telephonedirectories.us"

func StartTelephoneDirectoriesScrape(client *elasticsearch.Client) {
	stateUrlList, err := getStateUrlList()
	if err != nil {
		log.Fatal(err)
	}

	for _, stateUrl := range stateUrlList {
		countyUrlList, err := getStateCountyUrlList(stateUrl)
		if err != nil {
			log.Fatal(err)
		}
		for _, countyUrl := range countyUrlList {
			peopleList, err := getCountyPeopleInfo([]models.Person{}, countyUrl)
			if err != nil {
				// log.Fatal(err)
				continue // if single loop iteration doesn't work, just keep going
			}
			// add to elasticsearch after every iteration (should be a 5-20k records at a time) so we don't run out of memory
			elasticSearchBulkPeopleAdd(client, peopleList)
		}
	}
}

func getStateUrlList() ([]string, error) {
	stateUrlList := []string{}

	result, err := httphelp.HttpGetHeadersWait(telephoneDirectoriesRootUrl, httphelp.STANDARD_HEADERS, time.Duration(rand.Intn(TD_RAND_WAIT_TIME))*time.Second)
	if err != nil {
		return stateUrlList, err
	}
	defer result.Body.Close()

	doc, err := htmlquery.Parse(result.Body)
	if err != nil {
		return stateUrlList, err
	}

	stateNodes := htmlquery.Find(doc, `//a[@href][@title]/@href`)
	for _, stateNode := range stateNodes {
		stateUrlList = append(stateUrlList, htmlquery.InnerText(stateNode))
	}

	return stateUrlList, nil
}

func getStateCountyUrlList(stateUrl string) ([]string, error) {
	stateCountyUrlList := []string{}

	result, err := httphelp.HttpGetHeadersWait(telephoneDirectoriesRootUrl+stateUrl, httphelp.STANDARD_HEADERS, time.Duration(rand.Intn(TD_RAND_WAIT_TIME))*time.Second)
	if err != nil {
		return stateCountyUrlList, err
	}
	defer result.Body.Close()

	doc, err := htmlquery.Parse(result.Body)
	if err != nil {
		return stateCountyUrlList, err
	}

	countyNodes := htmlquery.Find(doc, `//table//tbody//tr[@role="row"]//td//a[@href]/@href`)
	for _, countyNode := range countyNodes {
		stateCountyUrlList = append(stateCountyUrlList, htmlquery.InnerText(countyNode))
	}

	return stateCountyUrlList, nil
}

func getCountyPeopleInfo(existingList []models.Person, countyUrl string) ([]models.Person, error) {
	result, err := httphelp.HttpGetHeadersWait(telephoneDirectoriesRootUrl+countyUrl, httphelp.STANDARD_HEADERS, time.Duration(rand.Intn(TD_RAND_WAIT_TIME))*time.Second)
	if err != nil {
		return existingList, err
	}
	defer result.Body.Close()

	doc, err := htmlquery.Parse(result.Body)
	if err != nil {
		return existingList, err
	}

	currentCounty := htmlhelp.GetXpathValue(doc, `//meta[@property="og:title"][@content]/@content`)
	currentCounty = strings.ReplaceAll(currentCounty, "▷ Telephone Directory of ", "")
	currentCounty = strings.ReplaceAll(currentCounty, ".", "")

	peopleOrBusinessList := htmlquery.Find(doc, `//div[contains(@id, "directory")]//li`)
	for _, personOrBusiness := range peopleOrBusinessList {
		person := models.Person{}
		person.Name = htmlhelp.GetXpathValue(personOrBusiness, `//span//em/text()`)

		personAddress := models.PersonAddress{}
		personAddress.Primary = true // all the addresses on the site are (or should be) primary
		personAddress.AddressString = htmlhelp.CleanString(htmlhelp.GetXpathValue(personOrBusiness, `//span//i/text()`)) + ", " + currentCounty
		person.Addresses = append(person.Addresses, personAddress)

		phoneNumber := models.PhoneNumber{}
		phoneNumber.Primary = true // phones on site are (or should be) primary
		phoneNumber.Number = phonenumber.ParseWithLandLine(htmlhelp.CleanString(htmlhelp.GetXpathValue(personOrBusiness, `//b/text()`)), "US")
		person.PhoneNumbers = append(person.PhoneNumbers, phoneNumber)

		existingList = append(existingList, person)
	}

	nextPageNode := htmlquery.FindOne(doc, `//a[@title="Go next"][@href]/@href`)
	if nextPageNode != nil {
		return getCountyPeopleInfo(existingList, htmlquery.InnerText(nextPageNode))
	} else {
		return existingList, nil
	}

}
