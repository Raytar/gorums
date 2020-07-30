// Code generated by protoc-gen-gorums. DO NOT EDIT.

package dev

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	gorums "github.com/relab/gorums"
)

// Unicast is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (n *Node) Unicast(ctx context.Context, in *Request) {

	cd := gorums.CallData{
		Manager:  n.mgr.Manager,
		Node:     n.Node,
		Message:  in,
		MethodID: unicastMethodID,
	}

	gorums.Unicast(ctx, cd)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ empty.Empty

// Unicast2 is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (n *Node) Unicast2(ctx context.Context, in *Request) {

	cd := gorums.CallData{
		Manager:  n.mgr.Manager,
		Node:     n.Node,
		Message:  in,
		MethodID: unicast2MethodID,
	}

	gorums.Unicast(ctx, cd)
}
