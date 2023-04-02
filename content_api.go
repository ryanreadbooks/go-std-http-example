package main

import (
	"net/http"
	"time"
)

type ServeVideoHandle struct{}

// 返回一个mp4视频文件
func (s *ServeVideoHandle) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	f := getOsFile(openedFiles, "assets/sample.mp4")
	if f == nil {
		MakeServerInternalError(w)
		return
	}
	// 开始serve content
	// 可以用来传视频数据
	http.ServeContent(w, req, "assets/sample.mp4", time.Now(), f)
}