package logs

import "github.com/gin-gonic/gin"

type ILogger interface {
	Info(args ...interface{})
	Debug(args ...interface{})
	Error(msg string, err error)
	Warn(msg string, err error)
	SetID(key string)
	GetID() (id string)
	SetCounter(counter int)
	GetCounter() (counter int)
	GetLogging(ctx *gin.Context) Logs
}
