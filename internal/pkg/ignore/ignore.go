package ignore

import (
	"strings"
	"sync"
)

type Ignore struct {
	config []string
	cache  sync.Map
}

func New(config []string) *Ignore {
	return &Ignore{
		config: config,
	}
}

// Match 检查参数 s 是否在名单中，在名单中返回 true，否则返回 false。
func (ig *Ignore) Match(s string) bool {
	value, found := ig.cache.Load(s)
	if found {
		return value.(bool)
	}

	var onlist = false
	for _, elem := range ig.config {
		if elem[len(elem)-1] == '*' {
			if strings.HasPrefix(s, elem[0:len(elem)-1]) {
				onlist = true
				break
			}
		} else if s == elem {
			onlist = true
			break
		}
	}
	ig.cache.Store(s, onlist)
	return onlist
}

func (ig *Ignore) Config() []string {
	return ig.config
}
