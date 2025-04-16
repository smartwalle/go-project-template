package pkg

import (
	"strings"
	"sync"
)

type Ignorelist struct {
	config []string
	cache  sync.Map
}

func NewIgnorelist(config []string) *Ignorelist {
	return &Ignorelist{
		config: config,
	}
}

// Exists 检查参数 s 是否在名单中，在名单中返回 true，否则返回 false。
func (ignore *Ignorelist) Exists(s string) bool {
	value, found := ignore.cache.Load(s)
	if found {
		return value.(bool)
	}

	var onlist = false
	for _, v := range ignore.config {
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
	ignore.cache.Store(s, onlist)
	return onlist
}

func (ignore *Ignorelist) Config() []string {
	return ignore.config
}
