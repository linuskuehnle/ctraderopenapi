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
	"github.com/linuskuehnle/ctraderopenapi/internal/datatypes"
	"github.com/linuskuehnle/ctraderopenapi/internal/messages"

	"fmt"
	"time"
)

type apiClientConfig struct {
	requestTimeout              time.Duration
	queueBufferSize             int
	tcpMessageBufferSize        int
	requestHeapIterationTimeout time.Duration
	concurrentEventEmits        bool
}

func defaultAPIClientConfig() apiClientConfig {
	return apiClientConfig{
		requestTimeout:              ConfigDefault_RequestTimeout,
		queueBufferSize:             ConfigDefault_QueueBufferSize,
		tcpMessageBufferSize:        ConfigDefault_TCPMessageBufferSize,
		requestHeapIterationTimeout: ConfigDefault_RequestHeapIterationTimeout,
		concurrentEventEmits:        ConfigDefault_ConcurrentEventEmits,
	}
}

// Used for referencing requests to the rate limit type (live/historical).
type rateLimitType int

// endpointAddress is the type used for endpoint address constants.
type endpointAddress string

// Environment represents the cTrader OpenAPI environment to connect to.
type Environment int

// GetAddress returns the endpoint address for the referenced environment.
func (env Environment) GetAddress() endpointAddress {
	switch env {
	case Environment_Demo:
		return EndpointAddress_Demo
	case Environment_Live:
		return EndpointAddress_Live
	default:
		return ""
	}
}

// CtraderAccountId represents a cTrader account identifier.
// It is a thin typed alias over int64 to make call sites explicit.
type CtraderAccountId = datatypes.CtraderAccountId

// AccessToken is the token that is being used to access a set of account ids.
// It is a thin typed alias over string to make call sites explicit.
type AccessToken = datatypes.AccessToken

// RefreshToken is the token that is being used to update the AccessToken
// when it expires.
// It is a thin typed alias over string to make call sites explicit.
type RefreshToken = datatypes.RefreshToken

// ApplicationCredentials holds the client credentials required to
// authenticate the application with the cTrader OpenAPI proxy.
// Both fields are required.
type ApplicationCredentials struct {
	ClientId     string
	ClientSecret string
}

// CheckError validates the ApplicationCredentials, returning a
// descriptive error when either the client id or client secret are
// missing.
func (c ApplicationCredentials) CheckError() error {
	if c.ClientId == "" {
		return fmt.Errorf("client ID must not be empty")
	}
	if c.ClientSecret == "" {
		return fmt.Errorf("client secret must not be empty")
	}
	return nil
}

// RequestData is the argument struct for SendRequest.
//   - Ctx: Request context. If context.Err() is not nil, the request response will not be awaited or if the
//     request is still in queue it will be removed from queue. Either way SendRequest will abort and return context.Err().
//   - ReqType: Payload type of the request message.
//   - Req: Request message. Must be a pointer to a protobuf message struct.
//   - ResType: Expected payload type of the response message.
//   - Res: Pointer to an empty protobuf message struct where the response will be unmarshalled into.
type RequestData = datatypes.RequestData

// Generic types for specifying underlying event types.

type eventType = datatypes.EventType
type apiEventType = eventType
type clientEventType = eventType

// Type for matching api events to listener channels via EventListener.
type apiEventKey = messages.APIEventKey

// BaseEvent is the interface that all events must implement.
type BaseEvent = datatypes.BaseEvent

// APIEvent marks protobuf API event message types that can be listened to
// (push-style events).
//
// These event types are delivered by the server independently from any
// single client's request; users can register listener channels of type
// `APIEvent` using `ListenToAPIEvent`.
type APIEvent = datatypes.APIEvent

// ClientEvent marks events emitted by the API client itself
// (for example connection loss, reconnect success, etc).
//
// These event types are emitted internally by the client based on the
// client behaviour and its connection status; users can register listener
// channels of type `ClientEvent` using `ListenToClientEvent`.
type ClientEvent = datatypes.ClientEvent

// DistributableEvent marks protobuf API event message types that can be
// listened to (push-style events) AND can be used in an EventListener
// instance.
//
// Distributable means you can add a specific event listener channel via
// ListenToAPIEvent and provide key data. That way each event is matched
// to the key data, and if it matches the event will only be sent to the
// specifically provided event listener channel for that key data.
type DistributableEvent = datatypes.DistributableEvent

// CastToAPIEventType attempts to cast a generic `APIEvent` to the
// concrete typed event `T`. It returns the typed value and `true` if the
// assertion succeeded; otherwise it returns the zero value and `false`.
//
// This helper is commonly used by small adapter goroutines that accept a
// generic `APIEvent` channel but want to call a typed handler
// (for example converting `APIEvent` to `*ProtoOASpotEvent`).
func CastToAPIEventType[T APIEvent](event APIEvent) (T, bool) {
	t, ok := event.(T)
	return t, ok
}

// CastToClientEventType attempts to cast a generic `ClientEvent` to the
// concrete typed event `T`. It returns the typed value and `true` if the
// assertion succeeded; otherwise it returns the zero value and `false`.
//
// This helper is commonly used by small adapter goroutines that accept a
// generic `ClientEvent` channel but want to call a typed handler
// (for example converting `ClientEvent` to `*ReconnectSuccessEvent`).
func CastToClientEventType[T ClientEvent](event ClientEvent) (T, bool) {
	t, ok := event.(T)
	return t, ok
}

// APIEventListenData describes the parameters for listening to API events.
//   - `EventType` selects which event to listen for.
//   - `EventCh` receives an `APIEvent` that will be delivered for each matching event.
//   - `EventKeyData` provides subscription-specific data for filtering (used by event listener).
//
// Mapping of `EventType` → concrete callback argument type:
//   - APIEventType_Spots                    -> ProtoOASpotEvent
//   - APIEventType_DepthQuotes              -> ProtoOADepthEvent
//   - APIEventType_TrailingSLChanged        -> ProtoOATrailingSLChangedEvent
//   - APIEventType_SymbolChanged            -> ProtoOASymbolChangedEvent
//   - APIEventType_TraderUpdated            -> ProtoOATraderUpdatedEvent
//   - APIEventType_Execution                -> ProtoOAExecutionEvent
//   - APIEventType_OrderError               -> ProtoOAOrderErrorEvent
//   - APIEventType_MarginChanged            -> ProtoOAMarginChangedEvent
//   - APIEventType_AccountsTokenInvalidated -> ProtoOAAccountsTokenInvalidatedEvent
//   - APIEventType_ClientDisconnect         -> ProtoOAClientDisconnectEvent
//   - APIEventType_AccountDisconnect        -> ProtoOAAccountDisconnectEvent
//   - APIEventType_MarginCallUpdate         -> ProtoOAMarginCallUpdateEvent
//   - APIEventType_MarginCallTrigger        -> ProtoOAMarginCallTriggerEvent
//
// EventCh behaviour:
//   - The provided channel is used by the library to deliver events of the
//     requested `EventType`. The library will close the channel when the
//     listener's context (`ctx`) is canceled or when the client is
//     disconnected. Callers should treat channel close as end-of-stream and
//     not attempt to write to the channel.
//
// EventKeyData usage:
//   - When nil, all events of the specified `EventType` are delivered to `EventCh`.
//   - When set to a concrete KeyData* type (e.g., KeyDataSpotEvent), only events matching
//     the key data's criteria are delivered. The event handler uses BuildKey() to generate
//     a key from the KeyData fields and routes only matching events to this listener.
//   - Example: setting EventKeyData to KeyDataSpotEvent{Ctid: 12345, SymbolId: 1} will
//     deliver only spot events for that specific account and symbol.
//   - This enables fine-grained event filtering without the need for listener-side filtering.
type APIEventListenData struct {
	EventType    apiEventType
	EventCh      chan APIEvent
	EventKeyData apiEventKeyData
}

// Checks if APIEventListenData is correctly configured.
//
// If ok, returns nil.
func (d *APIEventListenData) CheckError() error {
	if d.EventType == 0 {
		return fmt.Errorf("field EventType must not be 0")
	}
	if d.EventCh == nil {
		return fmt.Errorf("field EventCh must not be nil")
	}
	return nil
}

// ClientEventListenData describes the parameters for listening to client events.
//   - `EventType` selects which event to listen for.
//   - `EventCh` receives a `ClientEvent` that will be delivered for each matching event.
//
// Mapping of `EventType` → concrete callback argument type:
//   - ClientEventType_FatalErrorEvent       -> FatalErrorEvent
//   - ClientEventType_ConnectionLossEvent   -> ConnectionLossEvent
//   - ClientEventType_ReconnectSuccessEvent -> ReconnectSuccessEvent
//   - ClientEventType_ReconnectFailEvent    -> ReconnectFailEvent
//
// EventCh behaviuor:
//   - The provided channel is used by the library to deliver events of the
//     requested `eventType`. The library will close the channel when the
//     listener's context (`ctx`) is canceled or when the client is
//     disconnected. Callers should treat channel close as end-of-stream and
//     not attempt to write to the channel.
type ClientEventListenData struct {
	EventType clientEventType
	EventCh   chan ClientEvent
}

// Checks if ClientEventListenData is correctly configured.
//
// If ok, returns nil.
func (d *ClientEventListenData) CheckError() error {
	if d.EventType == 0 {
		return fmt.Errorf("field EventType must not be 0")
	}
	if d.EventCh == nil {
		return fmt.Errorf("field EventCh must not be nil")
	}
	return nil
}

// SubscriptionData describes data required to subscribe or unsubscribe
// to a subscription-based event (for example, account id and symbol
// ids for the Spots event). Implementations validate their fields via
// `CheckError` and are passed to `SubscribeAPIEvent` / `UnsubscribeAPIEvent`.
type SubscriptionData interface {
	// CheckError validates the subscription payload and returns a
	// non-nil error when the payload is invalid.
	CheckError() error
}

// APIEventSubData groups the event type and subscription-specific
// data for `SubscribeAPIEvent` and `UnsubscribeAPIEvent` calls.
//
//   - EventType selects which server-side event to subscribe/unsubscribe.
//   - SubcriptionData is a concrete struct implementing `SubscriptionData`
//     that provides the required parameters for that event (for example,
//     account id and symbol ids for Spot events).
type APIEventSubData struct {
	EventType       apiEventType
	SubcriptionData SubscriptionData
}

/*
	Subscription based events
*/
/**/

// Subscription data for the Spot Event (ProtoOASpotEvent).
type SpotEventData struct {
	CtraderAccountId CtraderAccountId
	SymbolIds        []int64
}

// CheckError checks if the provided subscription data is valid. Returns an error if it is invalid.
func (d *SpotEventData) CheckError() error {
	if err := d.CtraderAccountId.CheckError(); err != nil {
		return fmt.Errorf("field CtraderAccountId invalid: %w", err)
	}

	if len(d.SymbolIds) == 0 {
		return fmt.Errorf("field SymbolIds mustn't be empty")
	}

	return nil
}

// Subscription data for the Live Trendbars Event.
type LiveTrendbarEventData struct {
	CtraderAccountId CtraderAccountId
	SymbolId         int64
	Period           ProtoOATrendbarPeriod
}

// CheckError checks if the provided subscription data is valid. Returns an error if it is invalid.
func (d *LiveTrendbarEventData) CheckError() error {
	if err := d.CtraderAccountId.CheckError(); err != nil {
		return fmt.Errorf("field CtraderAccountId invalid: %w", err)
	}

	if d.SymbolId == 0 {
		return fmt.Errorf("field SymbolIds mustn't be 0")
	}
	if d.Period == 0 {
		return fmt.Errorf("field Period mustn't be 0")
	}

	return nil
}

// Subscription data for the Depth Quote Event (ProtoOADepthEvent).
type DepthQuoteEventData struct {
	CtraderAccountId CtraderAccountId
	SymbolIds        []int64
}

// CheckError checks if the provided subscription data is valid. Returns an error if it is invalid.
func (d *DepthQuoteEventData) CheckError() error {
	if err := d.CtraderAccountId.CheckError(); err != nil {
		return fmt.Errorf("field CtraderAccountId invalid: %w", err)
	}

	if len(d.SymbolIds) == 0 {
		return fmt.Errorf("field SymbolIds mustn't be empty")
	}

	return nil
}

/*
API Client event types
*/

// FatalErrorEvent is emitted when a fatal error occurs in the API client.
// The `Err` field contains the underlying error.
//
// Listeners may use this event to trigger application-level shutdown
// or cleanup logic.
type FatalErrorEvent struct {
	Err error
}

func (e *FatalErrorEvent) IsListenableClientEvent() {}
func (*FatalErrorEvent) ToDistributableEvent() DistributableEvent {
	return nil
}

// ConnectionLossEvent is emitted when the TCP connection to the
// cTrader OpenAPI server is unexpectedly lost. Listeners may use this
// to trigger local reconnection logic or cleanup.
type ConnectionLossEvent struct{}

func (e *ConnectionLossEvent) IsListenableClientEvent() {}
func (*ConnectionLossEvent) ToDistributableEvent() DistributableEvent {
	return nil
}

// ReconnectSuccessEvent is emitted after a successful reconnect and
// re-authentication cycle. It indicates the client is operational
// again.
type ReconnectSuccessEvent struct{}

func (e *ReconnectSuccessEvent) IsListenableClientEvent() {}
func (*ReconnectSuccessEvent) ToDistributableEvent() DistributableEvent {
	return nil
}

// ReconnectFailEvent is emitted when the client's reconnect loop fails
// repeatedly and gives up (or when an error occurs during reconnect).
// The `Err` field contains the underlying error.
type ReconnectFailEvent struct {
	Err error
}

func (e *ReconnectFailEvent) IsListenableClientEvent() {}
func (*ReconnectFailEvent) ToDistributableEvent() DistributableEvent {
	return nil
}

/*
Key Data generation section
*/

// Interface for api event key data structures. They must contain
// a BuildKey function to return an api event key based on the field
// values of the key data instance.
type apiEventKeyData = datatypes.APIEventKeyData

// KeyDataSpotEvent is the type for ProtoOASpotEvent key generation.
// Used in ListenToAPIEvent.
type KeyDataSpotEvent struct {
	Ctid     CtraderAccountId
	SymbolId int64
	Period   ProtoOATrendbarPeriod
}

func (k *KeyDataSpotEvent) BuildKey() apiEventKey {
	key := fmt.Sprintf("SpotEvent__Ctid:%d__SymbolId:%d__TrendbarPeriod:%d",
		k.Ctid,
		k.SymbolId,
		k.Period,
	)
	return apiEventKey(key)
}

// KeyDataDepthEvent is the type for ProtoOADepthEvent key generation.
// Used in ListenToAPIEvent.
type KeyDataDepthEvent struct {
	Ctid     CtraderAccountId
	SymbolId int64
}

func (k *KeyDataDepthEvent) BuildKey() apiEventKey {
	key := fmt.Sprintf("DepthEvent__Ctid:%d__SymbolId:%d",
		k.Ctid,
		k.SymbolId,
	)
	return apiEventKey(key)
}

// KeyDataTrailingSLChangedEvent is the type for ProtoOATrailingSLChangedEvent key generation.
// Used in ListenToAPIEvent.
type KeyDataTrailingSLChangedEvent struct {
	Ctid       CtraderAccountId
	PositionId int64
	OrderId    int64
}

func (k *KeyDataTrailingSLChangedEvent) BuildKey() apiEventKey {
	key := fmt.Sprintf("TrailingSLChangedEvent__Ctid:%d__PositionId:%d__OrderId:%d",
		k.Ctid,
		k.PositionId,
		k.OrderId,
	)
	return apiEventKey(key)
}

// KeyDataSymbolChangedEvent is the type for ProtoOASymbolChangedEvent key generation.
// Used in ListenToAPIEvent.
type KeyDataSymbolChangedEvent struct {
	Ctid CtraderAccountId
}

func (k *KeyDataSymbolChangedEvent) BuildKey() apiEventKey {
	key := fmt.Sprintf("SymbolChangedEvent__Ctid:%d",
		k.Ctid,
	)
	return apiEventKey(key)
}

// KeyDataTraderUpdatedEvent is the type for ProtoOATraderUpdatedEvent key generation.
// Used in ListenToAPIEvent.
type KeyDataTraderUpdatedEvent struct {
	Ctid CtraderAccountId
}

func (k *KeyDataTraderUpdatedEvent) BuildKey() apiEventKey {
	key := fmt.Sprintf("TraderUpdatedEvent__Ctid:%d",
		k.Ctid,
	)
	return apiEventKey(key)
}

// KeyDataExecutionEvent is the type for ProtoOAExecutionEvent key generation.
// Used in ListenToAPIEvent.
type KeyDataExecutionEvent struct {
	Ctid       CtraderAccountId
	PositionId int64
	OrderId    int64
	DealId     int64
}

func (k *KeyDataExecutionEvent) BuildKey() apiEventKey {
	key := fmt.Sprintf("ExecutionEvent__Ctid:%d__PositionId:%d__OrderId:%d__DealId:%d",
		k.Ctid,
		k.PositionId,
		k.OrderId,
		k.DealId,
	)
	return apiEventKey(key)
}

// KeyDataOrderErrorEvent is the type for ProtoOAOrderErrorEvent key generation.
// Used in ListenToAPIEvent.
type KeyDataOrderErrorEvent struct {
	Ctid       CtraderAccountId
	PositionId int64
	OrderId    int64
}

func (k *KeyDataOrderErrorEvent) BuildKey() apiEventKey {
	key := fmt.Sprintf("OrderErrorEvent__Ctid:%d__PositionId:%d__OrderId:%d",
		k.Ctid,
		k.PositionId,
		k.OrderId,
	)
	return apiEventKey(key)
}

// KeyDataMarginChangedEvent is the type for ProtoOAMarginChangedEvent key generation.
// Used in ListenToAPIEvent.
type KeyDataMarginChangedEvent struct {
	Ctid       CtraderAccountId
	PositionId int64
}

func (k *KeyDataMarginChangedEvent) BuildKey() apiEventKey {
	key := fmt.Sprintf("MarginChangedEvent__Ctid:%d__PositionId:%d",
		k.Ctid,
		k.PositionId,
	)
	return apiEventKey(key)
}

// KeyDataAccountDisconnectEvent is the type for ProtoOAAccountDisconnectEvent key generation.
// Used in ListenToAPIEvent.
type KeyDataAccountDisconnectEvent struct {
	Ctid CtraderAccountId
}

func (k *KeyDataAccountDisconnectEvent) BuildKey() apiEventKey {
	key := fmt.Sprintf("AccountDisconnectEvent__Ctid:%d",
		k.Ctid,
	)
	return apiEventKey(key)
}

// KeyDataMarginCallUpdateEvent is the type for ProtoOAMarginCallUpdateEvent key generation.
// Used in ListenToAPIEvent.
type KeyDataMarginCallUpdateEvent struct {
	Ctid CtraderAccountId
}

func (k *KeyDataMarginCallUpdateEvent) BuildKey() apiEventKey {
	key := fmt.Sprintf("MarginCallUpdateEvent__Ctid:%d",
		k.Ctid,
	)
	return apiEventKey(key)
}

// KeyDataMarginCallTriggerEvent is the type for ProtoOAMarginCallTriggerEvent key generation.
// Used in ListenToAPIEvent.
type KeyDataMarginCallTriggerEvent struct {
	Ctid CtraderAccountId
}

func (k *KeyDataMarginCallTriggerEvent) BuildKey() apiEventKey {
	key := fmt.Sprintf("MarginCallTriggerEvent__Ctid:%d",
		k.Ctid,
	)
	return apiEventKey(key)
}
