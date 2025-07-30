package common

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"slices"
	"strings"
)

func KSSignature(signContext map[string]any, secret, appId string) string {
	signContext["app_id"] = appId
	keyList := make([]string, 0)
	for key := range signContext {
		if key == "sign" {
			continue
		}
		keyList = append(keyList, key)
	}
	slices.Sort(keyList)
	var buf strings.Builder
	for i, k := range keyList {
		if reflect.ValueOf(signContext[k]).Kind() == reflect.String {
			if i == len(keyList)-1 {
				buf.WriteString(fmt.Sprintf("%s=%s", k, signContext[k]))
				break
			}
			buf.WriteString(fmt.Sprintf("%s=%s&", k, signContext[k]))
		} else {
			dByte, _ := json.Marshal(signContext[k])
			if i == len(keyList)-1 {
				buf.WriteString(fmt.Sprintf("%s=%s", k, string(dByte)))
				break
			}
			buf.WriteString(fmt.Sprintf("%s=%s&", k, string(dByte)))
		}
		// reflect.ValueOf()
		// buf.WriteString(fmt.Sprintf("%s=%s&", k, signContext[k]))
	}
	rawStr := buf.String()
	rawData := strings.TrimRight(rawStr, "&") + secret
	h := md5.New()
	io.WriteString(h, rawData)
	md5str := fmt.Sprintf("%x", h.Sum(nil))
	delete(signContext, "app_id")
	// fmt.Println(hex.EncodeToString(h.Sum(nil)))
	return md5str
}

func KsCheckSignature(checkContext, secret, md5Compare string) bool {
	toString := checkContext + secret
	h := md5.New()
	h.Write([]byte(toString))
	md5str := fmt.Sprintf("%x", h.Sum(nil))
	return md5Compare == md5str
}
