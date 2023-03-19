// 和dynamic apx相关的函数

package main

import (
	"net/http"
	"encoding/base64"
)

type Base64APIHandler struct {}

func (b *Base64APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addHeaderServer(w)
		defer r.Body.Close()
		if !ValidateMethod(w, r, http.MethodGet) {
			return
		}
		// 获取参数
		data := r.URL.String()[8:]
		// 当成base64的编码内容进行解码，然后返回
		dst, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			// internal error
			MakeServerInternalError(w)
			return
		}
		// decode ok
		addDateHeader(w)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(dst)
}