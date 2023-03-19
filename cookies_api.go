// 和cookies相关的api函数

package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

type GetCookiesHandler struct{}

func ConstructCookieToJson(cookies []*http.Cookie) []byte {
	var cookie_pairs map[string]string = make(map[string]string)
	for _, item := range cookies {
		cookie_pairs[item.Name] = item.Value
	}
	var c = make(map[string]map[string]string)
	c["cookies"] = cookie_pairs
	json_bytes, err := json.MarshalIndent(c, " ", "  ")
	if err != nil {
		return []byte{}
	}
	return json_bytes
}

func (g *GetCookiesHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	addHeaderServer(writer)
	defer req.Body.Close()
	if !ValidateMethod(writer, req, http.MethodGet) {
		return
	}
	// 获取请求的cookies，然后包装之后返回
	var cookies []*http.Cookie = req.Cookies()
	content := ConstructCookieToJson(cookies)
	if len(content) == 0 {
		MakeServerInternalError(writer)
		return
	}
	addDateHeader(writer)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(content)
}

type CookiesSetHandler struct{}

func (c *CookiesSetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addHeaderServer(w)
	defer r.Body.Close()
	if !ValidateMethod(w, r, http.MethodGet) {
		return
	}
	// 获取url中的查询参数，然后将他们设置为cookies
	var queries url.Values = r.URL.Query()
	if len(queries) == 0 {
		// 请求格式不符合要求，返回错误
		addDateHeader(w)
		w.Write([]byte("<h1>Invalid syntax when set cookies</h1>"))
		w.Write([]byte("<h1>syntax is like: /cookies/set?key1=value1&key2=value2</h1>"))
		return
	}
	// 格式符合要求
	// 给客户端设置cookie
	for k, v := range queries {
		if len(v) != 0 {
			// 每个cookie都设置有效期为0min
			http.SetCookie(w, &http.Cookie{Name: k, Value: v[0], Expires: time.Now().Add(10 * time.Minute)})
		}
	}

	// 重定向到/cookies
	addDateHeader(w)
	w.Header().Set("Location", "/cookies")
	w.WriteHeader(http.StatusFound)
}

type CookiesDeleteHandler struct{}

func (c *CookiesDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addHeaderServer(w)
		defer r.Body.Close()
		if !ValidateMethod(w, r, http.MethodGet) {
			return
		}
		// 检查url中的查询参数是否符合要求，然后获取到要删除的cookie的key
		queries := r.URL.Query()
		if len(queries) == 0 {
			// 请求格式不符合要求，返回错误
			addDateHeader(w)
			w.Write([]byte("<h1>Invalid syntax when set cookies</h1>"))
			w.Write([]byte("<h1>syntax is like: /cookies/set?key1=value1&key2=value2</h1>"))
			return
		}
		// 格式正确
		for k := range queries {
			http.SetCookie(w, &http.Cookie{Name: k, MaxAge: -1})
		}
		// 重定向到/cookies
		addDateHeader(w)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", "/cookies")
		w.WriteHeader(http.StatusFound)
}