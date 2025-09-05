package service

import (
	"fmt"
)

// 玩家获胜统计，获胜次数。about  true 左边，false，右边
func OpenIdWinCount(openId string, groupId string) error {
	if err := mysql.UpdateOpenWinCount(openId, groupId); err != nil {
		return fmt.Errorf("OpenIdWinCount err: %v", err)
	}
	return nil
}

// 组获胜统计，获胜次数。about  true 左边，false，右边
func GroupWinCount(groupId string) error {
	if err := mysql.UpdateGroupWinCount(groupId); err != nil {
		return fmt.Errorf("GroupWinCount err: %v", err)
	}
	return nil
}
