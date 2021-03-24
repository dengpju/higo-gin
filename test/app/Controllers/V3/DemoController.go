package V3

import (
	"fmt"
	"gitee.com/dengpju/higo-parameter/parameter"
	"github.com/dengpju/higo-annotation/annotation"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Exception"
	"github.com/dengpju/higo-gin/test/app/Models"
	"github.com/dengpju/higo-gin/test/app/Services"
	"github.com/dengpju/higo-ioc/injector"
	"github.com/dengpju/higo-throw/throw"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"log"
	"sync"
	"time"
)

type DemoController struct {
	Higo        *higo.Higo
	Age         *annotation.Value     `prefix:"user.age"`
	DemoService *Services.DemoService `inject:"Bean.DemoService()"`
	*higo.Gorm  `inject:"Bean.NewGorm()"`
	*redis.Pool `inject:"Bean.NewRedisPool()"`
}

type DemoController2 struct {
	Ttt string
}

func (this *DemoController) Self() higo.IClass {
	return this
}

var demoControllerOnce sync.Once
var demoControllerPointer *DemoController

func NewDemoController() *DemoController {
	demoControllerOnce.Do(func() {
		demoControllerPointer = &DemoController{}
		injector.BeanFactory.Apply(demoControllerPointer)
		injector.BeanFactory.Set(demoControllerPointer)
	})
	return demoControllerPointer
}

// 测试异常
func (this *DemoController) HttpsTestThrow(ctx *gin.Context) string {
	fmt.Println(ctx.Query("id"))
	fmt.Println(111)
	fmt.Println(&this)
	fmt.Println(this.Age.String())
	fmt.Println(this.Higo)
	var s []map[string]interface{}
	m1 := make(map[string]interface{})
	m1["jj"] = "m1jjj"
	m1["dd"] = "m1ddd"
	m2 := make(map[string]interface{})
	m2["jj"] = "m2jjj"
	m2["dd"] = "m2ddd"
	s = append(s, m1)
	s = append(s, m2)
	//测试自定义异常处理函数
	throw.Handle = func(p *parameter.Parameter) {
		if p.Name == throw.MESSAGE {
			throw.LogPayload.Msg = throw.ErrorToString(p.Value)
			throw.MapString.Put(p.Name, p.Value)
		}
	}
	Exception.BusinessException(throw.Code(2), throw.Message("v3 https 测试异常"), throw.Data(s))
	throw.Throw(throw.Message("v3 https 测试异常"), throw.Code(2), throw.Data(struct {
		Id   int
		Name string
	}{Id: 1, Name: "哦哦"}))
	return "v3 https_test_throw"
}

// 测试get请求
func (this *DemoController) HttpsTestGet(ctx *gin.Context) higo.Model {
	fmt.Println(injector.BeanFactory.Get(this))
	fmt.Println(this.DB)
	user := Models.NewUserModel()
	err := ctx.ShouldBindUri(user)
	if err != nil {
		log.Fatal("映射错误")
	}
	this.Table("ts_user").
		Where("id=?", 3).
		Find(user)
	higo.Task(this.TestTask, func() {
		this.TestTaskDone(3)
	}, user.Id)
	redisConn := this.Pool.Get()
	fmt.Println(redis.String(redisConn.Do("get", "name")))
	return user
}

// 测试post请求
func (this *DemoController) HttpsTestPost(ctx *gin.Context) string {
	return "v3 https_test_post"
}

// 测试异常
func (this *DemoController) HttpTestThrow(ctx *gin.Context) string {
	throw.Throw(throw.Message("v3 http 测试异常"), throw.Code(0))
	return "v3 http_test_throw"
}

// 测试get请求
func (this *DemoController) HttpTestGet(ctx *gin.Context) string {
	return "HttpTestGet"
}

// 测试post请求
func (this *DemoController) HttpTestPost(ctx *gin.Context) string {
	return "v3 http_test_post"
}

func (this *DemoController) Login(ctx *gin.Context) string {
	return "登录成功"
}

func (this *DemoController) TestTask(params ...interface{}) {
	time.Sleep(time.Second * 5)
	fmt.Println("测试task")
	fmt.Println(params)
}

func (this *DemoController) TestTaskDone(id int) {
	fmt.Println("测试task执行结束", id)
}
