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
func ServiceFuncSet(suf SingleUserFunc, sts SetWinScoreFunc, ltf LotteryFunc, webSocketFunc WebsocketFunc, interFunc InteractiveFunc, giftExtendInfo GiftExtendInfoFunc) {
	interactive = interFunc
	playerGroupAddinFunc = suf
	setWinnerScoreFunc = sts
	lotteryFunc = ltf
	otherWebsocketFunc = webSocketFunc
	giftExtendInfoFunc = giftExtendInfo
}

// 所有链接初始化
func ConnectInit(conf *conf.Config, etcClient *dao_etcd.Etcd, mysqlClient *dao_mysql.MysqlClient, redisClient dao_redis.RedisClient, logWirte *zilog.LogStruct) {
	config = conf
	etcdClient = etcClient
	battlematchv1.InitEtcd(etcClient)
	battlematchv1.InitProjectName(config.Project)
	mysql = mysqlClient
	rdb = redisClient
	ziLog = logWirte
}

// 初始化全局变量isPkMatch: 是否开启pk匹配, levelScorll: 是否开启等级滚动, storeLevel: 存储等级, liveLike: 直播点赞积分, versionTimeInterval: 版本时间间隔, groupIdList: 组名, weekSet: 多长时间滚动

func InitGlobalVar(isPkMatch, levelScorll bool, storeLevel int64, liveLike float64, versionTimeInterval int64, groupIdList []string, weekSet int) {
	is_pk_match = isPkMatch                     // 是否开启pk匹配
	is_level_scroll = levelScorll               // 是否开启等级滚动
	store_level = storeLevel                    // 存储等级
	live_like_score = liveLike                  // 直播点赞积分
	version_time_interval = versionTimeInterval // 版本滚动时间间隔
	groupid_list = groupIdList                  // 组名
	week_set = weekSet                          // 1 是一周，2是两周 ，3是三周 4是四周 0是一个月
}

// 初始化时间
func InitTime(d time.Weekday, h int, expireT time.Duration, day int) {
	scrollDay = d
	scrollHour = h
	expireTime = expireT
	month_day = day
}
