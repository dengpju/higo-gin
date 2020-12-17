package V3

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-ioc/injector"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"reflect"
	"sync"
)

type RedisController struct {
	*higo.Gorm  `inject:"Bean.NewGorm()"`
	*redis.Pool `inject:"Bean.NewRedisPool()"`
}

var redisControllerOnce sync.Once
var RedisControllerPointer *RedisController

func NewRedisController() *RedisController {
	redisControllerOnce.Do(func() {
		RedisControllerPointer = &RedisController{}
		injector.BeanFactory.Apply(RedisControllerPointer)
		injector.BeanFactory.Set(RedisControllerPointer)
	})
	return RedisControllerPointer
}

func (this *RedisController) Reflection() (reflect.Type, reflect.Value) {
	return reflect.TypeOf(this), reflect.ValueOf(this)
}

func (this *RedisController) Test(ctx *gin.Context) string {
	redisConn := this.Pool.Get()
	v, _ := redis.String(redisConn.Do("get","name"))
	fmt.Println(redis.String(redisConn.Do("get","name")))
	return v
}