// Code generated by protoc-gen-gorums. DO NOT EDIT.
// versions:
// 	protoc-gen-gorums v0.5.0-devel
// 	protoc            v3.16.0
// source: ordering/order.proto

package ordering

import (
	context "context"
	fmt "fmt"
	gorums "github.com/relab/gorums"
	encoding "google.golang.org/grpc/encoding"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = gorums.EnforceVersion(5 - gorums.MinVersion)
	// Verify that the gorums runtime is sufficiently up-to-date.
	_ = gorums.EnforceVersion(gorums.MaxVersion - 5)
)

// A Configuration represents a static set of nodes on which quorum remote
// procedure calls may be invoked.
type Configuration struct {
	gorums.Configuration
	qspec QuorumSpec
}

// Nodes returns a slice of each available node. IDs are returned in the same
// order as they were provided in the creation of the Manager.
func (c *Configuration) Nodes() []*Node {
	nodes := make([]*Node, 0, c.Size())
	for _, n := range c.Configuration {
		nodes = append(nodes, &Node{n})
	}
	return nodes
}

// And returns a NodeListOption that can be used to create a new configuration combining c and d.
func (c Configuration) And(d *Configuration) gorums.NodeListOption {
	return c.Configuration.And(d.Configuration)
}

// Except returns a NodeListOption that can be used to create a new configuration
// from c without the nodes in rm.
func (c Configuration) Except(rm *Configuration) gorums.NodeListOption {
	return c.Configuration.Except(rm.Configuration)
}

func init() {
	if encoding.GetCodec(gorums.ContentSubtype) == nil {
		encoding.RegisterCodec(gorums.NewCodec())
	}
}

// Manager maintains a connection pool of nodes on
// which quorum calls can be performed.
type Manager struct {
	*gorums.Manager
}

// NewManager returns a new Manager for managing connection to nodes added
// to the manager. This function accepts manager options used to configure
// various aspects of the manager.
func NewManager(opts ...gorums.ManagerOption) (mgr *Manager) {
	mgr = &Manager{}
	mgr.Manager = gorums.NewManager(opts...)
	return mgr
}

// NewConfiguration returns a configuration based on the provided list of nodes (required)
// and an optional quorum specification. The QuorumSpec is necessary for call types that
// must process replies. For configurations only used for unicast or multicast call types,
// a QuorumSpec is not needed. The QuorumSpec interface is also a ConfigOption.
// Nodes can be supplied using WithNodeMap or WithNodeList, or WithNodeIDs.
// A new configuration can also be created from an existing configuration,
// using the And, WithNewNodes, Except, and WithoutNodes methods.
func (m *Manager) NewConfiguration(opts ...gorums.ConfigOption) (c *Configuration, err error) {
	if len(opts) < 1 || len(opts) > 2 {
		return nil, fmt.Errorf("wrong number of options: %d", len(opts))
	}
	c = &Configuration{}
	for _, opt := range opts {
		switch v := opt.(type) {
		case gorums.NodeListOption:
			c.Configuration, err = gorums.NewConfiguration(m.Manager, v)
			if err != nil {
				return nil, err
			}
		case QuorumSpec:
			// Must be last since v may match QuorumSpec if it is interface{}
			c.qspec = v
		default:
			return nil, fmt.Errorf("unknown option type: %v", v)
		}
	}
	// return an error if the QuorumSpec interface is not empty and no implementation was provided.
	var test interface{} = struct{}{}
	if _, empty := test.(QuorumSpec); !empty && c.qspec == nil {
		return nil, fmt.Errorf("missing required QuorumSpec")
	}
	return c, nil
}

// Nodes returns a slice of available nodes on this manager.
// IDs are returned in the order they were added at creation of the manager.
func (m *Manager) Nodes() []*Node {
	gorumsNodes := m.Manager.Nodes()
	nodes := make([]*Node, 0, len(gorumsNodes))
	for _, n := range gorumsNodes {
		nodes = append(nodes, &Node{n})
	}
	return nodes
}

type Node struct {
	*gorums.Node
}

// QCAsync asynchronously invokes a quorum call on configuration c
// and returns a AsyncResponse, which can be used to inspect the quorum call
// reply and error when available.
func (c *Configuration) QCAsync(ctx context.Context, in *Request) *AsyncResponse {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "ordering.GorumsTest.QCAsync",
	}
	cd.QuorumFunction = func(req protoreflect.ProtoMessage, replies map[uint32]protoreflect.ProtoMessage) (protoreflect.ProtoMessage, bool) {
		r := make(map[uint32]*Response, len(replies))
		for k, v := range replies {
			r[k] = v.(*Response)
		}
		return c.qspec.QCAsyncQF(req.(*Request), r)
	}

	fut := c.Configuration.AsyncCall(ctx, cd)
	return &AsyncResponse{fut}
}

// QuorumSpec is the interface of quorum functions for GorumsTest.
type QuorumSpec interface {
	gorums.ConfigOption

	// QCQF is the quorum function for the QC
	// quorum call method. The in parameter is the request object
	// supplied to the QC method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	QCQF(in *Request, replies map[uint32]*Response) (*Response, bool)

	// QCAsyncQF is the quorum function for the QCAsync
	// asynchronous quorum call method. The in parameter is the request object
	// supplied to the QCAsync method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	QCAsyncQF(in *Request, replies map[uint32]*Response) (*Response, bool)
}

// QC is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) QC(ctx context.Context, in *Request) (resp *Response, err error) {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "ordering.GorumsTest.QC",
	}
	cd.QuorumFunction = func(req protoreflect.ProtoMessage, replies map[uint32]protoreflect.ProtoMessage) (protoreflect.ProtoMessage, bool) {
		r := make(map[uint32]*Response, len(replies))
		for k, v := range replies {
			r[k] = v.(*Response)
		}
		return c.qspec.QCQF(req.(*Request), r)
	}

	res, err := c.Configuration.QuorumCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*Response), err
}

// UnaryRPC is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (n *Node) UnaryRPC(ctx context.Context, in *Request) (resp *Response, err error) {
	cd := gorums.CallData{
		Message: in,
		Method:  "ordering.GorumsTest.UnaryRPC",
	}

	res, err := n.Node.RPCCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*Response), err
}

// GorumsTest is the server-side API for the GorumsTest Service
type GorumsTest interface {
	QC(ctx gorums.ServerCtx, request *Request) (response *Response, err error)
	QCAsync(ctx gorums.ServerCtx, request *Request) (response *Response, err error)
	UnaryRPC(ctx gorums.ServerCtx, request *Request) (response *Response, err error)
}

func RegisterGorumsTestServer(srv *gorums.Server, impl GorumsTest) {
	srv.RegisterHandler("ordering.GorumsTest.QC", func(ctx gorums.ServerCtx, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*Request)
		defer ctx.Release()
		resp, err := impl.QC(ctx, req)
		gorums.SendMessage(ctx, finished, gorums.WrapMessage(in.Metadata, resp, err))
	})
	srv.RegisterHandler("ordering.GorumsTest.QCAsync", func(ctx gorums.ServerCtx, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*Request)
		defer ctx.Release()
		resp, err := impl.QCAsync(ctx, req)
		gorums.SendMessage(ctx, finished, gorums.WrapMessage(in.Metadata, resp, err))
	})
	srv.RegisterHandler("ordering.GorumsTest.UnaryRPC", func(ctx gorums.ServerCtx, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*Request)
		defer ctx.Release()
		resp, err := impl.UnaryRPC(ctx, req)
		gorums.SendMessage(ctx, finished, gorums.WrapMessage(in.Metadata, resp, err))
	})
}

type internalResponse struct {
	nid   uint32
	reply *Response
	err   error
}

// AsyncResponse is a async object for processing replies.
type AsyncResponse struct {
	*gorums.Async
}

// Get returns the reply and any error associated with the called method.
// The method blocks until a reply or error is available.
func (f *AsyncResponse) Get() (*Response, error) {
	resp, err := f.Async.Get()
	if err != nil {
		return nil, err
	}
	return resp.(*Response), err
}
