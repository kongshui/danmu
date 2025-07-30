package service

import (
	"github.com/kongshui/danmu/sse"

	"github.com/gin-gonic/gin"
)

// 设置sseServer

func SseServerHandler(c *gin.Context) {
	sse.SseServer(c)
	// c.Writer.Header().Set("Content-Type", "text/event-stream")
	// type (
	// 	GetMessage struct {
	// 		Uid string `json:"uid"`
	// 	}
	// )
	// var gm GetMessage
	// err := c.ShouldBindJSON(&gm)
	// if err != nil {
	// 	c.JSON(400, gin.H{
	// 		"message": "参数错误",
	// 	})
	// 	return
	// }
	// ch := &ChanSet{
	// 	Ch:     make(chan string),
	// 	Status: true,
	// }
	// chanPool.Put(ch)
	// c.Header("Content-Type", "text/event-stream")
	// c.Header("Cache-Control", "no-cache")
	// c.Header("Connection", "keep-alive")
	// c.Header("Access-Control-Allow-Origin", "*")
	// for {
	// 	select {
	// 	case message := <-ch.Ch:
	// 		// 发送数据到客户端
	// 		_, err := c.Writer.Write([]byte(message))
	// 		if err != nil {
	// 			return // 如果有错误，停止发送数据
	// 		}
	// 		c.Writer.Flush()
	// 	case <-c.Request.Context().Done():
	// 		ch.Status = false
	// 		return
	// 	}
	// }
}

// 定义一个ChanPool结构体
// type (
// 	ChannelPool struct {
// 		pool    chan *ChanSet
// 		lock    sync.Mutex
// 		maxSize int
// 	}
// 	ChanSet struct {
// 		Ch     chan string
// 		Status bool
// 	}
// )

// // NewChanPool 创建一个新的ChanPool
// func NewChanPool(maxSize int) *ChannelPool {
// 	return &ChannelPool{
// 		pool:    make(chan *ChanSet, maxSize),
// 		maxSize: maxSize,
// 	}
// }

// // Get 从池中获取一个chan，如果没有可用的，则创建一个新的
// func (p *ChannelPool) Get() (*ChanSet, bool) {
// 	p.lock.Lock()
// 	defer p.lock.Unlock()
// 	ctx, cancel := context.WithTimeout(first_ctx, 50*time.Millisecond)
// 	defer cancel()
// 	select {
// 	case ch := <-p.pool:
// 		return ch, true
// 	case <-ctx.Done():
// 		return nil, false
// 	}
// }

// // Put 将一个chan放回池中
// func (p *ChannelPool) Put(ch *ChanSet) {
// 	p.lock.Lock()
// 	defer p.lock.Unlock()

// 	if len(p.pool) < p.maxSize {
// 		p.pool <- ch // 添加到池中

// 	} else {
// 		close(ch.Ch) // 如果池已满，关闭并丢弃chan（可选）
// 	}
// }

// // sse Send
// func sseSend(msgId pmsg.MessageId, uidStrList []string, data []byte) error {
// 	if len(uidStrList) < 1 {
// 		return errors.New("uid err: uid is nil")
// 	}
// 	dataBody := &pmsg.MessageBody{
// 		MessageId:   uint32(msgId),
// 		MessageType: msgId.String(),
// 		MessageData: data,
// 		Timestamp:   time.Now().UnixMilli(),
// 	}
// 	requestBody, err := proto.Marshal(dataBody)
// 	if err != nil {
// 		return err
// 	}
// 	sData := &pmsg.SseMessage{
// 		UidList: uidStrList,
// 		Data:    requestBody,
// 	}
// 	count := 0
// 	for {
// 		count++
// 		if count > 3 {
// 			return fmt.Errorf("get chan error, send data: %v", sData)
// 		}
// 		ch, ok := chanPool.Get()
// 		if !ok {
// 			time.Sleep(10 * time.Millisecond)
// 			continue
// 		}
// 		if ch.Status {
// 			ctx, cancel := context.WithTimeout(first_ctx, 100*time.Millisecond)
// 			defer cancel()
// 			select {
// 			case ch.Ch <- sData.String() + "\n":
// 			case <-ctx.Done():
// 				if !ch.Status {
// 					close(ch.Ch)
// 					continue
// 				}
// 			}
// 		} else {
// 			close(ch.Ch)
// 			continue
// 		}
// 		if ch.Status {
// 			chanPool.Put(ch)
// 		}
// 		return nil
// 	}
// }
