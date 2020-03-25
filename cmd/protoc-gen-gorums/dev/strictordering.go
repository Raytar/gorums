package dev

import (
	"context"
	"math"
	"math/rand"
	"net"
	"sync"
	atomic "sync/atomic"
	"time"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

type strictOrderingResult struct {
	nid   uint32
	reply *anypb.Any
	err   error
}

type requestHandler func(*gorums.Message) *gorums.Message

type strictOrderingServer struct {
	handlers map[string]requestHandler
}

func newStrictOrderingServer() *strictOrderingServer {
	return &strictOrderingServer{
		handlers: make(map[string]requestHandler),
	}
}

func (s *strictOrderingServer) registerHandler(url string, handler requestHandler) {
	s.handlers[url] = handler
}

func (s *strictOrderingServer) StrictOrdering(srv gorums.Gorums_StrictOrderingServer) error {
	for {
		req, err := srv.Recv()
		if err != nil {
			return err
		}
		// handle the request if a handler is available for this rpc
		if handler, ok := s.handlers[req.GetURL()]; ok {
			resp := handler(req)
			resp.ID = req.GetID()
			err = srv.Send(resp)
			if err != nil {
				return err
			}
		}
	}
}

type strictOrderingManager struct {
	srv          *strictOrderingServer
	gorumsServer *grpc.Server
	msgID        uint64
	recvQ        map[uint64]chan *strictOrderingResult
	recvQMut     sync.RWMutex
}

func newStrictOrderingManager() *strictOrderingManager {
	return &strictOrderingManager{
		recvQ: make(map[uint64]chan *strictOrderingResult),
	}
}

func (m *strictOrderingManager) serve(listener net.Listener) {
	m.gorumsServer = grpc.NewServer()
	gorums.RegisterGorumsServer(m.gorumsServer, m.srv)
	go m.gorumsServer.Serve(listener)
}

func (m *strictOrderingManager) createStream(node *Node, backoff *backoff.Config) *strictOrderingStream {
	return &strictOrderingStream{
		node:     node,
		backoff:  backoff,
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
		sendQ:    make(chan *gorums.Message),
		recvQ:    m.recvQ,
		recvQMut: &m.recvQMut,
	}
}

func (m *strictOrderingManager) nextMsgID() uint64 {
	return atomic.AddUint64(&m.msgID, 1)
}

func (m *strictOrderingManager) stop() {
	m.gorumsServer.GracefulStop()
}

type strictOrderingStream struct {
	// needed for ID and setLastError
	node         *Node
	backoff      *backoff.Config
	rand         *rand.Rand
	gorumsClient gorums.GorumsClient
	gorumsStream gorums.Gorums_StrictOrderingClient
	streamMut    sync.RWMutex
	streamBroken bool
	sendQ        chan *gorums.Message
	recvQ        map[uint64]chan *strictOrderingResult
	recvQMut     *sync.RWMutex
}

func (s *strictOrderingStream) connect(ctx context.Context, conn *grpc.ClientConn) error {
	var err error
	s.gorumsClient = gorums.NewGorumsClient(conn)
	s.gorumsStream, err = s.gorumsClient.StrictOrdering(ctx)
	if err != nil {
		return err
	}
	go s.sendMsgs()
	go s.recvMsgs(ctx)
	return nil
}

func (s *strictOrderingStream) sendMsgs() {
	for req := range s.sendQ {
		// return error if stream is broken
		if s.streamBroken {
			err := status.Errorf(codes.Unavailable, "stream is down")
			s.recvQMut.RLock()
			if c, ok := s.recvQ[req.GetID()]; ok {
				c <- &strictOrderingResult{nid: s.node.ID(), reply: nil, err: err}
			}
			s.recvQMut.RUnlock()
			continue
		}
		// else try to send message
		s.streamMut.RLock()
		err := s.gorumsStream.SendMsg(req)
		if err == nil {
			s.streamMut.RUnlock()
			continue
		}
		s.streamBroken = true
		s.streamMut.RUnlock()
		// return the error
		s.recvQMut.RLock()
		if c, ok := s.recvQ[req.GetID()]; ok {
			c <- &strictOrderingResult{nid: s.node.ID(), reply: nil, err: err}
		}
		s.recvQMut.RUnlock()
	}
}

func (s *strictOrderingStream) recvMsgs(ctx context.Context) {
	for {
		resp := new(gorums.Message)
		s.streamMut.RLock()
		err := s.gorumsStream.RecvMsg(resp)
		if err != nil {
			s.streamBroken = true
			s.streamMut.RUnlock()
			s.node.setLastErr(err)
			// attempt to reconnect
			s.reconnectStream(ctx)
		} else {
			s.streamMut.RUnlock()
			s.recvQMut.RLock()
			if c, ok := s.recvQ[resp.GetID()]; ok {
				c <- &strictOrderingResult{nid: s.node.ID(), reply: resp.GetData(), err: nil}
			}
			s.recvQMut.RUnlock()
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (s *strictOrderingStream) reconnectStream(ctx context.Context) {
	s.streamMut.Lock()
	defer s.streamMut.Unlock()

	var retries float64
	for {
		var err error
		s.gorumsStream, err = s.gorumsClient.StrictOrdering(ctx)
		if err == nil {
			s.streamBroken = false
			return
		}
		delay := float64(s.backoff.BaseDelay)
		max := float64(s.backoff.MaxDelay)
		if retries > 0 {
			delay = math.Pow(delay, retries)
			delay = math.Min(delay, max)
		}
		delay *= 1 + s.backoff.Jitter*(rand.Float64()*2-1)
		select {
		case <-time.After(time.Duration(delay)):
			retries++
		case <-ctx.Done():
			return
		}
	}
}
