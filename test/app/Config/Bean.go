package Config

import (
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/test/app/Services"
	"github.com/gomodule/redigo/redis"
)

type Bean struct {
}

func NewBean() *Bean {
	return &Bean{}
}

func (this *Bean) Provider() {

}

func (this *Bean) DemoService() *Services.DemoService {
	return Services.NewDemoService()
}

func (this *Bean) NewGorm() *higo.Gorm {
	return higo.NewGorm()
}

func (this *Bean) NewRedisPool() *redis.Pool {
	return higo.RedisPool
}

func (this *Bean) NewRedisAdapter() *higo.RedisAdapter {
	return higo.NewRedisAdapter()
}
