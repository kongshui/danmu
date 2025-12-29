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
	if level == 0 {
		return nil
	}
	score, err := rdb.ZIncrBy(level_db, float64(level), openId)
	if err != nil {
		return err
	}
	if score <= 0 {
		DeleteLevelInfo(openId)
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
func ScrollClearLevelInfo(version string) error {
	if !rdb.IsExistKey(level_db) {
		return nil
	}

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
	return nil
}
