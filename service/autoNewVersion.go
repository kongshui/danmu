package service

import (
	"context"
	"path"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 检测更新
// func AutoNewVersion() {
// 	//查看周几 Monday Tuesday Wednesday Thursday Friday Saturday Sunday
// 	isFirst := true
// 	t := time.NewTicker(1 * time.Minute)
// 	for {
// 		<-t.C
// 		fmt.Println("现在版本号为：", currentRankVersion)
// 		// 月榜初始化
// 		if time.Now().Hour() == 0 {
// 			// monthVersionSet()
// 		}
// 		if (week_set != 0 && time.Now().Weekday() == scrollDay && time.Now().Hour() == scrollHour && time.Now().Minute() <= 5) || (week_set == 0 && time.Now().Day() == month_day && time.Now().Hour() == scrollHour && time.Now().Minute() <= 5) {
// 			if isFirst {
// 				nowWorldRankVersion := time.Now().Format(version_time_layout)
// 				//比较版本号
// 				nowVersionT, _ := time.Parse(version_time_layout, nowWorldRankVersion)
// 				worldRankVersionT, _ := time.Parse(version_time_layout, currentRankVersion)
// 				//时间不够时间间隔，不轮转
// 				if nowVersionT.Unix()-worldRankVersionT.Unix() < version_time_interval || int64(week_set*7)*24*60*60 > nowVersionT.Unix()-worldRankVersionT.Unix() {
// 					continue
// 				}
// 				// 设置上期前100名名单列表
// 				// go top100Rank(currentVersionRankDb)
// 				// 设置isFirst状态
// 				// 设置版本号
// 				currentRankVersion = nowWorldRankVersion
// 				if !AutoNewVersionLock() {
// 					continue
// 				}
// 				isFirst = false
// 				// 设置level
// 				if is_level_scroll {
// 					// 等级滚动
// 					ScrollClearLevelInfo(currentRankVersion)
// 				}

// 				// 设置排行榜生效版本
// 				if !is_mock {
// 					WorldRankSet(currentRankVersion)
// 				}
// 				if scrollAuto != nil {
// 					scrollAuto(&currentRankVersion)
// 				}

// 				// 版本号初始化
// 				if err := worldRankVersionInit(); err != nil {
// 					ziLog.Error("autoNewVersion 初始化版本信息失败： "+err.Error(), debug)
// 					continue
// 				}

// 			}
// 		} else {
// 			isFirst = true
// 		}
// 	}
// }

// 使用etcd锁，避免多台服务器同时轮转
func AutoNewVersionLock() bool {
	listenId := etcdClient.NewLease(context.Background(), 3)

	ok := etcdClient.PutIfNotExist(context.Background(), path.Join("/", cfg.Project, monitor_auto_new_version_lock), "1", listenId)
	if ok {
		go func(id clientv3.LeaseID) {
			ctx, cancel := context.WithTimeout(context.Background(), 6*time.Minute)
			defer cancel()
			etcdClient.KeepLease(ctx, listenId)
		}(listenId)
	}
	return ok
}
