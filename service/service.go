package service

import (
	"log"
	"os"
	"time"

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
	is_mock = cfg.App.IsMock
	app_id = cfg.App.AppId         // appId
	app_secret = cfg.App.AppSecret // appSecret
	platform = cfg.App.PlatForm    // 平台
	debug = cfg.Web.Debug
	setUrl()
	// monthVersionSet()

	//设置uuid
	nodeUuid = uuid.New().String()

	if !is_mock {
		// 设置token
		if err := setToken(); err != nil {
			log.Println("设置token失败： ", err)
			ziLog.Error("设置token失败： "+err.Error(), debug)
			os.Exit(1)
		}
		go setAccessToken()
		//获取token
		if err := getToken(); err != nil {
			log.Println("获取token失败： ", err)
			ziLog.Error("获取token失败： "+err.Error(), debug)
			os.Exit(1)
		}
		go getAccessToken()
		//初始化世界排行版
		if err := worldRankInit(); err != nil {
			log.Println("初始化世界排行版失败： ", err)
			ziLog.Error("初始化世界排行版失败： "+err.Error(), debug)
			os.Exit(1)
		}
		// 失败消息获取
		go getFailMessage()
		// 检查断线状态
		go checkDisconnectRoomIdExpire()

		// 匹配心跳
		if is_pk_match {
			go matchV1HeardBeat()
		}
	} else {
		currentRankVersion = time.Now().Format(version_time_layout)
	}

	// 初始化etcd
	// etcdClient.InitEtcd(config.Etcd.Addr, config.Etcd.Username, config.Etcd.Password)

	//周二自动滚动
	if scrollAuto != nil {
		go scrollAuto(&currentRankVersion)
	}
	// go autoNewVersion()

	//注册后端域名
	go registerBackDomain(first_ctx)

	if getCfgFunc != nil {
		getCfgFunc(&cfgConfig)
	}
	//获取前端域名
	// go getFowardDomain(first_ctx)

	// 获取grpc
	// if config.Server.Grpc {
	// 	go getGrpcDomain(first_ctx)
	// 	// 匹配心跳
	// }
	// 初始化函数
	if initService != nil {
		initService(is_mock)
	}
	// 读取配置映射
	if err := configMapRead(); err != nil {
		ziLog.Error("读取配置映射失败： "+err.Error(), debug)
		os.Exit(1)
	}
	// 自动检测配置变更
	go autoDetectConfigChange()
	// 自动删除日志文件
	go autoDeleteLogFile()

	// 平台分开推送的内容
	switch platform {
	case "ks":
		if lottery == nil {
			ziLog.Error("快手抽奖函数未设置", debug)
			os.Exit(1)
		}
	case "dy":
		// 设置世界排行版生效版本
		if is_mock {
			break
		}
		WorldRankSet(currentRankVersion)
		if cfg.App.NoSend {
			break
		}
		go pushWorldRankDataEntry()
		go pushHistoryWorldRankDataEntry()
	}
}
