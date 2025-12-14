package middleware

import (
	"github.com/gin-gonic/gin"
	"project-template/pkg"
)

var (
	AuthIgnore *pkg.Ignore // 不用校验登录状态的接口
)

func init() {
	AuthIgnore = pkg.NewIgnore([]string{
		"/api/signin",
		"/api/signup",
	})
}

func AuthCheck() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 请求接口在“不需要校验登录状态接口名单中”，则不进行登录状态的验证
		if AuthIgnore.Exists(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 验证登录状态

		c.Next()
	}
}
