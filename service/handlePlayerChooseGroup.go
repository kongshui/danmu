package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kongshui/danmu/model/pmsg"
	"github.com/kongshui/danmu/sse"

	"github.com/kongshui/danmu/common"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func PlayerChooseGroupHandle(c *gin.Context) {
	if is_maintain {
		c.JSON(400, gin.H{
			"errcode": 90001,
			"errmsg":  "维护中",
		})
		return
	}
	type playerChooseGroup struct {
		AppId     string `json:"app_id"`
		OpenId    string `json:"open_id"`
		RoomId    string `json:"room_id"`
		GroupId   string `json:"group_id"`
		AvatarUrl string `json:"avatar_url"`
		NickName  string `json:"nickname"`
	}

	var (
		pCG playerChooseGroup
	)
	bodyByte := bytePool.Get().(*[]byte)
	defer bytePool.Put(bodyByte)
	*bodyByte, _ = c.GetRawData()
	headers := map[string]string{
		"x-nonce-str": c.GetHeader("x-nonce-str"),
		"x-timestamp": c.GetHeader("x-timestamp"),
		"x-roomid":    c.GetHeader("x-roomid"),
		"x-msg-type":  c.GetHeader("x-msg-type"),
	}
	if c.GetHeader("x-signature") != common.DySignature(headers, string(*bodyByte), config.App.ChooseGroupSecret) {
		ziLog.Error("PlayerChooseGroupHandle dy签名错误", debug)
		c.JSON(400, gin.H{
			"errcode": 11,
			"errmsg":  "签名错误",
		})
		return
	}
	if err := json.Unmarshal(*bodyByte, &pCG); err != nil {
		ziLog.Error(fmt.Sprintf("PlayerChooseGroupHandle 解析参数错误,err: %v", err), debug)
		c.JSON(400, gin.H{
			"errcode": 40001,
			"errmsg":  err.Error(),
		})
		return
	}
	if pCG.AppId != app_id {
		ziLog.Error(fmt.Sprintf("PlayerChooseGroupHandle 房间号不匹配,roomid: %v, getRoomId: %v", pCG.RoomId, c.GetHeader("X-Room-ID")), debug)
		c.JSON(400, gin.H{
			"errcode": 40001,
			"errmsg":  "roomid或者appid不匹配",
		})
		return
	}
	roundId, ok := queryRoomIdToRoundId(pCG.RoomId)
	if !ok {
		ziLog.Error("PlayerChooseGroupHandle 获取roundId失败", debug)
		c.JSON(400, gin.H{
			"errcode": 40001,
			"errmsg":  errors.New("PlayerChooseGroupHandle 获取roundId失败"),
		})
		return
	}
	uid := queryRoomIdToUid(pCG.RoomId)
	go playerGroupAdd(pCG.RoomId, uid, roundId, []*pmsg.SingleRoomAddGroupInfo{
		{
			GroupId:   pCG.GroupId,
			OpenId:    pCG.OpenId,
			AvatarUrl: pCG.AvatarUrl,
			NickName:  pCG.NickName,
		},
	}, true)
	go func(pcg playerChooseGroup, roundId int64) {

		anchorOpenid := QueryRoomIdInterconvertAnchorOpenId(pcg.RoomId)
		if anchorOpenid == "" {
			ziLog.Error(fmt.Sprintf("PlayerChooseGroupHandle 获取主播openid失败, roomId: %v, openid: %v", pcg.RoomId, pcg.OpenId), debug)
			c.JSON(400, gin.H{
				"errcode": 40001,
				"errmsg":  "获取主播openid失败",
			})
			return
		}
		sendUidList, _, _, _ := getUidListByOpenId(anchorOpenid)
		if len(sendUidList) == 0 {
			ziLog.Error(fmt.Sprintf("PlayerChooseGroupHandle sendUidList is nil, roomId: %v, anchorOpenid: %v, data: %v", pcg.RoomId, anchorOpenid, pcg), debug)
			c.JSON(400, gin.H{
				"errcode": 40001,
				"errmsg":  "获取主播uid失败",
			})
			return
		}
		gId, _, _ := getUserGroup(pcg.RoomId, pcg.OpenId)
		score, rank, _ := getPlayerWorldRankData(pcg.OpenId)
		coin, _ := QueryUserWinStreamCoin(pcg.OpenId)
		isConsume := queryIsConsume(pcg.OpenId)
		// 查询玩家等级
		level, _ := QueryLevelInfo(pcg.OpenId)
		sData := &pmsg.SingleUserAddGroupMessage{
			OpenId:            pcg.OpenId,
			AvatarUrl:         pcg.AvatarUrl,
			NickName:          pcg.NickName,
			GroupId:           gId,
			RoundId:           roundId,
			WorldScore:        score,
			WorldRank:         rank,
			WinningStreamCoin: coin,
			IsConsume:         !isConsume,
			Level:             level,
		}
		sDataByte, _ := proto.Marshal(sData)
		if err := sse.SseSend(pmsg.MessageId_SingleUserAddGroup, sendUidList, sDataByte); err != nil {
			ziLog.Error(fmt.Sprintf("PlayerChooseGroupHandle 推送玩家加入组信息失败: %v, data: %v", err, pcg), debug)
		}
	}(pCG, roundId)

	c.JSON(200, gin.H{
		"errcode": 0,
		"errmsg":  "success",
		"data": gin.H{
			"round_id":     roundId,
			"round_status": 1,
			"group_id":     pCG.GroupId,
		},
	})
}
