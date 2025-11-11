package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/kongshui/danmu/model/pmsg"
	"github.com/kongshui/danmu/sse"

	"github.com/kongshui/danmu/common"

	"google.golang.org/protobuf/proto"
)

// 测试test
func testMsg(msg *pmsg.MessageBody) error {
	return sse.SseSend(pmsg.MessageId_TestMsgAck, []string{msg.GetUuid()}, msg.GetMessageData())
}

// 对局start
func roundStart(msg *pmsg.MessageBody) error {
	syncGameStatusData := &pmsg.SyncGameStatusMessage{}
	err := proto.Unmarshal(msg.MessageData, syncGameStatusData)
	if err != nil {
		return errors.New("roundStart Unmarshal err: " + err.Error())
	}
	isMember := blackAnchorListIsMember(syncGameStatusData.GetAnchorOpenId())
	if isMember {
		return fmt.Errorf("anchor is black")
	}
	getRoundId, _ := queryRoomIdToRoundId(syncGameStatusData.GetRoomId())
	if syncGameStatusData.GetRoundId() == getRoundId {
		return nil
	}
	// ok, groupId := battlematchv1.IsInVs1Group(first_ctx, syncGameStatusData.GetAnchorOpenId())
	// if ok {
	// 	if !battlematchv1.QueryOpenIdInMatchDisconnect(first_ctx, groupId, syncGameStatusData.GetAnchorOpenId()) {
	// 		if err := battlematchv1.DisconnectMatchRegister(first_ctx, syncGameStatusData.GetAnchorOpenId()); err != nil {
	// 			ziLog.Error( fmt.Sprintf("roundStart 添加 QueryOpenIdInMatchDisconnect 失败，openId： %v, err: %v", syncGameStatusData.GetAnchorOpenId(), err), debug)
	// 		}
	// 	}
	// }
	time.Sleep(1 * time.Second)

	data := SyncGameStatusStruct{
		AnchorOpenId: syncGameStatusData.AnchorOpenId,
		AppId:        app_id,
		RoomId:       syncGameStatusData.RoomId,
		RoundId:      syncGameStatusData.RoundId,
		StartTime:    syncGameStatusData.StartTime,
		Status:       1,
	}
	//同步开始对局状态
	switch platform {
	case "ks":
		if err := ksSyncGameStatus(data, "start", true); err != nil {
			return fmt.Errorf("roundStart 同步开始对局状态失败,roomId: %v, roundId: %v, err: %v", data.RoomId, data.RoundId, err)
		}
	case "dy":
		if is_mock {
			if err := liveCurrentRoundAdd(syncGameStatusData.RoomId, syncGameStatusData.RoundId); err != nil {
				return errors.New("uplink liveCurrentRoundAdd err: " + err.Error())
			}
			break
		}
		if err := dySyncGameStatus(data); err != nil {
			return fmt.Errorf("roundStart 同步开始对局状态失败,roomId: %v, roundId: %v, err: %v", data.RoomId, data.RoundId, err)
		} else {
			if err := liveCurrentRoundAdd(syncGameStatusData.RoomId, syncGameStatusData.RoundId); err != nil {
				return errors.New("uplink liveCurrentRoundAdd err: " + err.Error())
			}
		}
	}

	score, err := GetIntegral(syncGameStatusData.AnchorOpenId)
	if err != nil {
		score = 0
	}
	sData := &pmsg.RoundReadyMessage{
		RoomId:        syncGameStatusData.RoomId,
		RoundId:       syncGameStatusData.RoundId,
		Timestamp:     syncGameStatusData.StartTime,
		LiveLikeScore: live_like_score,
		Integral:      int64(score),
	}

	msgData, err := proto.Marshal(sData)
	if err != nil {
		return errors.New("roundStart proto Marshal err: " + err.Error())
	}
	// 延迟调用选边
	switch platform {
	case "ks":
		t := time.NewTimer(time.Second * 2)
		<-t.C
		if interactive != nil {
			if !interactive(syncGameStatusData.RoomId, strconv.FormatInt(syncGameStatusData.RoundId, 10), 0) {
				ziLog.Error(fmt.Sprintf("round start, roomid:: %v,主播openid %v,设置自动选边互动失败", syncGameStatusData.RoomId, syncGameStatusData.AnchorOpenId), debug)
			} else {
				//添加roundid至CurrentRoundId
				if err := liveCurrentRoundAdd(syncGameStatusData.RoomId, syncGameStatusData.RoundId); err != nil {
					return errors.New("uplink liveCurrentRoundAdd err: " + err.Error())
				}
			}
		}
	case "dy":

	}
	ziLog.Info(fmt.Sprintf("roundStart, roomid:: %v,主播openid %v", syncGameStatusData.RoomId, syncGameStatusData.AnchorOpenId), debug)
	return sse.SseSend(pmsg.MessageId_SyncGameStartAck, []string{msg.Uuid}, msgData)
}

// testToekn
func TestTokenFunc(msg *pmsg.MessageBody) error {
	var result string = ""
	for range 19 {
		digit := rand.Intn(10)
		if digit == 0 {
			digit = 1
		}
		result += strconv.Itoa(digit)
	}
	roomInfo := &pmsg.AnchorInfoMessage{
		RoomId:       "1123456789876543212",
		AnchorOpenId: result,
		NickName:     result,
		AvatarUrl:    "",
	}
	var (
		data common.RoomRegister
	)
	data.Uuid = msg.Uuid
	data.RoomId = "1123456789876543212"
	data.UserId = result
	dataByte, err := json.Marshal(data)
	if err != nil {
		log.Println("json转换失败， info:", data, err, "err: ", err)
	}

	etcdClient.Client.Put(first_ctx, path.Join("/", cfg.Project, common.RoomInfo_Register_key, roomInfo.RoomId), string(dataByte))
	sData, err := proto.Marshal(roomInfo)
	if err != nil {
		return errors.New("token proto Marshal err: " + err.Error())
	}
	return sse.SseSend(pmsg.MessageId_TokenAck, []string{msg.Uuid}, sData)
}

// 对局结束
func roundEnd(msg *pmsg.MessageBody) error {
	data := &pmsg.SyncGameStatusMessage{}
	err := proto.Unmarshal(msg.MessageData, data)
	if err != nil {
		return errors.New("roundEnd Unmarshal err: " + err.Error())
	}

	var (
		groupId []GroupResultList = make([]GroupResultList, 0)
	)

	for _, v := range data.GroupResultList {
		groupId = append(groupId, GroupResultList{GroupId: v.GroupId, Result: int(v.Result)})
	}
	ziLog.Info(fmt.Sprintf("roundEnd roomid: %v, 主播openid: %v, roundId: %v, 组信息： %v", data.RoomId, data.AnchorOpenId, data.RoundId, groupId), debug)

	//json为map[string]any{}，mapp[data]为map[string]any, 胜利方为map[data]["winner"].(string),其他为map[string]int .openid:分数
	//删除CurrentRoundId中的roomid，暂时改为不删除
	// liveCurrentRoundDel(msg.RoomId)
	//同步结束对局状态
	useData := SyncGameStatusStruct{
		AnchorOpenId:    data.AnchorOpenId,
		AppId:           app_id,
		RoomId:          data.RoomId,
		RoundId:         data.RoundId,
		StartTime:       data.StartTime,
		EndTime:         data.EndTime,
		Status:          2,
		GroupResultList: groupId,
	}
	setRoundEndGroup(data.RoomId, data.RoundId, groupId)
	// 分平台
	switch platform {
	case "ks":
		if err := ksSyncGameStatus(useData, "stop", true); err != nil {
			return fmt.Errorf("roundEnd 同步结束对局状态失败,roomId: %v, roundId: %v, err: %v", data.RoomId, data.RoundId, err)
		}
	case "dy":
		if is_mock {
			return nil
		}
		if err := dySyncGameStatus(useData); err != nil {
			return fmt.Errorf("roundEnd 同步开始对局状态失败,roomId: %v, roundId: %v, err: %v", data.RoomId, data.RoundId, err)
		}
	}
	return nil
}

// 玩家加入组信息
func playerAddGroudId(msg *pmsg.MessageBody) error {
	data := &pmsg.SingleRoomAddGroupMessage{}
	err := proto.Unmarshal(msg.MessageData, data)
	if err != nil {
		return errors.New("player add group Unmarshal err: " + err.Error())
	}
	//获取roundId
	roundId, ok := queryRoomIdToRoundId(data.GetRoomId())
	if !ok {
		return errors.New("playerGroupAdd roundId 未查到")
	}
	if playerGroupAddin != nil {
		if err := playerGroupAddin(data.GetRoomId(), msg.GetUuid(), roundId, data.GetUserList(), false); err != nil {
			return errors.New("玩家加入组信息 err: " + err.Error())
		}
		return nil
	}
	return errors.New("playerAddGroudId playerGroupAddin is nil")
}

// 数据上报信息
func roundDataUpload(msg *pmsg.MessageBody) error {
	data := &pmsg.RoundUploadMessage{}
	err := proto.Unmarshal(msg.MessageData, data)
	if err != nil {
		return errors.New("roundDataUpload Unmarshal err: " + err.Error())
	}
	ziLog.Info(fmt.Sprintf("roundDataUpload roomid: %v, 主播openid: %v, roundId: %v, 组信息： %v, 玩家信息： %v", data.RoomId, data.GetAnchorOpenId(), data.RoundId, data.GroupResultList, data.GroupUserList), debug)
	//数据上报 =====================================================
	roundData := RoundUploadStruct{}
	roundData.RoundId = data.RoundId
	for _, v := range data.GroupResultList {
		roundData.GroupResultList = append(roundData.GroupResultList, GroupResultList{GroupId: v.GroupId, Result: int(v.Result)})
	}
	for _, v := range data.GroupUserList {
		roundData.GroupUserList = append(roundData.GroupUserList, UserUploadScoreStruct{OpenId: v.OpenId, GroupId: v.GroupId, Score: v.Score, RoundResult: int(v.RoundResult), Rank: v.Rank})
	}
	if is_mock {
		return nil
	}
	return usersRoundUpload(data.RoomId, data.GetAnchorOpenId(), roundData, data.GameType)
}

// get_user_worldinfo,获取玩家世界列表
func GetTopUserInfoTemplate(msg *pmsg.MessageBody) error {
	var (
		data  *pmsg.UserInfoListMessage
		msgId pmsg.MessageId
	)
	// 获取前一百名用户世界信息
	switch msg.MsgId {
	case pmsg.MessageId_GetVersionTopHundred:
		data = getTopWorldRankData(0, 99, false)
		msgId = pmsg.MessageId_GetVersionTopHundredAck
	case pmsg.MessageId_GetMonthTopHundred:
		data = getTopMonthRankData()
		msgId = pmsg.MessageId_GetMonthTopHundredAck
	default:
		return errors.New("get user world info Unmarshal err: " + msg.MsgId.String())
	}
	// sdatabyte
	data.Timestamp = time.Now().Unix()
	sDataByte, err := proto.Marshal(data)
	if err != nil {
		return errors.New("get user world info proto Marshal err: " + err.Error())
	}

	//获取用户世界信息
	if err := sse.SseSend(msgId, []string{msg.Uuid}, sDataByte); err != nil {
		return errors.New("玩家获取世界信息 err: " + err.Error())
	}
	return nil
}

// useWinningStreamCoin,使用奖池 连胜币
func useUserWinningStreamCoin(msg *pmsg.MessageBody) error {

	data := &pmsg.RequestwinnerstreamcoinMessage{}
	err := proto.Unmarshal(msg.MessageData, data)
	if err != nil {
		return errors.New("useWinningStreamCoin Unmarshal err: " + err.Error())
	}

	sData := UseWinningStreamCoin(data)
	if sData.CanUse {
		_, ok := winCoinToComment[data.UseNum]
		if ok {
			ziLog.Gift(fmt.Sprintf("useUserWinningStreamCoin giftId: %v, useNum: %v, openId: %v, roomId: %v, roundId: %v, timestamp: %v", data.GetGiftId(),
				data.GetUseNum(), data.GetOpenId(), data.GetRoomId(), data.GetRoundId(), data.GetTimeStamp()), debug)
		}
	}
	sDataByte, err := proto.Marshal(sData)
	if err != nil {
		return errors.New("useWinningStreamCoin proto Marshal err: " + err.Error())
	}
	if err := sse.SseSend(pmsg.MessageId_UseWinnerStreamCoinAck, []string{msg.Uuid}, sDataByte); err != nil {
		return errors.New("玩家使用连胜币 err: " + err.Error())
	}
	return nil
}

// 玩家获取连胜币
func addUsersWinningStreamCoin(msg *pmsg.MessageBody) error {
	data := &pmsg.AddWinnerStreamCoinMessage{}
	err := proto.Unmarshal(msg.MessageData, data)
	if err != nil {
		return errors.New("addUsersWinningStreamCoin Unmarshal err: " + err.Error())
	}
	sData := addWinningStreamCoin(data.GetUserList())
	if data.IsEnd {
		return nil
	}
	// 转变为二进制
	sDataByte, err := proto.Marshal(sData)
	if err != nil {
		return errors.New("addUsersWinningStreamCoin proto Marshal err: " + err.Error())
	}

	if err := sse.SseSend(pmsg.MessageId_UserAddWinnerStreamCoinAck, []string{msg.Uuid}, sDataByte); err != nil {
		return errors.New("玩家获取连胜币 err: " + err.Error())
	}
	return nil
}

// 查询连胜币
func queryUserWinningStreamCoin(msg *pmsg.MessageBody) error {
	useData := &pmsg.QueryWinnerStreamCoinMessage{}
	err := proto.Unmarshal(msg.MessageData, useData)
	if err != nil {
		return errors.New("queryUserWinningStreamCoin Unmarshal err: " + err.Error())
	}
	data := queryWinningStreamCoin(useData.GetUserList())
	data.Side = useData.GetSide()
	data.TimeStamp = time.Now().UnixMilli()
	sDataByte, err := proto.Marshal(data)
	if err != nil {
		return errors.New("queryUserWinningStreamCoin proto Marshal err: " + err.Error())
	}
	if err := sse.SseSend(pmsg.MessageId_QueryWinnerStreamCoinAck, []string{msg.Uuid}, sDataByte); err != nil {
		return errors.New("玩家查询连胜币 err: " + err.Error())
	}
	return nil
}

// 获取前端错误信息
func getFrontEndErrorInfo(msg *pmsg.MessageBody) error {
	ziLog.Error(fmt.Sprintf("getFrontEndErrorInfo err: %v", string(msg.MessageData)), debug)
	return nil
}

// 获取上期前100名
func GetLastTop100Rank(msg *pmsg.MessageBody) error {
	_, err := getTop100Rank()
	if err != nil {
		return errors.New("获取上期前100名 err: " + err.Error())
	}
	if err := sse.SseSend(pmsg.MessageId_GetMonthTopHundredAck, []string{msg.Uuid}, []byte{}); err != nil {
		return errors.New("玩家查询连胜币 err: " + err.Error())
	}
	return nil
}

// 是否是第一次送礼物
func consumeUse(msg *pmsg.MessageBody) error {
	data := &pmsg.IsFirstConsumeMessage{}
	err := proto.Unmarshal(msg.MessageData, data)
	if err != nil {
		return errors.New("consumeUse Unmarshal err: " + err.Error())
	}
	ok := compareIsConsume(data.GetOpenId(), data.GetTimeStamp())
	data.IsConsume = ok
	dataByte, err := proto.Marshal(data)
	if err != nil {
		return errors.New("consumeUse proto Marshal err: " + err.Error())
	}
	if err := sse.SseSend(pmsg.MessageId_IsFirstComsumeAck, []string{msg.Uuid}, dataByte); err != nil {
		return errors.New("玩家查询连胜币 err: " + err.Error())
	}
	return nil
}

// 断开连接
func disconnect(msg *pmsg.MessageBody) error {
	data := &pmsg.Disconnect{}
	err := proto.Unmarshal(msg.MessageData, data)
	if err != nil {
		return errors.New("Disconnect Unmarshal err: " + err.Error())
	}
	ziLog.Info(fmt.Sprintf("disconnect enter roomId:%v, userId: %v", data.GetRoomId(), data.GetUserId()), debug)
	endConnect(data.GetRoomId(), data.GetUserId())
	return nil
}

// 重新链接
func reconnect(msg *pmsg.MessageBody) error {
	data := &pmsg.Reconnect{}
	err := proto.Unmarshal(msg.MessageData, data)
	if err != nil {
		return errors.New("Disconnect Unmarshal err: " + err.Error())
	}
	// 删除掉线房间
	if err := disconnectRoomIdDelete(data.GetRoomId()); err != nil {
		return errors.New("Disconnect disconnectRoomIdDelete err: " + err.Error())
	}
	// 删除匹配掉线注册, 这个只能再匹配重连的时候使用
	// if err := battlematch.DeleteMatchDisconnectRegisterUser(first_ctx, data.GetUserId()); err != nil {
	// 	return errors.New("Disconnect DeleteMatchDisconnectRegisterUser err: " + err.Error())
	// }
	ziLog.Info(fmt.Sprintf("reconnect enter roomId:%v, userId: %v", data.GetRoomId(), data.GetUserId()), debug)
	return nil
}

// 快手bind
func KsBind(msg *pmsg.MessageBody, label string) error {
	data := &pmsg.KSBindReq{}
	err := proto.Unmarshal(msg.MessageData, data)
	if err != nil {
		return errors.New("KsBind Unmarshal err: " + err.Error())
	}
	ziLog.Info(fmt.Sprintf("KsBind label: %v, roomId: %v, openId: %v", label, data.GetRoomId(), data.GetOpenId()), debug)
	if data.GetRoomId() == "" {
		return errors.New("roomId is nil")
	}
	if platform == "dy" {
		isStart := false
		if label == "start" {
			isStart = true
		}
		dyStartPushTask(data.RoomId, data.OpenId, msg.Uuid, isStart)
		return nil
	}

	if err := ksStartFinishGameInfo(data.GetRoomId(), url_BindUrl, label, msg.GetUuid(), true); err != nil {
		return errors.New(label + " game 游戏失败, " + data.GetRoomId())
	}
	if label == "start" {
		if err := topGift(data.RoomId); err != nil {
			return errors.New("TopGift 置顶失败, " + data.RoomId)
		}
	}
	return nil
}

// 快手消息验证接收
func ksMsgAck(msg *pmsg.MessageBody) error {
	switch platform {
	case "ks":
		data := &pmsg.KsMsgAck{}
		err := proto.Unmarshal(msg.MessageData, data)
		if err != nil {
			return errors.New("ksMsgAck Unmarshal err: " + err.Error())
		}
		if cfg.App.IsOnline {
			ksMsgAckSend(data.GetRoomId(), "cpClientShow", KsMsgAckReceiveStruct{
				UniqueMessageId:  data.GetData().GetUniqueMessageId(),
				PushType:         "giftSend",
				CpClientShowTime: time.Now().UnixMilli() - 10,
			})
		}
	case "dy":
		data := &pmsg.DymsgAckMessage{}
		err := proto.Unmarshal(msg.MessageData, data)
		if err != nil {
			ziLog.Error("dytoken Unmarshal err: "+err.Error(), debug)
			return errors.New("dytoken Unmarshal err: " + err.Error())
		}
		sdata := MsgAckStruct{}
		sdata.AckType = 2
		sdata.AppId = app_id
		sdata.RoomId = data.GetRoomId()
		var dataList []MsgAckInfoStruct = make([]MsgAckInfoStruct, 0)
		for _, v := range data.Data {
			dataList = append(dataList, MsgAckInfoStruct{
				MsgId:      v.GetMsgId(),
				MsgType:    v.GetMsgType(),
				ClientTime: v.GetClientTime(),
			})
		}
		temp, err := json.Marshal(dataList)
		if err != nil {
			return err
		}
		sdata.Data = string(temp)
		return msgAckSend(sdata)
	}
	return nil
}

// 增加节点积分
func addIntegral(msg *pmsg.MessageBody) error {
	data := &pmsg.AddIntegralReq{}
	err := proto.Unmarshal(msg.MessageData, data)
	if err != nil {
		return errors.New("addIntegral Unmarshal err: " + err.Error())
	}
	switch data.GetAddType() {
	case 0:
		if _, err := addIntegralByNode(data.GetOpenId(), float64(nodeIdToIntegral[data.GetNodeId()])); err != nil {
			return errors.New("addIntegral err: " + err.Error())
		}
		return nil
	case 1:
		WorldRankNumerAdd(data.GetOpenId(), float64(data.GetIntegral()))
		return nil
	case 2:
		if _, err := addIntegralByNode(data.GetOpenId(), data.GetIntegral()); err != nil {
			return errors.New("addIntegralByNode err: " + err.Error())
		}
		return nil
	case 3:
		if err := addIntegralAndUserIntegral(data.GetAnchorOpenId(), data.GetOpenId(), data.GetIntegral()); err != nil {
			return errors.New("addIntegralAndUserIntegral err: " + err.Error())
		}
		return nil
	default:
		return errors.New("addIntegral err: AddType not support")
	}
}

// Dytoken
func dytoken(msg *pmsg.MessageBody) error {
	data := &pmsg.TokenMessage{}
	err := proto.Unmarshal(msg.MessageData, data)
	ziLog.Info(fmt.Sprintf("dytoken token: %v", data.Token), debug)
	if err != nil {
		return errors.New("dytoken Unmarshal err: " + err.Error())
	}
	return dyGetAnchorInfo(msg.Uuid, data.Token)
}

func levelQuery(msg *pmsg.MessageBody) error {
	data := &pmsg.QueryLevelMessage{}
	err := proto.Unmarshal(msg.MessageData, data)
	if err != nil {
		return errors.New("levelQuery Unmarshal err: " + err.Error())
	}
	sData := &pmsg.QueryLevelResponse{}
	// 查询
	for _, v := range data.GetOpenidList() {
		level, err := QueryLevelInfo(v)
		if err != nil {
			return errors.New("queryLevel err: " + err.Error())
		}
		sData.UserLevelList = append(sData.UserLevelList, &pmsg.UserLevelStruct{
			OpenId: v,
			Level:  level,
		})
	}
	sDataByte, err := proto.Marshal(sData)
	if err != nil {
		return errors.New("levelQuery proto Marshal err: " + err.Error())
	}
	if err := sse.SseSend(pmsg.MessageId_LevelQueryAck, []string{msg.Uuid}, sDataByte); err != nil {
		return errors.New("levelQuery 玩家查询等级 err: " + err.Error())
	}
	return nil
}

// 接收日志
func recvLog(msg *pmsg.MessageBody) error {
	code := 0
	ackMsg := "success"
	data := &pmsg.SendLogInfoMessage{}
	if err := proto.Unmarshal(msg.MessageData, data); err != nil {
		code = 1
		ackMsg = "recvLog Unmarshal err: " + err.Error()
	} else {
		// 写入文件
		dataDir := cfg.App.LogStoreDir
		if cfg.App.LogStoreDir == "" {
			dataDir, err = os.Executable()
			if err != nil {
				return errors.New("recvLog get log dir err: " + err.Error())
			}
			dataDir = filepath.Join(filepath.Dir(dataDir), "logstore")
			cfg.App.LogStoreDir = dataDir
		}
		// 检查目录是否存在
		if _, err := os.Stat(dataDir); os.IsNotExist(err) {
			// 创建目录
			if err := os.MkdirAll(dataDir, 0755); err != nil {
				return errors.New("recvLog create log dir err: " + err.Error())
			}
		}
		err := writeToFile(filepath.Join(dataDir, data.GetAnchorOpenId()+"_"+data.GetLogLabel()+".log"), data.GetLogContent())
		if err != nil {
			code = 2
			ackMsg = "writeToFile err: " + err.Error()
		}
	}
	// 发送日志信息返回
	sData := &pmsg.SendLogInfoAckMessage{
		AnchorOpenId: data.GetAnchorOpenId(),
		LogLabel:     data.GetLogLabel(),
		LogId:        data.GetLogId(),
		Code:         int32(code),
		Msg:          ackMsg,
	}
	// 发送日志信息返回
	sDataByte, err := proto.Marshal(sData)
	if err != nil {
		return errors.New("recvLog proto Marshal err: " + err.Error())
	}
	if err := sse.SseSend(pmsg.MessageId_SendLogInfoAck, []string{msg.Uuid}, sDataByte); err != nil {
		return errors.New("recvLog 发送日志信息返回 err: " + err.Error())
	}
	return nil
}

// 配置文件请求
func configMapRequest(msg *pmsg.MessageBody) error {
	data, nil := json.Marshal(&cfgConfig)
	if err := sse.SseSend(pmsg.MessageId_ConfigMapRequestAck, []string{msg.Uuid}, data); err != nil {
		return errors.New("configMapRequest 发送配置文件请求返回 err: " + err.Error())
	}
	return nil
}
