// 和缓存相关的api函数

package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type CacheHandler struct{}

func (c *CacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addHeaderServer(w)
	if !ValidateMethod(w, r, http.MethodGet) {
		return
	}
	// 获取header
	json_bytes, err := GenerateJsonResponse(w, r)
	if err != nil {
		MakeServerInternalError(w)
		return
	}
	// 没有
	if r.Header.Get("If-Modified-Since") != "" || r.Header.Get("If-None-Match") != "" {
		// 只要有这两个的其中一个，我们就返回304 Not Modified
		w.WriteHeader(http.StatusNotModified)
		// 根据HTTP的标准，304状态码是不允许有response body的，所以如果在304作为返回状态码的情况下，
		// 调用w.Write()是不会成功的
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_bytes)
	// 查看w真正的类型: *http.response，是一个内部的结构体
}

type CacheNHandler struct{}

func (c *CacheNHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addHeaderServer(w)
	if !ValidateMethod(w, r, http.MethodGet) {
		return
	}
	// 获取参数
	queries := r.URL.String()
	unescaped_q, _ := url.QueryUnescape(queries)
	param := unescaped_q[7:]
	n, err := strconv.Atoi(param)
	// 不能转化成数字，直接返回错误
	if err != nil {
		addDateHeader(w)
		w.Write([]byte("<h1>接口参数非法，正确格式为：/cache/n, 其中n必须为数字</h1>"))
		return
	}
	// 设置cache-control响应头给客户端
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", n))
	// 编辑content
	json_bytes, err := GenerateJsonResponse(w, r)
	if err != nil {
		MakeServerInternalError(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_bytes)
}

type ETagHandler struct{}

func (e *ETagHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addHeaderServer(w)
	if !ValidateMethod(w, r, http.MethodGet) {
		return
	}
	// 获取参数
	queries := r.URL.String()
	unescaped_q, _ := url.QueryUnescape(queries)
	etagFromQuery := unescaped_q[6:]
	if len(etagFromQuery) == 0 {
		addDateHeader(w)
		w.Write([]byte("<h1>接口参数非法：正确格式为 /etag/xxxxx, xxxxx为请求设置的etag值"))
		return
	}
	ifNoneMatchEtag := r.Header.Get("If-None-Match")
	if ifNoneMatchEtag == "" {
		// 返回200
		// 设置etag返回给客户端
		json_bytes, err := GenerateJsonResponse(w, r)
		if err != nil {
			MakeServerInternalError(w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("ETag", etagFromQuery)
		w.Write(json_bytes)
		return
	} else {
		// 检查请求的etag和ifnonematch的etag是否相同
		if etagFromQuery == ifNoneMatchEtag {
			// 相同就返回304
			w.WriteHeader(http.StatusNotModified)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}
