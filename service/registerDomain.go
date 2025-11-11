package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"errors"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 注册后端域名
func registerBackDomain(ctx context.Context) {
	time.Sleep(6 * time.Second)
	domain := "http://" + cfg.Server.Addr + ":" + cfg.Server.Port
	// 租约
	listenId := etcdClient.NewLease(ctx, 3)
	if listenId == 0 {
		ziLog.Error("registerBackDomain err: 创建租约失败", debug)
		os.Exit(10)
	}
	go func(doma string, ctx context.Context, id clientv3.LeaseID) {
		t := time.NewTicker(1 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				_, err := etcdClient.Client.Put(ctx, path.Join("/", cfg.Project, backend_domain_key, cfg.Server.Addr+":"+cfg.Server.Port), domain, clientv3.WithLease(id))
				if err != nil {
					log.Println("发送消息至etcd失败", err)
					continue
				}
				return
			case <-ctx.Done():
				return
			}
		}
	}(domain, ctx, listenId)

	etcdClient.KeepLease(ctx, listenId)
}

// Get  foward domain
func GetFowardDomain(ctx context.Context) {
	oneGetFowardDomain()
	respond := etcdClient.Client.Watch(ctx, path.Join("/", cfg.Project, forward_domain_key), clientv3.WithPrefix(), clientv3.WithPrevKV())
	for wresp := range respond {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				// 处理新增事件
				forward_domain.Add(string(ev.Kv.Value))
			case clientv3.EventTypeDelete:
				// 处理删除事件
				ziLog.Info(fmt.Sprintf("getFowardDomain 删除前端http节点信息, key: %v, value: %v", string(ev.Kv.Key), string(ev.Kv.Value)), debug)
				key := path.Base(string(ev.Kv.Key))
				forward_domain.Remove("http://" + key)
				oneGetFowardDomain()
				continue
			}
		}
	}
}

// 获取前端服务器
func oneGetFowardDomain() error {
	respond, err := etcdClient.Client.Get(first_ctx, path.Join("/", cfg.Project, forward_domain_key), clientv3.WithPrefix())
	if err != nil {
		log.Println("获取前端服务器失败", err)
		return errors.New("获取前端服务器失败: " + err.Error())
	}
	if len(respond.Kvs) == 0 {
		return errors.New("获取前端服务器失败: 前端服务器数量为零")
	}
	for _, kv := range respond.Kvs {
		forward_domain.Add(string(kv.Value))
	}
	return nil
}

// 获取grpc前端服务器
// 获取前端服务器
// func oneGetGrpcDomain() error {
// 	respond, err := etcdClient.Client.Get(first_ctx, path.Join("/", config.Project, grpc_domain_key), clientv3.WithPrefix())
// 	if err != nil {
// 		log.Println("获取前端服务器失败", err)
// 		return errors.New("获取前端服务器失败: " + err.Error())
// 	}
// 	if len(respond.Kvs) == 0 {
// 		return errors.New("获取前端服务器失败: 前端服务器数量为零")
// 	}
// 	for _, kv := range respond.Kvs {
// 		log.Println("获取到的grpc服务器", string(kv.Value))
// 		grpcConn(string(kv.Value))
// 	}
// 	if grpc_pool.Len() < 5 {
// 		oneGetGrpcDomain()
// 	}
// 	log.Println("获取到的grpc服务器数量", grpc_pool.Len())
// 	return nil
// }

// get grpc
// func getGrpcDomain(ctx context.Context) {
// 	oneGetGrpcDomain()
// 	respond := etcdClient.Client.Watch(ctx, path.Join("/", config.Project, grpc_domain_key), clientv3.WithPrefix(), clientv3.WithPrevKV())
// 	for wresp := range respond {
// 		for _, ev := range wresp.Events {
// 			switch ev.Type {
// 			case clientv3.EventTypePut:
// 				// 处理新增事件
// 				// grpc_domain.Add(string(ev.Kv.Value))
// 				grpcConn(string(ev.Kv.Value))
// 			case clientv3.EventTypeDelete:
// 				// 处理删除事件
// 				ziLog.Info(fmt.Sprintf("getFowardDomain 删除前端http节点信息, key: %v, value: %v", string(ev.Kv.Key), string(ev.Kv.Value)), debug)
// 				// key := path.Base(string(ev.Kv.Key))
// 				// grpc_domain.Remove(key)
// 				continue
// 			}
// 		}
// 	}
// }

// grpcConnTest
// func TestGrpcConn() {
// 	t := time.NewTicker(100 * time.Millisecond)
// 	defer t.Stop()
// 	groupGrpc := &pb.AddGiftToGroupReq{}
// 	groupGrpc.AnchorOpenId = "1"
// 	groupGrpc.AnchorOpenIdList = []string{"1", "2"}
// 	groupGrpc.GiftId = "0"
// 	groupGrpc.GroupId = "10000"
// 	groupGrpc.IsComment = true
// 	groupGrpc.OpenId = "666666"
// 	groupGrpc.IsJoin = true
// 	if err := grpcSend(groupGrpc, 0); err != nil {
// 		log.Println("TestGrpcConn 发送消息至etcd失败111111", err)
// 	}
// 	groupGrpc.IsJoin = false
// 	var count int64
// 	for {
// 		<-t.C
// 		count++
// 		groupGrpc.GiftNum = count
// 		if err := grpcSend(groupGrpc, 0); err != nil {
// 			log.Println("TestGrpcConn 发送消息至etcd失败", err)
// 			continue
// 		}
// 		if count > 200 {
// 			break
// 		}
// 	}
// }
