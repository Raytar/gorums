package dev

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
)

// Node encapsulates the state of a node on which a remote procedure call
// can be made.
type Node struct {
	// Only assigned at creation.
	id   uint32
	self bool
	addr string
	conn *grpc.ClientConn

	writeAsyncClient Register_WriteAsyncClient

	sync.Mutex
	lastErr error
	latency time.Duration
}

func (n *Node) connect(opts ...grpc.DialOption) error {
	var err error
	n.conn, err = grpc.Dial(n.addr, opts...)
	if err != nil {
		return fmt.Errorf("dialing node failed: %v", err)
	}
	client := NewRegisterClient(n.conn)
	n.writeAsyncClient, err = client.WriteAsync(context.Background())
	if err != nil {
		return fmt.Errorf("stream creation failed: %v", err)
	}
	return nil
}

func (n *Node) close() error {
	var err error
	_, err = n.writeAsyncClient.CloseAndRecv()
	err2 := n.conn.Close()
	if err != nil {
		return fmt.Errorf("stream close failed: %v", err)
	} else if err2 != nil {
		return fmt.Errorf("conn close failed: %v", err2)
	}
	return nil
}
