package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/kongshui/danmu/model/pmsg"
	"github.com/kongshui/danmu/sse"
	"google.golang.org/protobuf/proto"
)

// 玩家分组
func playerGroupAdd(roomId, uidStr string, roundId int64, userMap []*pmsg.SingleRoomAddGroupInfo, isChoose bool) error {

	//初始化sdata
	data := &pmsg.ResultUserAddGroupMessage{}
	//设置查询组的名称
	name := roomId + "_" + strconv.FormatInt(roundId, 10) + "_group"
	for _, v := range userMap {
		if _, err := rdb.HSetNX(name, v.GetOpenId(), v.GetGroupId()); err != nil {
			ziLog.Error(fmt.Sprintf("playerGroupAdd 设置组失败: %v,openId:%v, groupId: %v", err, v.GetOpenId(), v.GetGroupId()), debug)
		}

		// 其他前置处理
		if playerGroupAddin != nil {
			if err := playerGroupAddin(roomId, v.GetOpenId()); err != nil {
				ziLog.Error(fmt.Sprintf("playerGroupAdd playerGroupAddinFunc失败, err: %v,openId:%v, groupId: %v, roomId:%v", err, v.GetOpenId(), v.GetGroupId(), roomId), debug)

			}
		}
		go userInfoCompareStore(v.GetOpenId(), v.GetNickName(), v.GetAvatarUrl(), false)
		//
		go dyUploadUserGroup(roomId, v.GetOpenId(), v.GetGroupId(), roundId)
		// 是否是通过小摇杆加入
		if isChoose {
			continue
		}
		//获取玩家连胜币
		coin, _ := QueryUserWinStreamCoin(v.GetOpenId())
		// 查询玩家是否已经消费
		isConsume := queryIsConsume(v.OpenId)
		// 查询玩家等级
		level, _ := QueryLevelInfo(v.OpenId)

		score, rank, _ := getPlayerWorldRankData(v.OpenId)

		winningPoint, _ := QueryUserWinningPoint(v.OpenId)
		data.UserInfoList = append(data.UserInfoList, &pmsg.UserInfoStruct{
			OpenId:            v.OpenId,
			VersionScore:      score,
			VersionRank:       rank,
			WinningStreamCoin: coin,
			IsFirstConsume:    isConsume,
			Level:             level,
			WinningPoints:     winningPoint,
		})
	}
	ttl, _ := rdb.TTL(name)
	if ttl <= 0 {
		rdb.Expire(name, 21600*time.Second)
	}
	if isChoose {
		return nil
	}
	// fmt.Println(data)
	dataByte, err := proto.Marshal(data)
	if err != nil {
		return errors.New("playerGroupAdd proto Marshal err: " + err.Error())
	}
	if err := sse.SseSend(pmsg.MessageId_SingleRoomAddGroupAck, []string{uidStr}, dataByte); err != nil {
		ziLog.Error(fmt.Sprintf("playerGroupAdd 玩家加入组信息 err: %v", err), debug)
		return errors.New("玩家加入组信息 err: " + err.Error())
	}
	return nil
}
