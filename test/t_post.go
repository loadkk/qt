package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// PostForm Post 基本请求
func PostForm() {
	urlStr := "https://httpbin.org/post"
	data := url.Values{}
	data.Set("name", "loadkk")
	data.Set("age", "9527")
	response, _ := http.PostForm(urlStr, data) // 建立表单请求
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body) //读取buffer
	fmt.Printf("PostForm.body: %s\n", body)
}

// PostJson 发送 JSON 请求
func PostJson() {
	data := make(map[string]interface{})
	data["name"] = "loadkk"
	data["age"] = "9527"
	byteData, _ := json.Marshal(data) // 将json转化为字节数据
	response, _ := http.Post("https://httpbin.org/post", "application/html", bytes.NewReader(byteData))
	defer response.Body.Close() // 关闭连接
	body, _ := io.ReadAll(response.Body)
	fmt.Printf("PostJson.body: %s\n", body)
}
