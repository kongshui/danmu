package service

import "strconv"

// 查询是否在第一次消费列表中
func queryIsConsume(openId string) bool {
	ok, err := rdb.HExists(is_consume_db, openId)
	if err != nil {
		return false
	}
	return ok
}

//  对比第一次消费时间
func compareIsConsume(openId string, timeStamp int64) bool {
	if ok := queryIsConsume(openId); !ok {
		return false
	}
	sTimestamp, err := rdb.HGet(is_consume_db, openId)
	if err != nil {
		return false
	}
	if sTimestamp == strconv.FormatInt(timeStamp, 10) {
		return true
	}
	return false
}

// 设置第一次消费时间
func setIsConsume(openId string, timeStamp int64) (bool, error) {
	return rdb.HSetNX(is_consume_db, openId, timeStamp)
}
