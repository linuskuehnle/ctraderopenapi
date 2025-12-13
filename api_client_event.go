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

	"context"
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
)

/*
Pass through functions
*/
/**/

func (c *apiClient) SubscribeAPIEvent(eventData SubscribableAPIEventData) error {
	ctid, err := c.subscribeAPIEvent(eventData, false)
	if err != nil {
		return err
	}

	return c.accManager.AddEventSubscription(ctid, eventData.EventType, eventData.SubcriptionData)
}

func (c *apiClient) subscribeAPIEvent(eventData SubscribableAPIEventData, bypassReconnectBlock bool) (CtraderAccountId, error) {
	var ctid CtraderAccountId

	if eventData.SubcriptionData == nil {
		return ctid, &FunctionInvalidArgError{
			FunctionName: "SubscribeToEvent",
			Err:          errors.New("field SubscriptionData mustn't be nil"),
		}
	}
	if err := eventData.SubcriptionData.CheckError(); err != nil {
		return ctid, err
	}

	var reqData RequestData

	switch eventData.EventType {
	case APIEventType_Spots:
		subData, ok := eventData.SubcriptionData.(*SubscriptionDataSpotEvent)
		if !ok {
			return ctid, &FunctionInvalidArgError{
				FunctionName: "SubscribeToEvent",
				Err:          errors.New("invalid subscription data type for event type Spots"),
			}
		}

		ctid = subData.CtraderAccountId
		reqData = RequestData{
			ReqType: PROTO_OA_SUBSCRIBE_SPOTS_REQ,
			Req: &messages.ProtoOASubscribeSpotsReq{
				CtidTraderAccountId: (*int64)(&ctid),
				SymbolId:            subData.SymbolIds,
			},
			ResType: PROTO_OA_SUBSCRIBE_SPOTS_RES,
			Res:     &messages.ProtoOASubscribeSpotsRes{},
		}
	case APIEventType_LiveTrendbars:
		subData, ok := eventData.SubcriptionData.(*SubscriptionDataLiveTrendbarEvent)
		if !ok {
			return ctid, &FunctionInvalidArgError{
				FunctionName: "SubscribeToEvent",
				Err:          fmt.Errorf("invalid subscription data type for event type LiveTrendbars"),
			}
		}

		ctid = subData.CtraderAccountId
		reqData = RequestData{
			ReqType: PROTO_OA_SUBSCRIBE_LIVE_TRENDBAR_REQ,
			Req: &messages.ProtoOASubscribeLiveTrendbarReq{
				CtidTraderAccountId: (*int64)(&ctid),
				SymbolId:            &subData.SymbolId,
				Period:              (*ProtoOATrendbarPeriod)(&subData.Period),
			},
			ResType: PROTO_OA_SUBSCRIBE_LIVE_TRENDBAR_RES,
			Res:     &messages.ProtoOASubscribeLiveTrendbarRes{},
		}
	case APIEventType_DepthQuotes:
		subData, ok := eventData.SubcriptionData.(*SubscriptionDataDepthQuoteEvent)
		if !ok {
			return ctid, &FunctionInvalidArgError{
				FunctionName: "SubscribeToEvent",
				Err:          fmt.Errorf("invalid subscription data type for event type DepthQuotes"),
			}
		}

		ctid = subData.CtraderAccountId
		reqData = RequestData{
			ReqType: PROTO_OA_SUBSCRIBE_DEPTH_QUOTES_REQ,
			Req: &messages.ProtoOASubscribeDepthQuotesReq{
				CtidTraderAccountId: (*int64)(&ctid),
				SymbolId:            subData.SymbolIds,
			},
			ResType: PROTO_OA_SUBSCRIBE_DEPTH_QUOTES_RES,
			Res:     &messages.ProtoOASubscribeDepthQuotesRes{},
		}
	default:
		return ctid, &FunctionInvalidArgError{
			FunctionName: "SubscribeToEvent",
			Err:          fmt.Errorf("event type '%d' is not subscribable", eventData.EventType),
		}
	}

	if !bypassReconnectBlock {
		c.lifecycleData.BlockUntilReconnected(context.Background())
	}

	return ctid, c.sendRequest(reqData)
}

func (c *apiClient) UnsubscribeAPIEvent(eventData SubscribableAPIEventData) error {
	ctid, err := c.unsubscribeAPIEvent(eventData, false)
	if err != nil {
		return err
	}

	c.accManager.RemoveEventSubscription(ctid, eventData.EventType)
	return nil
}

func (c *apiClient) unsubscribeAPIEvent(eventData SubscribableAPIEventData, bypassReconnectBlock bool) (CtraderAccountId, error) {
	var ctid CtraderAccountId

	if eventData.SubcriptionData == nil {
		return ctid, &FunctionInvalidArgError{
			FunctionName: "SubscribeToEvent",
			Err:          errors.New("field SubscriptionData mustn't be nil"),
		}
	}
	if err := eventData.SubcriptionData.CheckError(); err != nil {
		return ctid, err
	}

	var reqData RequestData

	switch eventData.EventType {
	case APIEventType_Spots:
		subData, ok := eventData.SubcriptionData.(*SubscriptionDataSpotEvent)
		if !ok {
			return ctid, &FunctionInvalidArgError{
				FunctionName: "SubscribeToEvent",
				Err:          errors.New("invalid subscription data type for event type Spots"),
			}
		}

		ctid = subData.CtraderAccountId
		reqData = RequestData{
			ReqType: PROTO_OA_UNSUBSCRIBE_SPOTS_REQ,
			Req: &messages.ProtoOAUnsubscribeSpotsReq{
				CtidTraderAccountId: (*int64)(&ctid),
				SymbolId:            subData.SymbolIds,
			},
			ResType: PROTO_OA_UNSUBSCRIBE_SPOTS_RES,
			Res:     &messages.ProtoOAUnsubscribeSpotsRes{},
		}
	case APIEventType_LiveTrendbars:
		subData, ok := eventData.SubcriptionData.(*SubscriptionDataLiveTrendbarEvent)
		if !ok {
			return ctid, &FunctionInvalidArgError{
				FunctionName: "SubscribeToEvent",
				Err:          errors.New("invalid subscription data type for event type LiveTrendbars"),
			}
		}

		ctid = subData.CtraderAccountId
		reqData = RequestData{
			ReqType: PROTO_OA_UNSUBSCRIBE_LIVE_TRENDBAR_REQ,
			Req: &messages.ProtoOAUnsubscribeLiveTrendbarReq{
				CtidTraderAccountId: (*int64)(&ctid),
				SymbolId:            &subData.SymbolId,
				Period:              (*ProtoOATrendbarPeriod)(&subData.Period),
			},
			ResType: PROTO_OA_UNSUBSCRIBE_LIVE_TRENDBAR_RES,
			Res:     &messages.ProtoOAUnsubscribeLiveTrendbarRes{},
		}
	case APIEventType_DepthQuotes:
		subData, ok := eventData.SubcriptionData.(*SubscriptionDataDepthQuoteEvent)
		if !ok {
			return ctid, &FunctionInvalidArgError{
				FunctionName: "SubscribeToEvent",
				Err:          fmt.Errorf("invalid subscription data type for event type DepthQuotes"),
			}
		}

		ctid = subData.CtraderAccountId
		reqData = RequestData{
			ReqType: PROTO_OA_UNSUBSCRIBE_DEPTH_QUOTES_REQ,
			Req: &messages.ProtoOAUnsubscribeDepthQuotesReq{
				CtidTraderAccountId: (*int64)(&ctid),
				SymbolId:            subData.SymbolIds,
			},
			ResType: PROTO_OA_UNSUBSCRIBE_DEPTH_QUOTES_RES,
			Res:     &messages.ProtoOAUnsubscribeDepthQuotesRes{},
		}
	default:
		return ctid, &FunctionInvalidArgError{
			FunctionName: "SubscribeToEvent",
			Err:          fmt.Errorf("event type '%d' is not subscribable", eventData.EventType),
		}
	}

	if !bypassReconnectBlock {
		c.lifecycleData.BlockUntilReconnected(context.Background())
	}

	return ctid, c.sendRequest(reqData)
}

func (c *apiClient) ListenToAPIEvent(ctx context.Context, eventType eventType, onEventCh chan ListenableEvent) error {
	if onEventCh == nil {
		return &FunctionInvalidArgError{
			FunctionName: "ListenToAPIEvent",
			Err:          errors.New("onEventCh mustn't be nil"),
		}
	}
	if ctx == nil {
		ctx = context.Background()
	}

	var eventId datatypes.EventId
	var eventCallback func(event proto.Message)

	switch eventType {
	case APIEventType_Spots:
		eventId = datatypes.EventId(PROTO_OA_SPOT_EVENT)

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOASpotEvent)
		}
	case APIEventType_DepthQuotes:
		eventId = datatypes.EventId(PROTO_OA_DEPTH_EVENT)

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOADepthEvent)
		}

	case APIEventType_TrailingSLChanged:
		eventId = datatypes.EventId(PROTO_OA_TRAILING_SL_CHANGED_EVENT)

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOATrailingSLChangedEvent)
		}
	case APIEventType_SymbolChanged:
		eventId = datatypes.EventId(PROTO_OA_SYMBOL_CHANGED_EVENT)

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOASymbolChangedEvent)
		}
	case APIEventType_TraderUpdated:
		eventId = datatypes.EventId(PROTO_OA_TRADER_UPDATE_EVENT)

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOATraderUpdatedEvent)
		}
	case APIEventType_Execution:
		eventId = datatypes.EventId(PROTO_OA_EXECUTION_EVENT)

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOAExecutionEvent)
		}
	case APIEventType_OrderError:
		eventId = datatypes.EventId(PROTO_OA_ORDER_ERROR_EVENT)

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOAOrderErrorEvent)
		}
	case APIEventType_MarginChanged:
		eventId = datatypes.EventId(PROTO_OA_MARGIN_CHANGED_EVENT)

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOAMarginChangedEvent)
		}
	case APIEventType_AccountsTokenInvalidated:
		eventId = datatypes.EventId(PROTO_OA_ACCOUNTS_TOKEN_INVALIDATED_EVENT)

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOAAccountsTokenInvalidatedEvent)
		}
	case APIEventType_ClientDisconnect:
		eventId = datatypes.EventId(PROTO_OA_CLIENT_DISCONNECT_EVENT)

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOAClientDisconnectEvent)
		}
	case APIEventType_AccountDisconnect:
		eventId = datatypes.EventId(PROTO_OA_ACCOUNT_DISCONNECT_EVENT)

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOAAccountDisconnectEvent)
		}
	case APIEventType_MarginCallUpdate:
		eventId = datatypes.EventId(PROTO_OA_MARGIN_CALL_UPDATE_EVENT)

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOAMarginCallUpdateEvent)
		}
	case APIEventType_MarginCallTrigger:
		eventId = datatypes.EventId(PROTO_OA_MARGIN_CALL_TRIGGER_EVENT)

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOAMarginCallTriggerEvent)
		}
	default:
		return &FunctionInvalidArgError{
			FunctionName: "ListenToAPIEvent",
			Err:          fmt.Errorf("event type '%d' is not listenable", eventType),
		}
	}

	if err := c.apiEventHandler.AddEvent(eventId, eventCallback); err != nil {
		return err
	}

	removeCallback := func() error {
		err := c.apiEventHandler.RemoveEvent(eventId)
		close(onEventCh)
		return err
	}
	go c.runListenerRemove(ctx, removeCallback)

	return nil
}

func (c *apiClient) ListenToClientEvent(ctx context.Context, clientEventType clientEventType, onEventCh chan ListenableClientEvent) error {
	if onEventCh == nil {
		return &FunctionInvalidArgError{
			FunctionName: "ListenToClientEvent",
			Err:          errors.New("onEventCh mustn't be nil"),
		}
	}
	if ctx == nil {
		ctx = context.Background()
	}

	eventId := datatypes.EventId(clientEventType)
	eventCallback := func(event ListenableClientEvent) {
		onEventCh <- event
	}

	if err := c.clientEventHandler.AddEvent(eventId, eventCallback); err != nil {
		return err
	}

	removeCallback := func() error {
		err := c.clientEventHandler.RemoveEvent(eventId)

		close(onEventCh)
		return err
	}
	go c.runListenerRemove(ctx, removeCallback)

	return nil
}

/*
Event specific listen functions
*/
/**/

// runListenerRemove runs a goroutine that waits until either the listener context or the life cycle context is done
func (c *apiClient) runListenerRemove(listenerCtx context.Context, removeCallback func() error) {
	lifecycleCtx, _ := c.lifecycleData.GetContext()

	// Wait until either the listener context or the life cycle context is done
	// This ensures that if the api client is disconnected and all events are removed
	// that the listener goroutine also exits

	select {
	case <-listenerCtx.Done():
	case <-lifecycleCtx.Done():
	}

	// removeCallback does only contain RemoveEvent which in this context cannot return an error.
	// Hence we do not check the returned error value
	removeCallback()
}

// handleListenableEvent is called on an incoming event message that is listenable
func (c *apiClient) handleListenableAPIEvent(msgType ProtoOAPayloadType, protoMsg *messages.ProtoMessage) error {
	var eventId datatypes.EventId = datatypes.EventId(msgType)
	if !c.apiEventHandler.HasEvent(eventId) {
		// No listener for this event, ignore
		return nil
	}

	payloadBytes := protoMsg.GetPayload()

	var event proto.Message

	switch msgType {
	case PROTO_OA_SPOT_EVENT:
		var eventMsg ProtoOASpotEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA spot event",
				Err:         err,
			}
			return perr
		}

		event = &eventMsg
	case PROTO_OA_DEPTH_EVENT:
		var eventMsg ProtoOADepthEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA depth event",
				Err:         err,
			}
			return perr
		}

		event = &eventMsg

	case PROTO_OA_TRAILING_SL_CHANGED_EVENT:
		var eventMsg ProtoOATrailingSLChangedEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA trailing sl changed event",
				Err:         err,
			}
			return perr
		}

		event = &eventMsg
	case PROTO_OA_SYMBOL_CHANGED_EVENT:
		var eventMsg ProtoOASymbolChangedEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA symbol changed event",
				Err:         err,
			}
			return perr
		}

		event = &eventMsg
	case PROTO_OA_TRADER_UPDATE_EVENT:
		var eventMsg ProtoOATraderUpdatedEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA trader updated event",
				Err:         err,
			}
			return perr
		}

		event = &eventMsg
	case PROTO_OA_EXECUTION_EVENT:
		var eventMsg ProtoOAExecutionEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA execution event",
				Err:         err,
			}
			return perr
		}

		event = &eventMsg
	case PROTO_OA_ORDER_ERROR_EVENT:
		var eventMsg ProtoOAOrderErrorEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA order error event",
				Err:         err,
			}
			return perr
		}

		event = &eventMsg
	case PROTO_OA_MARGIN_CHANGED_EVENT:
		var eventMsg ProtoOAMarginChangedEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA margin changed event",
				Err:         err,
			}
			return perr
		}

		event = &eventMsg
	case PROTO_OA_ACCOUNTS_TOKEN_INVALIDATED_EVENT:
		var eventMsg ProtoOAAccountsTokenInvalidatedEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA accounts token invalidated event",
				Err:         err,
			}
			return perr
		}

		event = &eventMsg
	case PROTO_OA_CLIENT_DISCONNECT_EVENT:
		var eventMsg ProtoOAClientDisconnectEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA client disconnect event",
				Err:         err,
			}
			return perr
		}

		event = &eventMsg
	case PROTO_OA_ACCOUNT_DISCONNECT_EVENT:
		var depthEventMsg ProtoOAAccountDisconnectEvent

		if err := proto.Unmarshal(payloadBytes, &depthEventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA account disconnect event",
				Err:         err,
			}
			return perr
		}

		event = &depthEventMsg
	case PROTO_OA_MARGIN_CALL_UPDATE_EVENT:
		var eventMsg ProtoOAMarginCallUpdateEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA margin call update event",
				Err:         err,
			}
			return perr
		}

		event = &eventMsg
	case PROTO_OA_MARGIN_CALL_TRIGGER_EVENT:
		var eventMsg ProtoOAMarginCallTriggerEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA margin call trigger event",
				Err:         err,
			}
			return perr
		}

		event = &eventMsg
	default:
		return &UnexpectedMessageTypeError{
			MsgType: msgType,
		}
	}

	eventHandled := c.apiEventHandler.HandleEvent(eventId, event)
	if !eventHandled {
		return &IdNotIncludedError{
			Id: eventId,
		}
	}
	return nil
}

// SpawnAPIEventHandler starts a simple goroutine that reads `ListenableEvent`
// values from `onEventCh`, attempts to cast each value to the concrete
// type `T` and invokes `onEvent` for matching values.
//
// This helper is useful when callers want a lightweight adapter: it
// performs the type assertion and drops events that don't match `T`.
//
// Important notes:
//   - `ctx` controls the lifetime of the handler. When `ctx` is canceled
//     the goroutine exits. `onEventCh` must be closed by the sender to
//     terminate the loop if `ctx` remains active.
//   - `onEventCh` will be read-only for the goroutine; callers must not
//     close it while other readers expect it to remain open.
//   - The helper does not recover from panics in `onEvent`; if the typed
//     handler may panic wrap it appropriately.
func SpawnAPIEventHandler[T ListenableEvent](ctx context.Context, onEventCh chan ListenableEvent, onEvent func(T)) error {
	if ctx == nil {
		return &FunctionInvalidArgError{
			FunctionName: "SpawnAPIEventHandler",
			Err:          errors.New("ctx may not be nil"),
		}
	}
	if onEventCh == nil {
		return &FunctionInvalidArgError{
			FunctionName: "SpawnAPIEventHandler",
			Err:          errors.New("onEventCh may not be nil"),
		}
	}
	if onEvent == nil {
		return &FunctionInvalidArgError{
			FunctionName: "SpawnAPIEventHandler",
			Err:          errors.New("onEvent may not be nil"),
		}
	}

	go func(onEvent func(T)) {
		for event := range onEventCh {
			typedEvent, ok := CastToEventType[T](event)
			if ok {
				onEvent(typedEvent)
			}
		}
	}(onEvent)
	return nil
}

// SpawnClientEventHandler starts a simple goroutine that reads `ListenableClientEvent`
// values from `onEventCh`, attempts to cast each value to the concrete
// type `T` and invokes `onEvent` for matching values.
//
// This helper is useful when callers want a lightweight adapter: it
// performs the type assertion and drops events that don't match `T`.
//
// Important notes:
//   - `ctx` controls the lifetime of the handler. When `ctx` is canceled
//     the goroutine exits. `onEventCh` must be closed by the sender to
//     terminate the loop if `ctx` remains active.
//   - `onEventCh` will be read-only for the goroutine; callers must not
//     close it while other readers expect it to remain open.
//   - The helper does not recover from panics in `onEvent`; if the typed
//     handler may panic wrap it appropriately.
func SpawnClientEventHandler[T ListenableClientEvent](ctx context.Context, onEventCh chan ListenableClientEvent, onEvent func(T)) error {
	if ctx == nil {
		return &FunctionInvalidArgError{
			FunctionName: "SpawnClientEventHandler",
			Err:          errors.New("ctx may not be nil"),
		}
	}
	if onEventCh == nil {
		return &FunctionInvalidArgError{
			FunctionName: "SpawnClientEventHandler",
			Err:          errors.New("onEventCh may not be nil"),
		}
	}
	if onEvent == nil {
		return &FunctionInvalidArgError{
			FunctionName: "SpawnClientEventHandler",
			Err:          errors.New("onEvent may not be nil"),
		}
	}

	go func(onEvent func(T)) {
		for event := range onEventCh {
			typedEvent, ok := CastToClientEventType[T](event)
			if ok {
				onEvent(typedEvent)
			}
		}
	}(onEvent)
	return nil
}
