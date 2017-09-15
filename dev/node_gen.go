// Code generated by 'gorums' plugin for protoc-gen-go. DO NOT EDIT.
// Source file to edit is: node_tmpl

package dev

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
)

// Node encapsulates the state of a node on which a remote procedure call
// can be made.
type Node struct {
	// Only assigned at creation.
	id     uint32
	self   bool
	addr   string
	conn   *grpc.ClientConn
	logger *log.Logger

	StorageClient StorageClient

	WriteAsyncClient Storage_WriteAsyncClient

	mu      sync.Mutex
	lastErr error
	latency time.Duration
}

func (n *Node) connect(opts ...grpc.DialOption) error {
	var err error
	n.conn, err = grpc.Dial(n.addr, opts...)
	if err != nil {
		return fmt.Errorf("dialing node failed: %v", err)
	}

	n.StorageClient = NewStorageClient(n.conn)

	n.WriteAsyncClient, err = n.StorageClient.WriteAsync(context.Background())
	if err != nil {
		return fmt.Errorf("stream creation failed: %v", err)
	}

	return nil
}

func (n *Node) close() error {
	_, _ = n.WriteAsyncClient.CloseAndRecv()

	if err := n.conn.Close(); err != nil {
		if n.logger != nil {
			n.logger.Printf("%d: conn close error: %v", n.id, err)
		}
		return fmt.Errorf("%d: conn close error: %v", n.id, err)
	}
	return nil
}
