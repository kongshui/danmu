package service

import (
	"github.com/gin-gonic/gin"
)

func FinishGameHandle(c *gin.Context) {
	type RoomId struct {
		RoomID    string `json:"room_id"`
		TimeStamp int64  `json:"time_stamp"`
	}
	var (
		roomID RoomId
	)
	// 解析参数
	if err := c.ShouldBindJSON(&roomID); err != nil {
		c.JSON(400, gin.H{
			"code": 1,
			"msg":  "解析roomid失败",
		})
		return
	}

	// if !startFinishGameInfo(roomID.RoomID, url_BindUrl, "stop") {
	// 	log.Println("finish game 游戏结束失败", roomID.RoomID)
	// 	c.JSON(400, gin.H{
	// 		"code": 2,
	// 		"msg":  "游戏结束失败",
	// 	})
	// 	return
	// }
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "游戏开始",
	})
}
