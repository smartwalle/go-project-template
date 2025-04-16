package pkg

import (
	"strings"
	"sync"
)

type Whitelist struct {
	config []string
	cache  sync.Map
}

func NewWhitelist(config []string) *Whitelist {
	return &Whitelist{
		config: config,
	}
}

// Check 检查参数 uri 是否在白名单中，是返回 true，否则返回 false。
func (w *Whitelist) Check(uri string) bool {
	value, found := w.cache.Load(uri)
	if found {
		return value.(bool)
	}

	var onlist = false
	for _, v := range w.config {
		if v[len(v)-1] == '*' {
			if strings.HasPrefix(uri, v[0:len(v)-1]) {
				onlist = true
				break
			}
		} else {
			if uri == v {
				onlist = true
				break
			}
		}
	}
	w.cache.Store(uri, onlist)
	return onlist
}

func (w *Whitelist) Config() []string {
	return w.config
}
