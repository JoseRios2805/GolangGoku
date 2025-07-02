package logs

import (
	"RetoTecnicoDragonBall/internal/utils"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type Logs struct {
	ILogger
	ID      string
	Counter int
	Ctx     *gin.Context
}

func NewLog() ILogger {

	return &Logs{}
}

func (l *Logs) SetID(key string) {
	l.ID = key
}
func (l *Logs) GetID() string {
	return l.ID
}
func (l *Logs) GetCounter() int {
	return l.Counter
}
func (l *Logs) SetCounter(counter int) {
	l.Counter = counter
}

func (l *Logs) GetLogging(ctx *gin.Context) Logs {
	c, _ := ctx.Get(utils.C_RequestLog)
	d := c.(Logs)
	d.Ctx.Set(utils.C_Counter, 0)
	return d
}

var (
	infoLevel  = log.New(os.Stdout, "INFO: ", 0) // 0 para no incluir la fecha y hora
	errorLevel = log.New(os.Stderr, "ERROR:", 0) // 0 para no incluir la fecha y hora
	warnLevel  = log.New(os.Stderr, "WARN: ", 0) // 0 para no incluir la fecha y hora
)

func (l *Logs) Info(args ...interface{}) {
	infoLevel.Printf("[%s] %v", l.ID, args)
}

func (l *Logs) Debug(args ...interface{}) {
	infoLevel.Printf("[%s] %v", l.ID, args)
}

func (l *Logs) Error(msg string, err error) {
	errorLevel.Printf("[%s] [%s] [%v]", l.ID, msg, err)
}

func (l *Logs) Warn(msg string, err error) {
	warnLevel.Printf("[%s] [%s] [%v]", l.ID, msg, err)
}
