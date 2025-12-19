package tool_test

import (
	"project-template/internal/tool"
	"testing"
)

func TestWhitelist(t *testing.T) {
	var whitelist = tool.NewIgnore([]string{"1", "2", "3"})
	t.Log("在白名单上：", whitelist.Exists("1"))
	t.Log("在白名单上：", whitelist.Exists("2"))
	t.Log("在白名单上：", whitelist.Exists("4"))
}

func TestBlacklist(t *testing.T) {
	var blacklist = tool.NewIgnore([]string{"1", "2", "3"})
	t.Log("在黑名单上：", blacklist.Exists("1"))
	t.Log("在黑名单上：", blacklist.Exists("2"))
	t.Log("在黑名单上：", blacklist.Exists("4"))
}
