package service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/kongshui/danmu/model/pmsg"
	"github.com/kongshui/danmu/sse"

	"github.com/kongshui/danmu/common"

	"google.golang.org/protobuf/proto"
)

// 快手开始和结束游戏推送任务请求
func ksStartFinishGameInfo(roomId, url, label, uid string, isSend bool) bool {
	// fmt.Println(roomid)
	// if isMock {
	// 	return true
	// 	// roomInfo := &pmsg.AnchorInfoMessage{}
	// 	// roomInfo.RoomId = roomid
	// 	// roomInfo.AnchorOpenId = "_666666666666"
	// 	// roomInfo.NickName = "测试"
	// 	// roomInfo.AvatarUrl = ""
	// 	// dataByte, _ := json.Marshal(roomInfo)
	// 	// if err := pushDownLoadMessage(uint32(pmsg.MessageId_StartBindAck), pmsg.MessageId_StartBindAck.String(), uid, dataByte); err != nil {
	// 	// 	log.Println("pushDownLoadMessage err: ", err)
	// 	// }
	// 	// return true
	// }
	headers := map[string]string{
		"Content-Type": "application/json;charset=UTF-8",
	}
	urlPath := urlSet(url)
	if urlPath == "" {
		ziLog.Error( "startFinishGameInfo urlSet err: urlPath is nil", debug)
		return false
	}
	// fmt.Println(urlPath)
	data := map[string]any{}
	jsonData, _ := json.Marshal(data)
	response, err := common.HttpRespond("POST", urlPath, kuaiShouBindBodyToByte(roomId, "bind", label, string(jsonData)), headers)
	if err != nil {
		ziLog.Error( fmt.Sprintf("startFinishGameInfo response err: %v", err), debug)
		return false
	}
	defer response.Body.Close()
	var (
		request any
	)

	if err := json.NewDecoder(response.Body).Decode(&request); err != nil {
		ziLog.Error( fmt.Sprintf("startFinishGameInfo json.NewDecoder err: %v", err), debug)
		return false
	}

	if response.StatusCode != 200 {
		ziLog.Error( fmt.Sprintf("startFinishGameInfo response err ,statusCode 非200: %v, data: %v", response.StatusCode, request), debug)
		return false
	}
	if int64(request.(map[string]any)["result"].(float64)) != 1 {
		ziLog.Error( fmt.Sprintf("startFinishGameInfo err: %v,roomId: %v", request, roomId), debug)
		return false
	}
	if label == "start" {
		roomInfo := &pmsg.AnchorInfoMessage{}
		roomInfoJson := KsRoomInfoStruct{}
		dataStr := request.(map[string]any)["data"].(string)
		if err := json.Unmarshal([]byte(dataStr), &roomInfoJson); err != nil {
			ziLog.Error( fmt.Sprintf("startFinishGameInfo json.Unmarshal: %v", err), debug)
			return true
		}
		roomInfo.RoomId = roomId
		roomInfo.AnchorOpenId = roomInfoJson.UserId
		roomInfo.AvatarUrl = roomInfoJson.AvatarUrl
		roomInfo.NickName = roomInfoJson.NickName
		if !config.App.IsOnline {
			log.Println("startFinishGameInfo roomInfo: ", roomInfo.RoomId, roomInfo.AnchorOpenId, roomInfo.NickName)
		}
		dataByte, err := proto.Marshal(roomInfo)
		if err != nil {
			ziLog.Error( fmt.Sprintf("startFinishGameInfo proto.Marshal: %v", err), debug)
			return false
		}
		if isSend {
			go userInfoCompareStore(roomInfoJson.UserId, roomInfoJson.NickName, roomInfoJson.AvatarUrl)
			setRoomInfo(uid, roomInfo)
			connect(roomId, roomInfo.AnchorOpenId)
			if err := sse.SseSend(pmsg.MessageId_StartBindAck, []string{uid}, dataByte); err != nil {
				ziLog.Error( fmt.Sprintf("startFinishGameInfo pushDownLoadMessage: %v, roomId: %v", err, roomId), debug)
				return false
			}
		}
	}
	return true
}

// 查询快手推送任务状态
func getKsGameInfo(roomid, url string) int64 {
	headers := map[string]string{
		"Content-Type": "application/json;charset=UTF-8",
	}
	data := map[string]any{}
	jsonData, _ := json.Marshal(data)
	urlPath := urlSet(url)
	if urlPath == "" {
		ziLog.Error( "getGameInfo urlSet err, urlPath is nil", debug)
		return 0
	}
	response, err := common.HttpRespond("POST", urlPath, kuaiShouBindBodyToByte(roomid, "bind", "status", string(jsonData)), headers)
	if err != nil {
		ziLog.Error( fmt.Sprintf("getGameInfo response err: %v", err), debug)
		return 0
	}
	defer response.Body.Close()
	var (
		request any
	)

	if err := json.NewDecoder(response.Body).Decode(&request); err != nil {
		ziLog.Error( fmt.Sprintf("getGameInfo json.NewDecoder err: %v", err), debug)
		return 0
	}
	if response.StatusCode != 200 {
		return 0
	}
	if int64(request.(map[string]any)["result"].(float64)) != 1 {
		ziLog.Error( fmt.Sprintf("getGameInfo get err: %v", request), debug)
		return 0
	}
	dataStr := request.(map[string]any)["data"].(string)
	status := make(map[string]int64)
	if err := json.Unmarshal([]byte(dataStr), &status); err != nil {
		ziLog.Error( fmt.Sprintf("getGameInfo json.Unmarshal err: %v,data: %v", err, request), debug)
		return 0
	}
	return status["status"]
}

// 抖音开始和结束游戏推送任务请求
func dyStartFinishGameInfo(roomid, url, msgType string) bool {
	headers := map[string]string{
		"Content-Type": "application/json",
		"access-token": accessToken.Token,
	}
	body, err := json.Marshal(map[string]any{
		"roomid":   roomid,
		"appid":    app_id,
		"msg_type": msgType,
	})
	if err != nil {
		return false
	}
	response, err := common.HttpRespond("POST", url, body, headers)
	if err != nil {
		return false
	}
	defer response.Body.Close()
	var (
		request any
	)

	if err := json.NewDecoder(response.Body).Decode(&request); err != nil {
		log.Println(err)
		return false
	}

	if response.StatusCode != 200 {
		return false
	}
	return int64(request.(map[string]any)["err_no"].(float64)) == 0
}

// 查询抖音推送任务状态
func getDyGameInfo(roomid, url, msgType string) int64 {
	headers := map[string]string{
		"Content-Type": "application/json",
		"access-token": accessToken.Token,
	}
	query := map[string]string{
		"roomid":   roomid,
		"appid":    app_id,
		"msg_type": msgType,
	}
	checkUrl, err := common.GetUrl(url, query)
	if err != nil {
		return 0
	}
	response, err := common.HttpRespond("GET", checkUrl, nil, headers)
	if err != nil {
		return 0
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return 0
	}
	var (
		request any
	)
	if err := json.NewDecoder(response.Body).Decode(&request); err != nil {
		return 0
	}
	return int64(request.(map[string]any)["data"].(map[string]any)["status"].(float64))
}

// dy开启推送任务,房间Id，主播Id,uid为clientUuid，start是否是开始
func dyStartPushTask(roomId, openId, uid string, start bool) {
	if is_mock {
		if start {
			sse.SseSend(pmsg.MessageId_StartBindAck, []string{uid}, []byte{})
		}
		return
	}
	//开启推送任务
	pubsubNameList := []string{"live_comment", "live_gift", "live_like"}
	if start {
		for _, v := range pubsubNameList {
			if getDyGameInfo(roomId, url_check_push_url, v) != 3 {
				if !dyStartFinishGameInfo(roomId, url_start_push_url, v) {
					for _, vFinshed := range pubsubNameList {
						if getDyGameInfo(roomId, url_stop_push_url, vFinshed) == 3 {
							dyStartFinishGameInfo(roomId, url_start_push_url, vFinshed)
						}
					}
					break
				}
			}
		}
		connect(roomId, openId)
		if err := sse.SseSend(pmsg.MessageId_StartBindAck, []string{uid}, []byte{}); err != nil {
			ziLog.Error( fmt.Sprintf("startFinishGameInfo pushDownLoadMessage: %v, roomId: %v", err, roomId), debug)
		}
		return
	}
	for _, v := range pubsubNameList {
		if getDyGameInfo(roomId, url_check_push_url, v) == 3 {
			if !dyStartFinishGameInfo(roomId, url_stop_push_url, v) {
				ziLog.Error( "stop task Error,roomId: "+roomId, debug)
			}
		}
	}
}
