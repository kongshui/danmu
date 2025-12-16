package service

import (
	"errors"
	"fmt"

	"github.com/kongshui/danmu/model/pmsg"
)

// WebsocketMessageFunc   websocket
func WebsocketMessageFunc(msg *pmsg.MessageBody) error {
	// log.Println("websocketMessageFunc: ", msg.MessageType)
	ziLog.Debug(fmt.Sprintf("websocketMessageFunc: %v,uid: %v", msg.MsgId, msg.Uuid), debug)
	switch msg.MsgId {
	case pmsg.MessageId_StartBind: //快手绑定
		return KsBind(msg, "start")
	case pmsg.MessageId_StopBind: // 快手解除绑定
		return KsBind(msg, "stop")
	case pmsg.MessageId_SingleRoomAddGroup: // 加入房间
		return playerAddGroudId(msg)
	case pmsg.MessageId_RoundDataUpLoad: // 上传数据
		return roundDataUpload(msg)
	case pmsg.MessageId_SyncGameStart: // 开始游戏对局
		return roundStart(msg)
	case pmsg.MessageId_SyncGameEnd: // 结束游戏对局
		return roundEnd(msg)
	// case pmsg.MessageId_GetVersionTopHundred: // 获取版本前100名
	// 	return getTopUserInfo(msg)
	case pmsg.MessageId_TestMsg: // 测试消息
		return testMsg(msg)
	case pmsg.MessageId_UseWinnerStreamCoin: // 使用用户的获胜币
		return useUserWinningStreamCoin(msg)
	case pmsg.MessageId_UserAddWinnerStreamCoin: // 添加用户的获胜币
		return addUsersWinningStreamCoin(msg)
	case pmsg.MessageId_QueryWinnerStreamCoin: // 查询用户的获胜币
		return queryUserWinningStreamCoin(msg)
	case pmsg.MessageId_IsFirstComsume: // 是否第一次消费
		return consumeUse(msg)
	case pmsg.MessageId_DisConnect: // 断开连接
		return disconnect(msg)
	case pmsg.MessageId_MsgAckSendAck: // 快手、抖音信息消费接口
		return ksMsgAck(msg)
	case pmsg.MessageId_AddIntegral: // 添加选边接口
		return addIntegral(msg)
	case pmsg.MessageId_ReLogin: // 重连
		return reconnect(msg)
	case pmsg.MessageId_MatchBattleV1Apply: // 1v1匹配
		if is_pk_match {
			return matchV1(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleV1Ready: // 1v1匹配准备
		if is_pk_match {
			return readyV1(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleV1ReadyAck: // 1v1匹配准备返回
		if is_pk_match {
			return readyV1Ack(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleV1TimeCheckAck: // 1v1匹配时间确定
		if is_pk_match {
			return matchV1TimeAck(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleV1StartConfirm: // 1v1匹配开始确认
		if is_pk_match {
			return matchV1Confirm(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleV1StartConfirmAck: // 1v1匹配开始确认返回
		if is_pk_match {
			return matchV1ConfirmAck(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleV1AskRoundIdAck: // 1v1匹配获取对局id返回
		if is_pk_match {
			return askMatchV1RoundId(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleV1Start: // 1v1匹配开始
		if is_pk_match {
			return matchV1Start(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleV1End: // 1v1匹配结束
		if is_pk_match {
			return matchV1End(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleV1Sync: // pk同步数据
		if is_pk_match {
			return matchV1SyncData(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleV1SyncAck: // pk同步数据返回
		if is_pk_match {
			return matchV1SyncDataAck(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleV1RoundUpload: // pk上传数据
		if is_pk_match {
			return matchV1DataUpload(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleAddIntegral: // pk添加节点积分
		// if is_connect {
		// 	return MatchV1AddIntegral(msg)
		// }
		return nil
	case pmsg.MessageId_MatchBattleUseWinnerStreamCoin: // PK使用连胜币
		// if is_connect {
		// 	return MatchV1UseWinnerStreamCoin(msg)
		// }
		return nil
	case pmsg.MessageId_MatchBattleV1Cancel: // pk取消匹配
		if is_pk_match {
			return matchV1Cancel(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleAddStreamCoin: // pk添加连胜币
		// if is_connect {
		// 	return MatchV1AddStreamCoin(msg)
		// }
		return nil
	case pmsg.MessageId_MatchBattleStartGamedConfirmAck: // pk开始游戏确认
		if is_pk_match {
			return matchStartGamedConfirmAck(msg)
		}
		return nil
	case pmsg.MessageId_MatchBattleQuitWithError:
		if is_pk_match {
			return matchBattleQuitWithError(msg)
		}
		return nil
	case pmsg.MessageId_Token:
		return dytoken(msg)
	case pmsg.MessageId_LevelQuery: // 等级查询
		return levelQuery(msg)
	case pmsg.MessageId_FrontSendMessageError: // 前端发送消息错误
		return getFrontEndErrorInfo(msg)
	case pmsg.MessageId_SendLogInfo: // 发送日志信息
		return recvLog(msg)
	case pmsg.MessageId_ConfigMapRequest: // 配置文件请求
		return configMapRequest(msg)
	default:
		if otherWebsocket != nil {
			return otherWebsocket(msg)
		}
		return errors.New("websocket消息类型不存在")

	}
}
