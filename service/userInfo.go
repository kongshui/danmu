package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// 存储用户信息
func userInfoStore(user UserInfoStruct, isAnchor bool) error {
	var dbName string
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
	if isAnchor {
		dbName = anchor_info_db
	} else {
		dbName = user_info_db
	}
	return rdb.HSet(dbName, user.OpenId, data)
}

// 获取用户信息
func UserInfoGet(openId string, isAnchor bool) (UserInfoStruct, error) {
	var dbName string
	if isAnchor {
		dbName = anchor_info_db
	} else {
		dbName = user_info_db
	}
	userStr, err := rdb.HGet(dbName, openId)
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
func userInfoCompare(openId, NickName, AvatarUrl string, isAnchor bool) bool {
	var dbName string
	if isAnchor {
		dbName = anchor_info_db
	} else {
		dbName = user_info_db
	}
	ok, err := rdb.HExists(dbName, openId)
	if err != nil || !ok {
		return false
	}
	user, err := UserInfoGet(openId, isAnchor)
	if err != nil {
		return false
	}
	if user.NickName == NickName && user.AvatarUrl == AvatarUrl && user.OpenId == openId {
		return true
	}
	return false
}

// 对比后存储用户信息
func UserInfoCompareStore(openId, NickName, AvatarUrl string, isAnchor bool) {
	if ok, _ := rdb.SetKeyNX(openId+"_info_monitor", "1", time.Duration(config.App.UserChangeTime)*time.Second); !ok {
		return
	}

	// 对比后存储用户信息
	if !userInfoCompare(openId, NickName, AvatarUrl, isAnchor) {
		if err := userInfoStore(UserInfoStruct{OpenId: openId, NickName: NickName, AvatarUrl: AvatarUrl}, isAnchor); err != nil {
			ziLog.Error(fmt.Sprintf("UserInfoCompareStore InsertPlayerBaseInfo err: %v,openId： %v", err, openId), debug)
		}
	}

}

// 主播信息
func anchorInfoGet(openId string) (UserInfoStruct, error) {
	if !rdb.IsExistKey(anchor_info_db) {
		return UserInfoStruct{}, errors.New("anchorInfoGet err: anchor_info_db not exist")
	}
	userStr, err := rdb.HGet(anchor_info_db, openId)
	if err != nil {
		return UserInfoStruct{}, errors.New("anchorInfoGet err: " + err.Error())
	}
	var user UserInfoStruct
	if err := json.Unmarshal([]byte(userStr), &user); err != nil {
		return UserInfoStruct{}, errors.New("anchorInfoGet err: " + err.Error())
	}
	return user, nil
}
