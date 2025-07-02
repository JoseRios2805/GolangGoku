package dto

import "RetoTecnicoDragonBall/internal/models"

type SaveCharacterResponseDTO struct {
	models.CharacterModel
	Id int `json:"id"`
}
