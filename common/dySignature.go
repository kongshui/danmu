package common

import (
	"crypto/md5"
	"encoding/base64"
	"sort"
	"strings"
)

func DySignature(header map[string]string, bodyStr, secret string) string {
	keyList := make([]string, 0, 4)
	for key := range header {
		content := strings.ToLower(key)
		if content == "content-type" || content == "x-token" {
			continue
		}
		keyList = append(keyList, key)
	}
	sort.Slice(keyList, func(i, j int) bool {
		return keyList[i] < keyList[j]
	})
	kvList := make([]string, 0, 4)
	for _, key := range keyList {
		kvList = append(kvList, key+"="+header[key])
	}
	urlParams := strings.Join(kvList, "&")
	rawData := urlParams + bodyStr + secret
	md5Result := md5.Sum([]byte(rawData))
	return base64.StdEncoding.EncodeToString(md5Result[:])
}
