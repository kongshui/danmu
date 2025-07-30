package service

// 获取roomid接口
// func GetRoomIdHandle(c *gin.Context) {
// 	type Token struct {
// 		Token string `json:"token"`
// 		Uuid  string `json:"uuid"`
// 	}
// 	var token Token
// 	if err := c.ShouldBindJSON(&token); err != nil {
// 		log.Println(err)
// 		c.JSON(404, gin.H{
// 			"err": err,
// 		})
// 		return
// 	}
// 	roomInfo, err := getRoomId(token.Token, token.Uuid)
// 	if err != nil {
// 		log.Println(err)
// 		c.JSON(404, gin.H{
// 			"err": err,
// 		})
// 		return
// 	}

// 	c.JSON(200, RoomInfoStruct{
// 		RoomId:       roomInfo.RoomId,
// 		AnchorOpenId: roomInfo.AnchorOpenId,
// 		NickName:     roomInfo.NickName,
// 		AvatarUrl:    roomInfo.AvatarUrl,
// 	})
// }
