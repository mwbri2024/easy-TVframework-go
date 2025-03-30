package liveITV

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Itv struct{}

var logger = log.Default()

const MYSEPERETOR = "ThisIsMySeperator"

// 请求TS列表
func (i *Itv) HandleMainRequest(w http.ResponseWriter, r *http.Request, cdn string, id string, playseek string) {
	myid := strings.ReplaceAll(id, ".m3u8", "")
	matches := FindChannelsByTvgid(myid)
	if len(matches) == 0 {
		http.Error(w, "ChanneId not found!", http.StatusInternalServerError)
		return
	}

	Contentid := matches[0].Contentid
	Channecdn := matches[0].Cdn

	startUrl := "http://gslbserv.itv.cmvideo.cn:80/1.m3u8?channel-id=" + Channecdn + "&Contentid=" + Contentid + "&livemode=1&stbId=" + randomHexString(32)

	data, redirectURL, err := getHTTPResponse(startUrl, Channecdn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	redirectPrefix := redirectURL[:strings.LastIndex(redirectURL, "/")+1]

	if len(data) > 0 {
		// 替换TS文件的链接
		golang := "http://" + r.Host + r.URL.Path

		re := regexp.MustCompile(`((?i).*?\.ts)`)
		data = re.ReplaceAllStringFunc(data, func(match string) string {
			if strings.HasPrefix(match, "http") {
				return golang + "?ts=" + match
			} else {
				return golang + "?ts=" + redirectPrefix + match
			}
		})

		// 将&替换为自定义分隔符
		data = strings.ReplaceAll(data, "&", MYSEPERETOR)

		data = strings.ReplaceAll(data, "#EXTINF:10.0,", "#EXT-X-DISCONTINUITY#EXTINF:10.0,")

		w.Header().Set("Content-Disposition", "attachment;filename="+id)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
	} else {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// 请求TS数据
func (i *Itv) HandleTsRequest(w http.ResponseWriter, ts string) {
	// 将$替换回&
	ts = strings.ReplaceAll(ts, MYSEPERETOR, "&")

	// fmt.Println("ts=" + ts)

	// Read one piece and then write one piece
	w.Header().Set("Content-Type", "video/MP2T")
	_, _, err := handleTsHTTPResponse(ts, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// TS列表 请求方法
func getHTTPResponse(requestURL string, channecdn string) (string, string, error) {
	resp, err := FetchWithCustomResolver(requestURL)
	if err != nil || resp.StatusCode != 200 {
		maxRetries := 3
		retryDelay := 0 * time.Second

		for i := 0; i < maxRetries; i++ {
			resp, err = FetchWithCustomResolver(requestURL)
			if err != nil {
				if Debug {
					logger.Printf("GetM3U-%d Error: %v", i+2, err)
				}
			} else if resp.StatusCode == 200 {
				// 如果状态码是200，说明请求成功，退出循环
				break
			}

			if Debug {
				if resp != nil {
					logger.Printf("GetM3U-%d Code:%d", i+2, resp.StatusCode)
				} else {
					logger.Printf("GetM3U-%d Response is nil", i+2)
				}
			}

			// 如果响应状态码不是200，记录失败并增加失败计数
			switch {
			case strings.Contains(channecdn, "bestzb"):
				IncreaseFail(WorkIP_bestzb, IPList_bestzb)

			case strings.Contains(channecdn, "hnbblive"):
				IncreaseFail(WorkIP_hnbblive, IPList_hnbblive)

			case strings.Contains(channecdn, "FifastbLive"):
				IncreaseFail(WorkIP_fifalive, IPList_fifalive)

			}

			// 等待指定时间后再重试
			time.Sleep(retryDelay)
		}
	}

	if err != nil || resp.StatusCode != 200 {
		logger.Printf("GetM3U Failed")
		return "", "", err
	} else {
		redirectURL := resp.Header.Get("Location")
		if redirectURL == "" {
			redirectURL = requestURL
		}

		body, err := ReadResponseBody(resp)

		return body, redirectURL, err
	}
}

// TS数据 请求方法
func handleTsHTTPResponse(requestURL string, w http.ResponseWriter) (string, string, error) {
	dialer := &net.Dialer{
		Timeout: 10 * time.Second,
	}

	// 自定义resolver
	resolver := net.Resolver{
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			match, _ := regexp.MatchString(`^cache\.ott\..*\.itv\.cmvideo\.cn:80$`, address)
			if match {
				switch {
				case strings.Contains(address, "bestlive"):
					address = WorkIP_bestzb
				case strings.Contains(address, "hnbblive"):
					address = WorkIP_hnbblive
				case strings.Contains(address, "fifalive"):
					address = WorkIP_fifalive
				}

				if Debug {
					logger.Printf("TS:UseIP->%s", strings.Split(address, ":")[0])
				}
			}

			return dialer.DialContext(ctx, network, address)
		},
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
		Transport: &http.Transport{
			DialContext: resolver.Dial,
		},
	}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	resp, err := client.Do(req)

	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Printf("TS Request Fail:%d", resp.StatusCode)
		return "", "", http.ErrServerClosed
	}

	redirectURL := resp.Header.Get("Location")
	if redirectURL == "" {
		redirectURL = requestURL
	}

	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, resp.Body)

	if err != nil {
		return "", "", nil
	}

	return "", redirectURL, nil
}

// 使用自定义 DNS 解析器的 HTTP 请求方法，支持 302 重定向失败后重试 3 次
func FetchWithCustomResolver(requestURL string) (*http.Response, error) {
	isDebug := false
	if Debug && isDebug {
		logger.Printf("RequestURL: %s", requestURL)
	}

	dialer := &net.Dialer{
		Timeout: 10 * time.Second,
	}

	resolver := &net.Resolver{
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			match, _ := regexp.MatchString(`^cache\.ott\..*\.itv\.cmvideo\.cn:80$`, address)
			if match {
				switch {
				case strings.Contains(address, "bestlive"):
					address = GetBestIP(IPList_bestzb).Address + ":80"
					WorkIP_bestzb = address
				case strings.Contains(address, "hnbblive"):
					address = GetBestIP(IPList_hnbblive).Address + ":80"
					WorkIP_hnbblive = address
				case strings.Contains(address, "fifalive"):
					address = GetBestIP(IPList_fifalive).Address + ":80"
					WorkIP_fifalive = address
				}

				if Debug {
					logger.Printf("M3U:GetBestIP->%s", strings.Split(address, ":")[0])
				}
			}

			return dialer.DialContext(ctx, network, address)
		},
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if Debug && isDebug && len(via) > 0 {
				logger.Printf("Redirected to: %s", req.URL.String())
			}
			return http.ErrUseLastResponse // 禁止自动跳转
		},
		Transport: &http.Transport{
			DialContext: resolver.Dial,
		},
	}

	var lastResp *http.Response
	redirectCount := 0 // 记录 302 跳转次数

	// 重试
	for attempt := 1; attempt <= 3; attempt++ {
		req, err := http.NewRequest("GET", requestURL, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}

		lastResp, err = client.Do(req)
		if err != nil {
			if Debug && isDebug {
				logger.Printf("请求失败: %v (尝试 %d/3)", err, attempt)
			}
			time.Sleep(1 * time.Second) // 等待重试
			continue
		}

		// 处理 302/301 重定向
		for lastResp.StatusCode == http.StatusFound || lastResp.StatusCode == http.StatusMovedPermanently {
			redirectURL := lastResp.Header.Get("Location")
			if redirectURL == "" {
				return nil, fmt.Errorf("error: no Location header in response (attempt %d)", attempt)
			}

			redirectCount++
			if redirectCount > 20 {
				return nil, fmt.Errorf("too many redirects (max 6)")
			}

			if Debug && isDebug {
				logger.Printf("302 跳转 (%d): %s -> %s", redirectCount, requestURL, redirectURL)
			}

			requestURL = redirectURL // 更新 URL 继续请求
			req, err = http.NewRequest("GET", requestURL, nil)
			if err != nil {
				return nil, fmt.Errorf("error creating request after redirect: %w", err)
			}
			lastResp, err = client.Do(req)
			if err != nil {
				if Debug {
					logger.Printf("请求失败: %v (尝试 %d/3)", err, attempt)
				}
				time.Sleep(1 * time.Second)
				continue
			}
		}

		// 响应状态码不是 200 也重试
		if lastResp.StatusCode != http.StatusOK {
			if Debug && isDebug {
				logger.Printf("请求返回非200状态码: %d (尝试 %d/3)", lastResp.StatusCode, attempt)
			}
			continue
		}

		// 请求成功，返回响应
		return lastResp, nil
	}

	return nil, fmt.Errorf("exceeded maximum retries")
}

// 获取请求返回内容
func ReadResponseBody(resp *http.Response) (string, error) {
	var builder strings.Builder
	_, err := io.Copy(&builder, resp.Body)
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}

// 查找所有匹配 Tvgid 的记录
func FindChannelsByTvgid(tvgid string) []ChannelRecord {
	var results []ChannelRecord
	for _, channel := range ChannelList {
		if channel.Tvgid == tvgid {
			results = append(results, channel)
		}
	}
	return results
}

// 生成指定位数的随机字符
func randomHexString(length int) string {
	bytes := make([]byte, length/2)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
