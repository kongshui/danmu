package service

import (
	"encoding/json"
)

func anyToString(m any) string {
	strByte, _ := json.Marshal(m)
	return string(strByte)
}
