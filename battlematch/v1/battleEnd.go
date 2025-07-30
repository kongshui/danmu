package battlematch

import (
	"context"
	"path"
)

// 匹配组结束
func MatchGroupEnd(ctx context.Context, groupId string) bool {
	// 删除匹配组状态
	return etcdClient.PutIfNotExist(ctx, path.Join("/", projectName, match_battle_group_end, groupId), groupId, etcdClient.NewLease(ctx, 1800))
}
