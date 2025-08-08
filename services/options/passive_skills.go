package options

import (
	"palworld_tools/services/datamanage"
)

func GetPassiveSkills() []string {

	passiveSkills, err := datamanage.ReadPassiveSkills()
	if err != nil {
		return nil
	}

	skillNames := make([]string, 0)
	for _, skill := range passiveSkills {
		skillNames = append(skillNames, skill.Name)
	}
	return skillNames
}
