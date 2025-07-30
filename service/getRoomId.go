package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/kongshui/danmu/model/pmsg"
	"github.com/kongshui/danmu/sse"

	"github.com/kongshui/danmu/common"

	"google.golang.org/protobuf/proto"
)

// 获取roomId和anchorOpenId
// func getRoomId(token, uid string) (*pmsg.AnchorInfoMessage, error) {
// 	requestBody, err := json.Marshal(map[string]any{
// 		"token": token,
// 	})
// 	if err != nil {
// 		return &pmsg.AnchorInfoMessage{}, err
// 	}
// 	header := map[string]string{
// 		"Content-Type": "application/json",
// 	}

// 	if !isIndouyinyun {
// 		header["X-Token"] = accessToken.Token
// 	}
// 	response, err := common.HttpRespond("POST", url_GetLiveInfoUrl, requestBody, header)
// 	if err != nil {
// 		return &pmsg.AnchorInfoMessage{}, err
// 	}
// 	// 关闭响应体以释放资源
// 	defer response.Body.Close()
// 	if response.StatusCode != http.StatusOK {
// 		return &pmsg.AnchorInfoMessage{}, err
// 	}
// 	var (
// 		roomReponse any
// 		roomInfo    = &pmsg.AnchorInfoMessage{}
// 	)
// 	err = json.NewDecoder(response.Body).Decode(&roomReponse)
// 	if err != nil {
// 		return &pmsg.AnchorInfoMessage{}, err
// 	}
// 	if len(roomReponse.(map[string]any)["data"].(map[string]any)) == 0 {
// 		return &pmsg.AnchorInfoMessage{}, errors.New(roomReponse.(map[string]any)["errmsg"].(string))
// 	}
// 	for k, v := range roomReponse.(map[string]any)["data"].(map[string]any)["info"].(map[string]string) {
// 		switch k {
// 		case "room_id":
// 			roomInfo.RoomId = v
// 		case "anchor_open_id":
// 			roomInfo.AnchorOpenId = v
// 		case "avatar_url":
// 			roomInfo.AvatarUrl = v
// 		case "nick_name":
// 			roomInfo.NickName = v
// 		}
// 	}
// 	//存储用户信息
// 	setRoomIdToAnchorOpenId(roomInfo.RoomId, roomInfo.AnchorOpenId)
// 	//存储直播间信息
// 	var (
// 		data common.RoomRegister
// 	)
// 	level, _ := queryUserLevel(roomInfo.AnchorOpenId)
// 	data.Uuid = uid
// 	data.RoomId = roomInfo.RoomId
// 	data.UserId = roomInfo.AnchorOpenId
// 	data.GradeLevel = level
// 	dataByte, err := json.Marshal(data)
// 	if err != nil {
// 		log.Println("json转换失败， info:", data, err, "err: ", err)
// 	}

// 	etcdClient.Client.Put(first_ctx, common.RoomInfo_Register_key+"/"+roomInfo.RoomId, string(dataByte))
// 	return roomInfo, nil
// }

// // 查询玩家段位
// func queryUserLevel(userId string) (int64, error) {
// 	level, err := rdb.ZScore(player_grade_level_db, userId)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return int64(level), nil
// }

func setRoomInfo(uid string, roomInfo *pmsg.AnchorInfoMessage) (*pmsg.AnchorInfoMessage, error) {
	//存储用户信息
	setRoomIdToAnchorOpenId(roomInfo.RoomId, roomInfo.AnchorOpenId)
	if debug {
		log.Println(uid)
	}
	//存储直播间信息
	var (
		data common.RoomRegister
	)
	data.Uuid = uid
	data.RoomId = roomInfo.RoomId
	data.UserId = roomInfo.AnchorOpenId
	data.OpenId = roomInfo.AnchorOpenId
	dataByte, err := json.Marshal(data)
	if err != nil {
		ziLog.Error(fmt.Sprintf("json转换失败， info: %v，err: %v", data, err), debug)
	}

	etcdClient.Client.Put(first_ctx, path.Join("/", config.Project, common.RoomInfo_Register_key, roomInfo.RoomId), string(dataByte))
	return roomInfo, nil
}

// 抖音获取主播信息
func dyGetAnchorInfo(uid, token string) error {
	if is_mock {
		go userInfoCompareStore("123456789", "dytest", "http://tupian.geimian.com/pic/2016/11/2016-11-05_213153.jpg")
		data := &pmsg.AnchorInfoMessage{}
		data.AnchorOpenId = "123456789"
		data.AvatarUrl = "http://tupian.geimian.com/pic/2016/11/2016-11-05_213153.jpg"
		data.NickName = "dytest"
		data.RoomId = "987654321"
		databyte, _ := proto.Marshal(data)
		setRoomInfo(uid, data)
		connect(data.GetRoomId(), data.GetAnchorOpenId())
		if err := sse.SseSend(pmsg.MessageId_TokenAck, []string{uid}, databyte); err != nil {
			return fmt.Errorf("DyGetAnchorInfo pushDownLoadMessage err: %v", err)
		}
		return nil
	}
	requestBody, err := json.Marshal(map[string]any{
		"token": token,
	})
	if err != nil {
		return err
	}
	header := map[string]string{
		"Content-Type": "application/json",
		"X-Token":      accessToken.Token,
	}

	response, err := common.HttpRespond("POST", url_get_anchor_info_url, requestBody, header)
	if err != nil {
		return err
	}
	// 关闭响应体以释放资源
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return err
	}
	var (
		roomReponse GetDyAnchorInfoStruct
	)
	err = json.NewDecoder(response.Body).Decode(&roomReponse)
	if err != nil {
		return err
	}
	if roomReponse.Errcode != 0 {
		return fmt.Errorf("DyGetAnchorInfo Errcode err, code: %v, errmsg: %v", roomReponse.Errcode, roomReponse.Errmsg)
	}
	go userInfoCompareStore(roomReponse.Data.Info.AnchorOpenId, roomReponse.Data.Info.NickName, roomReponse.Data.Info.AvatarUrl)
	data := &pmsg.AnchorInfoMessage{}
	data.AnchorOpenId = roomReponse.Data.Info.AnchorOpenId
	data.AvatarUrl = roomReponse.Data.Info.AvatarUrl
	data.NickName = roomReponse.Data.Info.NickName
	data.RoomId = strconv.FormatInt(roomReponse.Data.Info.RoomId, 10)
	databyte, _ := proto.Marshal(data)
	setRoomInfo(uid, data)
	connect(data.GetRoomId(), data.GetAnchorOpenId())
	if err := sse.SseSend(pmsg.MessageId_TokenAck, []string{uid}, databyte); err != nil {
		return fmt.Errorf("DyGetAnchorInfo pushDownLoadMessage err: %v", err)
	}
	return nil
}
