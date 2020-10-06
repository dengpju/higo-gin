package higo

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	hg *Higo
	// 系统类型 Windows Or Linux
	SysType string
	// 路径分隔符
	PathSeparator string
	// ssl 证书
	SslOut, SslCrt, SslKey string
)

// 上下文
type IHiContext interface {
	OnRequest(*gin.Context) error
	OnResponse(result interface{}) (interface{}, error)
}

// http 服务结构体
type Hse struct {
	Config string
	Router IRouterLoader
	Serve  string
}

type Higo struct {
	*gin.Engine
	g            *gin.RouterGroup
	eg           errgroup.Group
	exprData     map[string]interface{}
	currentGroup string
	root         string
	containers   *Containers
	isAutoSsl    bool
	middle       []*IMiddleware
	serve        []*Hse
}

// 初始化
func Init() *Higo {
	hg = &Higo{
		Engine:     gin.New(),
		exprData:   map[string]interface{}{},
		containers: NewContainer(),
		middle:     make([]*IMiddleware, 0),
		serve:      make([]*Hse, 0),
	}

	// 全局异常
	hg.Engine.Use(NewRecover().RuntimeException(hg))
	// 系统类型
	SysType = runtime.GOOS
	// 初始分隔符
	if SysType == "windows" {
		PathSeparator = "\\"
	} else {
		PathSeparator = "/"
	}
	// 是否使用自带ssl测试https
	hg.isAutoSsl = false

	return hg
}

// 设置主目录
func (this *Higo) SetRoot(root string) *Higo {
	this.root = root
	return this
}

// 获取主目录
func (this *Higo) GetRoot() string {
	return utils.If(this.root == "", ROOT, this.root).(string)
}

// 配置
func (this *Higo) config() *Higo {
	// 获取主目录
	root := hg.GetRoot()
	// runtime目录
	runtimeDir := root + "runtime"
	if _, err := os.Stat(runtimeDir); os.IsNotExist(err) {
		if os.Mkdir(runtimeDir, os.ModePerm) != nil {}
	}
	// 日志
	Log(root)
	// 装载配置
	confDir := root + "conf"
	filepathErr := filepath.Walk(confDir,
		func(p string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				return nil
			}
			if path.Ext(p) == ".yaml" {
				fmt.Println("yaml file:", filepath.Base(p))
				yamlFile, _ := ioutil.ReadFile(p)
				yamlFileErr := yaml.Unmarshal(yamlFile, &Container().Configure)
				if yamlFileErr != nil {
					Throw(yamlFileErr,0)
				}
			}
			return nil
		})
	if filepathErr != nil {
		Throw(filepathErr,0)
	}
	mapSslConf := Container().Config("SSL")
	SslOut = root + mapSslConf["OUT"].(string) + fmt.Sprintf("%s", PathSeparator)
	SslCrt = mapSslConf["CRT"].(string)
	SslKey = mapSslConf["KEY"].(string)
	return this
}

// 中间件装载器
func (this *Higo) Middleware(imiddleware ...IMiddleware) *Higo {
	for _, middleware := range imiddleware {
		this.middle = append(this.middle, &middleware)
	}
	return this
}

// http服务
func (this *Higo) HttpServe(conf string, router IRouterLoader) *Higo {
	this.serve = append(this.serve, &Hse{Config: conf, Router: router, Serve: "http"})
	return this
}

// https服务
func (this *Higo) HttpsServe(conf string, router IRouterLoader) *Higo {
	this.serve = append(this.serve, &Hse{Config: conf, Router: router, Serve: "https"})
	return this
}

// websocket服务
func (this *Higo) WebsocketServe(conf string, router IRouterLoader) *Higo {
	this.serve = append(this.serve, &Hse{Config: conf, Router: router, Serve: "websocket"})
	return this
}

// 是否自动生成ssl证书
func (this *Higo) IsAutoGenerateSsl(isAuto bool) *Higo {
	this.isAutoSsl = isAuto
	return this
}

// 启动
func (this *Higo) Boot() {
	// 配置
	this.config()
	// 中间件
	for _,m := range this.middle {
		mp := *m
		this.Engine.Use(mp.Loader(this))
	}
	// 是否使用自带ssl测试https
	if this.isAutoSsl {
		// 生成ssl证书
		utils.NewSsl(SslOut, SslCrt, SslKey).Generate()
	}
	// 服务
	for _, s := range this.serve {
		config := Container().Config(s.Config)
		addr, _ := config["Addr"]
		rt, _ := config["ReadTimeout"]
		wt, _ := config["WriteTimeout"]
		readTimeout, _ := rt.(int)
		writeTimeout, _ := wt.(int)
		serve := &http.Server{
			Addr:         addr.(string),
			Handler:      s.Router.Loader(this),
			ReadTimeout:  time.Duration(readTimeout) * time.Second,
			WriteTimeout: time.Duration(writeTimeout) * time.Second,
		}
		if s.Serve == "http" {
			this.eg.Go(func() error {
				fmt.Println("http 启动成功")
				return serve.ListenAndServe()
			})
		}
		if s.Serve == "https" {
			this.eg.Go(func() error {
				fmt.Println("https 启动成功")
				return serve.ListenAndServeTLS(SslOut + SslCrt, SslOut + SslKey)
			})
		}
		if s.Serve == "websocket" {
			this.eg.Go(func() error {
				fmt.Println("websocket 启动成功")
				return serve.ListenAndServe()
			})
		}
	}

	fmt.Println("启动成功")

	if err := this.eg.Wait(); err != nil {
		Logrus.Fatal(err)
	}
}

// 容器
func Container() *Containers {
	return hg.containers
}

// 获取路由
func (this *Higo) GetRoute(relativePath string) (Route, bool) {
	return Container().GetRoute(relativePath), true
}

// 路由组
func (this *Higo) AddGroup(prefix string, routes ...Route) *Higo {
	this.g = this.Engine.Group(prefix)
	for _, route := range routes {
		// 判断空标记
		IsEmptyFlag(route)
		// 添加路由容器
		Container().AddRoutes(route.RelativePath, route)
		method := strings.ToUpper(route.Method)
		this.GroupHandle(method, route.RelativePath, route.Handle)
	}
	return this
}

// 路由
func (this *Higo) AddRoute(routes ...Route) *Higo {
	for _, route := range routes {
		// 判断空标记
		IsEmptyFlag(route)
		// 添加路由容器
		Container().AddRoutes(route.RelativePath, route)
		method := strings.ToUpper(route.Method)
		this.Handle(method, route.RelativePath, route.Handle)
	}
	return this
}

// 路由组Handle
func (this *Higo) GroupHandle(httpMethod, relativePath string, handler interface{}) *Higo {
	if h := Convert(handler); h != nil {
		this.g.Handle(httpMethod, relativePath, h)
	}
	return this
}

// 路由Handle
func (this *Higo) Handle(httpMethod, relativePath string, handler interface{}) *Higo {
	fmt.Printf("%T\n",handler)
	if h := Convert(handler); h != nil {
		this.Engine.Handle(httpMethod, relativePath, h)
	}
	return this
}

func (this *Higo) Mount(group string, icontroller ...IController) *Higo {
	this.g = this.Engine.Group(group)
	for _, controller := range icontroller {
		this.currentGroup = group
		controller.Controller(this)
	}
	return this
}