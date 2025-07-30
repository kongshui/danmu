package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func StartGameHandle(c *gin.Context) {
	if is_maintain {
		c.JSON(403, gin.H{
			"errcode": 90001,
			"errmsg":  "维护中",
		})
		return
	}
	type RoomId struct {
		RoomID    string `json:"room_id"`
		TimeStamp int64  `json:"time_stamp"`
	}
	var (
		roomID RoomId
	)
	// 解析参数
	if err := c.ShouldBindJSON(&roomID); err != nil {
		ziLog.Error(fmt.Sprintf("start game 解析参数错误,err: %v", err), debug)
		return
	}

	ok := rdb.IsExistKey(roomID.RoomID + "_round")
	if !ok {
		ziLog.Error(fmt.Sprintf("start game 房间不存在, roomId: %v", roomID.RoomID), debug)
		c.JSON(400, gin.H{
			"code": 1,
			"msg":  "房间不存在",
		})
		return
	}

	//开启推送任务

	// if !startFinishGameInfo(roomID.RoomID, url_BindUrl, "start") {
	// 	log.Println("start game 游戏开始失败", roomID.RoomID)
	// 	c.JSON(400, gin.H{
	// 		"code": 1,
	// 		"msg":  "游戏开始失败",
	// 	})
	// 	return
	// }

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "游戏开始",
	})
}
