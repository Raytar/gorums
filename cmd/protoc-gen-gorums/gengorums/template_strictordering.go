package gengorums

var strictOrderingVariables = `
{{$serv := serviceName .Method}}
{{$unexportMethod := unexport .Method.GoName}}
{{$marshalAny := use "ptypes.MarshalAny" .GenFile}}
{{$unmarshalAny := use "ptypes.UnmarshalAny" .GenFile}}
{{$errorf := use "fmt.Errorf" .GenFile}}
{{$gorumsMsg := use "gorums.Message" .GenFile}}
`

var strictOrderingPreamble = `
	{{- template "trace" .}}

	// get the ID which will be used to return the correct responses for a request
	msgID := c.mgr.nextMsgID()
	
	// set up a channel to collect replies
	replies := make(chan *strictOrderingResult, c.n)
	c.mgr.recvQMut.Lock()
	c.mgr.recvQ[msgID] = replies
	c.mgr.recvQMut.Unlock()
	
	defer func() {
		// remove the replies channel when we are done
		c.mgr.recvQMut.Lock()
		delete(c.mgr.recvQ, msgID)
		c.mgr.recvQMut.Unlock()
	}()
`

var strictOrderingLoop = `
{{if not (hasPerNodeArg .Method) -}}
	data, err := {{$marshalAny}}(in)
	if err != nil {
		return nil, {{$errorf}}("failed to marshal message: %w", err)
	}
	msg := &{{$gorumsMsg}}{
		ID: msgID,
		URL: "{{fullName .Method}}",
		Data: data,
	}
{{end -}}

	// push the message to the nodes
	expected := c.n
	for _, n := range c.nodes {
{{- if hasPerNodeArg .Method}}
		nodeArg := f(*in, n.ID())
		if nodeArg == nil {
			expected--
			continue
		}
		data, err := {{$marshalAny}}(nodeArg)
		if err != nil {
			return nil, {{$errorf}}("failed to marshal message: %w", err)
		}
		msg := &{{$gorumsMsg}}{
			ID: msgID,
			URL: "{{fullName .Method}}",
			Data: data,
		}
		nodeArg.{{msgIDField .Method}} = msgID
		n.strictOrdering.sendQ <- msg
{{- else}}
		n.strictOrdering.sendQ <- msg
{{- end}}
	}
`

var strictOrderingReply = `
	var (
		replyValues = make([]*{{$out}}, 0, expected)
		errs []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replies:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}
			{{template "traceLazyLog"}}
			reply := new({{$out}})
			err := {{$unmarshalAny}}(r.reply, reply)
			if err != nil {
				errs = append(errs, GRPCError{r.nid, {{$errorf}}("failed to unmarshal reply: %w", err)})
				break
			}
			replyValues = append(replyValues, reply)
			if resp, quorum = c.qspec.{{$method}}QF({{withQFArg .Method "in, "}}replyValues); quorum {
				return resp, nil
			}
		case <-ctx.Done():
			return resp, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
		}

		if len(errs)+len(replyValues) == expected {
			return resp, QuorumCallError{"incomplete call", len(replyValues), errs}
		}
	}
}
`

var strictOrderingHandler = `
// {{$method}}Handler is the server API for the {{$method}} rpc.
type {{$method}}Handler interface {
	{{$method}}(*{{$in}}) (*{{$out}})
}

func (m *strictOrderingManager) Register{{$method}}Handler(handler {{$method}}Handler) {
	m.srv.registerHandler("{{fullName .Method}}", func(in *{{$gorumsMsg}}) *{{$gorumsMsg}} {
		req := new({{$in}})
		err := {{$unmarshalAny}}(in.GetData(), req)
		// TODO: how to handle marshaling errors here
		if err != nil {
			return new({{$gorumsMsg}})
		}
		resp := handler.{{$method}}(req)
		data, err := {{$marshalAny}}(resp)
		if err != nil {
			return new({{$gorumsMsg}})
		}
		return &{{$gorumsMsg}}{Data: data, URL: in.GetURL()}
	})
}
`

var strictOrderingQuorum = commonVariables +
	quorumCallVariables +
	strictOrderingVariables +
	quorumCallComment +
	quorumCallSignature +
	strictOrderingPreamble +
	strictOrderingLoop +
	strictOrderingReply +
	strictOrderingHandler
