package web

import (
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"peoplefind/internal/pkg/models"
	"strconv"
	textTemplate "text/template"
)

func httpError(writer http.ResponseWriter, statusCode int) {
	http.Error(writer, http.StatusText(statusCode), statusCode)
}

func genericErrorHandler(writer http.ResponseWriter, request *http.Request, statusCode int) {
	writer.WriteHeader(statusCode)
	err := writeTemplateData(writer, request, "./web/template/error.html.tmpl", nil)
	if err != nil {
		log.Fatal(err)
		// httpError(writer, templateError)
		return
	}
}

func StaticHandler(writer http.ResponseWriter, request *http.Request) {
	// Basically just reads file of a given path, i.e. static/main.css
	path := "./web" + request.URL.Path
	file, err := os.Stat(path)
	if err == nil && !file.IsDir() {
		// writer.Header().Set("Vary", "Accept-Encoding")
		// writer.Header().Set("Cache-Control", "public, max-age=7776000")
		http.ServeFile(writer, request, path)
		return
	}

	http.NotFound(writer, request)
}

func writeTemplateData(writer http.ResponseWriter, request *http.Request, path string, data interface{}) error {
	// Initialize a slice containing the paths to the two files. Note that
	// home.page.tmpl file must be the *first* file in the slice.
	files := []string{
		path,
		"./web/template/base.html.tmpl",
	}

	// Use the template.ParseFiles() function to read the files and store the
	// templates in a template set. Notice that we can pass the slice of file paths
	// as a variadic parameter
	templateName := filepath.Base(path)
	template, err := template.New(templateName).ParseFiles(files...)
	if err != nil {
		return err
	}

	// write to template output to client
	err = template.ExecuteTemplate(writer, templateName, &data)
	if err != nil {
		return err
	}

	return nil
}

func writeTemplateDataText(writer http.ResponseWriter, request *http.Request, path string, data interface{}) error {
	// Initialize a slice containing the paths to the two files. Note that
	// home.page.tmpl file must be the *first* file in the slice.
	files := []string{
		path,
		"./web/template/base.html.tmpl",
	}

	// Use the template.ParseFiles() function to read the files and store the
	// templates in a template set. Notice that we can pass the slice of file paths
	// as a variadic parameter
	templateName := filepath.Base(path)
	template, err := textTemplate.New(templateName).ParseFiles(files...)
	if err != nil {
		return err
	}

	// write to template output to client
	err = template.ExecuteTemplate(writer, templateName, &data)
	if err != nil {
		return err
	}

	return nil
}

func GetSearchResults(writer http.ResponseWriter, request *http.Request) {

	offsetString := request.URL.Query().Get("offset")
	offsetInteger, err := strconv.Atoi(offsetString)
	if err != nil {
		offsetInteger = 0
	}

	queryResults := models.SiteSearchResult{}
	queryResults.Query = request.URL.Query().Get("q")
	queryResults.ElasticResult = ElasticSearchNameQuery(elasticSearchClient, offsetInteger, queryResults.Query)
	queryResults.Pagination = models.PaginationResult{}

	resultsCount := queryResults.ElasticResult.Hits.TotalResults.TotalValue
	if offsetInteger <= 0 {
		queryResults.Pagination.IsPreviousPage = false
	} else {
		queryResults.Pagination.IsPreviousPage = true
	}
	if offsetInteger >= (resultsCount-REQUEST_SIZE) || resultsCount < REQUEST_SIZE {
		queryResults.Pagination.IsNextPage = false
	} else {
		queryResults.Pagination.IsNextPage = true
	}

	queryResults.Pagination.Pages = int(math.Ceil(float64(resultsCount) / float64(REQUEST_SIZE)))

	queryResults.Pagination.Offset = offsetInteger
	queryResults.Pagination.PreviousPageOffset = offsetInteger - REQUEST_SIZE
	queryResults.Pagination.NextPageOffset = offsetInteger + REQUEST_SIZE

	err = writeTemplateData(writer, request, "./web/template/search.html.tmpl", queryResults)
	if err != nil {
		log.Fatal(err)
	}
}

func HomePageHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		genericErrorHandler(writer, request, 404)
		return
	}

	err := writeTemplateData(writer, request, "./web/template/home.html.tmpl", struct{}{})
	if err != nil {
		log.Fatal(err)
	}
}

func AboutPageHandler(writer http.ResponseWriter, request *http.Request) {
	err := writeTemplateData(writer, request, "./web/template/about.html.tmpl", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func TOSHandler(writer http.ResponseWriter, request *http.Request) {
	err := writeTemplateData(writer, request, "./web/template/tos.html.tmpl", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func PrivacyPolicyHandler(writer http.ResponseWriter, request *http.Request) {
	err := writeTemplateData(writer, request, "./web/template/privacy.html.tmpl", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func RobotsTxtHandler(writer http.ResponseWriter, request *http.Request) {
	err := writeTemplateData(writer, request, "./web/template/robots.txt.tmpl", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func SitemapXmlHandler(writer http.ResponseWriter, request *http.Request) {
	err := writeTemplateDataText(writer, request, "./web/template/sitemap.xml.tmpl", nil)
	if err != nil {
		log.Fatal(err)
	}
}
