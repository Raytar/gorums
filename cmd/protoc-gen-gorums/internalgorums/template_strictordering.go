package internalgorums

var strictOrderingVariables = `
{{$serv := serviceName .Method}}
{{$unexportMethod := unexport .Method.GoName}}
`

var strictOrderingPreamble = `
	{{- template "trace" .}}

	// get the ID which will be used to return the correct responses for a request
	msgID := {{use "atomic.AddUint64" .GenFile}}(&c.mgr.{{$unexportMethod}}ID, 1)
	in.{{msgIDField .Method}} = msgID
	
	// set up a channel to collect replies
	replies := make(chan *{{$intOut}}, c.n)
	c.mgr.{{$unexportMethod}}Lock.Lock()
	c.mgr.{{$unexportMethod}}Recv[msgID] = replies
	c.mgr.{{$unexportMethod}}Lock.Unlock()
	
	defer func() {
		// remove the replies channel when we are done
		c.mgr.{{$unexportMethod}}Lock.Lock()
		delete(c.mgr.{{$unexportMethod}}Recv, msgID)
		c.mgr.{{$unexportMethod}}Lock.Unlock()
	}()
`

var strictOrderingLoop = `
	// push the message to the nodes
	expected := c.n
	for _, n := range c.nodes {
{{- if hasPerNodeArg .Method}}
		nodeArg := f(*in, n.ID())
		if nodeArg == nil {
			expected--
			continue
		}
		nodeArg.{{msgIDField .Method}} = msgID
		n.{{$unexportMethod}}Send <- nodeArg
{{- else}}
		n.{{$unexportMethod}}Send <- in
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
			// TODO: An error from SendMsg/RecvMsg means that the stream has closed, so we probably don't need to check
			// for errors here.
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}
			{{template "traceLazyLog"}}
			replyValues = append(replyValues, r.reply)
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

var strictOrderingServerLoop = `
// {{$method}}ServerLoop is a helper function that will receive messages on srv,
// generate a response message using getResponse, and send the response back on
// srv. The function returns when the stream ends, and returns the error that
// caused it to end.
func {{$method}}ServerLoop(srv {{$serv}}_{{$method}}Server, getResponse func(*{{$in}}) *{{$out}}) error {
	ctx := srv.Context()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		req, err := srv.Recv()
		if err != nil {
			return err
		}
		resp := getResponse(req)
		resp.{{msgIDField .Method}} = req.{{msgIDField .Method}}
		err = srv.Send(resp)
		if err != nil {
			return err
		}
	}
}
`

var strictOrderingCall = commonVariables +
	quorumCallVariables +
	strictOrderingVariables +
	quorumCallComment +
	quorumCallSignature +
	strictOrderingPreamble +
	strictOrderingLoop +
	strictOrderingReply +
	strictOrderingServerLoop
