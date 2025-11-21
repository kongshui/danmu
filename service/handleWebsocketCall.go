package service

import (
	"fmt"
	"log"

	"github.com/kongshui/danmu/model/pmsg"
	"github.com/kongshui/danmu/sse"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

// WebsocketCallbackHandle，抖音云网关调用后端函数
func WebsocketCallbackHandle(c *gin.Context) {
	if is_maintain {
		c.JSON(200, gin.H{
			"errcode": 90001,
			"errmsg":  "维护中",
		})
		return
	}
	msgType := c.GetHeader("x-tt-event-type")
	switch msgType {
	case "connect":
		fmt.Println("connect")
	case "disconnect":
		fmt.Println("disconnect")
	case "uplink":
		{
			msg := msgBodyPool.Get().(*pmsg.MessageBody)
			data := bytePool.Get().(*[]byte)
			defer func() {
				msgBodyPool.Put(msg)
				bytePool.Put(data)
			}()
			// msg := &pmsg.MessageBody{}
			msg.Reset()
			*data, _ = c.GetRawData()
			if err := proto.Unmarshal(*data, msg); err != nil {
				log.Println("websocket uplink err: ", err)
				break
			}
			// 打印客户端上行消息
			// if debug {
			// 	log.Println("websocket_uplink_msg: ", msg)
			// }
			msg.Uuid = c.GetHeader("x-client-uuid")
			if err := websocketMessageFunc(msg); err != nil {
				ziLog.Error(fmt.Sprintf("websocketMessageFunc err,  websocket_uplink_msg: %s errorInfo: %s ", msg.String(), err), debug)
				if err := sse.SseSend(pmsg.MessageId_BackErrorSend, []string{c.GetHeader("x-client-uuid")}, []byte(err.Error()+"msg: "+msg.String())); err != nil {
					ziLog.Error(fmt.Sprintf("websocketMessageFunc uplink pushDownLoadMessage err: %v", err), debug)
				}
			}
		}
	default:
		break
	}
	c.JSON(200, gin.H{"code": 200, "msg": "connect success"})
}
