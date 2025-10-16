package service

import (
	"fmt"

	"github.com/go-redis/redis"
)

// 获取redis有序集合数据
func GetRedisZsetData(key string, startIndex int64, endIndex int64, reverse bool) []redis.Z {
	var (
		err        error
		openIdList []redis.Z
		end_index  int64
	)
	weekRankLen, _ := rdb.ZCard(key)
	if endIndex >= weekRankLen-1 {
		end_index = -1
	} else if startIndex > weekRankLen && weekRankLen != -1 {
		return openIdList
	} else if startIndex < 0 {
		return openIdList
	}
	if !reverse {
		openIdList, err = rdb.ZRevRangeWithScores(key, startIndex, end_index)
	} else {
		openIdList, err = rdb.ZRangeWithScores(key, startIndex, end_index)
	}
	if err != nil {
		ziLog.Error(fmt.Sprintf("GetRedisZsetData err: %v", err), debug)
		return openIdList
	}
	return openIdList
}
