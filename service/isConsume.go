package service

import "fmt"

// 查询是否在第一次消费列表中
func QueryIsConsume(openId string) bool {
	ok, err := rdb.SIsMember(is_consume_db, openId)
	if err != nil {
		return false
	}
	return ok
}

//  对比第一次消费时间
// func compareIsConsume(openId string, timeStamp int64) bool {
// 	if ok := QueryIsConsume(openId); !ok {
// 		return false
// 	}
// 	sTimestamp, err := rdb.HGet(is_consume_db, openId)
// 	if err != nil {
// 		return false
// 	}
// 	if sTimestamp == strconv.FormatInt(timeStamp, 10) {
// 		return true
// 	}
// 	return false
// }

// 设置第一次消费时间
func SetIsConsume(openId string) (int64, error) {
	return rdb.SAdd(is_consume_db, openId)
}

// 删除第一次消费时间
func DelIsConsume(openId string) error {
	return rdb.SRem(is_consume_db, openId)
}

// scroll 滚动删除过期的用户
func ScrollDelIsConsume(label string) {
	name := is_consume_db + "_" + label
	if err := rdb.Rename(is_consume_db, name); err != nil {
		ziLog.Error(fmt.Sprintf("ScrollDelIsConsume error: %v", err), debug)
	}
}
