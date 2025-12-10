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
//   - Fatal server-side or protocol errors are delivered via the
//     FatalErrorEvent client event; when such an error is emitted the
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

	// WithConfig updates the client configuration. It must be called while
	// the client is not connected (before `Connect`) and returns the same
	// client to allow fluent construction. To use the default configuration,
	// this function must not be called as the api client constructor already
	// applies the default config.
	WithConfig(APIClientConfig) APIClient

	// IsConnected does not reflect the actual connection state of the underlying
	// TCP client, but whether the API client has been connected via `Connect`
	// and not yet disconnected via `Disconnect`. The actual connection state can be
	// determined by listening to client events (ConnectionLossEvent,
	// ReconnectSuccessEvent, ReconnectFailEvent).
	IsConnected() bool
	// Connect connects to the cTrader OpenAPI server, authenticates the application and
	// starts the keepalive routine.
	Connect() error
	// Disconnect stops the keepalive routine and closes the connection.
	// Disconnect is a cleanup function. Event listeners are removed and any resources
	// held by the client are released. After calling Disconnect, the client
	// must not be used again unless Connect is called again.
	Disconnect()
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
	// listeners.
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
	// `onEventCh` receives a `ListenableEvent`.
	// `ctx` controls the lifetime of the listener: when `ctx` is canceled the
	// listener is removed. If `ctx` is nil, the listener will not be canceled
	// until the client is disconnected.
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
	ListenToEvent(eventType eventType, onEventCh chan ListenableEvent, ctx context.Context) error

	// ListenToClientEvent registers a long-running listener for listenable
	// api client events (connection loss / reconnect success / reconnect fail).
	// `eventType` selects which event to listen for.
	// `onEventCh` receives a `ListenToClientEvent`.
	// `ctx` controls the lifetime of the listener: when `ctx` is canceled the
	// listener is removed. If `ctx` is nil, the listener will not be canceled
	// until the client is disconnected.
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
	//  - ClientEventType_FatalErrorEvent       -> FatalErrorEvent
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

	activeSubscriptions map[eventType]SubscriptionData

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

	c := apiClient{
		mu: sync.RWMutex{},

		requestQueue:  datatypes.NewRequestQueue(),
		requestHeap:   datatypes.NewRequestHeap(),
		lifecycleData: datatypes.NewLifecycleData(),

		apiEventHandler:    datatypes.NewEventHandler[proto.Message](),
		clientEventHandler: datatypes.NewEventHandler[ListenableClientEvent](),

		cred: cred,

		activeSubscriptions: make(map[eventType]SubscriptionData),

		fatalErrCh: make(chan error),
	}

	c.tcpClient = tcp.NewTCPClient(string(env.GetAddress())).
		WithTLS(tlsConfig).
		WithTimeout(time.Millisecond*100).
		WithReconnectTimeout(time.Second*3).
		WithInfiniteReconnectAttempts().
		WithConnEventHooks(c.onConnectionLoss, c.onReconnectSuccess, c.onReconnectFail)

	c.WithConfig(DefaultAPIClientConfig())

	return &c, nil
}

func (c *apiClient) WithConfig(config APIClientConfig) APIClient {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lifecycleData.IsClientConnected() {
		return c
	}

	c.cfg = &config

	return c
}

func (c *apiClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	isConn := c.isConnected()
	return isConn
}

func (c *apiClient) isConnected() bool {
	return c.tcpClient.HasConn()
}

func (c *apiClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.connect()
}

func (c *apiClient) Disconnect() {
	c.mu.Lock()
	c.disconnect()
	c.mu.Unlock()
}

func (c *apiClient) SendRequest(reqData RequestData) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.sendRequest(reqData)
}

func (c *apiClient) sendRequest(reqData RequestData) error {
	if !c.isConnected() {
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

	if err := c.enqueueRequest(metaData); err != nil {
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

func (c *apiClient) connect() error {
	// Config specific setup
	c.queueDataCh = make(chan struct{}, c.cfg.QueueBufferSize)
	c.tcpMessageCh = make(chan []byte, c.cfg.TCPMessageBufferSize)
	c.requestHeap.SetIterationInterval(c.cfg.RequestHeapIterationTimeout)

	c.requestQueue.WithDataCallbackChan(c.queueDataCh)
	c.tcpClient.WithMessageCallbackChan(c.tcpMessageCh, c.onFatalError)

	if err := c.lifecycleData.Start(); err != nil {
		return err
	}

	ctx, err := c.lifecycleData.GetContext()
	if err != nil {
		return err
	}
	onMessageSend, err := c.lifecycleData.GetOnMessageSendCh()
	if err != nil {
		return err
	}

	// Execute heartbeat in separate goroutine
	go c.runHeartbeat(ctx, onMessageSend)

	// Start persistent goroutine that listens to queueDataCh and calls onQueueData
	go c.runQueueDataHandler(ctx, c.queueDataCh)

	// Start persistent goroutine that listens to tcpMessageCh and calls onTCPMessage
	go c.runTCPMessageHandler(ctx, c.tcpMessageCh)

	// Start persistent goroutine that listens to fatalErrCh and calls onFatalError
	go c.runFatalErrorHandler(ctx, c.fatalErrCh)

	// Start the request heap goroutine here
	if err := c.requestHeap.Start(); err != nil {
		return err
	}

	if err := c.tcpClient.OpenConn(); err != nil {
		return err
	}

	if err := c.authenticateApp(); err != nil {
		return err
	}

	c.lifecycleData.SetClientConnected(true)

	return nil
}

func (c *apiClient) disconnect() {
	c.requestQueue.Clear()

	c.lifecycleData.Stop()

	c.requestHeap.Stop()

	c.tcpClient.CloseConn()

	c.apiEventHandler.Clear()
	c.clientEventHandler.Clear()
	c.activeSubscriptions = make(map[eventType]SubscriptionData)
}

func (c *apiClient) authenticateApp() error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), appAuthenticationRequestTimeout)
	defer cancelCtx()

	reqData := RequestData{
		Ctx:     ctx,
		ReqType: PROTO_OA_APPLICATION_AUTH_REQ,
		Req: &messages.ProtoOAApplicationAuthReq{
			ClientId:     proto.String(c.cred.ClientId),
			ClientSecret: proto.String(c.cred.ClientSecret),
		},
		ResType: PROTO_OA_APPLICATION_AUTH_RES,
		Res:     &messages.ProtoOAApplicationAuthRes{},
	}

	return c.sendRequest(reqData)
}

func (c *apiClient) emitHeartbeat() error {
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

	if err := c.enqueueRequest(metaData); err != nil {
		return err
	}

	if err, ok := <-errCh; ok {
		return err
	}
	return nil
}

func (c *apiClient) handleEnqueue(calleeRef string, reqMetaData *datatypes.RequestMetaData) error {
	if !c.tcpClient.HasConn() {
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

func (c *apiClient) enqueueRequest(reqMetaData *datatypes.RequestMetaData) error {
	if err := c.handleEnqueue("enqueueRequest", reqMetaData); err != nil {
		return err
	}

	return c.requestQueue.Enqueue(reqMetaData)
}

func (c *apiClient) handleSendPayload(reqMetaData *datatypes.RequestMetaData) error {
	// Before sending the request first check if the request context has already expired
	if err := reqMetaData.Ctx.Err(); err != nil {
		return &RequestContextExpiredError{
			Err: err,
		}
	}

	return c.sendPayload(reqMetaData.Id, reqMetaData.Req, reqMetaData.ReqType)
}

func (c *apiClient) sendPayload(reqId datatypes.RequestId, msg proto.Message, payloadType ProtoOAPayloadType) error {
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
	c.lifecycleData.SignalMessageSend()

	return c.tcpClient.Send(reqBytes)
}

func (c *apiClient) runHeartbeat(ctx context.Context, onMessageSend chan struct{}) {
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
			if err = c.emitHeartbeat(); err != nil {
				/*
					Ignore EnqueueOnClosedConnError and NoConnectionError
				*/
				// Check if EnqueueOnClosedConnError occured right when the tcp client connection
				// has just been closed and before the lifecycle data has been updated to reflect
				// the disconnected state.
				var enqueueErr *EnqueueOnClosedConnError
				// Check if the queue data callback for heartbeat emit has been called
				// right when the tcp client connection has just been closed and before
				// the lifecycle data has been updated to reflect the disconnected state.
				var noConnErr *NoConnectionError

				if errors.As(err, &enqueueErr) || errors.As(err, &noConnErr) {
					timer.Reset(time.Second * Heartbeat_Timeout_Seconds)
					continue
				}

				c.fatalErrCh <- err
				return
			}

			timer.Reset(time.Second * Heartbeat_Timeout_Seconds)
		}
	}
}

func (c *apiClient) runQueueDataHandler(lifecycleCtx context.Context, queueDataCh <-chan struct{}) {
	for {
		select {
		case <-lifecycleCtx.Done():
			return
		case _, ok := <-queueDataCh:
			if !ok {
				return
			}
			c.onQueueData()
		}
	}
}

func (c *apiClient) runTCPMessageHandler(lifecycleCtx context.Context, tcpMessageCh <-chan []byte) {
	for {
		select {
		case <-lifecycleCtx.Done():
			return
		case msg, ok := <-tcpMessageCh:
			if !ok {
				return
			}
			c.onTCPMessage(msg)
		}
	}
}

func (c *apiClient) runFatalErrorHandler(lifecycleCtx context.Context, fatalErrCh <-chan error) {
	for {
		select {
		case <-lifecycleCtx.Done():
			return
		case err, ok := <-fatalErrCh:
			if !ok {
				return
			}
			c.onFatalError(err)
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
