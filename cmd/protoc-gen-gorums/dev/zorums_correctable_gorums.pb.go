// Code generated by protoc-gen-gorums. DO NOT EDIT.

package dev

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	trace "golang.org/x/net/trace"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	time "time"
)

// Correctable plain.
func (c *Configuration) Correctable(ctx context.Context, in *Request, opts ...grpc.CallOption) *CorrectableResponse {
	corr := &CorrectableResponse{
		level:   LevelNotSet,
		NodeIDs: make([]uint32, 0, c.n),
		donech:  make(chan struct{}),
	}
	go c.correctable(ctx, in, corr, opts...)
	return corr
}

func (c *Configuration) correctable(ctx context.Context, in *Request, resp *CorrectableResponse, opts ...grpc.CallOption) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "Correctable")
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
		go n.Correctable(ctx, in, replyChan)
	}

	var (
		replyValues = make([]*Response, 0, c.n)
		clevel      = LevelNotSet
		reply       *Response
		rlevel      int
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
			reply, rlevel, quorum = c.qspec.CorrectableQF(in, replyValues)
			if quorum {
				resp.set(reply, rlevel, nil, true)
				return
			}
			if rlevel > clevel {
				clevel = rlevel
				resp.set(reply, rlevel, nil, false)
			}
		case <-ctx.Done():
			resp.set(reply, clevel, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}, true)
			return
		}
		if len(errs)+len(replyValues) == expected {
			resp.set(reply, clevel, QuorumCallError{"incomplete call", len(replyValues), errs}, true)
			return
		}
	}
}

func (n *Node) Correctable(ctx context.Context, in *Request, replyChan chan<- internalResponse) {
	reply := new(Response)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ZorumsService/Correctable", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalResponse{n.id, reply, err}
}

// CorrectablePerNodeArg with per_node_arg option.
func (c *Configuration) CorrectablePerNodeArg(ctx context.Context, in *Request, f func(*Request, uint32) *Request, opts ...grpc.CallOption) *CorrectableResponse {
	corr := &CorrectableResponse{
		level:   LevelNotSet,
		NodeIDs: make([]uint32, 0, c.n),
		donech:  make(chan struct{}),
	}
	go c.correctablePerNodeArg(ctx, in, f, corr, opts...)
	return corr
}

func (c *Configuration) correctablePerNodeArg(ctx context.Context, in *Request, f func(*Request, uint32) *Request, resp *CorrectableResponse, opts ...grpc.CallOption) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "CorrectablePerNodeArg")
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
		go n.CorrectablePerNodeArg(ctx, nodeArg, replyChan)
	}

	var (
		replyValues = make([]*Response, 0, c.n)
		clevel      = LevelNotSet
		reply       *Response
		rlevel      int
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
			reply, rlevel, quorum = c.qspec.CorrectablePerNodeArgQF(in, replyValues)
			if quorum {
				resp.set(reply, rlevel, nil, true)
				return
			}
			if rlevel > clevel {
				clevel = rlevel
				resp.set(reply, rlevel, nil, false)
			}
		case <-ctx.Done():
			resp.set(reply, clevel, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}, true)
			return
		}
		if len(errs)+len(replyValues) == expected {
			resp.set(reply, clevel, QuorumCallError{"incomplete call", len(replyValues), errs}, true)
			return
		}
	}
}

func (n *Node) CorrectablePerNodeArg(ctx context.Context, in *Request, replyChan chan<- internalResponse) {
	reply := new(Response)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ZorumsService/CorrectablePerNodeArg", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalResponse{n.id, reply, err}
}

// CorrectableCustomReturnType with custom_return_type option.
func (c *Configuration) CorrectableCustomReturnType(ctx context.Context, in *Request, opts ...grpc.CallOption) *CorrectableMyResponse {
	corr := &CorrectableMyResponse{
		level:   LevelNotSet,
		NodeIDs: make([]uint32, 0, c.n),
		donech:  make(chan struct{}),
	}
	go c.correctableCustomReturnType(ctx, in, corr, opts...)
	return corr
}

func (c *Configuration) correctableCustomReturnType(ctx context.Context, in *Request, resp *CorrectableMyResponse, opts ...grpc.CallOption) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "CorrectableCustomReturnType")
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
		go n.CorrectableCustomReturnType(ctx, in, replyChan)
	}

	var (
		replyValues = make([]*Response, 0, c.n)
		clevel      = LevelNotSet
		reply       *MyResponse
		rlevel      int
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
			reply, rlevel, quorum = c.qspec.CorrectableCustomReturnTypeQF(in, replyValues)
			if quorum {
				resp.set(reply, rlevel, nil, true)
				return
			}
			if rlevel > clevel {
				clevel = rlevel
				resp.set(reply, rlevel, nil, false)
			}
		case <-ctx.Done():
			resp.set(reply, clevel, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}, true)
			return
		}
		if len(errs)+len(replyValues) == expected {
			resp.set(reply, clevel, QuorumCallError{"incomplete call", len(replyValues), errs}, true)
			return
		}
	}
}

func (n *Node) CorrectableCustomReturnType(ctx context.Context, in *Request, replyChan chan<- internalResponse) {
	reply := new(Response)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ZorumsService/CorrectableCustomReturnType", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalResponse{n.id, reply, err}
}

// CorrectableCombo with all supported options.
func (c *Configuration) CorrectableCombo(ctx context.Context, in *Request, f func(*Request, uint32) *Request, opts ...grpc.CallOption) *CorrectableMyResponse {
	corr := &CorrectableMyResponse{
		level:   LevelNotSet,
		NodeIDs: make([]uint32, 0, c.n),
		donech:  make(chan struct{}),
	}
	go c.correctableCombo(ctx, in, f, corr, opts...)
	return corr
}

func (c *Configuration) correctableCombo(ctx context.Context, in *Request, f func(*Request, uint32) *Request, resp *CorrectableMyResponse, opts ...grpc.CallOption) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "CorrectableCombo")
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
		go n.CorrectableCombo(ctx, nodeArg, replyChan)
	}

	var (
		replyValues = make([]*Response, 0, c.n)
		clevel      = LevelNotSet
		reply       *MyResponse
		rlevel      int
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
			reply, rlevel, quorum = c.qspec.CorrectableComboQF(in, replyValues)
			if quorum {
				resp.set(reply, rlevel, nil, true)
				return
			}
			if rlevel > clevel {
				clevel = rlevel
				resp.set(reply, rlevel, nil, false)
			}
		case <-ctx.Done():
			resp.set(reply, clevel, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}, true)
			return
		}
		if len(errs)+len(replyValues) == expected {
			resp.set(reply, clevel, QuorumCallError{"incomplete call", len(replyValues), errs}, true)
			return
		}
	}
}

func (n *Node) CorrectableCombo(ctx context.Context, in *Request, replyChan chan<- internalResponse) {
	reply := new(Response)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ZorumsService/CorrectableCombo", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalResponse{n.id, reply, err}
}

// CorrectableEmpty for testing imported message type.
func (c *Configuration) CorrectableEmpty(ctx context.Context, in *Request, opts ...grpc.CallOption) *CorrectableEmpty {
	corr := &CorrectableEmpty{
		level:   LevelNotSet,
		NodeIDs: make([]uint32, 0, c.n),
		donech:  make(chan struct{}),
	}
	go c.correctableEmpty(ctx, in, corr, opts...)
	return corr
}

func (c *Configuration) correctableEmpty(ctx context.Context, in *Request, resp *CorrectableEmpty, opts ...grpc.CallOption) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "CorrectableEmpty")
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
		go n.CorrectableEmpty(ctx, in, replyChan)
	}

	var (
		replyValues = make([]*empty.Empty, 0, c.n)
		clevel      = LevelNotSet
		reply       *empty.Empty
		rlevel      int
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
			reply, rlevel, quorum = c.qspec.CorrectableEmptyQF(in, replyValues)
			if quorum {
				resp.set(reply, rlevel, nil, true)
				return
			}
			if rlevel > clevel {
				clevel = rlevel
				resp.set(reply, rlevel, nil, false)
			}
		case <-ctx.Done():
			resp.set(reply, clevel, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}, true)
			return
		}
		if len(errs)+len(replyValues) == expected {
			resp.set(reply, clevel, QuorumCallError{"incomplete call", len(replyValues), errs}, true)
			return
		}
	}
}

func (n *Node) CorrectableEmpty(ctx context.Context, in *Request, replyChan chan<- internalEmpty) {
	reply := new(empty.Empty)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ZorumsService/CorrectableEmpty", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalEmpty{n.id, reply, err}
}

// CorrectableEmpty2 for testing imported message type; with same return
// type as Correctable: Response.
func (c *Configuration) CorrectableEmpty2(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) *CorrectableResponse {
	corr := &CorrectableResponse{
		level:   LevelNotSet,
		NodeIDs: make([]uint32, 0, c.n),
		donech:  make(chan struct{}),
	}
	go c.correctableEmpty2(ctx, in, corr, opts...)
	return corr
}

func (c *Configuration) correctableEmpty2(ctx context.Context, in *empty.Empty, resp *CorrectableResponse, opts ...grpc.CallOption) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "CorrectableEmpty2")
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
		go n.CorrectableEmpty2(ctx, in, replyChan)
	}

	var (
		replyValues = make([]*Response, 0, c.n)
		clevel      = LevelNotSet
		reply       *Response
		rlevel      int
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
			reply, rlevel, quorum = c.qspec.CorrectableEmpty2QF(in, replyValues)
			if quorum {
				resp.set(reply, rlevel, nil, true)
				return
			}
			if rlevel > clevel {
				clevel = rlevel
				resp.set(reply, rlevel, nil, false)
			}
		case <-ctx.Done():
			resp.set(reply, clevel, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}, true)
			return
		}
		if len(errs)+len(replyValues) == expected {
			resp.set(reply, clevel, QuorumCallError{"incomplete call", len(replyValues), errs}, true)
			return
		}
	}
}

func (n *Node) CorrectableEmpty2(ctx context.Context, in *empty.Empty, replyChan chan<- internalResponse) {
	reply := new(Response)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ZorumsService/CorrectableEmpty2", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalResponse{n.id, reply, err}
}

// CorrectableStream plain.
func (c *Configuration) CorrectableStream(ctx context.Context, in *Request, opts ...grpc.CallOption) *CorrectableStreamResponse {
	corr := &CorrectableStreamResponse{
		level:   LevelNotSet,
		NodeIDs: make([]uint32, 0, c.n),
		donech:  make(chan struct{}),
	}
	go c.correctableStream(ctx, in, corr, opts...)
	return corr
}

func (c *Configuration) correctableStream(ctx context.Context, in *Request, resp *CorrectableStreamResponse, opts ...grpc.CallOption) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "CorrectableStream")
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
		go n.CorrectableStream(ctx, in, replyChan)
	}

	var (
		//TODO(meling) don't recall why we need n*2 reply slots?
		replyValues = make([]*Response, 0, c.n*2)
		clevel      = LevelNotSet
		reply       *Response
		rlevel      int
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			resp.NodeIDs = appendIfNotPresent(resp.NodeIDs, r.nid)
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			reply, rlevel, quorum = c.qspec.CorrectableStreamQF(in, replyValues)
			if quorum {
				resp.set(reply, rlevel, nil, true)
				return
			}
			if rlevel > clevel {
				clevel = rlevel
				resp.set(reply, rlevel, nil, false)
			}
		case <-ctx.Done():
			resp.set(reply, clevel, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}, true)
			return
		}
		if len(errs) == expected { // Can't rely on reply count.
			resp.set(reply, clevel, QuorumCallError{"incomplete call", len(replyValues), errs}, true)
			return
		}
	}
}

func (n *Node) CorrectableStream(ctx context.Context, in *Request, replyChan chan<- internalResponse) {
	x := NewZorumsServiceClient(n.conn)
	y, err := x.CorrectableStream(ctx, in)
	if err != nil {
		replyChan <- internalResponse{n.id, nil, err}
		return
	}

	for {
		reply, err := y.Recv()
		if err == io.EOF {
			return
		}
		replyChan <- internalResponse{n.id, reply, err}
		if err != nil {
			return
		}
	}
}

// CorrectablePerNodeArg with per_node_arg option.
func (c *Configuration) CorrectableStreamPerNodeArg(ctx context.Context, in *Request, f func(*Request, uint32) *Request, opts ...grpc.CallOption) *CorrectableStreamResponse {
	corr := &CorrectableStreamResponse{
		level:   LevelNotSet,
		NodeIDs: make([]uint32, 0, c.n),
		donech:  make(chan struct{}),
	}
	go c.correctableStreamPerNodeArg(ctx, in, f, corr, opts...)
	return corr
}

func (c *Configuration) correctableStreamPerNodeArg(ctx context.Context, in *Request, f func(*Request, uint32) *Request, resp *CorrectableStreamResponse, opts ...grpc.CallOption) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "CorrectableStreamPerNodeArg")
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
		go n.CorrectableStreamPerNodeArg(ctx, nodeArg, replyChan)
	}

	var (
		//TODO(meling) don't recall why we need n*2 reply slots?
		replyValues = make([]*Response, 0, c.n*2)
		clevel      = LevelNotSet
		reply       *Response
		rlevel      int
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			resp.NodeIDs = appendIfNotPresent(resp.NodeIDs, r.nid)
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			reply, rlevel, quorum = c.qspec.CorrectableStreamPerNodeArgQF(in, replyValues)
			if quorum {
				resp.set(reply, rlevel, nil, true)
				return
			}
			if rlevel > clevel {
				clevel = rlevel
				resp.set(reply, rlevel, nil, false)
			}
		case <-ctx.Done():
			resp.set(reply, clevel, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}, true)
			return
		}
		if len(errs) == expected { // Can't rely on reply count.
			resp.set(reply, clevel, QuorumCallError{"incomplete call", len(replyValues), errs}, true)
			return
		}
	}
}

func (n *Node) CorrectableStreamPerNodeArg(ctx context.Context, in *Request, replyChan chan<- internalResponse) {
	x := NewZorumsServiceClient(n.conn)
	y, err := x.CorrectableStreamPerNodeArg(ctx, in)
	if err != nil {
		replyChan <- internalResponse{n.id, nil, err}
		return
	}

	for {
		reply, err := y.Recv()
		if err == io.EOF {
			return
		}
		replyChan <- internalResponse{n.id, reply, err}
		if err != nil {
			return
		}
	}
}

// CorrectableCustomReturnType with custom_return_type option.
func (c *Configuration) CorrectableStreamCustomReturnType(ctx context.Context, in *Request, opts ...grpc.CallOption) *CorrectableStreamMyResponse {
	corr := &CorrectableStreamMyResponse{
		level:   LevelNotSet,
		NodeIDs: make([]uint32, 0, c.n),
		donech:  make(chan struct{}),
	}
	go c.correctableStreamCustomReturnType(ctx, in, corr, opts...)
	return corr
}

func (c *Configuration) correctableStreamCustomReturnType(ctx context.Context, in *Request, resp *CorrectableStreamMyResponse, opts ...grpc.CallOption) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "CorrectableStreamCustomReturnType")
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
		go n.CorrectableStreamCustomReturnType(ctx, in, replyChan)
	}

	var (
		//TODO(meling) don't recall why we need n*2 reply slots?
		replyValues = make([]*Response, 0, c.n*2)
		clevel      = LevelNotSet
		reply       *MyResponse
		rlevel      int
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			resp.NodeIDs = appendIfNotPresent(resp.NodeIDs, r.nid)
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			reply, rlevel, quorum = c.qspec.CorrectableStreamCustomReturnTypeQF(in, replyValues)
			if quorum {
				resp.set(reply, rlevel, nil, true)
				return
			}
			if rlevel > clevel {
				clevel = rlevel
				resp.set(reply, rlevel, nil, false)
			}
		case <-ctx.Done():
			resp.set(reply, clevel, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}, true)
			return
		}
		if len(errs) == expected { // Can't rely on reply count.
			resp.set(reply, clevel, QuorumCallError{"incomplete call", len(replyValues), errs}, true)
			return
		}
	}
}

func (n *Node) CorrectableStreamCustomReturnType(ctx context.Context, in *Request, replyChan chan<- internalResponse) {
	x := NewZorumsServiceClient(n.conn)
	y, err := x.CorrectableStreamCustomReturnType(ctx, in)
	if err != nil {
		replyChan <- internalResponse{n.id, nil, err}
		return
	}

	for {
		reply, err := y.Recv()
		if err == io.EOF {
			return
		}
		replyChan <- internalResponse{n.id, reply, err}
		if err != nil {
			return
		}
	}
}

// CorrectableCombo with all supported options.
func (c *Configuration) CorrectableStreamCombo(ctx context.Context, in *Request, f func(*Request, uint32) *Request, opts ...grpc.CallOption) *CorrectableStreamMyResponse {
	corr := &CorrectableStreamMyResponse{
		level:   LevelNotSet,
		NodeIDs: make([]uint32, 0, c.n),
		donech:  make(chan struct{}),
	}
	go c.correctableStreamCombo(ctx, in, f, corr, opts...)
	return corr
}

func (c *Configuration) correctableStreamCombo(ctx context.Context, in *Request, f func(*Request, uint32) *Request, resp *CorrectableStreamMyResponse, opts ...grpc.CallOption) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "CorrectableStreamCombo")
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
		go n.CorrectableStreamCombo(ctx, nodeArg, replyChan)
	}

	var (
		//TODO(meling) don't recall why we need n*2 reply slots?
		replyValues = make([]*Response, 0, c.n*2)
		clevel      = LevelNotSet
		reply       *MyResponse
		rlevel      int
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			resp.NodeIDs = appendIfNotPresent(resp.NodeIDs, r.nid)
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			reply, rlevel, quorum = c.qspec.CorrectableStreamComboQF(in, replyValues)
			if quorum {
				resp.set(reply, rlevel, nil, true)
				return
			}
			if rlevel > clevel {
				clevel = rlevel
				resp.set(reply, rlevel, nil, false)
			}
		case <-ctx.Done():
			resp.set(reply, clevel, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}, true)
			return
		}
		if len(errs) == expected { // Can't rely on reply count.
			resp.set(reply, clevel, QuorumCallError{"incomplete call", len(replyValues), errs}, true)
			return
		}
	}
}

func (n *Node) CorrectableStreamCombo(ctx context.Context, in *Request, replyChan chan<- internalResponse) {
	x := NewZorumsServiceClient(n.conn)
	y, err := x.CorrectableStreamCombo(ctx, in)
	if err != nil {
		replyChan <- internalResponse{n.id, nil, err}
		return
	}

	for {
		reply, err := y.Recv()
		if err == io.EOF {
			return
		}
		replyChan <- internalResponse{n.id, reply, err}
		if err != nil {
			return
		}
	}
}

// CorrectableEmpty for testing imported message type.
func (c *Configuration) CorrectableStreamEmpty(ctx context.Context, in *Request, opts ...grpc.CallOption) *CorrectableStreamEmpty {
	corr := &CorrectableStreamEmpty{
		level:   LevelNotSet,
		NodeIDs: make([]uint32, 0, c.n),
		donech:  make(chan struct{}),
	}
	go c.correctableStreamEmpty(ctx, in, corr, opts...)
	return corr
}

func (c *Configuration) correctableStreamEmpty(ctx context.Context, in *Request, resp *CorrectableStreamEmpty, opts ...grpc.CallOption) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "CorrectableStreamEmpty")
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
		go n.CorrectableStreamEmpty(ctx, in, replyChan)
	}

	var (
		//TODO(meling) don't recall why we need n*2 reply slots?
		replyValues = make([]*empty.Empty, 0, c.n*2)
		clevel      = LevelNotSet
		reply       *empty.Empty
		rlevel      int
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			resp.NodeIDs = appendIfNotPresent(resp.NodeIDs, r.nid)
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			reply, rlevel, quorum = c.qspec.CorrectableStreamEmptyQF(in, replyValues)
			if quorum {
				resp.set(reply, rlevel, nil, true)
				return
			}
			if rlevel > clevel {
				clevel = rlevel
				resp.set(reply, rlevel, nil, false)
			}
		case <-ctx.Done():
			resp.set(reply, clevel, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}, true)
			return
		}
		if len(errs) == expected { // Can't rely on reply count.
			resp.set(reply, clevel, QuorumCallError{"incomplete call", len(replyValues), errs}, true)
			return
		}
	}
}

func (n *Node) CorrectableStreamEmpty(ctx context.Context, in *Request, replyChan chan<- internalEmpty) {
	x := NewZorumsServiceClient(n.conn)
	y, err := x.CorrectableStreamEmpty(ctx, in)
	if err != nil {
		replyChan <- internalEmpty{n.id, nil, err}
		return
	}

	for {
		reply, err := y.Recv()
		if err == io.EOF {
			return
		}
		replyChan <- internalEmpty{n.id, reply, err}
		if err != nil {
			return
		}
	}
}

// CorrectableEmpty2 for testing imported message type; with same return
// type as Correctable: Response.
func (c *Configuration) CorrectableStreamEmpty2(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) *CorrectableStreamResponse {
	corr := &CorrectableStreamResponse{
		level:   LevelNotSet,
		NodeIDs: make([]uint32, 0, c.n),
		donech:  make(chan struct{}),
	}
	go c.correctableStreamEmpty2(ctx, in, corr, opts...)
	return corr
}

func (c *Configuration) correctableStreamEmpty2(ctx context.Context, in *empty.Empty, resp *CorrectableStreamResponse, opts ...grpc.CallOption) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "CorrectableStreamEmpty2")
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
		go n.CorrectableStreamEmpty2(ctx, in, replyChan)
	}

	var (
		//TODO(meling) don't recall why we need n*2 reply slots?
		replyValues = make([]*Response, 0, c.n*2)
		clevel      = LevelNotSet
		reply       *Response
		rlevel      int
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			resp.NodeIDs = appendIfNotPresent(resp.NodeIDs, r.nid)
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			reply, rlevel, quorum = c.qspec.CorrectableStreamEmpty2QF(in, replyValues)
			if quorum {
				resp.set(reply, rlevel, nil, true)
				return
			}
			if rlevel > clevel {
				clevel = rlevel
				resp.set(reply, rlevel, nil, false)
			}
		case <-ctx.Done():
			resp.set(reply, clevel, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}, true)
			return
		}
		if len(errs) == expected { // Can't rely on reply count.
			resp.set(reply, clevel, QuorumCallError{"incomplete call", len(replyValues), errs}, true)
			return
		}
	}
}

func (n *Node) CorrectableStreamEmpty2(ctx context.Context, in *empty.Empty, replyChan chan<- internalResponse) {
	x := NewZorumsServiceClient(n.conn)
	y, err := x.CorrectableStreamEmpty2(ctx, in)
	if err != nil {
		replyChan <- internalResponse{n.id, nil, err}
		return
	}

	for {
		reply, err := y.Recv()
		if err == io.EOF {
			return
		}
		replyChan <- internalResponse{n.id, reply, err}
		if err != nil {
			return
		}
	}
}
