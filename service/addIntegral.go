package service

import (
	"errors"
	"strconv"
)

// 通过节点增加积分
func addIntegralByNode(key string, addscore float64) (float64, error) {
	score, err := rdb.IncrByFloat(integral_pool_Prefix+key, float64(addscore))
	if err != nil {
		return 0, errors.New("AddIntegral err: " + err.Error())
	}
	// if ttl, _ := rdb.TTL(integral_pool_Prefix + key); ttl < 0 {
	// 	rdb.Expire(integral_pool_Prefix+key, 10800*time.Second)
	// }
	return score, nil
}

// 通过分数增加积分
func AddIntegralByScore(key string, score float64) (float64, error) {
	score, err := rdb.IncrByFloat(integral_pool_Prefix+key, score)
	if err != nil {
		return 0, errors.New("addIntegralByScore err: " + err.Error())
	}
	// if ttl, _ := rdb.TTL(integral_pool_Prefix + key); ttl < 0 {
	// 	rdb.Expire(integral_pool_Prefix+key, 10800*time.Second)
	// }
	return score, nil
}

// 获取积分
func GetIntegral(key string) (float64, error) {
	score, err := rdb.Get(integral_pool_Prefix + key)
	if err != nil {
		return 0, errors.New("getIntegral err: " + err.Error())
	}
	strconv.ParseFloat(score, 64)
	return strconv.ParseFloat(score, 64)
}

// 删除积分池，只有匹配用
func delIntegral(key string) error {
	return rdb.Del(integral_pool_Prefix + key)
}

// 增加积分池和个人积分
func addIntegralAndUserIntegral(anchorOpenId, openId string, addscore float64) error {
	WorldRankNumerAdd(openId, addscore)
	_, err := rdb.IncrByFloat(integral_pool_Prefix+anchorOpenId, addscore)
	if err != nil {
		return errors.New("addIntegralAndUserIntegral err: " + err.Error())
	}
	// if ttl, _ := rdb.TTL(integral_pool_Prefix + key); ttl < 0 {
	// 	rdb.Expire(integral_pool_Prefix+key, 10800*time.Second)
	// }
	return nil
}
