package models

import "strings"

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

type PalSpecies struct {
	Name       string
	StoredPals []StoredPal
}

type StoredPal struct {
	ID     int
	Gender string

	PassiveSkills []string
}

func FindPal(pals []Pal, palName string) *Pal {
	for _, pal := range pals {
		if strings.ToLower(pal.Name) == strings.ToLower(palName) {
			return &pal
		}
	}
	return nil
}

func FindPalSpeciesFromStore(palSpecies []PalSpecies, palName string) *PalSpecies {
	for _, pal := range palSpecies {
		if strings.ToLower(pal.Name) == strings.ToLower(palName) {
			return &pal
		}
	}
	return nil
}
