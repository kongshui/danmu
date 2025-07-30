package service

import (
	"path"

	"github.com/kongshui/danmu/common"
)

// 通过openId获取uid
func queryUidByOpenid(openid string) string {
	res, err := etcdClient.Client.Get(first_ctx, path.Join("/", config.Project, common.OpenId_Register_Uid_key, openid))
	if err != nil {
		// ziLog.Write(logError, fmt.Sprintf("PlayerChooseGroupHandle 获取uid失败， err: %v", err), debug)
		return ""
	}
	if res.Count == 0 { // 没有用户掉线注册
		// ziLog.Write(logError, "PlayerChooseGroupHandle 获取uid失败， err: is nil", debug)
		return ""
	}
	return string(res.Kvs[0].Value)
}

// 通过uid获取openId
// func queryOpenidByUidStr(uidStr string) string {
// 	res, err := etcdClient.Client.Get(first_ctx, common.OpenId_Register_key+uidStr)
// 	if err != nil {
// 		ziLog.Write(logError, fmt.Sprintf("FindOpenidFromUidStr 获取uid失败， err: %v", err), debug)
// 		return ""
// 	}
// 	openid := res.Kvs[0].Value
// 	return string(openid)
// }
