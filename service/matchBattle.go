package service

import (
	"errors"
	"fmt"
	"log"
	"math"
	"path"
	"time"

	battlematchv1 "github.com/kongshui/danmu/battlematch/v1"

	"github.com/kongshui/danmu/model/pmsg"
	"github.com/kongshui/danmu/sse"

	"google.golang.org/protobuf/proto"
)

func MatchBattleV1(openId, matchNum string) {
	data := &pmsg.MatchBattleV1ApplyAckMessage{}
	matchOpenIdList, groupId, err := battlematchv1.MatchV1Battle(first_ctx, openId, matchNum)
	//匹配失败
	if err != nil {
		if err == errors.New("用户已取消匹配") {
			return
		}
		data.IsMatch = false
		ziLog.Error(fmt.Sprintf("MatchBattleV1 匹配失败, err: %v", err), debug)

	}
	if len(matchOpenIdList) == 0 || groupId == "" {
		return
	}
	//匹配成功
	data.MatchBattleRoomId = groupId
	data.OpenIdList = matchOpenIdList
	data.IsMatch = true
	ziLog.Info(fmt.Sprintf("MatchBattleV1 匹配成功, openId: %v,matchOpenIdList: %s, groupId: %s", openId, matchOpenIdList, groupId), debug)
	if err := battlematchv1.MatchGroupStatusSet(first_ctx, data.GetMatchBattleRoomId(), match_battle_status_ready); err != nil {
		if err == errors.New("equal") { // 状态相同
			return
		}
		// 匹配失败
		MatchErrorV1(data.GetOpenIdList(), "matchBattleV1 MatchGroupStatusSet err: "+err.Error(), data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_SetMatchStatusError))
		return

	}
	// 匹配成功
	for _, openIdStr := range data.GetOpenIdList() {
		user, err := UserInfoGet(openIdStr, false)
		if err != nil {
			ziLog.Error(fmt.Sprintf("MatchBattleV1 查询用户信息失败, err: %v, group: %v", err, groupId), debug)
		}
		ok, _ := battlematchv1.MatchBattleAnonymousGet(first_ctx, openIdStr)
		data.UsersInfo = append(data.UsersInfo, &pmsg.MatchBattleV1ApplyAckMessageUserInfo{
			OpenId:      openIdStr,
			NickName:    user.NickName,
			AvatarUrl:   user.AvatarUrl,
			IsAnonymous: ok,
		})
	}
	data.TimeStamp = time.Now().UnixMilli()
	sData, err := proto.Marshal(data)
	if err != nil {
		ziLog.Error(fmt.Sprintf("MatchBattleV1 序列化失败, err: %v", err), debug)
		// 匹配失败
		MatchErrorV1(data.GetOpenIdList(), fmt.Sprintf("MatchBattleV1 序列化失败, err: %v", err), data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_ProtoMarshalError))
		return
	}
	uid := queryUidByOpenid(openId)
	if uid == "" {
		ziLog.Error(fmt.Sprintf("MatchBattleV1 queryUidInterconvertOpenid 查询uid失败, openId: %s", openId), debug)
		// 匹配失败
		MatchErrorV1(data.GetOpenIdList(), fmt.Sprintf("MatchBattleV1 queryUidInterconvertOpenid 查询uid失败, openId: %s, err: %v", openId, err), data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_GetUuidError))
		return
	}
	sendUidList, _, _, _ := getUidListByGroupId(groupId)
	if len(sendUidList) == 0 {
		ziLog.Error(fmt.Sprintf("MatchBattleV1 查询uid失败, openId: %s", openId), debug)
		// 匹配失败
		MatchErrorV1(data.GetOpenIdList(), fmt.Sprintf("MatchBattleV1 查询uid失败, openId: %s, err: %v", openId, err), data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_GetUuidError))
		return
	}
	log.Println("MatchBattleV1 success", sendUidList)
	if err := sse.SseSend(pmsg.MessageId_MatchBattleV1ApplyAck, sendUidList, sData); err != nil {
		ziLog.Error(fmt.Sprintf("MatchBattleV1 发送消息失败, err: %v", err), debug)
		// 匹配失败
		MatchErrorV1(data.GetOpenIdList(), fmt.Sprintf("MatchBattleV1 发送消息失败, openId: %s, err: %v", openId, err), data.GetMatchBattleRoomId(),
			int32(pmsg.ErrorStatus_SendMessageError))
		return
	}
}

// 匹配错误消息发送
func MatchErrorV1(openIdList []string, msgErr, groupId string, errLabel int32) {
	defer battlematchv1.UnregisterBattleV1ByGroupId(first_ctx, groupId)
	data := &pmsg.MatchBattleErrorMessage{}
	data.TimeStamp = time.Now().UnixMilli()
	data.ErrorCode = errLabel
	data.ErrorMsg = msgErr
	data.MatchBattleRoomId = groupId
	for _, openId := range openIdList {
		uid := queryUidByOpenid(openId)
		if uid == "" {
			ziLog.Error(fmt.Sprintf("MatchErrorV1 查询uid失败, openId: %s", openId), debug)
			continue
		}
		data.OpenId = openId
		data.RoomId = QueryRoomIdInterconvertAnchorOpenId(openId) //查询房间id
		sData, err := proto.Marshal(data)
		if err != nil {
			ziLog.Error(fmt.Sprintf("MatchErrorV1 序列化失败, err: %v, data: %v", err, data.String()), debug)
			continue
		}
		if err := sse.SseSend(pmsg.MessageId_MatchBattleError, []string{uid}, sData); err != nil { // 推送消息给其他用户
			ziLog.Error(fmt.Sprintf("MatchErrorV1 pushDownLoadMessage err: %v, data: %v", err, data.String()), debug)
			continue
		}
	}
}

// 增加积分池积分
func matchAddIntrage(roomId, openId string, score float64) {
	// ctx := first_ctx
	WorldRankNumerAdd(openId, score)
	fastReturnAdd(roomId, openId, score)
	// go userInfoCompareStore(openId, nickName, avatarUrl, false)
}

// match瓜分积分池
func MatchBattleSetWinnerScore(battleGroupId string, result RoundUploadStruct) error {
	var (
		winDirection  string
		loseDirection string
		// winGroupList  []UserUploadScoreStruct //获胜组
		// loseGroupList []UserUploadScoreStruct //失败组
		winGroupList  []string
		loseGroupList []string
	)
	// 查询获胜组
	for _, v := range result.GroupResultList {
		if v.Result == 1 {
			winDirection = v.GroupId
			continue
		}
		loseDirection = v.GroupId
	}

	winGroupList, _ = rdb.ZRevRange(path.Join(battleGroupId, winDirection), 0, -1)
	loseGroupList, _ = rdb.ZRevRange(path.Join(battleGroupId, loseDirection), 0, -1)
	scoreFloat, err := rdb.ZScore(battleGroupId, group_integral_pool_key)
	if err != nil {
		return err
	}
	// 获取连胜币
	coin, _ := rdb.ZScore(battleGroupId, "coin"+winDirection)

	for i, u := range winGroupList {
		var score float64
		//设置为整数
		switch {
		case i == 0:
			score = math.Ceil(scoreFloat * 0.3)
			getCoin := coin / 2
			if _, err := AddUserWinStreamCoin(u, int64(getCoin)); err != nil {
				ziLog.Error(fmt.Sprintf("matchBattleSetWinnerScore addUserWinStreamCoin err: %v, openId: %v, coin: %v", err, u, getCoin), debug)
			}
		case i == 1:
			score = math.Ceil(scoreFloat * 0.2)
			getCoin := coin / 3
			if _, err := AddUserWinStreamCoin(u, int64(getCoin)); err != nil {
				ziLog.Error(fmt.Sprintf("matchBattleSetWinnerScore addUserWinStreamCoin err: %v, openId: %v, coin: %v", err, u, getCoin), debug)
			}
		case i == 2:
			score = math.Ceil(scoreFloat * 0.15)
			getCoin := coin / 6
			if _, err := AddUserWinStreamCoin(u, int64(getCoin)); err != nil {
				ziLog.Error(fmt.Sprintf("matchBattleSetWinnerScore addUserWinStreamCoin err: %v, openId: %v, coin: %v", err, u, getCoin), debug)
			}
		case i == 3:
			score = math.Ceil(scoreFloat * 0.1)
		case i == 4:
			score = math.Ceil(scoreFloat * 0.05)
		case i >= 5 && i <= 19:
			score = math.Ceil(scoreFloat * 0.012)
		}
		if _, err := AddUserWinStreamCoin(u, 1); err != nil {
			ziLog.Error(fmt.Sprintf("matchBattleSetWinnerScore addUserWinStreamCoin err: %v, openId: %v, coin: %v", err, u, 1), debug)
		}
		// 瓜分积分添加到世界排行榜
		if score > 0 {
			go WorldRankNumerAdd(u, score)
		}
	}
	for _, u := range loseGroupList {
		if _, err := AddUserWinStreamCoin(u, -1); err != nil {
			ziLog.Error(fmt.Sprintf("matchBattleSetWinnerScore addUserWinStreamCoin err: %v, openId: %v, coin: %v", err, u, -1), debug)
		}
	}
	// 删除积分池
	delIntegral(battleGroupId)
	return nil
}

// send开始游戏确认
func matchBattleStartGamedConfirm(groupId string, userIdList []string) {

	sendUidList, _, _, _ := getUidListByGroupId(groupId)
	if len(sendUidList) == 0 {
		MatchErrorV1(userIdList, "MatchBattleStartGamedConfirm queryOpenidFromUidStr nil", groupId,
			int32(pmsg.ErrorStatus_SendMessageError))
		return
	}
	data := &pmsg.MatchBattleStartGamedConfirmMessage{}
	data.MatchBattleRoomId = groupId
	data.OpenIdList = userIdList
	data.TimeStamp = time.Now().UnixMilli()
	dataByte, err := proto.Marshal(data)
	if err != nil {
		MatchErrorV1(userIdList, "MatchBattleStartGamedConfirm proto.Marshal err: "+err.Error(), groupId,
			int32(pmsg.ErrorStatus_SendMessageError))
		return
	}
	if err := sse.SseSend(pmsg.MessageId_MatchBattleStartGamedConfirm,
		sendUidList, dataByte); err != nil {
		MatchErrorV1(userIdList, "MatchBattleStartGamedConfirm pushDownLoadMessage err: "+err.Error(), groupId,
			int32(pmsg.ErrorStatus_SendMessageError))
		ziLog.Error(fmt.Sprintf("MatchBattleStartGamedConfirm send err: %v， group: %v", err, groupId), debug)
	}
}
