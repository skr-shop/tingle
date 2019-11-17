package main

import (
	"github.com/skr-shop/tingle"
)

func main() {
	// 创建一个注册了默认路由的tingle实例
	t := tingle.NewWithDefaultMW()

	// 注册一个路由
	t.Handle("get", "/ping", func(c *tingle.Context) error {
		// 输出Json响应内容
		c.JSON("Pong!")
		return nil
	})

	// 启动tingle服务
	t.Run(":4000")
}
