package service

func urlSet(url string) string {
	if url == "" {
		return ""
	}
	// uriMap := map[string]string{
	// 	"app_id":       app_id,
	// 	"access_token": accessToken.Token,
	// }
	// urlPath, err := common.GetUrl(url, uriMap)
	// if err != nil {
	// 	return ""
	// }
	urlPath := url + "?app_id=" + app_id + "&access_token=" + accessToken.Token
	return urlPath
}

// shezhi url
func setUrl() {
	switch platform {
	case "ks":
		// 快手设置
		url_GetAccessTokenUrl = "https://open.kuaishou.com/oauth2/access_token" //获取全局token的url

	case "dy":
		// 抖音设置
		url_GetAccessTokenUrl = "https://developer.toutiao.com/api/apps/v2/token" //获取全局token的url
	}
}
