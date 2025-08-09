package datamanage

import (
	"encoding/json"
	"os"
	"strings"
)

func RemovePal(palName string, id int) error {

	// read stored pal
	pals, err := ReadStoredPals()
	if err != nil {
		return err
	}

	// remove pal
	for i, pal := range pals {
		if strings.EqualFold(pal.Name, palName) {
			for j, storedPal := range pal.StoredPals {
				if storedPal.ID == id {
					pals[i].StoredPals = append(pals[i].StoredPals[:j], pals[i].StoredPals[j+1:]...)
					// update id
					for k := j; k < len(pals[i].StoredPals); k++ {
						pals[i].StoredPals[k].ID--
					}
					break

				}
			}
			if len(pals[i].StoredPals) == 0 {
				pals = append(pals[:i], pals[i+1:]...)
				i--
			}
		}
	}

	// update file
	// Convert the slice to JSON
	jsonData, err := json.MarshalIndent(pals, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON data to a file
	err = os.WriteFile("./data/stored_pals.json", jsonData, 0644)
	if err != nil {
		return err
	}

	return nil

}
