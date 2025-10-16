package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kongshui/danmu/model/pmsg"

	"github.com/kongshui/danmu/common"
	"github.com/kongshui/danmu/sse"

	"google.golang.org/protobuf/proto"
)

// 直接推送
func pushDyBasePayloayDirect(roomId, anchorOpenId, msgType string, data []byte) {
	var msgAck MsgAckStruct
	msgAck.RoomId = roomId
	msgAck.AppId = app_id
	msgAck.AckType = 1
	var dataList []MsgAckInfoStruct = make([]MsgAckInfoStruct, 0)
	var (
		score   float64
		isFirst bool = true
		// avatarUrl string
		nickName string
	)
	anchorName, err := UserInfoGet(anchorOpenId, true)
	if err != nil {
		_, nickName, err = mysql.QueryPlayerInfo(anchorOpenId)
		if err != nil {
			ziLog.Error(fmt.Sprintf("pushDyBasePayloayDirect get userInfoGet err:  %v", err), debug)
		}
		anchorName.NickName = nickName
		// anchorName.AvatarUrl = avatarUrl
	}
	roundId, _ := queryRoomIdToRoundId(roomId)
	switch msgType {
	case "live_comment":
		getData := []ContentPayloadStruct{}
		// if err := json.Unmarshal(data, &getData); err != nil {
		// 	ziLog.Error(fmt.Sprintf("PushDyBasePayloayDirect json.Unmarshal err: %v, data: %v", err, data), debug)
		// }
		if err := json.NewDecoder(strings.NewReader(string(data))).Decode(&getData); err != nil {
			ziLog.Error(fmt.Sprintf("PushDyBasePayloayDirect json.Unmarshal err: %v, data: %v", err, data), debug)
		}
		for _, v := range getData {
			var msgAckInfo MsgAckInfoStruct
			msgAckInfo.MsgId = v.MsgId
			msgAckInfo.MsgType = msgType
			msgAckInfo.ClientTime = time.Now().UnixMilli()
			dataList = append(dataList, msgAckInfo)
			if strings.HasPrefix(v.Content, "666") && strings.HasSuffix(v.Content, "666") {
				score = live_like_score
			}

			if score != 0 {
				// 添加积分
				go matchAddIntrage(roomId, v.SecOpenid, score)
			}
			dyPayloadSendMessage(v, pmsg.MessageId_LiveComment, roomId, anchorOpenId)
		}
	case "live_gift":
		getData := []GiftPayloadStruct{}
		if err := json.NewDecoder(strings.NewReader(string(data))).Decode(&getData); err != nil {
			ziLog.Error(fmt.Sprintf("PushDyBasePayloayDirect json.Unmarshal err: %v, data: %v", err, string(data)), debug)
		}
		ziLog.Gift(fmt.Sprintf("PushDyBasePayloayDirect getData: %v", string(data)), debug)
		for _, v := range getData {
			var msgAckInfo MsgAckInfoStruct
			msgAckInfo.MsgId = v.MsgId
			msgAckInfo.MsgType = msgType
			msgAckInfo.ClientTime = time.Now().UnixMilli()
			dataList = append(dataList, msgAckInfo)
			if !v.Test {
				score = float64(v.GiftNum) * giftToScoreMap[v.SecGiftId]
			} else {
				score = 0
			}
			// bByte, _ := json.Marshal(v)
			// var dataMap map[string]any
			// if err := json.Unmarshal(bByte, &dataMap); err != nil {
			// 	ziLog.Write(logError, fmt.Sprintf("PushDyBasePayloayDirect json.Unmarshal err: %v, data: %v", err, bByte), debug)
			// }
			if score != 0 {
				go matchAddIntrage(roomId, v.SecOpenid, score)
				// 数据到数据库中，防止数据丢失
				go mysql.InsertGiftData(roomId, anchorOpenId, anchorName.NickName, strconv.FormatInt(roundId, 10), v.SecOpenid, v.Nickname, v.MsgId, v.SecGiftId,
					v.GiftNum, v.GiftValue, v.Test)
				// 设置用户是否已经消费
				if isFirst {
					setIsConsume(v.SecOpenid, time.Now().UnixMilli())
					isFirst = false
				}
			}
			dyPayloadSendMessage(v, pmsg.MessageId_liveGift, roomId, anchorOpenId)
		}
	case "live_like":
		getData := []LiveLikePayloadStruct{}
		if err := json.NewDecoder(strings.NewReader(string(data))).Decode(&getData); err != nil {
			ziLog.Error(fmt.Sprintf("PushDyBasePayloayDirect json.Unmarshal err: %v, data: %v", err, data), debug)
		}
		for _, v := range getData {
			var msgAckInfo MsgAckInfoStruct
			msgAckInfo.MsgId = v.MsgId
			msgAckInfo.MsgType = msgType
			msgAckInfo.ClientTime = time.Now().UnixMilli()
			dataList = append(dataList, msgAckInfo)
			score = live_like_score
			if score != 0 {
				go matchAddIntrage(roomId, v.SecOpenid, score)
				// 送礼直接添加到世界排行榜
				// go worldRankNumerAdd(v.(map[string]any)["userInfo"].(map[string]any)["userId"].(string), score)
			}
			dyPayloadSendMessage(v, pmsg.MessageId_liveLike, roomId, anchorOpenId)
		}
	}
	//分数不为0时添加积分
	jsonData, err := json.Marshal(dataList)
	if err != nil {
		log.Println("json.Marshal err: ", err)
		return
	}
	msgAck.Data = string(jsonData)
	if err := msgAckSend(msgAck); err != nil {
		log.Println("MsgAck err: ", err)
	}
}

// d抖音回执
func msgAckSend(data MsgAckStruct) error {
	requestBody, err := json.Marshal(data)
	if err != nil {
		return err
	}
	header := map[string]string{
		"Content-Type": "application/json",
		"access-token": accessToken.Token,
	}
	_, err = common.HttpRespond("POST", url_live_data_ack_url, requestBody, header)
	return err
}

// dySendMessage

func dyPayloadSendMessage(v any, msgId pmsg.MessageId, roomId, anchorOpenid string) {
	sendUidList, _, _, _ := getUidListByOpenId(anchorOpenid)
	if len(sendUidList) == 0 {
		ziLog.Error(fmt.Sprintf("dyPayloadSendMessage sendUidList is nil, roomId: %v, anchorOpenid: %v, data: %v", roomId, anchorOpenid, v), debug)
		return
	}
	jData, err := json.Marshal(v)
	if err != nil {
		ziLog.Error(fmt.Sprintf("dyPayloadSendMessage jpushBasePayloayDirect json.Marshal err:  %v,失败数据为： %v", err, v), debug)
		return
	}
	endSendData := platFormPool.Get().(*pmsg.PlatFormDataSend)
	defer platFormPool.Put(endSendData)
	endSendData.Reset()
	endSendData.TimeStamp = time.Now().UnixMilli()
	endSendData.Data = jData
	endSendDatabyte, _ := proto.Marshal(endSendData)
	// 推送消息
	if err := sse.SseSend(msgId, sendUidList, endSendDatabyte); err != nil {
		ziLog.Error(fmt.Sprintf("dyPayloadSendMessage 推送消息失败:  %v,失败数据为： %v", err, v), debug)
	}
}

// 获取抖音失败信息
// 将失败数据推送到抖音云
func pushFailDataToDy(roomId string) error {
	// 获取失败数据
	failData, err := getPushFailData(roomId)
	if err != nil {
		return err
	}
	for _, v := range failData {
		anchorOpenId := QueryRoomIdInterconvertAnchorOpenId(roomId)
		msgType := v.(map[string]string)["msg_type"]
		data := v.(map[string]string)["payload"] // 将失败数据推送到抖音云
		// 完成失失败消息推送
		go pushDyBasePayloayDirect(roomId, anchorOpenId, msgType, []byte(data))
	}
	return nil
}

// 获取推送失败数据
func getPushFailData(roomId string) (map[string]any, error) {
	var getFailData GetLiveFailRequestStruct
	getFailData.Appid = app_id
	getFailData.Roomid = roomId
	getFailData.PageNum = "1"
	getFailData.PageSize = "50"
	getFailData.MsgType = "live_gift"
	query, err := common.StructToStringMap(getFailData)
	if err != nil {
		return nil, err
	}
	url, err := common.GetUrl(url_fail_get_live_data_url, query)
	if err != nil {
		return nil, err
	}
	headers := map[string]string{
		"Content-Type": "application/json",
		"access-token": accessToken.Token,
	}
	response, err := common.HttpRespond("GET", url, nil, headers)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, err
	}
	var (
		request map[string]any
	)
	if err := json.NewDecoder(response.Body).Decode(&request); err != nil {
		return nil, err
	}
	ziLog.Info(fmt.Sprintf("getPushFailData request: %v", request), debug)
	if int64(request["err_no"].(float64)) != 0 {
		return nil, errors.New(request["err_msg"].(string))
	}
	return request["data"].(map[string]any)["data_list"].(map[string]any), nil
}
