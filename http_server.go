// 这个文件中的函数样式net/http中的函数怎样使用，怎样使用来搭建一个http服务器
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

// 缓存打开的文件
var openedFiles map[string]*os.File = make(map[string]*os.File)

func http_server(host string, port uint, tls bool) {
	defer func() {
		if len(openedFiles) != 0 {
			for _, v := range openedFiles {
				v.Close()
			}
		}
	}()

	mux := http.NewServeMux()
	// 注册路由
	// HTTP方法接口
	mux.Handle("/get", &GetHandler{})
	mux.Handle("/post", &PostHandler{})

	// HTTP request inspect
	mux.Handle("/headers", &GetHeadersHandler{})
	mux.Handle("/user-agent", &GetUserAgentHandler{})

	// Dynamic data
	// 解码base64的内容
	mux.Handle("/base64/", &Base64APIHandler{})

	// cookieis
	// 这个接口用来设置和返回cookie信息
	mux.Handle("/cookies", &GetCookiesHandler{})
	// /cookies/set?xx=yy&xx=yy为设置cookie的接口方式
	mux.Handle("/cookies/set", &CookiesSetHandler{})
	// 删除一个cookie的接口
	mux.Handle("/cookies/delete", &CookiesDeleteHandler{})

	// images
	mux.Handle("/image/jpeg", &GetImageJpegHandler{})
	// /image/png接口，返回一张png格式的图片
	mux.Handle("/image/png", &GetImagePngHandler{})

	// redirects
	// 重定向接口的处理
	// 接口格式 //absolute-redirect?n=2 其中查询参数n表示重定向的次数
	mux.Handle("/absolute-redirect", &RedirectHandler{})

	// /echo/xxx
	mux.Handle("/echo/", &EchoHandler{})

	// 响应值的格式：压缩格式的接口
	mux.Handle("/deflate", &GetDeflateHandler{})
	mux.Handle("/gzip", &GetGZipHandler{})

	// 响应简单的html文件
	mux.Handle("/html", &GetHTMLHandler{})
	// 响应简单的xml文件
	mux.Handle("/xml", &GetXMLHandler{})

	// 和缓存相关的一些功能
	// 如果request header中存在If-Modified-Since或者If-None-Match，则返回304，否则返回200
	mux.Handle("/cache", &CacheHandler{})
	// 请求服务器设置一个cache-control， 其max-age为n
	mux.Handle("/cache/", &CacheNHandler{})
	// 模拟一个响应的资源存在其etag
	mux.Handle("/etag/", &ETagHandler{})

	// http身份认证相关
	// /basic-auth/name/passwd
	mux.Handle("/basic-auth/", &BasicAuthHandler{})

	// 文件
	// 这个接口返回单独一个文件
	mux.Handle("/text", &ServeTxtFile{})
	// 这个接口返回一个文件夹里面的内容，注意，如果要是返回文件夹的话，路由中必须以'/'结尾，否则会找不到资源从而返回404
	mux.Handle("/folder/", &ServeFolder{})

	addr := fmt.Sprintf("%s:%d", host, port)

	var err error
	// 开始监听
	if tls {
		// 提供https服务
		// 使用go自带的在crypto/tls中的genereate_cert.go函数生成cert.pem和key.pem两个文件
		if !CheckFileExists("cert.pem") || !CheckFileExists("key.pem") {
			fmt.Println("generating cert file")
			goRoot := runtime.GOROOT()
			scriptPath := path.Join(goRoot, "src", "crypto", "tls", "generate_cert.go")
			// 运行scriptPath生成两个文件
			// 无法校验颁发机构，，浏览器提示不安全。
			// 浏览器只认CA颁发的证书
			cmd := exec.Command("go", "run", scriptPath, "-host", host, "-ca", "true")
			// 将标准输出重定向
			var out strings.Builder
			cmd.Stdout = &out
			err = cmd.Run()
			if err != nil {
				panic(err)
			}
		}
		fmt.Println("tls enabled")
		err = http.ListenAndServeTLS(addr, "cert.pem", "key.pem", mux)
	} else {
		err = http.ListenAndServe(addr, mux)
	}
	if err != nil {
		log.Fatalf("Fatal err: %v\n", err)
	}
}
