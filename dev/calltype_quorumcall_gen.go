// DO NOT EDIT. Generated by 'gorums' plugin for protoc-gen-go
// Source file to edit is: calltype_quorumcall_tmpl

package dev

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

/* Methods on Configuration and the quorum call struct Read */

//TODO Make this a customizable struct that replaces FQRespName together with typedecl option in gogoprotobuf.
//(This file could maybe hold all types of structs for the different call semantics)

// ReadReply encapsulates the reply from a Read quorum call.
// It contains the id of each node of the quorum that replied and a single reply.
type ReadReply struct {
	// the actual reply
	*State
	NodeIDs []uint32
	err     error
}

func (r ReadReply) String() string {
	return fmt.Sprintf("node ids: %v | answer: %v", r.NodeIDs, r.State)
}

type readArg *ReadRequest

// Read is invoked as a quorum call on all nodes in configuration c,
// using the same argument arg, and returns the result as a ReadReply.
func (c *Configuration) Read(ctx context.Context, arg *ReadRequest) (*ReadReply, error) {
	return c.read(ctx, arg)
}

/* Methods on Manager for quorum call method Read */

type readReply struct {
	nid   uint32
	reply *State
	err   error
}

func (c *Configuration) read(ctx context.Context, a readArg) (resp *ReadReply, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.tr = trace.New("gorums."+c.tstring()+".Sent", "Read")
		defer ti.tr.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = deadline.Sub(time.Now())
		}
		ti.tr.LazyLog(&ti.firstLine, false)

		defer func() {
			ti.tr.LazyLog(&qcresult{
				ids:   resp.NodeIDs,
				reply: resp.State,
				err:   resp.err,
			}, false)
			if resp.err != nil {
				ti.tr.SetError()
			}
		}()
	}

	replyChan := make(chan readReply, c.n)

	if c.mgr.opts.trace {
		ti.tr.LazyLog(&payload{sent: true, msg: a}, false)
	}

	for _, n := range c.nodes {
		go callGRPCRead(ctx, n, a, replyChan)
	}

	resp = &ReadReply{NodeIDs: make([]uint32, 0, c.n)}
	var (
		replyValues = make([]*State, 0, c.n)
		errCount    int
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			resp.NodeIDs = append(resp.NodeIDs, r.nid)
			if r.err != nil {
				errCount++
				break
			}
			if c.mgr.opts.trace {
				ti.tr.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}
			replyValues = append(replyValues, r.reply)
			if resp.State, quorum = c.qspec.ReadQF(replyValues); quorum {
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

func callGRPCRead(ctx context.Context, node *Node, args *ReadRequest, replyChan chan<- readReply) {
	reply := new(State)
	start := time.Now()
	err := grpc.Invoke(
		ctx,
		"/dev.Register/Read",
		args,
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

/* Methods on Configuration and the quorum call struct ReadCustomReturn */

//TODO Make this a customizable struct that replaces FQRespName together with typedecl option in gogoprotobuf.
//(This file could maybe hold all types of structs for the different call semantics)

// ReadCustomReturnReply encapsulates the reply from a ReadCustomReturn quorum call.
// It contains the id of each node of the quorum that replied and a single reply.
type ReadCustomReturnReply struct {
	// the actual reply
	*State
	NodeIDs []uint32
	err     error
}

func (r ReadCustomReturnReply) String() string {
	return fmt.Sprintf("node ids: %v | answer: %v", r.NodeIDs, r.State)
}

type readCustomReturnArg *ReadRequest

// ReadCustomReturn is invoked as a quorum call on all nodes in configuration c,
// using the same argument arg, and returns the result as a ReadCustomReturnReply.
func (c *Configuration) ReadCustomReturn(ctx context.Context, arg *ReadRequest) (*ReadCustomReturnReply, error) {
	return c.readCustomReturn(ctx, arg)
}

/* Methods on Manager for quorum call method ReadCustomReturn */

type readCustomReturnReply struct {
	nid   uint32
	reply *State
	err   error
}

func (c *Configuration) readCustomReturn(ctx context.Context, a readCustomReturnArg) (resp *ReadCustomReturnReply, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.tr = trace.New("gorums."+c.tstring()+".Sent", "ReadCustomReturn")
		defer ti.tr.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = deadline.Sub(time.Now())
		}
		ti.tr.LazyLog(&ti.firstLine, false)

		defer func() {
			ti.tr.LazyLog(&qcresult{
				ids:   resp.NodeIDs,
				reply: resp.State,
				err:   resp.err,
			}, false)
			if resp.err != nil {
				ti.tr.SetError()
			}
		}()
	}

	replyChan := make(chan readCustomReturnReply, c.n)

	if c.mgr.opts.trace {
		ti.tr.LazyLog(&payload{sent: true, msg: a}, false)
	}

	for _, n := range c.nodes {
		go callGRPCReadCustomReturn(ctx, n, a, replyChan)
	}

	resp = &ReadCustomReturnReply{NodeIDs: make([]uint32, 0, c.n)}
	var (
		replyValues = make([]*State, 0, c.n)
		errCount    int
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			resp.NodeIDs = append(resp.NodeIDs, r.nid)
			if r.err != nil {
				errCount++
				break
			}
			if c.mgr.opts.trace {
				ti.tr.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}
			replyValues = append(replyValues, r.reply)
			if resp.State, quorum = c.qspec.ReadCustomReturnQF(replyValues); quorum {
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

func callGRPCReadCustomReturn(ctx context.Context, node *Node, args *ReadRequest, replyChan chan<- readCustomReturnReply) {
	reply := new(State)
	start := time.Now()
	err := grpc.Invoke(
		ctx,
		"/dev.Register/ReadCustomReturn",
		args,
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

/* Methods on Configuration and the quorum call struct Write */

//TODO Make this a customizable struct that replaces FQRespName together with typedecl option in gogoprotobuf.
//(This file could maybe hold all types of structs for the different call semantics)

// WriteReply encapsulates the reply from a Write quorum call.
// It contains the id of each node of the quorum that replied and a single reply.
type WriteReply struct {
	// the actual reply
	*WriteResponse
	NodeIDs []uint32
	err     error
}

func (r WriteReply) String() string {
	return fmt.Sprintf("node ids: %v | answer: %v", r.NodeIDs, r.WriteResponse)
}

type writeArg *State

// Write is invoked as a quorum call on all nodes in configuration c,
// using the same argument arg, and returns the result as a WriteReply.
func (c *Configuration) Write(ctx context.Context, arg *State) (*WriteReply, error) {
	return c.write(ctx, arg)
}

/* Methods on Manager for quorum call method Write */

type writeReply struct {
	nid   uint32
	reply *WriteResponse
	err   error
}

func (c *Configuration) write(ctx context.Context, a writeArg) (resp *WriteReply, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.tr = trace.New("gorums."+c.tstring()+".Sent", "Write")
		defer ti.tr.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = deadline.Sub(time.Now())
		}
		ti.tr.LazyLog(&ti.firstLine, false)

		defer func() {
			ti.tr.LazyLog(&qcresult{
				ids:   resp.NodeIDs,
				reply: resp.WriteResponse,
				err:   resp.err,
			}, false)
			if resp.err != nil {
				ti.tr.SetError()
			}
		}()
	}

	replyChan := make(chan writeReply, c.n)

	if c.mgr.opts.trace {
		ti.tr.LazyLog(&payload{sent: true, msg: a}, false)
	}

	for _, n := range c.nodes {
		go callGRPCWrite(ctx, n, a, replyChan)
	}

	resp = &WriteReply{NodeIDs: make([]uint32, 0, c.n)}
	var (
		replyValues = make([]*WriteResponse, 0, c.n)
		errCount    int
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			resp.NodeIDs = append(resp.NodeIDs, r.nid)
			if r.err != nil {
				errCount++
				break
			}
			if c.mgr.opts.trace {
				ti.tr.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}
			replyValues = append(replyValues, r.reply)
			if resp.WriteResponse, quorum = c.qspec.WriteQF(a, replyValues); quorum {
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

func callGRPCWrite(ctx context.Context, node *Node, args *State, replyChan chan<- writeReply) {
	reply := new(WriteResponse)
	start := time.Now()
	err := grpc.Invoke(
		ctx,
		"/dev.Register/Write",
		args,
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

/* Methods on Configuration and the quorum call struct WritePerNode */

//TODO Make this a customizable struct that replaces FQRespName together with typedecl option in gogoprotobuf.
//(This file could maybe hold all types of structs for the different call semantics)

// WritePerNodeReply encapsulates the reply from a WritePerNode quorum call.
// It contains the id of each node of the quorum that replied and a single reply.
type WritePerNodeReply struct {
	// the actual reply
	*WriteResponse
	NodeIDs []uint32
	err     error
}

func (r WritePerNodeReply) String() string {
	return fmt.Sprintf("node ids: %v | answer: %v", r.NodeIDs, r.WriteResponse)
}

type writePerNodeArg func(nodeID uint32) *State

// WritePerNode is invoked as a quorum call on each node in configuration c,
// with the argument returned by the provided perNode function and returns the
// result as a WritePerNodeReply. The perNode function returns a *State
// object to be passed to the given nodeID.
func (c *Configuration) WritePerNode(ctx context.Context, perNode func(nodeID uint32) *State) (*WritePerNodeReply, error) {
	return c.writePerNode(ctx, perNode)
}

/* Methods on Manager for quorum call method WritePerNode */

type writePerNodeReply struct {
	nid   uint32
	reply *WriteResponse
	err   error
}

func (c *Configuration) writePerNode(ctx context.Context, a writePerNodeArg) (resp *WritePerNodeReply, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.tr = trace.New("gorums."+c.tstring()+".Sent", "WritePerNode")
		defer ti.tr.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = deadline.Sub(time.Now())
		}
		ti.tr.LazyLog(&ti.firstLine, false)

		defer func() {
			ti.tr.LazyLog(&qcresult{
				ids:   resp.NodeIDs,
				reply: resp.WriteResponse,
				err:   resp.err,
			}, false)
			if resp.err != nil {
				ti.tr.SetError()
			}
		}()
	}

	replyChan := make(chan writePerNodeReply, c.n)

	if c.mgr.opts.trace {
		ti.tr.LazyLog(&payload{sent: true, msg: a}, false)
	}

	for _, n := range c.nodes {
		go callGRPCWritePerNode(ctx, n, a(n.id), replyChan)
	}

	resp = &WritePerNodeReply{NodeIDs: make([]uint32, 0, c.n)}
	var (
		replyValues = make([]*WriteResponse, 0, c.n)
		errCount    int
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			resp.NodeIDs = append(resp.NodeIDs, r.nid)
			if r.err != nil {
				errCount++
				break
			}
			if c.mgr.opts.trace {
				ti.tr.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}
			replyValues = append(replyValues, r.reply)
			if resp.WriteResponse, quorum = c.qspec.WritePerNodeQF(replyValues); quorum {
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

func callGRPCWritePerNode(ctx context.Context, node *Node, args *State, replyChan chan<- writePerNodeReply) {
	reply := new(WriteResponse)
	start := time.Now()
	err := grpc.Invoke(
		ctx,
		"/dev.Register/WritePerNode",
		args,
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
