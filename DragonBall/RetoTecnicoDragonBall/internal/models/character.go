package models

type CharacterModel struct {
	Name        string `json:"name"`
	IdCharacter int    `json:"id_character"`
	Description string `json:"description"`
}
