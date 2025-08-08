package scrapper

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func fetchDataToDoc(url string) (*goquery.Document, error) {
	// Fetch the HTML document
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch URL: status code %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
