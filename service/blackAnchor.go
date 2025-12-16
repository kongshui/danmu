package service

import (
	"encoding/json"

	"github.com/kongshui/danmu/model/pmsg"
)

// blackAnchorListAdd 添加黑名单主播
func blackAnchorListAdd(anchorOpenid string) error {
	_, err := rdb.SAdd(black_anchor_list_db, anchorOpenid)
	if err != nil {
		return err
	}
	return nil
}

// blackAnchorListDel 删除黑名单主播
func blackAnchorListDel(anchorOpenid string) error {
	if err := rdb.SRem(black_anchor_list_db, anchorOpenid); err != nil {
		return err
	}
	return nil
}

// 返回主播黑名单所有成员
func blackAnchorListMembers() ([]string, error) {
	members, err := rdb.SMembers(black_anchor_list_db)
	if err != nil {
		return nil, err
	}
	return members, nil
}

// blackAnchorListClear 清空黑名单主播
func blackAnchorListClear() error {
	if err := rdb.Del(black_anchor_list_db); err != nil {
		return err
	}
	return nil
}

// blackAnchorListIsMember 判断主播是否在黑名单
func blackAnchorListIsMember(anchorOpenid string) bool {
	isMember, _ := rdb.SIsMember(black_anchor_list_db, anchorOpenid)
	if !isMember {
		return false
	}
	data := &pmsg.BlackAnchorMessage{
		AnchorOpenId: anchorOpenid,
		Msg:          "您已被管理员移出房间，如有疑问请联系客服",
	}
	dataBytes, _ := json.Marshal(data)
	// 发送断线消息
	if err := SendMessage(pmsg.MessageId_BlackAnchorLogOff, []string{anchorOpenid}, dataBytes); err != nil {
		// 发送消息失败
		ziLog.Error("黑名单主播发送断线消息失败："+err.Error(), debug)
		return false
	}
	return true
}
