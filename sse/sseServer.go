package sse

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/kongshui/danmu/model/pmsg"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

// 设置sseServer

func SseServer(c *gin.Context) {
	// c.Writer.Header().Set("Content-Type", "text/event-stream")
	type (
		GetMessage struct {
			Uid string `json:"uid"`
		}
	)
	var gm GetMessage
	err := c.ShouldBindJSON(&gm)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "参数错误",
		})
		return
	}
	ch := &ChanSet{
		Ch:     make(chan string),
		Status: true,
	}
	ChanPool.Put(ch)
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	endCh := make(chan struct{})
	defer close(endCh)
	defer close(ch.Ch)
	go sseSendKeepAlive(ch, endCh)
	for {
		select {
		case message := <-ch.Ch:
			if message == "" {
				return
			}
			// 发送数据到客户端
			_, err := c.Writer.Write([]byte(message))
			if err != nil {
				endCh <- struct{}{}
				ch.Status = false
				return // 如果有错误，停止发送数据
			}
			c.Writer.Flush()
		case <-c.Request.Context().Done():
			endCh <- struct{}{}
			ch.Status = false
			return
		}
	}
}

// 定义一个ChanPool结构体
type (
	ChannelPool struct {
		pool    chan *ChanSet
		lock    sync.Mutex
		maxSize int
	}
	ChanSet struct {
		Ch     chan string
		Status bool
	}
)

// NewChanPool 创建一个新的ChanPool
func NewChanPool(maxSize int) *ChannelPool {
	return &ChannelPool{
		pool:    make(chan *ChanSet, maxSize),
		maxSize: maxSize,
	}
}

// Get 从池中获取一个chan，如果没有可用的，则创建一个新的-返回空
func (p *ChannelPool) Get() (*ChanSet, bool) {
	p.lock.Lock()
	defer p.lock.Unlock()
	ctx, cancel := context.WithTimeout(first_ctx, 100*time.Millisecond)
	defer cancel()
	select {
	case ch := <-p.pool:
		return ch, true
	case <-ctx.Done():
		return nil, false
	}
}

// Put 将一个chan放回池中
func (p *ChannelPool) Put(ch *ChanSet) {
	// p.lock.Lock()
	// defer p.lock.Unlock()
	if len(p.pool) < p.maxSize {
		p.pool <- ch // 添加到池中
	}
}

// sse Send
func SseSend(msgId pmsg.MessageId, uidStrList []string, data []byte) error {
	if len(uidStrList) < 1 {
		return errors.New("uid err: uid is nil")
	}
	dataBody := &pmsg.MessageBody{
		MessageId:   uint32(msgId),
		MessageType: msgId.String(),
		MessageData: data,
		Timestamp:   time.Now().UnixMilli(),
	}
	requestBody, err := proto.Marshal(dataBody)
	if err != nil {
		return err
	}
	sData := &pmsg.SseMessage{
		UidList: uidStrList,
		Data:    requestBody,
	}
	count := 0
	for {
		count++
		if count > 3 {
			return fmt.Errorf("get chan error, send data: %v", sData)
		}
		ch, ok := ChanPool.Get()
		if !ok {
			continue
		}
		if ch.Status {
			ch.Ch <- sData.String() + "\n"
		} else {
			ch.Status = false
			ch.Ch <- ""
			continue
		}
		if ch.Status {
			ChanPool.Put(ch)
		}
		return nil
	}
}

// 发送心跳包
func sseSendKeepAlive(ch *ChanSet, c chan struct{}) {
	t := time.NewTicker(5 * time.Second)
	defer t.Stop()
	sData := &pmsg.SseMessage{
		UidList:   []string{},
		MessageId: pmsg.MessageId_Ping,
		Data:      nil,
	}
	for {
		select {
		case <-t.C:
			if !ch.Status {
				return
			}
			ch.Ch <- sData.String() + "\n"
		case <-c:
			return
		}
	}
}
