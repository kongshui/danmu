package service

import (
	battlematchv1 "github.com/kongshui/danmu/battlematch/v1"
	conf "github.com/kongshui/danmu/conf/nodeweb"
	dao_etcd "github.com/kongshui/danmu/dao/etcd"
	dao_mysql "github.com/kongshui/danmu/dao/mysql"
	dao_redis "github.com/kongshui/danmu/dao/redis"
	"github.com/kongshui/danmu/zilog"
)

/*
函数初始化：
addGroupSetFunc：玩家入组时调用的函数
setWinningScoreFunc：对局结束后设置玩家的积分、胜点、连胜等
lotteryFunc：抽奖函数
webSocketFunc：websocket函数
interFunc：快手加边函数，其实可以包装更加深一点，暂时不动
giftExtendInfo: 礼物置顶函数
scrollfunc：自动滚动函数
initFunction：初始化函数
getCfgFunc：获取读取到的配置函数
*/
func ServiceFuncSet(addGroupSetFunc SingleUserAddGroupSetFunc, setWinningScoreFunc SetWinScoreFunc, lotteryFunc LotteryFunc, webSocketFunc WebsocketFunc, interFunc InteractiveFunc,
	giftExtendInfoFunc GiftExtendInfoFunc, scrollFunc ScrollFunc, initFunc InitFunc, setIntegralToRoundFunc SetIntegralToRoundFunc, gCfgFunc GetCfgConfigFunc,
	sendMsgToGatewayFunc SendMessageToGatewayFunc) {

	interactive = interFunc
	playerGroupAddin = addGroupSetFunc
	setWinnerScore = setWinningScoreFunc
	lottery = lotteryFunc
	otherWebsocket = webSocketFunc
	giftExtendInfos = giftExtendInfoFunc
	scrollAuto = scrollFunc
	initService = initFunc
	setIntegralToRound = setIntegralToRoundFunc
	getCfgFunc = gCfgFunc
	sendMessageToGateway = sendMsgToGatewayFunc
}

/*
所有链接初始化：
conf： 配置文件
etcClient：etcd链接
mysqlClient：mysql链接
redisClient：redis链接
logWirte：日志
*/
func ConnectInit(conf *conf.Config, etcClient *dao_etcd.Etcd, mysqlClient *dao_mysql.MysqlClient, redisClient dao_redis.RedisClient, logWirte *zilog.LogStruct) {
	cfg = conf
	etcdClient = etcClient
	battlematchv1.InitEtcd(etcClient)
	battlematchv1.InitProjectName(cfg.Project)
	mysql = mysqlClient
	rdb = redisClient
	ziLog = logWirte
}

/*
初始化全局变量：
isPkMatch: 是否开启pk匹配
levelScorll: 是否开启等级滚动
storeLevel: 存储等级
liveLike: 直播点赞积分
groupIdList: 组名
*/
func InitGlobalVar(isPkMatch bool, storeLevel int64, liveLike float64, groupIdList []string) {
	is_pk_match = isPkMatch    // 是否开启pk匹配
	store_level = storeLevel   // 存储等级
	live_like_score = liveLike // 直播点赞积分
	groupid_list = groupIdList // 组名
}

/*
初始化时间：
d: 星期几
h: 几点
expireT：过期时间，如果isCalculate为true，那么这个时间不能为0，但不参与运算
day: 是每月几号
isCalculate 即直至本期结束后才会过期，即为本期内不过期
weekSet: 多长时间滚动
versionTimeInterval: 滚动版本最小时间间隔
*/
func InitTime(setScrollTime *SetScrollTime) {
	// scrollDay = d
	// scrollHour = h
	// expireTime = expireT
	// month_day = day
	// isCalculateExpireTime = isCalculate
	// week_set = weekSet // 1 是一周，2是两周 ，3是三周 4是四周 0是一个月
	scrollTime = setScrollTime
}
