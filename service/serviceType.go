package service

import "sync"

type (
	//token消息结构体
	// LoginTokenStruct struct {
	// 	Token     string `json:"token"`      // token,客户端传过来的token
	// 	TimeStamp int64  `json:"time_stamp"` // 时间戳
	// }

	//断开消息结构体
	// DisconnectStruct struct {
	// 	Roomid    string `json:"room_id"`    // roomid
	// 	TimeStamp int64  `json:"time_stamp"` // 时间戳
	// }

	//错误信息结构体
	// ErrorStruct struct {
	// 	ErrorMsg string `json:"err_msg"` // 错误信息
	// }

	//用户加入分组消息结构体
	JoinGroupStruct struct {
		GroupId   string `json:"group_id"`   // 分组id
		OpenId    string `json:"open_id"`    // 用户open_id
		AvatarUrl string `json:"avatar_url"` // 评论用户头像地址
		NickName  string `json:"nick_name"`  // 评论用户昵称
	}

	//所有用户信息
	// AllUserStruct struct {
	// 	UserList []JoinGroupStruct `json:"user_list" required:"true"` // 用户open_id列表
	// }

	//roundReady消息结构体
	// RoundReadyStruct struct {
	// 	RoundId   int64 `json:"round_id"`   // roomid
	// 	TimeStamp int64 `json:"time_stamp"` // 时间戳
	// }

	// 自定义消息结构体
	// MessageGeneralStruct struct {
	// 	MsgId        string `json:"msg_id"`         // 消息ID
	// 	MsgType      string `json:"msg_type"`       // 消息类型
	// 	RoomId       string `json:"room_id"`        // 房间id
	// 	AnchorOpenId string `json:"anchor_open_id"` // 主播open_id
	// 	Data         any    `json:"data"`           // 消息数据
	// 	ExtraData    any    `json:"extra_data"`     // 额外数据
	// }

	//logToRedisStruct
	// LogToRedisStruct struct {
	// 	RoomId          string                  `json:"room_id"`           // 房间id
	// 	AnchorOpenId    string                  `json:"anchor_open_id"`    // 主播open_id
	// 	RoundId         int64                   `json:"round_id"`          // 对局id,同一个直播间room_id下，round_id需要是递增的，建议使用开始对局时的时间戳
	// 	GroupResultList []GroupResultList       `json:"group_result_list"` //对局结果列表
	// 	GroupUserList   []UserUploadScoreStruct `json:"group_user_list"`   //玩家信息列表
	// }

	//round_upload消息结构体
	RoundUploadStruct struct {
		RoundId         int64                   `json:"round_id"`          // 对局id,同一个直播间room_id下，round_id需要是递增的，建议使用开始对局时的时间戳
		GroupResultList []GroupResultList       `json:"group_result_list"` //对局结果列表
		GroupUserList   []UserUploadScoreStruct `json:"group_user_list"`   //玩家信息列表
	}

	//玩家worldInfo消息体
	WorldInfoStruct struct {
		OpenId            string `json:"open_id"`             // 房间id
		Score             int64  `json:"score"`               // 玩家世界分数
		Rank              int64  `json:"rank"`                //玩家世界排行
		AvatarUrl         string `json:"avatar_url"`          // 评论用户头像地址
		NickName          string `json:"nick_name"`           // 评论用户昵称
		WinningStreamCoin int64  `json:"winning_stream_coin"` // 连胜币多少
		// SwallowCount      int64  `json:"swallow_count"`       // 吞噬数量
	}

	//用户积分列表
	UserUploadScoreStruct struct {
		GroupId     string `json:"group_id"`     //分组id
		OpenId      string `json:"open_id"`      //用户open_id
		Rank        int64  `json:"rank"`         //用户排名，从1开始超过1000的，可以固定传递1000，抖音端会展示为 "999+"
		RoundResult int    `json:"round_result"` //对局结果（1=胜利、2=失败、3=平局）
		Score       int64  `json:"score"`        //核心数值,用户排名的依据
	}

	//评论消息 live_comment
	ContentPayloadStruct struct {
		MsgId     string `json:"msg_id"`     // 消息ID
		SecOpenid string `json:"sec_openid"` // 评论用户的加密openid, 当前其实没有加密
		Content   string `json:"content"`    // 评论内容
		AvatarUrl string `json:"avatar_url"` // 评论用户头像地址
		Nickname  string `json:"nickname"`   // 评论用户昵称
		TimeStamp int64  `json:"timestamp"`  // 时间戳
	}

	//送礼消息
	GiftPayloadStruct struct {
		MsgId             string `json:"msg_id"`               // 消息ID
		SecOpenid         string `json:"sec_openid"`           // 评论用户的加密openid, 当前其实没有加密
		SecGiftId         string `json:"sec_gift_id"`          // 加密的礼物id
		GiftNum           int    `json:"gift_num"`             // 礼物数量
		GiftValue         int    `json:"gift_value"`           // 礼物总价值，单位分
		AvatarUrl         string `json:"avatar_url"`           // 评论用户头像地址
		Nickname          string `json:"nickname"`             // 评论用户昵称
		TimeStamp         int64  `json:"timestamp"`            // 时间戳
		Test              bool   `json:"test"`                 // 如果是抖音平台的测试数据，则会下发该字段且值为 true。测试工具下发的送礼数据属于调试模式，不会有该字段
		AudienceSecOpenId string `json:"audience_sec_open_id"` // 被送礼的嘉宾openid，当前没有加密
		SecMagicGiftId    string `json:"sec_magic_gift_id"`    // 被送礼的嘉宾openid，当前没有加密(备用字段)
	}

	//点赞消息 live_like
	LiveLikePayloadStruct struct {
		MsgId     string `json:"msg_id"`     // 消息ID
		SecOpenid string `json:"sec_openid"` // 评论用户的加密openid, 当前其实没有加密
		LikeNum   int64  `json:"like_num"`   // 点赞数量，上游2s合并一次数据
		AvatarUrl string `json:"avatar_url"` // 评论用户头像地址
		Nickname  string `json:"nickname"`   // 评论用户昵称
		TimeStamp int64  `json:"timestamp"`  // 时间戳
	}

	// 粉丝团数据 live_fansclub
	FansPayloadStruct struct {
		MsgId              string `json:"msg_id"`               // 消息ID
		SecOpenid          string `json:"sec_openid"`           // 评论用户的加密openid, 当前其实没有加密
		AvatarUrl          string `json:"avatar_url"`           // 评论用户头像地址
		Nickname           string `json:"nickname"`             // 评论用户昵称
		TimeStamp          int64  `json:"timestamp"`            // 时间戳
		FansclubReasonType int    `json:"fansclub_reason_type"` //粉丝团消息类型：1-升级，2-加团
		FansclubLevel      int    `json:"fansclub_level"`       //用户粉丝团等级，加团消息下默认传1
	}

	// 推送失败请求结构体
	GetLiveFailRequestStruct struct {
		Roomid   string `json:"roomid"`    //直播间id
		Appid    string `json:"appid"`     //appid
		MsgType  string `json:"msg_type"`  //消息类型
		PageNum  string `json:"page_num"`  //页码：注意，需要从1开始
		PageSize string `json:"page_size"` //每页数据条数, 最大不超过100
	}

	// AccessToken全局token结构体
	AccessTokenStruct struct {
		Token     string        //全局token
		ExpiresIn int           //过期时间，暂时没用
		Lock      *sync.RWMutex //全局锁
	}

	// 请求快手接口获取AccessToken返回结构体
	AccessKsTokenRespondStruct struct {
		Result      int64  `json:"result"`       // 返回结果,1成功，其他失败
		AccessToken string `json:"access_token"` // access_token
		ExpiresIn   int64  `json:"expires_in"`   // 过期时间
		TokenType   string `json:"token_type"`   // token类型
	}

	// 请求抖音接口获取AccessToken返回结构体
	AccessDyTokenRespondStruct struct {
		ErrorNo   int64  `json:"err_no"`   //错误码
		ErrorTips string `json:"err_tips"` //错误提示
		Data      struct {
			AccessToken string `json:"access_token"` //全局token
			ExpiresIn   int64  `json:"expires_in"`   //过期时间
			ExpiresAt   int64  `json:"expires_at"`   //过期时间戳
		} `json:"data"` //数据
	}

	// 抖音获取主播信息结构体
	GetDyAnchorInfoStruct struct {
		Data struct {
			Ackcfg     []any          `json:"ack_cfg"`     // 预留信息，sdk接入使用，开发者不用感知
			LinkerInfo map[string]any `json:"linker_info"` // 连屏数据预留信息，开发者目前不用感知
			Info       struct {
				RoomId       int64  `json:"room_id"`        // 房间id
				NickName     string `json:"nick_name"`      // 昵称
				AvatarUrl    string `json:"avatar_url"`     // 头像url
				AnchorOpenId string `json:"anchor_open_id"` // 主播open_id
			} `json:"info"` // 主播信息
		} `json:"data"` // data信息
		Errcode    int64          `json:"errcode,omitempty"`     //错误信息
		Errmsg     string         `json:"errmsg,omitempty"`      // 错误信息
		Extra      map[string]any `json:"extra,omitempty"`       // 额外信息
		StatusCode int64          `json:"status_code,omitempty"` //状态码
	}

	// 请求抖音接口地址结构体
	// UrlStruct struct {
	// 	GetAccessTokenUrl          string `json:"get_access_token"`       //获取全局token的url
	// 	GetLiveInfoUrl             string `json:"get_liveinfo"`           //获取直播间信息的url
	// 	StartGameUrl               string `json:"start_game"`             //开始游戏url
	// 	FinishGameUrl              string `json:"finish_game"`            //结束游戏url
	// 	GetGameInfoUrl             string `json:"get_gameinfo"`           //获取游戏任务是否开启url
	// 	PushLivePayLoadUrl         string `json:"push_live_payload"`      //推送直播消息的url
	// 	DownLineMessageUrl         string `json:"down_line_message"`      //下线消息推送url
	// 	WorldRankVersionUrl        string `json:"world_rank"`             //世界榜url
	// 	WorldRankUploadUrl         string `json:"world_rank_list_upload"` //世界榜上传url
	// 	WorldRankCompleteUploadUrl string `json:"complete_world_rank"`    //完成用户世界榜单的累计战绩上报url
	// 	WorldRankUserTotalUrl      string `json:"world_rank_user_total"`  //世界榜用户累计战绩url
	// 	SyncGameStatusUrl          string `json:"sync_game_status"`       //同步对局状态url
	// 	UploadUserGameUrl          string `json:"upload_user_game"`       //上传用户对局数据url
	// 	UploadUserGameRankUrl      string `json:"upload_user_game_rank"`  //上传用户对局数据url
	// 	UploadUserCompeleteUrl     string `json:"upload_user_complete"`   //上传用户对局数据url
	// 	GetLiveFailUrl             string `json:"get_live_fail"`          //获取直播失败请求url
	// 	GetConnIdUrl               string `json:"get_conn_id"`            //获取直播失败请求url
	// 	UploadUserGroupUrl         string `json:"upload_user_group"`      //上传用户分组数据url
	// }

	//抖音直播间结构体
	// RoomInfoStruct struct {
	// 	RoomId       string `json:"room_id"`        // 房间id
	// 	NickName     string `json:"nick_name"`      // 昵称
	// 	AvatarUrl    string `json:"avatar_url"`     // 头像url
	// 	AnchorOpenId string `json:"anchor_open_id"` // 主播open_id
	// }

	//抖音小游戏结构体
	// AppInfoStruct struct {
	// 	// Appid              string `json:"appid"`                // appid
	// 	// AppSecret          string `json:"secret"`               // app_secret
	// 	LiveCommentSecret  string `json:"live_comment_secret"`  // 直播评论加密密钥
	// 	LiveGiftSecret     string `json:"live_gift_secret"`     // 直播送礼加密密钥
	// 	LiveLikeSecret     string `json:"live_like_secret"`     // 直播点赞加密密钥
	// 	LiveFansclubSecret string `json:"live_fansclub_secret"` // 直播粉丝团加密密钥
	// }

	// 用户世界榜单结构体
	UserRankStruct struct {
		OpenId             string `json:"open_id"`              // 用户open_id
		Rank               int64  `json:"rank"`                 // 世界榜单排名，从1开始
		Score              int64  `json:"score"`                // 当前用户的世界榜单积分
		WinningStreakCount int64  `json:"winning_streak_count"` // 当前用户的连胜次数，如果没有连胜记录传0
		WinningPoints      int64  `json:"winning_points"`       // 当前用户的胜点记录，如果没有胜点记录传0
	}

	// 世界榜单结构体
	WorldRankListStruct struct {
		AppId            string           `json:"app_id"`             // 小玩法app_id
		IsOnlineVersion  bool             `json:"is_online_version"`  // 是否是线上版本，默认为false，为false代表测试数据
		WorldRankVersion string           `json:"world_rank_version"` // 开发者指定上传到的榜单版本
		RankList         []UserRankStruct `json:"rank_list"`          // 世界榜单列表和用户累计战绩数据列表类型一样
	}

	// 世界榜单累计结构体
	WorldRankUserListStruct struct {
		AppId            string           `json:"app_id"`             // 小玩法app_id
		IsOnlineVersion  bool             `json:"is_online_version"`  // 是否是线上版本，默认为false，为false代表测试数据
		WorldRankVersion string           `json:"world_rank_version"` // 开发者指定上传到的榜单版本
		UserList         []UserRankStruct `json:"user_list"`          // 世界榜单列表和用户累计战绩数据列表类型一样
	}

	//同步对局状态结构体
	SyncGameStatusStruct struct {
		AnchorOpenId    string            `json:"anchor_open_id"`    // 主播open_id
		AppId           string            `json:"app_id"`            // 小玩法app_id
		RoomId          string            `json:"room_id"`           // 房间id
		RoundId         int64             `json:"round_id"`          // 对局id,同一个直播间room_id下，round_id需要是递增的，建议使用开始对局时的时间戳
		StartTime       int64             `json:"start_time"`        // 本局开始时间，秒级时间戳
		EndTime         int64             `json:"end_time"`          // 本局结束时间，秒级时间戳,同步的对局状态为对局结束时，该字段必传。
		Status          int               `json:"status"`            // 当前房间的游戏对局状态（1=已开始、2=已结束）
		GroupResultList []GroupResultList `json:"group_result_list"` //对局结果列表
	}

	//round_end 结构体
	//同步对局状态结构体
	RoundEndStruct struct {
		RoundId         int64             `json:"round_id"`          // 对局id,同一个直播间room_id下，round_id需要是递增的，建议使用开始对局时的时间戳
		StartTime       int64             `json:"start_time"`        // 本局开始时间，秒级时间戳
		EndTime         int64             `json:"end_time"`          // 本局结束时间，秒级时间戳,同步的对局状态为对局结束时，该字段必传。
		GroupResultList []GroupResultList `json:"group_result_list"` //对局结果列表
	}

	//玩家resultList
	GroupResultList struct {
		GroupId string `json:"group_id"` //分组id
		Result  int    `json:"result"`   //对局结果（1=胜利、2=失败、3=平局）
	}

	//上传用户对局数据结构体
	UploadUserGameStruct struct {
		AnchorOpenId string           `json:"anchor_open_id"` // 主播open_id
		AppId        string           `json:"app_id"`         // 小玩法app_id
		RoomId       string           `json:"room_id"`        // 房间id
		RoundId      int64            `json:"round_id"`       // 对局id,同一个直播间room_id下，round_id需要是递增的，建议使用开始对局时的时间戳
		UserList     []UserListStruct `json:"user_list"`      //用户对局数据列表
	}

	//上上报对局榜单列表结构体
	UploadRankGameStruct struct {
		AnchorOpenId string           `json:"anchor_open_id"` // 主播open_id
		AppId        string           `json:"app_id"`         // 小玩法app_id
		RoomId       string           `json:"room_id"`        // 房间id
		RoundId      int64            `json:"round_id"`       // 对局id,同一个直播间room_id下，round_id需要是递增的，建议使用开始对局时的时间戳
		RankList     []UserListStruct `json:"rank_list"`      //用户对局数据列表
	}

	//上传用户对局数据结构体
	UserListStruct struct {
		OpenId             string `json:"open_id"`              //用户open_id
		Rank               int64  `json:"rank"`                 //用户排名，从1开始超过1000的，可以固定传递1000，抖音端会展示为 "999+"
		RoundResult        int    `json:"round_result"`         //对局结果（1=胜利、2=失败、3=平局）
		Score              int64  `json:"score"`                //核心数值,用户排名的依据
		WinningPoints      int64  `json:"winning_points"`       //用户的胜点，如果没有胜点记录传0
		WinningStreakCount int64  `json:"winning_streak_count"` //用户的连胜次数，如果没有连胜记录传0
		GroupId            string `json:"group_id"`             //阵营Id，比如 red/blue
	}

	//完成用户对局数据上报结构体
	UploadUserGameCompleteStruct struct {
		AnchorOpenId string `json:"anchor_open_id"` // 主播open_id
		AppId        string `json:"app_id"`         // 小玩法app_id
		RoomId       string `json:"room_id"`        // 房间id
		RoundId      int64  `json:"round_id"`       // 对局id
		CompleteTime int64  `json:"complete_time"`  // 上传完成时间，由开发者传，秒级时间戳。默认为当前接口请求时间
	}

	// UserInfoStruct 用户信息
	UserInfoStruct struct {
		OpenId    string `json:"open_id"`    //用户open_id
		AvatarUrl string `json:"avatar_url"` // 评论用户头像地址
		NickName  string `json:"nick_name"`  // 评论用户昵称
	}

	// user_join_group 用户世界榜单信息 Message
	// UserChooseGroupStruct struct {
	// 	OpenId            string `json:"open_id"`             // 用户open_id
	// 	AvatarUrl         string `json:"avatar_url"`          // 评论用户头像地址
	// 	NickName          string `json:"nick_name"`           // 评论用户昵称
	// 	GroupId           string `json:"group_id"`            // 组信息
	// 	RoundId           int64  `json:"round_id"`            // 对局id
	// 	WorldScore        int64  `json:"world_score"`         // 世界排行榜分数
	// 	WorldRank         int64  `json:"world_rank"`          // 世界排行榜排名
	// 	WinningStreamCoin int64  `json:"winning_stream_coin"` // 连胜币多少
	// 	Isconsume         bool   `json:"is_consume"`          // 是否是第一次消费，ture为是，false为不是
	// }

	//result_group_user_worldinfo 返回玩家世界榜单信息 Message
	// ResultGroupUserWorldInfoStruct struct {
	// 	OpenId            string `json:"open_id"`             // 用户open_id
	// 	WorldScore        int64  `json:"world_score"`         // 世界排行榜分数
	// 	WorldRank         int64  `json:"world_rank"`          // 世界排行榜排名
	// 	WinningStreamCoin int64  `json:"winning_stream_coin"` // 连胜币多少
	// 	Isconsume         bool   `json:"is_consume"`          // 是否是第一次消费，ture为是，false为不是
	// }

	// use_winning_stream_coin 使用连胜币请求，Message
	// UseWinningStreamCoinStruct struct {
	// 	OpenId    string `json:"open_id"`    // 用户open_id
	// 	UseNum    int64  `json:"use_num"`    // 使用连胜币多少
	// 	RoundId   int64  `json:"round_id"`   // 对局id
	// 	GiftId    int    `json:"gift_id"`    // 礼物id
	// 	TimeStamp int64  `json:"time_stamp"` // 毫秒级时间戳
	// }

	// result_use_winning_stream_coin 返回使用连胜币请求，Message
	// ResultUseWinningStreamCoinStruct struct {
	// 	OpenId            string `json:"open_id"`             // 用户open_id
	// 	WinningStreamCoin int64  `json:"winning_stream_coin"` // 连胜币多少
	// 	CanUse            bool   `json:"can_use"`             // 是否能够使用
	// 	RoundId           int64  `json:"round_id"`            // 对局id
	// 	GiftId            int    `json:"gift_id"`             // 礼物id
	// 	TimeStamp         int64  `json:"time_stamp"`          // 毫秒级时间戳
	// }

	// get_winning_stream_coin 获得连胜币请求，Message
	// GetWinningStreamCoinStruct struct {
	// 	OpenId string `json:"open_id"` // 用户open_id
	// 	GetNum int64  `json:"get_num"` // 获得连胜币多少
	// }

	// result_get_winning_stream_coin 返回获得连胜币请求，Message
	// ResultGetWinningStreamCoinStruct struct {
	// 	OpenId            string `json:"open_id"`             // 用户open_id
	// 	WinningStreamCoin int64  `json:"winning_stream_coin"` // 连胜币多少
	// 	TimeStamp         int64  `json:"time_stamp"`          // 毫秒级时间戳
	// }

	// 维护信息
	WeihuStruct struct {
		IsMaintain  bool   `json:"is_maintain"`  // 是否维护中
		StartTime   string `json:"start_time"`   // 开始时间
		EndTime     string `json:"end_time"`     // 结束时间
		MaintainMsg string `json:"maintain_msg"` // 维护信息
	}

	//返回玩家上期前100排行榜信息 Message
	ResultGroupUserRankInfoStruct struct {
		OpenId string `json:"open_id"` // 用户open_id
		Rank   int    `json:"rank"`    // 排行名次
	}

	// get_winning_stream_coin 返回获得连胜币请求，Message
	// IsGetWinningStreamCoinStruct struct {
	// 	IsEnd    bool                         `json:"is_end"`    // 是否结束
	// 	UserList []GetWinningStreamCoinStruct `json:"user_list"` //需要获取连胜币的用户
	// }
	// open_id,time_stamp
	// OpenIdTimeStruct struct {
	// 	OpenId    string `json:"open_id"`    // 用户open_id
	// 	TimeStamp int64  `json:"time_stamp"` // 毫秒级时间戳
	// }

	// 快手直播间结构体
	KsRoomInfoStruct struct {
		UserId    string `json:"userId"`   // 房间id
		NickName  string `json:"userName"` // 昵称
		AvatarUrl string `json:"headUrl"`  // 头像url
	}

	//快手giftSend消息结构体
	KsGiftSendStruct struct {
		// UniqueMessageId string           `json:"uniqueMessageId"` // 唯一的消息id, 可用于幂等消费，第三方需要根据unique_message_id做好幂等控制
		UniqueNo       string           `json:"uniqueNo"`       // 一个单号代表一笔送收礼, 可用于幂等消费
		GiftId         string           `json:"giftId"`         // 礼物id, 可用于查询礼物信息
		GiftName       string           `json:"giftName"`       // 礼物名称, 可用于查询礼物信息
		GiftCount      int64            `json:"giftCount"`      // 礼物数量
		GiftUnitPrice  int64            `json:"giftUnitPrice"`  // 礼物单价, 单位快币 (1元=10快币)
		GiftTotalPrice int64            `json:"giftTotalPrice"` // 礼物总价, 单位快币 (1元=10快币)
		UserInfo       KsRoomInfoStruct `json:"userInfo"`       // 送礼者信息
		NormalGift     bool             `json:"normalGift"`     // 是否是普通礼物
	}

	//快手liveComment消息结构体
	KsLiveCommentStruct struct {
		Content  string           `json:"content"`  // 评论内容
		UserInfo KsRoomInfoStruct `json:"userInfo"` // 评论者信息
	}

	//快手liveLike消息结构体
	KsLiveLikeStruct struct {
		Count    int64            `json:"count"`    // 点赞数量
		UserInfo KsRoomInfoStruct `json:"userInfo"` // 评论者信息
	}

	// 快手回调消息data结构体
	KsCallbackDataStruct struct {
		UniqueMessageId string `json:"unique_message_id"` // 唯一的消息id, 可用于幂等消费，第三方需要根据unique_message_id做好幂等控制
		AuthorOpenId    string `json:"author_open_id"`    // 主播id
		RoomCode        string `json:"room_code"`         // 直播间id
		PushType        string `json:"push_type"`         // 消息类型
		Payload         []any  `json:"payload"`           // 消息体
	}

	// 快手礼物回调结构体

	// 快手回调消息data结构体
	KsCallbackQueryStruct struct {
		UniqueMessageId string             `json:"uniqueMessageId"` // 唯一的消息id, 可用于幂等消费，第三方需要根据unique_message_id做好幂等控制
		AuthorOpenId    string             `json:"authorOpenId"`    // 主播id
		RoomCode        string             `json:"roomCode"`        // 直播间id
		PushType        string             `json:"pushType"`        // 消息类型
		LiveTimeStamp   int64              `json:"liveTimestamp"`   //毫秒级时间戳
		Payload         []KsGiftSendStruct `json:"payload"`         // 消息体
	}

	//快手回调消息结构体
	KsCallbackStruct struct {
		MessageId string               `json:"message_id"` // 消息id，不幂等
		Event     string               `json:"event"`      // 事件名称，正式环境LIVE_INTERACTION_DATA  压测LIVE_INTERACTION_DATA_TEST
		AppId     string               `json:"app_id"`     // 小玩法app_id
		Data      KsCallbackDataStruct `json:"data"`       // 消息体
		TimeStamp int64                `json:"timestamp"`  // 时间戳，单位毫秒
	}

	// 快手回调返回结构体
	KsCallbackRespondStruct struct {
		Result   int64  `json:"result"`   // 如果result 不是 1， 会有error_msg
		ErrorMsg string `json:"errorMsg"` // 错误信息
	}

	// 快手消息验证结构体
	KsMsgAckReceiveStruct struct {
		UniqueMessageId     string `json:"uniqueMessageId"`                // 唯一的消息id, 可用于幂等消费，第三方需要根据unique_message_id做好幂等控制
		PushType            string `json:"pushType"`                       // 消息类型
		CpServerReceiveTime int64  `json:"cpServerReceiveTime,omitempty" ` // 消息接收时间，单位毫秒
		CpClientReceiveTime int64  `json:"cpClientReceiveTime,omitempty"`  // 消息接收时间，单位毫秒
		CpClientShowTime    int64  `json:"cpClientShowTime,omitempty"`     // 消息接收时间，单位毫秒
	}

	// 快手ack消息
	KsAckStruct struct {
		RoomCode  string `json:"roomCode"`  // 直播间id
		TimeStamp int64  `json:"timestamp"` // 时间戳，单位毫秒
		Sign      string `json:"sign"`      // 签名
		AckType   string `json:"ackType"`   // 消息类型
		Data      string `json:"data"`      // 消息体,receive时为KsMsgAckReceiveStruct，show时为KsMsgAckShowStruct
	}

	//快手TwoConnection结构体
	KsTwoConnectStruct struct {
		RoomCode  string `json:"roomCode"`  // 直播间id
		TimeStamp int64  `json:"timestamp"` // 时间戳，单位毫秒
		Sign      string `json:"sign"`      // 签名
		Method    string `json:"method"`    // 方法名
		BizType   string `json:"bizType"`   //业务类型
		Data      string `json:"data"`      // 消息体
	}
	// msg_ack 消息回执
	MsgAckStruct struct {
		RoomId  string `json:"room_id"`  // 房间id
		AppId   string `json:"app_id"`   // appid
		AckType int    `json:"ack_type"` // 消息回执类型
		Data    string `json:"data"`     // 消息回执数据
	}

	// msg_ack_Info
	MsgAckInfoStruct struct {
		MsgId      string `json:"msg_id"`      // 消息ID
		MsgType    string `json:"msg_type"`    // 消息类型
		ClientTime int64  `json:"client_time"` // 客户端接收时间
	}
)
