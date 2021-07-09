package {{.Package}}

import (
	"github.com/dengpju/higo-gin/higo"
)

{{range .TplFields}}
func With{{.Field}}(v {{.Type}}) higo.Property {
	return func(class higo.IClass) {
		class.(*{{$.ModelImpl}}).{{.Field}} = v
	}
}
{{end}}
