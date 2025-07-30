package dao

import "errors"

// 统计本周期内花费
func (m *MysqlClient) StatisticCost(openId, startTime, endTime string) (int64, error) {
	if !m.isUse {
		return 0, nil
	}

	var cost int64
	if startTime == "" {
		if err := m.Client.QueryRow("select sum(gift_value) from log_gift where create_at < ? and open_id = ?", endTime, openId).Scan(&cost); err != nil {
			return 0, errors.New("StatisticCost 0 err: " + err.Error())
		}
	} else {
		if err := m.Client.QueryRow("select sum(gift_value) from log_gift where create_at>=? and create_at < ? and open_id = ?", startTime, endTime, openId).Scan(&cost); err != nil {
			return 0, errors.New("StatisticCost err: " + err.Error())
		}
	}

	return cost, nil
}
