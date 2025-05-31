package scrapper

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func fetchDataToDoc(url string) (*goquery.Document, error) {
	// Fetch the HTML document
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch URL: status code %d", resp.StatusCode)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return doc, nil
}
