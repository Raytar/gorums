// Code generated by protoc-gen-gorums. DO NOT EDIT.

package benchmark

import (
	context "context"
	fmt "fmt"
	gorums "github.com/relab/gorums"
	encoding "google.golang.org/grpc/encoding"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	sort "sort"
	sync "sync"
)

// A Configuration represents a static set of nodes on which quorum remote
// procedure calls may be invoked.
type Configuration struct {
	id    uint32
	nodes []*gorums.Node
	mgr   *Manager
	qspec QuorumSpec
	errs  chan gorums.Error
}

// NewConfig returns a configuration for the given node addresses and quorum spec.
// The returned func() must be called to close the underlying connections.
// This is an experimental API.
func NewConfig(qspec QuorumSpec, opts ...gorums.ManagerOption) (*Configuration, func(), error) {
	man, err := NewManager(opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create manager: %v", err)
	}
	c, err := man.NewConfiguration(man.NodeIDs(), qspec)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create configuration: %v", err)
	}
	return c, func() { man.Close() }, nil
}

// ID reports the identifier for the configuration.
func (c *Configuration) ID() uint32 {
	return c.id
}

// NodeIDs returns a slice containing the local ids of all the nodes in the
// configuration. IDs are returned in the same order as they were provided in
// the creation of the Configuration.
func (c *Configuration) NodeIDs() []uint32 {
	ids := make([]uint32, len(c.nodes))
	for i, node := range c.nodes {
		ids[i] = node.ID()
	}
	return ids
}

// Nodes returns a slice of each available node. IDs are returned in the same
// order as they were provided in the creation of the Configuration.
func (c *Configuration) Nodes() []*Node {
	nodes := make([]*Node, 0, len(c.nodes))
	for _, n := range c.nodes {
		nodes = append(nodes, &Node{n})
	}
	return nodes
}

// Size returns the number of nodes in the configuration.
func (c *Configuration) Size() int {
	return len(c.nodes)
}

func (c *Configuration) String() string {
	return fmt.Sprintf("config-%d", c.id)
}

// Equal returns a boolean reporting whether a and b represents the same
// configuration.
func Equal(a, b *Configuration) bool { return a.id == b.id }

// SubError returns a channel for listening to individual node errors. Currently
// only a single listener is supported.
func (c *Configuration) SubError() <-chan gorums.Error {
	return c.errs
}

func init() {
	if encoding.GetCodec(gorums.ContentSubtype) == nil {
		encoding.RegisterCodec(gorums.NewCodec())
	}
}

func NewManager(opts ...gorums.ManagerOption) (mgr *Manager, err error) {
	mgr = &Manager{}
	mgr.Manager, err = gorums.NewManager(opts...)
	if err != nil {
		return nil, err
	}
	return mgr, nil
}

type Manager struct {
	*gorums.Manager
}

func (m *Manager) NewConfiguration(ids []uint32, qspec QuorumSpec) (*Configuration, error) {
	if len(ids) == 0 {
		return nil, gorums.IllegalConfigError("need at least one node")
	}

	var nodes []*gorums.Node
	unique := make(map[uint32]struct{})
	for _, nid := range ids {
		// ensure that identical IDs are only counted once
		if _, duplicate := unique[nid]; duplicate {
			continue
		}
		unique[nid] = struct{}{}

		node, found := m.Node(nid)
		if !found {
			return nil, gorums.NodeNotFoundError(nid)
		}

		i := sort.Search(len(nodes), func(i int) bool {
			return node.ID() < nodes[i].ID()
		})
		nodes = append(nodes, nil)
		copy(nodes[i+1:], nodes[i:])
		nodes[i] = node
	}

	c := &Configuration{
		nodes: nodes,
		mgr:   m,
		qspec: qspec,
	}
	return c, nil
}

// Nodes returns a slice of each available node. IDs are returned in the same
// order as they were provided in the creation of the Manager.
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
		Manager: c.mgr.Manager,
		Nodes:   c.nodes,
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

	fut := gorums.AsyncCall(ctx, cd)
	return &AsyncEcho{fut}
}

// Reference imports to suppress errors if they are not otherwise used.
var _ emptypb.Empty

// Multicast is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) Multicast(ctx context.Context, in *TimedMsg, opts ...gorums.CallOption) {

	cd := gorums.QuorumCallData{
		Manager: c.mgr.Manager,
		Nodes:   c.nodes,
		Message: in,
		Method:  "benchmark.Benchmark.Multicast",
	}

	gorums.Multicast(ctx, cd, opts...)
}

// QuorumSpec is the interface of quorum functions for Benchmark.
type QuorumSpec interface {

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
		Manager: c.mgr.Manager,
		Nodes:   c.nodes,
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

	res, err := gorums.QuorumCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*StartResponse), err
}

// StopServerBenchmark is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) StopServerBenchmark(ctx context.Context, in *StopRequest) (resp *Result, err error) {

	cd := gorums.QuorumCallData{
		Manager: c.mgr.Manager,
		Nodes:   c.nodes,
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

	res, err := gorums.QuorumCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*Result), err
}

// StartBenchmark is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) StartBenchmark(ctx context.Context, in *StartRequest) (resp *StartResponse, err error) {

	cd := gorums.QuorumCallData{
		Manager: c.mgr.Manager,
		Nodes:   c.nodes,
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

	res, err := gorums.QuorumCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*StartResponse), err
}

// StopBenchmark is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) StopBenchmark(ctx context.Context, in *StopRequest) (resp *MemoryStatList, err error) {

	cd := gorums.QuorumCallData{
		Manager: c.mgr.Manager,
		Nodes:   c.nodes,
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

	res, err := gorums.QuorumCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*MemoryStatList), err
}

// benchmarks
func (c *Configuration) QuorumCall(ctx context.Context, in *Echo) (resp *Echo, err error) {

	cd := gorums.QuorumCallData{
		Manager: c.mgr.Manager,
		Nodes:   c.nodes,
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

	res, err := gorums.QuorumCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*Echo), err
}

// SlowServer is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) SlowServer(ctx context.Context, in *Echo) (resp *Echo, err error) {

	cd := gorums.QuorumCallData{
		Manager: c.mgr.Manager,
		Nodes:   c.nodes,
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

	res, err := gorums.QuorumCall(ctx, cd)
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
