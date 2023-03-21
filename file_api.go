package main

import (
	"net/http"
)

type ServeTxtFile struct{}

// /file/text
func (s *ServeTxtFile) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !ValidateMethod(w, req, http.MethodGet) {
		return
	}
	// 返回一个文件
	// 这个函数会设置一个Last-Modified响应头，用来告诉客户端此文件的最后修改时间
	// 然后客户端在下一次请求这个文件的时候，就会带上If-Modified-Since请求头，来询问此文件是否有最新版本
	// 如果此文件没有更新的话，就会返回304 Not Modified，从而节省这个文件的再次传输
	http.ServeFile(w, req, "files/plain.txt")
}

type ServeFolder struct{}

func (s *ServeFolder) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !ValidateMethod(w, req, http.MethodGet) {
		return
	}
	// 返回文件夹
	http.ServeFile(w, req, "files/dest")
}