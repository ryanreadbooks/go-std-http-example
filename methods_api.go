// methods api

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"io"
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

func GenerateJsonResponse(writer http.ResponseWriter, req *http.Request) ([]byte, error) {
	content, err := io.ReadAll(req.Body)
	if err != nil {
		if len(content) != 0 {
			fmt.Printf("content= %s\n", content)
		}
	}
	addDateHeader(writer)
	// 获取一些信息
	r := GetResponse{
		Args:    ExtractArgs(req.URL.Query()),
		Headers: ExtractHeaders(&req.Header),
		Url:     req.URL.String(),
		Body:    string(content),
	}
	return json.MarshalIndent(r, " ", "  ")
}


// /get请求的处理
type GetHandler struct{}

func addHeaderServer(writer http.ResponseWriter) {
	writer.Header().Add("server", "GoStandardHTTPServer")
}

func (h *GetHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	addHeaderServer(writer)
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
	addHeaderServer(writer)
	defer req.Body.Close()
	// 如果这里并不是请求进来的并不是post请求，则需要返回错误405
	if !ValidateMethod(writer, req, http.MethodPost) {
		return
	}
	// 这里确认是请求时post请求，可以正常处理
	// 获取请求体的内容
	err := req.ParseForm()
	if err != nil {
		fmt.Println("can not parse form")
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
		addDateHeader(writer)
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(json_bytes)
	}
}

type GetHeadersHandler struct{}

func (g *GetHeadersHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	addHeaderServer(writer)
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
	addDateHeader(writer)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(json_bytes)
}

type GetUserAgentHandler struct{}

func (g* GetUserAgentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addHeaderServer(w)
		defer r.Body.Close()
		if !ValidateMethod(w, r, http.MethodGet) {
			return
		}
		user_agent := r.Header.Get("User-Agent")
		json_bytes, _ := json.MarshalIndent(map[string]string{"User-Agent": user_agent}, " ", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(json_bytes)
}