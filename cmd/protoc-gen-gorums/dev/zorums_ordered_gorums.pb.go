// Code generated by protoc-gen-gorums. DO NOT EDIT.

package dev

import (
	context "context"
	fmt "fmt"
	ordering "github.com/relab/gorums/ordering"
	trace "golang.org/x/net/trace"
	grpc "google.golang.org/grpc"
	time "time"
)

// OrderingQC is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) OrderingQC(ctx context.Context, in *Request) (resp *Response, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "OrderingQC")
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

	// get the ID which will be used to return the correct responses for a request
	msgID := c.mgr.nextMsgID()

	// set up a channel to collect replies
	replies := make(chan *orderingResult, c.n)
	c.mgr.putChan(msgID, replies)

	// remove the replies channel when we are done
	defer c.mgr.deleteChan(msgID)

	data, err := marshaler.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}
	msg := &ordering.Message{
		ID:       msgID,
		MethodID: orderingQCMethodID,
		Data:     data,
	}
	// push the message to the nodes
	expected := c.n
	for _, n := range c.nodes {
		n.sendQ <- msg
	}

	var (
		replyValues = make([]*Response, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replies:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			reply := new(Response)
			err := unmarshaler.Unmarshal(r.reply, reply)
			if err != nil {
				errs = append(errs, GRPCError{r.nid, fmt.Errorf("failed to unmarshal reply: %w", err)})
				break
			}
			replyValues = append(replyValues, reply)
			if resp, quorum = c.qspec.OrderingQCQF(in, replyValues); quorum {
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

// OrderingPerNodeArg is a quorum call invoked on each node in configuration c,
// with the argument returned by the provided function f, and returns the combined result.
// The per node function f receives a copy of the Request request argument and
// returns a Request manipulated to be passed to the given nodeID.
// The function f must be thread-safe.
func (c *Configuration) OrderingPerNodeArg(ctx context.Context, in *Request, f func(*Request, uint32) *Request) (resp *Response, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "OrderingPerNodeArg")
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

	// get the ID which will be used to return the correct responses for a request
	msgID := c.mgr.nextMsgID()

	// set up a channel to collect replies
	replies := make(chan *orderingResult, c.n)
	c.mgr.putChan(msgID, replies)

	// remove the replies channel when we are done
	defer c.mgr.deleteChan(msgID)

	// push the message to the nodes
	expected := c.n
	for _, n := range c.nodes {
		nodeArg := f(in, n.ID())
		if nodeArg == nil {
			expected--
			continue
		}
		data, err := marshaler.Marshal(nodeArg)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal message: %w", err)
		}
		msg := &ordering.Message{
			ID:       msgID,
			MethodID: orderingPerNodeArgMethodID,
			Data:     data,
		}
		n.sendQ <- msg
	}

	var (
		replyValues = make([]*Response, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replies:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			reply := new(Response)
			err := unmarshaler.Unmarshal(r.reply, reply)
			if err != nil {
				errs = append(errs, GRPCError{r.nid, fmt.Errorf("failed to unmarshal reply: %w", err)})
				break
			}
			replyValues = append(replyValues, reply)
			if resp, quorum = c.qspec.OrderingPerNodeArgQF(in, replyValues); quorum {
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

// OrderingCustomReturnType is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) OrderingCustomReturnType(ctx context.Context, in *Request) (resp *MyResponse, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "OrderingCustomReturnType")
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

	// get the ID which will be used to return the correct responses for a request
	msgID := c.mgr.nextMsgID()

	// set up a channel to collect replies
	replies := make(chan *orderingResult, c.n)
	c.mgr.putChan(msgID, replies)

	// remove the replies channel when we are done
	defer c.mgr.deleteChan(msgID)

	data, err := marshaler.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}
	msg := &ordering.Message{
		ID:       msgID,
		MethodID: orderingCustomReturnTypeMethodID,
		Data:     data,
	}
	// push the message to the nodes
	expected := c.n
	for _, n := range c.nodes {
		n.sendQ <- msg
	}

	var (
		replyValues = make([]*Response, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replies:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			reply := new(Response)
			err := unmarshaler.Unmarshal(r.reply, reply)
			if err != nil {
				errs = append(errs, GRPCError{r.nid, fmt.Errorf("failed to unmarshal reply: %w", err)})
				break
			}
			replyValues = append(replyValues, reply)
			if resp, quorum = c.qspec.OrderingCustomReturnTypeQF(in, replyValues); quorum {
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

// OrderingCombo is a quorum call invoked on each node in configuration c,
// with the argument returned by the provided function f, and returns the combined result.
// The per node function f receives a copy of the Request request argument and
// returns a Request manipulated to be passed to the given nodeID.
// The function f must be thread-safe.
func (c *Configuration) OrderingCombo(ctx context.Context, in *Request, f func(*Request, uint32) *Request) (resp *MyResponse, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "OrderingCombo")
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

	// get the ID which will be used to return the correct responses for a request
	msgID := c.mgr.nextMsgID()

	// set up a channel to collect replies
	replies := make(chan *orderingResult, c.n)
	c.mgr.putChan(msgID, replies)

	// remove the replies channel when we are done
	defer c.mgr.deleteChan(msgID)

	// push the message to the nodes
	expected := c.n
	for _, n := range c.nodes {
		nodeArg := f(in, n.ID())
		if nodeArg == nil {
			expected--
			continue
		}
		data, err := marshaler.Marshal(nodeArg)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal message: %w", err)
		}
		msg := &ordering.Message{
			ID:       msgID,
			MethodID: orderingComboMethodID,
			Data:     data,
		}
		n.sendQ <- msg
	}

	var (
		replyValues = make([]*Response, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replies:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			reply := new(Response)
			err := unmarshaler.Unmarshal(r.reply, reply)
			if err != nil {
				errs = append(errs, GRPCError{r.nid, fmt.Errorf("failed to unmarshal reply: %w", err)})
				break
			}
			replyValues = append(replyValues, reply)
			if resp, quorum = c.qspec.OrderingComboQF(in, replyValues); quorum {
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

// OrderingConcurrent is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) OrderingConcurrent(ctx context.Context, in *Request) (resp *Response, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "OrderingConcurrent")
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

	// get the ID which will be used to return the correct responses for a request
	msgID := c.mgr.nextMsgID()

	// set up a channel to collect replies
	replies := make(chan *orderingResult, c.n)
	c.mgr.putChan(msgID, replies)

	// remove the replies channel when we are done
	defer c.mgr.deleteChan(msgID)

	data, err := marshaler.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}
	msg := &ordering.Message{
		ID:       msgID,
		MethodID: orderingConcurrentMethodID,
		Data:     data,
	}
	// push the message to the nodes
	expected := c.n
	for _, n := range c.nodes {
		n.sendQ <- msg
	}

	var (
		replyValues = make([]*Response, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replies:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			reply := new(Response)
			err := unmarshaler.Unmarshal(r.reply, reply)
			if err != nil {
				errs = append(errs, GRPCError{r.nid, fmt.Errorf("failed to unmarshal reply: %w", err)})
				break
			}
			replyValues = append(replyValues, reply)
			if resp, quorum = c.qspec.OrderingConcurrentQF(in, replyValues); quorum {
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

// OrderingUnaryRPC is testing that we can create ordered Unary RPCs
func (n *Node) OrderingUnaryRPC(ctx context.Context, in *Request, opts ...grpc.CallOption) (resp *Response, err error) {

	// get the ID which will be used to return the correct responses for a request
	msgID := n.nextMsgID()

	// set up a channel to collect replies
	replies := make(chan *orderingResult, 1)
	n.putChan(msgID, replies)

	// remove the replies channel when we are done
	defer n.deleteChan(msgID)

	data, err := marshaler.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}
	msg := &ordering.Message{
		ID:       msgID,
		MethodID: orderingUnaryRPCMethodID,
		Data:     data,
	}
	n.sendQ <- msg

	select {
	case r := <-replies:
		if r.err != nil {
			return nil, r.err
		}
		reply := new(Response)
		err := unmarshaler.Unmarshal(r.reply, reply)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal reply: %w", err)
		}
		return reply, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// OrderingUnaryRPC is testing that we can create ordered Unary RPCs with
// concurrent handlers
func (n *Node) OrderingUnaryRPCConcurrent(ctx context.Context, in *Request, opts ...grpc.CallOption) (resp *Response, err error) {

	// get the ID which will be used to return the correct responses for a request
	msgID := n.nextMsgID()

	// set up a channel to collect replies
	replies := make(chan *orderingResult, 1)
	n.putChan(msgID, replies)

	// remove the replies channel when we are done
	defer n.deleteChan(msgID)

	data, err := marshaler.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}
	msg := &ordering.Message{
		ID:       msgID,
		MethodID: orderingUnaryRPCConcurrentMethodID,
		Data:     data,
	}
	n.sendQ <- msg

	select {
	case r := <-replies:
		if r.err != nil {
			return nil, r.err
		}
		reply := new(Response)
		err := unmarshaler.Unmarshal(r.reply, reply)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal reply: %w", err)
		}
		return reply, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// OrderingFuture asynchronously invokes a quorum call on configuration c
// and returns a FutureResponse, which can be used to inspect the quorum call
// reply and error when available.
func (c *Configuration) OrderingFuture(ctx context.Context, in *Request) *FutureResponse {
	fut := &FutureResponse{
		NodeIDs: make([]uint32, 0, c.n),
		c:       make(chan struct{}, 1),
	}
	// get the ID which will be used to return the correct responses for a request
	msgID := c.mgr.nextMsgID()

	// set up a channel to collect replies
	replyChan := make(chan *orderingResult, c.n)
	c.mgr.putChan(msgID, replyChan)

	expected := c.n

	var msg *ordering.Message
	data, err := marshaler.Marshal(in)
	if err != nil {
		// In case of a marshalling error, we should skip sending any messages
		fut.err = fmt.Errorf("failed to marshal message: %w", err)
		close(fut.c)
		return fut
	}
	msg = &ordering.Message{
		ID:       msgID,
		MethodID: orderingFutureMethodID,
		Data:     data,
	}

	// push the message to the nodes
	for _, n := range c.nodes {
		n.sendQ <- msg
	}

	go c.orderingFutureRecv(ctx, in, msgID, expected, replyChan, fut)

	return fut
}

func (c *Configuration) orderingFutureRecv(ctx context.Context, in *Request, msgID uint64, expected int, replyChan chan *orderingResult, fut *FutureResponse) {
	defer close(fut.c)

	if fut.err != nil {
		return
	}

	defer c.mgr.deleteChan(msgID)

	var (
		replyValues = make([]*Response, 0, c.n)
		reply       *Response
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			fut.NodeIDs = append(fut.NodeIDs, r.nid)
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}
			data := new(Response)
			err := unmarshaler.Unmarshal(r.reply, data)
			if err != nil {
				errs = append(errs, GRPCError{r.nid, fmt.Errorf("failed to unmarshal reply: %w", err)})
				break
			}
			replyValues = append(replyValues, data)
			if reply, quorum = c.qspec.OrderingFutureQF(in, replyValues); quorum {
				fut.Response, fut.err = reply, nil
				return
			}
		case <-ctx.Done():
			fut.Response, fut.err = reply, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
			return
		}
		if len(errs)+len(replyValues) == expected {
			fut.Response, fut.err = reply, QuorumCallError{"incomplete call", len(replyValues), errs}
			return
		}
	}
}

// OrderingFuturePerNodeArg asynchronously invokes a quorum call on each node in
// configuration c, with the argument returned by the provided function f
// and returns the result as a FutureResponse, which can be used to inspect
// the quorum call reply and error when available.
// The provide per node function f takes the provided Request argument
// and returns an Response object to be passed to the given nodeID.
// The per node function f should be thread-safe.
func (c *Configuration) OrderingFuturePerNodeArg(ctx context.Context, in *Request, f func(*Request, uint32) *Request) *FutureResponse {
	fut := &FutureResponse{
		NodeIDs: make([]uint32, 0, c.n),
		c:       make(chan struct{}, 1),
	}
	// get the ID which will be used to return the correct responses for a request
	msgID := c.mgr.nextMsgID()

	// set up a channel to collect replies
	replyChan := make(chan *orderingResult, c.n)
	c.mgr.putChan(msgID, replyChan)

	expected := c.n

	// push the message to the nodes
	for _, n := range c.nodes {
		nodeArg := f(in, n.ID())
		if nodeArg == nil {
			expected--
			continue
		}
		data, err := marshaler.Marshal(nodeArg)
		if err != nil {
			fut.err = fmt.Errorf("failed to marshal message: %w", err)
			close(fut.c)
			return fut
		}
		msg := &ordering.Message{
			ID:       msgID,
			MethodID: orderingFuturePerNodeArgMethodID,
			Data:     data,
		}
		n.sendQ <- msg
	}

	go c.orderingFuturePerNodeArgRecv(ctx, in, msgID, expected, replyChan, fut)

	return fut
}

func (c *Configuration) orderingFuturePerNodeArgRecv(ctx context.Context, in *Request, msgID uint64, expected int, replyChan chan *orderingResult, fut *FutureResponse) {
	defer close(fut.c)

	if fut.err != nil {
		return
	}

	defer c.mgr.deleteChan(msgID)

	var (
		replyValues = make([]*Response, 0, c.n)
		reply       *Response
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			fut.NodeIDs = append(fut.NodeIDs, r.nid)
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}
			data := new(Response)
			err := unmarshaler.Unmarshal(r.reply, data)
			if err != nil {
				errs = append(errs, GRPCError{r.nid, fmt.Errorf("failed to unmarshal reply: %w", err)})
				break
			}
			replyValues = append(replyValues, data)
			if reply, quorum = c.qspec.OrderingFuturePerNodeArgQF(in, replyValues); quorum {
				fut.Response, fut.err = reply, nil
				return
			}
		case <-ctx.Done():
			fut.Response, fut.err = reply, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
			return
		}
		if len(errs)+len(replyValues) == expected {
			fut.Response, fut.err = reply, QuorumCallError{"incomplete call", len(replyValues), errs}
			return
		}
	}
}

// OrderingFutureCustomReturnType asynchronously invokes a quorum call on configuration c
// and returns a FutureMyResponse, which can be used to inspect the quorum call
// reply and error when available.
func (c *Configuration) OrderingFutureCustomReturnType(ctx context.Context, in *Request) *FutureMyResponse {
	fut := &FutureMyResponse{
		NodeIDs: make([]uint32, 0, c.n),
		c:       make(chan struct{}, 1),
	}
	// get the ID which will be used to return the correct responses for a request
	msgID := c.mgr.nextMsgID()

	// set up a channel to collect replies
	replyChan := make(chan *orderingResult, c.n)
	c.mgr.putChan(msgID, replyChan)

	expected := c.n

	var msg *ordering.Message
	data, err := marshaler.Marshal(in)
	if err != nil {
		// In case of a marshalling error, we should skip sending any messages
		fut.err = fmt.Errorf("failed to marshal message: %w", err)
		close(fut.c)
		return fut
	}
	msg = &ordering.Message{
		ID:       msgID,
		MethodID: orderingFutureCustomReturnTypeMethodID,
		Data:     data,
	}

	// push the message to the nodes
	for _, n := range c.nodes {
		n.sendQ <- msg
	}

	go c.orderingFutureCustomReturnTypeRecv(ctx, in, msgID, expected, replyChan, fut)

	return fut
}

func (c *Configuration) orderingFutureCustomReturnTypeRecv(ctx context.Context, in *Request, msgID uint64, expected int, replyChan chan *orderingResult, fut *FutureMyResponse) {
	defer close(fut.c)

	if fut.err != nil {
		return
	}

	defer c.mgr.deleteChan(msgID)

	var (
		replyValues = make([]*Response, 0, c.n)
		reply       *MyResponse
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			fut.NodeIDs = append(fut.NodeIDs, r.nid)
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}
			data := new(Response)
			err := unmarshaler.Unmarshal(r.reply, data)
			if err != nil {
				errs = append(errs, GRPCError{r.nid, fmt.Errorf("failed to unmarshal reply: %w", err)})
				break
			}
			replyValues = append(replyValues, data)
			if reply, quorum = c.qspec.OrderingFutureCustomReturnTypeQF(in, replyValues); quorum {
				fut.MyResponse, fut.err = reply, nil
				return
			}
		case <-ctx.Done():
			fut.MyResponse, fut.err = reply, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
			return
		}
		if len(errs)+len(replyValues) == expected {
			fut.MyResponse, fut.err = reply, QuorumCallError{"incomplete call", len(replyValues), errs}
			return
		}
	}
}

// OrderingFutureConcurrent asynchronously invokes a quorum call on configuration c
// and returns a FutureResponse, which can be used to inspect the quorum call
// reply and error when available.
func (c *Configuration) OrderingFutureConcurrent(ctx context.Context, in *Request) *FutureResponse {
	fut := &FutureResponse{
		NodeIDs: make([]uint32, 0, c.n),
		c:       make(chan struct{}, 1),
	}
	// get the ID which will be used to return the correct responses for a request
	msgID := c.mgr.nextMsgID()

	// set up a channel to collect replies
	replyChan := make(chan *orderingResult, c.n)
	c.mgr.putChan(msgID, replyChan)

	expected := c.n

	var msg *ordering.Message
	data, err := marshaler.Marshal(in)
	if err != nil {
		// In case of a marshalling error, we should skip sending any messages
		fut.err = fmt.Errorf("failed to marshal message: %w", err)
		close(fut.c)
		return fut
	}
	msg = &ordering.Message{
		ID:       msgID,
		MethodID: orderingFutureConcurrentMethodID,
		Data:     data,
	}

	// push the message to the nodes
	for _, n := range c.nodes {
		n.sendQ <- msg
	}

	go c.orderingFutureConcurrentRecv(ctx, in, msgID, expected, replyChan, fut)

	return fut
}

func (c *Configuration) orderingFutureConcurrentRecv(ctx context.Context, in *Request, msgID uint64, expected int, replyChan chan *orderingResult, fut *FutureResponse) {
	defer close(fut.c)

	if fut.err != nil {
		return
	}

	defer c.mgr.deleteChan(msgID)

	var (
		replyValues = make([]*Response, 0, c.n)
		reply       *Response
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			fut.NodeIDs = append(fut.NodeIDs, r.nid)
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}
			data := new(Response)
			err := unmarshaler.Unmarshal(r.reply, data)
			if err != nil {
				errs = append(errs, GRPCError{r.nid, fmt.Errorf("failed to unmarshal reply: %w", err)})
				break
			}
			replyValues = append(replyValues, data)
			if reply, quorum = c.qspec.OrderingFutureConcurrentQF(in, replyValues); quorum {
				fut.Response, fut.err = reply, nil
				return
			}
		case <-ctx.Done():
			fut.Response, fut.err = reply, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
			return
		}
		if len(errs)+len(replyValues) == expected {
			fut.Response, fut.err = reply, QuorumCallError{"incomplete call", len(replyValues), errs}
			return
		}
	}
}

// OrderingFutureCombo asynchronously invokes a quorum call on each node in
// configuration c, with the argument returned by the provided function f
// and returns the result as a FutureMyResponse, which can be used to inspect
// the quorum call reply and error when available.
// The provide per node function f takes the provided Request argument
// and returns an Response object to be passed to the given nodeID.
// The per node function f should be thread-safe.
func (c *Configuration) OrderingFutureCombo(ctx context.Context, in *Request, f func(*Request, uint32) *Request) *FutureMyResponse {
	fut := &FutureMyResponse{
		NodeIDs: make([]uint32, 0, c.n),
		c:       make(chan struct{}, 1),
	}
	// get the ID which will be used to return the correct responses for a request
	msgID := c.mgr.nextMsgID()

	// set up a channel to collect replies
	replyChan := make(chan *orderingResult, c.n)
	c.mgr.putChan(msgID, replyChan)

	expected := c.n

	// push the message to the nodes
	for _, n := range c.nodes {
		nodeArg := f(in, n.ID())
		if nodeArg == nil {
			expected--
			continue
		}
		data, err := marshaler.Marshal(nodeArg)
		if err != nil {
			fut.err = fmt.Errorf("failed to marshal message: %w", err)
			close(fut.c)
			return fut
		}
		msg := &ordering.Message{
			ID:       msgID,
			MethodID: orderingFutureComboMethodID,
			Data:     data,
		}
		n.sendQ <- msg
	}

	go c.orderingFutureComboRecv(ctx, in, msgID, expected, replyChan, fut)

	return fut
}

func (c *Configuration) orderingFutureComboRecv(ctx context.Context, in *Request, msgID uint64, expected int, replyChan chan *orderingResult, fut *FutureMyResponse) {
	defer close(fut.c)

	if fut.err != nil {
		return
	}

	defer c.mgr.deleteChan(msgID)

	var (
		replyValues = make([]*Response, 0, c.n)
		reply       *MyResponse
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			fut.NodeIDs = append(fut.NodeIDs, r.nid)
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}
			data := new(Response)
			err := unmarshaler.Unmarshal(r.reply, data)
			if err != nil {
				errs = append(errs, GRPCError{r.nid, fmt.Errorf("failed to unmarshal reply: %w", err)})
				break
			}
			replyValues = append(replyValues, data)
			if reply, quorum = c.qspec.OrderingFutureComboQF(in, replyValues); quorum {
				fut.MyResponse, fut.err = reply, nil
				return
			}
		case <-ctx.Done():
			fut.MyResponse, fut.err = reply, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
			return
		}
		if len(errs)+len(replyValues) == expected {
			fut.MyResponse, fut.err = reply, QuorumCallError{"incomplete call", len(replyValues), errs}
			return
		}
	}
}