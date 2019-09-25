package main

import (
	"github.com/skr-shop/tingle"
)

// PingHandler Just for ping
type PingHandler struct {
	tingle.Next
}

// Do Ping
func (h *PingHandler) Do(c *tingle.Context) error {
	c.JSON("Pong!")

	return nil
}

func main() {
	t := tingle.NewWithDefaultMW()

	// 注册一个路由
	t.Handle("get", "/ping", &PingHandler{})

	t.Run(":4000")
}
