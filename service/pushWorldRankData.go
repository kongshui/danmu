package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/kongshui/danmu/model/pmsg"
)

// 推送世界榜单数据
func pushWorldRankData() error {
	var (
		//计数
		gCount int64
		//世界榜单
		worldRank WorldRankListStruct
	)
	//初始化app
	worldRank.AppId = app_id
	worldRank.IsOnlineVersion = config.App.IsOnline
	worldRank.WorldRankVersion = currentRankVersion
	worldRank.RankList = make([]UserRankStruct, 0)
	// //获取总长度
	userLen, err := rdb.ZCard(world_rank_week)
	if err != nil {
		return fmt.Errorf("获取总长度失败: %v", err)
	}
	if userLen == 0 {
		return nil
	}
	//查询世界榜单数据
	for range 3 {
		//是否跳出此循环
		isBreak := false
		//stop停止下表
		stop := gCount + 49
		if stop >= 149 || stop >= userLen {
			isBreak = true
			stop = -1
		}
		userScoreList, err := rdb.ZRevRangeWithScores(world_rank_week, gCount, stop)
		if err != nil {
			return fmt.Errorf("获取玩家数据失败: %v", err)
		}
		//添加用户到列表
		for _, v := range userScoreList {
			gCount++
			getUserData := getWorldPlayerData(v.Member.(string), gCount, int64(v.Score))
			if v.Score <= 100 {
				isBreak = true
				break
			}
			worldRank.RankList = append(worldRank.RankList, getUserData)
			if gCount == userLen || gCount == 150 {
				isBreak = true
				break
			}
		}
		if isBreak {
			break
		}
	}
	if len(worldRank.RankList) == 0 {
		return nil
	}
	//上报世界榜单数据
	if !dyWorldRankListUpload(worldRank, url_user_world_rank_upload_url) {
		ziLog.Error(fmt.Sprintf("推送世界榜单数据失败...%v", worldRank), debug)
		return errors.New("推送世界榜单数据失败")
	}
	return nil
}

// 推送世界榜单历史累计数据
func pushHistoryWorldRankData() error {
	var (
		//计数
		gCount int64
		//世界榜单
		worldRank WorldRankUserListStruct
	)
	//初始化app
	worldRank.AppId = app_id
	worldRank.IsOnlineVersion = config.App.IsOnline
	worldRank.WorldRankVersion = currentRankVersion
	// //获取总长度
	userLen, err := rdb.ZCard(world_rank_week)
	if err != nil {
		log.Println("pushHistoryWorldRankData获取总长度失败...", err, userLen)
		return err
	}
	if userLen == 0 {
		return nil
	}
	if debug {
		log.Println("pushHistoryWorldRankData获取总长度成功...", userLen)
	}
	//查询世界榜单数据
	for i := range int(math.Ceil(float64(userLen) / 50)) {
		worldRank.UserList = make([]UserRankStruct, 0)
		if (i%99 == 0) && (i != 0) {
			time.Sleep(time.Second * 1)
		}
		//是否跳出此循环
		isBreak := false
		//stop停止下表
		stop := gCount + 49
		if stop >= userLen {
			isBreak = true
			stop = -1
		}
		userScoreList, err := rdb.ZRevRangeWithScores(world_rank_week, gCount, stop)
		if err != nil {
			log.Println("pushHistoryWorldRankData获取玩家数据失败...", err, userScoreList)
			return err
		}
		//添加用户到列表
		for _, v := range userScoreList {
			gCount++
			getUserData := getWorldPlayerData(v.Member.(string), gCount, int64(v.Score))
			if v.Score <= 100 {
				isBreak = true
				break
			}
			worldRank.UserList = append(worldRank.UserList, getUserData)
			if gCount == userLen {
				isBreak = true
				break
			}
		}

		//上报世界榜单历史数据
		if len(worldRank.UserList) == 0 {
			log.Println("历史世界排行版数据为空...")
			break
		} else {
			if debug {
				log.Println("历史世界排行版数据长度...", len(worldRank.UserList), worldRank.UserList)
			}
			if !dyWorldRankListUpload(worldRank, url_world_rank_user_total_url) {
				return errors.New("推送世界历史榜单数据失败")
			}
		}
		//跳出
		if isBreak {
			break
		}
	}
	if debug {
		log.Println("推送世界榜单历史累计数据成功...", worldRank)
	}
	//上报世界榜单数据结束
	worldRankCompleteUpload()
	return nil
}

// 出啊hi话玩家世界排行数据
func getWorldPlayerData(openId string, rank, score int64) (userData UserRankStruct) {
	winningStreakCountInt, err := rdb.ZScore(winning_streak_coin_db, openId)
	if err != nil {
		winningStreakCountInt = 0
	}
	userData.OpenId = openId
	userData.Rank = rank
	userData.WinningStreakCount = int64(winningStreakCountInt)
	userData.Score = score
	return
}

// 当前世界数据定时推送入口
func pushWorldRankDataEntry() {
	ctx, cancel := context.WithCancel(first_ctx)
	defer cancel()
	t := time.NewTicker(time.Second * 35)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			ziLog.Info("开始推送世界榜单数据...", debug)
			ok, err := rdb.SetKeyNX(monitor_world_push_db, nodeUuid, 30*time.Second)
			if err != nil {
				ziLog.Error(fmt.Sprintf("推送世界榜单数据失败: %v", err), debug)
				continue
			}
			if ok {
				if err := pushWorldRankData(); err != nil {
					ziLog.Error(fmt.Sprintf("推送世界榜单数据失败: %v", err), debug)
				}
				rdb.Del(monitor_world_push_db)
			}
		}
	}
}

// 历史数据推送入口
func pushHistoryWorldRankDataEntry() {
	ctx, cancel := context.WithCancel(first_ctx)
	defer cancel()
	t := time.NewTicker(time.Minute * 10)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if debug {
				log.Println("开始推送历史世界榜单数据...")
			}
			ok, err := rdb.SetKeyNX(monitor_world_history_push_db, nodeUuid, 8*time.Minute)
			if err != nil {
				log.Println(err)
				continue
			}
			if ok {
				if err := pushHistoryWorldRankData(); err != nil {
					log.Println(err)
				}
				rdb.Del(monitor_world_history_push_db)
			}
		}
	}
}

// 根据用户返回用户世界排行数据
func GetWorldRankData(openIdList []string) []WorldInfoStruct {
	var worldInfoList []WorldInfoStruct
	for _, openId := range openIdList {
		score, rank, _ := getPlayerWorldRankData(openId)
		coin, _ := QueryUserWinStreamCoin(openId)
		worldInfoList = append(worldInfoList, WorldInfoStruct{
			OpenId:            openId,
			Rank:              rank,
			Score:             score,
			WinningStreamCoin: coin,
		})
	}
	return worldInfoList
}

// 返回世界排行榜前100名
func getTopWorldRankData() *pmsg.UserInfoListMessage {
	data := &pmsg.UserInfoListMessage{}
	openIdList, err := rdb.ZRevRangeWithScores(world_rank_week, 0, 99)
	if err != nil {
		ziLog.Error(fmt.Sprintf("getTopWorldRankData err: %v", err), debug)
		return data
	}
	for i, userInfo := range openIdList {
		openId := userInfo.Member.(string)
		user, _ := userInfoGet(openId)
		coin, _ := QueryUserWinStreamCoin(openId)
		level, _ := QueryLevelInfo(openId)
		if user.NickName == "" || user.AvatarUrl == "" {
			// 从数据库查询玩家信息
			avatarUrl, nickName, err := mysql.QueryPlayerInfo(openId)
			if err != nil {
				ziLog.Error(fmt.Sprintf("getTopWorldRankData QueryPlayerInfo err: %v,openId: %v", err, openId), debug)

			}
			user.NickName = nickName
			user.AvatarUrl = avatarUrl
		}

		data.UserInfoList = append(data.UserInfoList, &pmsg.UserInfo{
			OpenId:            openId,
			Rank:              int64(i + 1),
			Score:             int64(userInfo.Score),
			AvatarUrl:         user.AvatarUrl,
			NickName:          user.NickName,
			WinningStreamCoin: coin,
			Level:             level,
		})
	}
	return data
}

// 返回世界排行榜前100名
func getTopMonthRankData() *pmsg.UserInfoListMessage {
	data := &pmsg.UserInfoListMessage{}
	openIdList, err := rdb.ZRevRangeWithScores(monthVersionRankDb, 0, 99)
	if err != nil {
		return data
	}
	for i, userInfo := range openIdList {
		openId := userInfo.Member.(string)
		user, _ := userInfoGet(openId)
		coin, _ := QueryUserWinStreamCoin(openId)
		level, _ := QueryLevelInfo(openId)

		data.UserInfoList = append(data.UserInfoList, &pmsg.UserInfo{
			OpenId:            openId,
			Rank:              int64(i + 1),
			Score:             int64(userInfo.Score),
			AvatarUrl:         user.AvatarUrl,
			NickName:          user.NickName,
			WinningStreamCoin: coin,
			Level:             level,
		})
	}
	return data
}
