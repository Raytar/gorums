package gengorums

import (
	"google.golang.org/protobuf/compiler/protogen"
)

var nodeServices = `
type nodeServices struct {
	{{range .Services}}
	{{.GoName}}Client
	{{- $serviceName := .GoName}}
	{{- range streamMethods .Methods}}
	{{unexport .GoName}}Client {{$serviceName}}_{{.GoName}}Client
	{{- end -}}
	{{- end}}
}
`

var nodeConnectStream = `
{{$errorf := use "fmt.Errorf" .GenFile}}
func (n *Node) connectStream(ctx {{use "context.Context" .GenFile}}) (err error) {
	{{- range .Services}}
	n.{{.GoName}}Client = New{{.GoName}}Client(n.conn)
	{{- end}}

	{{- range .Services -}}
	{{$serviceName := .GoName}}
	{{- range streamMethods .Methods}}
	n.{{unexport .GoName}}Client, err = n.{{$serviceName}}Client.{{.GoName}}(ctx)
	if err != nil {
		return {{$errorf}}("stream creation failed: %v", err)
	}
	{{- end -}}

	{{range orderingMethods .Methods -}}
	{{$unexportMethod := unexport .GoName}}
	go n.{{$unexportMethod}}SendMsgs()
	go n.{{$unexportMethod}}RecvMsgs(ctx)
	{{- end -}}
	{{end}}
	return nil
}
`

var nodeCloseStream = `
func (n *Node) closeStream() (err error) {
	{{- range .Services -}}
	{{- range streamMethods .Methods}}
	{{- if not .Desc.IsStreamingServer}}
	_, err = n.{{unexport .GoName}}Client.CloseAndRecv()
	{{- end -}}
	{{- end -}}
	{{range orderingMethods .Methods}}
	close(n.{{unexport .GoName}}Send)
	{{- end -}}
	{{end}}
	return err
}
`

var node = nodeServices + nodeConnectStream + nodeCloseStream

// streamMethods returns all methods that support client streaming.
func streamMethods(methods []*protogen.Method) []*protogen.Method {
	var s []*protogen.Method
	for _, method := range methods {
		if method.Desc.IsStreamingClient() {
			s = append(s, method)
		}
	}
	return s
}
