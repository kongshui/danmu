package service

import (
	"encoding/json"
	"fmt"

	"github.com/kongshui/danmu/model/pmsg"

	"github.com/kongshui/danmu/common"

	"github.com/gin-gonic/gin"
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
	if c.GetHeader("x-signature") != common.DySignature(headers, string(*bodyByte), cfg.App.ChooseGroupSecret) {
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
		c.JSON(200, gin.H{
			"errcode": 1,
			"errmsg":  "参数不合法",
		})
		return
	}
	uid := queryRoomIdToUid(pCG.RoomId)
	if playerGroupAddin != nil {
		go playerGroupAddin(pCG.RoomId, uid, roundId, []*pmsg.SingleRoomAddGroupInfo{
			{
				GroupId:   pCG.GroupId,
				OpenId:    pCG.OpenId,
				AvatarUrl: pCG.AvatarUrl,
				NickName:  pCG.NickName,
			},
		})
	}

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
