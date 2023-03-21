// 定义一些工具函数
package main

import (
	"net/http"
	"strings"
	"time"
)

// 获取当前的GMT时区的时间（好像就是UTC时间）
func fetchNowGMT() string {
	nowUTC := time.Now().UTC()
	t := nowUTC.Format(time.RFC1123)
	return strings.Replace(t, "UTC", "GMT", 1)
}

// 添加Date头部
func addDateHeader(w http.ResponseWriter) {
	w.Header().Set("Date", fetchNowGMT())
}

func MakeServerInternalError(writer http.ResponseWriter) {
	addDateHeader(writer)
	http.Error(writer, "Server Internal Error", http.StatusInternalServerError) // 用http包自带的函数来返回错误，在这个函数里面帮你组织好了各种信息
	// writer.WriteHeader(http.StatusInternalServerError)
	// writer.Write([]byte("Server Internal Error"))
}

func MethodNotAllowed(writer http.ResponseWriter) {
	addDateHeader(writer)
	http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
	// writer.WriteHeader(http.StatusMethodNotAllowed)
	// writer.Write([]byte("<h1>Method Not Allowed</h1>"))
}

// 验证请求方式是否为允许的方法
func ValidateMethod(writer http.ResponseWriter, req *http.Request, allow string) bool {
	if req.Method != allow {
		MethodNotAllowed(writer)
		return false
	}
	return true
}

