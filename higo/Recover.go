package higo

import (
	"gitee.com/dengpju/higo-code/code"
	"github.com/dengpju/higo-logger/logger"
	"github.com/dengpju/higo-throw/exception"
	"github.com/dengpju/higo-utils/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

var (
	//Recover处理函数(可自定义替换)
	RecoverHandle RecoverFunc
	recoverOnce   sync.Once
)

func init() {
	recoverOnce.Do(func() {
		RecoverHandle = func(cxt *gin.Context, r interface{}) {

			//记录debug调用栈
			logger.LoggerStack(r, utils.GoroutineID())

			//封装通用json返回
			if h, ok := r.(gin.H); ok {
				cxt.JSON(http.StatusOK, h)
			} else if msg, ok := r.(*code.Code); ok {
				cxt.JSON(http.StatusOK, gin.H{
					"code":    msg.Code,
					"message": msg.Message,
					"data":    nil,
				})
			} else if MapString, ok := r.(utils.MapString); ok {
				cxt.JSON(http.StatusOK, MapString)
			} else if validate, ok := r.(*ValidateError); ok {
				cxt.JSON(http.StatusOK, gin.H{
					"code":    validate.Get().Code,
					"message": validate.Get().Message,
					"data":    nil,
				})
			} else {
				cxt.JSON(http.StatusOK, gin.H{
					"code":    0,
					"message": exception.ErrorToString(r),
					"data":    nil,
				})
			}
		}
	})
}

type RecoverFunc func(cxt *gin.Context, r interface{})

type Recover struct{}

func NewRecover() *Recover {
	return &Recover{}
}

func (this *Recover) Exception(hg *Higo) gin.HandlerFunc {
	return func(cxt *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				RecoverHandle(cxt, r)
				cxt.Abort()
			}
		}()
		cxt.Next()
	}
}
