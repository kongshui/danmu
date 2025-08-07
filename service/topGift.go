package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kongshui/danmu/common"
)

func topGift(roomId string) bool {
	// 初始化配置文件
	headers := map[string]string{
		"Content-Type": "application/json;charset=UTF-8",
	}
	var (
		giftIds string
	)
	for k := range giftIdToName {
		giftIds += k + ","
	}
	giftIds = strings.TrimRight(giftIds, ",")
	data := map[string]any{
		"giftList":       giftIds, //置顶礼物列表
		"giftExtendInfo": giftExtendInfo(),
	}
	jsonData, _ := json.Marshal(data)
	urlPath := KsUrlSet(url_TopGiftUrl)
	if urlPath == "" {
		ziLog.Error("TopGift err, urlPath is nil ", debug)
		return false
	}
	response, err := common.HttpRespond("POST", urlPath, kuaiShouBindBodyToByte(roomId, "gift", "top", string(jsonData)), headers)
	if err != nil {
		ziLog.Error(fmt.Sprintf("TopGift response err: %v", err), debug)
		return false
	}
	defer response.Body.Close()
	var (
		request any
	)

	if err := json.NewDecoder(response.Body).Decode(&request); err != nil {
		ziLog.Error(fmt.Sprintf("TopGift json.NewDecoder err: %v", err), debug)
		return false
	}
	if response.StatusCode != 200 {
		return false
	}
	if int64(request.(map[string]any)["result"].(float64)) != 1 {
		ziLog.Error(fmt.Sprintf("TopGift err, data: %v", request), debug)
		return false
	}
	return true
}

func giftExtendInfo() string {
	batchInfo := []map[string]any{
		{
			"count":            1,
			"action":           " ",
			"objectBeingActed": " ",
		},
		{
			"count":            10,
			"action":           " ",
			"objectBeingActed": " ",
		},
		{
			"count":            30,
			"action":           " ",
			"objectBeingActed": " ",
		},
		{
			"count":            50,
			"action":           " ",
			"objectBeingActed": " ",
		},
		{
			"count":            100,
			"action":           " ",
			"objectBeingActed": " ",
		},
	}
	batchNaiInfo := []map[string]any{
		{
			"count":            1,
			"action":           " ",
			"objectBeingActed": " ",
		},
		{
			"count":            5,
			"action":           " ",
			"objectBeingActed": " ",
		},
		{
			"count":            10,
			"action":           " ",
			"objectBeingActed": " ",
		},
		{
			"count":            30,
			"action":           " ",
			"objectBeingActed": " ",
		},
		{
			"count":            50,
			"action":           " ",
			"objectBeingActed": " ",
		},
		{
			"count":            100,
			"action":           " ",
			"objectBeingActed": " ",
		},
	}
	batchjiInfo := []map[string]any{
		{
			"count":            1,
			"action":           " ",
			"objectBeingActed": " ",
		},
		{
			"count":            7,
			"action":           " ",
			"objectBeingActed": " ",
		},
		{
			"count":            10,
			"action":           " ",
			"objectBeingActed": " ",
		},
		{
			"count":            30,
			"action":           " ",
			"objectBeingActed": " ",
		},
		{
			"count":            50,
			"action":           " ",
			"objectBeingActed": " ",
		},
		{
			"count":            100,
			"action":           " ",
			"objectBeingActed": " ",
		},
	}
	batchjiaInfo := []map[string]any{
		{
			"count":            1,
			"action":           " ",
			"objectBeingActed": " ",
		},
		{
			"count":            4,
			"action":           " ",
			"objectBeingActed": " ",
		},
		{
			"count":            10,
			"action":           " ",
			"objectBeingActed": " ",
		},
		{
			"count":            30,
			"action":           " ",
			"objectBeingActed": " ",
		},
		{
			"count":            50,
			"action":           " ",
			"objectBeingActed": " ",
		},
		{
			"count":            100,
			"action":           " ",
			"objectBeingActed": " ",
		},
	}
	giftExtend := []map[string]any{
		{
			"id":             11582,    // 礼物id
			"customGiftName": "永久腕力",   //自定义礼物名称，同一个礼物id对应相同的自定义礼物名称
			"effectDesc":     "永久增加腕力", // effectDesc为礼物面板左上角说明
			"batchBar": map[string]any{
				"subTitle":   "腕力",      //用于表示该礼物效果类型，选填字段，最多4个字
				"batchInfos": batchInfo, // 批量效果信息
			},
		},
		{
			"id":             11584,  // 礼物id
			"customGiftName": "随机惊喜", //自定义礼物名称，同一个礼物id对应相同的自定义礼物名称
			"effectDesc":     "随机礼物", // effectDesc为礼物面板左上角说明
			"batchBar": map[string]any{
				"subTitle":   "随机", //用于表示该礼物效果类型，选填字段，最多4个字
				"batchInfos": batchInfo,
			},
		},
		{
			"id":             12252,   // 礼物id
			"customGiftName": "健身卡",   //自定义礼物名称，同一个礼物id对应相同的自定义礼物名称
			"effectDesc":     "召唤健身卡", // effectDesc为礼物面板左上角说明
			"batchBar": map[string]any{
				"subTitle":   "健身卡", //用于表示该礼物效果类型，选填字段，最多4个字
				"batchInfos": batchInfo,
			},
		},
		{
			"id":             11606,  // 礼物id
			"customGiftName": "啵啵奶茶", //自定义礼物名称，同一个礼物id对应相同的自定义礼物名称
			"effectDesc":     "奶茶",   // effectDesc为礼物面板左上角说明
			"batchBar": map[string]any{
				"subTitle":   "奶茶", //用于表示该礼物效果类型，选填字段，最多4个字
				"batchInfos": batchNaiInfo,
			},
		},
		{
			"id":             11585,   // 礼物id
			"customGiftName": "机械臂",   //自定义礼物名称，同一个礼物id对应相同的自定义礼物名称
			"effectDesc":     "召唤机械臂", // effectDesc为礼物面板左上角说明
			"batchBar": map[string]any{
				"subTitle":   "机械臂", //用于表示该礼物效果类型，选填字段，最多4个字
				"batchInfos": batchjiInfo,
			},
		},
		{
			"id":             11586,   // 礼物id
			"customGiftName": "啦啦队",   //自定义礼物名称，同一个礼物id对应相同的自定义礼物名称
			"effectDesc":     "召唤啦啦队", // effectDesc为礼物面板左上角说明
			"batchBar": map[string]any{
				"subTitle":   "啦啦队", //用于表示该礼物效果类型，选填字段，最多4个字
				"batchInfos": batchInfo,
			},
		},
		{
			"id":             11587,  // 礼物id
			"customGiftName": "爱的搭子", //自定义礼物名称，同一个礼物id对应相同的自定义礼物名称
			"effectDesc":     "召唤搭子", // effectDesc为礼物面板左上角说明
			"batchBar": map[string]any{
				"subTitle":   "搭子", //用于表示该礼物效果类型，选填字段，最多4个字
				"batchInfos": batchInfo,
			},
		},
		{
			"id":             12720,  // 礼物id
			"customGiftName": "机甲矩阵", //自定义礼物名称，同一个礼物id对应相同的自定义礼物名称
			"effectDesc":     "召唤机甲", // effectDesc为礼物面板左上角说明
			"batchBar": map[string]any{
				"subTitle":   "机甲", //用于表示该礼物效果类型，选填字段，最多4个字
				"batchInfos": batchjiaInfo,
			},
		},
	}
	jsonData, _ := json.Marshal(giftExtend)
	return string(jsonData)
}
