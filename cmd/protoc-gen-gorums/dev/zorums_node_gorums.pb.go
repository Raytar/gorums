// Code generated by protoc-gen-gorums. DO NOT EDIT.

package dev

import (
	context "context"
	fmt "fmt"
	io "io"
	sync "sync"
)

type nodeServices struct {
	ReaderServiceClient
	readMulticastClient  ReaderService_ReadMulticastClient
	readMulticast2Client ReaderService_ReadMulticast2Client
	readOrderedClient    ReaderService_ReadOrderedClient
}

func (n *Node) connectStream() (err error) {
	n.ReaderServiceClient = NewReaderServiceClient(n.conn)
	n.readMulticastClient, err = n.ReaderServiceClient.ReadMulticast(context.Background())
	if err != nil {
		return fmt.Errorf("stream creation failed: %v", err)
	}
	n.readMulticast2Client, err = n.ReaderServiceClient.ReadMulticast2(context.Background())
	if err != nil {
		return fmt.Errorf("stream creation failed: %v", err)
	}
	n.readOrderedClient, err = n.ReaderServiceClient.ReadOrdered(context.Background())
	if err != nil {
		return fmt.Errorf("stream creation failed: %v", err)
	}
	go n.readOrderedSendMsgs()
	go n.readOrderedRecvMsgs()
	return nil
}

func (n *Node) closeStream() (err error) {
	_, err = n.readMulticastClient.CloseAndRecv()
	_, err = n.readMulticast2Client.CloseAndRecv()
	err = n.readOrderedClient.CloseSend()
	close(n.readOrderedSend)
	return err
}

func (n *Node) readOrderedSendMsgs() {
	for msg := range n.readOrderedSend {
		err := n.readOrderedClient.SendMsg(msg)
		if err != nil {
			if err != io.EOF {
				n.setLastErr(err)
			}
			return
		}
	}
}

func (n *Node) readOrderedRecvMsgs() {
	for {
		msg := new(ReadResponse)
		err := n.readOrderedClient.RecvMsg(msg)
		if err != nil {
			if err != io.EOF {
				n.setLastErr(err)
			}
			return
		}
		id := msg.MsgID
		n.readOrderedLock.RLock()
		if c, ok := n.readOrderedRecv[id]; ok {
			c <- &internalReadResponse{n.id, msg, nil}
		}
		n.readOrderedLock.RUnlock()
	}
}

type nodeData struct {
	readOrderedSend chan *ReadRequest
	readOrderedRecv map[uint64]chan *internalReadResponse
	readOrderedLock *sync.RWMutex
}
