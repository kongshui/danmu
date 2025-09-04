package service

import (
	"strconv"
	"time"
)

func setRoundEndGroup(roomId string, roundId int64, groupIdList []GroupResultList) {
	name := roomId + "_" + strconv.FormatInt(roundId, 10) + "_group"
	rdb.Expire(name, 21600*time.Second)
	// 初始化分组
	for _, v := range groupIdList {
		// 初始化分组
		rdb.HSet(name, v.GroupId, v.Result)
	}
}
