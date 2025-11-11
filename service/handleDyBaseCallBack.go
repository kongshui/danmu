package service

import (
	"github.com/kongshui/danmu/common"

	"github.com/gin-gonic/gin"
)

// 点赞评论等回调接口
func BaseDyCallBackHandle(c *gin.Context) {
	var (
		response KsCallbackRespondStruct
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
	var secret string
	switch headers["x-msg-type"] {
	case "live_comment":
		secret = cfg.App.CommentSecret
	case "live_gift":
		secret = cfg.App.GiftSecret
	case "live_like":
		secret = cfg.App.LikeSecret
	}
	if c.GetHeader("x-signature") != common.DySignature(headers, string(*bodyByte), secret) {
		ziLog.Error("BaseCallBackHandle dy签名错误", debug)
		response.Result = 11
		response.ErrorMsg = "签名错误"
		c.JSON(200, response)
		return
	}
	openId := QueryRoomIdInterconvertAnchorOpenId(headers["x-roomid"])
	if openId == "" {
		ziLog.Error("BaseCallBackHandle openId is nil", debug)
		response.Result = 12
		response.ErrorMsg = "获取openId null"
		c.JSON(200, response)
	}
	pushDyBasePayloayDirect(headers["x-roomid"], openId, headers["x-msg-type"], *bodyByte)
	response.Result = 1
	c.JSON(200, response)
}
