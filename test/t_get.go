package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Get 基本请求、处理JSON响应
func Get() {
	type JsonDec struct {
		Args    struct{} `json:"args"`
		Headers struct {
			AcceptEncoding string `json:"Accept-Encoding"`
			Host           string `json:"Host"`
			UserAgent      string `json:"User-Agent"`
			XAmznTraceId   string `json:"X-Amzn-Trace-Id"`
		} `json:"headers"`
		Origin string `json:"origin"`
		Url    string `json:"url"`
	} // JSON 解析

	response, _ := http.Get("https://httpbin.org/get")
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body) // 关闭连接

	body, _ := io.ReadAll(response.Body)
	var r JsonDec
	_ = json.Unmarshal(body, &r)
	fmt.Printf("body: %s\n\n", body)
}

// Get2 添加请求头
func Get2() {
	client := &http.Client{}

	newRequest, _ := http.NewRequest("GET", "https://httpbin.org/get", nil) // 获取请求信息
	newRequest.Header.Add("name", "loadkk")                                 // 添加请求头
	newRequest.Header.Add("age", "9527")                                    // 添加请求头
	response, _ := client.Do(newRequest)                                    // 发送请求
	body, _ := io.ReadAll(response.Body)
	fmt.Printf("Get2: %s \n", body)
}

// Get3 get? 后的参数拼接
func Get3() {
	params := url.Values{}
	params.Set("name", "loadkk")
	params.Set("age", "9527")
	Url, _ := url.Parse("https://httpbin.org/get")
	Url.RawQuery = params.Encode() // output age=9527&name=loadkk
	urlPath := Url.String()        // output https://httpbin.org/get?age=9527&name=loadkk
	response, _ := http.Get(urlPath)

	defer func(Body io.ReadCloser) { // 关闭连接
		_ = Body.Close()
	}(response.Body)

	body, _ := io.ReadAll(response.Body)
	fmt.Printf("Get3.body: %s\n", body)
}
