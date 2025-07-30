package service

import (
	"encoding/json"
	"time"

	"github.com/kongshui/danmu/common"
)

// 快手body转byte
func kuaiShouBindBodyToByte(roomid, moduleType, label, data string) []byte {
	bodyStruct := map[string]any{
		"roomCode":   roomid,
		"timestamp":  time.Now().UnixMilli(),
		"moduleType": moduleType,
		"actionType": label,
		"data":       data,
	}
	ms5Str := common.KSSignature(bodyStruct, app_secret, app_id)
	bodyStruct["sign"] = ms5Str
	body, err := json.Marshal(bodyStruct)
	if err != nil {
		return nil
	}
	// bodyJson, err := json.Marshal(body)
	// if err != nil {
	// 	log.Println("json.Marshal err: ", err)
	// 	return nil
	// }
	return body
}
