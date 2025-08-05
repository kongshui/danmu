package battlematch

const (
	match_battle_v1_lock         = "/match/lock/v1"          // etcd聊天匹配锁路径
	match_battle_group_v1_lock   = "/match/lock/group/"      // etcd聊天匹配组锁路径
	match_battle_roundid_v1      = "/match/roundid/v1/"      // etcd聊天匹配组轮次路径
	matcd_battle_store_v1        = "/match/battle/v1/"       // etcd战斗匹配路径
	matcd_battle_register_v1     = "/match/register/v1/"     // etcd战斗匹配key
	matcd_battle_num_register_v1 = "/match/register/num/v1/" // etcd战斗匹配key
	matcd_battle_disconnect      = "/match/disconnect/v1/"   // etcd战斗匹掉线
	match_battle_cancel_v1       = "/match/cancel/"          // 战斗匹配取消
	match_battle_group_status    = "/match/status/"          // 战斗匹配状态
	match_battle_group_time      = "/match/time/"            // 战斗匹配时间
	match_success_timeout        = 8000                      // 匹配成功超时时间 2小时15分钟
	// match_start_timeout               = 60                            // 匹配开始超时时间 60秒
	// match_cancel_timeout              = 60                            // 匹配取消超时时间 60秒
	match_register_timeout            = 120                           // 匹配注册超时时间 60秒
	match_battle_anonymous_status     = "/match/anonymous/v1/"        // 匹配匿名存储
	match_battle_divideup_coin_status = "/match/divideup/coin/status" // 瓜分积分状态
	match_battle_normal_coin_status   = "/match/normal/coin/status"   // 获取正常积分
	match_battle_group_end            = "/match/group/end/"           // 匹配组结束
)
