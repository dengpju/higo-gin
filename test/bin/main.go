package main

import (
	"fmt"
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Controllers/V3"
	"github.com/dengpju/higo-gin/test/app/Middlewares"
	"github.com/dengpju/higo-gin/test/providers"
	"github.com/dengpju/higo-gin/test/router"
	"github.com/dengpju/higo-ioc/injector"
)

func main()  {
	provider := providers.NewProvider()
	injector.BeanFactory.Config(provider)
	demoController := V3.NewDemoController()
	injector.BeanFactory.Apply(demoController)
	fmt.Println(demoController.DemoService)

	higo.Init().
		Middleware(Middlewares.NewAuth(), Middlewares.NewRunLog()).
		SetRoot(".\\test\\").
		//HttpServe("HTTP_HOST", router.NewHttp()).
		HttpsServe("HTTPS_HOST", router.NewHttps()).
		IsAutoGenerateSsl(true).
		Beans(higo.NewHgController(),V3.NewDemoController()).
		Boot()
}
