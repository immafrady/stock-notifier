package utils

import (
	"io"
	"log"
	"net/http"
)

// Request 简单的请求
func Request(url string) string {
	resp, err := http.Get(url)

	if err != nil {
		log.Println("[error]请求失败:", err)
		return ""
	}
	defer resp.Body.Close()

	// 读取响应内容
	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[error]读取响应失败:", err)
		return ""
	}

	return string(body)
}
