package helpers

import (
	"RetoTecnicoDragonBall/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
)

func ShouldBindData(c *gin.Context, s string) string {
	return c.MustGet(s).(string)
}

func GetJsonTimeOut() interface{} {
	obj := gin.H{"data": nil, "isSuccess": false, "isWarning": false, "errorCode": -1, "message": utils.MESSAGE_ERROR_TIMEOUT}
	return obj
}

func SerializeLogsStruct(s interface{}) string {
	str, _ := json.Marshal(s)
	strOut := string(str)
	if strings.Contains(strOut, "\"password\"") {
		re := regexp.MustCompile(`\"password\":\"[^\"]+\"`)
		strOut = re.ReplaceAllString(strOut, "\"password\":\"*****\"")
	}
	return strOut
}

func GetBaseDomain() string {
	environment := "/dev"

	base := fmt.Sprintf("%s/api", environment)
	domainUrl := "%s/dragonBall"
	return fmt.Sprintf(domainUrl, base)
}
