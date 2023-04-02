package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 所以使用http/httptest包的关键就是一个实现了http.ResponseWriter接口的结构体httptest.ResponseRecorder
// 然后http请求需要自己使用httptest.NewRequest进行构造
// 然后自己手动调用相应的需要进行测试的处理函数
// 随后自己从http.ResponseRecorder中通过取响应结果
// 自己通过响应结果判断这个处理函数是否执行正确

func TestContentApi(t *testing.T) {
	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()

	// 如果是带有参数请求的一些参数的话，就可以在创建request对象的时候就进行相关参数的设计了
	r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8080/gzip", nil)
	// 直接用对应的handlder来处理
	handler := &GetGZipHandler{}
	// http.ResponseWriter是一个interface，需要实现的方法有三个: Write(buf), WriteHeader(code), Header()
	// 刚好*httptest.ResponseRecorder实现了http.ResponseWriter接口
	handler.ServeHTTP(recorder, r)
	// 经过上面的模拟处理后，recorder就可以拿到模拟的响应结果
	var resp *http.Response = recorder.Result()
	// 然后就可以开始检验想要的测试结果
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type expect application/json, but got %s", resp.Header.Get("Content-Type"))
	}
	if resp.Header.Get("Content-Encoding") != "gzip" {
		t.Errorf("Content-Encoding expects gzip, but got %s", resp.Header.Get("Content-Encoding"))
	}
}

func TestEchoApi(t *testing.T) {
	recorder := httptest.NewRecorder()

	r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8080/echo/HelloWorld", nil)
	handler := &EchoHandler{}
	handler.ServeHTTP(recorder, r)

	// 验证响应结果
	resp := recorder.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status code expects %d, but got %d", http.StatusOK, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("io.ReadAll err: %v", err)
	}
	type tmpStruct struct {
		Data string `json:"echoing"`
	}
	var tt tmpStruct
	json.Unmarshal(body, &tt)
	if tt.Data != "HelloWorld" {
		t.Errorf("echo expects HellWorld, but got %s", tt.Data)
	}
}
