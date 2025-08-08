package datamanage

import (
	"encoding/json"
	"fmt"
	"os"
	"palworld_tools/models"
)

func AddPal(palName string, palGender string, passiveSkill []string) error {

	fmt.Println("Reading paldex and passive skills")
	pals, err := ReadPaldex()
	if err != nil {
		return err
	}
	passiveSkills, err := ReadPassiveSkills()
	if err != nil {
		return err
	}

	fmt.Println("Validate pal name")
	// validate pal name
	pal := models.FindPal(pals, palName)
	if pal == nil {
		return fmt.Errorf("pal name not found")
	}

	fmt.Println("Reading stored pals")
	// Read existing stored pals data or create new slice if file doesn't exist
	palStore, err := ReadStoredPals()
	if err != nil {
		return err
	}

	speciesStore := models.FindPalSpeciesFromStore(palStore, palName)
	if speciesStore != nil {
		fmt.Println("Pal species exists")
	} else {
		fmt.Println("New pal species")
	}

	fmt.Println("Validate passive skills")
	// validate passive skills
	registerPks := make([]models.PassiveSkill, 0)
	for _, skill := range passiveSkill {
		pks := models.FindPassiveSkill(passiveSkills, skill)
		if pks == nil {
			return fmt.Errorf("passive skill not found")
		}
		registerPks = append(registerPks, *pks)
	}

	skillNames := make([]string, len(registerPks))
	for i, skill := range registerPks {
		skillNames[i] = skill.Name
	}

	if speciesStore == nil {
		palSpecies := &models.PalSpecies{
			Name:       palName,
			StoredPals: []models.StoredPal{{ID: 1, Gender: palGender, PassiveSkills: skillNames}},
		}
		palStore = append(palStore, *palSpecies)
	} else {
		fmt.Println("Add new pal to species")

		for i := range palStore {
			if palStore[i].Name == palName {
				fmt.Println("Add new pal to specie: ", palStore[i].Name)

				// Update the existing species
				palStore[i].StoredPals = append(palStore[i].StoredPals,
					models.StoredPal{ID: len(palStore[i].StoredPals) + 1, Gender: palGender, PassiveSkills: skillNames})

				break
			}
		}
	}

	// Convert the slice to JSON
	jsonData, err := json.MarshalIndent(palStore, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON data to a file
	err = os.WriteFile("./data/stored_pals.json", jsonData, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Passive skills data saved to stored_pals.json result is", len(palStore))

	return nil

}
