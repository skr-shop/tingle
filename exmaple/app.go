package main

import (
	"net/http"

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

// TingleHandler Just for test
type TingleHandler struct {
	tingle.Next
}

// Do Ping
func (h *TingleHandler) Do(c *tingle.Context) error {
	c.JSON("Tingle!")

	return nil
}

func main() {
	t := tingle.NewWithDefaultMW()

	// 注册一个路由
	router := tingle.NewRouter()

	router.Handle(http.MethodGet, "/hello/world", &PingHandler{})
	router.Handle(http.MethodGet, "/hello/tingle", &TingleHandler{})
	t.SetRouter(router)

	t.Run(":4000")
}
