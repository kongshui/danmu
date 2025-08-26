package service

import (
	"fmt"
)

// 玩家获胜统计，获胜次数。about  true 左边，false，右边
func OpenIdWinCount(openId string, about bool) error {
	if err := mysql.UpdateOpenWinCount(openId, about); err != nil {
		ziLog.Error(fmt.Sprintf("OpenIdWinCount err: %v", err), debug)
		return err
	}
	return nil
}

// 组获胜统计，获胜次数。about  true 左边，false，右边
func GroupWinCount(groupId string) error {
	if err := mysql.UpdateGroupWinCount(groupId); err != nil {
		ziLog.Error(fmt.Sprintf("GroupWinCount err: %v", err), debug)
		return err
	}
	return nil
}
