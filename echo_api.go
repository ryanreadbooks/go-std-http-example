package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type EchoHandler struct{}

func (e *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addHeaderServer(w)
	defer r.Body.Close()
	if !ValidateMethod(w, r, http.MethodGet) {
		return
	}
	urlRaw := r.URL.Path
	data := urlRaw[6:]
	fmt.Printf("data=%v\n", data)
	// 这里可能需要将url进行escape处理
	data, err := url.PathUnescape(data)
	fmt.Printf("data=%v\n", data)
	if err != nil {
		// 无法unescape，出错, 简单地返回500错误
		MakeServerInternalError(w)
		return
	}
	json_bytes, _ := json.MarshalIndent(map[string]string{"echoing": data}, " ", "  ")
	w.Header().Set("Content-Type", "application/json")
	addDateHeader(w)
	w.Write(json_bytes)
}
