package service

import (
	"errors"
	"fmt"

	"github.com/kongshui/danmu/model/pmsg"
)

// websocketMessageFunc   websocket
func websocketMessageFunc(msg *pmsg.MessageBody) error {
	// log.Println("websocketMessageFunc: ", msg.MessageType)
	ziLog.Debug(fmt.Sprintf("websocketMessageFunc: %v,uid: %v", msg.MessageType, msg.Uuid), debug)
	switch msg.MessageType {
	case pmsg.MessageId_StartBind.String(): //快手绑定
		return KsBind(msg, "start")
	case pmsg.MessageId_StopBind.String(): // 快手解除绑定
		return KsBind(msg, "stop")
	case pmsg.MessageId_SingleRoomAddGroup.String(): // 加入房间
		return playerAddGroudId(msg)
	case pmsg.MessageId_RoundDataUpLoad.String(): // 上传数据
		return roundDataUpload(msg)
	case pmsg.MessageId_SyncGameStart.String(): // 开始游戏对局
		return roundStart(msg)
	case pmsg.MessageId_SyncGameEnd.String(): // 结束游戏对局
		return roundEnd(msg)
	case pmsg.MessageId_GetVersionTopHundred.String(): // 获取版本前100名
		return getUserWorldInfo(msg)
	case pmsg.MessageId_TestMsg.String(): // 测试消息
		return testMsg(msg)
	case pmsg.MessageId_UseWinnerStreamCoin.String(): // 使用用户的获胜币
		return useUserWinningStreamCoin(msg)
	case pmsg.MessageId_UserAddWinnerStreamCoin.String(): // 添加用户的获胜币
		return addUsersWinningStreamCoin(msg)
	case pmsg.MessageId_QueryWinnerStreamCoin.String(): // 查询用户的获胜币
		return queryUserWinningStreamCoin(msg)
	case pmsg.MessageId_IsFirstComsume.String(): // 是否第一次消费
		return consumeUse(msg)
	case pmsg.MessageId_DisConnect.String(): // 断开连接
		return disconnect(msg)
	case pmsg.MessageId_MsgAckSendAck.String(): // 快手、抖音信息消费接口
		return ksMsgAck(msg)
	case pmsg.MessageId_AddIntegral.String(): // 添加选边接口
		return addIntegral(msg)
	case pmsg.MessageId_ReLogin.String(): // 重连
		return reconnect(msg)
	case pmsg.MessageId_MatchBattleV1Apply.String(): // 1v1匹配
		if is_pk_match {
			return matchV1(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleV1Ready.String(): // 1v1匹配准备
		if is_pk_match {
			return readyV1(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleV1ReadyAck.String(): // 1v1匹配准备返回
		if is_pk_match {
			return readyV1Ack(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleV1TimeCheckAck.String(): // 1v1匹配时间确定
		if is_pk_match {
			return matchV1TimeAck(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleV1StartConfirm.String(): // 1v1匹配开始确认
		if is_pk_match {
			return matchV1Confirm(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleV1StartConfirmAck.String(): // 1v1匹配开始确认返回
		if is_pk_match {
			return matchV1ConfirmAck(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleV1AskRoundIdAck.String(): // 1v1匹配获取对局id返回
		if is_pk_match {
			return askMatchV1RoundId(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleV1Start.String(): // 1v1匹配开始
		if is_pk_match {
			return matchV1Start(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleV1End.String(): // 1v1匹配结束
		if is_pk_match {
			return matchV1End(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleV1Sync.String(): // pk同步数据
		if is_pk_match {
			return matchV1SyncData(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleV1SyncAck.String(): // pk同步数据返回
		if is_pk_match {
			return matchV1SyncDataAck(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleV1RoundUpload.String(): // pk上传数据
		if is_pk_match {
			return matchV1DataUpload(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleAddIntegral.String(): // pk添加节点积分
		// if is_connect {
		// 	return MatchV1AddIntegral(msg)
		// }
		return nil
	case pmsg.MessageId_MatchBattleUseWinnerStreamCoin.String(): // PK使用连胜币
		// if is_connect {
		// 	return MatchV1UseWinnerStreamCoin(msg)
		// }
		return nil
	case pmsg.MessageId_MatchBattleV1Cancel.String(): // pk取消匹配
		if is_pk_match {
			return matchV1Cancel(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleAddStreamCoin.String(): // pk添加连胜币
		// if is_connect {
		// 	return MatchV1AddStreamCoin(msg)
		// }
		return nil
	case pmsg.MessageId_MatchBattleStartGamedConfirmAck.String(): // pk开始游戏确认
		if is_pk_match {
			return matchStartGamedConfirmAck(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleQuitWithError.String():
		if is_pk_match {
			return matchBattleQuitWithError(msg)
		}
		return nil
	case pmsg.MessageId_Token.String():
		return dytoken(msg)
	case "err":
		getFrontEndErrorInfo(msg)
		return nil
	case "get_top_100_rank":
		return getLastTop100Rank(msg)
	case pmsg.MessageId_LevelQuery.String(): // 等级查询
		return levelQuery(msg)
	default:
		if otherWebsocket != nil {
			return otherWebsocket(msg)
		}
		return errors.New("websocket消息类型不存在")

	}
}
