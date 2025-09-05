package service

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/kongshui/danmu/common"
)

// 用户对局数据上传
func usersRoundUpload(roomid, anchorOpenId string, result RoundUploadStruct) error {
	var err error
	if setWinnerScore != nil {
		if err = setWinnerScore(anchorOpenId, result); err != nil {
			ziLog.Error(fmt.Sprintf("usersRoundUpload setWinnerScore err: %v", err), debug)
		}
	}

	switch platform {
	case "ks":
		return err
	case "dy":
		if is_mock {
			return nil
		}
		return dyUserRoundUpload(roomid, anchorOpenId, result)
	}
	return nil
}

// 抖音数据上传
func dyUserRoundUpload(roomId, anchorOpenId string, roundData RoundUploadStruct) error {
	var (
		//用户对局数据
		userData UploadUserGameStruct
		//对局榜单
		userListData []UserListStruct
		//修改后玩家数据存盘
		//原始数据存盘
	)

	//初始化用户对局数据
	userData.AnchorOpenId = anchorOpenId
	userData.AppId = app_id
	userData.RoomId = roomId
	userData.RoundId = roundData.RoundId
	//存储最终结果
	resultName := roomId + "_" + strconv.FormatInt(roundData.RoundId, 10) + "_result"
	for _, v := range roundData.GroupUserList {
		var (
			user UserListStruct
		)
		user.GroupId = v.GroupId
		user.OpenId = v.OpenId
		user.RoundResult = v.RoundResult
		user.Score = v.Score

		// 获取玩家连胜币
		coin, _ := QueryUserWinStreamCoin(v.OpenId)
		user.WinningStreakCount = coin
		userListData = append(userListData, user)
	}
	sort.Slice(userListData, func(i, j int) bool {
		return userListData[i].Score > userListData[j].Score
	})
	for i := range userListData {
		userListData[i].Rank = int64(i + 1)
	}

	//上报用户对局
	for i := range int(math.Ceil(float64(len(userListData)) / 50)) {
		var (
			count   int
			isBreak bool
		)
		//上报用户对局数据
		if (i+1)*50-1 >= len(userListData) {
			count = len(userListData)
			isBreak = true
		} else {
			count = (i+1)*50 - 1
		}
		userData.UserList = userListData[i*50 : count]
		if len(userData.UserList) == 0 {
			break
		}
		//开始上报用户对局数据，最多50条
		if err := uploadUserGameResult(userData, url_round_user_result_upload_url); err != nil {
			return fmt.Errorf("dyUserRoundUpload 上报用户对局数据失败: %v, err: %v", resultName, err)
		}
		//积分等于0分不上报
		if userData.UserList[len(userData.UserList)-1].Score == 0 {
			break
		}
		if isBreak {
			break
		}
	}

	//上报用户对局排行，最多150
	count := 0
	count = min(len(userListData), 150)
	if err := uploadUserGameResult(UploadRankGameStruct{
		AnchorOpenId: anchorOpenId,
		AppId:        app_id,
		RoomId:       roomId,
		RoundId:      roundData.RoundId,
		RankList:     userListData[0:count],
	}, url_round_user_rank_upload_url); err != nil {
		return fmt.Errorf("dyUserRoundUpload 上报用户对局排行失败: %v, err: %v", resultName, err)
	}
	if err := uploadUserGameResult(UploadUserGameCompleteStruct{
		AnchorOpenId: anchorOpenId,
		AppId:        app_id,
		RoomId:       roomId,
		RoundId:      roundData.RoundId,
		CompleteTime: time.Now().Unix(),
	}, url_round_user_upload_complete_url); err != nil {
		return fmt.Errorf("dyUserRoundUpload 上报用户对局完成失败: %v, err: %v", resultName, err)
	}
	return nil
}

// 上报用户对局数据,对局结束后，上报用户对局数据，分批次上报，单批次最多上报50个用户对局数据。
// 本接口是上传所有参与本局玩法并且有战绩的用户数据，批量上报，底层是以用户id+room_id+round_id为维度存储，用户id+room_id+round_id维度重复上报，是覆盖写的逻辑；
// 上报对局榜单列表也用此函数，对局结束后，上报对局榜单列表，榜单列表是指Top 150 的用户数据。如果接口调用的列表长度为 20，则表示榜单列表只展示Top 20的用户数据。
// 上报对局榜单列表时，当对局结束后，调用一次，一次性上报对局榜单Top 150，
func uploadUserGameResult(data any, url string) error {
	headers := map[string]string{
		"Content-Type": "application/json",
		"X-Token":      accessToken.Token,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("uploadUserGameResult json.Marshal err: %v", err)
	}

	response, err := common.HttpRespond("POST", url, body, headers)
	if err != nil {
		return fmt.Errorf("uploadUserGameResult response err: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return fmt.Errorf("uploadUserGameResult response.StatusCode: %v", response.StatusCode)
	}
	var (
		request any
	)
	if err := json.NewDecoder(response.Body).Decode(&request); err != nil {
		return fmt.Errorf("uploadUserGameResult json.NewDecoder err: %v", err)
	}

	if int64(request.(map[string]any)["errcode"].(float64)) != 0 {
		return fmt.Errorf("UploadUserGameResult3 err_no: %v,err: %v", request.(map[string]any)["errcode"], request.(map[string]any)["errmsg"])
	}
	return nil
}
