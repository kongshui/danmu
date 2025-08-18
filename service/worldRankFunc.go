package service

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// 初始化世界排行版相关信息
func worldRankInit() error {
	// 初始化月
	// MonthVersionSet()
	//初始化版本信息
	if err := worldRankVersionInit(); err != nil {
		ziLog.Error(fmt.Sprintf("初始化版本信息失败：: %v", err), debug)
		return err
	}
	//设置分组信息
	if err := playerMatchGroupAdd(); err != nil {
		ziLog.Error(fmt.Sprintf("设置分组信息失败: %v", err), debug)
		return err
	}
	return nil
}

// 当前生效的月版本
func monthVersionSet() error {
	if nowMonth == "" {
		nowMonth = time.Now().Format("01")
		monthVersionRankDb = "month_rank_" + nowMonth
	} else {
		if nowMonth != time.Now().Format("01") {
			nowMonth = time.Now().Format("01")
			monthVersionRankDb = "month_rank_" + nowMonth
		}
	}
	if monthVersionRankDb == "" {
		monthVersionRankDb = "month_rank_" + nowMonth
	}
	return nil
}

// 当前生效版本设置
func worldRankVersionInit() error {
	// 查看世界版本号是否存在
	if rdb.IsExistKey(world_rank_version_db) && currentRankVersion == "" {
		//获取世界版本列表
		version, err := rdb.Get(world_rank_version_db)
		if err != nil {
			return errors.New("worldRankVersionGet失败: " + err.Error())
		}
		currentRankVersion = version
		ok, err := rdb.SetKeyNX(monitor_version_scroll_db, "1", 3*time.Second)
		if err != nil {
			return errors.New("设置版本更新标识失败：" + err.Error())
		}
		if !ok {
			return nil
		}
		if !is_mock {
			if err := setHistoryVersion(); err != nil {
				return errors.New("设置世界历版本失败： " + err.Error())
			}
		}
		return nil
	}
	//设置世界版本列表
	if currentRankVersion == "" {
		currentRankVersion = time.Now().Format(version_time_layout)
	} else {
		version, err := rdb.Get(world_rank_version_db)
		if err != nil {
			return errors.New("worldRankVersionGet失败: " + err.Error())
		}
		if version == currentRankVersion {
			return nil
		}
	}

	return scrollow()
}

// 设置世界历版本
func setHistoryVersion() error {
	//设置版本列表
	length, err := rdb.LLen(world_rank_version_list_db)
	if err != nil {
		return errors.New("worldRankVersionSet 获取版本号长度失败: " + err.Error())
	}

	if length == 0 {
		if err := rdb.RPush(world_rank_version_list_db, currentRankVersion); err != nil {
			return errors.New("worldRankVersionSet 添加版本号失败1 : " + err.Error())
		}
		return nil
	}
	lVersion, err := rdb.LIndex(world_rank_version_list_db, length-1)
	if err != nil {
		return errors.New("worldRankVersionSet 获取版本号失败: " + err.Error())
	}
	if lVersion == currentRankVersion {
		return nil
	}
	if err := rdb.RPush(world_rank_version_list_db, currentRankVersion); err != nil {
		return errors.New("worldRankVersionSet 添加版本号失败: " + err.Error())
	}
	return nil
}

// 添加玩家数据至世界排行榜和历史世界排行榜
func WorldRankNumerAdd(openId string, score float64) error {
	if score <= 0 {
		return nil
	}
	if _, err := rdb.ZIncrBy(world_rank_week, score, openId); err != nil {
		ziLog.Error("添加玩家数据至世界排行榜失败，玩家OpenId： "+openId+",玩家获得的积分为："+strconv.FormatInt(int64(score), 10)+",err： "+err.Error(), debug)
		return err
	}
	if _, err := rdb.ZIncrBy(world_rank_historical_db, score, openId); err != nil {
		ziLog.Error("添加玩家数据至历史世界排行榜失败，玩家OpenId： "+openId+",玩家获得的积分为："+strconv.FormatInt(int64(score), 10)+",err： "+err.Error(), debug)
		return err
	}
	// 添加月榜
	if _, err := rdb.ZIncrBy(monthVersionRankDb, score, openId); err != nil {
		ziLog.Error("添加玩家数据至历史世界排行榜失败，玩家OpenId： "+openId+",玩家获得的积分为："+strconv.FormatInt(int64(score), 10)+",err： "+err.Error(), debug)
		return err
	}
	if err := mysql.UpdateRank(openId, int64(score)); err != nil {
		ziLog.Error("更新添加玩家数据至世界排行榜失败，玩家OpenId： "+openId+",玩家获得的积分为："+strconv.FormatInt(int64(score), 10)+",err： "+err.Error(), debug)
	}
	return nil
}

// // 添加玩家数据至历史世界排行榜
// func worldRankHistoricalAdd(openId string, score float64) error {
// 	if _, err := rdb.ZIncrBy(world_rank_historical_db, score, openId); err != nil {
// 		log.Println("添加玩家数据至历史世界排行榜失败...")
// 		return err
// 	}
// 	return nil
// }

// 设置玩家对局分组名称
func playerMatchGroupAdd() error {
	var (
		groups string
	)

	group := strings.SplitSeq(groups, ",")
	for v := range group {
		if _, err := rdb.SAdd(group_list_db, v); err != nil {
			return err
		}
	}
	return nil
}

/*
设置直播间当前对局名称,roomId+"_"+roundId+"_"+group+"_"+Rank ：内容openid ：score
设置直播间分组信息，分组名称为LiveCurrentRound，内容： 主播roomid ：roundid ，结束后设置为0
设置玩家分组信息，分组名称为：roomid + "_" + roundid + "_" + group，内容：openid ： group，最后设置wingroup
*/
func liveCurrentRoundAdd(roomId string, roundId int64) error {
	if err := rdb.Set(roomId+"_round", roundId, 0); err != nil {
		return err
	}
	return nil
}

// 删除直播间当前对局名称
func liveCurrentRoundDel(roomId string) error {
	if err := rdb.Del(roomId + "_round"); err != nil {
		return err
	}
	return nil
}

// 设置快速返回入口
func fastReturnAdd(roomId, openId string, score float64) {
	if _, ok := queryRoomIdToRoundId(roomId); ok {
		if err := setPlayerDataToRound(roomId, openId, score); err != nil {
			ziLog.Error(fmt.Sprintf("SetPlayer comment score to pool err: roomid: %v, openId: %v, score: %v, err: %v", roomId, openId, 0.5, err), debug)
		}
	}
}

// 如果在

// 设置玩家数据至对局排行榜和积分池,如果玩家暂时没加入组，则计入临时组
func setPlayerDataToRound(roomId, openId string, score float64) error {
	//存储到积分池
	anchorOpenId := QueryRoomIdInterconvertAnchorOpenId(roomId)
	if anchorOpenId == "" {
		return errors.New("setPlayerDataToRound anchorOpenId 未查找到")
	}
	if _, err := addIntegralByScore(anchorOpenId, score); err != nil {
		ziLog.Error("SetPlayerDataToRound err: roomid: "+roomId+", openId: "+openId+", score: "+strconv.FormatInt(int64(score), 10)+", err: "+err.Error(), debug)
	}
	// roundId, ok := queryRoomIdToRoundId(roomId)
	// if !ok {
	// 	return errors.New("setPlayerDataToRound roundId 未查找到")
	// }
	// roundIdStr := strconv.FormatInt(roundId, 10)
	// group, err := rdb.HGet(roomId+"_"+roundIdStr+"_"+"group", openId)
	// if err != nil {
	// 	return setTempGroup(roomId, openId, score)
	// }
	// //对局排行存储名称
	// name := roomId + "_" + roundIdStr + "_" + group + "_rank"
	// _, err = rdb.ZIncrBy(name, score, openId)
	// if err != nil {
	// 	ziLog.Error( "SetPlayerDataToRound rdb.ZIncrBy err: roomid: "+roomId+", openId: "+openId+", score: "+strconv.FormatInt(int64(score), 10)+", err: "+err.Error(), debug)
	// }
	// ttl, _ := rdb.TTL(name)
	// if ttl < 0 {
	// 	//设置过期时间为7天
	// 	rdb.Expire(name, time_interval*time.Second)
	// }
	return nil
}

// 设置临时组
func SetTempGroup(roomId, openId string, score float64) error {
	roundId, ok := queryRoomIdToRoundId(roomId)
	if !ok {
		return errors.New("setTempGroup roundId 未查找到")
	}
	name := roomId + "_" + strconv.FormatInt(roundId, 10) + "_temp_group_rank"
	_, err := rdb.ZIncrBy(name, score, openId)
	if err != nil {
		ziLog.Error(fmt.Sprintf("SetPlayerDataToRound err: name: %v, openId: %v, score: %v, err: %v", name, openId, score, err), debug)
	}
	rdb.Expire(name, 7200*time.Second)
	return nil
}

// 获取玩家世界排行版数据,分数，排名，错误
func getPlayerWorldRankData(openId string) (int64, int64, error) {
	score, err := rdb.ZScore(world_rank_week, openId)
	if err != nil {
		return 0, 0, err
	}
	rank, err := rdb.ZRevRank(world_rank_week, openId)
	if err != nil {
		return 0, 0, err
	}
	rank += 1
	if rank > 100 {
		rank = 0
	}
	return int64(score), rank, nil
}

// 查询玩家是否在分组中
func queryPlayerInGroup(roomId, openId string) (string, int64, bool, error) {
	var group string
	roundId, ok := queryRoomIdToRoundId(roomId)
	if !ok {
		return "", 0, false, errors.New("queryPlayerInGroup roundId 未找到")
	}
	//获取玩家分组
	ok, err := rdb.HExists(roomId+"_"+strconv.FormatInt(roundId, 10)+"_group", openId)
	if err != nil || !ok {
		gOk, _ := rdb.HExists(roomId+"_"+strconv.FormatInt(roundId, 10)+"_group", strings.Split(os.Getenv("GROUP"), ",")[0])
		return "", roundId, gOk, err
	}
	group, err = rdb.HGet(roomId+"_"+strconv.FormatInt(roundId, 10)+"_group", openId)
	if err != nil {
		return "", roundId, false, err
	}
	//获取游戏是否完成
	ok, err = rdb.HExists(roomId+"_"+strconv.FormatInt(roundId, 10)+"_group", group)
	if err != nil {
		return "", roundId, false, err
	}
	return group, roundId, ok, nil
}

// 获取当前版本Db
//
//	func getCurrentVersionDb() string {
//		return "world_rank_" + currentRankVersion
//	}
func scrollWorldRank(version string, count int) error {
	if !rdb.IsExistKey(world_rank_week) {
		return nil
	}

	//获取世界版本列表
	if err := rdb.Rename(world_rank_week, "world_rank_"+version); err != nil {
		ziLog.Error(fmt.Sprintf("scrollWorldRank 滚动世界榜单失败： version: %v, count: %v", version, count), debug)
		time.Sleep(time.Second * 1)
		if count > 60 {
			return err
		}
		count++
		return scrollWorldRank(version, count)
	}
	rdb.Expire("world_rank_"+version, 24*30*time.Hour)
	// rdb.Del(user_info_db)
	return nil
}

// 设置滚动
func scrollow() error {
	// 设置版本更新标识失败：
	ok, err := rdb.SetKeyNX(monitor_version_scroll_db, "1", 24*time.Hour)
	if err != nil {
		return errors.New("设置版本更新标识失败：" + err.Error())
	}
	if !ok {
		return nil
	}
	//初始化版本
	if err := rdb.Set(world_rank_version_db, currentRankVersion, 0); err != nil {
		return errors.New("worldRankVersionSet失败: " + err.Error() + " " + world_rank_version_db)
	}
	//mysql 当前排行版清零
	if err := mysql.ClearRank(); err != nil {
		ziLog.Error("mysql 当前排行版清零失败: "+err.Error(), debug)
	}
	// 滚动积分池
	if err := ScrollWinningStreamCoin(); err != nil {
		ziLog.Error(fmt.Sprintf("初始化连胜币失败: %v", err), debug)
	}
	//设置版本列表
	if err := setHistoryVersion(); err != nil {
		ziLog.Error("设置世界历版本失败： "+err.Error(), debug)
	}
	if err := scrollWorldRank(currentRankVersion, 0); err != nil {
		ziLog.Error("autoNewVersion 滚动世界榜单失败： "+err.Error(), debug)
	}
	// 统计
	go statistic()
	return nil
}

// 获取玩家分组信息,group, round是否结束，1为开始，2为结束
func getUserGroup(roomId, openId string) (string, int, error) {
	var (
		endStatus int = 1
	)
	roundId, err := rdb.Get(roomId + "_round")
	if err != nil {
		return "", 2, err
	}
	//获取玩家分组
	ok, err := rdb.HExists(roomId+"_"+roundId+"_group", openId)
	if err != nil || !ok {
		gOk, _ := rdb.HExists(roomId+"_"+roundId+"_group", strings.Split(os.Getenv("GROUP"), ",")[0])
		if gOk {
			endStatus = 2
		}
		return "", endStatus, err
	}
	gOk, _ := rdb.HExists(roomId+"_"+roundId+"_group", strings.Split(os.Getenv("GROUP"), ",")[0])
	if gOk {
		endStatus = 2
	}
	group, err := rdb.HGet(roomId+"_"+roundId+"_group", openId)
	if err != nil {
		return "", endStatus, err
	}
	return group, endStatus, nil
}
