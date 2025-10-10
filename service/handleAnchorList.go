package service

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kongshui/danmu/model/pmsg"
	"github.com/kongshui/danmu/sse"
)

// BlackAnchorDisconnectHandle 黑名单主播断线处理
func BlackAnchorDisconnectHandle(c *gin.Context) {
	type blackAnchorDisconnect struct {
		OpenId string `json:"open_id"`
		Msg    string `json:"msg"`
	}
	var (
		bAD blackAnchorDisconnect
	)
	bodyByte := bytePool.Get().(*[]byte)
	defer bytePool.Put(bodyByte)
	*bodyByte, _ = c.GetRawData()
	if err := json.Unmarshal(*bodyByte, &bAD); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errcode": 1,
			"errmsg":  "json unmarshal failed",
		})
		return
	}
	if bAD.OpenId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"errcode": 2,
			"errmsg":  "open_id or msg is empty",
		})
		return
	}
	if bAD.Msg == "" {
		bAD.Msg = "您已被管理员移出房间，如有疑问请联系客服"
	}
	data := &pmsg.BlackAnchorMessage{
		AnchorOpenId: bAD.OpenId,
		Msg:          bAD.Msg,
	}
	dataBytes, _ := json.Marshal(data)
	// 发送断线消息
	sse.SseSend(pmsg.MessageId_BlackAnchorLogOff, []string{bAD.OpenId}, dataBytes)
	c.JSON(200, gin.H{
		"errcode": 0,
		"errmsg":  "success",
	})
}

// BlackAnchorReconnectHandle 黑名单主播添加至列表处理
func BlackAnchorAddHandle(c *gin.Context) {
	type blackAnchorAdd struct {
		OpenId string `json:"open_id"`
	}
	var (
		bAA blackAnchorAdd
	)
	bodyByte := bytePool.Get().(*[]byte)
	defer bytePool.Put(bodyByte)
	*bodyByte, _ = c.GetRawData()
	if err := json.Unmarshal(*bodyByte, &bAA); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errcode": 1,
			"errmsg":  "json unmarshal failed",
		})
		return
	}
	if err := blackAnchorListAdd(bAA.OpenId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errcode": 2,
			"errmsg":  "black anchor list add failed",
		})
		return
	}
	c.JSON(200, gin.H{
		"errcode": 0,
		"errmsg":  "success",
	})
}

// BlackAnchorDelHandle 黑名单主播移除列表处理
func BlackAnchorDelHandle(c *gin.Context) {
	type blackAnchorDel struct {
		OpenId string `json:"open_id"`
	}
	var (
		bAD blackAnchorDel
	)
	bodyByte := bytePool.Get().(*[]byte)
	defer bytePool.Put(bodyByte)
	*bodyByte, _ = c.GetRawData()
	if err := json.Unmarshal(*bodyByte, &bAD); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errcode": 1,
			"errmsg":  "json unmarshal failed",
		})
		return
	}
	if err := blackAnchorListDel(bAD.OpenId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errcode": 2,
			"errmsg":  "black anchor list del failed",
		})
		return
	}
	c.JSON(200, gin.H{
		"errcode": 0,
		"errmsg":  "success",
	})
}

// BlackAnchorListMembersHandle 返回黑名单主播列表处理
func BlackAnchorListMembersHandle(c *gin.Context) {
	blackAnchorList, err := blackAnchorListMembers()
	if err != nil || len(blackAnchorList) == 0 {
		c.JSON(200, []string{})
		return
	}
	var (
		userInfos []UserInfoStruct
	)
	for _, anchorOpenid := range blackAnchorList {
		userInfo, _ := userInfoGet(anchorOpenid, true)
		userInfos = append(userInfos, userInfo)
	}
	c.JSON(200, userInfos)
}

// 清空黑名单主播列表处理
func BlackAnchorListClearHandle(c *gin.Context) {
	if err := blackAnchorListClear(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errcode": 2,
			"errmsg":  "black anchor list clear failed",
		})
		return
	}
	c.JSON(200, gin.H{
		"errcode": 0,
		"errmsg":  "success",
	})
}
