package gorums

import (
	"context"

	"github.com/relab/gorums/ordering"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type QuorumCallData struct {
	Manager        *Manager
	Nodes          []*Node
	Message        protoreflect.ProtoMessage
	MethodID       int32
	PerNodeArgFn   func(protoreflect.ProtoMessage, uint32) protoreflect.ProtoMessage
	QuorumFunction func(protoreflect.ProtoMessage, map[uint32]protoreflect.ProtoMessage) (protoreflect.ProtoMessage, bool)
}

func QuorumCall(ctx context.Context, d QuorumCallData) (resp protoreflect.ProtoMessage, err error) {
	msgID := d.Manager.nextMsgID()
	// set up channel to collect replies to this call.
	replyChan := make(chan *orderingResult, len(d.Nodes))
	d.Manager.putChan(msgID, replyChan)
	// and remove it when the call it scomplete
	defer d.Manager.deleteChan(msgID)

	md := &ordering.Metadata{
		MessageID: msgID,
		MethodID:  d.MethodID,
	}

	expected := len(d.Nodes)
	msgs := make([]*Message, len(d.Nodes))
	for i, n := range d.Nodes {
		msg := d.Message
		if d.PerNodeArgFn != nil {
			msg = d.PerNodeArgFn(d.Message, n.id)
			if msg == nil {
				expected--
				continue
			}
		}
		msgs[i] = &Message{Metadata: md, Message: msg}
	}

	ctx, cancel := context.WithCancel(ctx)
	resultChan := make(chan *orderingResult)

	nodeConfig(d.Nodes).sendMsgs(ctx, msgs)

	go func() {
		defer cancel()
		var (
			errs    []GRPCError
			quorum  bool
			replies = make(map[uint32]protoreflect.ProtoMessage)
		)

		for {
			select {
			case r := <-replyChan:
				if r.err != nil {
					errs = append(errs, GRPCError{r.nid, r.err})
					break
				}
				reply := r.reply
				replies[r.nid] = reply
				if resp, quorum = d.QuorumFunction(d.Message, replies); quorum {
					resultChan <- &orderingResult{reply: resp, err: nil}
					return
				}
			case <-ctx.Done():
				resultChan <- &orderingResult{reply: resp, err: QuorumCallError{"incomplete call", len(replies), errs}}
				return
			}
			if len(errs)+len(replies) == expected {
				resultChan <- &orderingResult{reply: resp, err: QuorumCallError{"incomplete call", len(replies), errs}}
				return
			}
		}
	}()

	result := <-resultChan
	return result.reply, result.err
}
