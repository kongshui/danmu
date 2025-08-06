package service

import (
	"log"
	"os"

	"github.com/google/uuid"
)

/*
需要解决问题：
pushBasePayloay.go : 推送数据需要进行更改，不是推送到抖音
pushDownLoadMessage : 推送下行消息,是否需要进行更改
commentkeyList : 评论关键字
accessToken : 全局token
*/

func ServiceInit() {
	// 设置全局变量
	is_mock = config.App.IsMock
	app_id = config.App.AppId         // appId
	app_secret = config.App.AppSecret // appSecret
	platform = config.App.PlatForm    // 平台
	debug = config.Server.Debug
	setUrl()
	monthVersionSet()
	//日志初始化
	// fmt.Println(config.Logging.LogPath, config.Logging.Level, config.Logging.MaxSize, config.Logging.MaxBackups, config.Logging.MaxAge)
	// ziLog.Init(config.Logging.LogPath, config.Logging.Level, config.Logging.MaxSize, config.Logging.MaxBackups, config.Logging.MaxAge)

	//设置uuid
	nodeUuid = uuid.New().String()
	//初始化mysql
	// mysql.MysqlInit(config.Mysql.Username, config.Mysql.Password, config.Mysql.Addr, config.Mysql.Db, config.Mysql.IsUse)

	//初始化推送消息
	// pubsubNameList = []string{"live_comment", "live_gift", "live_like"}

	// 设置token
	if !is_mock {
		if err := setToken(); err != nil {
			log.Println("设置token失败： ", err)
			ziLog.Error("设置token失败： "+err.Error(), config.Server.Debug)
			os.Exit(1)
		}
		go setAccessToken()
	}

	// 初始化礼物和名称
	// ksGiftIdNameInit()
	// //礼物积分初始胡
	// giftToScoreInit()
	// 获取token
	if !is_mock {
		if err := getToken(); err != nil {
			log.Println("获取token失败： ", err)
			ziLog.Error("获取token失败： "+err.Error(), config.Server.Debug)
			os.Exit(1)
		}
		go getAccessToken()
	}

	//初始化世界排行版
	if err := worldRankInit(); err != nil {
		log.Println("初始化世界排行版失败： ", err)
		ziLog.Error("初始化世界排行版失败： "+err.Error(), config.Server.Debug)
		os.Exit(1)
	}

	// 初始化etcd
	// etcdClient.InitEtcd(config.Etcd.Addr, config.Etcd.Username, config.Etcd.Password)

	//周二自动滚动
	go autoNewVersion()

	//注册后端域名
	go registerBackDomain(first_ctx)

	//获取前端域名
	// go getFowardDomain(first_ctx)

	// 获取grpc
	// if config.Server.Grpc {
	// 	go getGrpcDomain(first_ctx)
	// 	// 匹配心跳
	// }
	if is_pk_match {
		go matchV1HeardBeat()
	}
	// testChat
	// go TestChat()
	// 失败消息获取
	if !is_mock {
		go getFailMessage()
	}

	// 检查断线状态
	if !is_mock {
		go checkDisconnectRoomIdExpire()
	}

	// 平台分开推送的内容
	switch platform {
	case "ks":
		if lotteryFunc == nil {
			ziLog.Error("快手抽奖函数未设置", debug)
			os.Exit(1)
		}
	case "dy":
		// 设置世界排行版生效版本
		if is_mock {
			break
		}
		worldRankSet(currentRankVersion)
		go pushWorldRankDataEntry()
		go pushHistoryWorldRankDataEntry()
	}
}
