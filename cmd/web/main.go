package main

import (
	"log"
	"net/http"
	"os/signal"
	"peoplefind/web/app"
	"syscall"
	"time"
)

func main() {
	log.Print("Starting web server...")

	signal.Ignore(syscall.SIGPIPE) // ignore SIGPIPE

	serveMux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:         ":5000",
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Connection", "close") // cloudflare
			// writer.Header().Set("Cache-Control", "no-cache") // dev only, disable in production
			serveMux.ServeHTTP(writer, request)
		}),
	}

	serveMux.HandleFunc("/static/", web.StaticHandler)
	serveMux.HandleFunc("/search", web.GetSearchResults)
	serveMux.HandleFunc("/", web.HomePageHandler)
	serveMux.HandleFunc("/about", web.AboutPageHandler)
	serveMux.HandleFunc("/privacy-policy", web.PrivacyPolicyHandler)
	serveMux.HandleFunc("/terms-of-service", web.TOSHandler)
	serveMux.HandleFunc("/robots.txt", web.RobotsTxtHandler)
	serveMux.HandleFunc("/sitemap.xml", web.SitemapXmlHandler)

	log.Fatal(httpServer.ListenAndServe())
}
