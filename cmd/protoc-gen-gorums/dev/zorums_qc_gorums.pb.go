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

// QuorumCall plain.
func (c *Configuration) QuorumCall(ctx context.Context, in *Request, opts ...grpc.CallOption) (resp *Response, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "QuorumCall")
		defer ti.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = time.Until(deadline)
		}
		ti.LazyLog(&ti.firstLine, false)
		ti.LazyLog(&payload{sent: true, msg: in}, false)

		defer func() {
			ti.LazyLog(&qcresult{reply: resp, err: err}, false)
			if err != nil {
				ti.SetError()
			}
		}()
	}

	expected := c.n
	replyChan := make(chan internalResponse, expected)
	for _, n := range c.nodes {
		go n.QuorumCall(ctx, in, replyChan)
	}

	var (
		replyValues = make([]*Response, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			if resp, quorum = c.qspec.QuorumCallQF(replyValues); quorum {
				return resp, nil
			}
		case <-ctx.Done():
			return resp, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
		}
		if len(errs)+len(replyValues) == expected {
			return resp, QuorumCallError{"incomplete call", len(replyValues), errs}
		}
	}
}

func (n *Node) QuorumCall(ctx context.Context, in *Request, replyChan chan<- internalResponse) {
	reply := new(Response)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ZorumsService/QuorumCall", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalResponse{n.id, reply, err}
}

// QuorumCall with per_node_arg option.
func (c *Configuration) QuorumCallPerNodeArg(ctx context.Context, in *Request, f func(*Request, uint32) *Request, opts ...grpc.CallOption) (resp *Response, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "QuorumCallPerNodeArg")
		defer ti.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = time.Until(deadline)
		}
		ti.LazyLog(&ti.firstLine, false)
		ti.LazyLog(&payload{sent: true, msg: in}, false)

		defer func() {
			ti.LazyLog(&qcresult{reply: resp, err: err}, false)
			if err != nil {
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
		go n.QuorumCallPerNodeArg(ctx, nodeArg, replyChan)
	}

	var (
		replyValues = make([]*Response, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			if resp, quorum = c.qspec.QuorumCallPerNodeArgQF(replyValues); quorum {
				return resp, nil
			}
		case <-ctx.Done():
			return resp, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
		}
		if len(errs)+len(replyValues) == expected {
			return resp, QuorumCallError{"incomplete call", len(replyValues), errs}
		}
	}
}

func (n *Node) QuorumCallPerNodeArg(ctx context.Context, in *Request, replyChan chan<- internalResponse) {
	reply := new(Response)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ZorumsService/QuorumCallPerNodeArg", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalResponse{n.id, reply, err}
}

// QuorumCall with qf_with_req option.
func (c *Configuration) QuorumCallQFWithRequestArg(ctx context.Context, in *Request, opts ...grpc.CallOption) (resp *Response, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "QuorumCallQFWithRequestArg")
		defer ti.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = time.Until(deadline)
		}
		ti.LazyLog(&ti.firstLine, false)
		ti.LazyLog(&payload{sent: true, msg: in}, false)

		defer func() {
			ti.LazyLog(&qcresult{reply: resp, err: err}, false)
			if err != nil {
				ti.SetError()
			}
		}()
	}

	expected := c.n
	replyChan := make(chan internalResponse, expected)
	for _, n := range c.nodes {
		go n.QuorumCallQFWithRequestArg(ctx, in, replyChan)
	}

	var (
		replyValues = make([]*Response, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			if resp, quorum = c.qspec.QuorumCallQFWithRequestArgQF(in, replyValues); quorum {
				return resp, nil
			}
		case <-ctx.Done():
			return resp, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
		}
		if len(errs)+len(replyValues) == expected {
			return resp, QuorumCallError{"incomplete call", len(replyValues), errs}
		}
	}
}

func (n *Node) QuorumCallQFWithRequestArg(ctx context.Context, in *Request, replyChan chan<- internalResponse) {
	reply := new(Response)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ZorumsService/QuorumCallQFWithRequestArg", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalResponse{n.id, reply, err}
}

// QuorumCall with custom_return_type option.
func (c *Configuration) QuorumCallCustomReturnType(ctx context.Context, in *Request, opts ...grpc.CallOption) (resp *MyResponse, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "QuorumCallCustomReturnType")
		defer ti.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = time.Until(deadline)
		}
		ti.LazyLog(&ti.firstLine, false)
		ti.LazyLog(&payload{sent: true, msg: in}, false)

		defer func() {
			ti.LazyLog(&qcresult{reply: resp, err: err}, false)
			if err != nil {
				ti.SetError()
			}
		}()
	}

	expected := c.n
	replyChan := make(chan internalResponse, expected)
	for _, n := range c.nodes {
		go n.QuorumCallCustomReturnType(ctx, in, replyChan)
	}

	var (
		replyValues = make([]*Response, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			if resp, quorum = c.qspec.QuorumCallCustomReturnTypeQF(replyValues); quorum {
				return resp, nil
			}
		case <-ctx.Done():
			return resp, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
		}
		if len(errs)+len(replyValues) == expected {
			return resp, QuorumCallError{"incomplete call", len(replyValues), errs}
		}
	}
}

func (n *Node) QuorumCallCustomReturnType(ctx context.Context, in *Request, replyChan chan<- internalResponse) {
	reply := new(Response)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ZorumsService/QuorumCallCustomReturnType", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalResponse{n.id, reply, err}
}

// QuorumCallCombo with all supported options.
func (c *Configuration) QuorumCallCombo(ctx context.Context, in *Request, f func(*Request, uint32) *Request, opts ...grpc.CallOption) (resp *MyResponse, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "QuorumCallCombo")
		defer ti.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = time.Until(deadline)
		}
		ti.LazyLog(&ti.firstLine, false)
		ti.LazyLog(&payload{sent: true, msg: in}, false)

		defer func() {
			ti.LazyLog(&qcresult{reply: resp, err: err}, false)
			if err != nil {
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
		go n.QuorumCallCombo(ctx, nodeArg, replyChan)
	}

	var (
		replyValues = make([]*Response, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			if resp, quorum = c.qspec.QuorumCallComboQF(in, replyValues); quorum {
				return resp, nil
			}
		case <-ctx.Done():
			return resp, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
		}
		if len(errs)+len(replyValues) == expected {
			return resp, QuorumCallError{"incomplete call", len(replyValues), errs}
		}
	}
}

func (n *Node) QuorumCallCombo(ctx context.Context, in *Request, replyChan chan<- internalResponse) {
	reply := new(Response)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ZorumsService/QuorumCallCombo", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalResponse{n.id, reply, err}
}

// QuorumCallEmpty for testing imported message type.
func (c *Configuration) QuorumCallEmpty(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (resp *Response, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "QuorumCallEmpty")
		defer ti.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = time.Until(deadline)
		}
		ti.LazyLog(&ti.firstLine, false)
		ti.LazyLog(&payload{sent: true, msg: in}, false)

		defer func() {
			ti.LazyLog(&qcresult{reply: resp, err: err}, false)
			if err != nil {
				ti.SetError()
			}
		}()
	}

	expected := c.n
	replyChan := make(chan internalResponse, expected)
	for _, n := range c.nodes {
		go n.QuorumCallEmpty(ctx, in, replyChan)
	}

	var (
		replyValues = make([]*Response, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			if resp, quorum = c.qspec.QuorumCallEmptyQF(replyValues); quorum {
				return resp, nil
			}
		case <-ctx.Done():
			return resp, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
		}
		if len(errs)+len(replyValues) == expected {
			return resp, QuorumCallError{"incomplete call", len(replyValues), errs}
		}
	}
}

func (n *Node) QuorumCallEmpty(ctx context.Context, in *empty.Empty, replyChan chan<- internalResponse) {
	reply := new(Response)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ZorumsService/QuorumCallEmpty", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalResponse{n.id, reply, err}
}

// QuorumCallEmpty2 for testing imported message type.
func (c *Configuration) QuorumCallEmpty2(ctx context.Context, in *Request, opts ...grpc.CallOption) (resp *empty.Empty, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "QuorumCallEmpty2")
		defer ti.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = time.Until(deadline)
		}
		ti.LazyLog(&ti.firstLine, false)
		ti.LazyLog(&payload{sent: true, msg: in}, false)

		defer func() {
			ti.LazyLog(&qcresult{reply: resp, err: err}, false)
			if err != nil {
				ti.SetError()
			}
		}()
	}

	expected := c.n
	replyChan := make(chan internalEmpty, expected)
	for _, n := range c.nodes {
		go n.QuorumCallEmpty2(ctx, in, replyChan)
	}

	var (
		replyValues = make([]*empty.Empty, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			if resp, quorum = c.qspec.QuorumCallEmpty2QF(replyValues); quorum {
				return resp, nil
			}
		case <-ctx.Done():
			return resp, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
		}
		if len(errs)+len(replyValues) == expected {
			return resp, QuorumCallError{"incomplete call", len(replyValues), errs}
		}
	}
}

func (n *Node) QuorumCallEmpty2(ctx context.Context, in *Request, replyChan chan<- internalEmpty) {
	reply := new(empty.Empty)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ZorumsService/QuorumCallEmpty2", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalEmpty{n.id, reply, err}
}
