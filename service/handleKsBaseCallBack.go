package service

import (
	"encoding/json"
	"fmt"

	"github.com/kongshui/danmu/common"

	"github.com/gin-gonic/gin"
)

// 点赞评论等回调接口
func BaseKsCallBackHandle(c *gin.Context) {
	var (
		response KsCallbackRespondStruct
	)
	data := baseDataPool.Get().(*KsCallbackStruct)
	bodyByte := bytePool.Get().(*[]byte)
	defer func() {
		baseDataPool.Put(data)
		bytePool.Put(bodyByte)
	}()
	*bodyByte, _ = c.GetRawData()
	// if err != nil {
	// 	ziLog.Write(logError, fmt.Sprintf("BaseCallBackHandle 获取body错误,data: %v ,err: %v", string(bodyByte), err), debug)
	// 	response.Result = 10
	// 	response.ErrorMsg = "获取body错误：: " + err.Error()
	// 	c.JSON(200, response)
	// 	return
	// }
	if !common.KsCheckSignature(string(*bodyByte), app_secret, c.GetHeader("kwaisign")) {
		ziLog.Error("BaseCallBackHandle 签名错误,err", debug)
		response.Result = 11
		response.ErrorMsg = "签名错误"
		c.JSON(200, response)
		return
	}
	if err := json.Unmarshal(*bodyByte, data); err != nil {
		ziLog.Error(fmt.Sprintf("BaseCallBackHandle 解析参数错误,err: %v", err), debug)
		response.Result = 13
		response.ErrorMsg = "解析参数错误: " + err.Error()
		c.JSON(200, response)
	}
	ksPushBasePayloay(*data)

	response.Result = 1
	c.JSON(200, response)
}
