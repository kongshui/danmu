package service

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/kongshui/danmu/model/pmsg"
	"github.com/kongshui/danmu/sse"

	pb "github.com/kongshui/danmu/model/battlecalv1pb"

	"google.golang.org/protobuf/proto"
)

func ksPushBasePayloay(data KsCallbackStruct) {
	pushType := data.Data.PushType
	sendUidList, _, isGroup, _ := getUidListByOpenId(data.Data.AuthorOpenId)
	if len(sendUidList) == 0 {
		ziLog.Error(fmt.Sprintf("ksPushBasePayloay queryRoomIdToUid nil， roomId: %v, openId, %v, 数据为： %v", data.Data.RoomCode, data.Data.AuthorOpenId, data), debug)
		return
	}

	// 上报消息接收状态
	ack := KsMsgAckReceiveStruct{
		UniqueMessageId:     data.Data.UniqueMessageId,
		PushType:            data.Data.PushType,
		CpServerReceiveTime: time.Now().UnixMilli(),
	}
	var (
		isgift    bool = false
		isLottery bool = false
		isSendAck bool = false // 是否发送ack
		groupGrpc      = &pb.AddGiftToGroupReq{}
		isFirst   bool = true // 是否第一次
		openId    string
		nickName  string
		avatarUrl string
	)
	endSendData := platFormPool.Get().(*pmsg.PlatFormDataSend)
	defer platFormPool.Put(endSendData)
	endSendData.OpenId = data.Data.AuthorOpenId
	endSendData.RoomId = data.Data.RoomCode
	endSendData.PushType = data.Data.PushType
	// 添加uniqueMessageId到redis中，防止重复推送
	giftsendSet(data.Data.RoomCode, data.Data.UniqueMessageId)
	// 是否是group
	// if isGroup {
	// 	groupGrpc.GroupId = groupId
	// 	groupGrpc.AnchorOpenId = data.Data.AuthorOpenId
	// 	groupGrpc.AnchorOpenIdList = openIdList
	// }

	for _, v := range data.Data.Payload {
		var (
			score            float64
			msgId            pmsg.MessageId
			anchor_nick_name string
		)
		jsonByte, err := json.Marshal(v)
		if err != nil {
			ziLog.Error(fmt.Sprintf("ksPushBasePayloay json.Marshal err:  %v,失败数据为： %v", err, v), debug)
			continue
		}
		switch pushType {
		case "giftSend":
			isgift = true
			ziLog.Gift(fmt.Sprintf("ksPushBasePayloay giftdata： %v", string(jsonByte)), debug)
			// 送礼,正式上线前将test_去掉，用上一层的
			// if strings.HasPrefix(data.Data.UniqueMessageId, "stress_") || data.Event == "LIVE_INTERACTION_DATA_TEST" {
			gift := KsGiftSendStruct{}
			if err := json.NewDecoder(strings.NewReader(string(jsonByte))).Decode(&gift); err != nil {
				ziLog.Error(fmt.Sprintf("ksPushBasePayloay giftSend  json.Marshal err:  %v，失败数据为： %v", err, v), debug)
				continue
			}
			// 获取用户信息
			openId = gift.UserInfo.UserId
			nickName = gift.UserInfo.NickName
			avatarUrl = gift.UserInfo.AvatarUrl
			// 后端记录数据库
			anchorName, err := userInfoGet(data.Data.AuthorOpenId)
			if err != nil {
				_, anchor_nick_name, err = mysql.QueryPlayerInfo(data.Data.AuthorOpenId)
				if err != nil {
					ziLog.Error(fmt.Sprintf("ksPushBasePayloay giftSend userInfoGet err:  %v", err), debug)
				}
				anchorName.NickName = anchor_nick_name
			}
			roundId, _ := queryRoomIdToRoundId(data.Data.RoomCode)
			if !(strings.HasPrefix(data.Data.UniqueMessageId, "test_") || strings.HasPrefix(data.Data.UniqueMessageId, "stress_")) {
				// 数据到数据库中，防止数据丢失
				go mysql.InsertGiftData(data.Data.RoomCode, data.Data.AuthorOpenId, anchorName.NickName, strconv.FormatInt(roundId, 10), gift.UserInfo.UserId,
					gift.UserInfo.NickName, gift.UniqueNo, gift.GiftId, int(gift.GiftCount), int(gift.GiftTotalPrice), false)

				if isFirst {
					// 设置用户是否已经消费
					setIsConsume(gift.UserInfo.UserId, time.Now().UnixMilli())
					isFirst = false
				}
			}
			if gift.GiftTotalPrice > 0 {
				isSendAck = true
			}
			if gift.GiftId == "11584" {
				isLottery = true
				// 抽奖
				giftMap := lottery(data.Data.AuthorOpenId, gift.UserInfo.UserId, gift.GiftCount)
				ziLog.Gift(fmt.Sprintf("ksPushBasePayloay Lottery,火花数量：%v, giftdata： %v，用户Id： %v, 用户名称： %v", gift.GiftCount, giftMap, gift.UserInfo.UserId, gift.UserInfo.NickName), debug)
				for giftId, giftCount := range giftMap {
					score += giftToScoreMap[giftId] * float64(giftCount)
					if isGroup {
						groupGrpc.IsComment = false
						groupGrpc.OpenId = gift.UserInfo.UserId
						groupGrpc.GiftId = giftId
						groupGrpc.GiftNum = giftCount
						// go grpcSend(groupGrpc, 0)
					}
				}
				giftMapByte, _ := json.Marshal(giftMap)
				lotteryData := &pmsg.LotteryMsg{}
				lotteryData.OpenId = gift.UserInfo.UserId
				lotteryData.NickName = gift.UserInfo.NickName
				lotteryData.HeadImgUrl = gift.UserInfo.AvatarUrl
				lotteryData.LotteryMap = string(giftMapByte)
				lotteryData.Count = gift.GiftCount
				lotteryByte, _ := proto.Marshal(lotteryData)
				endSendData.Data = lotteryByte
				endSendData.TimeStamp = time.Now().UnixMilli()
				endSendData.PushType = "lottery"
				endSendDatabyte, _ := proto.Marshal(endSendData)
				if err := sse.SseSend(pmsg.MessageId_Lottery, sendUidList, endSendDatabyte); err != nil {
					ziLog.Error(fmt.Sprintf("ksPushBasePayloay 推送消息失败,用户Id： %v,用户名称： %v,err:  %v,内容为： %v", gift.UserInfo.UserId, gift.UserInfo.NickName, err, giftMap), debug)
				}
			} else {
				score = giftToScoreMap[gift.GiftId] * float64(gift.GiftCount)
				// if isGroup {
				// 	groupGrpc.IsComment = false
				// 	groupGrpc.OpenId = gift.UserInfo.UserId
				// 	groupGrpc.GiftId = gift.GiftId
				// 	groupGrpc.GiftNum = gift.GiftCount
				// 	go grpcSend(groupGrpc, 0)
				// }
			}
			if strings.HasPrefix(data.Data.UniqueMessageId, "test_") || strings.HasPrefix(data.Data.UniqueMessageId, "stress_") {
				score = 0
			}
			msgId = pmsg.MessageId_liveGift
		case "liveComment":
			commentData := KsLiveCommentStruct{}
			if err := json.NewDecoder(strings.NewReader(string(jsonByte))).Decode(&commentData); err != nil {
				ziLog.Error(fmt.Sprintf("ksPushBasePayloay json.Unmarshal err:  %v,失败数据为： %v", err, v), debug)
				continue
			}
			// 获取用户信息
			openId = commentData.UserInfo.UserId
			nickName = commentData.UserInfo.NickName
			avatarUrl = commentData.UserInfo.AvatarUrl
			// 评论
			if !(strings.HasPrefix(data.Data.UniqueMessageId, "stress_") && data.Event != "LIVE_INTERACTION_DATA_TEST") {
				if commentData.Content == "666" {
					score = live_like_score
					if isGroup {
						groupGrpc.IsComment = true
						groupGrpc.OpenId = commentData.UserInfo.UserId
						groupGrpc.GiftId = "0"
						groupGrpc.GiftNum = 1
						// go grpcSend(groupGrpc, 0)
					}
				} else {
					if isGroup {
						value, ok := commentTogiftId[commentData.Content]
						if ok && !strings.Contains(commentData.Content, "1") && strings.Contains(value, "1") {
							_, err := deleteUserWinStreamCoin(commentData.UserInfo.UserId, commentToCoin[commentData.Content])
							if err != nil {
								return
							}
							groupGrpc.IsComment = true
							groupGrpc.OpenId = commentData.UserInfo.UserId
							groupGrpc.GiftId = value
							groupGrpc.GiftNum = 1
							// go grpcSend(groupGrpc, 0)
						} else {
							isJoin1 := strings.HasPrefix(commentData.Content, "1")
							isJoin11 := strings.HasSuffix(commentData.Content, "1")
							isJoin2 := strings.HasPrefix(commentData.Content, "2")
							isJoin22 := strings.HasSuffix(commentData.Content, "2")
							isJoin3 := strings.HasPrefix(commentData.Content, "加入")
							if (isJoin1 && isJoin11) || (isJoin2 && isJoin22) || isJoin3 {
								groupGrpc.OpenId = commentData.UserInfo.UserId
								// go grpcSend(groupGrpc, 0)
							}
						}
					}
				}
			} else {
				fmt.Println("pushBasePayloayDirect 直播评论测试数据，跳过积分计算：", v)
			}
			msgId = pmsg.MessageId_LiveComment
		case "liveLike":
			// 点赞
			if strings.HasPrefix(data.Data.UniqueMessageId, "stress_") || data.Event == "LIVE_INTERACTION_DATA_TEST" {
				score = 0
			} else {
				liveLikeData := KsLiveLikeStruct{}
				if err := json.NewDecoder(strings.NewReader(string(jsonByte))).Decode(&liveLikeData); err != nil {
					ziLog.Error(fmt.Sprintf("ksPushBasePayloay json.Unmarshal err:  %v,失败数据为： %v", err, v), debug)
					continue
				}
				// 获取用户信息
				openId = liveLikeData.UserInfo.UserId
				nickName = liveLikeData.UserInfo.NickName
				avatarUrl = liveLikeData.UserInfo.AvatarUrl
				score = live_like_score
				if isGroup {
					groupGrpc.IsComment = false
					groupGrpc.OpenId = liveLikeData.UserInfo.UserId
					groupGrpc.GiftId = "0"
					groupGrpc.GiftNum = 1
					// go grpcSend(groupGrpc, 0)
				}
			}
			msgId = pmsg.MessageId_liveLike
		default:
			continue
		}
		//分数不为0时添加积分
		if score != 0 {
			go matchAddIntrage(data.Data.RoomCode, openId, nickName, avatarUrl, score)
			// 送礼直接添加到世界排行榜
			// go worldRankNumerAdd(v.(map[string]any)["userInfo"].(map[string]any)["userId"].(string), score)
		}
		// 格式化消息
		// fmt.Println("pushBasePayloayDirect v: ", v)
		if !isLottery {
			jData, err := json.Marshal(v)
			if err != nil {
				ziLog.Error(fmt.Sprintf("ksPushBasePayloay jpushBasePayloayDirect json.Marshal err:  %v,失败数据为： %v", err, v), debug)
				continue
			}
			endSendData.TimeStamp = time.Now().UnixMilli()
			endSendData.Data = jData
			endSendDatabyte, _ := proto.Marshal(endSendData)
			// 推送消息
			if err := sse.SseSend(msgId, sendUidList, endSendDatabyte); err != nil {
				ziLog.Error(fmt.Sprintf("ksPushBasePayloay 推送消息失败:  %v,失败数据为： %v", err, v), debug)
			}
		}
	}
	if isgift {
		if config.App.IsOnline && isSendAck {
			ksMsgAckSend(data.Data.RoomCode, "cpClientReceive", ack)
		}
		sendAck := &pmsg.KsMsgAck{}
		sendAck.Data = &pmsg.KsMsgAckData{}
		sendAck.RoomId = data.Data.RoomCode
		sendAck.Data.UniqueMessageId = data.Data.UniqueMessageId
		sendAck.Data.PushType = data.Data.PushType
		jData, err := proto.Marshal(sendAck)
		if err != nil {
			ziLog.Error(fmt.Sprintf("ksPushBasePayloay pushBasePayloayDirect proto.Marshal err:  %v,失败数据为： %v", err, ack), debug)
			return
		}
		uid := queryUidByOpenid(data.Data.AuthorOpenId)
		if uid == "" {
			return
		}
		if err := sse.SseSend(pmsg.MessageId_MsgAckSend, []string{uid}, jData); err != nil {
			ziLog.Error(fmt.Sprintf("ksPushBasePayloay sendAck 推送消息失败:  %v", err), debug)
		}
	}
}

// 推送礼物信息
func ksPushGiftSendPayloay(data KsCallbackQueryStruct) {
	endSendData := platFormPool.Get().(*pmsg.PlatFormDataSend)
	defer platFormPool.Put(endSendData)
	endSendData.OpenId = data.AuthorOpenId
	endSendData.RoomId = data.RoomCode
	endSendData.PushType = data.PushType
	ack := KsMsgAckReceiveStruct{
		UniqueMessageId:     data.UniqueMessageId,
		PushType:            data.PushType,
		CpServerReceiveTime: time.Now().UnixMilli(),
	}
	ziLog.Gift(fmt.Sprintf("ksPushGiftSendPayloay giftdata： %v", data), debug)
	// 获取房间信息
	sendUidList, _, isGroup, _ := getUidListByOpenId(data.AuthorOpenId)
	if len(sendUidList) == 0 {
		ziLog.Error(fmt.Sprintf("ksPushBasePayloay queryRoomIdToUid nil， roomId: %v, openId, %v, 数据为： %v", data.RoomCode, data.AuthorOpenId, data), debug)
		return
	}
	var (
		groupGrpc = &pb.AddGiftToGroupReq{}
		openId    string
		nickName  string
		avatarUrl string
	)
	// 添加uniqueMessageId到redis中，防止重复推送
	giftsendSet(data.RoomCode, data.UniqueMessageId)
	for _, v := range data.Payload {
		var (
			score            float64
			isLottery        bool           = false //是否抽奖
			msgId            pmsg.MessageId = pmsg.MessageId_liveGift
			anchor_nick_name string
		)
		jsonByte, err := json.Marshal(v)
		if err != nil {
			ziLog.Error(fmt.Sprintf("ksPushGiftSendPayloay json.Marshal err:  %v,失败数据为： %v", err, v), debug)
			continue
		}
		gift := KsGiftSendStruct{}
		if err := json.NewDecoder(strings.NewReader(string(jsonByte))).Decode(&gift); err != nil {
			ziLog.Error(fmt.Sprintf("ksPushGiftSendPayloay json.Unmarshal err:  %v,失败数据为： %v", err, v), debug)
			continue
		}
		// 获取用户信息
		openId = gift.UserInfo.UserId
		nickName = gift.UserInfo.NickName
		avatarUrl = gift.UserInfo.AvatarUrl
		// 后端记录数据库
		anchorName, err := userInfoGet(data.AuthorOpenId)
		if err != nil {
			_, anchor_nick_name, err = mysql.QueryPlayerInfo(data.AuthorOpenId)
			if err != nil {
				ziLog.Error(fmt.Sprintf("ksPushBasePayloay giftSend userInfoGet err:  %v", err), debug)
			}
			anchorName.NickName = anchor_nick_name
		}
		roundId, _ := queryRoomIdToRoundId(data.RoomCode)
		if !(strings.HasPrefix(data.UniqueMessageId, "test_") || strings.HasPrefix(data.UniqueMessageId, "stress_")) {
			// 设置用户是否已经消费
			setIsConsume(gift.UserInfo.UserId, time.Now().UnixMilli())
			//疾苦到数据库中，防止数据丢失
			go mysql.InsertGiftData(data.RoomCode, data.AuthorOpenId, anchorName.NickName, strconv.FormatInt(roundId, 10), gift.UserInfo.UserId,
				gift.UserInfo.NickName, gift.UniqueNo, gift.GiftId, int(gift.GiftCount), int(gift.GiftTotalPrice), false)
		}

		if gift.GiftId == "11584" {
			isLottery = true
			// 抽奖
			giftMap := lottery(data.AuthorOpenId, gift.UserInfo.UserId, gift.GiftCount)
			ziLog.Gift(fmt.Sprintf("ksPushBasePayloay Lottery giftdata： %v，用户Id： %v, 用户名称： %v", giftMap, gift.UserInfo.UserId, gift.UserInfo.NickName), debug)
			for giftId, giftCount := range giftMap {
				score += giftToScoreMap[giftId] * float64(giftCount)
				if isGroup {
					groupGrpc.IsComment = false
					groupGrpc.OpenId = gift.UserInfo.UserId
					groupGrpc.GiftId = giftId
					groupGrpc.GiftNum = giftCount
					// go grpcSend(groupGrpc, 0)
				}
			}
			giftMapByte, _ := json.Marshal(giftMap)
			lotteryData := &pmsg.LotteryMsg{}
			lotteryData.OpenId = gift.UserInfo.UserId
			lotteryData.NickName = gift.UserInfo.NickName
			lotteryData.HeadImgUrl = gift.UserInfo.AvatarUrl
			lotteryData.LotteryMap = string(giftMapByte)
			lotteryData.Count = gift.GiftCount
			lotteryByte, _ := proto.Marshal(lotteryData)
			endSendData.Data = lotteryByte
			endSendData.TimeStamp = time.Now().UnixMilli()
			endSendData.PushType = "lottery"
			endSendDatabyte, _ := proto.Marshal(endSendData)
			if err := sse.SseSend(pmsg.MessageId_Lottery, sendUidList, endSendDatabyte); err != nil {
				ziLog.Error(fmt.Sprintf("ksPushBasePayloay 推送消息失败,用户Id： %v,用户名称： %v,err:  %v,内容为： %v", gift.UserInfo.UserId, gift.UserInfo.NickName, err, giftMap), debug)
			}

		} else {
			score = giftToScoreMap[gift.GiftId] * float64(gift.GiftCount)
			if isGroup {
				groupGrpc.IsComment = false
				groupGrpc.OpenId = gift.UserInfo.UserId
				groupGrpc.GiftId = gift.GiftId
				groupGrpc.GiftNum = gift.GiftCount
				// go grpcSend(groupGrpc, 0)
			}
		}
		if strings.HasPrefix(data.UniqueMessageId, "test_") {
			score = 0
		}
		//分数不为0时添加积分
		if score != 0 {
			go matchAddIntrage(data.RoomCode, openId, nickName, avatarUrl, score)
		}
		if !isLottery {
			// 格式化消息
			jData, err := json.Marshal(v)
			if err != nil {
				ziLog.Error(fmt.Sprintf("ksPushGiftSendPayloay json.Marshal err： %v， data： %v", err, v), debug)
				return
			}
			endSendData.Data = jData
			endSendData.TimeStamp = time.Now().UnixMilli()
			endSendDatabyte, _ := proto.Marshal(endSendData)
			// 推送消息
			if err := sse.SseSend(msgId, sendUidList, endSendDatabyte); err != nil {
				ziLog.Error(fmt.Sprintf("ksPushGiftSendPayloay 推送消息失败： %v， data： %v", err, v), debug)
			}
		}

	}
	// 推送消息验证
	if config.App.IsOnline {
		ksMsgAckSend(data.RoomCode, "cpClientReceive", ack)
	}
	sendAck := &pmsg.KsMsgAck{}
	sendAck.Data = &pmsg.KsMsgAckData{}
	sendAck.RoomId = data.RoomCode
	sendAck.Data.UniqueMessageId = data.UniqueMessageId
	sendAck.Data.PushType = data.PushType
	jData, err := proto.Marshal(sendAck)
	if err != nil {
		ziLog.Error(fmt.Sprintf("ksPushGiftSendPayloay sendack proto.Marshal err： %v", err), debug)
		return
	}
	uid := queryUidByOpenid(data.AuthorOpenId)
	if uid == "" {
		return
	}
	if err := sse.SseSend(pmsg.MessageId_MsgAckSend, []string{uid}, jData); err != nil {
		ziLog.Error(fmt.Sprintf("ksPushGiftSendPayloay sendack 推送消息失败： %v", err), debug)
	}

}

// 消息推送
func ksPushCommentPayloay(data KsCallbackDataStruct) {
	endSendData := platFormPool.Get().(*pmsg.PlatFormDataSend)
	defer platFormPool.Put(endSendData)
	endSendData.OpenId = data.AuthorOpenId
	endSendData.RoomId = data.RoomCode
	endSendData.PushType = data.PushType
	for _, v := range data.Payload {
		var (
			msgId pmsg.MessageId = pmsg.MessageId_LiveComment
		)
		// 格式化消息
		jData, err := json.Marshal(v)
		if err != nil {
			ziLog.Error(fmt.Sprintf("ksPushCommentPayloay json.Marshal err: %v", err), debug)
			return
		}
		sendUidList, _, _, _ := getUidListByOpenId(data.AuthorOpenId)
		log.Println("ksPushCommentPayloay sendUidList", sendUidList)
		if len(sendUidList) == 0 {
			ziLog.Error(fmt.Sprintf("ksPushBasePayloay queryRoomIdToUid nil， roomId: %v, openId, %v, 数据为： %v", data.RoomCode, data.AuthorOpenId, data), debug)
			return
		}
		endSendData.Data = jData
		endSendData.TimeStamp = time.Now().UnixMilli()
		endSendDatabyte, _ := proto.Marshal(endSendData)
		// 推送消息
		if err := sse.SseSend(msgId, sendUidList, endSendDatabyte); err != nil {
			ziLog.Error(fmt.Sprintf("ksPushCommentPayloay 推送消息失败: %v, data: %v", err, v), debug)
		}
	}

}

// 消息推送
func ksPushLiveLikePayloay(data KsCallbackDataStruct) {
	endSendData := platFormPool.Get().(*pmsg.PlatFormDataSend)
	defer platFormPool.Put(endSendData)
	endSendData.OpenId = data.AuthorOpenId
	endSendData.RoomId = data.RoomCode
	endSendData.PushType = data.PushType
	for _, v := range data.Payload {
		var (
			msgId pmsg.MessageId = pmsg.MessageId_liveLike
		)
		// 格式化消息
		jData, err := json.Marshal(v)
		if err != nil {
			ziLog.Error(fmt.Sprintf("ksPushLiveLikePayloay json.Marshal err: %v", err), debug)
			return
		}
		sendUidList, _, _, _ := getUidListByOpenId(data.AuthorOpenId)
		log.Println("ksPushLiveLikePayloay sendUidList", sendUidList)
		if len(sendUidList) == 0 {
			ziLog.Error(fmt.Sprintf("ksPushBasePayloay queryRoomIdToUid nil， roomId: %v, openId, %v, 数据为： %v", data.RoomCode, data.AuthorOpenId, data), debug)
			return
		}
		endSendData.Data = jData
		endSendData.TimeStamp = time.Now().UnixMilli()
		endSendDatabyte, _ := proto.Marshal(endSendData)
		// 推送消息
		if err := sse.SseSend(msgId, sendUidList, endSendDatabyte); err != nil {
			ziLog.Error(fmt.Sprintf("ksPushLiveLikePayloay 推送消息失败: %v, data: %v", err, v), debug)
		}
	}

}

// giftSend set ttl
func giftsendSet(roomId, uniqueMessageId string) {
	key := roomId + "giftSend"
	rdb.SAdd(key, uniqueMessageId)
	ttl, _ := rdb.TTL(key)
	if ttl < 0 {
		rdb.Expire(key, time.Hour*2)
	}
}
