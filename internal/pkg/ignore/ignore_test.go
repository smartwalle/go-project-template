package ignore_test

import (
	"project-template/internal/pkg/ignore"
	"testing"
)

func TestWhitelist(t *testing.T) {
	var whitelist = ignore.New([]string{"1", "2", "3"})
	t.Log("在白名单上：", whitelist.Match("1"))
	t.Log("在白名单上：", whitelist.Match("2"))
	t.Log("在白名单上：", whitelist.Match("4"))
}

func TestBlacklist(t *testing.T) {
	var blacklist = ignore.New([]string{"1", "2", "3"})
	t.Log("在黑名单上：", blacklist.Match("1"))
	t.Log("在黑名单上：", blacklist.Match("2"))
	t.Log("在黑名单上：", blacklist.Match("4"))
}
