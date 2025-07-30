package service

import (
	"errors"
)

// 连胜次数
func winningStreamCountAdd(stat int, openId string) error {
	if stat == 1 {
		_, err := rdb.ZIncrBy(winning_streak_count_db, 1, openId)
		if err != nil {
			return errors.New("redis添加玩家连胜次数失败: " + err.Error())
		}
		if err := mysql.UpdateWin(openId, 1); err != nil {
			return errors.New("数据库更新玩家连胜次数失败: " + err.Error())
		}
	} else {
		err := rdb.ZRem(winning_streak_count_db, 0, openId)
		if err != nil {
			return errors.New("减少玩家连胜次数失败: " + err.Error())
		}
		if err := mysql.SetWin(openId, 0); err != nil {
			return errors.New("数据库更新玩家连胜次数失败: " + err.Error())
		}
	}
	return nil
}
