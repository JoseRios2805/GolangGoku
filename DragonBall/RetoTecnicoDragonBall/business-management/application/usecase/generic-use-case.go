package usecase

import (
	internal "RetoTecnicoDragonBall/internal/http"
	"RetoTecnicoDragonBall/internal/utils"
	"net/http"
)

func GetErrorResponseGeneric[Res any](response internal.ResponseAPI[Res], err error) (internal.ResponseAPI[Res], int, error) {
	response.Message = err.Error()
	response.ErrorCode = utils.EC_CRITICAL_ERROR
	response.IsSuccess = true
	response.IsWarning = true

	return response, http.StatusOK, err
}

func GetSuccessResponseGeneric[Res any](response internal.ResponseAPI[Res]) (internal.ResponseAPI[Res], int, error) {
	response.IsSuccess = true
	response.IsWarning = false
	response.Message = utils.MESSAGE_SUCCESSFUL

	return response, http.StatusOK, nil
}
