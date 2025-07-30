package service

import (
	"errors"
	"fmt"
	"time"
)

// 设置上一周期百强榜
func Top100Rank(key string) error {
	// 设置标识
	ok, err := rdb.SetKeyNX(monitor_top_100_ranking_db, "1", 24*time.Hour)
	if err != nil {
		ziLog.Error(fmt.Sprintf("Top100Rank 设置上一周期百强榜标识失败: %v", err), debug)
		return err
	}
	if !ok {
		return nil
	}
	// 获取前100名人员
	users, err := rdb.ZRevRangeWithScores(key, 0, 99)
	if err != nil {
		ziLog.Error(fmt.Sprintf("Top100Rank 获取前100名人员失败: %v", err), debug)
		return err
	}
	// 重命名
	reName := top_100_ranking + "_" + time.Now().Format("20060102")
	if err := rdb.Rename(top_100_ranking, reName); err != nil {
		ziLog.Error(fmt.Sprintf("Top100Rank 连胜币排行rename error: %v", err), debug)
		return errors.New("连胜币排行rename error")
	}
	for _, v := range users {
		if err := rdb.ZAdd(top_100_ranking, v.Score, v.Member.(string)); err != nil {
			ziLog.Error(fmt.Sprintf("Top100Rank 设置前100名人员失败，openId：: %v, 排名：%v,err: %v", v.Member.(string), v.Score, err), debug)
		}
	}
	return nil
}

// 获取上期百强榜
func getTop100Rank() ([]ResultGroupUserRankInfoStruct, error) {
	users, err := rdb.ZRevRange(top_100_ranking, 0, 99)
	if err != nil {
		return nil, err
	}
	var result []ResultGroupUserRankInfoStruct
	for i, v := range users {
		var user ResultGroupUserRankInfoStruct
		user.OpenId = v
		user.Rank = int(i + 1)
		result = append(result, user)
	}
	return result, nil
}
