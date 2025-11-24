package service

import "time"

// 获取到下一次清榜的时间
func GetNextExpireTime(month bool) time.Duration {
	if month {
		now := time.Now()
		if now.Day() >= scrollTime.ScrollTime.MonthDay {
			nextDay := now.AddDate(0, 1, 0)
			nextMonthDay := time.Date(nextDay.Year(), nextDay.Month(), scrollTime.ScrollTime.MonthDay, scrollTime.ScrollTime.MonthHour, 0, 0, 0, time.Now().Location())
			return time.Until(nextMonthDay)
		} else {
			nextMonthDay := time.Date(now.Year(), now.Month(), scrollTime.ScrollTime.MonthDay, scrollTime.ScrollTime.MonthHour, 0, 0, 0, time.Now().Location())
			return time.Until(nextMonthDay)
		}
	} else {
		daysUntilNextScroll := (int(scrollTime.ScrollTime.WeekDay) + 7*scrollTime.ScrollTime.WeekInterval) % 7 * scrollTime.ScrollTime.WeekInterval
		if daysUntilNextScroll == 0 {
			daysUntilNextScroll = 7 * scrollTime.ScrollTime.WeekInterval
		}
		nextScrollday := time.Now().AddDate(0, 0, daysUntilNextScroll)
		nextScrollDayTime := time.Date(nextScrollday.Year(), nextScrollday.Month(), nextScrollday.Day(), scrollTime.ScrollTime.WeekHour, 0, 0, 0, time.Now().Location())
		return time.Until(nextScrollDayTime)
	}
	// switch week_set {
	// case 0:
	// 	now := time.Now()
	// 	if now.Day() >= month_day {
	// 		nextDay := now.AddDate(0, 1, 0)
	// 		nextMonthDay := time.Date(nextDay.Year(), nextDay.Month(), month_day, scrollHour, 0, 0, 0, time.Now().Location())
	// 		return time.Until(nextMonthDay)
	// 	} else {
	// 		nextMonthDay := time.Date(now.Year(), now.Month(), month_day, scrollHour, 0, 0, 0, time.Now().Location())
	// 		return time.Until(nextMonthDay)
	// 	}
	// default:
	// 	daysUntilNextScroll := (int(scrollDay) + 7*week_set) % 7 * week_set
	// 	if daysUntilNextScroll == 0 {
	// 		daysUntilNextScroll = 7 * week_set
	// 	}
	// 	nextScrollday := time.Now().AddDate(0, 0, daysUntilNextScroll)
	// 	nextScrollDayTime := time.Date(nextScrollday.Year(), nextScrollday.Month(), nextScrollday.Day(), scrollHour, 0, 0, 0, time.Now().Location())
	// 	return time.Until(nextScrollDayTime)
	// }
}
