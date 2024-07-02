package auth_client

import (
	"fmt"
	app_http "gitee.com/go-nianxi/go-common/pkg/transports/http"
	"github.com/gin-gonic/gin"
	"strings"
)

var (
	staticExt     = [...]string{"css", "html", "jpg", "ttf", "woff", "woff2", "eot", "eot", "git", "png", "jpeg", "apk", "js", "bmp", "svg", "mp3", "mp4", "mpeg", "jpg", "bmp", "mjs", "ico", "json", "yaml", "yml", "xml"}
	constPathList = [...]string{"/auth", "/captcha", "/api-doc", "/ping", "info", "/test"}
)

func isStatic(path string) bool {
	if strings.HasPrefix(path, app_http.StaticUrlPrefix) {
		return true
	}

	path = strings.Split(path, "#")[0]
	for _, v := range staticExt {
		if strings.HasSuffix(path, "."+v) {
			return true
		}
	}

	path = strings.Split(path, "?")[0]
	for _, v := range staticExt {
		if strings.HasSuffix(path, "."+v) {
			return true
		}
	}
	return false
}

func (authClient *ClientContext) ignorePath(c *gin.Context) bool {
	//静态资源全部忽略
	path := c.Request.URL.Path
	if isStatic(path) {
		return true
	}

	for _, v := range constPathList {
		if strings.HasPrefix(path, v) {
			return true
		}
	}

	ignore := authClient.conf.O.Client.Ignore
	context := strings.Split(c.FullPath(), "*")[0]
	list := strings.Split(ignore, ",")
	path = strings.TrimPrefix(c.Request.URL.Path, context)
	//可以忽略
	for _, v := range list {
		v = strings.TrimPrefix(v, context)
		if v == path {
			return true
		}
		if strings.HasSuffix(v, "*") {
			v1 := strings.TrimSuffix(v, "*")
			if strings.HasSuffix(path, v1) {
				return true
			}
		}
	}

	//必须认证
	auth := authClient.conf.O.Client.Auth
	list = strings.Split(auth, ",")
	for _, v := range list {
		v = strings.TrimPrefix(v, context)
		if v == path {
			return false
		}
		if strings.HasSuffix(v, "*") {
			v1 := strings.TrimSuffix(v, "*")
			if strings.HasSuffix(path, v1) {
				return false
			}
		}
	}

	return false
}

func filterURl(urlStr string) string {
	url := []string{""}
	url = append(url, substr(urlStr, 0, strings.Index(urlStr, "?")))
	fmt.Println(url)
	urlStr = substr(urlStr, strings.Index(urlStr, "?")+1, len(urlStr))
	list := strings.Split(urlStr, "&")
	for _, v := range list {
		if !strings.HasPrefix(v, "sign=") && !strings.HasPrefix(v, "token=") && !strings.HasPrefix(v, "flag") {
			if len(url) == 1 {
				url = append(url, "?"+v)
			} else {
				url = append(url, "&"+v)
			}
		}
	}
	result := ""
	for _, v := range url {
		result += v
	}
	return result
}

func substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}
