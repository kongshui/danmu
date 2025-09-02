package service

import (
	"errors"
	"time"
)

// 查询玩家胜点
func QueryUserWinningPoint(openId string) (int64, error) {
	point, err := rdb.ZScore(winning_point_db, openId)
	if err != nil {
		return 0, err
	}
	if point < 0 {
		rdb.ZAdd(winning_point_db, 0, openId)
		return 0, nil
	}
	return int64(point), nil
}

// 增加玩家胜点
func AddUserWinningPoint(openId string, point int64) (int64, error) {
	newPoint, err := rdb.ZIncrBy(winning_point_db, float64(point), openId)
	if err != nil {
		return 0, err
	}
	if newPoint < 0 {
		rdb.ZAdd(winning_point_db, 0, openId)
		return 0, nil
	}
	return int64(newPoint), nil
}

// 滚动连胜币排行
func ScrollWinningPoint() error {
	ziLog.Info("开始滚动胜点", debug)
	l, err := rdb.ZCard(winning_point_db)
	if err != nil {
		return errors.New("胜点排行查询失败")
	}
	if l == 0 {
		return nil
	}
	// 重命名胜点排行
	reName := winning_point_db + "_" + time.Now().Format("20060102")
	if err := rdb.Rename(winning_point_db, reName); err != nil {
		return errors.New("胜点排行rename error")
	}
	rdb.Expire(reName, 720*time.Hour)
	return nil
}
