package httphelp

const USER_AGENT string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36"

var STANDARD_HEADERS = map[string]string{
	// https://techblog.willshouse.com/2012/01/03/most-common-user-agents/
	"User-Agent": USER_AGENT,
	"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
	// "Accept-Encoding": "gzip, deflate",
	"Accept-Language": "en-US,en;q=0.9",
	"Sec-Fetch-Site":  "Same-Origin",
	"Sec-Fetch-Mode":  "Navigate",
	"Sec-Fetch-User":  "?1",
}
