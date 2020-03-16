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
{{$context := use "context.Background" .GenFile}}
{{$errorf := use "fmt.Errorf" .GenFile}}
func (n *Node) connectStream() (err error) {
	{{- range .Services}}
	n.{{.GoName}}Client = New{{.GoName}}Client(n.conn)
	{{- end}}

	{{- range .Services -}}
	{{$serviceName := .GoName}}
	{{- range streamMethods .Methods}}
	n.{{unexport .GoName}}Client, err = n.{{$serviceName}}Client.{{.GoName}}({{$context}}())
	if err != nil {
		return {{$errorf}}("stream creation failed: %v", err)
	}
	{{- end -}}

	{{range orderingMethods .Methods -}}
	{{$unexportMethod := unexport .GoName}}
	go n.{{$unexportMethod}}SendMsgs()
	go n.{{$unexportMethod}}RecvMsgs()
	{{- end -}}
	{{end}}
	return nil
}
`

var nodeCloseStream = `
func (n *Node) closeStream() (err error) {
	{{- range .Services -}}
	{{- range streamMethods .Methods}}
	{{- if .Desc.IsStreamingServer}}
	err = n.{{unexport .GoName}}Client.CloseSend()
	{{- else}}
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

var nodeHandlers = `
{{$file := .GenFile -}}
{{$ioEOF := use "io.EOF" .GenFile -}}
{{range .Services -}}
{{range orderingMethods .Methods}}
{{$unexportMethod := unexport .GoName -}}
{{$method := .GoName -}}
{{$out := out $file . -}}
{{$intOut := printf "internal%s" $out -}}
func (n *Node) {{$unexportMethod}}SendMsgs() {
	for msg := range n.{{$unexportMethod}}Send {
		err := n.{{$unexportMethod}}Client.SendMsg(msg)
		if err != nil {
			if err != {{$ioEOF}} {
				n.setLastErr(err)
			}
			return
		}
	}
}

func (n *Node) {{$unexportMethod}}RecvMsgs() {
	for {
		msg := new({{$out}})
		err := n.{{$unexportMethod}}Client.RecvMsg(msg)
		if err != nil {
			if err != {{$ioEOF}} {
				n.setLastErr(err)
			}
			return
		}
		id := msg.{{msgIDField .}}
		n.{{$unexportMethod}}Lock.RLock()
		if c, ok := n.{{$unexportMethod}}Recv[id]; ok {
			c <- &{{$intOut}}{n.id, msg, nil}
		}
		n.{{$unexportMethod}}Lock.RUnlock()
	}
}
{{- end -}}
{{end -}}
`

var nodeData = `
{{$file := .GenFile}}
type nodeData struct {
{{range .Services}}
{{- range orderingMethods .Methods}}
{{$unexportMethod := unexport .GoName}}
{{$out := out $file .}}
{{$intOut := printf "internal%s" $out}}
	{{$unexportMethod}}Send chan *{{in $file .}}
	{{$unexportMethod}}Recv map[uint64]chan *{{$intOut}}
	{{$unexportMethod}}Lock *{{use "sync.RWMutex" $file}}
{{- end -}}
{{end}}
}
`

var node = nodeServices + nodeConnectStream + nodeCloseStream + nodeHandlers + nodeData

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
