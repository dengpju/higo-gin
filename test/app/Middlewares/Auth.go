package Middlewares

import (
	"gitee.com/dengpju/higo-code/code"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Consts"
	"github.com/dengpju/higo-throw/exception"
	"github.com/gin-gonic/gin"
)

// 鉴权
type Auth struct{}

// 构造函数
func NewAuth() *Auth {
	return &Auth{}
}

func (this Auth) Loader(hg *higo.Higo) gin.HandlerFunc {
	return func(c *gin.Context) {
		if route, ok := hg.GetRoute(c.Request.URL.Path); ok {
			// TODO::非静态页面需要鉴权
			if !higo.IsNotAuth(route.Flag()) && !route.IsStatic() {
				if "" == c.GetHeader("X-Token") {
					exception.Throw(exception.Message(code.Message(Consts.INVALID_TOKEN).Message), exception.Code(code.Message(Consts.INVALID_TOKEN).Code))
				}
			}
			c.Next()
		} else {
			exception.Throw(exception.Message(code.Message(Consts.INVALID_API).Message), exception.Code(code.Message(Consts.INVALID_API).Code))
		}
	}
}
