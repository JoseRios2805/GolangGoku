package helpers

import (
	"encoding/json"
)

func SerializeStruct(s interface{}) string {
	str, _ := json.Marshal(s)
	return string(str)
}
