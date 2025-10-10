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
		CreateRoomId(uid, "")
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
	isMember := blackAnchorListIsMember(roomReponse.Data.Info.AnchorOpenId)
	if isMember {
		return fmt.Errorf("anchor is black")
	}
	go userInfoCompareStore(roomReponse.Data.Info.AnchorOpenId, roomReponse.Data.Info.NickName, roomReponse.Data.Info.AvatarUrl, true)
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
