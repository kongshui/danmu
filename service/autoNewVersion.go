package service

import (
	"fmt"
	"time"
)

// 每周二0点检测更新
func autoNewVersion() {
	//查看周几 Monday Tuesday Wednesday Thursday Friday Saturday Sunday
	isFirst := true
	t := time.NewTicker(1 * time.Minute)
	for {
		<-t.C
		fmt.Println("现在版本号为：", currentRankVersion)
		// 月榜初始化
		if time.Now().Hour() == 0 {
			monthVersionSet()
		}
		if time.Now().Weekday() == scrollDay && time.Now().Hour() == scrollHour {
			if isFirst {
				isFirst = false
				nowWorldRankVersion := time.Now().Format(version_time_layout)
				//比较版本号
				nowVersionT, _ := time.Parse(version_time_layout, nowWorldRankVersion)
				worldRankVersionT, _ := time.Parse(version_time_layout, currentRankVersion)
				//时间不够7天，不轮转
				if nowVersionT.Unix()-worldRankVersionT.Unix() < version_time_interval {
					continue
				}
				// 设置上期前100名名单列表
				// go top100Rank(currentVersionRankDb)

				// 设置版本号
				currentRankVersion = nowWorldRankVersion
				// 设置是否是第一次进行
				if is_level_scroll {
					// 等级滚动
					scrollClearLevelInfo(currentRankVersion)
				}
				// 开启积分滚动
				if is_integral_scroll {
					// 积分滚动
					scrollClearIntegralInfo(currentRankVersion)
				}

				// 设置排行榜生效版本
				if !is_mock {
					worldRankSet(currentRankVersion)
				}
				// 版本号初始化
				if err := worldRankVersionInit(); err != nil {
					ziLog.Error("autoNewVersion 初始化版本信息失败： "+err.Error(), debug)
					continue
				}

			}
		} else {
			isFirst = true
		}
	}
}
