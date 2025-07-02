package api

import (
	"RetoTecnicoDragonBall/business-management/application/handler"
	"RetoTecnicoDragonBall/internal/middleware"
	"RetoTecnicoDragonBall/internal/utils/helpers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HandlerDto struct {
	dragonBallHandler handler.IDragonBallHandler
}

// @contact.email  rai.delgado@encora.com
func routes(handlers *HandlerDto) *gin.Engine {

	gin.SetMode(gin.ReleaseMode) //evaluate use var config environment
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	v1 := r.Group(helpers.GetBaseDomain())
	{
		v1.POST("/characters", middleware.DecryptNoAuth(), handlers.dragonBallHandler.SaveCharacter, middleware.EncryptNoAuth())

	}

	return r

}
