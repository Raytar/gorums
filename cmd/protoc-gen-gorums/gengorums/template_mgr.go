package gengorums

import (
	"github.com/relab/gorums"
	"google.golang.org/protobuf/compiler/protogen"
)

var managerDataVariables = `{{$file := .GenFile}}`

var managerData = `
type managerData struct {
{{range .Services}}
{{- range orderingMethods .Methods}}
{{$unexportMethod := unexport .GoName}}
{{$out := out $file .}}
{{$intOut := printf "internal%s" $out}}
	{{$unexportMethod}}ID uint64
	{{$unexportMethod}}Recv map[uint64]chan *{{$intOut}}
	{{$unexportMethod}}Lock {{use "sync.RWMutex" $file}}
{{- end -}}
{{end}}
}
`

var managerDataNew = `
func newManagerData() *managerData {
	return &managerData{
{{range .Services -}}
{{- range orderingMethods .Methods -}}
{{$unexportMethod := unexport .GoName -}}
{{$out := out $file . -}}
{{$intOut := printf "internal%s" $out -}}
		{{$unexportMethod}}Recv: make(map[uint64]chan *{{$intOut}}),
{{- end -}}
{{end}}
	}
}
`

var managerDataCreateNodeData = `
func (m *managerData) createNodeData() *nodeData {
	return &nodeData{
{{range .Services -}}
{{range orderingMethods .Methods -}}
{{$unexportMethod := unexport .GoName -}}
		{{$unexportMethod}}Send: make(chan *{{in $file .}}, 1),
		{{$unexportMethod}}Recv: m.{{$unexportMethod}}Recv,
		{{$unexportMethod}}MapLock: &m.{{$unexportMethod}}Lock,
{{end -}}
{{end}}
	}
}
`

var manager = managerDataVariables + managerData + managerDataNew + managerDataCreateNodeData

func orderingMethods(methods []*protogen.Method) []*protogen.Method {
	var ordering []*protogen.Method
	for _, m := range methods {
		if hasMethodOption(m, gorums.E_QcStrictOrdering) {
			ordering = append(ordering, m)
		}
	}
	return ordering
}
