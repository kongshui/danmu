package service

import (
	"encoding/json"
	"errors"
	"fmt"
)

// 存储用户信息
func userInfoStore(user UserInfoStruct) error {
	data, err := json.Marshal(user)
	if err != nil {
		return errors.New("UserInfoStore err: " + err.Error())
	}
	ok, err := mysql.IsPlayerExist(user.OpenId)
	if err != nil {
		ziLog.Error(fmt.Sprintf("UserInfoStore err: %v,openId： %v", err, user.OpenId), debug)
	}
	if !ok {
		if err := mysql.InsertPlayerBaseInfo(user.OpenId, user.AvatarUrl, user.NickName); err != nil {
			ziLog.Error(fmt.Sprintf("UserInfoStore InsertPlayerBaseInfo err: %v,openId： %v", err, user.OpenId), debug)
		}
	} else {
		if err := mysql.UpdatePlayerBaseInfo(user.OpenId, user.AvatarUrl, user.NickName); err != nil {
			ziLog.Error(fmt.Sprintf("userInfoStore UpdatePlayerBaseInfo err: %v,openId： %v", err, user.OpenId), debug)
		}
	}
	return rdb.HSet(user_info_db, user.OpenId, data)
}

// 获取用户信息
func userInfoGet(openId string) (UserInfoStruct, error) {
	userStr, err := rdb.HGet(user_info_db, openId)
	if err != nil {
		return UserInfoStruct{}, errors.New("UserInfoGet err: " + err.Error())
	}
	var user UserInfoStruct
	if err := json.Unmarshal([]byte(userStr), &user); err != nil {
		return UserInfoStruct{}, errors.New("UserInfoGet err: " + err.Error())
	}
	return user, nil
}

// 对比用户信息
func userInfoCompare(openId, NickName, AvatarUrl string) bool {
	ok, err := rdb.HExists(user_info_db, openId)
	if err != nil || !ok {
		return false
	}
	user, err := userInfoGet(openId)
	if err != nil {
		return false
	}
	if user.NickName == NickName && user.AvatarUrl == AvatarUrl && user.OpenId == openId {
		return true
	}
	return false
}

// 对比后存储用户信息
func userInfoCompareStore(openId, NickName, AvatarUrl string) {
	// 检测数据库中是否存在玩家，不存在就插入
	// ok, err := mysql.IsPlayerExist(openId)
	// if err != nil {
	// 	ziLog.Error( fmt.Sprintf("UserInfoCompareStore err:: %v,openId： %v", err, openId), debug)
	// }
	// if !ok {
	// 	if err := mysql.InsertPlayerBaseInfo(openId, AvatarUrl, NickName); err != nil {
	// 		ziLog.Error( fmt.Sprintf("UserInfoCompareStore InsertPlayerBaseInfo err: %v,openId： %v", err, openId), debug)
	// 	}
	// }
	// 对比后存储用户信息
	if !userInfoCompare(openId, NickName, AvatarUrl) {
		if err := userInfoStore(UserInfoStruct{OpenId: openId, NickName: NickName, AvatarUrl: AvatarUrl}); err != nil {
			ziLog.Error(fmt.Sprintf("UserInfoCompareStore InsertPlayerBaseInfo err: %v,openId： %v", err, openId), debug)
		}
	}

}
