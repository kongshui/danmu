package common

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func HttpRespond(method, url string, body []byte, header map[string]string) (*http.Response, error) {
	// //创建body
	// requestBody, err := json.Marshal(map[string]any{}{
	// 	"appid":      "",
	// 	"secret":     "",
	// 	"grant_type": "client_credential",
	// })
	// if err != nil {
	// 	return err
	// }
	// 创建一个新的HTTP请求对象，该对象用于设置请求的URL、请求方法、请求头、请求体等
	request, err := http.NewRequest(strings.ToUpper(method), url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	// 设置请求头
	for k, v := range header {
		request.Header.Set(k, v)
	}
	// request.Header.Set("Content-Type", "application/json")
	// 创建HTTP客户端并设置超时时间（可选）
	client := &http.Client{
		Timeout: time.Second * 10, // 设置超时时间为10秒
	}
	// 发送请求并获取响应
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func GetUrl(baseUrl string, query map[string]string) (string, error) {
	// 创建一个url.URL对象并设置其Path和RawQuery属性（注意：RawQuery会自动处理编码）
	u, err := url.Parse(baseUrl) // 先解析基础URL部分
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return "", err
	}
	q := u.Query() // 创建一个新的url.Values对象
	for k, v := range query {
		q.Set(k, v) // 将查询参数添加到url.Values对象中
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}
