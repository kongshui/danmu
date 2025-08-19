package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"path"
	"strconv"
	"time"

	"github.com/kongshui/danmu/model/pmsg"
	"github.com/kongshui/danmu/sse"

	"github.com/kongshui/danmu/common"

	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func HandleAddcoin(c *gin.Context) {
	// err := scrollWinningStreamCoin()
	type AddCoin struct {
		OpenId string `json:"open_id"`
		Coin   int64  `json:"coin"`
	}
	var addCoin AddCoin
	if err := c.ShouldBindJSON(&addCoin); err != nil {
		log.Println(err)
		c.JSON(404, gin.H{
			"err": err,
		})
		return
	}
	coin, err := AddUserWinStreamCoin(addCoin.OpenId, addCoin.Coin)
	if err != nil {
		log.Println(err)
		c.JSON(404, gin.H{
			"err": err,
		})
		return
	}
	c.JSON(200, gin.H{
		"coin": coin,
	})
}

func HandleAddUser(c *gin.Context) {
	var user UserInfoStruct
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println(err)
		c.JSON(404, gin.H{
			"err": err,
		})
		return
	}
	if err := userInfoStore(user); err != nil {
		log.Println(err)
		c.JSON(404, gin.H{
			"err": err,
		})
		return
	}
	c.JSON(200, "添加成功")
}

// 添加世界排行版
func HandleAddWorldRank(c *gin.Context) {
	type AddCoin struct {
		OpenId string `json:"open_id"`
		Score  int64  `json:"score"`
	}
	var addCoin AddCoin
	if err := c.ShouldBindJSON(&addCoin); err != nil {
		log.Println(err)
		c.JSON(404, gin.H{
			"err": err,
		})
		return
	}
	if err := WorldRankNumerAdd(addCoin.OpenId, float64(addCoin.Score)); err != nil {
		log.Println(err)
		c.JSON(404, gin.H{
			"err": err,
		})
		return
	}
	c.JSON(200, "添加成功")
}

// 添加连胜次数
func HandlestreamCount(c *gin.Context) {
	type AddCoin struct {
		OpenId string `json:"open_id"`
		Stats  int    `json:"stats"`
	}
	var addCoin AddCoin
	if err := c.ShouldBindJSON(&addCoin); err != nil {
		log.Println(err)
		c.JSON(404, gin.H{
			"err": err,
		})
		return
	}
	if err := winningStreamCountAdd(addCoin.Stats, addCoin.OpenId); err != nil {
		log.Println(err)
		c.JSON(404, gin.H{
			"err": err,
		})
		return
	}
	c.JSON(200, "添加成功")
}

// 发送虚拟评论
func HandleSendFakeComment(c *gin.Context) {
	type comment struct {
		AnchorOpenId string `json:"anchor_open_id"`
		Comment      string `json:"comment"`
		UserId       string `json:"user_id"`
		UserName     string `json:"user_name"`
		UserAvatar   string `json:"user_avatar"`
		TestCode     string `json:"test_code"`
	}
	var commentGet comment
	if err := c.ShouldBindJSON(&commentGet); err != nil {
		log.Println(err)
		c.JSON(404, gin.H{
			"err": err,
		})
		return
	}
	testCode := getTestCode()
	if testCode == "" {
		c.JSON(403, gin.H{
			"err": "获取测试验证码错误",
		})
		return
	}
	if commentGet.TestCode != testCode {
		c.JSON(403, gin.H{
			"err": "测试验证码错误",
		})
		return
	}

	roomId := QueryRoomIdInterconvertAnchorOpenId(commentGet.AnchorOpenId)
	if roomId == "" {
		c.JSON(404, gin.H{
			"err": "该主播未开播",
		})
		return
	}

	switch platform {
	case "ks":
		data := KsCallbackDataStruct{}
		commentData := KsLiveCommentStruct{}
		data.UniqueMessageId = strconv.FormatInt(time.Now().UnixMilli(), 10)
		data.PushType = "comment"
		data.RoomCode = roomId
		data.AuthorOpenId = commentGet.AnchorOpenId
		commentData.Content = commentGet.Comment
		commentData.UserInfo.NickName = commentGet.UserName
		commentData.UserInfo.UserId = commentGet.UserId
		commentData.UserInfo.AvatarUrl = commentGet.UserAvatar

		// switch commentGet.UserId {
		// case 0:
		// 	commentData.UserInfo.NickName = "cccc"
		// 	commentData.UserInfo.UserId = "1234567890565598"
		// 	commentData.UserInfo.AvatarUrl = "https://www.keaitupian.cn/cjpic/frombd/2/253/2107631312/3178897554.jpg"
		// case 1:
		// 	commentData.UserInfo.NickName = "dddd"
		// 	commentData.UserInfo.UserId = "1234567890565599"
		// 	commentData.UserInfo.AvatarUrl = "https://ts1.tc.mm.bing.net/th/id/OIP-C.mH9YLFEL5YdVxJM82mjVJQHaEo?w=280&h=211&c=8&rs=1&qlt=90&r=0&o=6&pid=3.1&rm=2"
		// }
		// fmt.Println(commentData.UserInfo)
		data.Payload = append(data.Payload, commentData)
		ksPushCommentPayloay(data)
	case "dy":
		data := []ContentPayloadStruct{}
		subData := ContentPayloadStruct{
			MsgId:     strconv.FormatInt(time.Now().UnixMilli(), 10),
			SecOpenid: commentGet.UserId,
			Content:   commentGet.Comment,
			AvatarUrl: commentGet.UserAvatar,
			Nickname:  commentGet.UserName,
			TimeStamp: time.Now().UnixMilli(),
		}
		data = append(data, subData)
		// switch commentGet.UserId {
		// case 0:
		// 	subData := ContentPayloadStruct{
		// 		MsgId:     strconv.FormatInt(time.Now().UnixMilli(), 10),
		// 		SecOpenid: "test123456",
		// 		Content:   commentGet.Comment,
		// 		AvatarUrl: "https://www.keaitupian.cn/cjpic/frombd/2/253/2107631312/3178897554.jpg",
		// 		Nickname:  "cccc",
		// 		TimeStamp: time.Now().UnixMilli(),
		// 	}
		// 	data = append(data, subData)
		// case 1:
		// 	subData := ContentPayloadStruct{
		// 		MsgId:     strconv.FormatInt(time.Now().UnixMilli(), 10),
		// 		SecOpenid: "test987654",
		// 		Content:   commentGet.Comment,
		// 		AvatarUrl: "https://ww2.sinaimg.cn/mw690/007ut4Uhly1hx4v37mpxcj30u017cgrv.jpg",
		// 		Nickname:  "dddd",
		// 		TimeStamp: time.Now().UnixMilli(),
		// 	}
		// 	data = append(data, subData)
		// }

		dataByte, _ := json.Marshal(data)
		pushDyBasePayloayDirect(roomId, commentGet.AnchorOpenId, "live_comment", dataByte)
	}
	c.JSON(200, "添加成功")
}

// 发送虚假礼物
func HandleSendFakeGift(c *gin.Context) {
	type gift struct {
		AnchorOpenId string `json:"anchor_open_id"`
		GiftId       string `json:"gift_id"`
		Num          int64  `json:"num"`
		UserId       string `json:"user_id"`
		UserName     string `json:"user_name"`
		UserAvatar   string `json:"user_avatar"`
		TestCode     string `json:"test_code"`
	}
	var giftGet gift
	if err := c.ShouldBindJSON(&giftGet); err != nil {
		log.Println(err)
		c.JSON(404, gin.H{
			"err": err,
		})
		return
	}
	testCode := getTestCode()
	if testCode == "" {
		c.JSON(403, gin.H{
			"err": "获取测试验证码错误",
		})
		return
	}
	if giftGet.TestCode != testCode {
		c.JSON(403, gin.H{
			"err": "测试验证码错误",
		})
		return
	}
	roomId := QueryRoomIdInterconvertAnchorOpenId(giftGet.AnchorOpenId)
	if roomId == "" {
		c.JSON(404, gin.H{
			"err": "该主播未开播",
		})
		return
	}
	switch platform {
	case "ks":
		data := KsCallbackQueryStruct{}
		giftData := KsGiftSendStruct{}
		data.AuthorOpenId = giftGet.AnchorOpenId
		data.RoomCode = roomId
		data.PushType = "giftSend"
		data.UniqueMessageId = strconv.FormatInt(time.Now().UnixMilli(), 10)
		data.LiveTimeStamp = time.Now().UnixMilli()
		giftData.UniqueNo = strconv.FormatInt(time.Now().Unix(), 10)
		giftData.GiftId = giftGet.GiftId
		giftData.GiftCount = giftGet.Num
		giftData.GiftName = giftIdToName[giftGet.GiftId]
		giftData.GiftUnitPrice = 1
		giftData.GiftTotalPrice = giftGet.Num
		// switch giftGet.UserId {
		// case 0:
		// 	giftData.UserInfo.NickName = "cccc"
		// 	giftData.UserInfo.UserId = "1234567890565598"
		// 	giftData.UserInfo.AvatarUrl = "https://www.keaitupian.cn/cjpic/frombd/2/253/2107631312/3178897554.jpg"
		// case 1:
		// 	giftData.UserInfo.NickName = "dddd"
		// 	giftData.UserInfo.UserId = "1234567890565599"
		// 	giftData.UserInfo.AvatarUrl = "https://ts1.tc.mm.bing.net/th/id/OIP-C.mH9YLFEL5YdVxJM82mjVJQHaEo?w=280&h=211&c=8&rs=1&qlt=90&r=0&o=6&pid=3.1&rm=2"
		// }
		giftData.UserInfo.NickName = giftGet.UserName
		giftData.UserInfo.UserId = giftGet.UserId
		giftData.UserInfo.AvatarUrl = giftGet.UserAvatar

		data.Payload = append(data.Payload, giftData)
		ksPushGiftSendPayloay(data)
	case "dy":
		data := []GiftPayloadStruct{}
		// switch giftGet.UserId {
		// case 0:
		// 	subData := GiftPayloadStruct{
		// 		MsgId:     strconv.FormatInt(time.Now().UnixMilli(), 10),
		// 		SecOpenid: "test123456",
		// 		SecGiftId: giftGet.GiftId,
		// 		GiftNum:   int(giftGet.Num),
		// 		AvatarUrl: "https://www.keaitupian.cn/cjpic/frombd/2/253/2107631312/3178897554.jpg",
		// 		Nickname:  "cccc",
		// 		TimeStamp: time.Now().UnixMilli(),
		// 	}
		// 	data = append(data, subData)
		// case 1:
		// 	subData := GiftPayloadStruct{
		// 		MsgId:     strconv.FormatInt(time.Now().UnixMilli(), 10),
		// 		SecOpenid: "test987654",
		// 		SecGiftId: giftGet.GiftId,
		// 		GiftNum:   int(giftGet.Num),
		// 		AvatarUrl: "https://ww2.sinaimg.cn/mw690/007ut4Uhly1hx4v37mpxcj30u017cgrv.jpg",
		// 		Nickname:  "dddd",
		// 		TimeStamp: time.Now().UnixMilli(),
		// 	}
		// 	data = append(data, subData)
		// }
		subData := GiftPayloadStruct{
			MsgId:     strconv.FormatInt(time.Now().UnixMilli(), 10),
			SecOpenid: giftGet.UserId,
			SecGiftId: giftGet.GiftId,
			GiftNum:   int(giftGet.Num),
			AvatarUrl: giftGet.UserAvatar,
			Nickname:  giftGet.UserName,
			TimeStamp: time.Now().UnixMilli(),
		}
		data = append(data, subData)

		dataByte, _ := json.Marshal(data)
		pushDyBasePayloayDirect(roomId, giftGet.AnchorOpenId, "live_gift", dataByte)
	}
	c.JSON(200, "添加成功")
}

// 发送虚假点赞
func HandleSendLiveLike(c *gin.Context) {
	type like struct {
		AnchorOpenId string `json:"anchor_open_id"`
		Num          int64  `json:"num"`
		UserId       string `json:"user_id"`
		UserName     string `json:"user_name"`
		UserAvatar   string `json:"user_avatar"`
		TestCode     string `json:"test_code"`
	}
	var likeGet like
	if err := c.ShouldBindJSON(&likeGet); err != nil {
		log.Println(err)
		c.JSON(404, gin.H{
			"err": err,
		})
		return
	}
	testCode := getTestCode()
	if testCode == "" {
		c.JSON(403, gin.H{
			"err": "获取测试验证码错误",
		})
		return
	}
	if likeGet.TestCode != testCode {
		c.JSON(403, gin.H{
			"err": "测试验证码错误",
		})
		return
	}
	roomId := QueryRoomIdInterconvertAnchorOpenId(likeGet.AnchorOpenId)
	if roomId == "" {
		c.JSON(404, gin.H{
			"err": "该主播未开播",
		})
		return
	}
	switch platform {
	case "ks":
		data := KsCallbackDataStruct{}
		likeData := KsLiveLikeStruct{}
		data.AuthorOpenId = likeGet.AnchorOpenId
		data.RoomCode = roomId
		data.PushType = "liveLike"
		data.UniqueMessageId = strconv.FormatInt(time.Now().UnixMilli(), 10)
		likeData.Count = likeGet.Num
		likeData.UserInfo.UserId = likeGet.UserId
		likeData.UserInfo.NickName = likeGet.UserName
		likeData.UserInfo.AvatarUrl = likeGet.UserAvatar
		data.Payload = append(data.Payload, likeData)
		ksPushLiveLikePayloay(data)
	case "dy":
		data := []LiveLikePayloadStruct{}
		subData := LiveLikePayloadStruct{
			MsgId:     strconv.FormatInt(time.Now().UnixMilli(), 10),
			SecOpenid: likeGet.UserId,
			LikeNum:   likeGet.Num,
			AvatarUrl: likeGet.UserAvatar,
			Nickname:  likeGet.UserName,
			TimeStamp: time.Now().UnixMilli(),
		}
		data = append(data, subData)

		dataByte, _ := json.Marshal(data)
		pushDyBasePayloayDirect(roomId, likeGet.AnchorOpenId, "live_like", dataByte)
	}
	c.JSON(200, "添加成功")
}

// 玩家加入组
func HandleAddGroup(c *gin.Context) {
	type AddGroup struct {
		RoomId  string            `json:"room_id"`
		OpenId  string            `json:"open_id"`
		RoundId int64             `json:"round_id"`
		UserMap []JoinGroupStruct `json:"user_map"`
	}
	var addGroup AddGroup
	if err := c.ShouldBindJSON(&addGroup); err != nil {
		log.Println(err)
		c.JSON(404, gin.H{
			"err": err,
		})
		return
	}
	data := make([]*pmsg.SingleRoomAddGroupInfo, 0)
	for _, v := range addGroup.UserMap {
		data = append(data, &pmsg.SingleRoomAddGroupInfo{OpenId: v.OpenId, GroupId: v.GroupId, AvatarUrl: v.AvatarUrl, NickName: v.NickName})
	}
	liveCurrentRoundAdd(addGroup.RoomId, addGroup.RoundId)
	err := playerGroupAdd(addGroup.RoomId, addGroup.OpenId, data, false)
	c.JSON(200, err)
}

// 获取所有直播间信息
func HandleGetAllRoom(c *gin.Context) {
	uidList, err := etcdClient.Client.Get(first_ctx, path.Join("/", config.Project, common.Uid_Register_RoomId_key), clientv3.WithPrefix())
	if err != nil {
		ziLog.Error(fmt.Sprintf("ksCallBackQueryToKs 查询直播房间号失败:  %v", err), debug)
		c.JSON(200, gin.H{
			"err": err,
		})
	}
	data := make(map[string]UserInfoStruct)
	for _, kv := range uidList.Kvs {
		roomId := string(kv.Value)
		anchorOpenId := QueryRoomIdInterconvertAnchorOpenId(roomId)
		userInfo, err := userInfoGet(anchorOpenId)
		if err != nil {
			log.Println("通过接口获取玩家信息失败", anchorOpenId, "err:", err)
		}
		data[roomId] = userInfo
	}
	c.JSON(200, data)
}

// 随机赠送礼物
func HandleSendFakeRandomGift(c *gin.Context) {
	groupId := c.Query("group_id")
	roomId := c.Query("room_id")
	// 1.查询现在直播房间Id

	if roomId == "" {
		roomIdList, err := rdb.SMembers(room_id_list_db)
		if err != nil && len(roomIdList) == 0 {
			log.Println(err)
			c.JSON(404, gin.H{
				"err": err,
			})
			return
		}
		roomId = roomIdList[len(roomIdList)-1]
	}
	// 2. 通过roomid查询直播roundid
	roundId, ok := queryRoomIdToRoundId(roomId)
	if !ok {
		c.JSON(404, gin.H{
			"err": ok,
		})
		return
	}
	// 3. 获取用户列表
	groupName := roomId + "_" + strconv.FormatInt(roundId, 10) + "_group"
	userGroupMap, err := rdb.HGetAll(groupName)
	if err != nil {
		log.Println(err)
		c.JSON(404, gin.H{
			"err": err,
		})
		return
	}
	var userList []string
	for k, v := range userGroupMap {
		log.Println("roomId和分组", k, v)
		if v == groupId {
			userList = append(userList, k)
		}
	}
	log.Println("userList:", userList)
	anchorOpenId := QueryRoomIdInterconvertAnchorOpenId(roomId)
	if len(userList) == 0 {
		var (
			context string
		)
		switch groupId {
		case "Left":
			context = "111"
		case "Right":
			context = "222"
		}
		var newOpenId string = "74"
		for range 17 {
			newOpenId += strconv.Itoa(rand.Intn(10))
		}
		// fakeSendMessage(roomId, anchorOpenId, "live_comment", createCommon(newOpenId, context))
		userList = append(userList, newOpenId)
		fmt.Println(createCommon(newOpenId, context))
	}
	giftIdList := make([]string, 0)
	// 4. 随机赠送礼物
	for k := range giftToScoreMap {
		giftIdList = append(giftIdList, k)
	}
	score := 0
	for score <= 2500 {
		rand.New(rand.NewSource(time.Now().UnixNano()))
		giftId := giftIdList[rand.Intn(len(giftIdList))]
		giftString := createGiftstring(userList[rand.Intn(len(userList))], giftId, 1)
		// fakeSendMessage(roomId, anchorOpenId, "live_gift", giftString)
		score += int(giftToScoreMap[giftId])
		fmt.Println(giftString)
	}
	fmt.Println(anchorOpenId)
	c.JSON(200, "发送成功")
}

func createGiftstring(openId, giftId string, number int) string {
	giftString := "{\"avatar_url\":\"https://p26.douyinpic.com/aweme/100x100/aweme-avatar/mosaic-legacy_3795_3033762272.jpeg?from=3067671334\",\"gift_num\":" + strconv.Itoa(number) + ",\"gift_value\":10,\"msg_id\":\"7482687722940601396\",\"nickname\":\"空空空\",\"sec_gift_id\":\"" + giftId + "\",\"sec_openid\":\"" + openId + "\",\"timestamp\":" + strconv.FormatInt(time.Now().UnixMilli(), 10) + "}"
	return giftString
}

func createCommon(openId, context string) string {
	commonString := "{\"avatar_url\":\"https://p3-developer-sign.bytemaimg.com/tos-cn-i-ke512zj2cu/fab2812ad6674d3d9feb87be98dc0a17~tplv-noop.jpeg?rk3s=3839646d\\u0026x-expires=1742867925\\u0026x-signature=ZuBUTcDFTV97gqHlFNS3H6VJP64%3D\",\"content\":\"" + context + "\",\"msg_id\":\"7482962512066384948\",\"nickname\":\"aaaa\",\"sec_openid\":\"" + openId + "\",\"timestamp\":" + strconv.FormatInt(time.Now().UnixMilli(), 10) + "}"
	return commonString
}

// 接收消息测试
func ReceiveMessageHandle(c *gin.Context) {
	//访问后端web服务器
	data, _ := c.GetRawData()
	if err := sse.SseSend(pmsg.MessageId_TestMsgAck, []string{c.GetHeader("x-client-uuid")}, data); err != nil {
		log.Println("访问前端web服务器失败", err)
	}
	c.JSON(200, gin.H{
		"err_code": 0,
		"err_msg":  "",
	})
}

// 获取测试验证码
func getTestCode() string {
	code, err := rdb.Get("test_code")
	if err != nil {
		log.Println("获取测试验证码失败", err)
		return ""
	}
	return code
}

// 获取所有主播
func GetAllAnchor(c *gin.Context) {
	req, err := etcdClient.Client.Get(context.Background(), path.Join("/"+config.Project+common.Uid_Register_OpenId_key), clientv3.WithPrefix())

	if err != nil {
		log.Println("获取所有主播失败", err)
		c.JSON(404, gin.H{
			"err": err,
		})
		return
	}
	if c.Query("verifycode") != getTestCode() {
		c.JSON(404, gin.H{
			"err": "验证码错误",
		})
		return
	}
	var anchorList []UserInfoStruct

	for _, resp := range req.Kvs {
		log.Println(string(resp.Key), string(resp.Value))
		userinfo, _ := userInfoGet(string(resp.Value))
		anchorList = append(anchorList, userinfo)
	}
	c.JSON(200, anchorList)
}
