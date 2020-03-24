package gengorums

import (
	"google.golang.org/protobuf/compiler/protogen"
)

var nodeVariables = `{{$context := use "context.Context" .GenFile}}`

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
func (n *Node) connectStream(ctx {{$context}}) (err error) {
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

var nodeHandlers = `
{{$file := .GenFile -}}
{{range .Services -}}
{{$serviceName := .GoName}}
{{range orderingMethods .Methods}}
{{$unexportMethod := unexport .GoName -}}
{{$method := .GoName -}}
{{$out := out $file . -}}
{{$intOut := printf "internal%s" $out -}}
func (n *Node) {{$unexportMethod}}SendMsgs() {
	for msg := range n.{{$unexportMethod}}Send {
		if n.{{$unexportMethod}}LinkBroken {
			id := msg.{{msgIDField .}}
			err := {{use "status.Errorf" $file}}({{use "codes.Unavailable" $file}}, "stream is down")
			n.{{$unexportMethod}}MapLock.RLock()
			if c, ok := n.{{$unexportMethod}}Recv[id]; ok {
				c <- &{{$intOut}}{n.id, nil, err}
			}
			n.{{$unexportMethod}}MapLock.RUnlock()
		}
		n.{{$unexportMethod}}StreamLock.RLock()
		 
		err := n.{{$unexportMethod}}Client.SendMsg(msg)
		if err == nil {
			n.{{$unexportMethod}}StreamLock.RUnlock()
			continue
		}
		n.{{$unexportMethod}}LinkBroken = true
		n.{{$unexportMethod}}StreamLock.RUnlock()
		// return the error
		id := msg.{{msgIDField .}}
		n.{{$unexportMethod}}MapLock.RLock()
		if c, ok := n.{{$unexportMethod}}Recv[id]; ok {
			c <- &{{$intOut}}{n.id, nil, err}
		}
		n.{{$unexportMethod}}MapLock.RUnlock()
	}
}

func (n *Node) {{$unexportMethod}}RecvMsgs(ctx {{$context}}) {
	for {
		msg := new({{$out}})
		n.{{$unexportMethod}}StreamLock.RLock()
		err := n.{{$unexportMethod}}Client.RecvMsg(msg)
		if err != nil {
			n.{{$unexportMethod}}LinkBroken = true
			n.{{$unexportMethod}}StreamLock.RUnlock()
			n.setLastErr(err)
			// reconnect
			n.{{$unexportMethod}}Reconnect(ctx)
		} else {
			n.{{$unexportMethod}}StreamLock.RUnlock()
			id := msg.{{msgIDField .}}
			n.{{$unexportMethod}}MapLock.RLock()
			if c, ok := n.{{$unexportMethod}}Recv[id]; ok {
				c <- &{{$intOut}}{n.id, msg, nil}
			}
			n.{{$unexportMethod}}MapLock.RUnlock()
		}
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (n *Node) {{$unexportMethod}}Reconnect(ctx {{$context}}) {
	// attempt to reconnect with exponential backoff
	// TODO: Allow using a custom config
	bc := {{use "backoff.DefaultConfig" $file}}
	retries := 0.0
	r := {{use "rand.New" $file}}({{use "rand.NewSource" $file}}({{use "time.Now" $file}}().UnixNano()))
	for {
		var err error
		n.{{$unexportMethod}}Client, err = n.{{$serviceName}}Client.{{.GoName}}(ctx)
		if err == nil {
			n.{{$unexportMethod}}LinkBroken = false
			return
		}
		delay := float64(bc.BaseDelay)
		max := float64(bc.MaxDelay)
		if retries > 0 {
			delay = {{use "math.Pow" $file}}(delay, retries)
			delay = {{use "math.Min" $file}}(delay, max)
		}
		delay *= 1 + bc.Jitter*(r.Float64()*2-1)
		{{use "time.Sleep" $file}}({{use "time.Duration" $file}}(delay))
		retries++
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}
{{- end -}}
{{end -}}
`

var nodeData = `
{{$file := .GenFile}}
{{$mutex := use "sync.RWMutex" .GenFile}}
type nodeData struct {
{{range .Services}}
{{- range orderingMethods .Methods}}
{{$unexportMethod := unexport .GoName}}
{{$out := out $file .}}
{{$intOut := printf "internal%s" $out}}
	{{$unexportMethod}}Send chan *{{in $file .}}
	{{$unexportMethod}}Recv map[uint64]chan *{{$intOut}}
	{{$unexportMethod}}MapLock *{{$mutex}}
	{{$unexportMethod}}StreamLock {{$mutex}}
	{{$unexportMethod}}LinkBroken bool
{{- end -}}
{{end}}
}
`

var node = nodeVariables + nodeServices + nodeConnectStream + nodeCloseStream + nodeHandlers + nodeData

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
