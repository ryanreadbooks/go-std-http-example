// 和重定向相关的api的函数

package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type RedirectHandler struct{}

func (r *RedirectHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	addHeaderServer(w)
	defer req.Body.Close()
	if !ValidateMethod(w, req, http.MethodGet) {
		return
	}
	num, err := strconv.Atoi(req.URL.Query().Get("n"))
	if err != nil || num < 0 {
		// 无法将转换为数字，直接可以返回错误
		addDateHeader(w)
		w.Write([]byte("<h1>接口参数非法</h1><p>正确格式为: /absolute-redirect?n=xxx, xxxxxx必须为大于等于零的数字</p>"))
		return
	}
	if num == 0 {
		w.Header().Add("Location", "/get")
		w.WriteHeader(http.StatusFound)
	} else {
		re_loc := fmt.Sprintf("/absolute-redirect?n=%d", num-1)
		w.Header().Add("Location", re_loc)
		w.WriteHeader(http.StatusFound)
	}
	addDateHeader(w)
}
