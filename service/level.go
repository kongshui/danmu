package service

import "time"

// 查询等级信息
func QueryLevelInfo(openId string) (int64, error) {
	level, err := rdb.ZScore(level_db, openId)
	if err != nil {
		return 0, err
	}
	return int64(level), nil
}

// 查询旧库信息
func QueryOldLevelInfo(openId, old_level_db string) (float64, error) {
	level, err := rdb.ZScore(old_level_db, openId)
	if err != nil {
		return 0, err
	}
	return level, nil
}

// 更新等级信息
func UpdateLevelInfo(openId string, level int64) error {
	_, err := rdb.ZIncrBy(level_db, float64(level), openId)
	if err != nil {
		return err
	}
	return nil
}

// 删除等级信息
func DeleteLevelInfo(openId string) error {
	err := rdb.ZRem(level_db, openId)
	if err != nil {
		return err
	}
	return nil
}

// scorll清除等级信息
func scrollClearLevelInfo(version string) error {
	key := level_db + version
	err := rdb.Rename(level_db, key)
	if err != nil {
		return err
	}
	// 旧信息设置过期时间
	err = rdb.Expire(key, 15*24*time.Hour) // 15天过期
	if err != nil {
		return err
	}
	// 读取前store_level名单
	if store_level > 0 {
		topLeverList, err := rdb.ZRevRange(world_rank_week, 0, store_level-1)
		if err != nil {
			ziLog.Error("scrollClearLevelInfo 读取前store_level名单失败： "+err.Error(), debug)
			return err
		}
		for _, openId := range topLeverList {
			level, err := QueryOldLevelInfo(openId, key)
			if err != nil {
				ziLog.Error("scrollClearLevelInfo 读取等级信息失败： "+err.Error(), debug)
			}
			_, err = rdb.ZIncrBy(level_db, float64(level), openId)
			if err != nil {
				ziLog.Error("scrollClearLevelInfo 等级滚动失败： "+err.Error(), debug)
				return err
			}
		}
	}
	return nil
}
