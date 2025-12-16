package service

import (
	"github.com/kongshui/danmu/model/pmsg"
	"github.com/kongshui/danmu/sse"
)

// SendMessageToGatewayFunc 发送消息到网关函数
func SendMessage(msgId pmsg.MessageId, uidList []string, data []byte) error {
	switch cfg.IsNode {
	case true:
		return sendMessageToGateway(msgId, uidList, data)
	default:
		return sse.SseSend(msgId, uidList, data)
	}

}
