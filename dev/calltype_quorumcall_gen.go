// DO NOT EDIT. Generated by 'gorums' plugin for protoc-gen-go
// Source file to edit is: calltype_quorumcall_tmpl

package dev

import (
	"sync"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

/* Exported types and methods for quorum call method Read */

// Read is invoked as a quorum call on all nodes in configuration c,
// using the same argument arg, and returns the result.
func (c *Configuration) Read(ctx context.Context, arg *ReadRequest) (*State, error) {
	return c.read(ctx, arg)
}

/* Unexported types and methods for quorum call method Read */

type readReply struct {
	nid   uint32
	reply *State
	err   error
}

func (c *Configuration) read(ctx context.Context, a *ReadRequest) (resp *State, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.tr = trace.New("gorums."+c.tstring()+".Sent", "Read")
		defer ti.tr.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = deadline.Sub(time.Now())
		}
		ti.tr.LazyLog(&ti.firstLine, false)
		ti.tr.LazyLog(&payload{sent: true, msg: a}, false)

		defer func() {
			ti.tr.LazyLog(&qcresult{
				reply: resp,
				err:   err,
			}, false)
			if err != nil {
				ti.tr.SetError()
			}
		}()
	}

	replyChan := make(chan readReply, c.n)
	var wg sync.WaitGroup
	wg.Add(c.n)
	for _, n := range c.nodes {
		go callGRPCRead(ctx, &wg, n, a, replyChan)
	}
	wg.Wait()

	var (
		replyValues = make([]*State, 0, c.n)
		errCount    int
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			if r.err != nil {
				errCount++
				break
			}
			if c.mgr.opts.trace {
				ti.tr.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}
			replyValues = append(replyValues, r.reply)
			if resp, quorum = c.qspec.ReadQF(replyValues); quorum {
				return resp, nil
			}
		case <-ctx.Done():
			return resp, QuorumCallError{ctx.Err().Error(), errCount, len(replyValues)}
		}

		if errCount+len(replyValues) == c.n {
			return resp, QuorumCallError{"incomplete call", errCount, len(replyValues)}
		}
	}
}

func callGRPCRead(ctx context.Context, wg *sync.WaitGroup, node *Node, arg *ReadRequest, replyChan chan<- readReply) {
	wg.Done()
	reply := new(State)
	start := time.Now()
	err := grpc.Invoke(
		ctx,
		"/dev.Register/Read",
		arg,
		reply,
		node.conn,
	)
	switch grpc.Code(err) { // nil -> codes.OK
	case codes.OK, codes.Canceled:
		node.setLatency(time.Since(start))
	default:
		node.setLastErr(err)
	}
	replyChan <- readReply{node.id, reply, err}
}

/* Exported types and methods for quorum call method ReadCustomReturn */

// ReadCustomReturn is invoked as a quorum call on all nodes in configuration c,
// using the same argument arg, and returns the result.
func (c *Configuration) ReadCustomReturn(ctx context.Context, arg *ReadRequest) (*MyState, error) {
	return c.readCustomReturn(ctx, arg)
}

/* Unexported types and methods for quorum call method ReadCustomReturn */

type readCustomReturnReply struct {
	nid   uint32
	reply *State
	err   error
}

func (c *Configuration) readCustomReturn(ctx context.Context, a *ReadRequest) (resp *MyState, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.tr = trace.New("gorums."+c.tstring()+".Sent", "ReadCustomReturn")
		defer ti.tr.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = deadline.Sub(time.Now())
		}
		ti.tr.LazyLog(&ti.firstLine, false)
		ti.tr.LazyLog(&payload{sent: true, msg: a}, false)

		defer func() {
			ti.tr.LazyLog(&qcresult{
				reply: resp,
				err:   err,
			}, false)
			if err != nil {
				ti.tr.SetError()
			}
		}()
	}

	replyChan := make(chan readCustomReturnReply, c.n)
	var wg sync.WaitGroup
	wg.Add(c.n)
	for _, n := range c.nodes {
		go callGRPCReadCustomReturn(ctx, &wg, n, a, replyChan)
	}
	wg.Wait()

	var (
		replyValues = make([]*State, 0, c.n)
		errCount    int
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			if r.err != nil {
				errCount++
				break
			}
			if c.mgr.opts.trace {
				ti.tr.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}
			replyValues = append(replyValues, r.reply)
			if resp, quorum = c.qspec.ReadCustomReturnQF(replyValues); quorum {
				return resp, nil
			}
		case <-ctx.Done():
			return resp, QuorumCallError{ctx.Err().Error(), errCount, len(replyValues)}
		}

		if errCount+len(replyValues) == c.n {
			return resp, QuorumCallError{"incomplete call", errCount, len(replyValues)}
		}
	}
}

func callGRPCReadCustomReturn(ctx context.Context, wg *sync.WaitGroup, node *Node, arg *ReadRequest, replyChan chan<- readCustomReturnReply) {
	wg.Done()
	reply := new(State)
	start := time.Now()
	err := grpc.Invoke(
		ctx,
		"/dev.Register/ReadCustomReturn",
		arg,
		reply,
		node.conn,
	)
	switch grpc.Code(err) { // nil -> codes.OK
	case codes.OK, codes.Canceled:
		node.setLatency(time.Since(start))
	default:
		node.setLastErr(err)
	}
	replyChan <- readCustomReturnReply{node.id, reply, err}
}

/* Exported types and methods for quorum call method Write */

// Write is invoked as a quorum call on all nodes in configuration c,
// using the same argument arg, and returns the result.
func (c *Configuration) Write(ctx context.Context, arg *State) (*WriteResponse, error) {
	return c.write(ctx, arg)
}

/* Unexported types and methods for quorum call method Write */

type writeReply struct {
	nid   uint32
	reply *WriteResponse
	err   error
}

func (c *Configuration) write(ctx context.Context, a *State) (resp *WriteResponse, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.tr = trace.New("gorums."+c.tstring()+".Sent", "Write")
		defer ti.tr.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = deadline.Sub(time.Now())
		}
		ti.tr.LazyLog(&ti.firstLine, false)
		ti.tr.LazyLog(&payload{sent: true, msg: a}, false)

		defer func() {
			ti.tr.LazyLog(&qcresult{
				reply: resp,
				err:   err,
			}, false)
			if err != nil {
				ti.tr.SetError()
			}
		}()
	}

	replyChan := make(chan writeReply, c.n)
	var wg sync.WaitGroup
	wg.Add(c.n)
	for _, n := range c.nodes {
		go callGRPCWrite(ctx, &wg, n, a, replyChan)
	}
	wg.Wait()

	var (
		replyValues = make([]*WriteResponse, 0, c.n)
		errCount    int
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			if r.err != nil {
				errCount++
				break
			}
			if c.mgr.opts.trace {
				ti.tr.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}
			replyValues = append(replyValues, r.reply)
			if resp, quorum = c.qspec.WriteQF(a, replyValues); quorum {
				return resp, nil
			}
		case <-ctx.Done():
			return resp, QuorumCallError{ctx.Err().Error(), errCount, len(replyValues)}
		}

		if errCount+len(replyValues) == c.n {
			return resp, QuorumCallError{"incomplete call", errCount, len(replyValues)}
		}
	}
}

func callGRPCWrite(ctx context.Context, wg *sync.WaitGroup, node *Node, arg *State, replyChan chan<- writeReply) {
	wg.Done()
	reply := new(WriteResponse)
	start := time.Now()
	err := grpc.Invoke(
		ctx,
		"/dev.Register/Write",
		arg,
		reply,
		node.conn,
	)
	switch grpc.Code(err) { // nil -> codes.OK
	case codes.OK, codes.Canceled:
		node.setLatency(time.Since(start))
	default:
		node.setLastErr(err)
	}
	replyChan <- writeReply{node.id, reply, err}
}

/* Exported types and methods for quorum call method WritePerNode */

// WritePerNode is invoked as a quorum call on each node in configuration c,
// with the argument returned by the provided perNode function and returns the
// result. The perNode function takes a request arg and
// returns a State object to be passed to the given nodeID.
// The perNode function should be thread-safe.
func (c *Configuration) WritePerNode(ctx context.Context, arg *State, perNode func(arg State, nodeID uint32) *State) (*WriteResponse, error) {
	return c.writePerNode(ctx, arg, perNode)
}

/* Unexported types and methods for quorum call method WritePerNode */

type writePerNodeReply struct {
	nid   uint32
	reply *WriteResponse
	err   error
}

func (c *Configuration) writePerNode(ctx context.Context, a *State, f func(arg State, nodeID uint32) *State) (resp *WriteResponse, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.tr = trace.New("gorums."+c.tstring()+".Sent", "WritePerNode")
		defer ti.tr.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = deadline.Sub(time.Now())
		}
		ti.tr.LazyLog(&ti.firstLine, false)
		ti.tr.LazyLog(&payload{sent: true, msg: a}, false)

		defer func() {
			ti.tr.LazyLog(&qcresult{
				reply: resp,
				err:   err,
			}, false)
			if err != nil {
				ti.tr.SetError()
			}
		}()
	}

	replyChan := make(chan writePerNodeReply, c.n)
	var wg sync.WaitGroup
	wg.Add(c.n)
	for _, n := range c.nodes {
		go callGRPCWritePerNode(ctx, &wg, n, f(*a, n.id), replyChan)
	}
	wg.Wait()

	var (
		replyValues = make([]*WriteResponse, 0, c.n)
		errCount    int
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			if r.err != nil {
				errCount++
				break
			}
			if c.mgr.opts.trace {
				ti.tr.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}
			replyValues = append(replyValues, r.reply)
			if resp, quorum = c.qspec.WritePerNodeQF(replyValues); quorum {
				return resp, nil
			}
		case <-ctx.Done():
			return resp, QuorumCallError{ctx.Err().Error(), errCount, len(replyValues)}
		}

		if errCount+len(replyValues) == c.n {
			return resp, QuorumCallError{"incomplete call", errCount, len(replyValues)}
		}
	}
}

func callGRPCWritePerNode(ctx context.Context, wg *sync.WaitGroup, node *Node, arg *State, replyChan chan<- writePerNodeReply) {
	wg.Done()
	reply := new(WriteResponse)
	start := time.Now()
	err := grpc.Invoke(
		ctx,
		"/dev.Register/WritePerNode",
		arg,
		reply,
		node.conn,
	)
	switch grpc.Code(err) { // nil -> codes.OK
	case codes.OK, codes.Canceled:
		node.setLatency(time.Since(start))
	default:
		node.setLastErr(err)
	}
	replyChan <- writePerNodeReply{node.id, reply, err}
}
