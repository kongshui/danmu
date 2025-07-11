package dao

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func NatsInit(addr []string) *nats.Conn {
	opts := nats.GetDefaultOptions()
	// opts.Url = "nats://localhost:4222"
	opts.Timeout = 10 * time.Second
	opts.MaxReconnect = 10
	opts.MaxPingsOut = 10
	opts.Servers = addr
	// opts.Dialer.KeepAliveConfig.Enable = true
	natsClient, err := opts.Connect()
	if err != nil {
		log.Println("初始化nats失败， err: ", err)
	}
	return natsClient
}
