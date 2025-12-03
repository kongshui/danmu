package dao

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type (
	RedisClient struct {
		Client *redis.Client
		Lock   *sync.RWMutex
	}
)

func init() {
}

// RedisInit 初始化 Redis 客户端
func (rdb *RedisClient) RedisInit(addr string, password string, db int, isCluster bool) {
	// var (
	// 	addr     string = "127.0.0.1:6379"
	// 	password string = ""
	// 	// err      error
	// )
	// addr = os.Getenv("REDIS_ADDRESS")
	// password = os.Getenv("REDIS_PASSWORD")
	// if addr == "" {
	// 	log.Println("从环境中获取redis地址失败...")
	// 	os.Exit(6)
	// }
	// s, _ := os.Getwd()
	// if strings.Contains(s, "d:\\OneDrive\\wendang\\program\\go\\src\\bced") {
	// 	addr = "127.0.0.1:6379"
	// 	password = ""
	// }

	rdb.Lock = &sync.RWMutex{}
	rdb.Lock.Lock()
	defer rdb.Lock.Unlock()
	rdb.Client = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     10,
		MinIdleConns: 1,
		DialTimeout:  5 * time.Second,
	})
}

func (rdb *RedisClient) RedisCheckPing(addr string, password string, db int, isCluster bool) {
	t := time.NewTicker(5 * time.Second)
	for {
		<-t.C
		_, err := rdb.Ping()
		if err != nil {
			log.Println("redis连接失败,重新初始化")
			rdb.RedisInit(addr, password, db, isCluster)
			continue
		}
	}
}

// redis key type
func (rdb *RedisClient) Type(key string) (string, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	// 检查连接是否正常
	return rdb.Client.Type(key).Result()
}

// 移除key 过期时间
func (rdb *RedisClient) Persist(key string) error {
	rdb.Lock.Lock()
	defer rdb.Lock.Unlock()
	// 移除键值
	return rdb.Client.Persist(key).Err()
}

// 为一个key增加一个float值
func (rdb *RedisClient) IncrByFloat(key string, value float64) (float64, error) {
	rdb.Lock.Lock()
	defer rdb.Lock.Unlock()
	// 为一个key增加一个float值
	return rdb.Client.IncrByFloat(key, value).Result()
}

// 判断redis中是否存在key
func (rdb *RedisClient) IsExistKey(key string) bool {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	// 判断键值是否存在
	return rdb.Client.Exists(key).Val() == 1
}

// GetRedisClient 获取 Redis 客户端
func GetRedisClient() *RedisClient {
	return &RedisClient{
		Lock: &sync.RWMutex{},
	}
}

// setkey 设置键值对到 Redis
func (rdb *RedisClient) Set(key string, value any, expiration time.Duration) error {
	rdb.Lock.Lock()
	defer rdb.Lock.Unlock()
	// 设置键值对，并设置过期时间
	err := rdb.Client.Set(key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("设置键值失败: %w", err)
	}
	return nil
}

// Get 从 Redis 获取键值
func (rdb *RedisClient) Get(key string) (string, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	// 获取键值
	return rdb.Client.Get(key).Result()
}

// Del 从 Redis 删除键值
func (rdb *RedisClient) Del(key string) error {
	rdb.Lock.Lock()
	defer rdb.Lock.Unlock()
	// 删除键值
	return rdb.Client.Del(key).Err()
}

// SetKeyNX 设置键值对到 Redis，如果键值已存在则不设置
func (rdb *RedisClient) SetKeyNX(key string, value any, expiration time.Duration) (bool, error) {
	rdb.Lock.Lock()
	defer rdb.Lock.Unlock()
	// 设置键值对，并设置过期时间
	return rdb.Client.SetNX(key, value, expiration).Result()
}

// 发布订阅消息
func (rdb *RedisClient) Publish(channel string, message any) error {
	rdb.Lock.Lock()
	defer rdb.Lock.Unlock()
	return rdb.Client.Publish(channel, message).Err()
}

// 查询TTl过期时间
func (rdb *RedisClient) TTL(key string) (time.Duration, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	return rdb.Client.TTL(key).Result()
}

// 设置过期时间
func (rdb *RedisClient) Expire(key string, expiration time.Duration) error {
	rdb.Lock.Lock()
	defer rdb.Lock.Unlock()
	return rdb.Client.Expire(key, expiration).Err()
}

// 设置过期时间戳
func (rdb *RedisClient) ExpireAt(key string, tm time.Time) error {
	rdb.Lock.Lock()
	defer rdb.Lock.Unlock()
	return rdb.Client.ExpireAt(key, tm).Err()
}

// hash 设置键值对到 Redis
func (rdb *RedisClient) HSet(key string, field string, value any) error {
	rdb.Lock.Lock()
	defer rdb.Lock.Unlock()
	return rdb.Client.HSet(key, field, value).Err()
}

// hash，hsetnx 设置键值对到 Redis，如果键值已存在则不设置
func (rdb *RedisClient) HSetNX(key string, field string, value any) (bool, error) {
	rdb.Lock.Lock()
	defer rdb.Lock.Unlock()
	return rdb.Client.HSetNX(key, field, value).Result()
}

// hash 多个hash值设置到redis
func (rdb *RedisClient) HMSet(key string, value map[string]any) error {
	rdb.Lock.Lock()
	defer rdb.Lock.Unlock()
	return rdb.Client.HMSet(key, value).Err()
}

// hash 获取键值
func (rdb *RedisClient) HGet(key string, field string) (string, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	// 获取键值
	return rdb.Client.HGet(key, field).Result()
}

// hash 获取键值是否存在
func (rdb *RedisClient) HExists(key string, field string) (bool, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	// 获取键值
	return rdb.Client.HExists(key, field).Result()
}

// hash 增加一个float值
func (rdb *RedisClient) HIncrByFloat(key string, field string, value float64) (float64, error) {
	rdb.Lock.Lock()
	defer rdb.Lock.Unlock()
	// 为一个key增加一个float值
	return rdb.Client.HIncrByFloat(key, field, value).Result()
}

// HScan 扫描hash键值对
func (rdb *RedisClient) HScan(key string, cursor uint64, match string, count int64) (keys []string, nextCursor uint64, err error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	// 扫描hash键值对
	return rdb.Client.HScan(key, cursor, match, count).Result()
}

// hash 查询多个hash值
func (rdb *RedisClient) HMGet(key string, fields ...string) ([]any, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	// 获取键值
	return rdb.Client.HMGet(key, fields...).Result()
}

// hash 删除键值
func (rdb *RedisClient) HDel(key string, field string) error {
	rdb.Lock.Lock()
	defer rdb.Lock.Unlock()
	// 获取键值
	return rdb.Client.HDel(key, field).Err()
}

// hash 获取所有key
func (rdb *RedisClient) HKeys(key string) ([]string, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	// 获取键值
	return rdb.Client.HKeys(key).Result()
}

// hash 获取所有键值
func (rdb *RedisClient) HGetAll(key string) (map[string]string, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	// 获取键值
	return rdb.Client.HGetAll(key).Result()
}

// hash 获取hash长度
func (rdb *RedisClient) HLen(key string) (int64, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	// 获取键值
	return rdb.Client.HLen(key).Result()
}

// ping redis
func (rdb *RedisClient) Ping() (string, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	return rdb.Client.Ping().Result()
}

// 添加一个值到集合中
func (rdb *RedisClient) SAdd(key string, value any) (int64, error) {
	rdb.Lock.Lock()
	defer rdb.Lock.Unlock()
	return rdb.Client.SAdd(key, value).Result()
}

// 获取集合中的所有值
func (rdb *RedisClient) SMembers(key string) ([]string, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	return rdb.Client.SMembers(key).Result()
}

// 删除集合中的某个值
func (rdb *RedisClient) SRem(key string, value any) error {
	rdb.Lock.Lock()
	defer rdb.Lock.Unlock()
	return rdb.Client.SRem(key, value).Err()
}

// 判断redis集合中中是否有某个key
func (rdb *RedisClient) SIsMember(key string, value any) (bool, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	return rdb.Client.SIsMember(key, value).Result()
}

// 添加有序集合
func (rdb *RedisClient) ZAdd(key string, score float64, member string) error {
	rdb.Lock.Lock()
	defer rdb.Lock.Unlock()
	return rdb.Client.ZAdd(key, redis.Z{Score: score, Member: member}).Err()
}

// 获取有序集合某个成员
func (rdb *RedisClient) ZScore(key string, member string) (float64, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	return rdb.Client.ZScore(key, member).Result()
}

// 从小到大获取有序集合
func (rdb *RedisClient) ZRange(key string, start, stop int64) ([]string, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	return rdb.Client.ZRange(key, start, stop).Result()
}

// 从小到大获取有序集合withscores
func (rdb *RedisClient) ZRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	return rdb.Client.ZRangeWithScores(key, start, stop).Result()
}

// 从大到小获取有序集合
func (rdb *RedisClient) ZRevRange(key string, start, stop int64) ([]string, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	return rdb.Client.ZRevRange(key, start, stop).Result()
}

// 从大到小获取有序集合withscores
func (rdb *RedisClient) ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	return rdb.Client.ZRevRangeWithScores(key, start, stop).Result()
}

// 为有序集合成员增加increment
func (rdb *RedisClient) ZIncrBy(key string, increment float64, member string) (float64, error) {
	rdb.Lock.Lock()
	defer rdb.Lock.Unlock()
	return rdb.Client.ZIncrBy(key, increment, member).Result()
}

// 移除有序集合中的成员
func (rdb *RedisClient) ZRem(key string, members ...any) error {
	rdb.Lock.Lock()
	defer rdb.Lock.Unlock()
	return rdb.Client.ZRem(key, members...).Err()
}

// 合并有序集合
func (rdb *RedisClient) ZUnionStore(dest string, store redis.ZStore, keys []string) error {
	rdb.Lock.Lock()
	defer rdb.Lock.Unlock()
	return rdb.Client.ZUnionStore(dest, store, keys...).Err()
}

// 有序集合长度
func (rdb *RedisClient) ZCard(key string) (int64, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	return rdb.Client.ZCard(key).Result()
}

// 获取member在有序集合中的位置
func (rdb *RedisClient) ZRank(key string, member string) (int64, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	return rdb.Client.ZRank(key, member).Result()
}

// 获取member在有序集合中的位置
func (rdb *RedisClient) ZRevRank(key string, member string) (int64, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	return rdb.Client.ZRevRank(key, member).Result()
}

// 插入数据到列表末尾
func (rdb *RedisClient) RPush(key string, value any) error {
	rdb.Lock.Lock()
	defer rdb.Lock.Unlock()
	return rdb.Client.RPush(key, value).Err()
}

// 获取列表长度
func (rdb *RedisClient) LLen(key string) (int64, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	return rdb.Client.LLen(key).Result()
}

// 获取列表中的值
func (rdb *RedisClient) LRange(key string, start, stop int64) ([]string, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	return rdb.Client.LRange(key, start, stop).Result()
}

// 通过下表获取列表中的值
func (rdb *RedisClient) LIndex(key string, index int64) (string, error) {
	rdb.Lock.RLock()
	defer rdb.Lock.RUnlock()
	return rdb.Client.LIndex(key, index).Result()
}

// 重命名key
func (rdb *RedisClient) Rename(key string, newKey string) error {
	rdb.Lock.Lock()
	defer rdb.Lock.Unlock()
	return rdb.Client.Rename(key, newKey).Err()
}
