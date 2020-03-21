// Code generated by protoc-gen-gorums. DO NOT EDIT.

package dev

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	trace "golang.org/x/net/trace"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	time "time"
)

// QuorumCallFuture plain.
func (c *Configuration) QuorumCallFuture(ctx context.Context, in *Request, opts ...grpc.CallOption) *FutureResponse {
	fut := &FutureResponse{
		NodeIDs: make([]uint32, 0, c.n),
		c:       make(chan struct{}, 1),
	}
	go func() {
		defer close(fut.c)
		c.quorumCallFuture(ctx, in, fut, opts...)
	}()
	return fut
}

// Get returns the reply and any error associated with the QuorumCallFuture.
// The method blocks until a reply or error is available.
func (f *FutureResponse) Get() (*Response, error) {
	<-f.c
	return f.Response, f.err
}

// Done reports if a reply and/or error is available for the QuorumCallFuture.
func (f *FutureResponse) Done() bool {
	select {
	case <-f.c:
		return true
	default:
		return false
	}
}

func (c *Configuration) quorumCallFuture(ctx context.Context, in *Request, resp *FutureResponse, opts ...grpc.CallOption) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "QuorumCallFuture")
		defer ti.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = time.Until(deadline)
		}
		ti.LazyLog(&ti.firstLine, false)
		ti.LazyLog(&payload{sent: true, msg: in}, false)

		defer func() {
			ti.LazyLog(&qcresult{ids: resp.NodeIDs, reply: resp.Response, err: resp.err}, false)
			if resp.err != nil {
				ti.SetError()
			}
		}()
	}

	expected := c.n
	replyChan := make(chan internalResponse, expected)
	for _, n := range c.nodes {
		go n.QuorumCallFuture(ctx, in, replyChan)
	}

	var (
		replyValues = make([]*Response, 0, c.n)
		reply       *Response
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			resp.NodeIDs = append(resp.NodeIDs, r.nid)
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			if reply, quorum = c.qspec.QuorumCallFutureQF(replyValues); quorum {
				resp.Response, resp.err = reply, nil
				return
			}
		case <-ctx.Done():
			resp.Response, resp.err = reply, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
			return
		}
		if len(errs)+len(replyValues) == expected {
			resp.Response, resp.err = reply, QuorumCallError{"incomplete call", len(replyValues), errs}
			return
		}
	}
}

func (n *Node) QuorumCallFuture(ctx context.Context, in *Request, replyChan chan<- internalResponse) {
	reply := new(Response)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ZorumsService/QuorumCallFuture", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalResponse{n.id, reply, err}
}

// QuorumCallFuturePerNodeArg with per_node_arg option.
func (c *Configuration) QuorumCallFuturePerNodeArg(ctx context.Context, in *Request, f func(*Request, uint32) *Request, opts ...grpc.CallOption) *FutureResponse {
	fut := &FutureResponse{
		NodeIDs: make([]uint32, 0, c.n),
		c:       make(chan struct{}, 1),
	}
	go func() {
		defer close(fut.c)
		c.quorumCallFuturePerNodeArg(ctx, in, f, fut, opts...)
	}()
	return fut
}

func (c *Configuration) quorumCallFuturePerNodeArg(ctx context.Context, in *Request, f func(*Request, uint32) *Request, resp *FutureResponse, opts ...grpc.CallOption) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "QuorumCallFuturePerNodeArg")
		defer ti.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = time.Until(deadline)
		}
		ti.LazyLog(&ti.firstLine, false)
		ti.LazyLog(&payload{sent: true, msg: in}, false)

		defer func() {
			ti.LazyLog(&qcresult{ids: resp.NodeIDs, reply: resp.Response, err: resp.err}, false)
			if resp.err != nil {
				ti.SetError()
			}
		}()
	}

	expected := c.n
	replyChan := make(chan internalResponse, expected)
	for _, n := range c.nodes {
		nodeArg := f(in, n.id)
		if nodeArg == nil {
			expected--
			continue
		}
		go n.QuorumCallFuturePerNodeArg(ctx, nodeArg, replyChan)
	}

	var (
		replyValues = make([]*Response, 0, c.n)
		reply       *Response
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			resp.NodeIDs = append(resp.NodeIDs, r.nid)
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			if reply, quorum = c.qspec.QuorumCallFuturePerNodeArgQF(replyValues); quorum {
				resp.Response, resp.err = reply, nil
				return
			}
		case <-ctx.Done():
			resp.Response, resp.err = reply, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
			return
		}
		if len(errs)+len(replyValues) == expected {
			resp.Response, resp.err = reply, QuorumCallError{"incomplete call", len(replyValues), errs}
			return
		}
	}
}

func (n *Node) QuorumCallFuturePerNodeArg(ctx context.Context, in *Request, replyChan chan<- internalResponse) {
	reply := new(Response)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ZorumsService/QuorumCallFuturePerNodeArg", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalResponse{n.id, reply, err}
}

// QuorumCallfutureQFWithRequestArg with qf_with_req option.
func (c *Configuration) QuorumCallfutureQFWithRequestArg(ctx context.Context, in *Request, opts ...grpc.CallOption) *FutureResponse {
	fut := &FutureResponse{
		NodeIDs: make([]uint32, 0, c.n),
		c:       make(chan struct{}, 1),
	}
	go func() {
		defer close(fut.c)
		c.quorumCallfutureQFWithRequestArg(ctx, in, fut, opts...)
	}()
	return fut
}

func (c *Configuration) quorumCallfutureQFWithRequestArg(ctx context.Context, in *Request, resp *FutureResponse, opts ...grpc.CallOption) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "QuorumCallfutureQFWithRequestArg")
		defer ti.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = time.Until(deadline)
		}
		ti.LazyLog(&ti.firstLine, false)
		ti.LazyLog(&payload{sent: true, msg: in}, false)

		defer func() {
			ti.LazyLog(&qcresult{ids: resp.NodeIDs, reply: resp.Response, err: resp.err}, false)
			if resp.err != nil {
				ti.SetError()
			}
		}()
	}

	expected := c.n
	replyChan := make(chan internalResponse, expected)
	for _, n := range c.nodes {
		go n.QuorumCallfutureQFWithRequestArg(ctx, in, replyChan)
	}

	var (
		replyValues = make([]*Response, 0, c.n)
		reply       *Response
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			resp.NodeIDs = append(resp.NodeIDs, r.nid)
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			if reply, quorum = c.qspec.QuorumCallfutureQFWithRequestArgQF(in, replyValues); quorum {
				resp.Response, resp.err = reply, nil
				return
			}
		case <-ctx.Done():
			resp.Response, resp.err = reply, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
			return
		}
		if len(errs)+len(replyValues) == expected {
			resp.Response, resp.err = reply, QuorumCallError{"incomplete call", len(replyValues), errs}
			return
		}
	}
}

func (n *Node) QuorumCallfutureQFWithRequestArg(ctx context.Context, in *Request, replyChan chan<- internalResponse) {
	reply := new(Response)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ZorumsService/QuorumCallfutureQFWithRequestArg", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalResponse{n.id, reply, err}
}

// QuorumCallFutureCustomReturnType with custom_return_type option.
func (c *Configuration) QuorumCallFutureCustomReturnType(ctx context.Context, in *Request, opts ...grpc.CallOption) *FutureMyResponse {
	fut := &FutureMyResponse{
		NodeIDs: make([]uint32, 0, c.n),
		c:       make(chan struct{}, 1),
	}
	go func() {
		defer close(fut.c)
		c.quorumCallFutureCustomReturnType(ctx, in, fut, opts...)
	}()
	return fut
}

func (c *Configuration) quorumCallFutureCustomReturnType(ctx context.Context, in *Request, resp *FutureMyResponse, opts ...grpc.CallOption) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "QuorumCallFutureCustomReturnType")
		defer ti.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = time.Until(deadline)
		}
		ti.LazyLog(&ti.firstLine, false)
		ti.LazyLog(&payload{sent: true, msg: in}, false)

		defer func() {
			ti.LazyLog(&qcresult{ids: resp.NodeIDs, reply: resp.MyResponse, err: resp.err}, false)
			if resp.err != nil {
				ti.SetError()
			}
		}()
	}

	expected := c.n
	replyChan := make(chan internalResponse, expected)
	for _, n := range c.nodes {
		go n.QuorumCallFutureCustomReturnType(ctx, in, replyChan)
	}

	var (
		replyValues = make([]*Response, 0, c.n)
		reply       *MyResponse
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			resp.NodeIDs = append(resp.NodeIDs, r.nid)
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			if reply, quorum = c.qspec.QuorumCallFutureCustomReturnTypeQF(replyValues); quorum {
				resp.MyResponse, resp.err = reply, nil
				return
			}
		case <-ctx.Done():
			resp.MyResponse, resp.err = reply, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
			return
		}
		if len(errs)+len(replyValues) == expected {
			resp.MyResponse, resp.err = reply, QuorumCallError{"incomplete call", len(replyValues), errs}
			return
		}
	}
}

func (n *Node) QuorumCallFutureCustomReturnType(ctx context.Context, in *Request, replyChan chan<- internalResponse) {
	reply := new(Response)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ZorumsService/QuorumCallFutureCustomReturnType", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalResponse{n.id, reply, err}
}

// QuorumCallFutureCombo with all supported options.
func (c *Configuration) QuorumCallFutureCombo(ctx context.Context, in *Request, f func(*Request, uint32) *Request, opts ...grpc.CallOption) *FutureMyResponse {
	fut := &FutureMyResponse{
		NodeIDs: make([]uint32, 0, c.n),
		c:       make(chan struct{}, 1),
	}
	go func() {
		defer close(fut.c)
		c.quorumCallFutureCombo(ctx, in, f, fut, opts...)
	}()
	return fut
}

func (c *Configuration) quorumCallFutureCombo(ctx context.Context, in *Request, f func(*Request, uint32) *Request, resp *FutureMyResponse, opts ...grpc.CallOption) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "QuorumCallFutureCombo")
		defer ti.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = time.Until(deadline)
		}
		ti.LazyLog(&ti.firstLine, false)
		ti.LazyLog(&payload{sent: true, msg: in}, false)

		defer func() {
			ti.LazyLog(&qcresult{ids: resp.NodeIDs, reply: resp.MyResponse, err: resp.err}, false)
			if resp.err != nil {
				ti.SetError()
			}
		}()
	}

	expected := c.n
	replyChan := make(chan internalResponse, expected)
	for _, n := range c.nodes {
		nodeArg := f(in, n.id)
		if nodeArg == nil {
			expected--
			continue
		}
		go n.QuorumCallFutureCombo(ctx, nodeArg, replyChan)
	}

	var (
		replyValues = make([]*Response, 0, c.n)
		reply       *MyResponse
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			resp.NodeIDs = append(resp.NodeIDs, r.nid)
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			if reply, quorum = c.qspec.QuorumCallFutureComboQF(in, replyValues); quorum {
				resp.MyResponse, resp.err = reply, nil
				return
			}
		case <-ctx.Done():
			resp.MyResponse, resp.err = reply, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
			return
		}
		if len(errs)+len(replyValues) == expected {
			resp.MyResponse, resp.err = reply, QuorumCallError{"incomplete call", len(replyValues), errs}
			return
		}
	}
}

func (n *Node) QuorumCallFutureCombo(ctx context.Context, in *Request, replyChan chan<- internalResponse) {
	reply := new(Response)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ZorumsService/QuorumCallFutureCombo", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalResponse{n.id, reply, err}
}

// QuorumCallFuture2 plain; with same return type: Response.
func (c *Configuration) QuorumCallFuture2(ctx context.Context, in *Request, opts ...grpc.CallOption) *FutureResponse {
	fut := &FutureResponse{
		NodeIDs: make([]uint32, 0, c.n),
		c:       make(chan struct{}, 1),
	}
	go func() {
		defer close(fut.c)
		c.quorumCallFuture2(ctx, in, fut, opts...)
	}()
	return fut
}

func (c *Configuration) quorumCallFuture2(ctx context.Context, in *Request, resp *FutureResponse, opts ...grpc.CallOption) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "QuorumCallFuture2")
		defer ti.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = time.Until(deadline)
		}
		ti.LazyLog(&ti.firstLine, false)
		ti.LazyLog(&payload{sent: true, msg: in}, false)

		defer func() {
			ti.LazyLog(&qcresult{ids: resp.NodeIDs, reply: resp.Response, err: resp.err}, false)
			if resp.err != nil {
				ti.SetError()
			}
		}()
	}

	expected := c.n
	replyChan := make(chan internalResponse, expected)
	for _, n := range c.nodes {
		go n.QuorumCallFuture2(ctx, in, replyChan)
	}

	var (
		replyValues = make([]*Response, 0, c.n)
		reply       *Response
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			resp.NodeIDs = append(resp.NodeIDs, r.nid)
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			if reply, quorum = c.qspec.QuorumCallFuture2QF(replyValues); quorum {
				resp.Response, resp.err = reply, nil
				return
			}
		case <-ctx.Done():
			resp.Response, resp.err = reply, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
			return
		}
		if len(errs)+len(replyValues) == expected {
			resp.Response, resp.err = reply, QuorumCallError{"incomplete call", len(replyValues), errs}
			return
		}
	}
}

func (n *Node) QuorumCallFuture2(ctx context.Context, in *Request, replyChan chan<- internalResponse) {
	reply := new(Response)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ZorumsService/QuorumCallFuture2", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalResponse{n.id, reply, err}
}

// QuorumCallFutureEmpty for testing imported message type.
func (c *Configuration) QuorumCallFutureEmpty(ctx context.Context, in *Request, opts ...grpc.CallOption) *FutureEmpty {
	fut := &FutureEmpty{
		NodeIDs: make([]uint32, 0, c.n),
		c:       make(chan struct{}, 1),
	}
	go func() {
		defer close(fut.c)
		c.quorumCallFutureEmpty(ctx, in, fut, opts...)
	}()
	return fut
}

func (c *Configuration) quorumCallFutureEmpty(ctx context.Context, in *Request, resp *FutureEmpty, opts ...grpc.CallOption) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "QuorumCallFutureEmpty")
		defer ti.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = time.Until(deadline)
		}
		ti.LazyLog(&ti.firstLine, false)
		ti.LazyLog(&payload{sent: true, msg: in}, false)

		defer func() {
			ti.LazyLog(&qcresult{ids: resp.NodeIDs, reply: resp.Empty, err: resp.err}, false)
			if resp.err != nil {
				ti.SetError()
			}
		}()
	}

	expected := c.n
	replyChan := make(chan internalEmpty, expected)
	for _, n := range c.nodes {
		go n.QuorumCallFutureEmpty(ctx, in, replyChan)
	}

	var (
		replyValues = make([]*empty.Empty, 0, c.n)
		reply       *empty.Empty
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			resp.NodeIDs = append(resp.NodeIDs, r.nid)
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			if reply, quorum = c.qspec.QuorumCallFutureEmptyQF(replyValues); quorum {
				resp.Empty, resp.err = reply, nil
				return
			}
		case <-ctx.Done():
			resp.Empty, resp.err = reply, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
			return
		}
		if len(errs)+len(replyValues) == expected {
			resp.Empty, resp.err = reply, QuorumCallError{"incomplete call", len(replyValues), errs}
			return
		}
	}
}

func (n *Node) QuorumCallFutureEmpty(ctx context.Context, in *Request, replyChan chan<- internalEmpty) {
	reply := new(empty.Empty)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ZorumsService/QuorumCallFutureEmpty", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalEmpty{n.id, reply, err}
}

// QuorumCallFutureEmpty2 for testing imported message type; with same return
// type as QuorumCallFuture: Response.
func (c *Configuration) QuorumCallFutureEmpty2(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) *FutureResponse {
	fut := &FutureResponse{
		NodeIDs: make([]uint32, 0, c.n),
		c:       make(chan struct{}, 1),
	}
	go func() {
		defer close(fut.c)
		c.quorumCallFutureEmpty2(ctx, in, fut, opts...)
	}()
	return fut
}

func (c *Configuration) quorumCallFutureEmpty2(ctx context.Context, in *empty.Empty, resp *FutureResponse, opts ...grpc.CallOption) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "QuorumCallFutureEmpty2")
		defer ti.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = time.Until(deadline)
		}
		ti.LazyLog(&ti.firstLine, false)
		ti.LazyLog(&payload{sent: true, msg: in}, false)

		defer func() {
			ti.LazyLog(&qcresult{ids: resp.NodeIDs, reply: resp.Response, err: resp.err}, false)
			if resp.err != nil {
				ti.SetError()
			}
		}()
	}

	expected := c.n
	replyChan := make(chan internalResponse, expected)
	for _, n := range c.nodes {
		go n.QuorumCallFutureEmpty2(ctx, in, replyChan)
	}

	var (
		replyValues = make([]*Response, 0, c.n)
		reply       *Response
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			resp.NodeIDs = append(resp.NodeIDs, r.nid)
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			if reply, quorum = c.qspec.QuorumCallFutureEmpty2QF(replyValues); quorum {
				resp.Response, resp.err = reply, nil
				return
			}
		case <-ctx.Done():
			resp.Response, resp.err = reply, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
			return
		}
		if len(errs)+len(replyValues) == expected {
			resp.Response, resp.err = reply, QuorumCallError{"incomplete call", len(replyValues), errs}
			return
		}
	}
}

func (n *Node) QuorumCallFutureEmpty2(ctx context.Context, in *empty.Empty, replyChan chan<- internalResponse) {
	reply := new(Response)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ZorumsService/QuorumCallFutureEmpty2", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalResponse{n.id, reply, err}
}
