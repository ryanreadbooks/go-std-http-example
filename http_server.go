// 这个文件中的函数样式net/http中的函数怎样使用，怎样使用来搭建一个http服务器
package main

import (
	"compress/flate"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// get请求的响应内容
type GetResponse struct {
	Args    map[string]string `json:"args"`
	Headers map[string]string `json:"headers"`
	Url     string            `json:"url"`
	Body    string            `json:"body"`
}

// post请求的响应内容
type PostResponse struct {
	GetResponse
	Forms map[string]string `json:"form"`
}

func Extract(data map[string][]string) map[string]string {
	var ret map[string]string = make(map[string]string)
	for k, v := range data {
		if len(v) == 1 {
			ret[k] = v[0]
		} else {
			var s string
			for idx, item := range v {
				s += item
				if idx != len(v)-1 {
					s += ", "
				}
			}
		}
	}
	return ret
}

func ExtractArgs(query url.Values) map[string]string {
	// url.Values本质是map[string][]string类型的
	// 我们需要对query的内容进行转义
	var ret map[string]string = make(map[string]string)
	for k, v := range query {
		if len(v) == 1 {
			unescaped_v, _ := url.QueryUnescape(v[0])
			ret[k] = unescaped_v
		} else {
			var s string
			for idx, item := range v {
				unescaped_item, _ := url.QueryUnescape(item)
				s += unescaped_item
				if idx != len(v)-1 {
					s += ", "
				}
			}
		}
	}
	return ret
}

func ExtractHeaders(header *http.Header) map[string]string {
	// url.Values本质是map[string][]string类型的
	return Extract(*header)
}

func ExtractForm(form url.Values) map[string]string {
	return Extract(form)
}

func MakeServerInternalError(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write([]byte("Server Internal Error"))
}

func MethodNotAllowed(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusMethodNotAllowed)
	writer.Write([]byte("<h1>Method Not Allowed</h1>"))
}

func ValidateMethod(writer http.ResponseWriter, req *http.Request, allow string) bool {
	if req.Method != allow {
		MethodNotAllowed(writer)
		return false
	}
	return true
}

// /get请求的处理
type GetHandler struct{}

func GenerateJsonResponse(writer http.ResponseWriter, req *http.Request) ([]byte, error) {
	content, err := io.ReadAll(req.Body)
	if err != nil {
		if len(content) != 0 {
			fmt.Printf("content= %s\n", content)
		}
	}
	writer.Header().Add("server", "GoStandardHTTPServer")
	// 获取一些信息
	r := GetResponse{
		Args:    ExtractArgs(req.URL.Query()),
		Headers: ExtractHeaders(&req.Header),
		Url:     req.URL.String(),
		Body:    string(content),
	}
	return json.MarshalIndent(r, " ", "  ")
}

func (h *GetHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	if !ValidateMethod(writer, req, http.MethodGet) {
		return
	}
	json_bytes, err := GenerateJsonResponse(writer, req)
	if err != nil {
		MakeServerInternalError(writer)
		return
	} else {
		// marshal ok
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(json_bytes)
	}
}

type PostHandler struct{}

func (p *PostHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	// 如果这里并不是请求进来的并不是post请求，则需要返回错误405
	if !ValidateMethod(writer, req, http.MethodPost) {
		return
	}
	// 这里确认是请求时post请求，可以正常处理
	// 获取请求体的内容
	err := req.ParseForm()
	if err != nil {
		log.Println("can not parse form")
		MakeServerInternalError(writer)
		return
	}
	r := PostResponse{
		GetResponse: GetResponse{
			Args:    ExtractArgs(req.URL.Query()),
			Headers: ExtractHeaders(&req.Header),
			Url:     req.URL.String(),
		},
		Forms: ExtractForm(req.Form),
	}
	json_bytes, err := json.MarshalIndent(r, " ", "  ")
	if err != nil {
		MakeServerInternalError(writer)
		return
	} else {
		// marshal ok
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(json_bytes)
	}
}

type GetHeadersHandler struct{}

func (g *GetHeadersHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	if !ValidateMethod(writer, req, http.MethodGet) {
		return
	}
	headers := ExtractHeaders(&req.Header)
	json_bytes, err := json.MarshalIndent(headers, " ", "  ")
	if err != nil {
		MakeServerInternalError(writer)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(json_bytes)
}

type GetCookiesHandler struct{}

func ConstructCookieToJson(cookies []*http.Cookie) []byte {
	var cookie_pairs map[string]string = make(map[string]string)
	for _, item := range cookies {
		cookie_pairs[item.Name] = item.Value
	}
	var c = make(map[string]map[string]string)
	c["cookies"] = cookie_pairs
	json_bytes, err := json.MarshalIndent(c, " ", "  ")
	if err != nil {
		return []byte{}
	}
	return json_bytes
}

func (g *GetCookiesHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	if !ValidateMethod(writer, req, http.MethodGet) {
		return
	}
	// 获取请求的cookies，然后包装之后返回
	var cookies []*http.Cookie = req.Cookies()
	content := ConstructCookieToJson(cookies)
	if len(content) == 0 {
		MakeServerInternalError(writer)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(content)
}

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

func http_server() {
	mux := http.NewServeMux()
	// 注册路由
	// HTTP方法接口
	mux.Handle("/get", &GetHandler{})
	mux.Handle("/post", &PostHandler{})

	// HTTP request inspect
	mux.Handle("/headers", &GetHeadersHandler{})
	mux.HandleFunc("/user-agent", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if !ValidateMethod(w, r, http.MethodGet) {
			return
		}
		user_agent := r.Header.Get("User-Agent")
		json_bytes, _ := json.MarshalIndent(map[string]string{"User-Agent": user_agent}, " ", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(json_bytes)
	})

	// Dynamic data
	// 解码base64的内容
	mux.HandleFunc("/base64/", func(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(dst)
	})

	// cookieis
	// 这个接口用来设置和返回cookie信息
	mux.Handle("/cookies", &GetCookiesHandler{})
	// /cookies/set?xx=yy&xx=yy为设置cookie的接口方式
	mux.HandleFunc("/cookies/set", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if !ValidateMethod(w, r, http.MethodGet) {
			return
		}
		// 获取url中的查询参数，然后将他们设置为cookies
		var queries url.Values = r.URL.Query()
		if len(queries) == 0 {
			// 请求格式不符合要求，返回错误
			w.Write([]byte("<h1>Invalid syntax when set cookies</h1>"))
			w.Write([]byte("<h1>syntax is like: /cookies/set?key1=value1&key2=value2</h1>"))
			return
		}
		// 格式符合要求
		// 给客户端设置cookie
		for k, v := range queries {
			if len(v) != 0 {
				// 每个cookie都设置有效期为0min
				http.SetCookie(w, &http.Cookie{Name: k, Value: v[0], Expires: time.Now().Add(10 * time.Minute)})
			}
		}

		// 重定向到/cookies
		w.Header().Set("Location", "/cookies")
		w.WriteHeader(http.StatusFound)
	})

	// 删除一个cookie的接口
	mux.HandleFunc("/cookies/delete", func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()
		if !ValidateMethod(w, r, http.MethodGet) {
			return
		}
		// 检查url中的查询参数是否符合要求，然后获取到要删除的cookie的key
		queries := r.URL.Query()
		if len(queries) == 0 {
			// 请求格式不符合要求，返回错误
			w.Write([]byte("<h1>Invalid syntax when set cookies</h1>"))
			w.Write([]byte("<h1>syntax is like: /cookies/set?key1=value1&key2=value2</h1>"))
			return
		}
		// 格式正确
		for k := range queries {
			http.SetCookie(w, &http.Cookie{Name: k, MaxAge: -1})
		}
		// 重定向到/cookies
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", "/cookies")
		w.WriteHeader(http.StatusFound)
	})

	var openedFiles map[string]*os.File = make(map[string]*os.File)
	defer func() {
		if len(openedFiles) != 0 {
			for _, v := range openedFiles {
				v.Close()
			}
		}
	}()

	// images
	mux.HandleFunc("/image/jpeg", func(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Set("Content-Type", "image/jpeg")
		_, err := io.Copy(w, img_file)
		img_file.Seek(0, io.SeekStart) // 文件指针回到文件开头的位置
		if err != nil {
			fmt.Printf("err when io.Copy: %v\n", err)
		}
	})
	// /image/png接口，返回一张png格式的图片
	mux.HandleFunc("/image/png", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if !ValidateMethod(w, r, http.MethodGet) {
			return
		}
		img_file := getOsFile(openedFiles, "images/cat_png.png")
		if img_file == nil {
			MakeServerInternalError(w)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		_, err := io.Copy(w, img_file)
		img_file.Seek(0, io.SeekStart) // 文件指针回到文件开头的位置
		if err != nil {
			fmt.Printf("err when io.Copy: %v\n", err)
		}
	})

	// redirects
	// 重定向接口的处理
	// 接口格式 //absolute-redirect?n=2 其中查询参数n表示重定向的次数
	mux.HandleFunc("/absolute-redirect", func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		if !ValidateMethod(w, req, http.MethodGet) {
			return
		}
		num, err := strconv.Atoi(req.URL.Query().Get("n"))
		if err != nil || num < 0 {
			// 无法将转换为数字，直接可以返回错误
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
	})

	// /echo/xxx
	mux.HandleFunc("/echo/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if !ValidateMethod(w, r, http.MethodGet) {
			return
		}
		urlRaw := r.URL.String()
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
		w.Write(json_bytes)
	})

	// 响应值的格式：压缩格式的接口
	mux.HandleFunc("/deflate", func(w http.ResponseWriter, req *http.Request) {
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
	})

	mux.HandleFunc("/gzip", func(w http.ResponseWriter, r *http.Request) {
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
	})

	// 响应简单的html文件
	mux.HandleFunc("/html", func(w http.ResponseWriter, r *http.Request) {
		if !ValidateMethod(w, r, http.MethodGet) {
			return
		}
		img_file := getOsFile(openedFiles, "assets/home.html")
		if img_file == nil {
			MakeServerInternalError(w)
			return
		}
		// 将img_file的内容复制到http.ResponseWriter中
		w.Header().Set("Content-Type", "text/html")
		_, err := io.Copy(w, img_file)
		img_file.Seek(0, io.SeekStart) // 文件指针回到文件开头的位置
		if err != nil {
			fmt.Printf("err when io.Copy: %v\n", err)
		}
	})
	// 响应简单的xml文件

	mux.HandleFunc("/xml", func(w http.ResponseWriter, r *http.Request) {
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
		_, err := io.Copy(w, img_file)
		img_file.Seek(0, io.SeekStart) // 文件指针回到文件开头的位置
		if err != nil {
			fmt.Printf("err when io.Copy: %v\n", err)
		}
	})

	// 开始监听
	err := http.ListenAndServe("127.0.0.1:8080", mux)
	if err != nil {
		log.Fatalf("Fatal err: %v\n", err)
	}

}
