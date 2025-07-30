package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	battlematchv1 "github.com/kongshui/danmu/battlematch/v1"

	"github.com/kongshui/danmu/model/pmsg"

	"github.com/kongshui/danmu/common"

	"google.golang.org/protobuf/proto"
)

// 发送下行消息
func PushDownLoadMessage(msgId pmsg.MessageId, uidStrList []string, data []byte) error {
	// openId := queryRoomIdToAnchorOpenId(roomId)
	// uidStr := findUidStrFromOpenid(openId)
	if len(uidStrList) < 1 {
		return errors.New("uid err: uid is nil")
	}
	dataBody := &pmsg.MessageBody{
		MessageId:   uint32(msgId),
		MessageType: msgId.String(),
		MessageData: data,
		Timestamp:   time.Now().UnixMilli(),
	}
	requestBody, err := proto.Marshal(dataBody)
	if err != nil {
		return err
	}
	// Marshal将slice转换为JSON格式的字节切片
	// jsonData, err := json.Marshal(openIdList)
	// if err != nil {
	// 	return err
	// }
	header := map[string]string{
		"Content-Type": "application/json",
		"x-event-type": msgId.String(),
	}
	// 转成json格式的字符串
	uidStrListbyte, err := json.Marshal(uidStrList)
	if err != nil {
		return errors.New("pushDownLoadMessage json.Marshal uidStrList err: " + err.Error())
	}
	// 将字节切片转换为字符串
	header["x-client-uuid"] = string(uidStrListbyte)
	if forward_domain.Len() == 0 {
		if err := oneGetFowardDomain(); err != nil {
			return errors.New("pushDownLoadMessage 前端服务器数量为零")
		}
		ziLog.Info("pushDownLoadMessage get forward_domain, success ", debug)
	}
	var index int
	if forward_domain.Len() > 1 {
		index = int(msgId) % forward_domain.Len()
	}
	domain := forward_domain.Get(index)
	if domain == "" {
		if err := oneGetFowardDomain(); err != nil {
			return errors.New("pushDownLoadMessage 前端服务器数量为零")
		}
		domain = forward_domain.Get(0)
		if domain == "" {
			return errors.New("pushDownLoadMessage domain 前端服务器数量为零")
		}
	}
	forwardPath, _ := url.JoinPath(domain, forward_domain_uri)
	response, err := common.HttpRespond("POST", forwardPath, requestBody, header)
	if err != nil {
		forward_domain.Remove(domain)
		time.Sleep(100 * time.Millisecond) // 等待100毫秒后重试
		ziLog.Info("pushDownLoadMessage delete forward_domain: "+forward_domain.Get(index), debug)
		return PushDownLoadMessage(msgId, uidStrList, data)
	}
	// 关闭响应体以释放资源
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return err
	}
	var result map[string]any
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return err
	}
	if int64(result["err_code"].(float64)) != 0 {
		return errors.New(result["err_msg"].(string))
	}
	return nil
}

// 通过openId查看uidList,第二个列表是openId列表，第三个为是否为group,第四个为组名
func getUidListByOpenId(openId string) ([]string, []string, bool, string) {
	var uidStrList []string = make([]string, 0)
	// 先获取单个
	uid := queryUidByOpenid(openId)
	if uid != "" {
		uidStrList = append(uidStrList, uid)
	}
	ok, groupId := battlematchv1.IsInVs1Group(first_ctx, openId)
	if !ok {
		return uidStrList, []string{}, false, ""
	}
	if battlematchv1.QueryOpenIdInMatchDisconnect(first_ctx, groupId, openId) {
		return uidStrList, []string{}, false, ""
	}
	return getUidListByGroupId(groupId)
}

// 通过groupId查看uidList
func getUidListByGroupId(groupId string) ([]string, []string, bool, string) {
	var uidStrList []string = make([]string, 0)
	openIdList, err := battlematchv1.QueryVs1GroupInfo(first_ctx, groupId)
	if err != nil {
		return uidStrList, []string{}, true, groupId
	}
	for _, openIdStr := range openIdList {
		if battlematchv1.QueryOpenIdInMatchDisconnect(first_ctx, groupId, openIdStr) {
			continue
		}
		uid := queryUidByOpenid(openIdStr)
		if uid == "" {
			continue
		}
		uidStrList = append(uidStrList, uid)
	}
	return uidStrList, openIdList, true, groupId
}
