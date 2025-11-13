package service

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 获取指定用户排名
func GetPlayerTopHandle(c *gin.Context) {
	type playerTop struct {
		StartIndex int64  `json:"start_index"`
		EndIndex   int64  `json:"end_index"`
		TestCode   string `json:"test_code"`
		Revrse     bool   `json:"reverse"` // 是否为倒序，默认正序
	}
	var (
		pt playerTop
	)
	bodyByte := bytePool.Get().(*[]byte)
	defer bytePool.Put(bodyByte)
	*bodyByte, _ = c.GetRawData()
	if err := json.Unmarshal(*bodyByte, &pt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errcode": 1,
			"errmsg":  "json unmarshal failed",
		})
		return
	}
	if !compareTestCode(pt.TestCode) {
		c.JSON(http.StatusBadRequest, gin.H{
			"errcode": 4,
			"errmsg":  "测试验证码错误",
		})
		return
	}
	// 校验参数
	if pt.StartIndex < 0 || (pt.EndIndex != -1 && pt.EndIndex < pt.StartIndex) {
		c.JSON(http.StatusBadRequest, gin.H{
			"errcode": 2,
			"errmsg":  "start_index or end_index is invalid",
		})
		return
	}
	// 获取排行榜数据和总长度
	data := getTopWorldRankData(pt.StartIndex, pt.EndIndex, pt.Revrse)
	total, _ := getTop100RankLen()

	c.JSON(200, gin.H{
		"errcode": 0,
		"errmsg":  "success",
		"data":    data,
		"total":   total,
	})
}

// 更改玩家积分
func ChangePlayerTopHandle(c *gin.Context) {
	type changePlayerTop struct {
		OpenId   string `json:"open_id"`
		Score    int64  `json:"score"`
		TestCode string `json:"test_code"`
	}
	var (
		ct changePlayerTop
	)
	bodyByte := bytePool.Get().(*[]byte)
	defer bytePool.Put(bodyByte)
	*bodyByte, _ = c.GetRawData()
	if err := json.Unmarshal(*bodyByte, &ct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errcode": 1,
			"errmsg":  "json unmarshal failed",
		})
		return
	}
	if !compareTestCode(ct.TestCode) {
		c.JSON(http.StatusBadRequest, gin.H{
			"errcode": 4,
			"errmsg":  "测试验证码错误",
		})
		return
	}
	// 校验参数
	if ct.OpenId == "" || ct.Score <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"errcode": 2,
			"errmsg":  "open_id or score is invalid",
		})
		return
	}
	if ct.Score <= 0 {
		rdb.ZRem(world_rank_week, ct.OpenId)
	} else {
		if err := rdb.ZAdd(world_rank_week, float64(ct.Score), ct.OpenId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errcode": 3,
				"errmsg":  "zincrby failed",
			})
			return
		}
	}
	// 如果积分小于等于0，从排行榜中移除

	c.JSON(200, gin.H{
		"errcode": 0,
		"errmsg":  "success",
	})
}

// 通过openId获取玩家信息
func GetPlayerInfoByOpenIdHandle(c *gin.Context) {
	type playerInfoByOpenId struct {
		OpenId   string `json:"open_id"`
		TestCode string `json:"test_code"`
	}
	var (
		pio playerInfoByOpenId
	)
	bodyByte := bytePool.Get().(*[]byte)
	defer bytePool.Put(bodyByte)
	*bodyByte, _ = c.GetRawData()
	if err := json.Unmarshal(*bodyByte, &pio); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errcode": 1,
			"errmsg":  "json unmarshal failed",
		})
		return
	}
	if !compareTestCode(pio.TestCode) {
		c.JSON(http.StatusBadRequest, gin.H{
			"errcode": 4,
			"errmsg":  "测试验证码错误",
		})
		return
	}
	// 校验参数
	if pio.OpenId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"errcode": 2,
			"errmsg":  "open_id is empty",
		})
		return
	}
	// 获取玩家信息
	setUserInfos := make([]map[string]any, 1)
	score, rank, _ := GetPlayerWorldRankData(pio.OpenId)
	userInfo, _ := UserInfoGet(pio.OpenId, false)
	userInfos := map[string]any{
		"open_id":    pio.OpenId,
		"rank":       rank,
		"score":      score,
		"avatar_url": userInfo.AvatarUrl,
		"nick_name":  userInfo.NickName,
	}
	setUserInfos[0] = userInfos
	c.JSON(200, gin.H{
		"errcode": 0,
		"errmsg":  "success",
		"data":    setUserInfos,
		"total":   0,
	})
}
