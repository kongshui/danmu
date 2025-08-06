package service

import (
	"time"

	battlematchv1 "github.com/kongshui/danmu/battlematch/v1"
	conf "github.com/kongshui/danmu/conf/web"
	dao_etcd "github.com/kongshui/danmu/dao/etcd"
	dao_mysql "github.com/kongshui/danmu/dao/mysql"
	dao_redis "github.com/kongshui/danmu/dao/redis"
	"github.com/kongshui/danmu/zilog"
)

// 函数初始化
func ServiceFuncSet(suf SingleUserFunc, sts SetWinScoreFunc, ltf LotteryFunc, webSocketFunc WebsocketFunc) {
	playerGroupAddinFunc = suf
	setWinnerScoreFunc = sts
	lotteryFunc = ltf
	otherWebsocketFunc = webSocketFunc
}

// 所有链接初始化
func ConnectInit(conf *conf.Config, etcClient *dao_etcd.Etcd, mysqlClient *dao_mysql.MysqlClient, redisClient dao_redis.RedisClient, logWirte zilog.LogStruct) {
	config = conf
	etcdClient = etcClient
	battlematchv1.InitEtcd(etcClient)
	battlematchv1.InitProjectName(config.Project)
	mysql = mysqlClient
	rdb = redisClient
	ziLog = logWirte
}

// 初始化全局变量
func InitGlobalVar(isPkMatch, levelScorll bool, storeLevel int64) {
	is_pk_match = isPkMatch       // 是否开启pk匹配
	is_level_scroll = levelScorll // 是否开启等级滚动
	store_level = storeLevel      // 存储等级
}

// 初始化时间
func InitTime(d time.Weekday, h int, expireT time.Duration) {
	scrollDay = d
	scrollHour = h
	expireTime = expireT
}
