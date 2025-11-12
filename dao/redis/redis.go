package dao

import (
	"time"

	sRedis "github.com/kongshui/danmu/dao/redis/redis"
	cRedis "github.com/kongshui/danmu/dao/redis/rediscluster"

	"github.com/go-redis/redis"
)

type RedisClient interface {
	Persist(key string) error
	IncrByFloat(key string, value float64) (float64, error)
	RedisInit(addr string, password string, db int, isCluster bool)
	RedisCheckPing(addr string, password string, db int, isCluster bool)
	IsExistKey(key string) bool
	Set(key string, value any, expiration time.Duration) error
	Get(key string) (string, error)
	Del(key string) error
	SetKeyNX(key string, value any, expiration time.Duration) (bool, error)
	Publish(channel string, message any) error
	TTL(key string) (time.Duration, error)
	Expire(key string, expiration time.Duration) error
	ExpireAt(key string, tm time.Time) error
	HSet(key string, field string, value any) error
	HSetNX(key string, field string, value any) (bool, error)
	HMSet(key string, value map[string]any) error
	HGet(key string, field string) (string, error)
	HExists(key string, field string) (bool, error)
	HMGet(key string, fields ...string) ([]any, error)
	HDel(key string, field string) error
	HKeys(key string) ([]string, error)
	HGetAll(key string) (map[string]string, error)
	HLen(key string) (int64, error)
	Ping() (string, error)
	SAdd(key string, value any) (int64, error)
	SMembers(key string) ([]string, error)
	SRem(key string, value any) error
	SIsMember(key string, value any) (bool, error)
	ZAdd(key string, score float64, member string) error
	ZScore(key string, member string) (float64, error)
	ZRange(key string, start, stop int64) ([]string, error)
	ZRangeWithScores(key string, start, stop int64) ([]redis.Z, error)
	ZRevRange(key string, start, stop int64) ([]string, error)
	ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error)
	ZIncrBy(key string, increment float64, member string) (float64, error)
	ZRem(key string, members ...any) error
	ZUnionStore(dest string, store redis.ZStore, keys []string) error
	ZCard(key string) (int64, error)
	ZRank(key string, member string) (int64, error)
	ZRevRank(key string, member string) (int64, error)
	RPush(key string, value any) error
	LLen(key string) (int64, error)
	LRange(key string, start, stop int64) ([]string, error)
	LIndex(key string, index int64) (string, error)
	Rename(key string, newKey string) error
}

func GetRedisClient(addr string, password string, db int, isCluster bool, isNil bool) RedisClient {
	if isNil {
		return nil
	}
	if isCluster {
		rdb := &cRedis.RedisClient{}
		rdb.RedisInit(addr, password, db, isCluster)
		return rdb
	} else {
		rdb := &sRedis.RedisClient{}
		rdb.RedisInit(addr, password, db, isCluster)
		return rdb
	}
}
