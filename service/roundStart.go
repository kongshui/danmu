package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kongshui/danmu/common"
)

// 同步对局状态,开发者在对局开始时调用，同步对局开始事件；在对局结束时再次调用，同步对局结束事件。
func ksSyncGameStatus(data SyncGameStatusStruct, label string, isFirst bool) error {
	if is_mock {
		return nil
	}
	headers := map[string]string{
		"Content-Type": "application/json;charset=UTF-8",
	}
	// data.AppId = app_id
	// body, err := json.Marshal(data)
	// if err != nil {
	// 	log.Println("SyncGameStatus1 err: ", err)
	// 	return false
	// }
	var (
		sdata any
	)
	if label == "start" {
		sdata = map[string]any{
			"roundId":    data.RoundId,
			"roundType":  "singleGroup",
			"bulletPlay": map[string]any{},
		}
	} else if label == "stop" {
		var (
			group      string = ""
			groupIdAll string
		)
		for _, v := range data.GroupResultList {
			if v.Result == 1 {
				group = v.GroupId
				break
			}
		}
		// groupIdAll = "\""
		for _, v := range groupid_list {
			groupIdAll += v + ","
		}
		groupIdAll = groupIdAll[:len(groupIdAll)-1]
		// groupIdAll += "\""
		sdata = map[string]any{
			"roundId":    data.RoundId,
			"roundType":  "singleGroup",
			"roundGroup": groupIdAll,
			"result": map[string]string{
				"singleGroupRoundResult": group,
			},
			"bulletPlay": map[string]any{},
		}

	}

	sdataByte, err := json.Marshal(sdata)
	if err != nil {
		return fmt.Errorf("ksSyncGameStatus err: %v", err)
	}
	urlPath := KsUrlSet(url_SyncGameStatusUrl)
	if urlPath == "" {
		return fmt.Errorf("ksSyncGameStatusurlSet err: urlPath is nil")
	}
	response, err := common.HttpRespond("POST", urlPath, kuaiShouBindBodyToByte(data.RoomId, "round", label, string(sdataByte)), headers)
	if err != nil {
		return fmt.Errorf("ksSyncGameStatus err: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return fmt.Errorf("ksSyncGameStatus status err: %v", response.Status)
	}
	var (
		request any
	)
	if err := json.NewDecoder(response.Body).Decode(&request); err != nil {
		return fmt.Errorf("ksSyncGameStatus json.NewDecoder err: %v", err)
	}
	if int64(request.(map[string]any)["result"].(float64)) != 1 {
		if request.(map[string]any)["errorMsg"].(string) == "直播已关播" || request.(map[string]any)["errorMsg"].(string) == "无游戏中记录" ||
			request.(map[string]any)["errorMsg"].(string) == "直播不存在" || label == "stop" {
			ziLog.Info(fmt.Sprintf("ksSyncGameStatus 直播已关播, roomId: %v, 用户Id： %v ", data.RoomId, data.AnchorOpenId), debug)
			// 后续清理
			endClean(data.RoomId, data.AnchorOpenId)
			return nil
		}
		if label == "start" && isFirst {
			ziLog.Error(fmt.Sprintf("ksSyncGameStatus result err: %v, roomId: %v, 用户Id： %v ", request, data.RoomId, data.AnchorOpenId), debug)
			if int64(request.(map[string]any)["result"].(float64)) == 221283 {
				if isFirst && rdb.IsExistKey(data.RoomId+"_round") {
					// 获取roundId
					roundId, ok := queryRoomIdToRoundId(data.RoomId)
					if ok {
						ksSyncGameStatus(SyncGameStatusStruct{
							AnchorOpenId:    data.AnchorOpenId,
							AppId:           app_id,
							RoomId:          data.RoomId,
							RoundId:         roundId,
							StartTime:       time.Now().UnixMilli() - 1000,
							EndTime:         time.Now().UnixMilli(),
							Status:          2,
							GroupResultList: []GroupResultList{{GroupId: groupid_list[0], Result: 1}},
						}, "stop", false)
					}
					t := time.NewTimer(3 * time.Second)
					<-t.C
					// 再次同步对局状态
					ksSyncGameStatus(data, label, false)
				}
			}
		}
	}
	return nil
}

// 抖音同步对局状态,开发者在对局开始时调用，同步对局开始事件；在对局结束时再次调用，同步对局结束事件。
func dySyncGameStatus(data SyncGameStatusStruct) error {
	headers := map[string]string{
		"Content-Type": "application/json",
		"X-Token":      accessToken.Token,
	}
	data.AppId = app_id
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("dySyncGameStatus err: %v", err)
	}
	response, err := common.HttpRespond("POST", url_round_sync_status, body, headers)
	if err != nil {
		return fmt.Errorf("dySyncGameStatus err: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return fmt.Errorf("dySyncGameStatus status err: %v", response.Status)
	}
	var (
		request any
	)
	if err := json.NewDecoder(response.Body).Decode(&request); err != nil {
		return fmt.Errorf("dySyncGameStatus json.NewDecoder err: %v", err)
	}
	if int64(request.(map[string]any)["errcode"].(float64)) != 0 {
		return fmt.Errorf("dySyncGameStatus err_no: %v,err: %v", request.(map[string]any)["errcode"], request.(map[string]any)["errmsg"])
	}
	return nil
}
