package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	battlematchv1 "github.com/kongshui/danmu/battlematch/v1"

	"github.com/kongshui/danmu/model/pmsg"

	"google.golang.org/protobuf/proto"
)

// 接收匹配消息
func matchV1(msg *pmsg.MessageBody) error {
	data := &pmsg.MatchBattleV1ApplyMessage{}
	if err := proto.Unmarshal(msg.MessageData, data); err != nil {
		return errors.New("MatchV1 Unmarshal err: " + err.Error())
	}
	if data.GetIsAnonymous() {
		if err := battlematchv1.MatchBattleAnonymousSet(first_ctx, data.GetOpenId()); err != nil {
			ziLog.Error(fmt.Sprintf("设置匿名失败，匿名Id: %v, 匿名err: %v", data.GetOpenId(), err), debug)
		}
	}

	go MatchBattleV1(data.GetOpenId(), data.GetMatchNum())
	return nil
}

// 接收准备消息
func readyV1(msg *pmsg.MessageBody) error {
	data := &pmsg.MatchBattleV1ReadyMessage{}
	if err := proto.Unmarshal(msg.MessageData, data); err != nil {
		MatchErrorV1(data.GetOpenIdList(), "readyV1 Unmarshal err: "+err.Error(), data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_ProtoUnmarshalError))
		return errors.New("ReadyV1 Unmarshal err: " + err.Error())
	}
	for _, openId := range data.GetOpenIdList() {
		// 如果是自己发送的消息，则不需要推送给自己
		if openId == data.GetOpenId() {
			continue
		}
		// 查询uid
		uid := queryUidByOpenid(openId)
		if uid == "" {
			MatchErrorV1(data.GetOpenIdList(), "ReadyV1 queryOpenidFromUidStr nil", data.GetMatchBattleRoomId(),
				int32(pmsg.ErrorStatus_SendMessageError))
			return errors.New("ReadyV1 queryOpenidFromUidStr nil")
		}
		// 发送给前端
		if err := sendMessage(pmsg.MessageId_MatchBattleV1Ready, []string{uid}, msg.GetMessageData()); err != nil { // 推送消息给其他用户
			// 发送消息失败
			MatchErrorV1(data.GetOpenIdList(), "ReadyV1 pushDownLoadMessage err: "+err.Error(), data.GetMatchBattleRoomId(),
				int32(pmsg.ErrorStatus_SendMessageError))
			return errors.New("ReadyV1 pushDownLoadMessage err: " + err.Error() + "data: " + data.String())
		}
	}
	return nil
}

// 接收准备返回消息
func readyV1Ack(msg *pmsg.MessageBody) error {
	data := &pmsg.MatchBattleV1ReadyAckMessage{}
	if err := proto.Unmarshal(msg.MessageData, data); err != nil {
		MatchErrorV1(data.GetOpenIdList(), "readyV1Ack Unmarshal err: "+err.Error(), data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_ProtoUnmarshalError))
		return errors.New("readyV1Ack Unmarshal err: " + err.Error())
	}
	if err := battlematchv1.MatchGroupStatusSet(first_ctx, data.GetMatchBattleRoomId(), match_battle_status_set_time); err != nil {
		if err.Error() == "equal" { // 状态相同
			return nil
		}
		MatchErrorV1(data.GetOpenIdList(), "readyV1Ack MatchGroupStatusSet err: "+err.Error(), data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_SetMatchStatusError))
		return errors.New("readyV1Ack MatchGroupStatusSet err: " + err.Error())
	}
	sData := &pmsg.MatchBattleV1TimeCheckMessage{}
	sData.MatchBattleRoomId = data.GetMatchBattleRoomId()
	sData.TimeStamp = time.Now().UnixMilli()
	sData.CheckTimeStamp = time.Now().Unix() + 15 // 设置开始时间戳为当前时间戳加10秒
	sDataByte, err := proto.Marshal(sData)
	if err != nil {
		MatchErrorV1(data.GetOpenIdList(), "readyV1Ack Marshal err: "+err.Error(), data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_ProtoMarshalError))
		return errors.New("readyV1Ack Marshal err: " + err.Error())
	}
	// 查询uid
	uid := queryUidByOpenid(data.GetOpenId())
	if uid == "" {
		MatchErrorV1(data.GetOpenIdList(), "ReadyV1 queryOpenidFromUidStr nil", data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_SendMessageError))
		return errors.New("ReadyV1 queryOpenidFromUidStr nil")
	}
	sendUidList, _, _, _ := getUidListByGroupId(data.GetMatchBattleRoomId())
	if len(sendUidList) == 0 {
		MatchErrorV1(data.GetOpenIdList(), "ReadyV1 queryOpenidFromUidStr nil", data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_SendMessageError))
		return errors.New("ReadyV1 queryOpenidFromUidStr nil")
	}
	if err := sendMessage(pmsg.MessageId_MatchBattleV1TimeCheck, sendUidList, sDataByte); err != nil { // 推送消息给其他用户
		MatchErrorV1(data.GetOpenIdList(), "ReadyV1 pushDownLoadMessage err: "+err.Error(), data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_SetMatchStatusError))
		return errors.New("ReadyV1 pushDownLoadMessage err: " + err.Error() + "data: " + data.String())
	}
	// 设置匹配时间
	if err := battlematchv1.MatchBattleTimeSet(first_ctx, data.GetMatchBattleRoomId(), sData.GetCheckTimeStamp()); err != nil {
		MatchErrorV1(data.GetOpenIdList(), "ReadyV1 MatchBattleTimeSet err: "+err.Error(), data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_MatchBattleV1SetStartTimeError))
	}
	return nil
}

// 互相接收彼此匹配时间
func matchV1TimeAck(msg *pmsg.MessageBody) error {
	data := &pmsg.MatchBattleV1TimeCheckAckMessage{}
	if err := proto.Unmarshal(msg.MessageData, data); err != nil {
		MatchErrorV1(data.GetOpenIdList(), "matchTimeAck Unmarshal err: "+err.Error(), data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_ProtoUnmarshalError))
		return errors.New("matchTimeAck Unmarshal err: " + err.Error())
	}
	for _, openId := range data.GetOpenIdList() {
		// 如果是自己发送的消息，则不需要推送给自己
		if openId == data.GetOpenId() {
			continue
		}
		// 查询uid
		uid := queryUidByOpenid(openId)
		if uid == "" {
			MatchErrorV1(data.GetOpenIdList(), "matchTimeAck queryOpenidFromUidStr nil", data.GetMatchBattleRoomId(),
				int32(pmsg.ErrorStatus_SendMessageError))
			return errors.New("matchTimeAck queryOpenidFromUidStr nil")
		}

		if err := sendMessage(pmsg.MessageId_MatchBattleV1TimeCheckAck, []string{uid}, msg.GetMessageData()); err != nil { // 推送消息给其他用户
			// 发送消息失败
			MatchErrorV1(data.GetOpenIdList(), "matchTimeAck pushDownLoadMessage err: "+err.Error(), data.GetMatchBattleRoomId(),
				int32(pmsg.ErrorStatus_SendMessageError))
			return errors.New("matchTimeAck pushDownLoadMessage err: " + err.Error() + "data: " + data.String())
		}
	}
	return nil
}

// 匹配确认消息，场景准备好以后发送给对方，对方收到后再发送给自己，然后再发送给其他用户
func matchV1Confirm(msg *pmsg.MessageBody) error {
	data := &pmsg.MatchBattleV1StartConfirmMessage{}
	if err := proto.Unmarshal(msg.MessageData, data); err != nil {
		MatchErrorV1(data.GetOpenIdList(), "matchV1Confirm Unmarshal err: "+err.Error(), data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_ProtoUnmarshalError))
		return errors.New("MatchV1Confirm Unmarshal err: " + err.Error())
	}
	// 转发消息给其他用户
	for _, openId := range data.GetOpenIdList() {
		// 如果是自己发送的消息，则不需要推送给自己
		if openId == data.GetOpenId() {
			continue
		}
		// 查询uid
		uid := queryUidByOpenid(openId)
		if uid == "" {
			MatchErrorV1(data.GetOpenIdList(), "matchV1Confirm queryOpenidFromUidStr nil", data.GetMatchBattleRoomId(),
				int32(pmsg.ErrorStatus_SendMessageError))
			return errors.New("matchV1Confirm queryOpenidFromUidStr nil")
		}
		if err := sendMessage(pmsg.MessageId_MatchBattleV1StartConfirm, []string{uid}, msg.GetMessageData()); err != nil { // 推送消息给其他用户
			// 发送消息失败
			MatchErrorV1(data.GetOpenIdList(), "matchV1Confirm pushDownLoadMessage err: "+err.Error(), data.GetMatchBattleRoomId(),
				int32(pmsg.ErrorStatus_SendMessageError))
			return errors.New("matchV1Confirm pushDownLoadMessage err: " + err.Error() + "data: " + data.String())
		}
	}
	return nil
}

// 匹配确认放回消息
func matchV1ConfirmAck(msg *pmsg.MessageBody) error {
	data := &pmsg.MatchBattleV1StartConfirmAckMessage{}
	if err := proto.Unmarshal(msg.MessageData, data); err != nil {
		MatchErrorV1(data.GetOpenIdList(), "matchV1ConfirmAck Unmarshal err: "+err.Error(), data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_ProtoUnmarshalError))
		return errors.New("matchV1ConfirmAck Unmarshal err: " + err.Error())
	}
	if err := battlematchv1.MatchGroupStatusSet(first_ctx, data.GetMatchBattleRoomId(), match_battle_status_Confirm); err != nil {
		if err.Error() == "equal" { // 状态相同
			return nil
		}
		MatchErrorV1(data.GetOpenIdList(), "matchV1ConfirmAck MatchGroupStatusSet err: "+err.Error(), data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_SetMatchStatusError))
		return errors.New("matchV1ConfirmAck MatchGroupStatusSet err: " + err.Error())
	}
	sdata := &pmsg.MatchBattleV1SendRoundIdMessage{}
	sdata.MatchBattleRoomId = data.GetMatchBattleRoomId()
	sdata.RoundId = time.Now().UnixMicro()
	sdata.TimeStamp = time.Now().UnixMilli()
	// 序列化数据
	sdataByte, err := proto.Marshal(sdata)
	if err != nil {
		MatchErrorV1(data.GetOpenIdList(), "matchV1ConfirmAck Marshal err: "+err.Error(), data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_ProtoMarshalError))
		return errors.New("matchV1ConfirmAck Marshal err: " + err.Error())
	}
	// 设置对局Id
	if err := battlematchv1.MatchV1GroupRoundIdSet(first_ctx, data.GetMatchBattleRoomId(), sdata.GetRoundId()); err != nil {
		MatchErrorV1(data.GetOpenIdList(), "matchV1ConfirmAck MatchV1GroupRoundIdSet err: "+err.Error(), data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_SetMatchStatusError))
		return errors.New("matchV1ConfirmAck MatchV1GroupRoundIdSet err: " + err.Error())
	}
	// 通过openId获取group用户
	for _, openId := range data.GetOpenIdList() {
		// 查询roomId
		roomId := QueryRoomIdInterconvertAnchorOpenId(openId)
		if roomId != "" {
			liveCurrentRoundAdd(roomId, sdata.GetRoundId()) // 添加当前对局Id
		}
	}
	// 查询uid
	sendUidList, _, _, _ := getUidListByGroupId(data.GetMatchBattleRoomId())
	if len(sendUidList) == 0 {
		MatchErrorV1(data.GetOpenIdList(), "askMatchV1RoundId queryOpenidFromUidStr nil", data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_SendMessageError))
		return errors.New("askMatchV1RoundId queryOpenidFromUidStr nil")
	}
	// 推送消息给所有用户
	if err := sendMessage(pmsg.MessageId_MatchBattleV1SendRoundId, sendUidList, sdataByte); err != nil { // 推送消息给其他用户
		// 发送消息失败
		MatchErrorV1(data.GetOpenIdList(), "askMatchV1RoundId pushDownLoadMessage err: "+err.Error(), data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_SendMessageError))
		return errors.New("askMatchV1RoundId pushDownLoadMessage err: " + err.Error() + "data: " + data.String())
	}
	// 计算roomId
	roomId := QueryRoomIdInterconvertAnchorOpenId(data.GetOpenIdList()[0]) // 记录对方的roomId
	if roomId == "" {
		return errors.New("MatchV1Start roomId2 err:  is nil")
	}
	roomId2 := QueryRoomIdInterconvertAnchorOpenId(data.GetOpenIdList()[1]) // 记录对方的roomId
	if roomId2 == "" {
		return errors.New("MatchV1Start roomId2 err:  is nil")
	}
	// 推送配置
	if ok, _ := battlematchv1.MatchBattleAnonymousGetByGroupId(first_ctx, data.GetMatchBattleRoomId()); !ok {
		if err := TwoConnect("config", roomId, roomId2, data.GetMatchBattleRoomId()); err != nil {
			ziLog.Error("matchV1ConfirmAck twoConnect err:"+err.Error(), debug)
		}
		if err := TwoConnect("config", roomId2, roomId, data.GetMatchBattleRoomId()); err != nil {
			ziLog.Error("matchV1ConfirmAck twoConnect err:"+err.Error(), debug)
		}
	}
	return nil
}

// // 客户端收到对局信息返回消息
// func matchV1RoundIdAck(msg *pmsg.MessageBody) error {
// 	data := &pmsg.MatchBattleV1SendRoundIdAckMessage{}
// 	err := proto.Unmarshal(msg.MessageData, data)
// 	if err != nil {
// 		return errors.New("askMatchV1RoundIdAck Unmarshal err: " + err.Error())
// 	}
// 	sData := &pmsg.MatchBattleV1SendUserInfoMessage{}
// 	sData.MatchBattleRoomId = data.GetMatchBattleRoomId()
// 	sData.RoundId = data.GetRoundId()
// 	sData.TimeStamp = time.Now().UnixMilli()
// 	sData.AvatarUrl = data.GetAvatarUrl()
// 	sData.NickName = data.GetNickName()
// 	sData.OpenId = data.GetOpenId()
// 	for _, openId := range data.GetOpenIdList() { // 通过openId获取group用户
// 		if openId == data.GetOpenId() { // 如果是自己发送的消息，则不需要推送给自己
// 			continue
// 		}
// 		sDataByte, err := proto.Marshal(sData)
// 		if err != nil {
// 			return errors.New("askMatchV1RoundIdAck Marshal err: " + err.Error() + "data: " + data.String())
// 		}
// 		// 查询uid
// 		uid := queryUidInterconvertOpenid(openId)
// 		if err := pushDownLoadMessage(uint32(pmsg.MessageId_MatchBattleV1SendUserInfo), pmsg.MessageId_MatchBattleV1SendUserInfo.String(),
// 			uid, sDataByte); err != nil { // 推送消息给其他用户
// 			// 发送消息失败
// 			return errors.New("askMatchV1RoundIdAck pushDownLoadMessage err: " + err.Error() + "data: " + data.String())
// 		}
// 	}
// 	return nil
// }

// askRoundId 询问对局Id
func askMatchV1RoundId(msg *pmsg.MessageBody) error {
	data := &pmsg.MatchBattleV1AskRoundIdMessage{}
	if err := proto.Unmarshal(msg.MessageData, data); err != nil {
		return errors.New("askMatchV1RoundId Unmarshal err: " + err.Error())
	}
	// 查询对局Id
	roundId, err := battlematchv1.MatchV1GroupRoundIdGet(first_ctx, data.GetMatchBattleRoomId())
	if err != nil {
		return errors.New("askMatchV1RoundId MatchV1GroupRoundIdGet err: " + err.Error())
	}
	// 序列化数据
	sdata := &pmsg.MatchBattleV1AskRoundIdAckMessage{}
	sdata.MatchBattleRoomId = data.GetMatchBattleRoomId()
	sdata.RoundId = roundId
	sdata.TimeStamp = time.Now().UnixMilli()
	sdataByte, err := proto.Marshal(sdata)
	if err != nil {
		return errors.New("askMatchV1RoundId Marshal err: " + err.Error())
	}
	// 查询uid
	uid := queryUidByOpenid(data.GetOpenId())
	if uid == "" {
		return errors.New("askMatchV1RoundId queryOpenidFromUidStr nil")
	}
	// 推送消息给用户
	if err := sendMessage(pmsg.MessageId_MatchBattleV1AskRoundIdAck, []string{uid}, sdataByte); err != nil { // 推送消息给其他用户
		// 发送消息失败
		return errors.New("askMatchV1RoundId pushDownLoadMessage err: " + err.Error() + "data: " + data.String())
	}
	return nil
}

// 游戏开始
func matchV1Start(msg *pmsg.MessageBody) error {
	data := &pmsg.MatchBattleV1StartMessage{}
	if err := proto.Unmarshal(msg.MessageData, data); err != nil {
		return errors.New("MatchV1Start Unmarshal err: " + err.Error())
	}
	status, err := battlematchv1.MatchGroupStatusGet(first_ctx, data.GetMatchBattleRoomId())
	if err != nil {
		return fmt.Errorf("matchV1Start MatchGroupStatusGet error, err: %v", err)
	}
	if status != match_battle_status_Confirm {
		return fmt.Errorf("matchV1Start MatchGroupStatusGet error, get status err, status: %v", status)
	}
	// 绑定roundId
	ksSyncGameStatus(SyncGameStatusStruct{
		RoomId:       data.GetRoomId(),
		RoundId:      data.GetRoundId(),
		AnchorOpenId: data.GetOpenId(),
	}, "start", true)
	if data.OpenId != data.GetOpenIdList()[0] {
		if interactive != nil {
			interactive(data.GetRoomId(), strconv.FormatInt(data.GetRoundId(), 10), 2)
		}
		return nil
	}
	if interactive != nil {
		interactive(data.GetRoomId(), strconv.FormatInt(data.GetRoundId(), 10), 1)
	}
	if err := battlematchv1.MatchGroupStatusSet(first_ctx, data.GetMatchBattleRoomId(), match_battle_status_start); err != nil {
		if err.Error() == "equal" { // 状态相同
			return nil
		}
		return errors.New("matchV1ConfirmAck MatchGroupStatusSet err: " + err.Error())
	}
	// 计算roomId
	roomId := QueryRoomIdInterconvertAnchorOpenId(data.GetOpenIdList()[0]) // 记录对方的roomId
	if roomId == "" {
		return errors.New("MatchV1Start roomId2 err:  is nil")
	}
	roomId2 := QueryRoomIdInterconvertAnchorOpenId(data.GetOpenIdList()[1]) // 记录对方的roomId
	if roomId2 == "" {
		return errors.New("MatchV1Start roomId2 err:  is nil")
	}
	// 发送消息
	sData := &pmsg.MatchBattleV1StartAckMessage{}
	sData.MatchBattleRoomId = data.GetMatchBattleRoomId()
	sData.OpenIdList = data.GetOpenIdList()
	sData.RoundId = data.GetRoundId()
	sData.TimeStamp = time.Now().UnixMilli()
	sDataByte, err := proto.Marshal(sData)
	if err != nil {
		MatchErrorV1(data.GetOpenIdList(), "matchV1Start send Marshal err: "+err.Error(), data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_ProtoMarshalError))
		return errors.New("matchV1Start send Marshal err: " + err.Error())
	}
	// 查询uid
	sendUidList, _, _, _ := getUidListByGroupId(data.GetMatchBattleRoomId())
	if len(sendUidList) == 0 {
		MatchErrorV1(data.GetOpenIdList(), "askMatchV1RoundId queryOpenidFromUidStr nil", data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_SendMessageError))
		return errors.New("askMatchV1RoundId queryOpenidFromUidStr nil")
	}
	// 发送消息
	if err := sendMessage(pmsg.MessageId_MatchBattleV1StartAck, sendUidList, sDataByte); err != nil {
		MatchErrorV1(data.GetOpenIdList(), "matchV1Start send pushDownLoadMessage err: "+err.Error(), data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_ProtoMarshalError))
		return errors.New("matchV1Start send pushDownLoadMessage err: " + err.Error())
	}
	//
	// matchbattlecal.AddToMonitorGroup(data.GetMatchBattleRoomId())
	// 是否匿名，连线暂时关闭
	// if ok, _ := battlematchv1.MatchBattleAnonymousGetByGroupId(first_ctx, data.GetMatchBattleRoomId()); !ok {
	// 	// log.Println("matchV1Start enter twoConnect set")
	// 	// if err := twoConnect("config", data.GetRoomId(), data.GetRoomId(), roomId2, data.GetMatchBattleRoomId()); err != nil {
	// 	// 	ziLog.Error( "matchV1ConfirmAck twoConnect err:"+err.Error(), debug)
	// 	// }
	// 	twoConnect("start", roomId, roomId2, data.GetMatchBattleRoomId())
	// }
	return nil
}

// 游戏结束
func matchV1End(msg *pmsg.MessageBody) error {
	data := &pmsg.MatchBattleV1EndMessage{}
	if err := proto.Unmarshal(msg.MessageData, data); err != nil {
		MatchErrorV1(data.GetOpenIdList(), "matchV1ConfirmAck Marshal err: "+err.Error(), data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_ProtoMarshalError))
		return errors.New("matchV1ConfirmAck Marshal err: " + err.Error())
	}
	ksSyncGameStatus(SyncGameStatusStruct{
		RoomId:       data.GetRoomId(),
		RoundId:      data.GetRoundId(),
		AnchorOpenId: data.GetOpenId(),
	}, "stop", true)
	if err := battlematchv1.MatchGroupStatusSet(first_ctx, data.GetMatchBattleRoomId(), match_battle_status_stop); err != nil {
		if err.Error() == "equal" { // 状态相同
			return nil
		}
		return errors.New("matchV1ConfirmAck MatchGroupStatusSet err: " + err.Error())
	}
	// pk断线
	roomid := QueryRoomIdInterconvertAnchorOpenId(data.GetOpenIdList()[0])
	roomid2 := QueryRoomIdInterconvertAnchorOpenId(data.GetOpenIdList()[0])
	TwoConnect("stop", roomid, roomid2, data.GetMatchBattleRoomId())
	// 结束
	err := battlematchv1.UnregisterBattleV1ByGroupId(first_ctx, data.GetMatchBattleRoomId())
	if err != nil {
		return errors.New("matchV1ConfirmAck UnregisterBattleV1ByGroupId err: " + err.Error())
	}
	return nil
}

// 传输数据
func matchV1SyncData(msg *pmsg.MessageBody) error {
	data := &pmsg.MatchBattleV1SyncMessage{}
	if err := proto.Unmarshal(msg.MessageData, data); err != nil {
		return errors.New("MatchV1SyncData Unmarshal err: " + err.Error())
	}
	// 传输数据给其他用户
	for _, openId := range data.GetOpenIdList() {
		// 如果是自己发送的消息，则不需要推送给自己
		if openId == data.GetOpenId() {
			continue
		}
		if battlematchv1.QueryOpenIdInMatchDisconnect(first_ctx, data.GetMatchBattleRoomId(), openId) {
			continue
		}
		// 查询uid
		uid := queryUidByOpenid(openId)
		if err := sendMessage(pmsg.MessageId_MatchBattleV1Sync, []string{uid}, msg.GetMessageData()); err != nil { // 推送消息给其他用户
			// 发送消息失败
			return errors.New("matchV1SyncData pushDownLoadMessage err: " + err.Error() + "data: " + data.String())
		}
	}
	return nil
}

// 传输数据返回
func matchV1SyncDataAck(msg *pmsg.MessageBody) error {
	data := &pmsg.MatchBattleV1SyncAckMessage{}
	if err := proto.Unmarshal(msg.MessageData, data); err != nil {
		return errors.New("MatchV1SyncData Unmarshal err: " + err.Error())
	}
	// 传输数据给其他用户
	// 查询uid
	uid := queryUidByOpenid(data.GetOpenId())
	if err := sendMessage(pmsg.MessageId_MatchBattleV1Sync, []string{uid}, msg.GetMessageData()); err != nil { // 推送消息给其他用户
		// 发送消息失败
		return errors.New("matchV1SyncData pushDownLoadMessage err: " + err.Error() + "data: " + data.String())
	}
	return nil
}

// // 请求对局主播信息
// func askMatchV1AnchorInfo(msg *pmsg.MessageBody) error {
// 	data := &pmsg.MatchBattleV1AskUserInfoMessage{}
// 	err := proto.Unmarshal(msg.MessageData, data)
// 	if err != nil {
// 		return errors.New("askMatchV1AnchorInfo Unmarshal err: " + err.Error())
// 	}
// 	// 通过openId查看玩家信息
// 	user, err := userInfoGet(data.GetAskOpenId())
// 	if err != nil {
// 		return errors.New("askMatchV1AnchorInfo userInfoGet err: " + err.Error())
// 	}
// 	sData := &pmsg.MatchBattleV1AskUserInfoAckMessage{}
// 	sData.MatchBattleRoomId = data.GetMatchBattleRoomId()
// 	sData.AvatarUrl = user.AvatarUrl
// 	sData.NickName = user.NickName
// 	sData.OpenId = data.GetAskOpenId()
// 	sData.TimeStamp = time.Now().UnixMilli()
// 	sDataByte, err := proto.Marshal(sData)
// 	if err != nil {
// 		return errors.New("askMatchV1AnchorInfo Marshal err: " + err.Error())
// 	}
// 	// 查询uid
// 	uid := queryUidInterconvertOpenid(data.GetOpenId())
// 	if err := pushDownLoadMessage(uint32(pmsg.MessageId_MatchBattleV1AskUserInfoAck), pmsg.MessageId_MatchBattleV1AskUserInfoAck.String(),
// 		uid, sDataByte); err != nil {
// 		// 发送消息失败
// 		return errors.New("askMatchV1AnchorInfo pushDownLoadMessage err: " + err.Error() + "data: " + data.String())
// 	}
// 	return nil
// }

// 匹配上报信息
func matchV1DataUpload(msg *pmsg.MessageBody) error {
	data := &pmsg.MatchBattleV1Upload{}
	if err := proto.Unmarshal(msg.MessageData, data); err != nil {
		return errors.New("matchV1DataUpload Unmarshal err: " + err.Error())
	}
	// if err := battlematchv1.MatchGroupStatusSet(first_ctx, data.GetMatchBattleRoomId(), match_battle_status_roundup); err != nil {
	// 	if err.Error() == "equal" { // 状态相同
	// 		return nil
	// 	}
	// 	return errors.New("matchV1DataUpload MatchGroupStatusSet err: " + err.Error())
	// }
	//数据上报 =====================================================
	roundData := RoundUploadStruct{}
	if len(data.GetGroupUserList()) == 0 {
		return nil
	}
	for _, v := range data.GroupResultList {
		if v.GetGroupId() == "Left" && v.GetResult() == 1 {
			mysql.SetPkScores(data.GetOpenIdList()[0], true)
		} else if v.GetGroupId() == "Right" && v.GetResult() == 1 {
			mysql.SetPkScores(data.GetOpenIdList()[1], true)
		} else if v.GetGroupId() == "Left" && v.GetResult() != 1 {
			mysql.SetPkScores(data.GetOpenIdList()[0], false)
		} else if v.GetGroupId() == "Right" && v.GetResult() != 1 {
			mysql.SetPkScores(data.GetOpenIdList()[1], false)
		} else {
			ziLog.Error("matchV1DataUpload other err", debug)
		}
		roundData.GroupResultList = append(roundData.GroupResultList, GroupResultList{GroupId: v.GroupId, Result: int(v.Result)})
	}
	if !battlematchv1.MatchGroupEnd(first_ctx, data.GetMatchBattleRoomId()) {
		return nil
	}
	for _, v := range data.GroupUserList {
		roundData.GroupUserList = append(roundData.GroupUserList, UserUploadScoreStruct{OpenId: v.OpenId, GroupId: v.GroupId, Score: v.Score})
	}
	if err := MatchBattleSetWinnerScore(data.GetMatchBattleRoomId(), roundData); err != nil {
		ziLog.Error(fmt.Sprintf("matchV1DataUpload matchBattleSetWinnerScore err: %v", err), debug)
		return fmt.Errorf("matchV1DataUpload matchBattleSetWinnerScore err: %v", err)
	}
	time.Sleep(5 * time.Second)
	return battlematchv1.MatchGroupStatusDel(first_ctx, data.GetMatchBattleRoomId())
}

// 取消匹配
func matchV1Cancel(msg *pmsg.MessageBody) error {
	// 数据初始化
	data := &pmsg.MatchBattleV1CancelMessage{}
	if err := proto.Unmarshal(msg.MessageData, data); err != nil {
		return errors.New("MatchV1Cancel Unmarshal err: " + err.Error())
	}
	// 取消数据初始化
	cancelData := &pmsg.MatchBattleV1CancelAckMessage{}
	cancelData.IsSuccess = true
	cancelData.OpenId = data.GetOpenId()
	cancelData.RoomId = data.GetRoomId()
	// 判断情况
	if battlematchv1.IsInResiter(first_ctx, data.GetMatchNum(), data.GetOpenId()) {
		if err := battlematchv1.UnregisterBattleByOpenId(first_ctx, data.GetMatchNum(), data.GetOpenId()); err != nil {
			cancelData.IsSuccess = false
			cancelData.Error = "MatchV1Cancel IsInResiter 取消失败，openId： " + data.GetOpenId()
			ziLog.Error("MatchV1Cancel IsInResiter 取消失败，openId： "+data.GetOpenId(), debug)
		}
	} else {
		if err := battlematchv1.AddToCancelMatchV1Battle(first_ctx, data.MatchNum, data.GetOpenId()); err != nil {
			cancelData.IsSuccess = true
			cancelData.Error = "MatchV1Cancel 取消失败，openId： " + data.GetOpenId()
			ziLog.Error("MatchV1Cancel 取消失败，openId： "+data.GetOpenId(), debug)
		}
	}
	cancelData.TimeStamp = time.Now().UnixMilli()
	cancelDataByte, _ := proto.Marshal(cancelData)
	uid := queryUidByOpenid(data.GetOpenId())
	if err := sendMessage(pmsg.MessageId_MatchBattleV1CancelAck, []string{uid}, cancelDataByte); err != nil { // 推送消息给其他用户
		// 发送消息失败
		return errors.New("matchV1SyncData pushDownLoadMessage err: " + err.Error() + "data: " + data.String())
	}
	return nil
}

// 增加节点积分
func MatchBattleAddIntegral(msg *pmsg.MessageBody) error {
	data := &pmsg.MatchBattleAddIntegralMessage{}
	if err := proto.Unmarshal(msg.MessageData, data); err != nil {
		return errors.New("matchBattleAddIntegral Unmarshal err: " + err.Error())
	}
	if _, err := addIntegralByNode(data.GetMatchBattleRoomId(), 500); err != nil {
		return errors.New("addIntegral err: " + err.Error())
	}
	return nil
}

// 匹配使用连胜币
func MatchBattleUseStreamCoin(msg *pmsg.MessageBody) error {
	data := &pmsg.MatchBattleUseWinnerStreamCoinMessage{}
	if err := proto.Unmarshal(msg.MessageData, data); err != nil {
		return errors.New("matchBattleUseStreamCoin Unmarshal err: " + err.Error())
	}
	comment, ok := commentTogiftId[data.GetComment()]
	if !ok {
		return errors.New("matchBattleUseStreamCoin use WinStreamCoin err: can not find giftId ")
	}
	useCoin, ok := commentToCoin[comment]
	if !ok {
		return errors.New("matchBattleUseStreamCoin use WinStreamCoin err: can not find comment ")
	}
	sData := &pmsg.MatchBattleUseWinnerStreamCoinAckMessage{}
	sData.IsUse = true
	coin, err := deleteUserWinStreamCoin(data.OpenId, useCoin)
	if err != nil {
		sData.IsUse = false
	}
	sData.GiftId = data.GetComment()
	sData.OpenId = data.GetOpenId()
	sData.MatchBattleRoomId = data.GetMatchBattleRoomId()
	sData.OpenIdList = data.GetOpenIdList()
	sData.UseSide = data.GetUseSide()
	sData.WinningStreamCoin = coin
	sData.TimeStamp = time.Now().UnixMilli()
	sData.RoomId = data.GetRoomId()
	sendUidList := make([]string, 0)
	if sData.GetIsUse() {
		sendUidList, _, _, _ = getUidListByGroupId(data.GetMatchBattleRoomId())
	} else {
		uid := queryRoomIdToUid(data.GetRoomId())
		sendUidList = append(sendUidList, uid)
	}
	if len(sendUidList) == 0 {
		return errors.New("matchBattleUseStreamCoin get uid is nil")
	}
	sDataByte, _ := proto.Marshal(sData)
	if err := sendMessage(pmsg.MessageId_MatchBattleUseWinnerStreamCoinAck, sendUidList, sDataByte); err != nil { // 推送消息给其他用户
		// 发送消息失败
		return errors.New("matchBattleUseStreamCoin pushDownLoadMessage err: " + err.Error() + "data: " + data.String())
	}
	return nil
}

// 瓜分连胜币
func MatchAddStreamCoin(msg *pmsg.MessageBody) error {
	data := &pmsg.MatchBattleAddStreamCoinMessage{}
	if err := proto.Unmarshal(msg.MessageData, data); err != nil {
		return errors.New("matchAddStreamCoin Unmarshal err: " + err.Error())
	}
	if !battlematchv1.MatchAddStreamCoinStatusSet(first_ctx, data.GetMatchBattleRoomId(), data.GetType()) {
		return nil
	}
	// log.Println(data.String(), 1111111111)
	addWinningStreamCoin(data.GetUserList())
	return nil
}

// 开始确认返回
func matchStartGamedConfirmAck(msg *pmsg.MessageBody) error {
	data := &pmsg.MatchBattleStartGamedConfirAckMessage{}
	if err := proto.Unmarshal(msg.MessageData, data); err != nil {
		MatchErrorV1(data.GetOpenIdList(), "MatchStartGamedConfirm Marshal err: "+err.Error(), data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_ProtoMarshalError))
		return errors.New("MatchStartGamedConfirm Unmarshal err: " + err.Error())
	}
	// 转发消息给其他用户
	for _, openId := range data.GetOpenIdList() {
		// 如果是自己发送的消息，则不需要推送给自己
		if openId == data.GetOpenId() {
			continue
		}
		// 查询uid
		uid := queryUidByOpenid(openId)
		if uid == "" {
			MatchErrorV1(data.GetOpenIdList(), "MatchStartGamedConfirm queryOpenidFromUidStr nil", data.GetMatchBattleRoomId(),
				int32(pmsg.ErrorStatus_SendMessageError))
			return errors.New("MatchStartGamedConfirm queryOpenidFromUidStr nil")
		}
		if err := sendMessage(pmsg.MessageId_MatchBattleStartGamedConfirmAck, []string{uid}, msg.GetMessageData()); err != nil { // 推送消息给其他用户
			// 发送消息失败
			MatchErrorV1(data.GetOpenIdList(), "MatchStartGamedConfirm pushDownLoadMessage err: "+err.Error(), data.GetMatchBattleRoomId(),
				int32(pmsg.ErrorStatus_SendMessageError))
			return errors.New("MatchStartGamedConfirm pushDownLoadMessage err: " + err.Error() + "data: " + data.String())
		}
	}
	return nil
}

// 设置匹配组掉线
func matchBattleQuitWithError(msg *pmsg.MessageBody) error {
	data := &pmsg.MatchBattleQuitWithErrorMessage{}
	if err := proto.Unmarshal(msg.MessageData, data); err != nil {
		return errors.New("matchBattleQuitWithError Unmarshal err: " + err.Error())
	}
	ziLog.Error(fmt.Sprintf("matchBattleQuitWithError errMsg: %v", data.GetErrorMsg()), debug)
	if err := battlematchv1.DisconnectMatchRegister(first_ctx, data.GetOpenId()); err != nil {
		ziLog.Error(fmt.Sprintf("matchBattleQuitWithError DisconnectMatchRegister err: %v", err), debug)
		return fmt.Errorf("matchBattleQuitWithError DisconnectMatchRegister err: %v", err)
	}

	return nil
}
