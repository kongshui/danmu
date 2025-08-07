package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"time"

	battlematchv1 "github.com/kongshui/danmu/battlematch/v1"

	"github.com/kongshui/danmu/common"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 快手连线回调
func TwoConnect(label, roomId, roomId2, pkId string) error {
	var (
		data map[string]any = make(map[string]any)
	)
	headers := map[string]string{
		"Content-Type": "application/json;charset=UTF-8",
	}
	data["roomCode"] = roomId
	data["method"] = label
	data["bizType"] = "chat"
	data["timestamp"] = time.Now().UnixMilli()
	switch label {
	case "config":
		dataMap := map[string]any{
			"version": "1",
			"chatLayoutList": []map[string]int{
				{
					"windowId": 1,
					"x":        200,
					"y":        50,
					"width":    160,
					"height":   160,
				},
				{
					"windowId": 2,
					"x":        720,
					"y":        50,
					"width":    160,
					"height":   160,
				},
			},
		}
		data["data"] = anyToString(dataMap)
	case "start":
		dataMap := map[string]any{
			"cpPkId": pkId,
			"hostAuthor2Layout": map[string]any{
				roomId: map[string]int{
					roomId:  1,
					roomId2: 2,
				},
				roomId2: map[string]int{
					roomId:  1,
					roomId2: 2,
				},
			},
		}
		data["data"] = anyToString(dataMap)
	case "heartbeat", "stop":
		dataMap := map[string]any{
			"cpPkId":    pkId,
			"roomCodes": []string{roomId, roomId2},
			"timestamp": time.Now().UnixMilli(),
		}
		data["data"] = anyToString(dataMap)
	}
	data["sign"] = common.KSSignature(data, app_secret, app_id)
	urlPath := KsUrlSet(url_ChatUrl)
	if urlPath == "" {
		return errors.New("TwoConnect urlSet err: urlPath is nil")
	}
	jsonData, err := json.Marshal(data)
	// fmt.Println("TwoConnect jsonData: ", string(jsonData), "type: ", label)
	if err != nil {
		return errors.New("TwoConnect json.Marshal err: " + err.Error() + ", data: " + anyToString(data))
	}
	response, err := common.HttpRespond("POST", urlPath, jsonData, headers)
	if err != nil {
		return errors.New("TwoConnect HttpRespond err: " + err.Error())
	}
	defer response.Body.Close()
	var (
		request any
	)

	if err := json.NewDecoder(response.Body).Decode(&request); err != nil {
		return errors.New("TwoConnect json.NewDecoder err: " + err.Error())
	}
	if response.StatusCode != 200 {
		return errors.New("TwoConnect response err,statusCode 非200: " + strconv.Itoa(response.StatusCode))
	}
	if int64(request.(map[string]any)["result"].(float64)) != 1 {
		return errors.New("TwoConnect err: " + anyToString(request) + "roomId:" + roomId + roomId2)
	}
	return nil
}

func matchV1HeardBeat() {
	t := time.NewTicker(6 * time.Second)
	defer t.Stop()
	// var (
	// 	roomId1 string
	// 	roomId2 string
	// )
	for {
		<-t.C
		res, err := etcdClient.Client.Get(first_ctx, path.Join("/", config.Project, match_battle_group_time), clientv3.WithPrefix())
		if err != nil {
			ziLog.Error(fmt.Sprintf("MatchV1HeardBeat etcdClient.Client.Get err: %v", err), debug)
			continue
		}
		// log.Println("matchV1HeardBeat :", res.Count)
		if res.Count == 0 {
			continue
		}
		for _, v := range res.Kvs {
			startTime := string(v.Value)
			startTimeInt, err := strconv.ParseInt(startTime, 10, 64)
			if err != nil {
				ziLog.Error(fmt.Sprintf("MatchV1HeardBeat strconv.ParseInt err: %v", err), debug)
				continue
			}

			// log.Println("matchV1HeardBeat time err:", startTime)
			// 获取时间戳
			if time.Now().Unix() < startTimeInt { // 未到开始时间
				continue
			}
			groupId := filepath.Base(string(v.Key))
			// log.Println("matchV1HeardBeat groupId:", groupId)
			// 如果匿名，则返回
			if ok, _ := battlematchv1.MatchBattleAnonymousGetByGroupId(first_ctx, groupId); ok {
				continue
			}
			//
			sendUidList, _, _, _ := getUidListByGroupId(groupId)
			if len(sendUidList) == 0 {
				battlematchv1.UnregisterBattleV1ByGroupId(first_ctx, groupId)
				continue
			}
			// 获取状态
			status, err := battlematchv1.MatchGroupStatusGet(first_ctx, groupId)
			if err != nil {
				ziLog.Error(fmt.Sprintf("MatchV1HeardBeat MatchGroupStatusGet err: %v， group: %v", err, groupId), debug)
				continue
			}
			// 通过groupId 获取UserId
			userIdList, err := battlematchv1.QueryVs1GroupInfo(first_ctx, groupId)
			if err != nil {
				ziLog.Error(fmt.Sprintf("MatchV1HeardBeat QueryVs1GroupInfo err: %v， group: %v", err, groupId), debug)
				continue
			}
			switch status {
			case match_battle_status_Confirm:
				// 查询uid
				matchBattleStartGamedConfirm(groupId, userIdList)
				continue
			case match_battle_status_start:
				continue
			default:
				continue
			}
			// 后面暂时不用，等以后上连线再开启
			// for i, v := range userIdList {
			// 	switch i {
			// 	case 0:
			// 		roomId1 = queryRoomIdInterconvertAnchorOpenId(v)
			// 		if roomId1 == "" {
			// 			log.Println("roomId1 is nil")
			// 			continue
			// 		}
			// 	case 1:
			// 		roomId2 = queryRoomIdInterconvertAnchorOpenId(v)
			// 		if roomId2 == "" {
			// 			log.Println("roomId1 is nil")
			// 			continue
			// 		}
			// 	}
			// }
			// if err := twoConnect("heartbeat", roomId1, roomId2, groupId); err != nil {
			// 	// 错误可能后续不需要记录
			// 	ziLog.Error( fmt.Sprintf("MatchV1HeardBeat heartbeat err: %v， group: %v", err, groupId), debug)
			// 	battlematchv1.MatchGroupStatusSet(first_ctx, groupId, match_battle_status_error)
			// }
		}
	}
}
