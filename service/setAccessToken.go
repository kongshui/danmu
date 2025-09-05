package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/kongshui/danmu/common"
)

func GetAccessToken() string {
	if accessToken.Token == "" {
		token, err := rdb.Get("access_token")
		if err != nil {
			return ""
		}
		//设置token
		accessToken.Lock.Lock()
		accessToken.Token = token
		accessToken.Lock.Unlock()
	}
	return accessToken.Token
}

func setKsGlobalAccessToken() error {
	//创建请求头
	var header map[string]string = map[string]string{
		"Content-Type": "x-www-form-urlencoded",
	}
	//创建body
	uri := map[string]string{
		"app_id":     app_id,
		"app_secret": app_secret,
		"grant_type": "client_credentials",
	}
	pathUrl, err := common.GetUrl(url_GetAccessTokenUrl, uri)
	if err != nil {
		ziLog.Error(fmt.Sprintf("SetGlobalAccessToken SetAccessToken Url err: %v", err), debug)
	}
	// fmt.Println(pathUrl, 2222222222222)
	response, err := common.HttpRespond("GET", pathUrl, nil, header)
	if err != nil {
		return fmt.Errorf("SetGlobalAccessToken response err: %v, path: %v", err, pathUrl)
	}
	// 关闭响应体以释放资源
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("SetGlobalAccessToken status err: %v", response.Status)
	}
	// 读取响应体并解析为JSON对象
	var (
		result AccessKsTokenRespondStruct
	)

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("SetGlobalAccessToken json.NewDecoder err: %v", err)
	}
	if result.Result != 1 {
		return errors.New(strconv.FormatInt(result.Result, 10))
	}
	// fmt.Println(result, 66666666)
	//设置access_token
	accessToken.Lock.Lock()
	accessToken.Token = result.AccessToken
	accessToken.Lock.Unlock()
	err = rdb.Set(access_token_db, result.AccessToken, time.Duration(result.ExpiresIn)*time.Second)
	if err != nil {
		return fmt.Errorf("SetGlobalAccessToken  rdb set err: %v", err)
	}
	return nil
}

// 设置抖音登录
func setDyGlobalAccessToken() error {
	//创建请求头
	var header map[string]string = map[string]string{
		"Content-Type": "application/json",
	}
	//创建body
	body := map[string]string{
		"appid":      app_id,
		"secret":     app_secret,
		"grant_type": "client_credential",
	}
	bodyByte, _ := json.Marshal(body)
	response, err := common.HttpRespond("POST", url_GetAccessTokenUrl, bodyByte, header)
	if err != nil {
		return fmt.Errorf("SetGlobalAccessToken response err: %v", err)
	}
	// 关闭响应体以释放资源
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("SetGlobalAccessToken status err: %v", response.Status)
	}
	// 读取响应体并解析为JSON对象
	var (
		result AccessDyTokenRespondStruct
	)

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("SetGlobalAccessToken json.NewDecoder err: %v", err)
	}
	if result.ErrorNo != 0 {
		return errors.New(strconv.FormatInt(result.ErrorNo, 10))
	}
	// fmt.Println(result, 66666666)
	//设置access_token
	accessToken.Lock.Lock()
	accessToken.Token = result.Data.AccessToken
	accessToken.Lock.Unlock()
	err = rdb.Set(access_token_db, result.Data.AccessToken, time.Duration(result.Data.ExpiresIn)*time.Second)
	if err != nil {
		return fmt.Errorf("SetGlobalAccessToken  rdb set err: %v", err)
	}
	return nil
}

func setToken() error {
	var (
		isSet     bool
		timeCheck time.Duration
		function  func() error
	)
	switch platform {
	case "ks":
		timeCheck = 24 * time.Hour
		function = setKsGlobalAccessToken
	case "dy":
		timeCheck = 50 * time.Minute
		function = setDyGlobalAccessToken
	}
	if rdb.IsExistKey(access_token_db) {
		timeLeave, _ := rdb.TTL(access_token_db)
		if timeLeave < timeCheck {
			isSet = true
		}
	} else {
		isSet = true
	}
	if isSet {
		ok, err := rdb.SetKeyNX(monitor_access_token_db, nodeUuid, timeCheck)
		if err != nil {
			ziLog.Error(fmt.Sprintf("setAccessToken 设置全局Access token推送标识失败: %v", err), debug)
			return err
		}
		if ok {
			count := 0
			for {
				if count >= 5 {
					return errors.New("setAccessToken 失败")
				}
				if err := function(); err != nil {
					count++
					time.Sleep(1 * time.Second)
					fmt.Println("setAccessToken: ", err)
				} else {
					break
				}
			}
		}

	}
	return nil
}
func setAccessToken() {
	var (
		t         *time.Ticker
		timeCheck time.Duration
		function  func() error
	)
	switch platform {
	case "ks":
		timeCheck = 24 * time.Hour
		function = setKsGlobalAccessToken
	case "dy":
		timeCheck = 50 * time.Minute
		function = setDyGlobalAccessToken
	}
	// if rdb.IsExistKey(access_token_db) {
	// 	timeLeave, _ := rdb.TTL(access_token_db)
	// 	if timeLeave < timeCheck {
	// 		isSet = true
	// 	}
	// } else {
	// 	isSet = true
	// }
	// if isSet {
	// 	for {
	// 		if err := function(); err != nil {
	// 			time.Sleep(1 * time.Second)
	// 			fmt.Println("setAccessToken: ", err)
	// 		} else {
	// 			break
	// 		}
	// 	}
	// }
	t = time.NewTicker(55 * time.Minute)
	for {
		<-t.C
		count := 0
		ok, err := rdb.SetKeyNX(monitor_access_token_db, nodeUuid, timeCheck)
		if err != nil {
			ziLog.Error(fmt.Sprintf("setAccessToken 设置全局Access token推送标识失败: %v", err), debug)
			continue
		}
		if ok {
			for {
				if count >= 13800 {
					count = 0
					break
				}
				if err := function(); err != nil {
					time.Sleep(6 * time.Second)
					ziLog.Error(fmt.Sprintf("setAccessToken 设置全局Access token失败: %v", err), debug)
					count++
				} else {
					break
				}
			}
		}
	}

}

// 获取token
func getAccessToken() {
	t := time.NewTicker(time.Minute * 5)
	for {
		<-t.C
		getToken()
	}
}

func getToken() error {
	token, err := rdb.Get(access_token_db)
	if err != nil {
		ziLog.Error(fmt.Sprintf("getToken 获取全局Access token失败: %v", err), debug)
		return fmt.Errorf("getToken 获取全局Access token失败: %v", err)
	}
	accessToken.Lock.Lock()
	accessToken.Token = token
	accessToken.Lock.Unlock()
	return nil
}
