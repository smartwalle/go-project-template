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

// Exists 检查参数 uri 是否在名单中，是返回 true，否则返回 false。
func (ignore *Ignorelist) Exists(uri string) bool {
	value, found := ignore.cache.Load(uri)
	if found {
		return value.(bool)
	}

	var onlist = false
	for _, v := range ignore.config {
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
	ignore.cache.Store(uri, onlist)
	return onlist
}

func (ignore *Ignorelist) Config() []string {
	return ignore.config
}
