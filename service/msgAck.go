package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kongshui/danmu/common"
)

func ksMsgAckSend(roomid, ackType string, data KsMsgAckReceiveStruct) {
	defer RecoverFunc()
	var (
		url string = url_CpShowAckUrl
	)
	if ackType == "cpClientReceive" {
		url = url_CpReceiveAckUrl
		data.CpClientReceiveTime = time.Now().UnixMilli()
	}
	dataByte, _ := json.Marshal(data)
	sdata := KsAckStruct{}
	sdata.RoomCode = roomid
	sdata.AckType = ackType
	sdata.Data = string(dataByte)
	sdata.TimeStamp = time.Now().UnixMilli()
	var (
		request map[string]any
	)
	sdataByte, _ := json.Marshal(sdata)
	if err := json.Unmarshal(sdataByte, &request); err != nil {
		panic(fmt.Sprintf("KsMsgAckSend json.Unmarshal err: %v,ackType: %v", err, ackType))
	}
	sdata.Sign = common.KSSignature(request, app_secret, app_id)
	requestBody, err := json.Marshal(sdata)
	if err != nil {
		panic(fmt.Sprintf("KsMsgAckSend json.Marshal err: %v,ackType: %v", err, ackType))
	}
	header := map[string]string{
		"Content-Type": "application/json;charset=UTF-8",
	}
	urlPath := KsUrlSet(url)
	if urlPath == "" {
		panic("KsMsgAckSend urlSet err: urlPath is nil")
	}
	response, err := common.HttpRespond("POST", urlPath, requestBody, header)
	if err != nil {
		panic(fmt.Sprintf("KsMsgAckSend HttpRespond err: %v,ackType: %v", err, ackType))
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		panic(fmt.Sprintf("KsMsgAckSend HttpStatus err: %v,ackType: %v", response.StatusCode, ackType))
	}
	request = map[string]any{}
	if err := json.NewDecoder(response.Body).Decode(&request); err != nil {
		panic(fmt.Sprintf("KsMsgAckSend json.NewDecoder err: %v,ackType: %v", err, ackType))
	}
	if int64(request["result"].(float64)) != 1 && int64(request["result"].(float64)) != 220372 {
		panic(fmt.Sprintf("KsMsgAckSend result err，ackType: %v, sendData: %v", ackType, request))
	}
}
