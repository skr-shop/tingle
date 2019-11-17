// bee run -main=demo.go

package main

import (
	"fmt"
	"time"

	"github.com/skr-shop/tingle"
)

func main() {
	// 创建一个注册了默认路由的tingle实例
	t := tingle.NewWithDefaultMW()

	// 注册一个路由
	t.Get(
		"/ping",
		// 注册业务
		func(c *tingle.Context) error {
			// 输出Json响应内容
			return c.JSON("Pong!")
		},
		// 注册启动前中间件
		func(t *tingle.Tingle) error {
			ticker := time.NewTicker(1 * time.Second)
			for {
				select {
				case t := <-ticker.C:
					fmt.Println(tingle.FormatTimeToStr(&t), "API Ping Caching")
				}
			}
		})

	// 注册一个路由
	t.Get(
		"/pong",
		// 注册业务
		func(c *tingle.Context) error {
			// 输出Json响应内容
			return c.JSON("Boom!")
		},
		// 注册启动前中间件
		func(t *tingle.Tingle) error {
			ticker := time.NewTicker(1 * time.Second)
			for {
				select {
				case t := <-ticker.C:
					fmt.Println(tingle.FormatTimeToStr(&t), "API Pong Caching")
				}
			}
		})

	// 启动tingle服务
	t.Run(":4000")
}
