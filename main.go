package main

import (
	"easy-itv/config"
	"easy-itv/list"
	"easy-itv/liveITV"
	"easy-itv/livePhoenix"
	"easy-itv/update"
	"easy-itv/utils"
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var logger = log.Default()

func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()        // 不使用 gin.Default()，这样不会自动加载 Logger 和 Recovery 中间件
	r.Use(gin.Recovery()) // 恢复中间件，不使用 gin.Default()时必须手动加载

	r.HEAD("/", func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, utils.GetFormattedTime())
	})

	r.GET("/tv.m3u", func(c *gin.Context) {

		clientIP := c.ClientIP()
		urltoken := utils.DefaultQuery(c.Request, "token", "")

		if config.Token == "" || urltoken == config.Token {
			itvm3uobj := &list.Tvm3u{}
			if !config.Debug {
				c.Writer.Header().Set("Content-Type", "application/octet-stream")
				c.Writer.Header().Set("Content-Disposition", "attachment; filename=tv.m3u")
			}

			itvm3uobj.GetTvM3u(c.Writer, c.Request.Host)
		} else {
			logger.Printf("Token error:%s clientIP:%s GetM3U", urltoken, clientIP)
			c.String(http.StatusUnauthorized, "token error")
		}

	})

	r.GET("/:path/:rid", func(c *gin.Context) {
		clientIP := c.ClientIP()
		path := c.Param("path")
		rid := c.Param("rid")
		ts := utils.DefaultQuery(c.Request, "ts", "")
		playtoken := utils.DefaultQuery(c.Request, "token", "")

		// token 验证
		if ts == "" && (config.Token == "" || playtoken != config.Token) {
			logger.Printf("Token error:%s clientIP:%s VIEW:%s", playtoken, clientIP, rid)
			c.String(http.StatusUnauthorized, "token error")
			return
		}

		switch path {
		case "TVOD":
			if !liveITV.Enable {
				c.String(http.StatusNotFound, "err!")
				return
			}

			itvobj := &liveITV.Itv{}
			cdn := utils.DefaultQuery(c.Request, "cdn", "")
			playseek := utils.DefaultQuery(c.Request, "playseek", "")
			ind := strings.Index(cdn, "?playseek=")

			//rebuild playseek and cdn
			if ind > -1 {
				playseek = cdn[ind+10:]
				cdn = cdn[0:ind]
			}
			if ts == "" {
				// 获取TS列表
				logger.Printf("Client IP:%s VIEW:%s", clientIP, rid)
				itvobj.HandleMainRequest(c.Writer, c.Request, cdn, rid, playseek)
			} else {
				// 获取TS数据
				itvobj.HandleTsRequest(c.Writer, ts)
			}

		case "phoenix":
			if !livePhoenix.Enable {
				c.String(http.StatusNotFound, "err!")
				return
			}

			logger.Printf("Client IP:%s VIEW:%s", clientIP, rid)
			// phxvobj := &livePhoenix.Phoenix{}
			// phxvobj.HandleMainRequest(c.Writer, c.Request, rid)
		}
	})
	return r
}

func main() {
	fmt.Println("Version:", config.Version)

	// 接收token参数，允许为空
	flag.StringVar(&config.Token, "token", "", "Your token")
	flag.Parse()

	// 验证
	if config.Token != "" {
		// 正则表达式：只允许大小写字母和数字
		match, err := regexp.MatchString("^[a-zA-Z0-9]+$", config.Token)
		if err != nil {
			log.Fatal("Error in regex match:", err)
		}

		// 如果不匹配，打印错误信息并退出
		if !match {
			log.Fatal("Invalid token: Token can only contain letters and numbers")
		}

		fmt.Println("Local Token:", config.Token)
	}

	// 获取运行环境IP地址
	ip, province, operator := utils.GetIPInfo()
	fmt.Println("Local IPaddress:", ip)
	fmt.Println("Local Province:", province)
	fmt.Println("Local Operator:", operator)

	// 启动一个异步任务 例：获取全网通IP、YSP初始化缓存、……
	go func() {
		// do some thing

		fmt.Println("AsyncTask:Done")
	}()

	// 创建通道用于停止定时任务
	done1 := make(chan bool)
	done2 := make(chan bool)

	// 启动定时任务
	if config.Update {
		go timedFunction1(done1, 1*time.Hour) //运行间隔 1小时
	}
	go timedFunction2(done2, 30*time.Minute) //运行间隔 30分钟

	fmt.Println("Started:", utils.GetFormattedTime())

	// 启动WEB服务
	r := setupRouter()
	r.Run(":" + config.Port)

	// 发送停止信号
	done1 <- true
	done2 <- true
}

// 定时任务1 自动更新
func timedFunction1(done <-chan bool, interval time.Duration) {
	// 创建一个定时器，每隔 interval 触发一次
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			// 收到停止信号，退出函数
			return
		case <-ticker.C:
			if config.Debug {
				logger.Printf("检查更新...")
			} else {
				logger.Printf("......")
			}
			latest, err := update.GetLatestVersionInfo()
			if err != nil {
				if config.Debug {
					logger.Printf("获取最新版本失败:%v", err)
				}
				continue
			}

			// 检查是否需要更新
			if latest.Version <= config.Version {
				if config.Debug {
					logger.Printf("当前已是最新版本:%s", config.Version)
				}
				continue
			}

			if config.Debug {
				logger.Printf("发现新版本: %s ，正在下载...", latest.Version)
			}
			// 获取平台哈希值
			platformBinary := update.GetPlatformBinaryName(latest.Version)
			expectedHash, exists := latest.Hash[platformBinary]
			if !exists {
				if config.Debug {
					logger.Printf("未找到该平台的哈希值:%s", platformBinary)
				}
				continue
			}

			newBinary, err := update.DownloadNewBinary(latest.Version)
			if err != nil {
				if config.Debug {
					logger.Printf("下载失败:%v", err)
				}
				continue
			}

			// 替换并重启
			if config.Debug {
				logger.Printf("开始更新...")
			}
			err = update.ReplaceAndRestart(newBinary, config.Version, expectedHash)
			if err != nil {
				if config.Debug {
					logger.Printf("更新失败:%v", err)
				}
			}

		}
	}
}

// 定时任务2 例：更新全网通IP、更新YSP缓存、……
func timedFunction2(done <-chan bool, interval time.Duration) {
	// 创建一个定时器，每隔 interval 触发一次
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			// 收到停止信号，退出函数
			return
		case <-ticker.C:
			// do some thing
			// liveITV.UpdateIPList(ipRecords, &liveITV.IPList_bestzb)

			logger.Println("Scheduled task has completed.")
		}
	}
}
