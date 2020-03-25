// Package gengorums is internal to the gorums protobuf module.
package gengorums

import (
	"fmt"
	"log"
	"os"

	"github.com/relab/gorums"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
)

// TODO(meling) replace github.com/relab/gorums with gorums.io as import package

// GenerateFile generates a _gorums.pb.go file containing Gorums service definitions.
func GenerateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Services) == 0 || !checkMethodOptions(file.Services, gorumsCallTypes...) {
		// there is nothing for this plugin to do
		return nil
	}
	if len(file.Services) > 1 {
		// To build multiple services, make separate proto files and
		// run the plugin separately for each proto file.
		// These cannot share the same Go package.
		log.Fatalln("Gorums does not support multiple services in the same proto file.")
	}
	// TODO(meling) make this more generic; figure out what are the reserved types from the static files.
	for _, msg := range file.Messages {
		msgName := fmt.Sprintf("%v", msg.Desc.Name())
		for _, reserved := range []string{"Configuration", "Node", "Manager", "ManagerOption"} {
			if msgName == reserved {
				log.Fatalf("%v.proto: contains message %s, which is a reserved Gorums type.\n", file.GeneratedFilenamePrefix, msgName)
			}
		}
	}

	filename := file.GeneratedFilenamePrefix + "_gorums.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by protoc-gen-gorums. DO NOT EDIT.")
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	g.P(staticCode)
	g.P()
	for path, ident := range pkgIdentMap {
		addImport(path, ident, g)
	}
	GenerateFileContent(gen, file, g)
	return g
}

// GenerateFileContent generates the Gorums service definitions, excluding the package statement.
func GenerateFileContent(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile) {
	data := servicesData{g, file.Services}
	g.P(mustExecute(parseTemplate("Manager", manager), data))
	g.P()
	g.P(mustExecute(parseTemplate("Node", node), data))
	g.P()
	g.P(mustExecute(parseTemplate("QuorumSpec", qspecInterface), data))
	g.P()
	g.P(mustExecute(parseTemplate("DataTypes", datatypes), data))
	g.P()
	genGorumsMethods(data, gorumsCallTypes...)
	g.P()
}

func genGorumsMethods(data servicesData, methodOptions ...*protoimpl.ExtensionInfo) {
	g := data.GenFile
	for _, service := range data.Services {
		for _, method := range service.Methods {
			if hasMethodOption(method, methodOptions...) {
				fmt.Fprintf(os.Stderr, "processing %s\n", method.GoName)
				g.P(genGorumsMethod(g, method))
			}
		}
	}
}

func genGorumsMethod(g *protogen.GeneratedFile, method *protogen.Method) string {
	methodOption := validateMethodExtensions(method)
	if template, ok := gorumsCallTypeTemplates[methodOption]; ok {
		return mustExecute(parseTemplate(methodOption.Name, template), methodData{g, method})
	}
	panic(fmt.Sprintf("unknown method type %s\n", method.GoName))
}

func callTypeName(method *protogen.Method) string {
	methodOption := validateMethodExtensions(method)
	if callTypeName, ok := gorumsCallTypeNames[methodOption]; ok {
		return callTypeName
	}
	panic(fmt.Sprintf("unknown method type %s\n", method.GoName))
}

type servicesData struct {
	GenFile  *protogen.GeneratedFile
	Services []*protogen.Service
}

type methodData struct {
	GenFile *protogen.GeneratedFile
	Method  *protogen.Method
}

// hasGorumsType returns true if one of the service methods specify
// the given gorums type.
func hasGorumsType(services []*protogen.Service, gorumsType string) bool {
	// TODO(meling) try to avoid this loop slice; reuse devTypes??
	for _, gType := range []string{"node", "qspec", "types", "manager"} {
		if gorumsType == gType {
			return true
		}
	}
	if methodOption, ok := gorumsTypes[gorumsType]; ok {
		return checkMethodOptions(services, methodOption)
	}
	return false
}

// compute index to start of option name
const index = len("gorums.")

// name to method option mapping
var gorumsTypes = map[string]*protoimpl.ExtensionInfo{
	gorums.E_Qc.Name[index:]:                gorums.E_Qc,
	gorums.E_QcFuture.Name[index:]:          gorums.E_QcFuture,
	gorums.E_Correctable.Name[index:]:       gorums.E_Correctable,
	gorums.E_CorrectableStream.Name[index:]: gorums.E_CorrectableStream,
	gorums.E_Multicast.Name[index:]:         gorums.E_Multicast,
	gorums.E_StrictOrdering.Name[index:]:    gorums.E_StrictOrdering,
}

var gorumsCallTypeTemplates = map[*protoimpl.ExtensionInfo]string{
	gorums.E_Qc:                quorumCall,
	gorums.E_QcFuture:          futureCall,
	gorums.E_Correctable:       correctableCall,
	gorums.E_CorrectableStream: correctableStreamCall,
	gorums.E_Multicast:         multicastCall,
	gorums.E_StrictOrdering:    strictOrderingQuorum,
}

var gorumsCallTypeNames = map[*protoimpl.ExtensionInfo]string{
	gorums.E_Qc:                "quorum",
	gorums.E_QcFuture:          "asynchronous quorum",
	gorums.E_Correctable:       "correctable quorum",
	gorums.E_CorrectableStream: "correctable stream quorum",
	gorums.E_Multicast:         "multicast",
	gorums.E_StrictOrdering:    "strict ordering quorum",
}

// gorumsCallTypes should list all available call types supported by Gorums.
// These are considered mutually incompatible.
var gorumsCallTypes = []*protoimpl.ExtensionInfo{
	gorums.E_Qc,
	gorums.E_QcFuture,
	gorums.E_Correctable,
	gorums.E_CorrectableStream,
	gorums.E_Multicast,
	gorums.E_StrictOrdering,
}

// callTypesWithInternal should list all available call types that
// has a quorum function and hence need an internal type that wraps
// the return type with additional information.
var callTypesWithInternal = []*protoimpl.ExtensionInfo{
	gorums.E_Qc,
	gorums.E_QcFuture,
	gorums.E_Correctable,
	gorums.E_CorrectableStream,
	gorums.E_StrictOrdering,
}

// callTypesWithPromiseObject lists all call types that returns
// a promise (future or correctable) object.
var callTypesWithPromiseObject = []*protoimpl.ExtensionInfo{
	gorums.E_QcFuture,
	gorums.E_Correctable,
	gorums.E_CorrectableStream,
}

// hasGorumsCallType returns true if the given method has specified
// one of the call types supported by Gorums.
func hasGorumsCallType(method *protogen.Method) bool {
	return hasMethodOption(method, gorumsCallTypes...)
}

// checkMethodOptions returns true if one of the methods provided by
// the given services has one of the given options.
func checkMethodOptions(services []*protogen.Service, methodOptions ...*protoimpl.ExtensionInfo) bool {
	for _, service := range services {
		for _, method := range service.Methods {
			if hasMethodOption(method, methodOptions...) {
				return true
			}
		}
	}
	return false
}

// hasMethodOption returns true if the method has one of the given method options.
func hasMethodOption(method *protogen.Method, methodOptions ...*protoimpl.ExtensionInfo) bool {
	ext := protoimpl.X.MessageOf(method.Desc.Options()).Interface()
	for _, callType := range methodOptions {
		if proto.HasExtension(ext, callType) {
			return true
		}
	}
	return false
}

// validateMethodExtensions returns the method option for the
// call type of the given method. If the method specifies multiple
// call types, validation will fail with a panic.
func validateMethodExtensions(method *protogen.Method) *protoimpl.ExtensionInfo {
	methExt := protoimpl.X.MessageOf(method.Desc.Options()).Interface()
	var firstOption *protoimpl.ExtensionInfo
	for _, callType := range gorumsCallTypes {
		if proto.HasExtension(methExt, callType) {
			if firstOption != nil {
				log.Fatalf("%s.%s: cannot combine options: '%s' and '%s'",
					method.Parent.Desc.Name(), method.Desc.Name(), firstOption.Name, callType.Name)
			}
			firstOption = callType
		}
	}

	isQuorumCallVariant := hasMethodOption(method, callTypesWithInternal...)
	switch {
	case !isQuorumCallVariant && proto.GetExtension(methExt, gorums.E_CustomReturnType) != "":
		// Only QC variants can define custom return type
		// (we don't support rewriting the plain gRPC methods.)
		log.Fatalf(
			"%s.%s: cannot combine non-quorum call method with the '%s' option",
			method.Parent.Desc.Name(), method.Desc.Name(), gorums.E_CustomReturnType.Name)

	case !isQuorumCallVariant && hasMethodOption(method, gorums.E_QfWithReq):
		// Only QC variants need to process replies.
		log.Fatalf(
			"%s.%s: cannot combine non-quorum call method with the '%s' option",
			method.Parent.Desc.Name(), method.Desc.Name(), gorums.E_QfWithReq.Name)

	case !hasMethodOption(method, gorums.E_Multicast, gorums.E_StrictOrdering) && method.Desc.IsStreamingClient():
		log.Fatalf(
			"%s.%s: client-server streams is only valid with the '%s' or '%s' options",
			method.Parent.Desc.Name(), method.Desc.Name(), gorums.E_Multicast.Name, gorums.E_StrictOrdering.Name)

	case hasMethodOption(method, gorums.E_Multicast) && !method.Desc.IsStreamingClient():
		log.Fatalf(
			"%s.%s: '%s' option is only valid for client-server streams methods",
			method.Parent.Desc.Name(), method.Desc.Name(), gorums.E_Multicast.Name)

	case !hasMethodOption(method, gorums.E_CorrectableStream, gorums.E_StrictOrdering) && method.Desc.IsStreamingServer():
		log.Fatalf(
			"%s.%s: server-client streams is only valid with the '%s' or '%s' options",
			method.Parent.Desc.Name(), method.Desc.Name(), gorums.E_CorrectableStream.Name, gorums.E_StrictOrdering.Name)

	case hasMethodOption(method, gorums.E_CorrectableStream) && !method.Desc.IsStreamingServer():
		log.Fatalf(
			"%s.%s: '%s' option is only valid for server-client streams",
			method.Parent.Desc.Name(), method.Desc.Name(), gorums.E_CorrectableStream.Name)

	case hasMethodOption(method, gorums.E_StrictOrdering) && (!method.Desc.IsStreamingClient() || !method.Desc.IsStreamingServer()):
		log.Fatalf(
			"%s.%s: '%s' option is only valid for bidirectional streams",
			method.Parent.Desc.Name(), method.Desc.Name(), gorums.E_StrictOrdering.Name)
	}

	return firstOption
}
