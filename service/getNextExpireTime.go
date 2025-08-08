package service

import "time"

// 获取到下一次清榜的时间
func GetNextExpireTime() time.Duration {
	daysUntilNextScroll := (int(scrollDay) + 7) % 7
	if daysUntilNextScroll == 0 {
		daysUntilNextScroll = 7
	}
	nextScrollday := time.Now().AddDate(0, 0, daysUntilNextScroll)
	nextScrollDayTime := time.Date(nextScrollday.Year(), nextScrollday.Month(), nextScrollday.Day(), scrollHour, 0, 0, 0, time.Now().Location())
	return time.Until(nextScrollDayTime)
}
