package scrapper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PassiveSkill struct {
	Name   string
	Effect string
	Tier   int
}

func ScrapperPassiveSkill() error {
	// URL of the Game8 Palworld Passive Skills page
	url := "https://game8.co/games/Palworld/archives/439667"

	// Read existing passive skills data or create new slice if file doesn't exist
	var passiveSkills []PassiveSkill
	data, err := os.ReadFile("passive_skills.json")
	if err == nil {
		err = json.Unmarshal(data, &passiveSkills)
		if err != nil {
			fmt.Println("Error parsing existing passive skill data.json:", err)
			return err
		}
	}

	// Fetch the HTML document
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch URL: status code %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the heading with id="hm_1"
	doc.Find("h3.a-header--3").Each(func(i int, s *goquery.Selection) {
		if id, exists := s.Attr("id"); exists && id == "hm_1" {
			fmt.Printf("📘 Found 'All Passive Skills' section from %s\n", url)

			// Go to the next table following the <h3>
			next := s.Next()
			for {
				if goquery.NodeName(next) == "table" {
					break
				}
				next = next.Next()
				if next == nil {
					fmt.Println("⚠️ No table found after the heading.")
					return
				}
			}

			// Parse the rows in the table
			next.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
				name := strings.TrimSpace(row.Find("td").Eq(0).Text())
				if !isPassiveSkillExists(passiveSkills, name) {
					effect := strings.TrimSpace(row.Find("td").Eq(1).Text())
					if effect == "" {
						effect = strings.TrimSpace(row.Find("td").Eq(2).Text())
					}
					tierStr := strings.TrimSpace(row.Find("td").Eq(3).Text())
					// Convert tier string to integer
					tier, _ := strconv.ParseInt(strings.Split(tierStr, " ")[1], 10, 64)

					// Create a PassiveSkill object and append it to the slice
					passiveSkill := PassiveSkill{
						Name:   name,
						Effect: effect,
						Tier:   int(tier),
					}
					passiveSkills = append(passiveSkills, passiveSkill)
				}

			})
		}
	})

	// Sort passive skills by tier
	sort.Slice(passiveSkills, func(i, j int) bool {
		return passiveSkills[i].Tier < passiveSkills[j].Tier
	})

	// Convert the slice to JSON
	jsonData, err := json.MarshalIndent(passiveSkills, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Write the JSON data to a file
	err = os.WriteFile("./data/passive_skills.json", jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Passive skills data saved to passive_skills.json result is", len(passiveSkills))

	return nil

}

func isPassiveSkillExists(passiveSkills []PassiveSkill, pk string) bool {
	for _, passiveSkill := range passiveSkills {
		if passiveSkill.Name == pk {
			return true
		}
	}
	return false
}
