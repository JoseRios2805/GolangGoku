package mapper

import (
	"RetoTecnicoDragonBall/business-management/application/dto"
	"RetoTecnicoDragonBall/business-management/domain/entity"
	"RetoTecnicoDragonBall/internal/models"
)

func InsertDataRequestDtoToEntity(source []*entity.GetCharacterByNameResponse) *entity.SaveCharacterRequestDTO {
	if len(source) == 0 {
		return nil
	}

	sourceFirst := source[0]

	return &entity.SaveCharacterRequestDTO{
		CharacterModel: models.CharacterModel{
			Name:        sourceFirst.Name,
			IdCharacter: sourceFirst.ID,
			Description: sourceFirst.Description,
		},
	}
}

func GetDataByNameDBResponseEntityToDTO(source *entity.GetCharacterByNameDBResponseEntity) *dto.SaveCharacterResponseDTO {
	return &dto.SaveCharacterResponseDTO{
		CharacterModel: models.CharacterModel{
			Name:        source.Name,
			IdCharacter: source.IdCharacter,
			Description: source.Description,
		},
		Id: source.ID,
	}
}

func GetDataByNameResponseEntityToDTO(source *entity.SaveCharacterRequestDTO, id int) *dto.SaveCharacterResponseDTO {
	return &dto.SaveCharacterResponseDTO{
		CharacterModel: models.CharacterModel{
			Name:        source.Name,
			IdCharacter: source.IdCharacter,
			Description: source.Description,
		},
		Id: id,
	}
}
