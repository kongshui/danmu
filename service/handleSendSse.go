package service

import (
	"fmt"
	"net/http"

	"github.com/kongshui/danmu/model/pmsg"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

// 其他模块通过web发送sse
func OtherSendSse(c *gin.Context) {
	// func(c *gin.Context) {
	// 	data, err := c.GetRawData()
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{
	// 			"message": "GetRawData error",
	// 		})
	// 		ziLog.Write(zilog.Error, fmt.Sprintf("OtherSendSse GetRawData err %v", err), debug)
	// 		return
	// 	}
	// 	log.Println(string(data))
	// }(c)
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "GetRawData error",
		})
		ziLog.Error(fmt.Sprintf("OtherSendSse GetRawData err %v", err), debug)
		return
	}
	sData := &pmsg.SseMessage{}
	if err := proto.Unmarshal(data, sData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unmarshal error",
		})
		ziLog.Error(fmt.Sprintf("OtherSendSse Unmarshal err %v", err), debug)
		return
	}
	if err := sendMessage(sData.GetMessageId(), sData.GetUidList(), sData.GetData()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "send error",
		})
		ziLog.Error(fmt.Sprintf("OtherSendSse sseSend err %v", err), debug)
		return
	}
}
