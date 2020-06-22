package gengorums

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var nodeStreamProto = `
syntax = "proto3";

package {{ .GoPackageName }}.gengorums;

option go_package = {{ .GoImportPath }};

import "{{.Desc.Path}}";
{{- range fileImports .}}
import "{{.Path}}";
{{- end}}

service Gorums { rpc NodeStream(stream Message) returns (stream Message); }

message Message {
	oneof Data {
	{{- range $index, $type := fileMessages .}}
		{{$type.Message.Desc.FullName}} {{$type.Message.Desc.FullName.Name}} = {{inc $index}};
	{{- end}}
	}
	uint64 ID = 20000;
	int32 MethodID = 20001;
}
`

type messageType struct {
	Message  *protogen.Message
	Imported bool
}

func fileMessages(file *protogen.File) []messageType {
	msgs := make(map[*protogen.Message]struct{})
	for _, s := range file.Services {
		for _, m := range s.Methods {
			msgs[m.Input] = struct{}{}
			msgs[m.Output] = struct{}{}
		}
	}
	var unique []messageType
	for m := range msgs {
		imported := true
		// search all top-level messages to see if m is imported
		for i := range file.Messages {
			if file.Messages[i] == m {
				imported = false
			}
		}
		unique = append(unique, messageType{Message: m, Imported: imported})
	}
	return unique
}

func fileImports(file *protogen.File) []protoreflect.FileDescriptor {
	var imports []protoreflect.FileDescriptor
	for i := 0; i < file.Desc.Imports().Len(); i++ {
		imports = append(imports, file.Desc.Imports().Get(i).FileDescriptor)
	}
	return imports
}
