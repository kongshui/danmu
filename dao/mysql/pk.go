package dao

import (
	"errors"
	"log"
)

// 设置pk获得分数
func (m *MysqlClient) SetPkScores(openId string, isWinner bool) error {

	if isWinner {
		if _, err := m.Client.Exec("insert into pk_count (open_id,winner_count) values(?,1) on DUPLICATE key update winner_count=winner_count+1", openId); err != nil {
			log.Println("mysql SetPkScores", err)
			return errors.New("mysql winner_count SetPkScores: " + err.Error())
		}
	} else {
		if _, err := m.Client.Exec("insert into pk_count (open_id,loser_count) values(?,1) on DUPLICATE key update loser_count=loser_count+1", openId); err != nil {
			log.Println("mysql SetPkScores", err)
			return errors.New("mysql loser_count SetPkScores: " + err.Error())
		}
	}

	return nil
}

// 获取Pk得分
func (m *MysqlClient) GetPkScores(openId string) (int, int, error) {
	winnerCount := 0
	loserCount := 0
	err := m.Client.QueryRow("select winner_count,loser_count from pk_count where open_id=?", openId).Scan(&winnerCount, &loserCount)
	if err != nil {
		log.Println("mysql GetPkScores", err)
		return 0, 0, errors.New("mysql GetPkScores: " + err.Error())
	}
	return winnerCount, loserCount, nil
}
