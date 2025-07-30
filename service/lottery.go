package service

// import (
// 	"fmt"
// 	"math/rand/v2"
// 	"path"

// 	clientv3 "go.etcd.io/etcd/client/v3"
// )

// func lottery(anchorOpenId, openid string, count int64) map[string]int64 {
// 	giftMap := make(map[string]int64)
// 	isFirst := true
// 	pcg, err := getPcg(anchorOpenId, openid)
// 	if err != nil {
// 		ziLog.Write(logError, fmt.Sprintf("获取抽奖节点失败: %v", err), debug)
// 		pcg = newPcg() // 如果获取失败，则生成新的随机数生成器
// 	} else {
// 		isFirst = false
// 	}

// 	r := rand.New(&pcg) // 使用获取到的或新生成的随机数生成器
// 	for range count {
// 		lottery := r.Uint64N(100)
// 		switch {
// 		case lottery < 5:
// 			giftMap["11582"] += 10
// 		case lottery >= 5 && lottery < 49:
// 			giftMap["12252"] += 1
// 		case lottery >= 49 && lottery < 79:
// 			giftMap["11606"] += 1
// 		case lottery >= 79 && lottery < 89:
// 			giftMap["11585"] += 1
// 		case lottery >= 89 && lottery < 94:
// 			giftMap["11586"] += 1
// 		case lottery >= 94 && lottery < 98:
// 			giftMap["11587"] += 1
// 		case lottery >= 98:
// 			giftMap["12720"] += 1
// 		default:
// 			giftMap["11582"] += 10
// 		}
// 	}
// 	if count >= 50 {
// 		all, ok := giftMap["12720"]
// 		if !ok || all < count/50 {
// 			jJCount := count/50 - all
// 			giftMap["12720"] = count / 50
// 			if giftMap["11582"] > jJCount*10 {
// 				giftMap["11582"] -= jJCount * 10
// 			} else if giftMap["11582"] < jJCount*10 {
// 				rCount := giftMap["11582"] / 10
// 				gCount := jJCount - rCount
// 				if giftMap["12252"] > gCount {
// 					giftMap["12252"] -= gCount
// 				}
// 				delete(giftMap, "11582")
// 			} else if giftMap["11582"] == jJCount*10 {
// 				delete(giftMap, "11582")
// 			}

// 		}
// 	}
// 	// 存储抽奖节点Pcg
// 	if err := storePcg(anchorOpenId, openid, pcg, isFirst); err != nil {
// 		ziLog.Write(logError, fmt.Sprintf("存储抽奖节点失败: %v", err), debug)
// 	}
// 	return giftMap
// }

// // 存储抽奖节点Pcg
// func storePcg(anchorOpenId, openid string, pcg rand.PCG, isFirst bool) error {
// 	pcgBytes, err := pcg.MarshalBinary()
// 	if err != nil {
// 		return fmt.Errorf("序列化错误: %w", err)
// 	}
// 	var id clientv3.LeaseID
// 	key := path.Join("/", config.Project, lottery_seed_key, anchorOpenId, openid)
// 	if isFirst {
// 		id = etcdClient.NewLease(first_ctx, 7200)
// 	} else {
// 		id, _ = etcdClient.GetLeaseByKey(first_ctx, key)
// 	}
// 	if _, err := etcdClient.Client.Put(first_ctx, key, string(pcgBytes), clientv3.WithLease(id)); err != nil {
// 		return fmt.Errorf("存储抽奖节点失败: %w", err)
// 	}
// 	return nil
// }

// // 获取抽奖节点Pcg
// func getPcg(anchorOpenId, openid string) (rand.PCG, error) {
// 	var pcg rand.PCG
// 	res, err := etcdClient.Client.Get(first_ctx, path.Join("/", config.Project, lottery_seed_key, anchorOpenId, openid))
// 	if err != nil {
// 		return pcg, fmt.Errorf("获取抽奖节点失败: %w", err)
// 	}
// 	if len(res.Kvs) == 0 {
// 		return pcg, fmt.Errorf("抽奖节点不存在")
// 	}
// 	if err := pcg.UnmarshalBinary([]byte(res.Kvs[0].Value)); err != nil {
// 		return pcg, fmt.Errorf("反序列化错误: %w", err)
// 	}
// 	return pcg, nil
// }

// // 生成随机数生成器
// func newPcg() rand.PCG {
// 	seed1 := rand.Uint64N(1 << 63) // 生成一个64位的随机数作为种子
// 	seed2 := rand.Uint64N(1 << 63) // 生成另一个64位的随机数作为种子
// 	randSource := rand.NewPCG(seed1, seed2)
// 	return *randSource
// }
