package tingle

import (
	"net/http"
	"sync"
)

const (
	// DefalutPort 默认端口
	DefalutPort = "8088"
)

// Tingle Golang Framework
// 名称的灵感来自于《蜘蛛侠》中的 “peter tingle”
type Tingle struct {
	router                   *Router
	logger                   *Logger
	server                   *http.Server
	contextPool              *sync.Pool
	commonMiddlewares        []Handler
	beforeStartupMiddlewares []BeforeStartupHandler
}

// handle 注册用户路由请求
// method http method
// path http path
// handlerFunc UserHandlerFuncr
// beforeStartupHandler BeforeStartupHandler
func (tingle *Tingle) handle(method string, path string, handlerFunc HandlerFunc, bsHandlers ...BeforeStartupHandler) {
	tingle.router.Add(method, path, handlerFunc)
	// 注册启动前中间件
	tingle.RegisterBeforeStartupMW(bsHandlers...)
}

// Get 注册用户路由请求
// method http method
// path http path
// handlerFunc UserHandlerFuncr
// beforeStartupHandler BeforeStartupHandler
func (tingle *Tingle) Get(path string, handlerFunc HandlerFunc, bsHandlers ...BeforeStartupHandler) {
	tingle.handle(http.MethodGet, path, handlerFunc, bsHandlers...)
}

// Post 注册用户路由请求
// method http method
// path http path
// handlerFunc UserHandlerFuncr
// beforeStartupHandler BeforeStartupHandler
func (tingle *Tingle) Post(path string, handlerFunc HandlerFunc, bsHandlers ...BeforeStartupHandler) {
	tingle.handle(http.MethodPost, path, handlerFunc, bsHandlers...)
}

// Put 注册用户路由请求
// method http method
// path http path
// handlerFunc UserHandlerFuncr
// beforeStartupHandler BeforeStartupHandler
func (tingle *Tingle) Put(path string, handlerFunc HandlerFunc, bsHandlers ...BeforeStartupHandler) {
	tingle.handle(http.MethodPut, path, handlerFunc, bsHandlers...)
}

// Delete 注册用户路由请求
// method http method
// path http path
// handlerFunc UserHandlerFuncr
// beforeStartupHandler BeforeStartupHandler
func (tingle *Tingle) Delete(path string, handlerFunc HandlerFunc, bsHandlers ...BeforeStartupHandler) {
	tingle.handle(http.MethodDelete, path, handlerFunc, bsHandlers...)
}

// Patch 注册用户路由请求
// method http method
// path http path
// handlerFunc UserHandlerFuncr
// beforeStartupHandler BeforeStartupHandler
func (tingle *Tingle) Patch(path string, handlerFunc HandlerFunc, bsHandlers ...BeforeStartupHandler) {
	tingle.handle(http.MethodPatch, path, handlerFunc, bsHandlers...)
}

// Head 注册用户路由请求
// method http method
// path http path
// handlerFunc UserHandlerFuncr
// beforeStartupHandler BeforeStartupHandler
func (tingle *Tingle) Head(path string, handlerFunc HandlerFunc, bsHandlers ...BeforeStartupHandler) {
	tingle.handle(http.MethodHead, path, handlerFunc, bsHandlers...)
}

// Options 注册用户路由请求
// method http method
// path http path
// handlerFunc UserHandlerFuncr
// beforeStartupHandler BeforeStartupHandler
func (tingle *Tingle) Options(path string, handlerFunc HandlerFunc, bsHandlers ...BeforeStartupHandler) {
	tingle.handle(http.MethodOptions, path, handlerFunc, bsHandlers...)
}

// Trace 注册用户路由请求
// method http method
// path http path
// handlerFunc UserHandlerFuncr
// beforeStartupHandler BeforeStartupHandler
func (tingle *Tingle) Trace(path string, handlerFunc HandlerFunc, bsHandlers ...BeforeStartupHandler) {
	tingle.handle(http.MethodTrace, path, handlerFunc, bsHandlers...)
}

// Connect 注册用户路由请求
// method http method
// path http path
// handlerFunc UserHandlerFuncr
// beforeStartupHandler BeforeStartupHandler
func (tingle *Tingle) Connect(path string, handlerFunc HandlerFunc, bsHandlers ...BeforeStartupHandler) {
	tingle.handle(http.MethodConnect, path, handlerFunc, bsHandlers...)
}

// RegisterCommonMW 注册公共中间件
func (tingle *Tingle) RegisterCommonMW(handlers ...Handler) {
	tingle.commonMiddlewares = append(tingle.commonMiddlewares, handlers...)
}

// RegisterBeforeStartupMW 注册启动前中间件
func (tingle *Tingle) RegisterBeforeStartupMW(handler ...BeforeStartupHandler) {
	tingle.beforeStartupMiddlewares = append(tingle.beforeStartupMiddlewares, handler...)
}

// StartupBeforeStartupMW 启动启动前中间件
func (tingle *Tingle) StartupBeforeStartupMW() {
	if len(tingle.beforeStartupMiddlewares) == 0 {
		return
	}

	// 启动
	for _, bsHandler := range tingle.beforeStartupMiddlewares {
		go func(t *Tingle, bh BeforeStartupHandler) {
			// todo recover

			// 坑，这里不能写 bsHandler(t)，因为是并发的bsHandler可能被其他gouroutine修改
			bh(t)
		}(tingle, bsHandler)
	}
}

// Run 启动框架
func (tingle *Tingle) Run(addr string) {
	// 启动启动前中间件
	tingle.StartupBeforeStartupMW()

	// 启动服务
	if addr == "" {
		addr = DefalutPort
	}
	tingle.server.Addr = addr
	tingle.server.Handler = tingle
	if err := tingle.server.ListenAndServe(); err != nil {
		panic(err)
	}
}

// ServeHTTP 实现http.handler接口
func (tingle *Tingle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := tingle.contextPool.Get().(*Context)
	context.Request = r
	context.Response = w
	tingle.handleHTTPRequest(context)
}

// handleHTTPRequest 执行http请求
func (tingle *Tingle) handleHTTPRequest(context *Context) {
	handler := tingle.router.GetHandler(context.Request.Method, context.Request.URL.Path)
	if handler == nil {
		context.Response.WriteHeader(404)
		return
	}

	// 执行中间件
	var nullHandler Handler
	if len(tingle.commonMiddlewares) == 0 {
		// todo
		panic("commonMiddlewares is empty")
	}
	// 责任链，链式调用
	for k, handler := range tingle.commonMiddlewares {
		if k == 0 {
			nullHandler = handler
			continue
		}
		tingle.commonMiddlewares[k-1].SetNext(handler)
	}
	nullHandler.Run(&Context{})

	// 执行用户注册的handler
	handler.Do(context)
}

// New 创建Tingle框架实例
func New() *Tingle {
	t := &Tingle{
		router: &Router{
			Trees: make(map[string]*node),
		},
		logger:      new(Logger),
		server:      &http.Server{},
		contextPool: new(sync.Pool),
	}
	t.contextPool.New = func() interface{} {
		return new(Context)
	}
	return t
}

// NewWithDefaultMW 创建Tingle框架实例并注册默认的中间件
// 1. 默认注册goroutine panic recover中间件
// 2. 默认注册请求访问日志(access log)中间件
func NewWithDefaultMW() *Tingle {
	t := New()
	t.RegisterCommonMW(
		&NullHandler{},
		&RecoverHandler{},
		&AccessLogHandler{})
	return t
}
