package models

import "strings"

type PassiveSkill struct {
	Name   string
	Effect string
	Tier   int
}

func FindPassiveSkill(passiveSkill []PassiveSkill, skillName string) *PassiveSkill {
	for _, v := range passiveSkill {
		if strings.ToLower(v.Name) == strings.ToLower(skillName) {
			return &v
		}
	}

	return nil
}
