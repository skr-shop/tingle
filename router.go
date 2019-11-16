package tingle

type Router struct {
	trees map[string]*node
}

// NewRouter 创建路由
func NewRouter() *Router {
	return &Router{
		trees: map[string]*node{},
	}
}

// Handle 为路由增加处理函数
func (r *Router) Handle(method, path string, handle Handler) {
	root := r.trees[method]
	if root == nil {
		root = new(node)
		r.trees[method] = root
	}

	root.addRoute(path, handle)
}

// getHandlers 读取路由
func (r *Router) getHandlers(method, path string) Handler {
	t := r.trees[method]
	if t != nil {
		return t.getValue(path)
	}

	return nil
}
