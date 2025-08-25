package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kongshui/danmu/common"
)

// 同步对局状态,开发者在对局开始时调用，同步对局开始事件；在对局结束时再次调用，同步对局结束事件。
func ksSyncGameStatus(data SyncGameStatusStruct, label string, isFirst bool) bool {
	if is_mock {
		return true
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
			group string = "绝对中立"
		)
		for _, v := range data.GroupResultList {
			if v.Result == 1 {
				group = v.GroupId
				break
			}
		}
		sdata = map[string]any{
			"roundId":    data.RoundId,
			"roundType":  "singleGroup",
			"roundGroup": "守序善良,中立善良,混乱善良,绝对中立,守序邪恶,中立邪恶,混乱邪恶",
			"result": map[string]string{
				"singleGroupRoundResult": group,
			},
			"bulletPlay": map[string]any{},
		}

	}

	sdataByte, err := json.Marshal(sdata)
	if err != nil {
		ziLog.Error(fmt.Sprintf("SyncGameStatus1 err: %v", err), debug)
		return false
	}
	urlPath := KsUrlSet(url_SyncGameStatusUrl)
	if urlPath == "" {
		ziLog.Error("urlSet err: urlPath is nil", debug)
		return false
	}
	response, err := common.HttpRespond("POST", urlPath, kuaiShouBindBodyToByte(data.RoomId, "round", label, string(sdataByte)), headers)
	if err != nil {
		ziLog.Error(fmt.Sprintf("SyncGameStatus2 err: %v", err), debug)
		return false
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		ziLog.Error(fmt.Sprintf("SyncGameStatus3 err: %v", response.StatusCode), debug)
		return false
	}
	var (
		request any
	)
	if err := json.NewDecoder(response.Body).Decode(&request); err != nil {
		ziLog.Error(fmt.Sprintf("SyncGameStatus json.NewDecoder err: %v", err), debug)
		return false
	}
	if int64(request.(map[string]any)["result"].(float64)) != 1 {
		if request.(map[string]any)["errorMsg"].(string) == "直播已关播" || request.(map[string]any)["errorMsg"].(string) == "无游戏中记录" ||
			request.(map[string]any)["errorMsg"].(string) == "直播不存在" || label == "stop" {
			ziLog.Info(fmt.Sprintf("SyncGameStatus 直播已关播, roomId: %v, 用户Id： %v ", data.RoomId, data.AnchorOpenId), debug)
			// 后续清理
			endClean(data.RoomId, data.AnchorOpenId)
			return true
		}
		if label == "start" && isFirst {
			ziLog.Error(fmt.Sprintf("SyncGameStatus result err: %v, roomId: %v, 用户Id： %v ", request, data.RoomId, data.AnchorOpenId), debug)
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
							GroupResultList: []GroupResultList{{GroupId: "绝对中立", Result: 1}},
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
	return true
}

// 抖音同步对局状态,开发者在对局开始时调用，同步对局开始事件；在对局结束时再次调用，同步对局结束事件。
func dySyncGameStatus(data SyncGameStatusStruct) bool {
	headers := map[string]string{
		"Content-Type": "application/json",
		"X-Token":      accessToken.Token,
	}
	data.AppId = app_id
	body, err := json.Marshal(data)
	if err != nil {
		ziLog.Error(fmt.Sprintf("SyncGameStatus1 err: %v", err), debug)
		return false
	}
	response, err := common.HttpRespond("POST", url_round_sync_status, body, headers)
	if err != nil {
		ziLog.Error(fmt.Sprintf("SyncGameStatus2 err: %v", err), debug)
		return false
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		ziLog.Error(fmt.Sprintf("SyncGameStatus status err: %v", response.Status), debug)
		return false
	}
	var (
		request any
	)
	if err := json.NewDecoder(response.Body).Decode(&request); err != nil {
		ziLog.Error(fmt.Sprintf("SyncGameStatus json.NewDecoder err: %v", err), debug)
		return false
	}
	if int64(request.(map[string]any)["errcode"].(float64)) != 0 {
		ziLog.Error(fmt.Sprintf("dySyncGameStatus err_no: %v,err: %v", request.(map[string]any)["errcode"], request.(map[string]any)["errmsg"]), debug)
		return false
	}
	return true
}

// 上传世界榜单列表数据,接口是一次性上传排好序的榜单Top 150的用户数据，本接口是上传所有参与玩法并且有战绩的用户数据，
// 批量上报，底层是以用户id+榜单版本world_rank_version为维度存储，用户id+榜单版本world_rank_version维度重复上报，是覆盖写的逻辑；
// 也用作上报用户世界榜单的累计战绩，单次调用最多上报50个用户战绩。
func dyWorldRankListUpload(data any, url string) bool {
	headers := map[string]string{
		"Content-Type": "application/json",
		"X-Token":      accessToken.Token,
	}
	body, err := json.Marshal(data)
	if err != nil {
		ziLog.Error(fmt.Sprintf("WorldRankListUpload1 err: %v", err), debug)
		return false
	}
	response, err := common.HttpRespond("POST", url, body, headers)
	if err != nil {
		ziLog.Error(fmt.Sprintf("WorldRankListUpload2 err: %v", err), debug)
		return false
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		ziLog.Error(fmt.Sprintf("WorldRankListUpload status err: %v", response.Status), debug)
		return false
	}
	var (
		request any
	)
	if err := json.NewDecoder(response.Body).Decode(&request); err != nil {
		// 关闭响应体以释放资源
		ziLog.Error(fmt.Sprintf("WorldRankListUpload3 err: %v", err), debug)
		return false
	}

	errCode := int64(request.(map[string]any)["errcode"].(float64))
	if errCode != 0 {
		ziLog.Error(fmt.Sprintf("WorldRankListUpload errCode: %v, errmsg: %s", errCode, request.(map[string]any)["errmsg"].(string)), debug)
		return false
	}
	return true
}

// 世界排行榜设置
func worldRankSet(worldRankVersion string) bool {
	headers := map[string]string{
		"Content-Type": "application/json",
		"X-Token":      accessToken.Token,
	}
	body, err := json.Marshal(map[string]any{
		"app_id":             app_id,
		"is_online_version":  config.App.IsOnline,
		"world_rank_version": worldRankVersion,
	})
	if err != nil {
		ziLog.Error(fmt.Sprintf("WorldRankSet1 err: %v", err), debug)
		return false
	}
	response, err := common.HttpRespond("POST", url_set_world_rank_version_url, body, headers)
	if err != nil {
		ziLog.Error(fmt.Sprintf("WorldRankSet2 err: %v", err), debug)
		return false
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		ziLog.Error(fmt.Sprintf("WorldRankSet status err: %v", response.Status), debug)
		return false
	}
	var (
		request any
	)
	if err := json.NewDecoder(response.Body).Decode(&request); err != nil {
		ziLog.Error(fmt.Sprintf("WorldRankSet3 err: %v", err), debug)
		return false
	}

	errCode := int64(request.(map[string]any)["errcode"].(float64))
	if errCode != 0 {
		ziLog.Error(fmt.Sprintf("WorldRankSet errCode: %v, errmsg: %s", errCode, request.(map[string]any)["errmsg"].(string)), debug)
	}
	return errCode == 0
}

// 完成用户世界榜单的累计战绩上报
// 当到达截榜时间且有线上生效的世界榜单版本时，在完成用户世界榜单的累计战绩上报后，调用本接口，标记本次截榜时间内的世界榜单累计战绩完成了上报。
func worldRankCompleteUpload() bool {
	headers := map[string]string{
		"Content-Type": "application/json",
		"X-Token":      accessToken.Token,
	}
	body, err := json.Marshal(map[string]any{
		"app_id":             app_id,
		"is_online_version":  config.App.IsOnline,
		"world_rank_version": currentRankVersion,
		"complete_time":      time.Now().Unix(),
	})
	if err != nil {
		ziLog.Error(fmt.Sprintf("WorldRankCompleteUpload1 err: %v", err), debug)
		return false
	}
	response, err := common.HttpRespond("POST", url_complete_upload_url, body, headers)
	if err != nil {
		ziLog.Error(fmt.Sprintf("WorldRankCompleteUpload2 err: %v", err), debug)
		return false
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return false
	}
	var (
		request any
	)
	if err := json.NewDecoder(response.Body).Decode(&request); err != nil {
		ziLog.Error(fmt.Sprintf("WorldRankCompleteUpload3 err: %v", err), debug)
		return false
	}
	if debug {
		ziLog.Info(fmt.Sprintf("世界排行榜完成上传完成 worldRankCompleteUpload: %v", request), debug)
	}
	return int64(request.(map[string]any)["errcode"].(float64)) == 0
}
