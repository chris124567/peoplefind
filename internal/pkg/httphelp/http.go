package httphelp

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var EMPTY_HTTP_RESPONSE = &http.Response{}

func HttpPostHeaders(link string, headers map[string]string, data url.Values) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("POST", link, strings.NewReader(data.Encode()))
	if err != nil {
		return EMPTY_HTTP_RESPONSE, err
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, err := client.Do(request)
	if err != nil {
		return EMPTY_HTTP_RESPONSE, err
	}

	return response, nil
}

func HttpGetHeaders(link string, headers map[string]string) (*http.Response, error) {
	log.Print("Requesting " + link)

	client := &http.Client{}
	request, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return EMPTY_HTTP_RESPONSE, err
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, err := client.Do(request)
	if err != nil {
		return EMPTY_HTTP_RESPONSE, err
	}

	return response, nil
}

func HttpGetHeadersWait(link string, headers map[string]string, length time.Duration) (*http.Response, error) {
	time.Sleep(length)

	response, err := HttpGetHeaders(link, headers)
	if err != nil {
		return EMPTY_HTTP_RESPONSE, err
	}

	if response.StatusCode != 200 {
		if response.StatusCode >= 400 && response.StatusCode <= 500 { // if client error, wait (longer) and try again 
			return HttpGetHeadersWait(link, headers, length*2)
		} else {
			return EMPTY_HTTP_RESPONSE, errors.New("Got non 200 response, " + strconv.Itoa(response.StatusCode))
		}
	}

	return response, nil
}
