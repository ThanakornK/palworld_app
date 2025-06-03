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

type ComboPassiveSkill struct {
	Name   string
	Skills []string
}

func BestComboPassiveSkill() error {
	// URL of the Game8 Palworld Best Combo Passive Skills page
	url := "https://game8.co/games/Palworld/archives/440414"

	topicCombo := map[string]interface{}{
		"combat": "hs_4",
		"work":   "hs_5",
		"mount":  "hs_6",
	}

	// Read existing combo passive skills data or create new slice if file doesn't exist
	var comboPks []ComboPassiveSkill

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

	for k, v := range topicCombo {
		comboObj := make([]string, 0)
		comboName := strings.Replace(k, string(k[0]), strings.ToUpper(string(k[0])), 1)
		doc.Find("h4.a-header--4").Each(func(i int, s *goquery.Selection) {
			if id, exists := s.Attr("id"); exists && id == v {

				fmt.Printf("📘 Found 'Best Passive Skill Combos for %s Pals' section from %s\n", comboName, url)

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
					if len(comboObj) < 4 {
						row.Find("td").Each(func(i int, col *goquery.Selection) {
							pkName := strings.TrimSpace(col.Text())
							if pkName != "" {
								comboObj = append(comboObj, pkName)
							}
						})
					}

				})
			}
		})
		comboPks = append(comboPks, ComboPassiveSkill{Name: comboName, Skills: comboObj})
	}

	// Convert the slice to JSON
	jsonData, err := json.MarshalIndent(comboPks, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Write the JSON data to a file
	err = os.WriteFile("./data/passive_skill_combos.json", jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Combo passive skills data saved to passive_skill_combos.json result is", len(comboPks))

	return nil
}
