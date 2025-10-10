package service

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

func blackAnchorListClear() error {
	if err := rdb.Del(black_anchor_list_db); err != nil {
		return err
	}
	return nil
}
