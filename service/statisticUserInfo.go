package service

import (
	"fmt"
	"os"
	"time"

	"github.com/kongshui/danmu/common"

	"github.com/xuri/excelize/v2"
)

// 统计信息
type StatisticInfo struct {
	OpenId    string `json:"open_id"`    // 用户Id
	NickName  string `json:"nick_name"`  // 昵称
	Rank      int    `json:"rank"`       // 排名
	Scores    int64  `json:"scores"`     // 总得分
	Coin      int64  `json:"coin"`       // 总连胜币
	TotalCost int64  `json:"total_cost"` // 总花费
}

// 统计
func StatisticTemplate() {
	var (
		startTime string
		endTime   string
	)
	//获取到本期版本号
	length, err := rdb.LLen(world_rank_version_list_db)
	if err != nil {
		ziLog.Error("Statistic rdb LLen err: "+err.Error(), debug)
		return
	}
	if length == 0 {
		ziLog.Error("Statistic length is nil", debug)
		return
	}
	if length == 1 {
		// 只有一个版本号，直接统计
	} else {
		startTimeStamp, _ := rdb.LIndex(world_rank_version_list_db, length-2)
		startTimeParse, _ := time.Parse(version_time_layout, startTimeStamp)
		startTimeParse = time.Date(startTimeParse.Year(), startTimeParse.Month(), startTimeParse.Day(), scrollTime.ScrollTime.WeekHour, 0, 0, 0, startTimeParse.Location())
		startTime = startTimeParse.Format(mysql_query_time_layout)
	}
	endTimeParse, _ := time.Parse(version_time_layout, currentRankVersion)
	endTimeParse = time.Date(endTimeParse.Year(), endTimeParse.Month(), endTimeParse.Day(), scrollTime.ScrollTime.WeekHour, 0, 0, 0, endTimeParse.Location())
	endTime = endTimeParse.Format(mysql_query_time_layout)
	top100, err := rdb.ZRevRangeWithScores("world_rank_"+currentRankVersion, 0, 100)
	if err != nil {
		ziLog.Error("Statistic rdb ZRevRangeWithScores err: "+err.Error(), debug)
		return
	}
	userInfos := make([][]any, 0)
	userInfos = append(userInfos, []any{"用户id", "昵称", "排名", "积分", "连胜币", "胜点", "充值总额"})
	for i, v := range top100 {
		// 统计用户信息
		userInfo := make([]any, 0)
		fmt.Println(i, v.Member, v.Score)
		openId := v.Member.(string)
		user, _ := UserInfoGet(openId, false)
		coin, _ := QueryUserWinStreamCoin(openId)
		cost, _ := mysql.StatisticCost(openId, startTime, endTime)
		winPoint, _ := QueryUserWinningPoint(openId)
		userInfo = append(userInfo, openId)
		userInfo = append(userInfo, user.NickName)
		userInfo = append(userInfo, i+1)
		userInfo = append(userInfo, v.Score)
		userInfo = append(userInfo, coin)
		userInfo = append(userInfo, winPoint)
		userInfo = append(userInfo, cost)
		userInfos = append(userInfos, userInfo)
	}
	writeExcel(userInfos)
}

// 向excel中填写数据
func writeExcel(userInfos [][]any) {

	f := excelize.NewFile()
	// 创建一个工作表
	defer func() {
		if err := f.Close(); err != nil {
			ziLog.Error("WriteExcel Close err: "+err.Error(), debug)
		}
	}()
	for idx, row := range userInfos {
		cell, err := excelize.CoordinatesToCellName(1, idx+1)
		if err != nil {
			ziLog.Error("WriteExcel write err: "+err.Error(), debug)
			return
		}
		f.SetSheetRow("Sheet1", cell, &row)
	}
	if !common.PathExists("./static") {
		os.Mkdir("./static", 0755)
	}
	// 保存文件
	err := f.SaveAs("./static/" + currentRankVersion + ".xlsx")
	if err != nil {
		ziLog.Error("WriteExcel save err: "+err.Error(), debug)
		return
	}
	ziLog.Info("WriteExcel success", debug)
}
