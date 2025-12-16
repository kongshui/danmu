package service

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/kongshui/danmu/model/pmsg"
	"google.golang.org/protobuf/proto"
)

// 创建num位字符串
func CreateIntStr(num int) string {
	var numStr string
	// 先创建数字字符串
	r := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().Day())))

	for range num {
		randNum := r.Int64N(10)
		numStr += strconv.FormatInt(randNum, 10)

	}
	return numStr
}

// 创建随机字符串
func CreateRandomStr(n int) string {
	r := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(n)))
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, n)
	for i := range b {
		b[i] = letters[r.IntN(len(letters))]
	}
	return string(b)
}

// create room
func CreateRoomId(uid, roomId string) {
	switch platform {
	case "dy":
		data := &pmsg.AnchorInfoMessage{}
		data.AnchorOpenId = CreateIntStr(10)
		data.AvatarUrl = "http://tupian.geimian.com/pic/2016/11/2016-11-05_213153.jpg"
		data.NickName = CreateRandomStr(6)
		data.RoomId = CreateIntStr(12)
		databyte, _ := proto.Marshal(data)
		setRoomInfo(uid, data)
		connect(data.GetRoomId(), data.GetAnchorOpenId())
		if err := sendMessage(pmsg.MessageId_TokenAck, []string{uid}, databyte); err != nil {
			ziLog.Error(fmt.Sprintf("DyGetAnchorInfo pushDownLoadMessage err: %v", err), debug)
		}
	case "ks":
		roomInfo := &pmsg.AnchorInfoMessage{}
		roomInfo.RoomId = roomId
		roomInfo.AnchorOpenId = CreateIntStr(10)
		roomInfo.NickName = CreateRandomStr(6)
		roomInfo.AvatarUrl = "http://tupian.geimian.com/pic/2016/11/2016-11-05_213153.jpg"
		dataByte, _ := json.Marshal(roomInfo)
		if err := sendMessage(pmsg.MessageId_StartBindAck, []string{uid}, dataByte); err != nil {
			ziLog.Error(fmt.Sprintf("KsCreateRoomId pushDownLoadMessage err: %v", err), debug)
		}
	default:
	}
}
