package datamanage

import (
	"encoding/json"
	"fmt"
	"os"
	"palworld_tools/models"
)

func ReadPaldex() ([]models.Pal, error) {
	// Read existing pals info data or create new slice if file doesn't exist
	var pals []models.Pal
	data, err := os.ReadFile("./data/pals.json")
	if err == nil {
		err = json.Unmarshal(data, &pals)
		if err != nil {
			fmt.Println("Error parsing existing pals.json:", err)
			return nil, err
		}
	}

	return pals, nil
}

func ReadPassiveSkills() ([]models.PassiveSkill, error) {
	// Read existing passive skills data or create new slice if file doesn't exist
	var passiveSkills []models.PassiveSkill
	data, err := os.ReadFile("./data/passive_skills.json")
	if err == nil {
		err = json.Unmarshal(data, &passiveSkills)
		if err != nil {
			fmt.Println("Error parsing existing passive_skills.json:", err)
			return nil, err
		}
	}

	return passiveSkills, nil
}

func ReadStoredPals() ([]models.PalSpecies, error) {
	// Read existing pals info data or create new slice if file doesn't exist
	var storePals []models.PalSpecies
	data, err := os.ReadFile("./data/stored_pals.json")
	if err == nil {
		err = json.Unmarshal(data, &storePals)
		if err != nil {
			fmt.Println("Error parsing existing stored_pals.json:", err)
			return nil, err
		}
	}

	return storePals, nil

}
