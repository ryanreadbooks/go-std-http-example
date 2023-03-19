// 和auth认证有关的api函数

package main

import (
	"net/http"
	"strings"
	"net/url"
	"fmt"
	"encoding/base64"
	"encoding/json"
)

// 处理HTTP简单的认证功能
// basic-auth
type BasicAuthHandler struct{}

func (b *BasicAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addHeaderServer(w)
	if !ValidateMethod(w, r, http.MethodGet) {
		return
	}
	// 格式检查
	queries, err := url.QueryUnescape(r.URL.String())
	if err != nil || len(queries) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// 将username和password分开
	infos := strings.Split(queries[12:], "/")
	if len(infos) != 2 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("<h1>接口参数非法</h1>"))
		w.Write([]byte("<p>正确格式为： /basic-auth/user/passwd"))
		return
	}
	username, password := infos[0], infos[1]
	if username == "" || password == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("<h1>接口参数非法</h1>"))
		w.Write([]byte("<p>正确格式为： /basic-auth/user/passwd"))
		return
	}
	// 判断是不是已经登陆过的
	// 根据请求头中的Authorization字段来判断
	authCode := r.Header.Get("Authorization")
	if authCode != "" {
		// 进一步验证是否正确 Authorization: Basic yyyy
		if strings.EqualFold(authCode[:6], "basic ") {
			auth := authCode[6:]
			if auth != "" {
				// base64解码
				authRaw, err := base64.StdEncoding.DecodeString(auth)
				if err == nil {
					// 解码成功，和username和password匹配
					// 标准规定，密码是Basic认证的密码是明文传输的，只是将其进行了base64编码
					// 编码的格式为 base64(user:passwd)
					authTemplate := fmt.Sprintf("%s:%s", username, password)
					if string(authRaw) == authTemplate {
						// 验证成功，可以返回响应数据
						res := map[string]string{
							"authenticated": "true",
							"user":          username,
						}
						json_bytes, _ := json.MarshalIndent(res, " ", "  ")
						w.Write(json_bytes)
						return
					}
				}
			}
		}
	}
	// 不成功，返回401,并且指定认证方式为basic
	w.Header().Set("WWW-Authenticate", "Basic realm=\"Strict realm\"")
	w.WriteHeader(http.StatusUnauthorized)
}

// TODO : 处理使用digest的HTTP认证
type DigestAuthHandler struct{}