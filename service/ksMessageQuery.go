package service

import (
	"encoding/json"
	"fmt"
	"path"
	"time"

	"github.com/kongshui/danmu/common"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 快手消息回调
func ksMessageQuery(roomId string) {
	request := map[string]any{
		"roomCode":  roomId,
		"timestamp": time.Now().UnixMilli(),
		"pushType":  "giftSend",
		"data":      "{}",
	}
	request["sign"] = common.KSSignature(request, app_secret, app_id)
	requestBody, err := json.Marshal(request)
	if err != nil {
		ziLog.Error(fmt.Sprintf("KsMessageQuery json marshal err:  %v", err), debug)
		return
	}
	header := map[string]string{
		"Content-Type": "application/json;charset=UTF-8",
	}
	urlPath := KsUrlSet(url_MessageQueryUrl)
	if urlPath == "" {
		ziLog.Error("KsMessageQuery err: urlPath is nil,ackType: ", debug)
		return
	}
	response, err := common.HttpRespond("POST", urlPath, requestBody, header)
	if err != nil {
		ziLog.Error(fmt.Sprintf("KsMessageQuery HttpRespond err:  %v", err), debug)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		ziLog.Error(fmt.Sprintf("KsMessageQuery HttpRespond code err:  %v", response.StatusCode), debug)
		return
	}
	request = map[string]any{}
	if err := json.NewDecoder(response.Body).Decode(&request); err != nil {
		ziLog.Error(fmt.Sprintf("KsMessageQuery json.NewDecoder err:  %v", err), debug)
		return
	}
	if int64(request["result"].(float64)) != 1 {
		ziLog.Error(fmt.Sprintf("KsMessageQuery err:  %v， data： %v", err, request), debug)
		return
	}
	go checkQuery(request["data"].(string))
}

// 查询快手礼物消息
func checkQuery(dataList string) {
	var (
		data []KsCallbackQueryStruct
	)
	dateListByte, err := json.Marshal(dataList)
	if err != nil {
		ziLog.Error(fmt.Sprintf("CheckQuery json marshal err:  %v", err), debug)
		return
	}
	if err := json.Unmarshal(dateListByte, &data); err != nil {
		ziLog.Error(fmt.Sprintf("CheckQuery json Unmarshal err:  %v", err), debug)
		return
	}
	for _, v := range data {
		code, err := rdb.SAdd(v.RoomCode+"giftSend", v.UniqueMessageId)
		if err != nil {
			continue
		}
		if code == 1 {
			go ksPushGiftSendPayloay(v)
		}
	}
}

// 失败消息获取
func getFailMessage() {
	t := time.NewTicker(time.Second * 15)
	isQueryOk := true
	for {
		<-t.C
		isQueryOk = false
		ok, err := rdb.SetKeyNX(monitor_fail_message_push_db, nodeUuid, 10*time.Second)
		if err != nil {
			ziLog.Error(fmt.Sprintf("ksCallBackQueryToKs 查询直播房间号失败:  %v", err), debug)
			isQueryOk = true
			continue
		}
		if !isQueryOk || !ok {
			continue
		}
		uidList, err := etcdClient.Client.Get(first_ctx, path.Join("/", config.Project, common.Uid_Register_RoomId_key), clientv3.WithPrefix())
		if err != nil {
			ziLog.Error(fmt.Sprintf("ksCallBackQueryToKs 查询直播房间号失败:  %v", err), debug)
			isQueryOk = true
			continue
		}
		if uidList.Count == 0 {
			continue
		}
		for i, kv := range uidList.Kvs {
			switch platform {
			case "ks":
				// 快手
				go ksMessageQuery(string(kv.Value))
			case "dy":
				// 抖音
				go pushFailDataToDy(string(kv.Value))
				if i > 0 && i%9 == 0 {
					time.Sleep(time.Second * 1)
				}
			}
		}
		isQueryOk = true
	}
}
