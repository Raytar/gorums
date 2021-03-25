package gorums

import (
	"context"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/relab/gorums/ordering"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type request struct {
	ctx  context.Context
	msg  *Message
	opts callOptions
}

type response struct {
	nid uint32
	msg protoreflect.ProtoMessage
	err error
}

type receiveQueue struct {
	msgID    uint64
	recvQ    map[uint64]chan *response
	recvQMut sync.RWMutex
}

func newReceiveQueue() *receiveQueue {
	return &receiveQueue{
		recvQ: make(map[uint64]chan *response),
	}
}

// newCall returns unique metadata for a method call.
func (m *receiveQueue) newCall(method string) (md *ordering.Metadata) {
	msgID := atomic.AddUint64(&m.msgID, 1)
	return &ordering.Metadata{
		MessageID: msgID,
		Method:    method,
	}
}

// newReply returns a channel for receiving replies
// and a done function to be called for clean up.
func (m *receiveQueue) newReply(md *ordering.Metadata, maxReplies int) (replyChan chan *response, done func()) {
	replyChan = make(chan *response, maxReplies)
	m.recvQMut.Lock()
	m.recvQ[md.MessageID] = replyChan
	m.recvQMut.Unlock()
	done = func() {
		m.recvQMut.Lock()
		delete(m.recvQ, md.MessageID)
		m.recvQMut.Unlock()
	}
	return
}

func (m *receiveQueue) putResult(id uint64, result *response) {
	m.recvQMut.RLock()
	c, ok := m.recvQ[id]
	m.recvQMut.RUnlock()
	if ok {
		c <- result
	}
}

type Channel struct {
	sendQ        chan request
	node         *Node // needed for ID and setLastError
	backoff      backoff.Config
	rand         *rand.Rand
	gorumsClient ordering.GorumsClient
	gorumsStream ordering.Gorums_NodeStreamClient
	streamMut    sync.RWMutex
	streamBroken bool
	parentCtx    context.Context
	close        context.CancelFunc
	streamCtx    context.Context
	cancelStream context.CancelFunc
}

func (s *Channel) connect() error {
	var err error
	s.parentCtx, s.close = context.WithCancel(s.node.chanCtx)
	s.streamCtx, s.cancelStream = context.WithCancel(s.parentCtx)
	s.gorumsClient = ordering.NewGorumsClient(s.node.conn)
	s.gorumsStream, err = s.gorumsClient.NodeStream(s.streamCtx)
	if err != nil {
		return err
	}
	go s.sendMsgs()
	go s.recvMsgs()
	return nil
}

func (s *Channel) Close() {
	s.close()
}

func (s *Channel) sendMsg(req request) (err error) {
	// unblock the waiting caller unless noSendWaiting is enabled
	defer func() {
		if req.opts.callType == E_Multicast || req.opts.callType == E_Unicast && !req.opts.noSendWaiting {
			s.node.putResult(req.msg.Metadata.MessageID, &response{})
		}
	}()

	// don't send if context is already cancelled.
	if req.ctx.Err() != nil {
		return req.ctx.Err()
	}

	s.streamMut.RLock()
	defer s.streamMut.RUnlock()

	c := make(chan struct{}, 1)

	// wait for either the message to be sent, or the request context being cancelled.
	// if the request context was cancelled, then we most likely have a blocked stream.
	go func() {
		select {
		case <-c:
		case <-req.ctx.Done():
			s.cancelStream()
		}
	}()

	err = s.gorumsStream.SendMsg(req.msg)
	if err != nil {
		s.node.setLastErr(err)
		s.streamBroken = true
	}
	c <- struct{}{}

	return err
}

func (s *Channel) sendMsgs() {
	var req request
	for {
		select {
		case <-s.parentCtx.Done():
			return
		case req = <-s.sendQ:
		}
		// return error if stream is broken
		if s.streamBroken {
			err := status.Errorf(codes.Unavailable, "stream is down")
			s.node.putResult(req.msg.Metadata.MessageID, &response{nid: s.node.ID(), msg: nil, err: err})
			continue
		}
		// else try to send message
		err := s.sendMsg(req)
		if err != nil {
			// return the error
			s.node.putResult(req.msg.Metadata.MessageID, &response{nid: s.node.ID(), msg: nil, err: err})
		}
	}
}

func (s *Channel) recvMsgs() {
	for {
		resp := newMessage(responseType)
		s.streamMut.RLock()
		err := s.gorumsStream.RecvMsg(resp)
		if err != nil {
			s.streamBroken = true
			s.streamMut.RUnlock()
			s.node.setLastErr(err)
			// attempt to reconnect
			s.reconnect()
		} else {
			s.streamMut.RUnlock()
			err := status.FromProto(resp.Metadata.GetStatus()).Err()
			s.node.putResult(resp.Metadata.MessageID, &response{nid: s.node.ID(), msg: resp.Message, err: err})
		}

		select {
		case <-s.parentCtx.Done():
			return
		default:
		}
	}
}

func (s *Channel) reconnect() {
	s.streamMut.Lock()
	defer s.streamMut.Unlock()

	var retries float64
	for {
		var err error

		s.streamCtx, s.cancelStream = context.WithCancel(s.parentCtx)
		s.gorumsStream, err = s.gorumsClient.NodeStream(s.streamCtx)
		if err == nil {
			s.streamBroken = false
			return
		}
		s.cancelStream()
		s.node.setLastErr(err)
		delay := float64(s.backoff.BaseDelay)
		max := float64(s.backoff.MaxDelay)
		for r := retries; delay < max && r > 0; r-- {
			delay *= s.backoff.Multiplier
		}
		delay = math.Min(delay, max)
		delay *= 1 + s.backoff.Jitter*(rand.Float64()*2-1)
		select {
		case <-time.After(time.Duration(delay)):
			retries++
		case <-s.parentCtx.Done():
			return
		}
	}
}
