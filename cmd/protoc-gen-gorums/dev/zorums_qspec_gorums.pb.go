// Code generated by protoc-gen-gorums. DO NOT EDIT.

package dev

import (
	empty "github.com/golang/protobuf/ptypes/empty"
)

// QuorumSpec is the interface of quorum functions for ZorumsService.
type QuorumSpec interface {

	// QuorumCallQF is the quorum function for the QuorumCall
	// quorum call method.
	QuorumCallQF(replies []*Response) (*Response, bool)

	// QuorumCallPerNodeArgQF is the quorum function for the QuorumCallPerNodeArg
	// quorum call method.
	QuorumCallPerNodeArgQF(replies []*Response) (*Response, bool)

	// QuorumCallQFWithRequestArgQF is the quorum function for the QuorumCallQFWithRequestArg
	// quorum call method.
	QuorumCallQFWithRequestArgQF(in *Request, replies []*Response) (*Response, bool)

	// QuorumCallCustomReturnTypeQF is the quorum function for the QuorumCallCustomReturnType
	// quorum call method.
	QuorumCallCustomReturnTypeQF(replies []*Response) (*MyResponse, bool)

	// QuorumCallComboQF is the quorum function for the QuorumCallCombo
	// quorum call method.
	QuorumCallComboQF(in *Request, replies []*Response) (*MyResponse, bool)

	// QuorumCallEmptyQF is the quorum function for the QuorumCallEmpty
	// quorum call method.
	QuorumCallEmptyQF(replies []*Response) (*Response, bool)

	// QuorumCallEmpty2QF is the quorum function for the QuorumCallEmpty2
	// quorum call method.
	QuorumCallEmpty2QF(replies []*empty.Empty) (*empty.Empty, bool)

	// QuorumCallFutureQF is the quorum function for the QuorumCallFuture
	// asynchronous quorum call method.
	QuorumCallFutureQF(replies []*Response) (*Response, bool)

	// QuorumCallFuturePerNodeArgQF is the quorum function for the QuorumCallFuturePerNodeArg
	// asynchronous quorum call method.
	QuorumCallFuturePerNodeArgQF(replies []*Response) (*Response, bool)

	// QuorumCallFutureQFWithRequestArgQF is the quorum function for the QuorumCallFutureQFWithRequestArg
	// asynchronous quorum call method.
	QuorumCallFutureQFWithRequestArgQF(in *Request, replies []*Response) (*Response, bool)

	// QuorumCallFutureCustomReturnTypeQF is the quorum function for the QuorumCallFutureCustomReturnType
	// asynchronous quorum call method.
	QuorumCallFutureCustomReturnTypeQF(replies []*Response) (*MyResponse, bool)

	// QuorumCallFutureComboQF is the quorum function for the QuorumCallFutureCombo
	// asynchronous quorum call method.
	QuorumCallFutureComboQF(in *Request, replies []*Response) (*MyResponse, bool)

	// QuorumCallFuture2QF is the quorum function for the QuorumCallFuture2
	// asynchronous quorum call method.
	QuorumCallFuture2QF(replies []*Response) (*Response, bool)

	// QuorumCallFutureEmptyQF is the quorum function for the QuorumCallFutureEmpty
	// asynchronous quorum call method.
	QuorumCallFutureEmptyQF(replies []*empty.Empty) (*empty.Empty, bool)

	// QuorumCallFutureEmpty2QF is the quorum function for the QuorumCallFutureEmpty2
	// asynchronous quorum call method.
	QuorumCallFutureEmpty2QF(replies []*Response) (*Response, bool)

	// CorrectableQF is the quorum function for the Correctable
	// correctable quorum call method.
	CorrectableQF(replies []*Response) (*Response, int, bool)

	// CorrectablePerNodeArgQF is the quorum function for the CorrectablePerNodeArg
	// correctable quorum call method.
	CorrectablePerNodeArgQF(replies []*Response) (*Response, int, bool)

	// CorrectableQFWithRequestArgQF is the quorum function for the CorrectableQFWithRequestArg
	// correctable quorum call method.
	CorrectableQFWithRequestArgQF(in *Request, replies []*Response) (*Response, int, bool)

	// CorrectableCustomReturnTypeQF is the quorum function for the CorrectableCustomReturnType
	// correctable quorum call method.
	CorrectableCustomReturnTypeQF(replies []*Response) (*MyResponse, int, bool)

	// CorrectableComboQF is the quorum function for the CorrectableCombo
	// correctable quorum call method.
	CorrectableComboQF(in *Request, replies []*Response) (*MyResponse, int, bool)

	// CorrectableEmptyQF is the quorum function for the CorrectableEmpty
	// correctable quorum call method.
	CorrectableEmptyQF(replies []*empty.Empty) (*empty.Empty, int, bool)

	// CorrectableEmpty2QF is the quorum function for the CorrectableEmpty2
	// correctable quorum call method.
	CorrectableEmpty2QF(replies []*Response) (*Response, int, bool)

	// CorrectableStreamQF is the quorum function for the CorrectableStream
	// correctable stream quorum call method.
	CorrectableStreamQF(replies []*Response) (*Response, int, bool)

	// CorrectableStreamPerNodeArgQF is the quorum function for the CorrectableStreamPerNodeArg
	// correctable stream quorum call method.
	CorrectableStreamPerNodeArgQF(replies []*Response) (*Response, int, bool)

	// CorrectableStreamQFWithRequestArgQF is the quorum function for the CorrectableStreamQFWithRequestArg
	// correctable stream quorum call method.
	CorrectableStreamQFWithRequestArgQF(in *Request, replies []*Response) (*Response, int, bool)

	// CorrectableStreamCustomReturnTypeQF is the quorum function for the CorrectableStreamCustomReturnType
	// correctable stream quorum call method.
	CorrectableStreamCustomReturnTypeQF(replies []*Response) (*MyResponse, int, bool)

	// CorrectableStreamComboQF is the quorum function for the CorrectableStreamCombo
	// correctable stream quorum call method.
	CorrectableStreamComboQF(in *Request, replies []*Response) (*MyResponse, int, bool)

	// CorrectableStreamEmptyQF is the quorum function for the CorrectableStreamEmpty
	// correctable stream quorum call method.
	CorrectableStreamEmptyQF(replies []*empty.Empty) (*empty.Empty, int, bool)

	// CorrectableStreamEmpty2QF is the quorum function for the CorrectableStreamEmpty2
	// correctable stream quorum call method.
	CorrectableStreamEmpty2QF(replies []*Response) (*Response, int, bool)

	// StrictOrderingQCQF is the quorum function for the StrictOrderingQC
	// ordered quorum call method.
	StrictOrderingQCQF(replies []*Response) (*Response, bool)

	// StrictOrderingPerNodeArgQF is the quorum function for the StrictOrderingPerNodeArg
	// ordered quorum call method.
	StrictOrderingPerNodeArgQF(replies []*Response) (*Response, bool)

	// StrictOrderingQFWithReqQF is the quorum function for the StrictOrderingQFWithReq
	// ordered quorum call method.
	StrictOrderingQFWithReqQF(in *Request, replies []*Response) (*Response, bool)

	// StrictOrderingCustomReturnTypeQF is the quorum function for the StrictOrderingCustomReturnType
	// ordered quorum call method.
	StrictOrderingCustomReturnTypeQF(replies []*Response) (*MyResponse, bool)

	// StrictOrderingCombiQF is the quorum function for the StrictOrderingCombi
	// ordered quorum call method.
	StrictOrderingCombiQF(in *Request, replies []*Response) (*MyResponse, bool)
}
