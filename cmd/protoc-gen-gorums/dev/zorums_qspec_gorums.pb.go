// Code generated by protoc-gen-gorums. DO NOT EDIT.

package dev

import (
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// QuorumSpec is the interface of quorum functions for ZorumsService.
type QuorumSpec interface {

	// QuorumCallQF is the quorum function for the QuorumCall
	// quorum call method. The in parameter is the request object
	// supplied to the QuorumCall method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	QuorumCallQF(in *Request, replies map[uint32]*Response) (*Response, bool)

	// QuorumCallPerNodeArgQF is the quorum function for the QuorumCallPerNodeArg
	// quorum call method. The in parameter is the request object
	// supplied to the QuorumCallPerNodeArg method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	QuorumCallPerNodeArgQF(in *Request, replies map[uint32]*Response) (*Response, bool)

	// QuorumCallCustomReturnTypeQF is the quorum function for the QuorumCallCustomReturnType
	// quorum call method. The in parameter is the request object
	// supplied to the QuorumCallCustomReturnType method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	QuorumCallCustomReturnTypeQF(in *Request, replies map[uint32]*Response) (*MyResponse, bool)

	// QuorumCallComboQF is the quorum function for the QuorumCallCombo
	// quorum call method. The in parameter is the request object
	// supplied to the QuorumCallCombo method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	QuorumCallComboQF(in *Request, replies map[uint32]*Response) (*MyResponse, bool)

	// QuorumCallEmptyQF is the quorum function for the QuorumCallEmpty
	// quorum call method. The in parameter is the request object
	// supplied to the QuorumCallEmpty method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *emptypb.Empty'.
	QuorumCallEmptyQF(in *emptypb.Empty, replies map[uint32]*Response) (*Response, bool)

	// QuorumCallEmpty2QF is the quorum function for the QuorumCallEmpty2
	// quorum call method. The in parameter is the request object
	// supplied to the QuorumCallEmpty2 method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	QuorumCallEmpty2QF(in *Request, replies map[uint32]*emptypb.Empty) (*emptypb.Empty, bool)

	// QuorumCallFutureQF is the quorum function for the QuorumCallFuture
	// asynchronous quorum call method. The in parameter is the request object
	// supplied to the QuorumCallFuture method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	QuorumCallFutureQF(in *Request, replies map[uint32]*Response) (*Response, bool)

	// QuorumCallFuturePerNodeArgQF is the quorum function for the QuorumCallFuturePerNodeArg
	// asynchronous quorum call method. The in parameter is the request object
	// supplied to the QuorumCallFuturePerNodeArg method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	QuorumCallFuturePerNodeArgQF(in *Request, replies map[uint32]*Response) (*Response, bool)

	// QuorumCallFutureCustomReturnTypeQF is the quorum function for the QuorumCallFutureCustomReturnType
	// asynchronous quorum call method. The in parameter is the request object
	// supplied to the QuorumCallFutureCustomReturnType method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	QuorumCallFutureCustomReturnTypeQF(in *Request, replies map[uint32]*Response) (*MyResponse, bool)

	// QuorumCallFutureComboQF is the quorum function for the QuorumCallFutureCombo
	// asynchronous quorum call method. The in parameter is the request object
	// supplied to the QuorumCallFutureCombo method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	QuorumCallFutureComboQF(in *Request, replies map[uint32]*Response) (*MyResponse, bool)

	// QuorumCallFuture2QF is the quorum function for the QuorumCallFuture2
	// asynchronous quorum call method. The in parameter is the request object
	// supplied to the QuorumCallFuture2 method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	QuorumCallFuture2QF(in *Request, replies map[uint32]*Response) (*Response, bool)

	// QuorumCallFutureEmptyQF is the quorum function for the QuorumCallFutureEmpty
	// asynchronous quorum call method. The in parameter is the request object
	// supplied to the QuorumCallFutureEmpty method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	QuorumCallFutureEmptyQF(in *Request, replies map[uint32]*emptypb.Empty) (*emptypb.Empty, bool)

	// QuorumCallFutureEmpty2QF is the quorum function for the QuorumCallFutureEmpty2
	// asynchronous quorum call method. The in parameter is the request object
	// supplied to the QuorumCallFutureEmpty2 method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *emptypb.Empty'.
	QuorumCallFutureEmpty2QF(in *emptypb.Empty, replies map[uint32]*Response) (*Response, bool)

	// CorrectableQF is the quorum function for the Correctable
	// correctable quorum call method. The in parameter is the request object
	// supplied to the Correctable method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	CorrectableQF(in *Request, replies map[uint32]*Response) (*Response, int, bool)

	// CorrectablePerNodeArgQF is the quorum function for the CorrectablePerNodeArg
	// correctable quorum call method. The in parameter is the request object
	// supplied to the CorrectablePerNodeArg method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	CorrectablePerNodeArgQF(in *Request, replies map[uint32]*Response) (*Response, int, bool)

	// CorrectableCustomReturnTypeQF is the quorum function for the CorrectableCustomReturnType
	// correctable quorum call method. The in parameter is the request object
	// supplied to the CorrectableCustomReturnType method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	CorrectableCustomReturnTypeQF(in *Request, replies map[uint32]*Response) (*MyResponse, int, bool)

	// CorrectableComboQF is the quorum function for the CorrectableCombo
	// correctable quorum call method. The in parameter is the request object
	// supplied to the CorrectableCombo method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	CorrectableComboQF(in *Request, replies map[uint32]*Response) (*MyResponse, int, bool)

	// CorrectableEmptyQF is the quorum function for the CorrectableEmpty
	// correctable quorum call method. The in parameter is the request object
	// supplied to the CorrectableEmpty method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	CorrectableEmptyQF(in *Request, replies map[uint32]*emptypb.Empty) (*emptypb.Empty, int, bool)

	// CorrectableEmpty2QF is the quorum function for the CorrectableEmpty2
	// correctable quorum call method. The in parameter is the request object
	// supplied to the CorrectableEmpty2 method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *emptypb.Empty'.
	CorrectableEmpty2QF(in *emptypb.Empty, replies map[uint32]*Response) (*Response, int, bool)

	// CorrectableStreamQF is the quorum function for the CorrectableStream
	// correctable stream quorum call method. The in parameter is the request object
	// supplied to the CorrectableStream method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	CorrectableStreamQF(in *Request, replies map[uint32]*Response) (*Response, int, bool)

	// CorrectableStreamPerNodeArgQF is the quorum function for the CorrectableStreamPerNodeArg
	// correctable stream quorum call method. The in parameter is the request object
	// supplied to the CorrectableStreamPerNodeArg method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	CorrectableStreamPerNodeArgQF(in *Request, replies map[uint32]*Response) (*Response, int, bool)

	// CorrectableStreamCustomReturnTypeQF is the quorum function for the CorrectableStreamCustomReturnType
	// correctable stream quorum call method. The in parameter is the request object
	// supplied to the CorrectableStreamCustomReturnType method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	CorrectableStreamCustomReturnTypeQF(in *Request, replies map[uint32]*Response) (*MyResponse, int, bool)

	// CorrectableStreamComboQF is the quorum function for the CorrectableStreamCombo
	// correctable stream quorum call method. The in parameter is the request object
	// supplied to the CorrectableStreamCombo method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	CorrectableStreamComboQF(in *Request, replies map[uint32]*Response) (*MyResponse, int, bool)

	// CorrectableStreamEmptyQF is the quorum function for the CorrectableStreamEmpty
	// correctable stream quorum call method. The in parameter is the request object
	// supplied to the CorrectableStreamEmpty method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Request'.
	CorrectableStreamEmptyQF(in *Request, replies map[uint32]*emptypb.Empty) (*emptypb.Empty, int, bool)

	// CorrectableStreamEmpty2QF is the quorum function for the CorrectableStreamEmpty2
	// correctable stream quorum call method. The in parameter is the request object
	// supplied to the CorrectableStreamEmpty2 method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *emptypb.Empty'.
	CorrectableStreamEmpty2QF(in *emptypb.Empty, replies map[uint32]*Response) (*Response, int, bool)
}
