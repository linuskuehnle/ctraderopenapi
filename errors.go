// Copyright 2025 Linus Kühnle
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ctraderopenapi

import (
	"github.com/linuskuehnle/ctraderopenapi/datatypes"
	"github.com/linuskuehnle/ctraderopenapi/tcp"

	"fmt"
)

// ApiHandlerNotConnectedError is returned when an operation that requires
// an active connection is attempted while the API handler is not
// connected. The CallContext field contains a short description of the
// call that produced the error (for example "SendRequest").
type ApiHandlerNotConnectedError struct {
	CallContext string
}

func (e *ApiHandlerNotConnectedError) Error() string {
	return fmt.Sprintf("%s: api handler not connected", e.CallContext)
}

// ResponseError represents a structured error returned by the server in
// response to a request. It includes an ErrorCode, optional human-readable
// Description, a maintenance end timestamp (if the server indicates a
// planned maintenance window) and an optional RetryAfter value indicating
// how many seconds the client should wait before retrying.
type ResponseError struct {
	ErrorCode               string
	Description             string
	MaintenanceEndTimestamp int64
	RetryAfter              uint64
}

func (e *ResponseError) Error() string {
	errText := e.ErrorCode
	if e.Description != "" {
		errText += fmt.Sprintf(". %s", e.Description)
	}
	if e.MaintenanceEndTimestamp != 0 {
		errText += fmt.Sprintf(". maintenance end timestamp: %d",
			e.MaintenanceEndTimestamp)
	}
	if e.RetryAfter != 0 {
		errText += fmt.Sprintf(". retry again in %d seconds",
			e.RetryAfter)
	}
	return errText
}

// GenericResponseError is used when the server returns an error that
// cannot be mapped to a more specific client-side error type. It carries
// an ErrorCode and optional Description and MaintenanceEndTimestamp to
// help callers decide how to react (retry, wait, surface to user).
type GenericResponseError struct {
	ErrorCode               string
	Description             string
	MaintenanceEndTimestamp uint64
}

func (e *GenericResponseError) Error() string {
	errText := fmt.Sprintf("generic request error: %s", e.ErrorCode)
	if e.Description != "" {
		errText += fmt.Sprintf(". %s", e.Description)
	}
	if e.MaintenanceEndTimestamp != 0 {
		errText += fmt.Sprintf(". maintenance end timestamp: %d",
			e.MaintenanceEndTimestamp)
	}
	return errText
}

// UnexpectedMessageTypeError is returned when an incoming message from
// the server has a payload type that the client code does not expect
// or cannot handle. It includes the numeric message type for
// debugging/logging.
type UnexpectedMessageTypeError struct {
	MsgType ProtoOAPayloadType
}

func (e *UnexpectedMessageTypeError) Error() string {
	return fmt.Sprintf("unexpected message type: %d", e.MsgType)
}

// ProtoUnmarshalError wraps errors produced while unmarshalling a
// protobuf payload. CallContext describes the logical context (for
// example "proto OA spot event") to aid debugging and logging.
type ProtoUnmarshalError struct {
	CallContext string
	Err         error
}

func (e *ProtoUnmarshalError) Error() string {
	errText := fmt.Sprintf("error unmarshalling %s", e.CallContext)
	if e.Err != nil {
		errText += fmt.Sprintf(": %v", e.Err)
	}
	return errText
}

// ProtoMarshalError wraps errors produced while marshalling a
// protobuf payload. CallContext describes the logical context (for
// example "request bytes") to aid debugging and logging.
type ProtoMarshalError struct {
	CallContext string
	Err         error
}

func (e *ProtoMarshalError) Error() string {
	errText := fmt.Sprintf("error marshalling %s", e.CallContext)
	if e.Err != nil {
		errText += fmt.Sprintf(": %v", e.Err)
	}
	return errText
}

/*
	Sub-package error assertions
*/
/**/

/*
./datatypes
*/
// FunctionInvalidArgError is an alias for `datatypes.FunctionInvalidArgError`.
// It is returned by library functions when a caller provides an argument
// that fails validation (for example a nil pointer where a struct is
// expected, or an argument with invalid contents). The embedded
// `FunctionName` and `Err` fields provide the callee context and the
// underlying validation error.
type FunctionInvalidArgError = datatypes.FunctionInvalidArgError

// RequestContextExpiredError is an alias for `datatypes.RequestContextExpiredError`.
// It is used by the request queue/mapper to indicate that a request's
// provided `context.Context` has been cancelled or has timed out before
// the request could be sent or responded to. The wrapped `Err` field
// contains the original context error (for example `context.Canceled`).
type RequestContextExpiredError = datatypes.RequestContextExpiredError

// IdAlreadyIncludedError is an alias for `datatypes.IdAlreadyIncludedError`.
// The event handler returns this error when attempting to register a
// callback for an event id that is already registered.
type IdAlreadyIncludedError = datatypes.IdAlreadyIncludedError

// IdNotIncludedError is an alias for `datatypes.IdNotIncludedError`.
// It is returned when trying to remove or dispatch to an event id that
// has no registered handler (unless the handler was constructed with
// `WithIgnoreIdsNotIncluded`, in which case dispatch skips silently).
type IdNotIncludedError = datatypes.IdNotIncludedError

// LifeCycleAlreadyRunningError indicates an attempt to start the internal
// lifecycle (heartbeat, channel handlers) while it is already running.
type LifeCycleAlreadyRunningError = datatypes.LifeCycleAlreadyRunningError

// LifeCycleNotRunningError indicates an operation that requires the
// lifecycle to be running was called while it was stopped.
type LifeCycleNotRunningError = datatypes.LifeCycleNotRunningError

// RequestHeapAlreadyRunningError indicates an attempt to start the
// request heap goroutine when it is already active.
type RequestHeapAlreadyRunningError = datatypes.RequestHeapAlreadyRunningError

// RequestHeapNotRunningError indicates operations on the request heap
// were attempted while the heap worker was not running.
type RequestHeapNotRunningError = datatypes.RequestHeapNotRunningError

// RequestHeapNodeNotIncludedError is returned when attempting to remove
// a request node that is not managed by the heap.
type RequestHeapNodeNotIncludedError = datatypes.RequestHeapNodeNotIncludedError

/*
./tcp
*/
// OpenConnectionError is an alias for `tcp.OpenConnectionError` and is
// returned when the underlying TCP client fails to open a connection
// (for example due to network errors or a refused connection).
type OpenConnectionError = tcp.OpenConnectionError

// CloseConnectionError is an alias for `tcp.CloseConnectionError` and
// is returned when closing the underlying connection fails.
type CloseConnectionError = tcp.CloseConnectionError

// NoConnectionError is an alias for `tcp.NoConnectionError` and is used
// to indicate that an operation that requires an open connection was
// attempted when no connection is available (for example `Send` or
// `Read` when the client is disconnected).
type NoConnectionError = tcp.NoConnectionError

// InvalidAddressError is an alias for `tcp.InvalidAddressError` and is
// returned when the configured TCP address is invalid or cannot be
// parsed.
type InvalidAddressError = tcp.InvalidAddressError

// OperationBlockedError is an alias for `tcp.OperationBlockedError`.
// It signals that an operation cannot be performed because another
// mutually-exclusive mode is active (for example calling `Read` when a
// message handler is set on the TCP client).
type OperationBlockedError = tcp.OperationBlockedError
