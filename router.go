package tingle

// HandlerFunc 注册路由时的闭包
type HandlerFunc func(c *Context) error

// BeforeStartupHandler 前置**启动**handler
// 可以定义一些前置(同步或异步任务或同步+异步)
// 比如异步更新内存缓存
type BeforeStartupHandler func(tingle *Tingle) error

// BeforeRequestHandler 前置请求handler
// 可以定义一些接口请求的前置逻辑(同步或异步任务或同步+异步)
// 比如校验用户是否登陆逻辑
type BeforeRequestHandler func(c *Context) error

// AfterRequestHandler 后置请求handler
// 可以定义一些接口请求的后置逻辑(同步或异步任务或同步+异步) 比如对一致性要求不高的 异步刷新缓存到db
type AfterRequestHandler func(c *Context) error

// tree 路由树
type tree struct {
	Method     string
	Path       string
	UserHandle Handler
}

// Router 路由结构体
type Router struct {
	Trees                map[string]*node
	BeforeStartupHandles []BeforeStartupHandler
	BeforeRequestHandles []BeforeRequestHandler
	AfterRequestHandles  []AfterRequestHandler
}

// Add 绑定路由
func (router *Router) Add(method string, path string, handlerFunc HandlerFunc) {
	root := router.Trees[method]
	if root == nil {
		root = new(node)
		router.Trees[method] = root
	}

	root.addRoute(path, &TemplateHandler{
		handlerFunc: handlerFunc,
	})
}

// GetHandler 读取路由
func (router *Router) GetHandler(method, path string) Handler {
	t := router.Trees[method]
	if t != nil {
		return t.getValue(path)
	}

	return nil
}

// TemplateHandler 注册路由的模板handler
type TemplateHandler struct {
	Next
	handlerFunc HandlerFunc
}

// Do 模板handler执行注册路由时的业务闭包
func (h *TemplateHandler) Do(c *Context) error {
	return h.handlerFunc(c)
}
