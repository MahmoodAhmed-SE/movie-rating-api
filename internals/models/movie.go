package models

type Movie struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	DescriptionVec string `json:"description_vec"`
}
