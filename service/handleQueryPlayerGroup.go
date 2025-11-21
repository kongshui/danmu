package service

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/kongshui/danmu/common"

	"github.com/gin-gonic/gin"
)

func QueryPlayerGroupHandler(c *gin.Context) {
	type queryPlayerGroup struct {
		APPId  string `json:"app_id"`
		OpenId string `json:"open_id"`
		RoomId string `json:"room_id"`
	}
	var (
		queryInfo queryPlayerGroup
		endGame   int
		userGroup string
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
	// 验签
	if c.GetHeader("x-signature") != common.DySignature(headers, string(*bodyByte), cfg.App.QueryGroupSecret) {
		ziLog.Error("QueryPlayerGroupHandler dy签名错误", debug)
		c.JSON(200, gin.H{
			"errcode": 11,
			"errmsg":  "签名错误",
		})
		return
	}
	// 解析
	if err := json.Unmarshal(*bodyByte, &queryInfo); err != nil {
		ziLog.Error(fmt.Sprintf("QueryPlayerGroupHandler 解析参数错误, err: %v", err), debug)
		c.JSON(200, gin.H{
			"errcode": 40001,
			"errmsg":  err.Error(),
		})
		return
	}
	if queryInfo.APPId != app_id {
		ziLog.Error(fmt.Sprintf("QueryPlayerGroupHandler roomid或者appid不匹配, queryId: %v, appId: %v", queryInfo.APPId, app_id), debug)
		c.JSON(200, gin.H{
			"errcode": 40001,
			"errmsg":  "roomid或者appid不匹配",
		})
		return
	}
	group, roundId, _, err := queryPlayerInGroup(queryInfo.RoomId, queryInfo.OpenId)
	if err != nil {
		ziLog.Error(fmt.Sprintf("QueryPlayerGroupHandler 查询玩家所在组失败, group: %v, roundId: %v, roomId: %v, openId： %v, err: %v",
			group, roundId, queryInfo.RoomId, queryInfo.OpenId, err), debug)
		log.Println("QueryPlayerGroupHandler 查询玩家所在组失败", group, roundId, err)
		c.JSON(200, gin.H{
			"errcode": 1,
			"errmsg":  "参数不合法",
		})
		return
	}
	if group == "" {
		userGroup = groupid_list[0]
	}
	//获取游戏是否完成
	ok, _ := rdb.HExists(queryInfo.RoomId+"_"+strconv.FormatInt(roundId, 10)+"_group", group)
	if ok {
		endGame = 2
	} else {
		endGame = 1
	}
	c.JSON(200, gin.H{
		"errcode": 0,
		"errmsg":  "success",
		"data": gin.H{
			"round_id":          roundId,
			"group_id":          group,
			"user_group_status": userGroup,
			"round_status":      endGame,
		},
	})
}
