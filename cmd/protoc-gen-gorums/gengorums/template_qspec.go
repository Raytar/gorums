package gengorums

import (
	"github.com/relab/gorums"
	"google.golang.org/protobuf/compiler/protogen"
)

var qspecInterface = `
{{$genFile := .GenFile}}
{{- range qspecServices .Services}}
// QuorumSpec is the interface of quorum functions for {{.GoName}}.
type QuorumSpec interface {
{{range qspecMethods .Methods -}}
	{{/* Below . is the method object */}}
	{{$method := .GoName}}
	{{$in := in $genFile .}}
	{{$out := out $genFile .}}
	{{$customOut := customOut $genFile .}}
	// {{$method}}QF is the quorum function for the {{$method}}
	// {{docName .}} call method. The in parameter is the request object
	// supplied to the {{$method}} method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should declare your quorum function with '_ *{{$in}}'.
	{{$method}}QF(in *{{$in}}, replies []*{{$out}}) (*{{$customOut}}{{withCorrectable . ", int"}}, bool)
{{end}}
}
{{end}}
`

// qspecMethods returns all Gorums methods that require
// a quorum function; that is, all except multicast and plain gRPC methods.
func qspecMethods(methods []*protogen.Method) (s []*protogen.Method) {
	for _, method := range methods {
		if hasMethodOption(method, gorums.E_Multicast) || !hasGorumsCallType(method) {
			// ignore multicast and non-Gorums methods
			continue
		}
		s = append(s, method)
	}
	return s
}

// qspecServices returns all services that have Gorums methods.
func qspecServices(services []*protogen.Service) (s []*protogen.Service) {
	for _, service := range services {
		if len(qspecMethods(service.Methods)) < 1 {
			// ignore services without Gorums methods
			continue
		}
		s = append(s, service)
	}
	return s
}
