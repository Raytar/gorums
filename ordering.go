package gorums

import (
	"context"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/relab/gorums/ordering"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type MethodInfo struct {
	RequestType  protoreflect.Message
	ResponseType protoreflect.Message
}

type orderingResult struct {
	nid   uint32
	reply protoreflect.ProtoMessage
	err   error
}

type receiveQueue struct {
	msgID    uint64
	recvQ    map[uint64]chan *orderingResult
	recvQMut sync.RWMutex
}

func newReceiveQueue() *receiveQueue {
	return &receiveQueue{
		recvQ: make(map[uint64]chan *orderingResult),
	}
}

func (m *receiveQueue) nextMsgID() uint64 {
	return atomic.AddUint64(&m.msgID, 1)
}

func (m *receiveQueue) putChan(id uint64, c chan *orderingResult) {
	m.recvQMut.Lock()
	m.recvQ[id] = c
	m.recvQMut.Unlock()
}

func (m *receiveQueue) deleteChan(id uint64) {
	m.recvQMut.Lock()
	delete(m.recvQ, id)
	m.recvQMut.Unlock()
}

func (m *receiveQueue) putResult(id uint64, result *orderingResult) {
	m.recvQMut.RLock()
	c, ok := m.recvQ[id]
	m.recvQMut.RUnlock()
	if ok {
		c <- result
	}
}

type orderedNodeStream struct {
	*receiveQueue
	sendQ        chan *Message
	node         *Node // needed for ID and setLastError
	backoff      backoff.Config
	rand         *rand.Rand
	gorumsClient ordering.GorumsClient
	gorumsStream ordering.Gorums_NodeStreamClient
	streamMut    sync.RWMutex
	streamBroken uint32
}

func (s *orderedNodeStream) connectOrderedStream(ctx context.Context, conn *grpc.ClientConn) error {
	var err error
	s.gorumsClient = ordering.NewGorumsClient(conn)
	s.gorumsStream, err = s.gorumsClient.NodeStream(ctx)
	if err != nil {
		return err
	}
	go s.sendMsgs(ctx)
	go s.recvMsgs(ctx)
	return nil
}

func (s *orderedNodeStream) sendMsgs(ctx context.Context) {
	var req *Message
	for {
		select {
		case <-ctx.Done():
			return
		case req = <-s.sendQ:
		}
		// return error if stream is broken
		if atomic.LoadUint32(&s.streamBroken) == 1 {
			err := status.Errorf(codes.Unavailable, "stream is down")
			s.putResult(req.Metadata.MessageID, &orderingResult{nid: s.node.ID(), reply: nil, err: err})
			continue
		}
		// else try to send message
		s.streamMut.RLock()
		err := s.gorumsStream.SendMsg(req)
		if err == nil {
			s.streamMut.RUnlock()
			continue
		}
		atomic.StoreUint32(&s.streamBroken, 1)
		s.streamMut.RUnlock()
		s.node.setLastErr(err)
		// return the error
		s.putResult(req.Metadata.MessageID, &orderingResult{nid: s.node.ID(), reply: nil, err: err})
	}
}

func (s *orderedNodeStream) recvMsgs(ctx context.Context) {
	for {
		resp := newGorumsMessage(gorumsResponse)
		s.streamMut.RLock()
		err := s.gorumsStream.RecvMsg(resp)
		if err != nil {
			atomic.StoreUint32(&s.streamBroken, 1)
			s.streamMut.RUnlock()
			s.node.setLastErr(err)
			// attempt to reconnect
			s.reconnectStream(ctx)
		} else {
			s.streamMut.RUnlock()
			err := status.FromProto(resp.Metadata.GetStatus()).Err()
			s.putResult(resp.Metadata.MessageID, &orderingResult{nid: s.node.ID(), reply: resp.Message, err: err})
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (s *orderedNodeStream) reconnectStream(ctx context.Context) {
	s.streamMut.Lock()
	defer s.streamMut.Unlock()

	var retries float64
	for {
		var err error
		s.gorumsStream, err = s.gorumsClient.NodeStream(ctx)
		if err == nil {
			atomic.StoreUint32(&s.streamBroken, 0)
			return
		}
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
		case <-ctx.Done():
			return
		}
	}
}