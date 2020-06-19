package gengorums

var nodeStreamProto = `
syntax = "proto3";

package gorums.{{ .GoPackageName }};

option go_package = {{ .GoImportPath }};

service Gorums { rpc NodeStream(stream Message) returns (stream Message); }

message Message {
	oneof Data {
		{{- range $index, $type := allTypes .Services}}
		{{$type.Desc.FullName}} {{$type.GoIdent.GoName}} = {{$index}};
		{{- end}}
	}
	uint64 ID = 20000;
	int32 MethodID = 20001;
}
`
