package dao

import (
	"context"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Etcd struct {
	Client *clientv3.Client
}

func init() {

}

// new
func NewEtcd() *Etcd {
	return &Etcd{}
}

// 初始化
func (etcd *Etcd) InitEtcd(endpoints []string, username, password string) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints, // 你的etcd服务器地址
		DialTimeout: 10 * time.Second,
		Username:    username,
		Password:    password,
	})
	if err != nil {
		log.Println("etcd初始化失败", err)
	}
	etcd.Client = cli
	go etcd.isAlived(endpoints, username, password)
}

// 检测是否存活
func (etcd *Etcd) isAlived(endpoints []string, username, password string) {
	t := time.NewTicker(3 * time.Second)
	for {
		<-t.C
		count := 0
		for _, v := range endpoints {
			_, err := etcd.Client.Status(context.Background(), v)
			if err != nil {
				log.Println("获取状态失败，节点： ", v)
				count++
			}
			if count == len(endpoints) {
				// 重新初始胡
				etcd.InitEtcd(endpoints, username, password)
			}
		}
	}
}

// 创建租约lease，设置ttl,秒级，默认可以设置为3600s
// func (etcd *Etcd) NewLease(ctx context.Context, ttl int64) clientv3.LeaseID {
// 	leaseGrantResp, err := etcd.Client.Grant(ctx, ttl)
// 	if err != nil {
// 		log.Println("创建租约失败", err)
// 	}
// 	return leaseGrantResp.ID
// }

// // 续约
// func (etcd *Etcd) KeepAlive(ctx context.Context, leaseID clientv3.LeaseID) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
// 	keepRespChan, err := etcd.Client.KeepAlive(ctx, leaseID)
// 	if err != nil {
// 		log.Println("续约失败", err)
// 	}
// 	return keepRespChan, err
// }

// // 释放租约
//
//	func (etcd *Etcd) Revoke(ctx context.Context, leaseID clientv3.LeaseID) {
//		_, err := etcd.Client.Revoke(ctx, leaseID)
//		if err != nil {
//			log.Println("释放租约失败", err)
//		}
//	}
//
// newLease
func (etcd *Etcd) NewLease(ctx context.Context, ttl int64) clientv3.LeaseID {
	leaseGrantResp, err := etcd.Client.Grant(ctx, ttl)
	if err != nil {
		log.Println("创建租约失败", err)
	}
	return leaseGrantResp.ID
}

// 续约
func (etcd *Etcd) KeepAlive(ctx context.Context, leaseID clientv3.LeaseID) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	keepRespChan, err := etcd.Client.KeepAlive(ctx, leaseID)
	if err != nil {
		log.Println("续约失败", err)
	}
	return keepRespChan, err
}

// 释放租约
func (etcd *Etcd) Revoke(ctx context.Context, leaseID clientv3.LeaseID) {
	_, err := etcd.Client.Revoke(ctx, leaseID)
	if err != nil {
		log.Println("释放租约失败", err)
	}
}

// 保持租约
func (etcd *Etcd) KeepLease(ctx context.Context, leaseId clientv3.LeaseID) {
	ch, err := etcd.KeepAlive(ctx, leaseId)
	if err != nil {
		log.Println("续约失败", err)
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case keepResp := <-ch:
			if keepResp == nil {
				return
			}
		}
	}
}

// 不存在则创建
func (etcd *Etcd) PutIfNotExist(ctx context.Context, key string, value string, lease clientv3.LeaseID) bool {
	if lease == 0 {
		a, _ := etcd.Client.Txn(ctx).If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).Then(clientv3.OpPut(key, value)).Commit()
		return a.Succeeded
	}
	a, _ := etcd.Client.Txn(ctx).If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).Then(clientv3.OpPut(key, value, clientv3.WithLease(lease))).Commit()
	return a.Succeeded
}

// 根据key获取租约
func (etcd *Etcd) GetLeaseByKey(ctx context.Context, key string) (clientv3.LeaseID, error) {
	resp, err := etcd.Client.Get(ctx, key)
	if err != nil {
		log.Println("获取租约失败", err)
		return 0, err
	}
	if len(resp.Kvs) == 0 {
		return 0, nil // key不存在
	}
	return clientv3.LeaseID(resp.Kvs[0].Lease), nil
}
