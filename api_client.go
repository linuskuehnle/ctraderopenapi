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

// APIClient is the main entry point for interacting with the cTrader
// OpenAPI. It exposes methods to connect/disconnect, send RPC-like
// requests, and subscribe to or listen for server-side events.
//
// Concurrency / lifecycle notes:
//   - Implementations are safe for concurrent use by multiple goroutines.
//   - The client maintains an internal lifecycle (started via Connect,
//     stopped via Disconnect). Certain operations assume a running
//     lifecycle (for example, sending requests or listening to events).
//
// Error handling:
//   - Fatal server-side or protocol errors are delivered via the channel
//     returned by MakeFatalErrChan(); when such an error is emitted the
//     client will be disconnected and cleaned up. the channel stays open
//     until a fatal error occurs no matter connect/disconnect state.
//
// Usage summary:
//   - Create the client with NewApiClient, call Connect, then use
//     SendRequest, SubscribeEvent/UnsubscribeEvent, or ListenToEvent as
//     required. Call Disconnect when finished.
type APIClient interface {
	/*
		Build functions
	*/
	/**/

	// MakeFatalErrChan creates a channel that will receive fatal errors from the client.
	// Errors that are sent to this channel cannot be recovered from, and the client will be
	// in a disconnected and memory-cleaned state after a fatal error is sent.
	// MakeFatalErrChan must be called again once a fatal error has occured.
	MakeFatalErrChan() (chan error, error)

	// WithConfig updates the client configuration. It must be called while
	// the client is not connected (before `Connect`) and returns the same
	// client to allow fluent construction.
	WithConfig(APIClientConfig) APIClient

	// IsConnected returns true if the client is currently connected to the server.
	IsConnected() bool
	// Connect connects to the cTrader OpenAPI server, authenticates the application and starts the keepalive routine.
	Connect() error
	// Disconnect logs the user out, stops the keepalive routine and closes the connection. It does not perform
	// memory cleanup or deallocation, the caller is responsible for disposing of the entire client.
	Disconnect() error
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
	// the library's internal event clients and any registered
	// callbacks.
	SubscribeEvent(eventData SubscribableEventData) error

	// UnsubscribeEvent removes a previously created subscription. The
	// `eventData` must match the original subscription parameters used
	// to subscribe. This call sends the corresponding unsubscribe
	// request to the server and returns any transport or validation
	// error.
	UnsubscribeEvent(eventData SubscribableEventData) error

	// ListenToEvent registers a long-running listener for listenable
	// events (server-initiated push events).
	// `eventType` selects which event to listen for.
	// `onEventCh` receives a `ListenableEvent` (concrete types live in
	// the `messages` package).
	// `ctx` controls the lifetime of the listener: when `ctx` is canceled the
	// listener is removed. If `ctx` is nil, a background context is used.
	//
	// onEventCh behavior:
	//  - The provided channel is used by the library to deliver events of the
	//    requested `eventType`. The library will close the channel when the
	//    listener's context (`ctx`) is canceled or when the client is
	//    disconnected. Callers should treat channel close as end-of-stream and
	//    not attempt to write to the channel.
	//
	// Note: callbacks that need typed event values can use `CastToEventType`
	// or the `SpawnEventHandler` helper to adapt typed functions.
	//
	// Mapping of `eventType` → concrete callback argument type:
	//  - EventType_Spots                     -> ProtoOASpotEvent
	//  - EventType_DepthQuotes               -> ProtoOADepthEvent
	//  - EventType_TrailingSLChanged        -> ProtoOATrailingSLChangedEvent
	//  - EventType_SymbolChanged            -> ProtoOASymbolChangedEvent
	//  - EventType_TraderUpdated            -> ProtoOATraderUpdatedEvent
	//  - EventType_Execution                -> ProtoOAExecutionEvent
	//  - EventType_OrderError               -> ProtoOAOrderErrorEvent
	//  - EventType_MarginChanged            -> ProtoOAMarginChangedEvent
	//  - EventType_AccountsTokenInvalidated -> ProtoOAAccountsTokenInvalidatedEvent
	//  - EventType_ClientDisconnect         -> ProtoOAClientDisconnectEvent
	//  - EventType_AccountDisconnect        -> ProtoOAAccountDisconnectEvent
	//  - EventType_MarginCallUpdate         -> ProtoOAMarginCallUpdateEvent
	//  - EventType_MarginCallTrigger        -> ProtoOAMarginCallTriggerEvent
	//
	// To register a typed callback without writing a manual wrapper, use
	// `SpawnEventHandler` which starts a small goroutine that adapts the
	// generic `ListenableEvent` channel to a typed client. Example:
	//
	//   onSpot := func(e *ProtoOASpotEvent) { fmt.Println(e) }
	//   onEventCh := make(chan ListenableEvent)
	//   if err := h.ListenToEvent(EventType_Spots, onEventCh, ctx); err != nil { /* ... */ }
	//   if err := SpawnEventHandler(ctx, onEventCh, onSpot); err != nil { /* ... */ }
	//
	// Note: the adapter performs a runtime type assertion and will panic if
	// the handler's type does not match the actual delivered event. This is a
	// caller error.
	ListenToEvent(eventType eventType, onEventCh chan ListenableEvent, ctx context.Context) error

	// ListenToClientEvent registers a long-running listener for listenable
	// api client events (connection loss / reconnect success / reconnect fail).
	// `eventType` selects which event to listen for.
	// `onEventCh` receives a `ListenToClientEvent` (concrete types live in the
	// `messages` package).
	// `ctx` controls the lifetime of the listener: when `ctx` is canceled the
	// listener is removed. If `ctx` is nil, a background context is used.
	//
	// onEventCh behavior:
	//  - The provided channel is used by the library to deliver events of the
	//    requested `eventType`. The library will close the channel when the
	//    listener's context (`ctx`) is canceled or when the client is
	//    disconnected. Callers should treat channel close as end-of-stream and
	//    not attempt to write to the channel.
	//
	// Note: callbacks that need typed event values can use `CastToClientEventType`
	// or the `SpawnClientEventHandler` helper to adapt typed functions.
	//
	// Mapping of `clientEventType` → concrete callback argument type:
	//  - ClientEventType_ConnectionLossEvent   -> ConnectionLossEvent
	//  - ClientEventType_ReconnectSuccessEvent -> ReconnectSuccessEvent
	//  - ClientEventType_ReconnectFailEvent    -> ReconnectFailEvent
	//
	// To register a typed callback without writing a manual wrapper, use
	// `SpawnClientEventHandler` which starts a small goroutine that adapts the
	// generic `ListenToClientEvent` channel to a typed client. Example:
	//
	//   onConnLoss := func(e *ConnectionLossEvent) { fmt.Println(e) }
	//   onEventCh := make(chan ListenableClientEvent)
	//   if err := h.ListenToClientEvent(ClientEventType_ConnectionLossEvent, onEventCh, ctx); err != nil { /* ... */ }
	//   if err := SpawnClientEventHandler(ctx, onEventCh, onConnLoss); err != nil { /* ... */ }
	//
	// Note: the adapter performs a runtime type assertion and will panic if
	// the handler's type does not match the actual delivered event. This is a
	// caller error.
	ListenToClientEvent(clientEventType clientEventType, onEventCh chan ListenableClientEvent, ctx context.Context) error
}

type apiClient struct {
	mu  sync.RWMutex
	cfg *APIClientConfig

	requestQueue  datatypes.RequestQueue
	requestHeap   datatypes.RequestHeap
	lifecycleData datatypes.LifecycleData
	tcpClient     tcp.TCPClient

	apiEventHandler    datatypes.EventHandler[proto.Message]
	clientEventHandler datatypes.EventHandler[ListenableClientEvent]

	cred ApplicationCredentials

	queueDataCh  chan struct{}
	tcpMessageCh chan []byte

	fatalErrCh chan error
}

// NewAPIClient creates a new API client interface instance. The returned client is not connected automatically.
func NewAPIClient(cred ApplicationCredentials, env Environment) (APIClient, error) {
	return newApiClient(cred, env)
}

func newApiClient(cred ApplicationCredentials, env Environment) (*apiClient, error) {
	if err := cred.CheckError(); err != nil {
		return nil, err
	}

	tlsConfig, err := tcp.NewSystemCertTLSConfig()
	if err != nil {
		return nil, err
	}

	h := apiClient{
		mu: sync.RWMutex{},

		requestQueue:  datatypes.NewRequestQueue(),
		requestHeap:   datatypes.NewRequestHeap(),
		lifecycleData: datatypes.NewLifecycleData(),

		apiEventHandler: datatypes.NewEventHandler[proto.Message]().
			WithIgnoreIdsNotIncluded(),
		clientEventHandler: datatypes.NewEventHandler[ListenableClientEvent]().
			WithIgnoreIdsNotIncluded(),

		cred: cred,
	}

	h.tcpClient = tcp.NewTCPClient(string(env.GetAddress())).
		WithTLS(tlsConfig).
		WithTimeout(time.Millisecond*100).
		WithReconnectDelay(time.Millisecond*500).
		WithInfiniteReconnectAttempts().
		WithReconnectCallback(h.onReconnect).
		WithConnEventHooks(h.connectionLossCallback, h.reconnectSuccessCallback, h.reconnectFailCallback)

	h.WithConfig(DefaultAPIClientConfig())

	return &h, nil
}

func (h *apiClient) MakeFatalErrChan() (chan error, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.lifecycleData.IsClientConnected() {
		return nil, &LifeCycleAlreadyRunningError{
			CallContext: "error creating fatal error channel",
		}
	}

	h.fatalErrCh = make(chan error, 1)
	return h.fatalErrCh, nil
}

func (h *apiClient) WithConfig(config APIClientConfig) APIClient {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.lifecycleData.IsClientConnected() {
		return h
	}

	h.cfg = &config

	return h
}

func (h *apiClient) onQueueData() {
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
		var noConnErr *NoConnectionError
		if !errors.As(err, &noConnErr) {
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

func (h *apiClient) onTCPMessage(msgBytes []byte) {
	if h.lifecycleData.IsClientConnected() {
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
			var reqHeapErr *RequestHeapNodeNotIncludedError
			if errors.As(err, &reqHeapErr) {
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

func (h *apiClient) onFatalError(err error) {
	if h.lifecycleData.IsClientConnected() {
		// Only lock if life cycle is running already, on Connect the mutex is already locked
		h.mu.Lock()
		defer h.mu.Unlock()
	}

	h.disconnect()

	if h.fatalErrCh == nil {
		panic(err)
	}

	h.fatalErrCh <- err
	close(h.fatalErrCh)
	h.fatalErrCh = nil
}

func (h *apiClient) onReconnect(errCh chan error) {
	defer close(errCh)

	if err := h.authenticateApp(); err != nil {
		errCh <- err
		return
	}

	errCh <- nil
}

func (h *apiClient) connectionLossCallback() {
	event := &ConnectionLossEvent{}
	h.clientEventHandler.HandleEvent(datatypes.EventId(ClientEventType_ConnectionLossEvent), event)
}

func (h *apiClient) reconnectSuccessCallback() {
	event := &ReconnectSuccessEvent{}
	h.clientEventHandler.HandleEvent(datatypes.EventId(ClientEventType_ReconnectSuccessEvent), event)
}

func (h *apiClient) reconnectFailCallback(err error) {
	event := &ReconnectFailEvent{
		Err: err,
	}
	h.clientEventHandler.HandleEvent(datatypes.EventId(ClientEventType_ReconnectFailEvent), event)
}

func (h *apiClient) IsConnected() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	isConn := h.isConnected()
	return isConn
}

func (h *apiClient) isConnected() bool {
	return h.tcpClient.HasConn()
}

func (h *apiClient) Connect() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.connect()
}

func (h *apiClient) Disconnect() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.disconnect()
}

func (h *apiClient) SendRequest(reqData RequestData) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.sendRequest(reqData)
}

func (h *apiClient) sendRequest(reqData RequestData) error {
	if !h.isConnected() {
		return &APIClientNotConnectedError{
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

func (h *apiClient) connect() error {
	// Config specific setup
	h.queueDataCh = make(chan struct{}, h.cfg.QueueBufferSize)
	h.tcpMessageCh = make(chan []byte, h.cfg.TCPMessageBufferSize)
	h.requestHeap.SetIterationInterval(h.cfg.RequestHeapIterationTimeout)

	h.requestQueue.WithDataCallbackChan(h.queueDataCh)
	h.tcpClient.WithMessageCallbackChan(h.tcpMessageCh, h.onFatalError)

	if err := h.lifecycleData.Start(); err != nil {
		return err
	}

	ctx, err := h.lifecycleData.GetContext()
	if err != nil {
		return err
	}
	onMessageSend, err := h.lifecycleData.GetOnMessageSendCh()
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

	h.lifecycleData.SetClientConnected(true)

	return nil
}

func (h *apiClient) disconnect() error {
	h.requestQueue.Clear()
	h.apiEventHandler.Clear()

	h.lifecycleData.Stop()

	h.requestHeap.Stop()

	if err := h.tcpClient.CloseConn(); err != nil {
		return err
	}

	return nil
}

func (h *apiClient) authenticateApp() error {
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

func (h *apiClient) emitHeartbeat() error {
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

	if err := h.enqueueRequest(metaData); err != nil {
		return err
	}

	if err, ok := <-errCh; ok {
		return err
	}
	return nil
}

func (h *apiClient) handleEnqueue(calleeRef string, reqMetaData *datatypes.RequestMetaData) error {
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

func (h *apiClient) enqueueRequest(reqMetaData *datatypes.RequestMetaData) error {
	if err := h.handleEnqueue("enqueueRequest", reqMetaData); err != nil {
		return err
	}

	return h.requestQueue.Enqueue(reqMetaData)
}

func (h *apiClient) handleSendPayload(reqMetaData *datatypes.RequestMetaData) error {
	// Before sending the request first check if the request context has already expired
	if err := reqMetaData.Ctx.Err(); err != nil {
		return &RequestContextExpiredError{
			Err: err,
		}
	}

	return h.sendPayload(reqMetaData.Id, reqMetaData.Req, reqMetaData.ReqType)
}

func (h *apiClient) sendPayload(reqId datatypes.RequestId, msg proto.Message, payloadType ProtoOAPayloadType) error {
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
	h.lifecycleData.SignalMessageSend()

	return h.tcpClient.Send(reqBytes)
}

func (h *apiClient) runHeartbeat(ctx context.Context, onMessageSend chan struct{}) {
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

func (h *apiClient) runQueueDataHandler(lifecycleCtx context.Context, queueDataCh <-chan struct{}) {
	for {
		select {
		case <-lifecycleCtx.Done():
			return
		case _, ok := <-queueDataCh:
			if !ok {
				return
			}
			h.onQueueData()
		}
	}
}

func (h *apiClient) runTCPMessageHandler(lifecycleCtx context.Context, tcpMessageCh <-chan []byte) {
	for {
		select {
		case <-lifecycleCtx.Done():
			return
		case msg, ok := <-tcpMessageCh:
			if !ok {
				return
			}
			h.onTCPMessage(msg)
		}
	}
}

func (h *apiClient) runFatalErrorHandler(lifecycleCtx context.Context, fatalErrCh <-chan error) {
	for {
		select {
		case <-lifecycleCtx.Done():
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
