package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/kongshui/danmu/model/pmsg"
)

// 查询玩家连胜币
func QueryUserWinStreamCoin(openId string) (int64, error) {
	coin, err := rdb.ZScore(winning_streak_coin_db, openId)
	if err != nil {
		return 0, err
	}
	if coin < 0 {
		rdb.ZRem(winning_streak_coin_db, openId)
		return 0, nil
	}
	return int64(coin), nil
}

// 添加玩家连胜币
func AddUserWinStreamCoin(openId string, coin int64) (int64, error) {
	if coin == 0 {
		nCoin, _ := rdb.ZScore(winning_streak_coin_db, openId)
		// if err != nil {
		// 	return 0, err
		// }
		return int64(nCoin), nil
	}
	nCoin, err := rdb.ZIncrBy(winning_streak_coin_db, float64(coin), openId)
	if err != nil {
		ziLog.Error(fmt.Sprintf("添加玩家连胜币失败，玩家OpenId： %v,玩家获得的连胜币为： %v,err： %v", openId, coin, err), debug)
		return 0, err
	}
	if nCoin < 0 {
		rdb.ZRem(winning_streak_coin_db, openId)
		if err := mysql.UpdateCoin(openId, 0); err != nil {
			ziLog.Error(fmt.Sprintf("更新添加玩家连胜币失败，玩家OpenId： %v,玩家获得的连胜币为： %v,err： %v", openId, coin, err), debug)
		}
		return 0, nil
	}
	if err := mysql.UpdateCoin(openId, coin); err != nil {
		ziLog.Error(fmt.Sprintf("更新添加玩家连胜币失败，玩家OpenId： %v,玩家获得的连胜币为： %v,err： %v", openId, coin, err), debug)
	}
	return int64(nCoin), nil
}

// 玩家使用连胜币
func deleteUserWinStreamCoin(openId string, coin int64) (int64, error) {
	if coin <= 0 {
		nCoin, err := rdb.ZScore(winning_streak_coin_db, openId)
		if err != nil {
			return 0, err
		}
		return int64(nCoin), errors.ErrUnsupported
	}
	nCoin, err := rdb.ZIncrBy(winning_streak_coin_db, float64(-coin), openId)
	if err != nil {
		return 0, err
	}
	if int64(nCoin) < 0 {
		sCoin, err := rdb.ZIncrBy(winning_streak_coin_db, float64(coin), openId)
		if err != nil {
			ziLog.Error(fmt.Sprintf("错误删除玩家连胜币之后无法加回，玩家OpenId： %v,玩家消耗的连胜币为： %v", openId, coin), debug)
			return 0, err
		}
		return int64(sCoin), errors.New("cannotUse")
	}
	if err := mysql.UpdateCoin(openId, -coin); err != nil {
		ziLog.Error("更新删除玩家连胜币失败，玩家OpenId： "+openId+",玩家消耗的连胜币为： "+strconv.FormatInt(coin, 10)+",err： "+err.Error(), debug)
	}
	return int64(nCoin), nil
}

// 使用连胜币
func UseWinningStreamCoin(data *pmsg.RequestwinnerstreamcoinMessage) *pmsg.ResponsewinnerstreamcoinMessage {
	canUse := true
	coin, err := deleteUserWinStreamCoin(data.OpenId, data.UseNum)
	if err != nil {
		canUse = false
	}
	useCoin := &pmsg.ResponsewinnerstreamcoinMessage{}
	useCoin.CanUse = canUse
	useCoin.WinningStreamCoin = coin
	useCoin.OpenId = data.OpenId
	useCoin.RoundId = data.RoundId
	useCoin.GiftId = data.GiftId
	useCoin.TimeStamp = data.TimeStamp
	useCoin.RoomId = data.RoomId
	return useCoin
}

// 获得连胜币
func addWinningStreamCoin(data []*pmsg.AddWinnerStreamCoin) *pmsg.ResponseAddWinnerStreamCoinMessage {
	resultData := &pmsg.ResponseAddWinnerStreamCoinMessage{}
	for _, v := range data {
		// if v.AddNum <= 0 {
		// 	continue
		// }
		ziLog.Gift(fmt.Sprintf("玩家id为： %s，获得连胜币为： %v", v.OpenId, v.GetAddNum()), debug)
		coin, err := AddUserWinStreamCoin(v.OpenId, v.GetAddNum())
		if err != nil {
			ziLog.Error(fmt.Sprintf("玩家id为： %s，： %d，获得连胜币失败: %v", v.OpenId, v.GetAddNum(), err), debug)
		}
		resultData.UserList = append(resultData.UserList, &pmsg.ResponseAddWinnerStreamCoin{OpenId: v.OpenId, WinningStreamCoin: coin})
	}
	resultData.TimeStamp = time.Now().UnixMilli()
	return resultData
}

// 查询连胜币信息
func queryWinningStreamCoin(openIdList []string) *pmsg.ResponseAddWinnerStreamCoinMessage {
	resultData := &pmsg.ResponseAddWinnerStreamCoinMessage{}
	for _, openId := range openIdList {
		coin, err := QueryUserWinStreamCoin(openId)
		if err != nil {
			ziLog.Error(fmt.Sprintf("查询玩家id为： %s，连胜币失败: %v", openId, err), debug)
		}
		resultData.UserList = append(resultData.UserList, &pmsg.ResponseAddWinnerStreamCoin{OpenId: openId, WinningStreamCoin: coin})
	}
	return resultData
}

// 滚动连胜币排行
func ScrollWinningStreamCoin() error {
	ziLog.Info("开始滚动连胜币", debug)
	l, err := rdb.ZCard(winning_streak_coin_db)
	if err != nil {
		return errors.New("连胜币排行查询失败")
	}
	if l == 0 {
		return nil
	}
	// 重命名连胜币排行
	reName := winning_streak_coin_db + "_" + time.Now().Format("20060102")
	if err := rdb.Rename(winning_streak_coin_db, reName); err != nil {
		return errors.New("连胜币排行rename error")
	}
	rdb.Expire(reName, 720*time.Hour)
	// 数据库清空连胜币
	if err := mysql.ClearCoin(); err != nil {
		log.Println(err)
	}
	return nil
}
