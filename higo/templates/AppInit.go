package templates

import (
	"bufio"
	"io"
	"os"
	"strings"
	"sync"
)

const (
	module = "module "
	funcDecl = "func (this *%s) %s() *%s%s"
)

var (
	moduleName  = ""
	moduleOnce  sync.Once
)

func init() {
	moduleOnce.Do(func() {
		pwd, _ := os.Getwd()
		gomodfile, err := os.Open(pwd + "/go.mod")
		if err != nil {
			panic(err)
		}
		defer gomodfile.Close()
		gomodbr := bufio.NewReader(gomodfile)
		for {
			a, _, c := gomodbr.ReadLine()
			if c == io.EOF {
				break
			}
			if strings.Contains(string(a), module) {
				moduleName = strings.TrimLeft(string(a), module)
				break
			}
		}
	})
}

func GetModName() string {
	return moduleName
}
