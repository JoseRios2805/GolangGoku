package handler

import (
	"RetoTecnicoDragonBall/business-management/application/dto"
	"RetoTecnicoDragonBall/business-management/application/usecase"
	internal "RetoTecnicoDragonBall/internal/http"
	"RetoTecnicoDragonBall/internal/logs"
	"RetoTecnicoDragonBall/internal/utils"
	"RetoTecnicoDragonBall/internal/utils/helpers"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IDragonBallHandler interface {
	SaveCharacter(ctx *gin.Context)
}

type dragonBallService struct {
	_log               logs.ILogger
	_dragonBallUseCase usecase.IDragonBallUseCase
}

func NewDragonBallHandler(dragonBallUseCase usecase.IDragonBallUseCase, log logs.ILogger) IDragonBallHandler {
	return &dragonBallService{
		_log:               log,
		_dragonBallUseCase: dragonBallUseCase,
	}
}

func (d *dragonBallService) SaveCharacter(ctx *gin.Context) {
	var objResponse internal.Response
	var data *dto.SaveCharacterRequestDTO

	objRequest := helpers.ShouldBindData(ctx, utils.C_Request)
	if err := json.Unmarshal([]byte(objRequest), &data); err != nil {
		objResponse.Status = http.StatusBadRequest
		objResponse.Data = gin.H{"message": "Bad request", "code": http.StatusBadRequest}
		ctx.Set(utils.C_Response, helpers.SerializeStruct(objResponse))
		return
	}

	result, status, err := d._dragonBallUseCase.SaveCharacter(data)
	objResponse.Status = status

	if err != nil {
		d._log.Error("dragonBall-handler getCharacterByID - Error :", err)
		objResponse.Data = result
		ctx.Set(utils.C_Response, helpers.SerializeStruct(objResponse))
		return
	}

	objResponse.Data = result
	ctx.Set(utils.C_Response, helpers.SerializeStruct(objResponse))
}
