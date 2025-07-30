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
		// {
		// 	if debug {
		// 		log.Println("websocket connect: ", c.Request.Header)
		// 	}
		// 	// 客户端建连
		// 	if c.GetHeader("X-Room-ID") == "" {
		// 		log.Println("websocket connect err: ", "X-Room-ID is null")
		// 		break
		// 	} else {
		// 		//设置积分池为零，删除CurrentRoundId内roomid信息，初始化
		// 		if !connect(c.GetHeader("X-Room-ID"), c.GetHeader("X-Anchor-Openid")) {
		// 			log.Println("链接后房间初始化失败")
		// 		}
		// 		//存储用户信息
		// 		setRoomIdToAnchorOpenId(c.GetHeader("X-Room-ID"), c.GetHeader("X-Anchor-Openid"))
		// 	}
		// 	if err := userInfoStore(UserInfoStruct{
		// 		OpenId:    c.GetHeader("X-Anchor-Openid"),
		// 		AvatarUrl: c.GetHeader("X-Avatar-Url"),
		// 		NickName:  c.GetHeader("X-Nick-Name"),
		// 	}); err != nil {
		// 		log.Println("connect userInfoStore err: ", err)
		// 	}
		// 	data := pmsg.AnchorInfoMessage{
		// 		RoomId:       c.GetHeader("X-Room-ID"),
		// 		AnchorOpenId: c.GetHeader("X-Anchor-Openid"),
		// 		NickName:     c.GetHeader("X-Nick-Name"),
		// 		AvatarUrl:    c.GetHeader("X-Avatar-Url"),
		// 	}
		// 	sData, _ := proto.Marshal(&data)
		// 	if err := pushDownLoadMessage(1, "connect", c.GetHeader("x-client-uuid"), sData); err != nil {
		// 		log.Println("connect pushDownLoadMessage err: ", err)
		// 	}
		// }
	case "disconnect":
		fmt.Println("disconnect")
		// {
		// 	if debug {
		// 		log.Println("websocket disconnect: ", c.Request.Header)
		// 	}

		// 	// 客户端断连
		// 	if c.GetHeader("X-Room-ID") == "" {
		// 		log.Println("websocket disconnect err: ", "X-Room-ID is null")
		// 	}
		// 	//设置积分池为零，删除CurrentRoundId内roomid信息
		// 	endConnect(c.GetHeader("X-Room-ID"), c.GetHeader("X-Anchor-Openid"))
		// 	//删除用户信息
		// 	// DelRoomIdToAnchorOpenId(c.GetHeader("X-Room-ID"))
		// }
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
				ziLog.Error(fmt.Sprintf("websocketMessageFunc err,  websocket_uplink_msg: %s errorInfo: %s ", msg, err), debug)
				if err := sse.SseSend(pmsg.MessageId_BackErrorSend, []string{c.GetHeader("x-client-uuid")}, []byte(err.Error())); err != nil {
					ziLog.Error(fmt.Sprintf("websocketMessageFunc uplink pushDownLoadMessage err: %v", err), debug)
				}
			}
		}
	default:
		break
	}
	c.JSON(200, gin.H{"code": 200, "msg": "connect success"})
}
