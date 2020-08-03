package gorums

import (
	"context"

	"github.com/relab/gorums/ordering"
)

func Unicast(ctx context.Context, d CallData) {
	msgID := d.Manager.nextMsgID()

	md := &ordering.Metadata{
		MessageID: msgID,
		MethodID:  d.MethodID,
	}

	select {
	case d.Node.sendQ <- &Message{Metadata: md, Message: d.Message}:
	case <-ctx.Done():
	}
}
