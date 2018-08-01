package scrapers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func loadImage(u string) ([]byte, error) {

	imageURL, err := url.Parse(u)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	// By default, the received uri is formatted like "//imgs.xkcd.com/comics/....png"
	// This creates an error, that's why we parse it as "https://imgs.xkcd.com/comics...."
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://%v%v", imageURL.Host, imageURL.RequestURI()),
		nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.89 Safari/537.36")
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}
