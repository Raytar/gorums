package metadata

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	status "google.golang.org/grpc/status"
)

type testSrv struct{}

func (srv testSrv) IDFromMD(ctx context.Context, _ *empty.Empty, ret func(*NodeID, error)) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		ret(nil, status.Error(codes.NotFound, "Metadata unavailable"))
		return
	}
	v := md.Get("id")
	if len(v) < 1 {
		ret(nil, status.Error(codes.NotFound, "ID field missing"))
		return
	}
	id, err := strconv.Atoi(v[0])
	if err != nil {
		ret(nil, status.Errorf(codes.InvalidArgument, "Got '%s', but could not convert to integer", v[0]))
		return
	}
	ret(&NodeID{ID: uint32(id)}, nil)
}

func (srv testSrv) WhatIP(ctx context.Context, _ *empty.Empty, ret func(*IPAddr, error)) {
	peerInfo, ok := peer.FromContext(ctx)
	if !ok {
		ret(nil, fmt.Errorf("Peer info unavailable"))
		return
	}
	ret(&IPAddr{Addr: peerInfo.Addr.String()}, nil)
}

func initServer(t *testing.T) *GorumsServer {
	srv := NewGorumsServer()
	srv.RegisterMetadataTestServer(&testSrv{})
	return srv
}

func TestMetadata(t *testing.T) {
	addrs, teardown := gorums.TestSetup(t, 1, func() interface{} { return initServer(t) })
	defer teardown()

	md := metadata.New(map[string]string{
		"id": "1",
	})

	mgr, err := NewManager(WithNodeList(addrs), WithMetadata(md), WithDialTimeout(1*time.Second), WithGrpcDialOptions(
		grpc.WithBlock(),
		grpc.WithInsecure(),
	))
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	node := mgr.Nodes()[0]
	resp, err := node.IDFromMD(context.Background(), &empty.Empty{})
	if err != nil {
		t.Fatalf("Failed to execute RPC: %v", err)
	}

	if resp.GetID() != 1 {
		t.Fatalf("Wrong id")
	}
}

func TestPerNodeMetadata(t *testing.T) {
	addrs, teardown := gorums.TestSetup(t, 2, func() interface{} { return initServer(t) })
	defer teardown()

	perNodeMD := func(nid uint32) metadata.MD {
		return metadata.New(map[string]string{
			"id": fmt.Sprintf("%d", nid),
		})
	}

	mgr, err := NewManager(WithNodeList(addrs), WithPerNodeMetadata(perNodeMD), WithDialTimeout(1*time.Second), WithGrpcDialOptions(
		grpc.WithBlock(),
		grpc.WithInsecure(),
	))
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	for _, node := range mgr.Nodes() {
		resp, err := node.IDFromMD(context.Background(), &empty.Empty{})
		if err != nil {
			t.Fatalf("Failed to execute RPC: %v", err)
		}

		if resp.GetID() != node.ID() {
			t.Fatalf("Wrong message")
		}
	}
}

func TestCanGetPeerInfo(t *testing.T) {
	addrs, teardown := gorums.TestSetup(t, 1, func() interface{} { return initServer(t) })
	defer teardown()

	mgr, err := NewManager(WithNodeList(addrs), WithDialTimeout(1*time.Second), WithGrpcDialOptions(
		grpc.WithBlock(),
		grpc.WithInsecure(),
	))
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	node := mgr.Nodes()[0]

	ip, err := node.WhatIP(context.Background(), &empty.Empty{})
	if err != nil {
		t.Fatalf("Failed to execute RPC: %v", err)
	}

	if ip.GetAddr() == "" {
		t.Fatalf("No data returned")
	}
}