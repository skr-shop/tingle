package tingle

import (
	"strings"
)

type node struct {
	path     string  // 当前节点的路径
	indices  string  // 子节点路径第一个字符
	children []*node // 子节点
	handle   Handler // 当前节点处理函数
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// addRoute 增加路由
func (n *node) addRoute(path string, handle Handler) error {
	// 检查n节点是否为空节点
	if len(n.path) == 0 {
		return n.insertChild(path, handle)
	}

	for {
		// 查找最长公共前缀
		i := 0
		j := min(len(path), len(n.path))
		for i < j && path[i] == n.path[i] {
			i++
		}

		// 切割最长公共前缀
		if i < len(n.path) {
			n.children = []*node{
				&node{
					path:     n.path[i:],
					indices:  n.indices,
					children: n.children,
					handle:   n.handle,
				},
			}
			n.indices = string(n.path[i])
			n.path = path[:i]
			n.handle = nil
		}

		// 处理剩余path部分
		if i < len(path) {
			path = path[i:]

			c := path[0]
			if idx := strings.IndexByte(n.indices, c); idx != -1 {
				n = n.children[i]
				continue
			}

			// 处理indices部分
			n.indices += string(c)
			child := &node{}
			n.children = append(n.children, child)
			n = child

			return n.insertChild(path, handle)
		}
	}
}

// insertChild 插入新节点
func (n *node) insertChild(path string, handle Handler) error {
	n.path = path
	n.handle = handle

	return nil
}

// getValue 读取路由
func (n *node) getValue(path string) Handler {
	// 通过验证n.path是否师path的前缀来验证路由是否匹配
	for len(path) >= len(n.path) && path[:len(n.path)] == n.path {
		// 完全匹配
		if len(path) == len(n.path) {
			return n.handle
		}

		path = path[len(n.path):]

		// 验证indices
		c := path[0]
		if idx := strings.IndexByte(n.indices, c); idx != -1 {
			n = n.children[idx]
			continue
		}

		return nil
	}

	return nil
}
