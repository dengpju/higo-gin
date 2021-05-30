package UserModel

import (
	"gitee.com/dengpju/higo-code/code"
	"github.com/dengpju/higo-gin/higo"
)

type UserModelImpl struct {
	Id    int    `gorm:"column:id" json:"id" binding:"required"`
	Utel  string `gorm:"column:u_tel" json:"utel" binding:"Utel"`
	Uname string `gorm:"column:uname" json:"uname" binding:"UserName"`
}

func init() {
	//初始化校验器
	u := &UserModelImpl{}
	u.InitValidator()
}

func New(attrs ...higo.Property) *UserModelImpl {
	u := &UserModelImpl{}
	higo.Propertys(attrs).Apply(u)
	return u
}

func (this *UserModelImpl) New() higo.IClass {
	return New()
}

func (this *UserModelImpl) Mutate(attrs ...higo.Property) higo.Model {
	higo.Propertys(attrs).Apply(this)
	return this
}

func (this *UserModelImpl) InitValidator() higo.Valid {
	return higo.RegisterValid(this).
		Tag("UserName",
			higo.Rule("required", code.Message("20000@UserName必须填")),
			higo.Rule("min=5", code.Message("20000@UserName必须填大于5"))).
		Tag("Utel",
			higo.Rule("required", code.Message("20000@Utel必须填")),
			higo.Rule("min=4", code.Message("20000@Utel大于4")))
}
