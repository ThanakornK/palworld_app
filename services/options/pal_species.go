package options

import "palworld_tools/services/datamanage"

func GetPalSpecies() []string {
	palSpecies, err := datamanage.ReadPaldex()
	if err != nil {
		return nil
	}

	speciesNames := make([]string, 0)
	for _, species := range palSpecies {
		speciesNames = append(speciesNames, species.Name)
	}
	return speciesNames
}
