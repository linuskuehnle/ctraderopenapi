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
	"github.com/linuskuehnle/ctraderopenapi/messages"
	"github.com/linuskuehnle/ctraderopenapi/tcp"

	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

// APIHandler is the main entry point for interacting with the cTrader
// OpenAPI. It exposes methods to connect/disconnect, send RPC-like
// requests, and subscribe to or listen for server-side events.
//
// Concurrency / lifecycle notes:
//   - Implementations are safe for concurrent use by multiple goroutines.
//   - The handler maintains an internal lifecycle (started via Connect,
//     stopped via Disconnect). Certain operations assume a running
//     lifecycle (for example, sending requests or listening to events).
//
// Error handling:
//   - Fatal server-side or protocol errors are delivered via the channel
//     returned by MakeFatalErrChan(); when such an error is emitted the
//     handler will be disconnected and cleaned up. the channel stays open
//     until a fatal error occurs no matter connect/disconnect state.
//
// Usage summary:
//   - Create the handler with NewApiHandler, call Connect, then use
//     SendRequest, SubscribeEvent/UnsubscribeEvent, or ListenToEvent as
//     required. Call Disconnect when finished.
type APIHandler interface {
	/*
		Build functions
	*/
	/**/

	// MakeFatalErrChan creates a channel that will receive fatal errors from the handler.
	// Errors that are sent to this channel cannot be recovered from, and the handler will be
	// in a disconnected and memory-cleaned state after a fatal error is sent.
	// MakeFatalErrChan must be called again once a fatal error has occured.
	MakeFatalErrChan() (chan error, error)

	// WithConfig updates the handler configuration. It must be called while
	// the handler is not connected (before `Connect`) and returns the same
	// handler to allow fluent construction.
	WithConfig(APIHandlerConfig) APIHandler

	// IsConnected returns true if the handler is currently connected to the server.
	IsConnected() bool
	// Connect connects to the cTrader OpenAPI server, authenticates the application and starts the keepalive routine.
	Connect() error
	// Disconnect logs the user out, stops the keepalive routine and closes the connection. It does not perform
	// memory cleanup or deallocation, the caller is responsible for disposing of the entire handler.
	Disconnect() error
	// Reconnect performs a memory and connection cleanup followed by a connect. It can be used to recover from non-fatal connection issues.
	Reconnect() error
	// SendRequest sends a request to the server and waits for the response.
	// The response is written to the Res field of the provided RequestData struct.
	SendRequest(RequestData) error

	// SubscribeEvent subscribes the currently authenticated client for
	// the provided subscription-based event. The `eventData` argument
	// chooses the event type and provides the subscription parameters
	// (for example, account id and symbol ids for spot quotes).
	//
	// The call will send the corresponding subscribe request to the
	// server and return any transport or validation error. When the
	// server starts sending matching events, they will be dispatched to
	// the library's internal event handlers and any registered
	// callbacks.
	SubscribeEvent(eventData SubscribableEventData) error

	// UnsubscribeEvent removes a previously created subscription. The
	// `eventData` must match the original subscription parameters used
	// to subscribe. This call sends the corresponding unsubscribe
	// request to the server and returns any transport or validation
	// error.
	UnsubscribeEvent(eventData SubscribableEventData) error

	// ListenToEvent registers a long-running listener for listenable
	// events (server-initiated push events). `eventType` selects which
	// event to listen for. `onEvent` is called for each delivered
	// event; it receives a `ListenableEvent` (concrete types live in
	// the `messages` package). `ctx` controls the lifetime of the
	// listener: when `ctx` is canceled the listener is removed. If
	// `ctx` is nil, a background context is used.
	//
	// onEventCh behavior:
	//  - The provided channel is used by the library to deliver events of the
	//    requested `eventType`. The library will close the channel when the
	//    listener's context (`ctx`) is canceled or when the handler is
	//    disconnected. Callers should treat channel close as end-of-stream and
	//    not attempt to write to the channel.
	//
	// Note: callbacks that need typed event values can use `CastToEventType`
	// or the `SpawnEventHandler` helper to adapt typed functions.
	//
	// Mapping of `eventType` → concrete callback argument type:
	//  - EventType_Spots                  -> ProtoOASpotEvent
	//  - EventType_DepthQuotes            -> ProtoOADepthEvent
	//  - EventType_TrailingSLChanged     -> ProtoOATrailingSLChangedEvent
	//  - EventType_SymbolChanged         -> ProtoOASymbolChangedEvent
	//  - EventType_TraderUpdated         -> ProtoOATraderUpdatedEvent
	//  - EventType_Execution             -> ProtoOAExecutionEvent
	//  - EventType_OrderError            -> ProtoOAOrderErrorEvent
	//  - EventType_MarginChanged         -> ProtoOAMarginChangedEvent
	//  - EventType_AccountsTokenInvalidated -> ProtoOAAccountsTokenInvalidatedEvent
	//  - EventType_ClientDisconnect      -> ProtoOAClientDisconnectEvent
	//  - EventType_AccountDisconnect     -> ProtoOAAccountDisconnectEvent
	//  - EventType_MarginCallUpdate      -> ProtoOAMarginCallUpdateEvent
	//  - EventType_MarginCallTrigger     -> ProtoOAMarginCallTriggerEvent
	//
	// To register a typed callback without writing a manual wrapper, use
	// `SpawnEventHandler` which starts a small goroutine that adapts the
	// generic `ListenableEvent` channel to a typed handler. Example:
	//
	//   onSpot := func(e *messages.ProtoOASpotEvent) { fmt.Println(e) }
	//   onEventCh := make(chan ListenableEvent)
	//   if err := h.ListenToEvent(EventType_Spots, onEventCh, ctx); err != nil { /* ... */ }
	//   if err := SpawnEventHandler(ctx, onEventCh, onSpot); err != nil { /* ... */ }
	//
	// Note: the adapter performs a runtime type assertion and will panic if
	// the handler's type does not match the actual delivered event. This is a
	// caller error.
	ListenToEvent(eventType eventType, onEventCh chan ListenableEvent, ctx context.Context) error
}

type apiHandler struct {
	mu  sync.RWMutex
	cfg *APIHandlerConfig

	requestQueue  datatypes.RequestQueue
	requestHeap   datatypes.RequestHeap
	eventHandler  datatypes.EventHandler
	lifeCycleData datatypes.LifeCycleData
	tcpClient     tcp.TCPClient

	cred ApplicationCredentials

	queueDataCh  chan struct{}
	tcpMessageCh chan []byte

	fatalErrCh chan error
}

// NewApiHandler creates a new API handler interface instance. The returned handler is not connected automatically.
func NewApiHandler(cred ApplicationCredentials, env Environment) (APIHandler, error) {
	return newApiHandler(cred, env)
}

func newApiHandler(cred ApplicationCredentials, env Environment) (*apiHandler, error) {
	if err := cred.CheckError(); err != nil {
		return nil, err
	}

	tlsConfig, err := tcp.NewSystemCertTLSConfig()
	if err != nil {
		return nil, err
	}

	h := apiHandler{
		mu: sync.RWMutex{},

		requestQueue: datatypes.NewRequestQueue(),
		requestHeap:  datatypes.NewRequestHeap(),
		eventHandler: datatypes.NewEventHandler().
			WithIgnoreIdsNotIncluded(),
		lifeCycleData: datatypes.NewLifeCycleData(),
		tcpClient: tcp.NewTCPClient(string(env.GetAddress())).
			WithTLS(tlsConfig),

		cred: cred,
	}

	h.WithConfig(DefaultAPIHandlerConfig())

	return &h, nil
}

func (h *apiHandler) MakeFatalErrChan() (chan error, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.lifeCycleData.IsClientConnected() {
		return nil, &LifeCycleAlreadyRunningError{
			CallContext: "error creating fatal error channel",
		}
	}

	h.fatalErrCh = make(chan error, 1)
	return h.fatalErrCh, nil
}

func (h *apiHandler) WithConfig(config APIHandlerConfig) APIHandler {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.lifeCycleData.IsClientConnected() {
		return h
	}

	h.cfg = &config

	return h
}

func (h *apiHandler) onQueueData() {
	reqMetaData, _ := h.requestQueue.Dequeue()
	// Don't check the error of Dequeue since by the time this
	// callback goroutine is called, all previously enqueued requests might
	// have been context cancelled. Hence we check if theres non-nil returned reqMetaData
	if reqMetaData == nil {
		return
	}

	reqErrCh := reqMetaData.ErrCh

	// Check if a response is expected
	expectRes := reqMetaData.ResDataCh != nil
	if expectRes {
		// Add the request meta data to the request heap
		if err := h.requestHeap.AddNode(reqMetaData); err != nil {
			reqErrCh <- err
			close(reqErrCh)
			return
		}
	} else {
		// If no response is expected, close the error channel immediately after dispatching the request
		defer close(reqErrCh)
	}

	if err := h.handleSendPayload(reqMetaData); err != nil {
		// Check if the queue data callback has been called when the tcp client connection has just been closed
		if _, ok := err.(*NoConnectionError); !ok {
			if reqErrCh != nil {
				reqErrCh <- err

				// This check is required to avoid closing a closed channel when no response is expected
				if expectRes {
					close(reqErrCh)
				}
			}
		}
		return
	}
}

func (h *apiHandler) onTCPMessage(msgBytes []byte) {
	if h.lifeCycleData.IsClientConnected() {
		// Only lock if life cycle is running already, on Connect the mutex is already locked
		h.mu.RLock()
		defer h.mu.RUnlock()
	}

	var protoMsg messages.ProtoMessage
	if err := proto.Unmarshal(msgBytes, &protoMsg); err != nil {
		perr := &ProtoUnmarshalError{
			CallContext: "proto message",
			Err:         err,
		}
		h.fatalErrCh <- perr
		return
	}

	msgPayloadType := messages.ProtoPayloadType(protoMsg.GetPayloadType())

	switch msgPayloadType {
	case messages.ProtoPayloadType_ERROR_RES:
		var protoErrorRes messages.ProtoErrorRes
		if err := proto.Unmarshal(msgBytes, &protoErrorRes); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto error response",
				Err:         err,
			}
			h.fatalErrCh <- perr
			return
		}

		genericResErr := &GenericResponseError{
			ErrorCode:               protoErrorRes.GetErrorCode(),
			Description:             protoErrorRes.GetDescription(),
			MaintenanceEndTimestamp: protoErrorRes.GetMaintenanceEndTimestamp(),
		}

		// Since this response cannot be mapped to any request, make it a fatal error
		h.fatalErrCh <- genericResErr
		return
	case messages.ProtoPayloadType_HEARTBEAT_EVENT:
		// Ignore heartbeat events
		return
	default:
		msgOAPayloadType := ProtoOAPayloadType(msgPayloadType)
		if isListenableEvent[msgOAPayloadType] {
			if err := h.handleListenableEvent(msgOAPayloadType, &protoMsg); err != nil {
				h.fatalErrCh <- err
			}
			return
		}

		if protoMsg.ClientMsgId == nil {
			h.fatalErrCh <- errors.New("invalid proto message on response: field ClientMsgId type is nil")
			return
		}
		reqId := datatypes.RequestId(protoMsg.GetClientMsgId())

		// Remove the request meta data from the request heap
		reqMetaData, err := h.requestHeap.RemoveNode(reqId)
		if err != nil {
			if _, ok := err.(*RequestHeapNodeNotIncludedError); ok {
				return
			}
			h.fatalErrCh <- err
			return
		}

		close(reqMetaData.ErrCh)

		reqMetaData.ResDataCh <- &datatypes.ResponseData{
			ProtoMsg:    &protoMsg,
			PayloadType: msgOAPayloadType,
		}
		return
	}
}

func (h *apiHandler) onFatalError(err error) {
	if h.lifeCycleData.IsClientConnected() {
		// Only lock if life cycle is running already, on Connect the mutex is already locked
		h.mu.Lock()
		defer h.mu.Unlock()
	}

	if !h.tcpClient.HasConn() {
		// Caught concurrent function call with another fatal error that cannot be processed
		return
	}

	h.disconnect()

	if h.fatalErrCh == nil {
		panic(err)
	}

	h.fatalErrCh <- err
	close(h.fatalErrCh)
	h.fatalErrCh = nil
}

func (h *apiHandler) IsConnected() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	isConn := h.isConnected()
	return isConn
}

func (h *apiHandler) isConnected() bool {
	return h.tcpClient.HasConn()
}

func (h *apiHandler) Connect() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.connect()
}

func (h *apiHandler) Disconnect() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.disconnect()
}

func (h *apiHandler) Reconnect() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if err := h.disconnect(); err != nil {
		if _, ok := err.(*tcp.CloseConnectionError); !ok {
			return err
		} // else: ignore close connection on already closed connection error on reconnect
	}

	return h.connect()
}

func (h *apiHandler) SendRequest(reqData RequestData) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.sendRequest(reqData)
}

func (h *apiHandler) sendRequest(reqData RequestData) error {
	if !h.isConnected() {
		return &ApiHandlerNotConnectedError{
			CallContext: "SendRequest",
		}
	}

	if reqData.Ctx == nil {
		reqData.Ctx = context.Background()
	}

	expectedResType, exists := resTypeByReqType[reqData.ReqType]
	if !exists {
		return &FunctionInvalidArgError{
			FunctionName: "SendRequest",
			Err:          fmt.Errorf("provided unknown request type %d", reqData.ReqType),
		}
	}
	if expectedResType != reqData.ResType {
		return &FunctionInvalidArgError{
			FunctionName: "SendRequest",
			Err: fmt.Errorf("expected response type %d, got %d for request type %d",
				expectedResType, reqData.ResType, reqData.ReqType,
			),
		}
	}

	errCh := make(chan error)
	resDataCh := make(chan *datatypes.ResponseData)

	metaData, err := datatypes.NewRequestMetaData(&reqData, errCh, resDataCh)
	if err != nil {
		return err
	}

	if err := h.enqueueRequest(metaData); err != nil {
		return err
	}

	if err, ok := <-errCh; ok {
		return err
	}

	resData := <-resDataCh

	payload := resData.ProtoMsg.GetPayload()
	if err := checkResponseForError(payload, resData.PayloadType); err != nil {
		return err
	}

	if resData.PayloadType != reqData.ResType {
		return fmt.Errorf("unexpected response payload type: got %d, expected %d",
			resData.PayloadType, reqData.ResType,
		)
	}

	// Unmarshal payload into provided response struct
	if err := proto.Unmarshal(payload, reqData.Res); err != nil {
		return &ProtoUnmarshalError{
			CallContext: fmt.Sprintf("proto response [%d]", reqData.ResType),
			Err:         err,
		}
	}

	return nil
}

func (h *apiHandler) connect() error {
	// Config specific setup
	h.queueDataCh = make(chan struct{}, h.cfg.QueueBufferSize)
	h.tcpMessageCh = make(chan []byte, h.cfg.TCPMessageBufferSize)
	h.requestHeap.SetIterationInterval(h.cfg.RequestHeapIterationTimeout)

	h.requestQueue.WithDataCallbackChan(h.queueDataCh)
	h.tcpClient.WithMessageCallbackChan(h.tcpMessageCh, h.onFatalError)

	if err := h.lifeCycleData.Start(); err != nil {
		return err
	}

	ctx, err := h.lifeCycleData.GetContext()
	if err != nil {
		return err
	}
	onMessageSend, err := h.lifeCycleData.GetOnMessageSendCh()
	if err != nil {
		return err
	}

	// Execute heartbeat in separate goroutine
	go h.runHeartbeat(ctx, onMessageSend)

	// Start persistent goroutine that listens to queueDataCh and calls onQueueData
	go h.runQueueDataHandler(ctx, h.queueDataCh)

	// Start persistent goroutine that listens to tcpMessageCh and calls onTCPMessage
	go h.runTCPMessageHandler(ctx, h.tcpMessageCh)

	// Start persistent goroutine that listens to fatalErrCh and calls onFatalError
	go h.runFatalErrorHandler(ctx, h.fatalErrCh)

	// Start the request heap goroutine here
	if err := h.requestHeap.Start(); err != nil {
		return err
	}

	if err := h.tcpClient.OpenConn(); err != nil {
		return err
	}

	if err := h.authenticateApp(); err != nil {
		return err
	}

	h.lifeCycleData.SetClientConnected(true)

	return nil
}

func (h *apiHandler) disconnect() error {
	h.requestQueue.Clear()
	h.eventHandler.Clear()

	h.lifeCycleData.Stop()

	h.requestHeap.Stop()

	if err := h.tcpClient.CloseConn(); err != nil {
		return err
	}

	return nil
}

func (h *apiHandler) authenticateApp() error {
	reqData := RequestData{
		Ctx:     context.Background(),
		ReqType: PROTO_OA_APPLICATION_AUTH_REQ,
		Req: &messages.ProtoOAApplicationAuthReq{
			ClientId:     proto.String(h.cred.ClientId),
			ClientSecret: proto.String(h.cred.ClientSecret),
		},
		ResType: PROTO_OA_APPLICATION_AUTH_RES,
		Res:     &messages.ProtoOAApplicationAuthRes{},
	}

	return h.sendRequest(reqData)
}

func (h *apiHandler) emitHeartbeat() error {
	errCh := make(chan error)

	reqData := RequestData{
		Ctx:     context.Background(),
		ReqType: ProtoOAPayloadType(messages.ProtoPayloadType_HEARTBEAT_EVENT),
		Req:     &messages.ProtoHeartbeatEvent{},
	}

	metaData, err := datatypes.NewRequestMetaData(&reqData, errCh, nil)
	if err != nil {
		return err
	}

	if err := h.enqueueRequestFront(metaData); err != nil {
		return err
	}

	if err, ok := <-errCh; ok {
		return err
	}
	return nil
}

func (h *apiHandler) handleEnqueue(calleeRef string, reqMetaData *datatypes.RequestMetaData) error {
	if !h.tcpClient.HasConn() {
		return &EnqueueOnClosedConnError{}
	}

	if reqMetaData.ErrCh == nil {
		return &FunctionInvalidArgError{
			FunctionName: calleeRef,
			Err:          errors.New("field ErrCh of request meta data cannot be nil"),
		}
	}

	return nil
}

func (h *apiHandler) enqueueRequest(reqMetaData *datatypes.RequestMetaData) error {
	if err := h.handleEnqueue("enqueueRequest", reqMetaData); err != nil {
		return err
	}

	return h.requestQueue.Enqueue(reqMetaData)
}

func (h *apiHandler) enqueueRequestFront(reqMetaData *datatypes.RequestMetaData) error {
	if err := h.handleEnqueue("enqueueRequestFront", reqMetaData); err != nil {
		return err
	}

	return h.requestQueue.EnqueueFront(reqMetaData)
}

func (h *apiHandler) handleSendPayload(reqMetaData *datatypes.RequestMetaData) error {
	// Before sending the request first check if the request context has already expired
	if err := reqMetaData.Ctx.Err(); err != nil {
		return &RequestContextExpiredError{
			Err: err,
		}
	}

	return h.sendPayload(reqMetaData.Id, reqMetaData.Req, reqMetaData.ReqType)
}

func (h *apiHandler) sendPayload(reqId datatypes.RequestId, msg proto.Message, payloadType ProtoOAPayloadType) error {
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return &ProtoMarshalError{
			CallContext: "proto message bytes",
			Err:         err,
		}
	}

	reqIdStr := string(reqId)
	wrappedMsg := messages.ProtoMessage{
		ClientMsgId: &reqIdStr,
		PayloadType: proto.Uint32(uint32(payloadType)),
		Payload:     msgBytes,
	}

	reqBytes, err := proto.Marshal(&wrappedMsg)
	if err != nil {
		return &ProtoMarshalError{
			CallContext: "request bytes",
			Err:         err,
		}
	}

	// Signal the keepalive that a message is being sent
	h.lifeCycleData.SignalMessageSend()

	return h.tcpClient.Send(reqBytes)
}

func (h *apiHandler) runHeartbeat(ctx context.Context, onMessageSend chan struct{}) {
	timer := time.NewTimer(time.Second * Heartbeat_Timeout_Seconds)
	defer timer.Stop()

	var err error

	for {
		select {
		case <-ctx.Done():
			return
		case <-onMessageSend:
			// Reset timer on message send
			if !timer.Stop() {
				<-timer.C
			}

			timer.Reset(time.Second * Heartbeat_Timeout_Seconds)
		case <-timer.C:
			// Heartbeat timeout reached, send heartbeat
			if err = h.emitHeartbeat(); err != nil {
				h.fatalErrCh <- err
				return
			}

			timer.Reset(time.Second * Heartbeat_Timeout_Seconds)
		}
	}
}

func (h *apiHandler) runQueueDataHandler(lifeCycleCtx context.Context, queueDataCh <-chan struct{}) {
	for {
		select {
		case <-lifeCycleCtx.Done():
			return
		case _, ok := <-queueDataCh:
			if !ok {
				return
			}
			h.onQueueData()
		}
	}
}

func (h *apiHandler) runTCPMessageHandler(lifeCycleCtx context.Context, tcpMessageCh <-chan []byte) {
	for {
		select {
		case <-lifeCycleCtx.Done():
			return
		case msg, ok := <-tcpMessageCh:
			if !ok {
				return
			}
			h.onTCPMessage(msg)
		}
	}
}

func (h *apiHandler) runFatalErrorHandler(lifeCycleCtx context.Context, fatalErrCh <-chan error) {
	for {
		select {
		case <-lifeCycleCtx.Done():
			return
		case err, ok := <-fatalErrCh:
			if !ok {
				return
			}
			h.onFatalError(err)
		}
	}
}

func checkResponseForError(payloadBytes []byte, payloadType ProtoOAPayloadType) error {
	if payloadType != PROTO_OA_ERROR_RES {
		return nil
	}

	var errorMsg messages.ProtoOAErrorRes
	if err := proto.Unmarshal(payloadBytes, &errorMsg); err != nil {
		return &ProtoUnmarshalError{
			CallContext: "proto OA error response",
			Err:         err,
		}
	}

	return &ResponseError{
		ErrorCode:               errorMsg.GetErrorCode(),
		Description:             errorMsg.GetDescription(),
		MaintenanceEndTimestamp: errorMsg.GetMaintenanceEndTimestamp(),
		RetryAfter:              errorMsg.GetRetryAfter(),
	}
}
