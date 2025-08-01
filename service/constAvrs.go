package service

import (
	"time"
)

const (
	winning_streak_coin_db        = "winning_stream_coin_store"  //存储玩家连胜币
	world_rank_version_db         = "world_rank_version"         //世界排行版版本信息
	world_rank_version_list_db    = "world_rank_version_list"    //世界排行版版本列表
	world_rank_historical_db      = "world_rank_historical"      //历史排行版
	winning_streak_count_db       = "winning_streak_count"       //连胜统计
	world_rank_week               = "world_rank_week"            //周排行版
	group_list_db                 = "group_list"                 //分组列表
	integral_pool_Prefix          = "integral_pool_"             //积分池前缀
	room_id_list_db               = "roomid_list"                //房间列表
	user_info_db                  = "user_info"                  //用户信息
	access_token_db               = "access_token"               //全局token
	monitor_world_push_db         = "monitor_world_push"         //世界排行版推送锁
	monitor_world_history_push_db = "monitor_world_history_push" //历史世界榜单推送锁
	roomdid_to_anchoropenid_db    = "roomd_id_to_anchor_openid"  //房间id对应主播id
	disconnect_roomid_db          = "disconnect_roomid"          //断线房间id列表
	monitor_fail_message_push_db  = "monitor_fail_message_push"  //失败消息推送锁
	monitor_version_scroll_db     = "monitor_version_scroll"     //版本滚动推送锁
	monitor_top_100_ranking_db    = "monitor_top_100_ranking"    //世界排行版前100名推送锁
	monitor_access_token_db       = "monitor_access_token"       //全局token推送锁
	monitor_disconnect_roomid_db  = "monitor_disconnect_roomid"  //断线房间id列表推送锁
	version_time_layout           = "20060102"                   //版本号时间格式
	mysql_query_time_layout       = "2006-01-02"                 // 数据库版本查询时间格式
	version_time_interval         = 259200                       //时间间隔
	top_100_ranking               = "top_100_ranking"            //世界排行版前100名
	is_consume_db                 = "is_consume"                 //是否是第一次消费
	forward_domain_key            = "/forward_domain"            //前端http节点
	backend_domain_key            = "/backend_domain"            //后端http节点
	grpc_domain_key               = "/grpc_domain"               //grpc节点
	forward_domain_uri            = "/ws"                        //前端uri
	match_battle_status_ready     = 1                            // 匹配准备状态
	match_battle_status_set_time  = 2                            // 匹配设置时间状态
	match_battle_status_Confirm   = 3                            // 匹配确认状态
	match_battle_status_start     = 4                            // 匹配开始状态
	match_battle_status_stop      = 5                            // 匹配结束状态
	match_battle_group_time       = "/match/time/"               // 战斗匹配时间
	group_integral_pool_key       = "groupIntegralPool"          // pk积分池
	scroll_time_hours             = 20                           // 清榜时间小时,晚8点
	scroll_time_week              = time.Thursday                // 清榜时间周
	live_like_score               = 1                            // 点赞或者666的积分
	// 抖音相关接口
	url_start_push_url                 = "https://webcast.bytedance.com/api/live_data/task/start"                              // 开始推送地址
	url_stop_push_url                  = "https://webcast.bytedance.com/api/live_data/task/stop"                               // 停止推送地址
	url_check_push_url                 = "https://webcast.bytedance.com/api/live_data/task/get"                                // 查询推送状态地址
	url_get_anchor_info_url            = "https://webcast.bytedance.com/api/webcastmate/info"                                  // 获取主播信息地址
	url_round_sync_status              = "https://webcast.bytedance.com/api/gaming_con/round/sync_status"                      // 同步对局状态，开始和结束后调用
	url_upload_user_group_url          = "https://webcast.bytedance.com/api/gaming_con/round/upload_user_group_info"           // 上传玩家组信息
	url_user_world_rank_upload_url     = "https://webcast.bytedance.com/api/gaming_con/world_rank/upload_rank_list"            //上传世界排行版数据
	url_round_user_result_upload_url   = "https://webcast.bytedance.com/api/gaming_con/round/upload_user_result"               //上传对局用户数据
	url_round_user_rank_upload_url     = "https://webcast.bytedance.com/api/gaming_con/round/upload_rank_list"                 //上传对局用户排行数据
	url_round_user_upload_complete_url = "https://webcast.bytedance.com/api/gaming_con/round/complete_upload_user_result"      //上传对局用户完成数据
	url_fail_get_live_data_url         = "https://webcast.bytedance.com/api/live_data/task/fail_data/get"                      //获取失败数据地址
	url_set_world_rank_version_url     = "https://webcast.bytedance.com/api/gaming_con/world_rank/set_valid_version"           //设置世界排行版版本地址
	url_live_data_ack_url              = "https://webcast.bytedance.com/api/live_data/ack"                                     // 抖音ack地址
	url_world_rank_user_total_url      = "https://webcast.bytedance.com/api/gaming_con/world_rank/upload_user_result"          //上传世界排行版用户累计战绩
	url_complete_upload_url            = "https://webcast.bytedance.com/api/gaming_con/world_rank/complete_upload_user_result" //完成用户世界榜单的累计战绩上报
	// ks相关接口
	url_BindUrl           = "https://open.kuaishou.com/openapi/developer/live/smallPlay/bind"                  //绑定直播间的url
	url_SyncGameStatusUrl = "https://open.kuaishou.com/openapi/developer/live/smallPlay/round"                 //同步对局状态url
	url_CpReceiveAckUrl   = "https://open.kuaishou.com/openapi/developer/live/data/interactive/ack/receive"    //cp客户端收到消息ack接口
	url_CpShowAckUrl      = "https://open.kuaishou.com/openapi/developer/live/data/interactive/ack/show"       //cp客户端展示消息ack接口
	url_MessageQueryUrl   = "https://open.kuaishou.com/openapi/developer/live/data/interactive/pushdata/query" //消息回查接口
	url_TopGiftUrl        = "https://open.kuaishou.com/openapi/developer/live/smallPlay/gift"                  //礼物置顶接口
	url_InteractiveUrl    = "https://open.kuaishou.com/openapi/developer/live/data/interactive/start"          //快捷选边
	url_ChatUrl           = "https://open.kuaishou.com/openapi/developer/live/data/interactive/action/chat"    // 连线对局url
	// is_maintain_db                = "is_maintain"                              //维护开关
	// integral_pool_db               = "integral_pool"                                                                       //积分池
)
