package service

//初始化积分
func GiftToScoreInit(giftScore map[string]float64, giftToName map[string]string, commentToGift map[string]string) {

	//初始化礼物兑换积分
	giftToScoreMap = giftScore
	giftIdToName = giftToName
	commentTogiftId = commentToGift
	//仙女棒

	// //黄色仙女棒
	// giftToScoreMap["eplFUy7i0B0fiv0Iym1MpOZa5XmUE8g/WUAyJ6Tc+UJJDpcs7pzclNOz/WM="] = 3.0
	// //蓝色仙女棒
	// giftToScoreMap["4I66OIE1HKWfM7PNvAHtAgYUSNlggSEgcpo3ai8GYQXAWqjrDuH8NtjsWEQ="] = 3.0
	// //绿色仙女棒
	// giftToScoreMap["XHS+QR5Cv0b9ydsZZ5mLkrhPtMTTdrLsgfWNU/QX2IfUy2P6dmDaRHJT+0U="] = 3.0
	// //紫色仙女棒
	// giftToScoreMap["gs+95ujNzXXSCtLTv97fWgbApTQi0sqz1BULB+7w62g+v4sFxINvxOIrXCw="] = 3.0

	// //派对话筒
	// giftToScoreMap["YbLESoUj053FWVYPWUNOAtp4FYnb+/eZbyrLi7ndArVFz14rivgxf0cFrKs="] = 100.00
	// //神秘空投
	// giftToScoreMap["pGLo7HKNk1i4djkicmJXf6iWEyd+pfPBjbsHmd3WcX0Ierm2UdnRR7UINvI="] = 2500
	// //草莓甜点
	// giftToScoreMap["rHSYLYmWZvy2eal+foarBgE5qDJe3j2e96icNOhTSO/t2G0W8UVnXN3mltM="] = 100.00
	// //超级空投
	// giftToScoreMap["lsEGaeC++k/yZbzTU2ST64EukfpPENQmqEZxaK9v1+7etK+qnCRKOnDyjsE="] = 100.00
	// //生命药水
	// giftToScoreMap["CwaE5Zyq2ZvwgA4S9udpz1BTWMAM41FRkMN5WI5J8E7QJEUT6VFZK0USxig="] = 100.00
	// //超能喷射
	// giftToScoreMap["P7zDZzpeO215SpUptB+aURb1+zC14UC9MY1+MHszKoF0p5gzYk8CNEbey60="] = 100.00
	// //稀有宝箱
	// giftToScoreMap["rROiXLcY2saGvxHt3fAkYbWvbbikhEzbo0wpI794zEv+A2SCLrkNKYZEVuE="] = 100.00

	// switch platform {
	// case "ks":
	// 	// 快手设置
	// 	giftToScoreMap["11582"] = 1
	// 	//能力药丸
	// 	giftToScoreMap["12252"] = 30
	// 	//魔法镜
	// 	giftToScoreMap["11606"] = 70
	// 	//甜甜圈
	// 	giftToScoreMap["11585"] = 200
	// 	//能量电池
	// 	giftToScoreMap["11586"] = 600
	// 	//爱的爆炸
	// 	giftToScoreMap["11587"] = 1600
	// 	//魔法相机
	// 	giftToScoreMap["12720"] = 2500
	// case "dy":
	// 	//仙女棒
	// 	giftToScoreMap["n1/Dg1905sj1FyoBlQBvmbaDZFBNaKuKZH6zxHkv8Lg5x2cRfrKUTb8gzMs="] = 1
	// 	//能力药丸
	// 	giftToScoreMap["28rYzVFNyXEXFC8HI+f/WG+I7a6lfl3OyZZjUS+CVuwCgYZrPrUdytGHu0c="] = 30
	// 	//魔法镜
	// 	giftToScoreMap["fJs8HKQ0xlPRixn8JAUiL2gFRiLD9S6IFCFdvZODSnhyo9YN8q7xUuVVyZI="] = 70
	// 	//甜甜圈
	// 	giftToScoreMap["PJ0FFeaDzXUreuUBZH6Hs+b56Jh0tQjrq0bIrrlZmv13GSAL9Q1hf59fjGk="] = 200
	// 	//能量电池
	// 	giftToScoreMap["IkkadLfz7O/a5UR45p/OOCCG6ewAWVbsuzR/Z+v1v76CBU+mTG/wPjqdpfg="] = 600
	// 	//爱的爆炸
	// 	giftToScoreMap["gx7pmjQfhBaDOG2XkWI2peZ66YFWkCWRjZXpTqb23O/epru+sxWyTV/3Ufs="] = 1600
	// 	//神秘空投
	// 	giftToScoreMap["pGLo7HKNk1i4djkicmJXf6iWEyd+pfPBjbsHmd3WcX0Ierm2UdnRR7UINvI="] = 2500
	// }
}

// 快手初始化礼物Id和name
// func ksGiftIdNameInit() {
// 	giftIdToName = make(map[string]string)
// 	lotteryMap = make(map[string]string)
// 	commentTogiftId = make(map[string]string)
// 	switch platform {
// 	case "ks":
// 		giftIdToName["11582"] = "助威火炬"
// 		giftIdToName["11584"] = "火花"
// 		giftIdToName["11585"] = "雪糕刺客"
// 		commentTogiftId["用雪糕"] = "11585"
// 		commentTogiftId["11585"] = "用雪糕"
// 		giftIdToName["11586"] = "冲鸭"
// 		commentTogiftId["用冲鸭"] = "11586"
// 		commentTogiftId["11586"] = "用冲鸭"
// 		giftIdToName["11587"] = "浪漫满屏"
// 		commentTogiftId["用浪漫"] = "11587"
// 		commentTogiftId["11587"] = "用浪漫"
// 		giftIdToName["11606"] = "助威火箭"
// 		commentTogiftId["用火箭"] = "11606"
// 		commentTogiftId["11606"] = "用火箭"
// 		giftIdToName["12252"] = "一叶知秋"
// 		commentTogiftId["用叶子"] = "12252"
// 		commentTogiftId["12252"] = "用叶子"
// 		// giftIdToName["12719"] = "爱心礼盒"
// 		giftIdToName["12720"] = "魔法相机"
// 		commentTogiftId["用相机"] = "12720"
// 		commentTogiftId["12720"] = "用相机"
// 		// giftIdToName["12721"] = "超酷跑车"
// 		// giftIdToName["12722"] = "蛋糕皇座"
// 		// giftIdToName["12723"] = "玩法之星"
// 		// giftIdToName["13549"] = "红药水"
// 		// giftIdToName["13550"] = "蓝药水"
// 		// giftIdToName["13551"] = "黄药水"
// 		// giftIdToName["13552"] = "绿药水"
// 		// giftIdToName["13585"] = "蓝钥匙"
// 		// giftIdToName["13586"] = "黄钥匙"
// 		// giftIdToName["13587"] = "粉钥匙"
// 		// giftIdToName["13588"] = "绿钥匙"
// 		// LotteryMap 抽奖
// 		lotteryMap["11582"] = "助威火炬"
// 		lotteryMap["12252"] = "一叶知秋"
// 		lotteryMap["11606"] = "助威火箭"
// 		lotteryMap["11585"] = "雪糕刺客"
// 		lotteryMap["11586"] = "冲鸭"
// 		lotteryMap["11587"] = "浪漫满屏"
// 		lotteryMap["12720"] = "魔法相机"
// 	case "dy":
// 		giftIdToName["n1/Dg1905sj1FyoBlQBvmbaDZFBNaKuKZH6zxHkv8Lg5x2cRfrKUTb8gzMs="] = "仙女棒"

// 		giftIdToName["28rYzVFNyXEXFC8HI+f/WG+I7a6lfl3OyZZjUS+CVuwCgYZrPrUdytGHu0c="] = "能力药丸"
// 		commentTogiftId["用药丸"] = "28rYzVFNyXEXFC8HI+f/WG+I7a6lfl3OyZZjUS+CVuwCgYZrPrUdytGHu0c="
// 		commentTogiftId["28rYzVFNyXEXFC8HI+f/WG+I7a6lfl3OyZZjUS+CVuwCgYZrPrUdytGHu0c="] = "用药丸"

// 		giftIdToName["fJs8HKQ0xlPRixn8JAUiL2gFRiLD9S6IFCFdvZODSnhyo9YN8q7xUuVVyZI="] = "魔法镜"
// 		commentTogiftId["用镜子"] = "fJs8HKQ0xlPRixn8JAUiL2gFRiLD9S6IFCFdvZODSnhyo9YN8q7xUuVVyZI="
// 		commentTogiftId["fJs8HKQ0xlPRixn8JAUiL2gFRiLD9S6IFCFdvZODSnhyo9YN8q7xUuVVyZI="] = "用镜子"

// 		giftIdToName["PJ0FFeaDzXUreuUBZH6Hs+b56Jh0tQjrq0bIrrlZmv13GSAL9Q1hf59fjGk="] = "甜甜圈"
// 		commentTogiftId["用甜甜圈"] = "PJ0FFeaDzXUreuUBZH6Hs+b56Jh0tQjrq0bIrrlZmv13GSAL9Q1hf59fjGk="
// 		commentTogiftId["PJ0FFeaDzXUreuUBZH6Hs+b56Jh0tQjrq0bIrrlZmv13GSAL9Q1hf59fjGk="] = "用甜甜圈"

// 		giftIdToName["IkkadLfz7O/a5UR45p/OOCCG6ewAWVbsuzR/Z+v1v76CBU+mTG/wPjqdpfg="] = "能量电池"
// 		commentTogiftId["用电池"] = "IkkadLfz7O/a5UR45p/OOCCG6ewAWVbsuzR/Z+v1v76CBU+mTG/wPjqdpfg="
// 		commentTogiftId["IkkadLfz7O/a5UR45p/OOCCG6ewAWVbsuzR/Z+v1v76CBU+mTG/wPjqdpfg="] = "用电池"

// 		giftIdToName["gx7pmjQfhBaDOG2XkWI2peZ66YFWkCWRjZXpTqb23O/epru+sxWyTV/3Ufs="] = "爱的爆炸"
// 		commentTogiftId["用炸弹"] = "gx7pmjQfhBaDOG2XkWI2peZ66YFWkCWRjZXpTqb23O/epru+sxWyTV/3Ufs="
// 		commentTogiftId["gx7pmjQfhBaDOG2XkWI2peZ66YFWkCWRjZXpTqb23O/epru+sxWyTV/3Ufs="] = "用炸弹"

// 		giftIdToName["pGLo7HKNk1i4djkicmJXf6iWEyd+pfPBjbsHmd3WcX0Ierm2UdnRR7UINvI="] = "神秘空投"
// 		commentTogiftId["用空投"] = "pGLo7HKNk1i4djkicmJXf6iWEyd+pfPBjbsHmd3WcX0Ierm2UdnRR7UINvI="
// 		commentTogiftId["pGLo7HKNk1i4djkicmJXf6iWEyd+pfPBjbsHmd3WcX0Ierm2UdnRR7UINvI="] = "用空投"
// 	}
// }
