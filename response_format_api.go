// 和响应格式有关的api的函数实现

package main

import (
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

type GetDeflateHandler struct{}

func (g *GetDeflateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	addHeaderServer(w)
	defer req.Body.Close()
	if !ValidateMethod(w, req, http.MethodGet) {
		return
	}

	json_bytes, err := GenerateJsonResponse(w, req)
	if err != nil {
		MakeServerInternalError(w)
		return
	}
	// 将json内容进行压缩,
	flateWriter, err := flate.NewWriter(w, flate.BestSpeed)
	defer func() { _ = flateWriter.Close() }()
	if err != nil {
		MakeServerInternalError(w)
		return
	} else {
		// 响应头最好是在写入压缩的内容之前填入
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "deflate")
		// 往flateWriter里面些内容
		_, err := flateWriter.Write(json_bytes)
		if err != nil {
			MakeServerInternalError(w)
			return
		}
		// marshal ok
		flateWriter.Flush() // 将flateWriter底层的writer(也就是使用flate.NewWriter创建的时候传入的参数io.Writer接口)中的内容马上写出去
	}
}

type GetGZipHandler struct{}

func (g *GetGZipHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addHeaderServer(w)
	defer r.Body.Close()
	if !ValidateMethod(w, r, http.MethodGet) {
		return
	}
	json_bytes, err := GenerateJsonResponse(w, r)
	if err != nil {
		MakeServerInternalError(w)
		return
	}
	// 响应头最好是在填入压缩的内容之前设置
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "application/json")
	gzWriter := gzip.NewWriter(w)
	defer func() { _ = gzWriter.Close() }()
	_, err = gzWriter.Write(json_bytes)
	if err != nil {
		MakeServerInternalError(w)
		return
	}
	// 添加响应头
	gzWriter.Flush()
}

type GetHTMLHandler struct{}

func (g *GetHTMLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addHeaderServer(w)
	if !ValidateMethod(w, r, http.MethodGet) {
		return
	}
	img_file := getOsFile(openedFiles, "assets/home.html")
	if img_file == nil {
		MakeServerInternalError(w)
		return
	}
	// 将img_file的内容复制到http.ResponseWriter中
	addDateHeader(w)
	w.Header().Set("Content-Type", "text/html")
	_, err := io.Copy(w, img_file)
	img_file.Seek(0, io.SeekStart) // 文件指针回到文件开头的位置
	if err != nil {
		fmt.Printf("err when io.Copy: %v\n", err)
	}
}

type GetXMLHandler struct{}

func (g *GetXMLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addHeaderServer(w)
	if !ValidateMethod(w, r, http.MethodGet) {
		return
	}
	img_file := getOsFile(openedFiles, "assets/test.xml")
	if img_file == nil {
		MakeServerInternalError(w)
		return
	}
	// 将img_file的内容复制到http.ResponseWriter中
	w.Header().Set("Content-Type", "text/xml")
	addDateHeader(w)
	_, err := io.Copy(w, img_file)
	img_file.Seek(0, io.SeekStart) // 文件指针回到文件开头的位置
	if err != nil {
		fmt.Printf("err when io.Copy: %v\n", err)
	}
}
