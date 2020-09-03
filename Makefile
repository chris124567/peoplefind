SCRAPER_OUTPUT_FILE := scraper
WEB_OUTPUT_FILE := webserver
SCRAPER_FILE := cmd/scraper/main.go
WEB_FILE := cmd/web/main.go

# all:
# 	gofmt -s -w .
# 	reset
# 	sudo systemctl start elasticsearch
# 	go run ${SCRAPER_FILE}

all:
	gofmt -s -w .
	reset
	go run ${SCRAPER_FILE}

runscraper:
	gofmt -s -w .
	reset
	go run ${SCRAPER_FILE}

buildscraper:
	gofmt -s -w .
	go build -o ${SCRAPER_OUTPUT_FILE} ${SCRAPER_FILE}

buildweb:
	gofmt -s -w .
	go build -o ${WEB_OUTPUT_FILE} ${WEB_FILE}

clean:
	rm ${MAIN_FILE}