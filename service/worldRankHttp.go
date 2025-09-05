package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kongshui/danmu/common"
)

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
		ziLog.Error(fmt.Sprintf("WorldRankListUpload err: %v, data: %s", request, string(body)), debug)
		return false
	}
	return true
}

// 世界排行榜设置
func worldRankSet(worldRankVersion string) {
	defer RecoverFunc()
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
		panic(fmt.Sprintf("WorldRankSet1 err: %v", err))
	}
	response, err := common.HttpRespond("POST", url_set_world_rank_version_url, body, headers)
	if err != nil {
		panic(fmt.Sprintf("WorldRankSet2 err: %v", err))
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		panic(fmt.Sprintf("WorldRankSet status err: %v", response.Status))
	}
	var (
		request any
	)
	if err := json.NewDecoder(response.Body).Decode(&request); err != nil {
		panic(fmt.Sprintf("WorldRankSet3 err: %v", err))
	}

	errCode := int64(request.(map[string]any)["errcode"].(float64))
	if errCode != 0 {
		panic(fmt.Sprintf("WorldRankSet errCode: %v, errmsg: %s", errCode, request.(map[string]any)["errmsg"].(string)))
	}
}

// 完成用户世界榜单的累计战绩上报
// 当到达截榜时间且有线上生效的世界榜单版本时，在完成用户世界榜单的累计战绩上报后，调用本接口，标记本次截榜时间内的世界榜单累计战绩完成了上报。
func worldRankCompleteUpload() {
	defer RecoverFunc()
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
		panic(fmt.Sprintf("WorldRankCompleteUpload1 err: %v", err))
	}
	response, err := common.HttpRespond("POST", url_complete_upload_url, body, headers)
	if err != nil {
		panic(fmt.Sprintf("WorldRankCompleteUpload2 err: %v", err))
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		panic(fmt.Sprintf("WorldRankCompleteUpload status err: %v", response.Status))
	}
	var (
		request any
	)
	if err := json.NewDecoder(response.Body).Decode(&request); err != nil {
		panic(fmt.Sprintf("WorldRankCompleteUpload3 err: %v", err))
	}
	if debug {
		ziLog.Info(fmt.Sprintf("世界排行榜完成上传完成 worldRankCompleteUpload: %v", request), debug)
	}
	if int64(request.(map[string]any)["errcode"].(float64)) != 0 {
		panic(fmt.Sprintf("WorldRankCompleteUpload errCode: %v, errmsg: %s", int64(request.(map[string]any)["errcode"].(float64)), request.(map[string]any)["errmsg"].(string)))
	}
}
