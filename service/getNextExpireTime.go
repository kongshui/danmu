package service

import "time"

// 获取到下一次清榜的时间
func GetNextExpireTime() time.Duration {
	switch week_set {
	case 0:
		return 0
	default:
		daysUntilNextScroll := (int(scrollDay) + 7*week_set) % 7 * week_set
		if daysUntilNextScroll == 0 {
			daysUntilNextScroll = 7 * week_set
		}
		nextScrollday := time.Now().AddDate(0, 0, daysUntilNextScroll)
		nextScrollDayTime := time.Date(nextScrollday.Year(), nextScrollday.Month(), nextScrollday.Day(), scrollHour, 0, 0, 0, time.Now().Location())
		return time.Until(nextScrollDayTime)
	}

}
