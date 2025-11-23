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

	"fmt"
	"time"
)

// APIClientConfig configures runtime buffer sizes and timeouts used by
// the API client. These values provide reasonable defaults but can be
// adjusted with `WithConfig` before calling `Connect`.
//
// Fields:
//   - QueueBufferSize: number of queued requests that may be buffered by
//     the internal request queue before backpressure applies.
//   - TCPMessageBufferSize: size of the channel used to receive inbound
//     TCP messages from the network reader.
//   - RequestHeapIterationTimeout: interval used by the request heap to
//     periodically check for expired request contexts.
type APIClientConfig struct {
	QueueBufferSize             int
	TCPMessageBufferSize        int
	RequestHeapIterationTimeout time.Duration
}

// DefaultAPIClientConfig returns a configuration with conservative,
// production-suitable defaults. Callers may modify the returned struct
// and pass it to `WithConfig` on the API client before connecting.
func DefaultAPIClientConfig() APIClientConfig {
	return APIClientConfig{
		QueueBufferSize:             DefaultQueueBufferSize,
		TCPMessageBufferSize:        DefaultTCPMessageBufferSize,
		RequestHeapIterationTimeout: DefaultRequestHeapIterationTimeout,
	}
}

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
type CtraderAccountId int64

// CheckError validates the account id. It returns a non-nil error when
// the id is invalid (zero). Additional validation rules can be added
// later if required.
func (id CtraderAccountId) CheckError() error {
	if id == 0 {
		return fmt.Errorf("cTrader account ID must not be empty")
	}

	// Further validation can be added here if needed (like length or format checks)

	return nil
}

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

// Event types for ListenToEvent and ListenToClientEvent
// Use these typed aliases to make intent explicit in public APIs.
type eventType datatypes.EventId
type clientEventType datatypes.EventId

// RequestData is the argument struct for SendRequest.
//   - Ctx: Request context. If context.Err() is not nil, the request response will not be awaited or if the
//     request is still in queue it will be removed from queue. Either way SendRequest will abort and return context.Err().
//   - ReqType: Payload type of the request message.
//   - Req: Request message. Must be a pointer to a protobuf message struct.
//   - ResType: Expected payload type of the response message.
//   - Res: Pointer to an empty protobuf message struct where the response will be unmarshalled into.
type RequestData = datatypes.RequestData

// ListenableEvent marks protobuf message types that can be listened to
// (push-style events). Implemented by generated protobuf message types
// in the `messages` package with a zero-sized method.
//
// These event types are delivered by the server independently from any
// single client's request; users can register callbacks using
// `ListenToEvent` which will receive a `ListenableEvent` value when a
// matching event occurs.
type ListenableEvent interface {
	IsListenableEvent()
}

type ListenableClientEvent interface {
	IsListenableClientEvent()
}

// CastToEventType attempts to cast a generic `ListenableEvent` to the
// concrete typed event `T`. It returns the typed value and `true` if the
// assertion succeeded; otherwise it returns the zero value and `false`.
//
// This helper is commonly used by small adapter goroutines that accept a
// generic `ListenableEvent` channel but want to call a typed handler
// (for example converting `ListenableEvent` to `*ProtoOASpotEvent`).
func CastToEventType[T ListenableEvent](event ListenableEvent) (T, bool) {
	t, ok := event.(T)
	return t, ok
}

// CastToClientEventType attempts to cast a generic `ListenableClientEvent` to the
// concrete typed event `T`. It returns the typed value and `true` if the
// assertion succeeded; otherwise it returns the zero value and `false`.
//
// This helper is commonly used by small adapter goroutines that accept a
// generic `ListenableClientEvent` channel but want to call a typed handler
// (for example converting `ListenableClientEvent` to `*ReconnectSuccessEvent`).
func CastToClientEventType[T ListenableClientEvent](event ListenableClientEvent) (T, bool) {
	t, ok := event.(T)
	return t, ok
}

// SubscriptionData describes data required to subscribe or unsubscribe
// to a subscription-based event (for example, account id and symbol
// ids for the Spots event). Implementations validate their fields via
// `CheckError` and are passed to `SubscribeEvent` / `UnsubscribeEvent`.
type SubscriptionData interface {
	// CheckError validates the subscription payload and returns a
	// non-nil error when the payload is invalid.
	CheckError() error
}

// SubscribableEventData groups the event type and subscription-specific
// data for `SubscribeEvent` and `UnsubscribeEvent` calls.
//
//   - EventType selects which server-side event to subscribe/unsubscribe.
//   - SubcriptionData is a concrete struct implementing `SubscriptionData`
//     that provides the required parameters for that event (for example,
//     account id and symbol ids for Spot events).
type SubscribableEventData struct {
	EventType       eventType
	SubcriptionData SubscriptionData
}

/*
	Subscription based events
*/
/**/

// Subscription data for the Spot Event (ProtoOASpotEvent).
type SubscriptionDataSpotEvent struct {
	CtraderAccountId CtraderAccountId
	SymbolIds        []int64
}

// CheckError checks if the provided subscription data is valid. Returns an error if it is invalid.
func (d *SubscriptionDataSpotEvent) CheckError() error {
	if err := d.CtraderAccountId.CheckError(); err != nil {
		return fmt.Errorf("field CtraderAccountId invalid: %w", err)
	}

	if len(d.SymbolIds) == 0 {
		return fmt.Errorf("field SymbolIds mustn't be empty")
	}

	return nil
}

// Subscription data for the Live Trendbars Event.
type SubscriptionDataLiveTrendbarEvent struct {
	CtraderAccountId CtraderAccountId
	SymbolId         int64
	Period           ProtoOATrendbarPeriod
}

// CheckError checks if the provided subscription data is valid. Returns an error if it is invalid.
func (d *SubscriptionDataLiveTrendbarEvent) CheckError() error {
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
type SubscriptionDataDepthQuoteEvent struct {
	CtraderAccountId CtraderAccountId
	SymbolIds        []int64
}

// CheckError checks if the provided subscription data is valid. Returns an error if it is invalid.
func (d *SubscriptionDataDepthQuoteEvent) CheckError() error {
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

// ConnectionLossEvent is emitted when the TCP connection to the
// cTrader OpenAPI server is unexpectedly lost. Listeners may use this
// to trigger local reconnection logic or cleanup.
type ConnectionLossEvent struct{}

func (e *ConnectionLossEvent) IsListenableClientEvent() {}

// ReconnectSuccessEvent is emitted after a successful reconnect and
// re-authentication cycle. It indicates the client is operational
// again.
type ReconnectSuccessEvent struct{}

func (e *ReconnectSuccessEvent) IsListenableClientEvent() {}

// ReconnectFailEvent is emitted when the client's reconnect loop fails
// repeatedly and gives up (or when an error occurs during reconnect).
// The `Err` field contains the underlying error.
type ReconnectFailEvent struct {
	Err error
}

func (e *ReconnectFailEvent) IsListenableClientEvent() {}
