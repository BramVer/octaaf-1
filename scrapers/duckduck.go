package scrapers

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const ddgo_uri = "https://duckduckgo.com/lite?k1=-1&q="

// Search searches on duckduck go & returns the first url
func Search(query string, nsfw bool) (string, bool) {
	if len(query) == 0 {
		return "", false
	}

	bang_uri, is_bang := bang(query)
	if is_bang {
		return bang_uri, true
	}

	url := ddgo_uri + query

	if nsfw {
		url += "&kp=-2"
	}

	res, err := fetch(url)

	if err != nil {
		return "", false
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return "", false
	}

	return doc.Find(".result-link").First().Attr("href")
}

// Detect if the query is a bang query, return the duckduckgo redirect
// Eg: !arch amdgpu
func bang(query string) (string, bool) {
	if len(query) > 1 && query[0] == '!' {
		resp, err := http.Get(ddgo_uri + query)
		if err != nil {
			return "", false
		}

		uri := resp.Request.URL.String()
		defer resp.Body.Close()
		return uri, true
	}
	return "", false
}
