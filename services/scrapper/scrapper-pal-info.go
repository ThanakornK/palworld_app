package scrapper

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"palworld_tools/models"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ScrapperPalInfo() error {
	// URL of the Game8 Palworld Pals info page
	url := "https://game8.co/games/Palworld/archives/439556"

	// Read existing pals info data or create new slice if file doesn't exist
	var pals []models.Pal
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

	// First find table pal list - target the main Pal table
	doc.Find("table.a-table.flexible-cell").Each(func(i int, s *goquery.Selection) {
		fmt.Printf("Found table %d\n", i)
		// Find each row
		s.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
			id := strings.TrimSpace(row.Find("th").Eq(0).Text())
			name := row.Find("td a").Eq(0).Text()
			// Skip empty rows or non-Pal rows
			if name == "" || id == "" {
				return
			}
			fmt.Printf("Processing Pal: ID=%s, Name=%s\n", id, name)
			
			existingPal := findPalByName(pals, name)
			if existingPal == nil {
				// Create new Pal entry
				imageUrl := getImageFromPalworldWiki(name)
				suitabilities := getSuitabilityCol(row)
				children := getChildrenCol(row)

				pal := models.Pal{
					Id:          id,
					Name:        name,
					ImageUrl:    imageUrl,
					Suitability: suitabilities,
					Children:    children,
				}

				pals = append(pals, pal)
			} else if existingPal.ImageUrl == "" {
				// Update existing Pal with image URL if it's missing
				fmt.Printf("Updating image for existing Pal: %s\n", name)
				existingPal.ImageUrl = getImageFromPalworldWiki(name)
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

		if len(matches1) == 0 || len(matches2) == 0 {
			return pals[i].Id < pals[j].Id
		}

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

func getSuitabilityCol(row *goquery.Selection) []models.Suitability {
	var suitabilities []models.Suitability
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
				suitability := models.Suitability{Work: work, Level: level}
				suitabilities = append(suitabilities, suitability)
			})
		}
	})

	return suitabilities
}

func getChildrenCol(row *goquery.Selection) []models.Child {
	var children []models.Child
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
							childObj := models.Child{
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

func isPalExists(pals []models.Pal, name string) bool {
	for _, pal := range pals {
		if pal.Name == name {
			return true
		}
	}
	return false
}

// findPalByName returns a pointer to the Pal with the given name, or nil if not found
func findPalByName(pals []models.Pal, name string) *models.Pal {
	for i := range pals {
		if pals[i].Name == name {
			return &pals[i]
		}
	}
	return nil
}

// getImageFromPalworldWiki fetches the image URL from the Palworld wiki page for a specific Pal
func getImageFromPalworldWiki(palName string) string {
	// Construct the wiki URL using the Pal name
	wikiURL := fmt.Sprintf("https://palworld.wiki.gg/wiki/%s", strings.ReplaceAll(palName, " ", "_"))
	
	fmt.Printf("Fetching image from: %s\n", wikiURL)
	
	// Add a small delay to be respectful to the server
	time.Sleep(1 * time.Second)
	
	// Fetch the wiki page
	doc, err := fetchDataToDoc(wikiURL)
	if err != nil {
		fmt.Printf("Error fetching wiki page for %s: %v (trying alternative naming)\n", palName, err)
		// Try alternative naming patterns for special variants
		var alternateURL string
		if strings.Contains(palName, "Special") {
			// Try format: BaseName_(Special)
			baseName := strings.ReplaceAll(palName, "Special ", "")
			baseName = strings.ReplaceAll(baseName, " Special", "")
			alternateURL = fmt.Sprintf("https://palworld.wiki.gg/wiki/%s_(Special)", strings.ReplaceAll(baseName, " ", "_"))
		} else if strings.Contains(palName, " Lux") {
			// Handle Lux variants
			alternateURL = wikiURL // Keep original for now
		} else if strings.Contains(palName, " Cryst") {
			// Handle Cryst variants
			alternateURL = wikiURL // Keep original for now
		} else if strings.Contains(palName, " Ignis") {
			// Handle Ignis variants
			alternateURL = wikiURL // Keep original for now
		} else {
			// Try removing spaces and special characters
			cleanName := strings.ReplaceAll(palName, " ", "")
			alternateURL = fmt.Sprintf("https://palworld.wiki.gg/wiki/%s", cleanName)
		}
		
		if alternateURL != wikiURL {
			fmt.Printf("Trying alternate URL: %s\n", alternateURL)
			doc, err = fetchDataToDoc(alternateURL)
			if err != nil {
				fmt.Printf("Alternative URL also failed for %s: %v\n", palName, err)
				return ""
			}
		} else {
			return ""
		}
	}
	
	// Look for the main Pal image in the infobox or main content area
	// Try multiple selectors to find the image
	var imageUrl string
	
	// Try infobox image first
	doc.Find(".infobox img, .portable-infobox img").Each(func(i int, img *goquery.Selection) {
		if imageUrl == "" {
			if src, exists := img.Attr("src"); exists {
				// Skip small icons and thumbnails
				if !strings.Contains(src, "thumb") || strings.Contains(src, "150px") {
					imageUrl = src
				}
			}
		}
	})
	
	// If no infobox image, try other image selectors
	if imageUrl == "" {
		doc.Find("img").Each(func(i int, img *goquery.Selection) {
			if imageUrl == "" {
				if src, exists := img.Attr("src"); exists {
					alt, _ := img.Attr("alt")
					// Look for images that likely represent the Pal
					if strings.Contains(strings.ToLower(alt), strings.ToLower(palName)) {
						imageUrl = src
					}
				}
			}
		})
	}
	
	// Convert relative URLs to absolute URLs
	if imageUrl != "" && strings.HasPrefix(imageUrl, "/") {
		imageUrl = "https://palworld.wiki.gg" + imageUrl
	}
	
	fmt.Printf("Found image URL for %s: %s\n", palName, imageUrl)
	return imageUrl
}
