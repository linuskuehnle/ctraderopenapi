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

func (h *apiHandler) SubscribeEvent(eventData SubscribableEventData) error {
	if eventData.SubcriptionData == nil {
		return &FunctionInvalidArgError{
			FunctionName: "SubscribeToEvent",
			Err:          errors.New("field SubscriptionData mustn't be nil"),
		}
	}
	if err := eventData.SubcriptionData.CheckError(); err != nil {
		return err
	}

	var reqData RequestData

	switch eventData.EventType {
	case EventType_Spots:
		subData, ok := eventData.SubcriptionData.(*SubscriptionDataSpotEvent)
		if !ok {
			return &FunctionInvalidArgError{
				FunctionName: "SubscribeToEvent",
				Err:          errors.New("invalid subscription data type for event type Spots"),
			}
		}

		reqData = RequestData{
			ReqType: PROTO_OA_SUBSCRIBE_SPOTS_REQ,
			Req: &messages.ProtoOASubscribeSpotsReq{
				CtidTraderAccountId: (*int64)(&subData.CtraderAccountId),
				SymbolId:            subData.SymbolIds,
			},
			ResType: PROTO_OA_SUBSCRIBE_SPOTS_RES,
			Res:     &messages.ProtoOASubscribeSpotsRes{},
		}
	case EventType_LiveTrendbars:
		subData, ok := eventData.SubcriptionData.(*SubscriptionDataLiveTrendbarEvent)
		if !ok {
			return &FunctionInvalidArgError{
				FunctionName: "SubscribeToEvent",
				Err:          fmt.Errorf("invalid subscription data type for event type LiveTrendbars"),
			}
		}

		reqData = RequestData{
			ReqType: PROTO_OA_SUBSCRIBE_LIVE_TRENDBAR_REQ,
			Req: &messages.ProtoOASubscribeLiveTrendbarReq{
				CtidTraderAccountId: (*int64)(&subData.CtraderAccountId),
				SymbolId:            &subData.SymbolId,
				Period:              (*messages.ProtoOATrendbarPeriod)(&subData.Period),
			},
			ResType: PROTO_OA_SUBSCRIBE_LIVE_TRENDBAR_RES,
			Res:     &messages.ProtoOASubscribeLiveTrendbarRes{},
		}
	case EventType_DepthQuotes:
		subData, ok := eventData.SubcriptionData.(*SubscriptionDataDepthQuoteEvent)
		if !ok {
			return &FunctionInvalidArgError{
				FunctionName: "SubscribeToEvent",
				Err:          fmt.Errorf("invalid subscription data type for event type DepthQuotes"),
			}
		}

		reqData = RequestData{
			ReqType: PROTO_OA_SUBSCRIBE_DEPTH_QUOTES_REQ,
			Req: &messages.ProtoOASubscribeDepthQuotesReq{
				CtidTraderAccountId: (*int64)(&subData.CtraderAccountId),
				SymbolId:            subData.SymbolIds,
			},
			ResType: PROTO_OA_SUBSCRIBE_DEPTH_QUOTES_RES,
			Res:     &messages.ProtoOASubscribeDepthQuotesRes{},
		}
	default:
		return &FunctionInvalidArgError{
			FunctionName: "SubscribeToEvent",
			Err:          fmt.Errorf("event type '%d' is not subscribable", eventData.EventType),
		}
	}

	return h.SendRequest(reqData)
}

func (h *apiHandler) UnsubscribeEvent(eventData SubscribableEventData) error {
	if eventData.SubcriptionData == nil {
		return &FunctionInvalidArgError{
			FunctionName: "SubscribeToEvent",
			Err:          errors.New("field SubscriptionData mustn't be nil"),
		}
	}
	if err := eventData.SubcriptionData.CheckError(); err != nil {
		return err
	}

	var reqData RequestData

	switch eventData.EventType {
	case EventType_Spots:
		subData, ok := eventData.SubcriptionData.(*SubscriptionDataSpotEvent)
		if !ok {
			return &FunctionInvalidArgError{
				FunctionName: "SubscribeToEvent",
				Err:          errors.New("invalid subscription data type for event type Spots"),
			}
		}

		reqData = RequestData{
			ReqType: PROTO_OA_UNSUBSCRIBE_SPOTS_REQ,
			Req: &messages.ProtoOAUnsubscribeSpotsReq{
				CtidTraderAccountId: (*int64)(&subData.CtraderAccountId),
				SymbolId:            subData.SymbolIds,
			},
			ResType: PROTO_OA_UNSUBSCRIBE_SPOTS_RES,
			Res:     &messages.ProtoOAUnsubscribeSpotsRes{},
		}
	case EventType_LiveTrendbars:
		subData, ok := eventData.SubcriptionData.(*SubscriptionDataLiveTrendbarEvent)
		if !ok {
			return &FunctionInvalidArgError{
				FunctionName: "SubscribeToEvent",
				Err:          errors.New("invalid subscription data type for event type LiveTrendbars"),
			}
		}

		reqData = RequestData{
			ReqType: PROTO_OA_UNSUBSCRIBE_LIVE_TRENDBAR_REQ,
			Req: &messages.ProtoOAUnsubscribeLiveTrendbarReq{
				CtidTraderAccountId: (*int64)(&subData.CtraderAccountId),
				SymbolId:            &subData.SymbolId,
				Period:              (*messages.ProtoOATrendbarPeriod)(&subData.Period),
			},
			ResType: PROTO_OA_UNSUBSCRIBE_LIVE_TRENDBAR_RES,
			Res:     &messages.ProtoOAUnsubscribeLiveTrendbarRes{},
		}
	case EventType_DepthQuotes:
		subData, ok := eventData.SubcriptionData.(*SubscriptionDataDepthQuoteEvent)
		if !ok {
			return &FunctionInvalidArgError{
				FunctionName: "SubscribeToEvent",
				Err:          fmt.Errorf("invalid subscription data type for event type DepthQuotes"),
			}
		}

		reqData = RequestData{
			ReqType: PROTO_OA_UNSUBSCRIBE_DEPTH_QUOTES_REQ,
			Req: &messages.ProtoOAUnsubscribeDepthQuotesReq{
				CtidTraderAccountId: (*int64)(&subData.CtraderAccountId),
				SymbolId:            subData.SymbolIds,
			},
			ResType: PROTO_OA_UNSUBSCRIBE_DEPTH_QUOTES_RES,
			Res:     &messages.ProtoOAUnsubscribeDepthQuotesRes{},
		}
	default:
		return &FunctionInvalidArgError{
			FunctionName: "SubscribeToEvent",
			Err:          fmt.Errorf("event type '%d' is not subscribable", eventData.EventType),
		}
	}

	return h.SendRequest(reqData)
}

func (h *apiHandler) ListenToEvent(eventType eventType, onEventCh chan ListenableEvent, ctx context.Context) error {
	if onEventCh == nil {
		return &FunctionInvalidArgError{
			FunctionName: "ListenToEvent",
			Err:          errors.New("onEventCh mustn't be nil"),
		}
	}
	if ctx == nil {
		ctx = context.Background()
	}

	var eventId ProtoOAPayloadType
	var eventCallback func(event proto.Message)

	switch eventType {
	case EventType_Spots:
		eventId = PROTO_OA_SPOT_EVENT

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*messages.ProtoOASpotEvent)
		}
	case EventType_DepthQuotes:
		eventId = PROTO_OA_DEPTH_EVENT

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*messages.ProtoOADepthEvent)
		}

	case EventType_TrailingSLChanged:
		eventId = PROTO_OA_TRAILING_SL_CHANGED_EVENT

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOATrailingSLChangedEvent)
		}
	case EventType_SymbolChanged:
		eventId = PROTO_OA_SYMBOL_CHANGED_EVENT

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOASymbolChangedEvent)
		}
	case EventType_TraderUpdated:
		eventId = PROTO_OA_TRADER_UPDATE_EVENT

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOATraderUpdatedEvent)
		}
	case EventType_Execution:
		eventId = PROTO_OA_EXECUTION_EVENT

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOAExecutionEvent)
		}
	case EventType_OrderError:
		eventId = PROTO_OA_ORDER_ERROR_EVENT

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOAOrderErrorEvent)
		}
	case EventType_MarginChanged:
		eventId = PROTO_OA_MARGIN_CHANGED_EVENT

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOAMarginChangedEvent)
		}
	case EventType_AccountsTokenInvalidated:
		eventId = PROTO_OA_ACCOUNTS_TOKEN_INVALIDATED_EVENT

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOAAccountsTokenInvalidatedEvent)
		}
	case EventType_ClientDisconnect:
		eventId = PROTO_OA_CLIENT_DISCONNECT_EVENT

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOAClientDisconnectEvent)
		}
	case EventType_AccountDisconnect:
		eventId = PROTO_OA_ACCOUNT_DISCONNECT_EVENT

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOAAccountDisconnectEvent)
		}
	case EventType_MarginCallUpdate:
		eventId = PROTO_OA_MARGIN_CALL_UPDATE_EVENT

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOAMarginCallUpdateEvent)
		}
	case EventType_MarginCallTrigger:
		eventId = PROTO_OA_MARGIN_CALL_TRIGGER_EVENT

		eventCallback = func(event proto.Message) {
			onEventCh <- event.(*ProtoOAMarginCallTriggerEvent)
		}
	default:
		return &FunctionInvalidArgError{
			FunctionName: "ListenToEvent",
			Err:          fmt.Errorf("event type '%d' is not listenable", eventType),
		}
	}

	if err := h.eventHandler.AddEvent(eventId, eventCallback); err != nil {
		return err
	}

	removeCallback := func() error {
		close(onEventCh)
		return h.eventHandler.RemoveEvent(eventId)
	}
	go h.runListenerRemove(ctx, removeCallback)

	return nil
}

/*
Event specific listen functions
*/
/**/

// runListenerRemove runs a goroutine that waits until either the listener context or the life cycle context is done
func (h *apiHandler) runListenerRemove(listenerCtx context.Context, removeCallback func() error) {
	lifeCycleCtx, _ := h.lifeCycleData.GetContext()

	// Wait until either the listener context or the life cycle context is done
	// This ensures that if the api handler is disconnected and all events are removed
	// that the listener goroutine also exits

	select {
	case <-listenerCtx.Done():
	case <-lifeCycleCtx.Done():
	}

	// removeCallback does only contain RemoveEvent which in this context cannot return an error.
	// Hence we do not check the returned error value
	removeCallback()
}

// handleListenableEvent is called on an incoming event message that is listenable
func (h *apiHandler) handleListenableEvent(msgType ProtoOAPayloadType, protoMsg *messages.ProtoMessage) error {
	eventId := msgType
	if !h.eventHandler.HasEvent(eventId) {
		// No listener for this event, ignore
		return nil
	}

	payloadBytes := protoMsg.GetPayload()

	var event datatypes.Event

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

	return h.eventHandler.HandleEvent(eventId, event)
}

// SpawnEventHandler starts a simple goroutine that reads `ListenableEvent`
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
func SpawnEventHandler[T ListenableEvent](ctx context.Context, onEventCh chan ListenableEvent, onEvent func(T)) error {
	if ctx == nil {
		return &FunctionInvalidArgError{
			FunctionName: "SpawnEventHandler",
			Err:          errors.New("ctx may not be nil"),
		}
	}
	if onEventCh == nil {
		return &FunctionInvalidArgError{
			FunctionName: "SpawnEventHandler",
			Err:          errors.New("onEventCh may not be nil"),
		}
	}
	if onEvent == nil {
		return &FunctionInvalidArgError{
			FunctionName: "SpawnEventHandler",
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
