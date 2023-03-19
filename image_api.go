// 和图片apy有关的函数

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type GetImageJpegHandler struct{}

// 复制函数
func getOsFile(openedFiles map[string]*os.File, name string) *os.File {
	if _, ok := openedFiles[name]; !ok {
		f, err := os.Open(name)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return nil
		}
		openedFiles[name] = f
	}
	return openedFiles[name]
}

func (g *GetImageJpegHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addHeaderServer(w)
	defer r.Body.Close()
	if !ValidateMethod(w, r, http.MethodGet) {
		return
	}
	img_file := getOsFile(openedFiles, "images/leaves_jpeg.jpg")
	if img_file == nil {
		MakeServerInternalError(w)
		return
	}
	// 将img_file的内容复制到http.ResponseWriter中
	addDateHeader(w)
	w.Header().Set("Content-Type", "image/jpeg")
	_, err := io.Copy(w, img_file)
	img_file.Seek(0, io.SeekStart) // 文件指针回到文件开头的位置
	if err != nil {
		fmt.Printf("err when io.Copy: %v\n", err)
	}
}

type GetImagePngHandler struct{}

func (g *GetImagePngHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addHeaderServer(w)
	defer r.Body.Close()
	if !ValidateMethod(w, r, http.MethodGet) {
		return
	}
	img_file := getOsFile(openedFiles, "images/cat_png.png")
	if img_file == nil {
		MakeServerInternalError(w)
		return
	}
	addDateHeader(w)
	w.Header().Set("Content-Type", "image/png")
	_, err := io.Copy(w, img_file)
	img_file.Seek(0, io.SeekStart) // 文件指针回到文件开头的位置
	if err != nil {
		fmt.Printf("err when io.Copy: %v\n", err)
	}
}
