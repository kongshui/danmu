package battlematch

import (
	"context"
	"path"
)

// 设置连胜币状态
func MatchAddStreamCoinStatusSet(ctx context.Context, groupId, label string) bool {
	switch label {
	case "divideup":
		return etcdClient.PutIfNotExist(ctx, path.Join("/", projectName, match_battle_divideup_coin_status, groupId), groupId, etcdClient.NewLease(ctx, 180))

	case "normal":
		return etcdClient.PutIfNotExist(ctx, path.Join("/", projectName, match_battle_normal_coin_status, groupId), groupId, etcdClient.NewLease(ctx, 180))
	default:
		return false
	}
}
