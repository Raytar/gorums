// Code generated by protoc-gen-gorums. DO NOT EDIT.
// versions:
// 	protoc-gen-gorums v0.3.0-devel
// 	protoc            v3.15.5
// source: benchmark/benchmark.proto

package benchmark

import (
	context "context"
	fmt "fmt"
	gorums "github.com/relab/gorums"
	encoding "google.golang.org/grpc/encoding"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = gorums.EnforceVersion(gorums.GenVersion - gorums.MinVersion)
	// Verify that gorums runtime is sufficiently up-to-date.
	_ = gorums.EnforceVersion(gorums.MaxVersion - gorums.GenVersion)
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
	gorumsNodes := c.Configuration.Nodes()
	nodes := make([]*Node, 0, len(gorumsNodes))
	for _, n := range gorumsNodes {
		nodes = append(nodes, &Node{n})
	}
	return nodes
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
// and an optional quorum specification. The QuorumSpec is require for call types that
// must process replies. For configurations only used for unicast or multicast call types,
// a QuorumSpec is not needed. The QuorumSpec interface is also a ConfigOption.
// Nodes can be supplied using WithNodeMap or WithNodeList or WithNodeIDs.
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

// AsyncQuorumCall asynchronously invokes a quorum call on configuration c
// and returns a AsyncEcho, which can be used to inspect the quorum call
// reply and error when available.
func (c *Configuration) AsyncQuorumCall(ctx context.Context, in *Echo) *AsyncEcho {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "benchmark.Benchmark.AsyncQuorumCall",
	}
	cd.QuorumFunction = func(req protoreflect.ProtoMessage, replies map[uint32]protoreflect.ProtoMessage) (protoreflect.ProtoMessage, bool) {
		r := make(map[uint32]*Echo, len(replies))
		for k, v := range replies {
			r[k] = v.(*Echo)
		}
		return c.qspec.AsyncQuorumCallQF(req.(*Echo), r)
	}

	fut := c.Configuration.AsyncCall(ctx, cd)
	return &AsyncEcho{fut}
}

// Reference imports to suppress errors if they are not otherwise used.
var _ emptypb.Empty

// Multicast is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) Multicast(ctx context.Context, in *TimedMsg, opts ...gorums.CallOption) {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "benchmark.Benchmark.Multicast",
	}

	c.Configuration.Multicast(ctx, cd, opts...)
}

// QuorumSpec is the interface of quorum functions for Benchmark.
type QuorumSpec interface {
	gorums.ConfigOption

	// StartServerBenchmarkQF is the quorum function for the StartServerBenchmark
	// quorum call method. The in parameter is the request object
	// supplied to the StartServerBenchmark method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *StartRequest'.
	StartServerBenchmarkQF(in *StartRequest, replies map[uint32]*StartResponse) (*StartResponse, bool)

	// StopServerBenchmarkQF is the quorum function for the StopServerBenchmark
	// quorum call method. The in parameter is the request object
	// supplied to the StopServerBenchmark method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *StopRequest'.
	StopServerBenchmarkQF(in *StopRequest, replies map[uint32]*Result) (*Result, bool)

	// StartBenchmarkQF is the quorum function for the StartBenchmark
	// quorum call method. The in parameter is the request object
	// supplied to the StartBenchmark method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *StartRequest'.
	StartBenchmarkQF(in *StartRequest, replies map[uint32]*StartResponse) (*StartResponse, bool)

	// StopBenchmarkQF is the quorum function for the StopBenchmark
	// quorum call method. The in parameter is the request object
	// supplied to the StopBenchmark method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *StopRequest'.
	StopBenchmarkQF(in *StopRequest, replies map[uint32]*MemoryStat) (*MemoryStatList, bool)

	// QuorumCallQF is the quorum function for the QuorumCall
	// quorum call method. The in parameter is the request object
	// supplied to the QuorumCall method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Echo'.
	QuorumCallQF(in *Echo, replies map[uint32]*Echo) (*Echo, bool)

	// AsyncQuorumCallQF is the quorum function for the AsyncQuorumCall
	// asynchronous quorum call method. The in parameter is the request object
	// supplied to the AsyncQuorumCall method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Echo'.
	AsyncQuorumCallQF(in *Echo, replies map[uint32]*Echo) (*Echo, bool)

	// SlowServerQF is the quorum function for the SlowServer
	// quorum call method. The in parameter is the request object
	// supplied to the SlowServer method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Echo'.
	SlowServerQF(in *Echo, replies map[uint32]*Echo) (*Echo, bool)
}

// StartServerBenchmark is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) StartServerBenchmark(ctx context.Context, in *StartRequest) (resp *StartResponse, err error) {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "benchmark.Benchmark.StartServerBenchmark",
	}
	cd.QuorumFunction = func(req protoreflect.ProtoMessage, replies map[uint32]protoreflect.ProtoMessage) (protoreflect.ProtoMessage, bool) {
		r := make(map[uint32]*StartResponse, len(replies))
		for k, v := range replies {
			r[k] = v.(*StartResponse)
		}
		return c.qspec.StartServerBenchmarkQF(req.(*StartRequest), r)
	}

	res, err := c.Configuration.QuorumCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*StartResponse), err
}

// StopServerBenchmark is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) StopServerBenchmark(ctx context.Context, in *StopRequest) (resp *Result, err error) {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "benchmark.Benchmark.StopServerBenchmark",
	}
	cd.QuorumFunction = func(req protoreflect.ProtoMessage, replies map[uint32]protoreflect.ProtoMessage) (protoreflect.ProtoMessage, bool) {
		r := make(map[uint32]*Result, len(replies))
		for k, v := range replies {
			r[k] = v.(*Result)
		}
		return c.qspec.StopServerBenchmarkQF(req.(*StopRequest), r)
	}

	res, err := c.Configuration.QuorumCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*Result), err
}

// StartBenchmark is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) StartBenchmark(ctx context.Context, in *StartRequest) (resp *StartResponse, err error) {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "benchmark.Benchmark.StartBenchmark",
	}
	cd.QuorumFunction = func(req protoreflect.ProtoMessage, replies map[uint32]protoreflect.ProtoMessage) (protoreflect.ProtoMessage, bool) {
		r := make(map[uint32]*StartResponse, len(replies))
		for k, v := range replies {
			r[k] = v.(*StartResponse)
		}
		return c.qspec.StartBenchmarkQF(req.(*StartRequest), r)
	}

	res, err := c.Configuration.QuorumCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*StartResponse), err
}

// StopBenchmark is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) StopBenchmark(ctx context.Context, in *StopRequest) (resp *MemoryStatList, err error) {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "benchmark.Benchmark.StopBenchmark",
	}
	cd.QuorumFunction = func(req protoreflect.ProtoMessage, replies map[uint32]protoreflect.ProtoMessage) (protoreflect.ProtoMessage, bool) {
		r := make(map[uint32]*MemoryStat, len(replies))
		for k, v := range replies {
			r[k] = v.(*MemoryStat)
		}
		return c.qspec.StopBenchmarkQF(req.(*StopRequest), r)
	}

	res, err := c.Configuration.QuorumCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*MemoryStatList), err
}

// benchmarks
func (c *Configuration) QuorumCall(ctx context.Context, in *Echo) (resp *Echo, err error) {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "benchmark.Benchmark.QuorumCall",
	}
	cd.QuorumFunction = func(req protoreflect.ProtoMessage, replies map[uint32]protoreflect.ProtoMessage) (protoreflect.ProtoMessage, bool) {
		r := make(map[uint32]*Echo, len(replies))
		for k, v := range replies {
			r[k] = v.(*Echo)
		}
		return c.qspec.QuorumCallQF(req.(*Echo), r)
	}

	res, err := c.Configuration.QuorumCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*Echo), err
}

// SlowServer is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) SlowServer(ctx context.Context, in *Echo) (resp *Echo, err error) {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "benchmark.Benchmark.SlowServer",
	}
	cd.QuorumFunction = func(req protoreflect.ProtoMessage, replies map[uint32]protoreflect.ProtoMessage) (protoreflect.ProtoMessage, bool) {
		r := make(map[uint32]*Echo, len(replies))
		for k, v := range replies {
			r[k] = v.(*Echo)
		}
		return c.qspec.SlowServerQF(req.(*Echo), r)
	}

	res, err := c.Configuration.QuorumCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*Echo), err
}

// Benchmark is the server-side API for the Benchmark Service
type Benchmark interface {
	StartServerBenchmark(context.Context, *StartRequest, func(*StartResponse, error))
	StopServerBenchmark(context.Context, *StopRequest, func(*Result, error))
	StartBenchmark(context.Context, *StartRequest, func(*StartResponse, error))
	StopBenchmark(context.Context, *StopRequest, func(*MemoryStat, error))
	QuorumCall(context.Context, *Echo, func(*Echo, error))
	AsyncQuorumCall(context.Context, *Echo, func(*Echo, error))
	SlowServer(context.Context, *Echo, func(*Echo, error))
	Multicast(context.Context, *TimedMsg)
}

func RegisterBenchmarkServer(srv *gorums.Server, impl Benchmark) {
	srv.RegisterHandler("benchmark.Benchmark.StartServerBenchmark", func(ctx context.Context, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*StartRequest)
		once := new(sync.Once)
		f := func(resp *StartResponse, err error) {
			once.Do(func() {
				finished <- gorums.WrapMessage(in.Metadata, resp, err)
			})
		}
		impl.StartServerBenchmark(ctx, req, f)
	})
	srv.RegisterHandler("benchmark.Benchmark.StopServerBenchmark", func(ctx context.Context, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*StopRequest)
		once := new(sync.Once)
		f := func(resp *Result, err error) {
			once.Do(func() {
				finished <- gorums.WrapMessage(in.Metadata, resp, err)
			})
		}
		impl.StopServerBenchmark(ctx, req, f)
	})
	srv.RegisterHandler("benchmark.Benchmark.StartBenchmark", func(ctx context.Context, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*StartRequest)
		once := new(sync.Once)
		f := func(resp *StartResponse, err error) {
			once.Do(func() {
				finished <- gorums.WrapMessage(in.Metadata, resp, err)
			})
		}
		impl.StartBenchmark(ctx, req, f)
	})
	srv.RegisterHandler("benchmark.Benchmark.StopBenchmark", func(ctx context.Context, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*StopRequest)
		once := new(sync.Once)
		f := func(resp *MemoryStat, err error) {
			once.Do(func() {
				finished <- gorums.WrapMessage(in.Metadata, resp, err)
			})
		}
		impl.StopBenchmark(ctx, req, f)
	})
	srv.RegisterHandler("benchmark.Benchmark.QuorumCall", func(ctx context.Context, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*Echo)
		once := new(sync.Once)
		f := func(resp *Echo, err error) {
			once.Do(func() {
				finished <- gorums.WrapMessage(in.Metadata, resp, err)
			})
		}
		impl.QuorumCall(ctx, req, f)
	})
	srv.RegisterHandler("benchmark.Benchmark.AsyncQuorumCall", func(ctx context.Context, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*Echo)
		once := new(sync.Once)
		f := func(resp *Echo, err error) {
			once.Do(func() {
				finished <- gorums.WrapMessage(in.Metadata, resp, err)
			})
		}
		impl.AsyncQuorumCall(ctx, req, f)
	})
	srv.RegisterHandler("benchmark.Benchmark.SlowServer", func(ctx context.Context, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*Echo)
		once := new(sync.Once)
		f := func(resp *Echo, err error) {
			once.Do(func() {
				finished <- gorums.WrapMessage(in.Metadata, resp, err)
			})
		}
		impl.SlowServer(ctx, req, f)
	})
	srv.RegisterHandler("benchmark.Benchmark.Multicast", func(ctx context.Context, in *gorums.Message, _ chan<- *gorums.Message) {
		req := in.Message.(*TimedMsg)
		impl.Multicast(ctx, req)
	})
}

type internalEcho struct {
	nid   uint32
	reply *Echo
	err   error
}

type internalMemoryStat struct {
	nid   uint32
	reply *MemoryStat
	err   error
}

type internalResult struct {
	nid   uint32
	reply *Result
	err   error
}

type internalStartResponse struct {
	nid   uint32
	reply *StartResponse
	err   error
}

// AsyncEcho is a async object for processing replies.
type AsyncEcho struct {
	*gorums.Async
}

// Get returns the reply and any error associated with the called method.
// The method blocks until a reply or error is available.
func (f *AsyncEcho) Get() (*Echo, error) {
	resp, err := f.Async.Get()
	if err != nil {
		return nil, err
	}
	return resp.(*Echo), err
}
