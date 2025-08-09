package dto

type AddPalRequest struct {
	Name          string   `json:"name"`
	Gender        string   `json:"gender"`
	PassiveSkills []string `json:"passive_skills"`
}

type Pal struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
	Gender   string `json:"gender"`

	PassiveSkills []PassiveSkill `json:"passive_skills"`
}

type PassiveSkill struct {
	Name string `json:"name"`
}

type PalSpecies struct {
	Name string `json:"name"`
}

type RemovePalRequest struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}
