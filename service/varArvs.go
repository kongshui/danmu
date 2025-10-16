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
	SingleUserAddGroupSetFunc func(roomId, uidStr string, roundId int64, userMap []*pmsg.SingleRoomAddGroupInfo, isChoose bool) error
	SetWinScoreFunc           func(string, RoundUploadStruct, int64) error
	LotteryFunc               func(string, string, int64) map[string]int64
	WebsocketFunc             func(msg *pmsg.MessageBody) error
	InteractiveFunc           func(roomId, roundId string, label int) bool //自动选边
	GiftExtendInfoFunc        func() string
	ScrollFunc                func(string)
	InitFunc                  func(bool)
	SetIntegralToRoundFunc    func(roomId, openId string, score float64) error
)

var (
	playerGroupAddin   SingleUserAddGroupSetFunc
	setWinnerScore     SetWinScoreFunc
	lottery            LotteryFunc
	otherWebsocket     WebsocketFunc
	interactive        InteractiveFunc
	giftExtendInfos    GiftExtendInfoFunc
	scrollAuto         ScrollFunc
	initService        InitFunc
	setIntegralToRound SetIntegralToRoundFunc

	is_mock bool

	config      *conf.Config
	accessToken *AccessTokenStruct = &AccessTokenStruct{Lock: &sync.RWMutex{}} //全局token使用
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
	is_pk_match bool   //是否开启pk匹配
	// pubsubNameList       []string                 // 推送数据的关键字
	forward_domain = common.NewStringList() //前端转发域名
	//grpc域名
	// grpc_domain = common.NewStringList() //前端转发域名
	// grpcpool
	// grpc_pool = newTCPConnectionPool(20)
	// testMode             bool                     //测试模式
	ziLog *zilog.LogStruct
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
	first_ctx             = context.Background()
	scrollDay             time.Weekday
	scrollHour            int
	is_level_scroll       bool     //是否开启等级滚动
	store_level           int64    //存储等级
	live_like_score       float64  //直播点赞积分
	version_time_interval int64    //版本时间间隔
	groupid_list          []string //组名
	week_set              int      // 1 是一周，2是两周 ，3是三周 4是四周 0是一个月
	month_day             int      // 每月第几天滚动
	isCalculateExpireTime bool     //是否计算过期时间

)
