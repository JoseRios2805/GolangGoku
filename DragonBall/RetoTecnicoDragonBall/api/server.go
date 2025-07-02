package api

import (
	"RetoTecnicoDragonBall/internal/logs"
	"github.com/gin-gonic/gin"
)

type server struct {
	_log logs.ILogger
}

func newServer(log logs.ILogger) *server {
	return &server{_log: log}
}

func (s *server) Start(port string, gin *gin.Engine) {
	gin.Run(port)
}
