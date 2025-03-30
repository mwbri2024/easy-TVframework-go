package utils

import (
	"bytes"
	"crypto/md5"
	"easy-itv/config"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var logger = log.Default()

// 获取本机信息
func GetIPInfo() (string, string, string) {
	// 发起请求
	client := &http.Client{
		Timeout: 30 * time.Second, // 设置超时时间为10秒
	}

	var resp *http.Response
	var err error
	maxRetries := 3 // 最大重试次数
	for i := 0; i < maxRetries; i++ {
		resp, err = client.Get("http://myip.ipip.net")
		if err == nil {
			defer resp.Body.Close()

			// 读取返回的内容
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return "异常", "异常", "异常"
			}

			// 将响应内容转换为字符串
			ipInfo := string(body)

			// 正则表达式提取 IP 地址
			ipRegex := regexp.MustCompile(`当前 IP：(\d+\.\d+\.\d+\.\d+)`)
			ipMatches := ipRegex.FindStringSubmatch(ipInfo)
			if len(ipMatches) < 2 {
				// fmt.Println("解析 IP 地址失败")
				return "未知", "未知", "未知"
			}

			// 获取 IP 地址
			ip := ipMatches[1]
			if ip == "" {
				// fmt.Println("IP 地址为空")
				return "未知", "未知", "未知"
			}

			// 查找 "来自于：" 之后的部分
			parts := strings.Split(ipInfo, "来自于：")
			if len(parts) < 2 {
				// fmt.Println("解析位置信息失败")
				return ip, "未知", "未知"
			}

			// 获取位置并去除多余的空格
			location := strings.TrimSpace(parts[1])
			if location == "" {
				location = "未知"
			}

			province := "未知"
			if key, _, found := GetFromList(location, config.ProvinceList); found {
				province = key
				config.Province = province
			}

			operator := "未知"
			if key, _, found := GetFromList(location, config.OperatorList); found {
				operator = key
				config.Operator = key
			}

			return ip, province, operator
		}
		fmt.Printf("%d/%d: %v\n", i+1, maxRetries, err)
		time.Sleep(2 * time.Second) // 每次重试间隔 2 秒
	}
	return "", "", ""
}

// 通用查询列表
func GetFromList(SearchKey string, SearchList map[string]string) (string, string, bool) {
	for key, code := range SearchList {
		if strings.Contains(SearchKey, key) {
			return key, code, true
		}
	}
	return "", "", false
}

// GetFormattedTime 获取当前时间并格式化为 "yyyy-MM-dd HH:mm:ss" 格式
func GetFormattedTime() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02 15:04:05")
}

// 获取 HTTP 请求的查询参数（query 参数）
func DefaultQuery(r *http.Request, name string, defaultValue string) string {
	param := r.URL.Query().Get(name)
	if param == "" {
		return defaultValue
	}
	return param
}

// sendHttpRequest 发送 HTTP 请求
func SendHttpRequest(requestUrl, method string, requestBody interface{}, requestHeader map[string]string, debugMode bool, isFormUrlEncoded bool) (int, string, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	var reqBody io.Reader
	if requestBody != nil {
		if isFormUrlEncoded {
			// 处理 x-www-form-urlencoded
			data := url.Values{}
			if bodyMap, ok := requestBody.(map[string]string); ok {
				for key, value := range bodyMap {
					data.Set(key, value)
				}
			}
			reqBody = strings.NewReader(data.Encode())
		} else {
			// 默认 JSON 处理
			jsonData, err := json.Marshal(requestBody)
			if err != nil {
				return 0, "", err
			}
			reqBody = bytes.NewBuffer([]byte(jsonData))
		}
	}

	// 创建请求
	req, err := http.NewRequest(strings.ToUpper(method), requestUrl, reqBody)
	if err != nil {
		logger.Printf("NewRequest Error: %v", err)
		return 0, "", err
	}

	// 设置请求头
	for key, value := range requestHeader {
		// req.Header.Set(key, value)
		req.Header[key] = []string{value} // 直接赋值，避免 net/http 规范化
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Printf("io.ReadAll Error: %v", err)
		return resp.StatusCode, "", err
	}

	if debugMode {
		fmt.Println("Request URL:", requestUrl)
		fmt.Println("Method:", method)
		fmt.Println("Request Header:", requestHeader)
		fmt.Println("Request Body:", requestBody)
		fmt.Println("HTTP Code:", resp.StatusCode)
		fmt.Println("Response Body:", string(respBody))
		fmt.Println("------------------------")
	}

	return resp.StatusCode, string(respBody), nil
}

// Md5Encrypt 生成字符串的 MD5 值
func Md5Encrypt(input string) string {
	hash := md5.New()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}

// Base64 解码函数
func DecodeBase64(encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
