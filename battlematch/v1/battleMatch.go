package battlematch

import (
	"context"

	dao_etcd "github.com/kongshui/danmu/dao/etcd"
)

var (
	etcdClient  = dao_etcd.NewEtcd()
	first_ctx   context.Context
	projectName string
)

func init() {

}

// etcd 初始化
func InitEtcd(client *dao_etcd.Etcd) {
	etcdClient = client
}

// 初始化项目名称
func InitProjectName(project string) {
	projectName = project
}
