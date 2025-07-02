package middleware

import (
	internal "RetoTecnicoDragonBall/internal/http"
	"RetoTecnicoDragonBall/internal/logs"
	"RetoTecnicoDragonBall/internal/utils"
	"RetoTecnicoDragonBall/internal/utils/helpers"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

const (
	NameRequestID = "RequestID"
)

func EncryptNoAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logs.NewLog()
		log.SetID(c.Request.Header.Get(NameRequestID))
		var response internal.Response

		objRequest := helpers.ShouldBindData(c, utils.C_Response)
		json.Unmarshal([]byte(objRequest), &response)

		c.JSON(response.Status, response.Data)

		c.Next()
	}
}

func DecryptNoAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logs.NewLog()
		log.SetID(c.Request.Header.Get(NameRequestID))

		var objdata interface{}
		c.ShouldBindJSON(&objdata)
		c.Set(utils.C_Request, helpers.SerializeStruct(objdata))

		c.Next()
	}

}
