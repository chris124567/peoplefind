/*
BOT DETECTION KEEPS GETTING TRIGGERED.
Determined this was unfeasible.  Would require scraping at least ~50 million records (most likely much more).
At 1000 requests/minutes going 24/7 (assuming free proxies are even that reliable), that would imply:
	50,000,000 / 1000 = 50000 minutes
	50000 minutes / 60 = 833.333333 hours
	833.333333 hours / 24 hours/day = 34.7222222 days of continuous scraping
*/

package scrapers

import (
	"encoding/json"
	"github.com/antchfx/htmlquery"
	"github.com/dongri/phonenumber"
	"log"
	"math/rand"
	"peoplefind/internal/pkg/htmlhelp"
	"peoplefind/internal/pkg/httphelp"
	"peoplefind/internal/pkg/models"
	"regexp"
	"strings"
	"time"
)

const FPS_RAND_WAIT_TIME = 15
const fpsRootUrl string = "https://www.fastpeoplesearch.com"

var SINCE_REGEX = regexp.MustCompile(`Since .* \d{4}`)

func StartFastPeopleSearchScrape() {
	alphabetUrls, err := getAlphabetDirectoryUrls()
	if err != nil {
		log.Fatal(err, "Failed to find alphabetical directory URLs")
	}

	for _, alphabetUrl := range alphabetUrls {
		firstSubDirectoryUrls, err := getSubDirectoryUrls([]string{}, fpsRootUrl+alphabetUrl)
		if err != nil {
			log.Fatal(err, "Failed to find first name layer directory URLs")
		}
		for _, firstSubDirectoryUrl := range firstSubDirectoryUrls {
			secondSubDirectoryUrls, err := getSubDirectoryUrls([]string{}, fpsRootUrl+firstSubDirectoryUrl)
			if err != nil {
				log.Fatal(err, "Failed to find second name layer directory URLs")
			}
			for _, secondSubDirectoryUrl := range secondSubDirectoryUrls {
				thirdSubDirectoryUrls, err := getSubDirectoryUrls([]string{}, fpsRootUrl+secondSubDirectoryUrl)
				if err != nil {
					log.Fatal(err, "Failed to find second name layer directory URLs")
				}
				for _, thirdSubDirectoryUrl := range thirdSubDirectoryUrls {
					peopleUrls, err := getPeopleUrls([]string{}, fpsRootUrl+thirdSubDirectoryUrl)
					if err != nil {
						log.Fatal(err, "Couldn't get people URLs")
					}
					for _, peopleUrl := range peopleUrls {
						person, err := getPeopleInfo(fpsRootUrl + peopleUrl)
						if err != nil {
							log.Fatal(err, "Couldn't get person info")
						}
						jsonBytes, err := json.MarshalIndent(person, "", "    ")
						if err != nil {
							log.Fatal(err, "Couldn't serialize JSON")
						}
						log.Print(string(jsonBytes))
					}
				}
			}
		}
	}

	// DEMO
	// testCases := []string{"https://www.fastpeoplesearch.com/gladston-arnold_id_G-3956511822158245495", "https://www.fastpeoplesearch.com/reta-miller_id_G-6903234668264814379", "https://www.fastpeoplesearch.com/john-smith_id_G-3513550719393192335", "https://www.fastpeoplesearch.com/alexander-aab_id_G6823521758004453108"}
	// testCases := []string{"https://www.fastpeoplesearch.com/ernest-coakley_id_G-4099424036927694009", "https://www.fastpeoplesearch.com/evangeline-santos_id_G5596241005371903465", "https://www.fastpeoplesearch.com/latia-oakley_id_G-7915682135058384874", "https://www.fastpeoplesearch.com/lynn-lambert_id_G-2188804341903744325"}
	// for _, testCase := range testCases {
	// 	person, err := getPeopleInfo(testCase)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	jsonString, err := json.MarshalIndent(person, "", "    ")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	log.Print(string(jsonString))
	// }
}

func getAlphabetDirectoryUrls() ([]string, error) {
	directoryList := []string{}

	result, err := httphelp.HttpGetHeadersWait(fpsRootUrl, httphelp.STANDARD_HEADERS, time.Duration(rand.Intn(FPS_RAND_WAIT_TIME))*time.Second)
	if err != nil {
		return directoryList, err
	}
	defer result.Body.Close()

	doc, err := htmlquery.Parse(result.Body)
	if err != nil {
		return directoryList, err
	}

	nodeList := htmlquery.Find(doc, `//a[@class="directory-letter-link"][contains(@title, "whose last name starts")][@href]/@href`)

	for _, node := range nodeList {
		directoryList = append(directoryList, htmlquery.InnerText(node))
	}

	return directoryList, nil
}

func getSubDirectoryUrls(existingList []string, directoryUrl string) ([]string, error) {
	result, err := httphelp.HttpGetHeadersWait(directoryUrl, httphelp.STANDARD_HEADERS, time.Duration(rand.Intn(FPS_RAND_WAIT_TIME))*time.Second)
	if err != nil {
		return existingList, err
	}
	defer result.Body.Close()

	doc, err := htmlquery.Parse(result.Body)
	if err != nil {
		return existingList, err
	}

	nodeList := htmlquery.Find(doc, `//li[@class="col-sm-12 col-md-6 col-lg-4"]//a[@href][@title]/@href`)

	for _, node := range nodeList {
		existingList = append(existingList, htmlquery.InnerText(node))
	}

	nextPageNode := htmlquery.FindOne(doc, `//a[@href][contains(text(), "Next Page")]/@href`)
	if nextPageNode != nil {
		return getSubDirectoryUrls(existingList, htmlquery.InnerText(nextPageNode))
	}
	return existingList, nil
}

func getPeopleUrls(existingList []string, nameUrl string) ([]string, error) {
	result, err := httphelp.HttpGetHeadersWait(nameUrl, httphelp.STANDARD_HEADERS, time.Duration(rand.Intn(FPS_RAND_WAIT_TIME))*time.Second)
	if err != nil {
		return existingList, err
	}
	defer result.Body.Close()

	doc, err := htmlquery.Parse(result.Body)
	if err != nil {
		return existingList, err
	}

	nodeList := htmlquery.Find(doc, `//a[@class="btn btn-primary link-to-details"][@href]/@href`)

	for _, node := range nodeList {
		existingList = append(existingList, htmlquery.InnerText(node))
	}

	nextPageNode := htmlquery.FindOne(doc, `//a[@href][contains(text(), "Next Page")]/@href`)
	if nextPageNode != nil {
		return getPeopleUrls(existingList, htmlquery.InnerText(nextPageNode))
	}

	return existingList, nil
}

func getPeopleInfo(personUrl string) (models.Person, error) {
	person := models.Person{}

	result, err := httphelp.HttpGetHeadersWait(personUrl, httphelp.STANDARD_HEADERS, time.Duration(rand.Intn(FPS_RAND_WAIT_TIME))*time.Second)
	if err != nil {
		return person, err
	}
	defer result.Body.Close()

	doc, err := htmlquery.Parse(result.Body)
	if err != nil {
		return person, err
	}

	backgroundCheckString := htmlhelp.GetXpathValue(doc, `//div[@id="background_report_section"]`)
	if strings.Contains(backgroundCheckString, " passed away in") {
		person.Deceased = true

		birthDateString := strings.ReplaceAll(htmlhelp.GetStringInBetween(backgroundCheckString, " was born in ", ","), " of", "")
		birthDateParsed, err := time.Parse("January 2006", birthDateString)
		if err == nil {
			person.BirthDate = birthDateParsed
		}

		deathDateString := strings.ReplaceAll(htmlhelp.GetStringInBetween(backgroundCheckString, "passed away in ", "."), " of", "")
		deathDateParsed, err := time.Parse("January 2006", deathDateString)
		if err == nil {
			person.DeathDate = deathDateParsed
		}

	} else {
		birthDateString := strings.ReplaceAll(htmlhelp.GetStringInBetween(backgroundCheckString, " years old and was born in ", "."), " of", "")
		birthDateParsed, err := time.Parse("January 2006", birthDateString)
		if err == nil {
			person.BirthDate = birthDateParsed
		}
	}

	primaryAddressDivNode := htmlquery.FindOne(doc, `//div[@id="current_address_section"]`)
	if primaryAddressDivNode != nil {
		primaryAddress := models.PersonAddress{}
		primaryAddress.Primary = true

		primaryAddress.AddressString = htmlhelp.CleanString(htmlhelp.GetXpathValue(primaryAddressDivNode, `//a[contains(@title, "Search people living at ")]`))

		primaryAddressDateString := (strings.ReplaceAll(SINCE_REGEX.FindString(htmlhelp.GetXpathValue(primaryAddressDivNode, `//div[@class="detail-box-content"]`)), "Since ", ""))
		primaryAddressSinceDateParsed, err := time.Parse("January 2006", primaryAddressDateString)
		if err == nil {
			primaryAddress.Recorded = primaryAddressSinceDateParsed
		}
		person.Addresses = append(person.Addresses, primaryAddress)
	}

	addressNodes := htmlquery.Find(doc, `//dl[@class="col-sm-12 col-md-6"]//dt[@class="address-link"]/ancestor::dl`)
	for _, addressNode := range addressNodes {
		personAddress := models.PersonAddress{}
		personAddress.AddressString = htmlhelp.GetXpathValue(addressNode, `//dt[@class="address-link"]//a[@title][@href]`)

		if personAddress.AddressString == "" { // if we don't get an address, start over
			continue
		}

		recordedTime := strings.ReplaceAll(htmlhelp.GetXpathValue(addressNode, `//dd[contains(text(), "Recorded ")]`), "Recorded ", "")
		recordedTimeParsed, err := time.Parse("January 2006", recordedTime)
		if err == nil {
			personAddress.Recorded = recordedTimeParsed
		}

		personAddress.Primary = false
		person.Addresses = append(person.Addresses, personAddress)
	}

	emailNodes := htmlquery.Find(doc, `//div[@id="email_section"]//h3[@class="col-sm-12 col-md-6"][contains(text(), "@")]`)
	for _, emailNode := range emailNodes {
		person.Emails = append(person.Emails, htmlquery.InnerText(emailNode))
	}

	alsoKnownAsNodes := htmlquery.Find(doc, `//div[@id="aka-links"]//h3[@class="col-sm-12 col-md-6"]`)
	for _, alsoKnownAsNode := range alsoKnownAsNodes {
		person.AlsoKnownAs = append(person.AlsoKnownAs, htmlquery.InnerText(alsoKnownAsNode))
	}

	person.FullName = htmlhelp.GetXpathValue(doc, `//span[@class="fullname"]`)

	phoneNumberNodes := htmlquery.Find(doc, `//dl//a[contains(@title, "Search people associated with the phone number")][@href]/ancestor::dl`)
	for _, phoneNumberNode := range phoneNumberNodes {
		phoneNumber := models.PhoneNumber{}
		phoneNumber.Number = phonenumber.ParseWithLandLine(htmlhelp.GetXpathValue(phoneNumberNode, `//a[contains(@title, "Search people associated with the phone number")][@href]`), "US")
		phoneNumber.Primary = htmlhelp.GetXpathValue(phoneNumberNode, `//span[@class="nowrap"][@style][contains(text(), "Primary")]`) != ""

		ddTags := htmlquery.Find(phoneNumberNode, `//dd`)
		for _, ddTag := range ddTags {
			tagValue := htmlquery.InnerText(ddTag)
			switch {
			case tagValue == "Landline":
				// phoneNumber.Landline = true
				break
			case tagValue == "Wireless":
				// phoneNumber.Landline = false
				break
			case strings.Contains(tagValue, "First reported "):
				reportedDateString := strings.ReplaceAll(tagValue, "First reported ", "")
				reportedDateParsed, err := time.Parse("January 2006", reportedDateString)
				if err == nil {
					phoneNumber.ReportedDate = reportedDateParsed
				}
				break
			}
		}

		person.PhoneNumbers = append(person.PhoneNumbers, phoneNumber)
	}

	relativeNodes := htmlquery.Find(doc, `//dl[@class="col-sm-12 col-md-4"]`)
	for _, relativeNode := range relativeNodes {
		relative := models.Relative{}
		relative.Name = htmlhelp.CleanString(htmlhelp.GetXpathValue(relativeNode, `//a`))

		ddTags := htmlquery.Find(relativeNode, `//dd`)
		for _, ddTag := range ddTags {
			tagValue := htmlquery.InnerText(ddTag)
			if strings.Contains(tagValue, "Age ") {
				birthDateString := htmlhelp.GetStringInBetween(tagValue, "(", ")")
				birthDateParsed, err := time.Parse("Jan 2006", birthDateString)
				if err == nil {
					relative.BirthDate = birthDateParsed
				}
			} else if strings.Contains(tagValue, "Spouse") {
				relative.Spouse = true
			}
		}

		person.Relatives = append(person.Relatives, relative)
	}

	associateNodes := htmlquery.Find(doc, `//dl[@class="col-sm-6 col-md-4"]`)
	for _, associateNode := range associateNodes {
		associate := models.Associate{}
		associate.Name = htmlhelp.CleanString(htmlhelp.GetXpathValue(associateNode, `//a`))

		ageValue := htmlhelp.CleanString(htmlhelp.GetXpathValue(associateNode, `//dd`))
		birthDateString := htmlhelp.GetStringInBetween(ageValue, "(", ")")
		birthDateParsed, err := time.Parse("Jan 2006", birthDateString)

		if err == nil {
			associate.BirthDate = birthDateParsed
		}

		person.Associates = append(person.Associates, associate)
	}

	businessNodes := htmlquery.Find(doc, `//div[@class="detail-box-content"]//div[@class="detail-box-business"]//dt//h3`)
	for _, businessNode := range businessNodes {
		person.Businesses = append(person.Businesses, htmlhelp.CleanString(htmlquery.InnerText(businessNode)))
	}

	return person, nil
}
