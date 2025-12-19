package tool

import (
	"strings"
	"sync"
)

type Ignore struct {
	config []string
	cache  sync.Map
}

func NewIgnore(config []string) *Ignore {
	return &Ignore{
		config: config,
	}
}

// Exists 检查参数 s 是否在名单中，在名单中返回 true，否则返回 false。
func (ig *Ignore) Exists(s string) bool {
	value, found := ig.cache.Load(s)
	if found {
		return value.(bool)
	}

	var onlist = false
	for _, v := range ig.config {
		if v[len(v)-1] == '*' {
			if strings.HasPrefix(s, v[0:len(v)-1]) {
				onlist = true
				break
			}
		} else {
			if s == v {
				onlist = true
				break
			}
		}
	}
	ig.cache.Store(s, onlist)
	return onlist
}

func (ig *Ignore) Config() []string {
	return ig.config
}
