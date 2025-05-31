package scrapper

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Pal struct {
	Id          string
	Name        string
	ImageUrl    string
	Suitability []Suitability
	Children    []Child
}

type Suitability struct {
	Work  string
	Level int
}

type Child struct {
	Parent string
	Child  string
}

func ScrapperPalInfo() error {
	// URL of the Game8 Palworld Pals info page
	url := "https://game8.co/games/Palworld/archives/439556"

	// Read existing pals info data or create new slice if file doesn't exist
	var pals []Pal
	data, err := os.ReadFile("./data/pals.json")
	if err == nil {
		err = json.Unmarshal(data, &pals)
		if err != nil {
			fmt.Println("Error parsing existing pals.json:", err)
			return err
		}
	}

	// Fetch the HTML doc
	doc, err := fetchDataToDoc(url)
	if err != nil {
		return err
	}

	// First find table pal list
	doc.Find("table.a-table.a-table.flexible-cell.a-table").Each(func(i int, s *goquery.Selection) {
		// Find each row
		s.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
			id := strings.TrimSpace(row.Find("th").Eq(0).Text())
			name := row.Find("td a").Eq(0).Text()
			if !isPalExists(pals, name) {

				// imageUrl, _ := row.Find("td.center img").Attr("src")
				suitabilities := getSuitabilityCol(row)
				children := getChildrenCol(row)

				pal := Pal{
					Id:   id,
					Name: name,
					// ImageUrl:    imageUrl,
					Suitability: suitabilities,
					Children:    children,
				}

				pals = append(pals, pal)
			}

		})
	})

	// Sort the slice by ID
	sort.Slice(pals, func(i, j int) bool {
		less := false
		// split id and "B"
		re := regexp.MustCompile(`(\d+)([A-Z]?)`)
		matches1 := re.FindStringSubmatch(pals[i].Id)
		matches2 := re.FindStringSubmatch(pals[j].Id)

		// Convert the numeric part to an integer
		id1, _ := strconv.Atoi(matches1[1])
		id2, _ := strconv.Atoi(matches2[1])

		// Compare the numeric part
		// return pals[i].Id < pals[j].Id
		if id1 < id2 {
			less = true
		} else if id1 == id2 {
			// if matches1 have B matches2 will be first
			if matches1[2] == "B" {
				less = false
			} else if matches2[2] == "B" {
				less = true
			}
		} else {
			less = false
		}

		return less

	})

	// Convert the slice to JSON
	jsonData, err := json.MarshalIndent(pals, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// // Write the JSON data to a file
	err = os.WriteFile("./data/pals.json", jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nPal data saved to pals.json result is", len(pals))

	return nil
}

func getSuitabilityCol(row *goquery.Selection) []Suitability {
	var suitabilities []Suitability
	row.Find("td").Each(func(i int, col *goquery.Selection) {
		if i == 2 {
			col.Find(".align").Each(func(i int, div *goquery.Selection) {
				// Extract the text content of the div
				text := div.Text()
				// Split the text by the colon to separate the work and level
				parts := strings.Split(text, " ")
				work := parts[0]
				level, _ := strconv.Atoi(parts[2])

				// Create a Suitability struct and append it to the slice
				suitability := Suitability{Work: work, Level: level}
				suitabilities = append(suitabilities, suitability)
			})
		}
	})

	return suitabilities
}

func getChildrenCol(row *goquery.Selection) []Child {
	var children []Child
	row.Find("td").Each(func(i int, col *goquery.Selection) {
		if i == 0 {
			palUrl, _ := col.Find("a").Eq(0).Attr("href")

			// Delay for 2 seconds to be polite with the server and avoid being blocked
			time.Sleep(2 * time.Second)

			// Fetch Pal page
			doc, err := fetchDataToDoc(palUrl)
			if err != nil {
				fmt.Println("Error fetching pal page:", err)
				return
			}

			// Find breed table
			doc.Find("h3.a-header--3").Each(func(i int, s *goquery.Selection) {
				if strings.Contains(s.Text(), "Best Ways to") {

					// Go to breed table
					next := s.Next()
					for {
						if goquery.NodeName(next) == "table" {
							break
						}
						next = next.Next()
						if next == nil {
							log.Fatal("Cannot find breed table")
							return
						}
					}

					// Parse the rows in the table
					next.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
						parent := strings.TrimSpace(row.Find("td").Eq(2).Text())
						child := strings.TrimSpace(row.Find("td").Eq(4).Text())

						if parent != "" && child != "" {
							childObj := Child{
								Parent: parent,
								Child:  child,
							}

							children = append(children, childObj)
						}

					})
				}
			})

		}
	})

	return children
}

func isPalExists(pals []Pal, name string) bool {
	for _, pal := range pals {
		if pal.Name == name {
			return true
		}
	}
	return false
}
