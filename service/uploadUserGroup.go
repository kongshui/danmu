package service

import (
	"encoding/json"
	"fmt"

	"github.com/kongshui/danmu/common"
)

// 上传完玩家组信息
func dyUploadUserGroup(roomId, openId, groupId string, roundId int64) {
	defer RecoverFunc()
	headers := map[string]string{
		"Content-Type": "application/json",
		"X-Token":      accessToken.Token,
	}
	bodydata := map[string]any{
		"app_id":   app_id,
		"group_id": groupId,
		"open_id":  openId,
		"room_id":  roomId,
		"round_id": roundId,
	}
	body, err := json.Marshal(bodydata)
	if err != nil {
		panic("上传完玩家组信息解析失败: " + err.Error())
	}
	headers["access-token"] = accessToken.Token
	response, err := common.HttpRespond("POST", url_upload_user_group_url, body, headers)
	if err != nil {
		panic("上传完玩家组信息请求失败: " + err.Error())
	}
	defer response.Body.Close()
	var (
		request any
	)

	if err := json.NewDecoder(response.Body).Decode(&request); err != nil {
		panic("上传完玩家组信息解析返回值失败: " + err.Error())
	}

	if response.StatusCode != 200 {
		panic("上传完玩家组信息解析状态码有误: " + string(body))
	}
	if int64(request.(map[string]any)["errcode"].(float64)) != 0 {
		panic(fmt.Sprintf("上传完玩家组信息解析返回值有误: %v, 数据为： %v", request, body))
	}
}
