package service

// 自动选边
// func interactive(roomId, roundId string, label int) bool {
// 	headers := map[string]string{
// 		"Content-Type": "application/json;charset=UTF-8",
// 	}
// 	comment := map[string]any{
// 		"closeQuickInfoBySendGift": false,
// 		"quickInfos": []map[string]any{
// 			{
// 				"buttonColor": "#FE3666",
// 				"buttonText":  "免费加入守序善良",
// 				"commentText": "1",
// 				"extraText":   "扣1加入守序善良",
// 			},
// 			{
// 				"buttonColor": "#326BFB",  // 加方按钮颜色,可选:#FE3666,#326BFB,#17C5A2,#F6C000,#FE6636,#7140FF
// 				"buttonText":  "免费加入中立善良", //加方按钮文字，自定义，根据业务引导用户点击
// 				"commentText": "2",        //用来通过评论回传，必须是订阅的评论关键字
// 				"extraText":   "扣2加入中立善良", //额外拓展字段，如果不为空，则按钮变为双行展示，限制13个字
// 			},
// 			{
// 				"buttonColor": "#17C5A2",  // 加方按钮颜色,可选:#FE3666,#326BFB,#17C5A2,#F6C000,#FE6636,#7140FF
// 				"buttonText":  "免费加入混乱善良", //加方按钮文字，自定义，根据业务引导用户点击
// 				"commentText": "2",        //用来通过评论回传，必须是订阅的评论关键字
// 				"extraText":   "扣2加入混乱善良", //额外拓展字段，如果不为空，则按钮变为双行展示，限制13个字
// 			},
// 			{
// 				"buttonColor": "#FED236",  // 加方按钮颜色,可选:#FE3666,#326BFB,#17C5A2,#F6C000,#FE6636,#7140FF
// 				"buttonText":  "免费加入守序邪恶", //加方按钮文字，自定义，根据业务引导用户点击
// 				"commentText": "2",        //用来通过评论回传，必须是订阅的评论关键字
// 				"extraText":   "扣2加入守序邪恶", //额外拓展字段，如果不为空，则按钮变为双行展示，限制13个字
// 			},
// 			{
// 				"buttonColor": "#FE6636",  // 加方按钮颜色,可选:#FE3666,#326BFB,#17C5A2,#F6C000,#FE6636,#7140FF
// 				"buttonText":  "免费加入中立邪恶", //加方按钮文字，自定义，根据业务引导用户点击
// 				"commentText": "2",        //用来通过评论回传，必须是订阅的评论关键字
// 				"extraText":   "扣2加入中立邪恶", //额外拓展字段，如果不为空，则按钮变为双行展示，限制13个字
// 			},
// 			{
// 				"buttonColor": "#6A36FE",  // 加方按钮颜色,可选:#FE3666,#326BFB,#17C5A2,#F6C000,#FE6636,#7140FF
// 				"buttonText":  "免费加入混乱邪恶", //加方按钮文字，自定义，根据业务引导用户点击
// 				"commentText": "2",        //用来通过评论回传，必须是订阅的评论关键字
// 				"extraText":   "扣2加入混乱邪恶", //额外拓展字段，如果不为空，则按钮变为双行展示，限制13个字
// 			},
// 		},
// 	}
// 	switch label {
// 	case 1:
// 		comment["quickInfos"] = []map[string]any{
// 			{
// 				"buttonColor": "#326BFB",
// 				"buttonText":  "免费加入蓝队",
// 				"commentText": "1",
// 				"extraText":   "扣1加入蓝队",
// 			},
// 		}
// 	case 2:
// 		comment["quickInfos"] = []map[string]any{
// 			{
// 				"buttonColor": "#FE3666", // 加方按钮颜色,可选:#FE3666,#326BFB,#17C5A2,#F6C000,#FE6636,#7140FF
// 				"buttonText":  "免费加入红队",  //加方按钮文字，自定义，根据业务引导用户点击
// 				"commentText": "2",       //用来通过评论回传，必须是订阅的评论关键字
// 				"extraText":   "扣2加入红队",  //额外拓展字段，如果不为空，则按钮变为双行展示，限制13个字
// 			},
// 		}
// 	default:
// 		break
// 	}

// 	commentByte, _ := json.Marshal(comment)
// 	bodyStruct := map[string]any{
// 		"roomCode":  roomId,
// 		"timestamp": time.Now().UnixMilli(),
// 		"roundId":   roundId,
// 		"type":      "1",
// 		"data":      string(commentByte),
// 	}
// 	ms5Str := common.KSSignature(bodyStruct, app_secret, app_id)
// 	bodyStruct["sign"] = ms5Str
// 	body, err := json.Marshal(bodyStruct)
// 	if err != nil {
// 		ziLog.Error(fmt.Sprintf("interactive json.Marshal err: %v", err), debug)
// 		return false
// 	}
// 	urlPath := urlSet(url_InteractiveUrl)
// 	if urlPath == "" {
// 		ziLog.Error("interactive err, urlPath is nil", debug)
// 		return false
// 	}
// 	response, err := common.HttpRespond("POST", urlPath, body, headers)
// 	if err != nil {
// 		ziLog.Error(fmt.Sprintf("interactive response err:  %v", err), debug)
// 		return false
// 	}
// 	defer response.Body.Close()
// 	var (
// 		request any
// 	)

// 	if err := json.NewDecoder(response.Body).Decode(&request); err != nil {
// 		ziLog.Error(fmt.Sprintf("interactive NewDecoder err:  %v", err), debug)
// 		return false
// 	}
// 	if response.StatusCode != 200 {
// 		return false
// 	}
// 	if int64(request.(map[string]any)["result"].(float64)) != 1 {
// 		ziLog.Error(fmt.Sprintf("interactive err,roomid:%v, context:%v err:  %v", roomId, request, err), debug)
// 		return false
// 	}
// 	return true

// }
