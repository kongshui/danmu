package service

import (
	"context"
	"sync"
	"time"

	"github.com/kongshui/danmu/model/pmsg"

	"github.com/kongshui/danmu/zilog"

	dao_redis "github.com/kongshui/danmu/dao/redis"

	dao_mysql "github.com/kongshui/danmu/dao/mysql"

	conf "github.com/kongshui/danmu/conf/web"

	dao_etcd "github.com/kongshui/danmu/dao/etcd"

	"github.com/kongshui/danmu/common"
)

type (
	SingleUserFunc  func(string, string) error
	SetWinScoreFunc func(string, RoundUploadStruct) error
	LotteryFunc     func(string, string, int64) map[string]int64
)

var (
	playerGroupAddinFunc SingleUserFunc
	setWinnerScoreFunc   SetWinScoreFunc
	lotteryFunc          LotteryFunc
	is_mock              bool
	config               *conf.Config
	accessToken          *AccessTokenStruct = &AccessTokenStruct{Lock: &sync.RWMutex{}} //全局token使用
	// isNotMock            bool                                                           //是否不模拟
	debug                 bool               //是否调试
	giftToScoreMap        map[string]float64 //礼物对应的积分
	commentToScore        map[string]float64 //评论对应的积分
	winCoinToComment      map[int64]string   //连胜币对应的评论
	commentToCoin         map[string]int64   //连胜币对应的评论
	commentTogiftId       map[string]string  //连胜币对应的礼物Id
	giftIdToName          map[string]string  //礼物id对应的礼物名称
	nodeIdToIntegral      map[int64]int64    //节点id对应的积分
	expireTime            time.Duration      //过期时间
	currentRankVersion    string             //世界排行版version
	nowMonth              string             //当前月
	monthVersionRankDb    string             //月排行版db名称
	app_id                string             // appId
	app_secret            string             // appSecret
	platform              string             // 平台
	url_GetAccessTokenUrl string             // 获取全局token的url
	// rdb                   = dao_redis.GetRedisClient(config.Redis.Addr, config.Redis.Password, config.Redis.IsCluster, false)
	rdb         dao_redis.RedisClient
	mysql       = dao_mysql.NewMysqlClient()
	nodeUuid    string //全局uid，识别客户端
	is_maintain bool   //是否维护
	is_connect  bool   //是否连接
	// pubsubNameList       []string                 // 推送数据的关键字
	forward_domain = common.NewStringList() //前端转发域名
	//grpc域名
	// grpc_domain = common.NewStringList() //前端转发域名
	// grpcpool
	grpc_pool = newTCPConnectionPool(20)
	// testMode             bool                     //测试模式
	ziLog zilog.LogStruct
	// isMock          bool    = true //是否模拟
	// testChat []string = make([]string, 0)
	msgBodyPool sync.Pool = sync.Pool{New: func() any {
		return &pmsg.MessageBody{}
	}}
	baseDataPool sync.Pool = sync.Pool{New: func() any {
		return &KsCallbackStruct{}
	}}

	bytePool sync.Pool = sync.Pool{New: func() any {
		out := make([]byte, 0)
		return &out
	}}
	platFormPool sync.Pool = sync.Pool{New: func() any {
		return &pmsg.PlatFormDataSend{}
	}}
	etcdClient = dao_etcd.NewEtcd()
	// chanPool   = NewChanPool(10)
	first_ctx  = context.Background()
	scrollDay  time.Weekday
	scrollHour int
)
