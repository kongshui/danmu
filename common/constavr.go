package common

const (
	// Match_group_prefix = "match_group_" //匹配组前缀
	Uuid_Online_key            = "/uuid_online_key/"            // client uuid在线状态
	Uuid_Register_key          = "/uuid_register_key/"          // uuid注册,clientuuid和gateway uuid注册
	Node_Register_key          = "/node_register_key/"          //节点注册，node节点注册key
	RoomId_Register_Uid_key    = "/roomid_register_uid_key/"    // roomId注册，roomId和client Uuid信息
	Uid_Register_RoomId_key    = "/uid_register_roomid_key/"    // roomId注册，roomId和client Uuid信息
	GroupId_Register_key       = "/groupid_register_key/"       // groupId注册，注册groupId和client Uuid
	UserId_Register_Uid_key    = "/userid_register_uid_key/"    // userId注册，注册userId和client Uuid
	Uid_Register_UserId_key    = "/uid_register_userid_key/"    // userId注册，注册userId和client Uuid
	OpenId_Register_Uid_key    = "/openid_register_uid_key/"    // openId注册，注册openId和client Uuid
	Uid_Register_OpenId_key    = "/uid_register_openid_key/"    // openId注册，注册openId和client Uuid
	RoomInfo_Register_key      = "/roominfo_register_key/"      // 直播间信息注册，注册直播间信息和client Uuid
	Group_UserId_Register_key  = "/group_userid_register_key/"  // 群成员注册，注册群成员和client Uuid
	Udp_Register_key           = "/udp_register_key/"           // udp注册，注册udp和client Uuid
	RoomId_OpenId_Register_key = "/roomid_openid_register_key/" // roomId和openId注册，注册roomId和openId信息
	// Token_Register_key        = "/token_register_key"        // token注册，注册token和client Uuid
)
