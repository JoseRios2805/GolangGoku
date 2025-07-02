package usecase

import (
	"RetoTecnicoDragonBall/business-management/application/dto"
	"RetoTecnicoDragonBall/business-management/application/mapper"
	"RetoTecnicoDragonBall/business-management/domain/entity"
	"RetoTecnicoDragonBall/business-management/infrastructure/repository"
	internal "RetoTecnicoDragonBall/internal/http"
	"RetoTecnicoDragonBall/internal/logs"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type IDragonBallUseCase interface {
	SaveCharacter(source *dto.SaveCharacterRequestDTO) (internal.ResponseAPI[*dto.SaveCharacterResponseDTO], int, error)
}

type dragonBallUseCase struct {
	_logger               logs.ILogger
	_dragonBallRepository repository.IDragonBallRepository
	_clientRest           *internal.ClientRest
}

func NewDragonBallUseCase(dragonBallRepository repository.IDragonBallRepository, clientRest *internal.ClientRest, log logs.ILogger) IDragonBallUseCase {
	return &dragonBallUseCase{
		_logger:               log,
		_dragonBallRepository: dragonBallRepository,
		_clientRest:           clientRest,
	}
}

func (d *dragonBallUseCase) SaveCharacter(source *dto.SaveCharacterRequestDTO) (internal.ResponseAPI[*dto.SaveCharacterResponseDTO], int, error) {
	var response internal.ResponseAPI[*dto.SaveCharacterResponseDTO]

	character, errGetByID := d.getCharacterByID(source.Name)
	if errGetByID != nil {
		d._logger.Error("UseCaseDragonBall - UseCaseDragonBall:", errGetByID)
		return GetErrorResponseGeneric[*dto.SaveCharacterResponseDTO](response, errGetByID)
	}
	if character != nil {
		response.Data = mapper.GetDataByNameDBResponseEntityToDTO(character)
		return GetSuccessResponseGeneric[*dto.SaveCharacterResponseDTO](response)
	}

	characterWeb, err := d.getCharacterByName(source.Name)

	if err != nil {
		d._logger.Error("UseCaseDragonBall-UseCaseDragonBall", err)
		return GetErrorResponseGeneric[*dto.SaveCharacterResponseDTO](response, err)
	}

	dataEntity := mapper.InsertDataRequestDtoToEntity(characterWeb)
	if dataEntity == nil {
		errNotFount := errors.New("Not Found Character")
		d._logger.Error("UseCaseDragonBall-UseCaseDragonBall", errNotFount)
		return GetErrorResponseGeneric[*dto.SaveCharacterResponseDTO](response, errNotFount)
	}
	idInsert, err := d._dragonBallRepository.Insert(dataEntity)

	if err != nil {
		d._logger.Error("UseCaseDragonBall - UseCaseDragonBall:", err)
		return GetErrorResponseGeneric[*dto.SaveCharacterResponseDTO](response, err)
	}

	response.Data = mapper.GetDataByNameResponseEntityToDTO(dataEntity, idInsert)

	return GetSuccessResponseGeneric[*dto.SaveCharacterResponseDTO](response)
}

func (d *dragonBallUseCase) getCharacterByName(name string) ([]*entity.GetCharacterByNameResponse, error) {
	urlService := fmt.Sprintf("https://dragonball-api.com/api/characters?name=%s", name)

	result, _, err := d._clientRest.InvokeAPI(urlService, http.MethodGet, nil)

	if err != nil {
		d._logger.Error("InvokeAPI - getCharacterByName:", err)
		return []*entity.GetCharacterByNameResponse{}, err
	}

	response := []*entity.GetCharacterByNameResponse{}
	err = json.Unmarshal(result, &response)
	if err != nil {
		d._logger.Error("Unmarshal - getCharacterByName:", err)
		return []*entity.GetCharacterByNameResponse{}, err
	}

	return response, nil
}

func (d *dragonBallUseCase) getCharacterByID(name string) (*entity.GetCharacterByNameDBResponseEntity, error) {

	result, err := d._dragonBallRepository.SelectByNameInDB(name)
	if err != nil {
		d._logger.Error("UseCaseDragonBall - getCharacterByID:", err)
		return nil, err
	}

	return result, nil
}
