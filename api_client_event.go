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

func (c *apiClient) SubscribeAPIEvent(eventData APIEventSubData) error {
	ctid, err := c.subscribeAPIEvent(eventData, false)
	if err != nil {
		return err
	}

	return c.accManager.AddEventSubscription(ctid, eventData.EventType, eventData.SubcriptionData)
}

func (c *apiClient) subscribeAPIEvent(eventData APIEventSubData, bypassReconnectBlock bool) (CtraderAccountId, error) {
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
		subData, ok := eventData.SubcriptionData.(*SpotEventData)
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
		subData, ok := eventData.SubcriptionData.(*LiveTrendbarEventData)
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
		subData, ok := eventData.SubcriptionData.(*DepthQuoteEventData)
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

func (c *apiClient) UnsubscribeAPIEvent(eventData APIEventSubData) error {
	ctid, err := c.unsubscribeAPIEvent(eventData, false)
	if err != nil {
		return err
	}

	c.accManager.RemoveEventSubscription(ctid, eventData.EventType)
	return nil
}

func (c *apiClient) unsubscribeAPIEvent(eventData APIEventSubData, bypassReconnectBlock bool) (CtraderAccountId, error) {
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
		subData, ok := eventData.SubcriptionData.(*SpotEventData)
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
		subData, ok := eventData.SubcriptionData.(*LiveTrendbarEventData)
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
		subData, ok := eventData.SubcriptionData.(*DepthQuoteEventData)
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

func (c *apiClient) ListenToAPIEvent(ctx context.Context, eventData APIEventListenData) error {
	if err := eventData.CheckError(); err != nil {
		return err
	}
	if ctx == nil {
		ctx = context.Background()
	}

	var removeCallback func() error

	if eventData.EventKeyData != nil {
		if err := c.apiEventHandler.SetListener(eventData.EventType, eventData.EventKeyData, eventData.EventCh); err != nil {
			return err
		}

		removeCallback = func() error {
			return c.apiEventHandler.UnsetListener(eventData.EventType, eventData.EventKeyData)
		}
	} else {
		c.apiEventHandler.SetDefaultListener(eventData.EventType, eventData.EventCh)

		removeCallback = func() error {
			c.apiEventHandler.UnsetDefaultListener(eventData.EventType)
			return nil
		}
	}

	go c.runListenerRemove(ctx, removeCallback)

	return nil
}

func (c *apiClient) ListenToClientEvent(ctx context.Context, eventData ClientEventListenData) error {
	if err := eventData.CheckError(); err != nil {
		return err
	}
	if ctx == nil {
		ctx = context.Background()
	}

	c.clientEventHandler.SetDefaultListener(eventData.EventType, eventData.EventCh)

	removeCallback := func() error {
		c.clientEventHandler.UnsetDefaultListener(eventData.EventType)
		return nil
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

// handleAPIEvent is called on an incoming event message that is listenable
func (c *apiClient) handleAPIEvent(msgType ProtoOAPayloadType, protoMsg *messages.ProtoMessage) error {
	var eventType apiEventType = apiEventType(msgType)

	hasEvent := c.apiEventHandler.HasListenerSource(eventType)
	if !hasEvent && !hasHookForAPIEvent[msgType] {
		// No listener and hook for this event, so ignore it
		return nil
	}

	payloadBytes := protoMsg.GetPayload()

	var events []proto.Message

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

		if eventMsg.Trendbar != nil {
			// Filter trendbars by included periods
			trendbarsByPeriod := make(map[ProtoOATrendbarPeriod][]*messages.ProtoOATrendbar)
			for _, trendbar := range eventMsg.GetTrendbar() {
				period := trendbar.GetPeriod()

				trendbarsByPeriod[period] = append(trendbarsByPeriod[period], trendbar)
			}

			for period, trendbars := range trendbarsByPeriod {
				events = append(events, &ProtoOASpotEvent{
					PayloadType:         eventMsg.PayloadType,
					CtidTraderAccountId: eventMsg.CtidTraderAccountId,
					SymbolId:            eventMsg.SymbolId,
					Bid:                 eventMsg.Bid,
					Ask:                 eventMsg.Ask,
					Trendbar:            trendbars,
					SessionClose:        eventMsg.SessionClose,
					Timestamp:           eventMsg.Timestamp,
					TrendbarPeriod:      &period,
				})
			}
		} else {
			events = append(events, &eventMsg)
		}

		events = append(events, &eventMsg)
	case PROTO_OA_DEPTH_EVENT:
		var eventMsg ProtoOADepthEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA depth event",
				Err:         err,
			}
			return perr
		}

		events = append(events, &eventMsg)

	case PROTO_OA_TRAILING_SL_CHANGED_EVENT:
		var eventMsg ProtoOATrailingSLChangedEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA trailing sl changed event",
				Err:         err,
			}
			return perr
		}

		events = append(events, &eventMsg)
	case PROTO_OA_SYMBOL_CHANGED_EVENT:
		var eventMsg ProtoOASymbolChangedEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA symbol changed event",
				Err:         err,
			}
			return perr
		}

		events = append(events, &eventMsg)
	case PROTO_OA_TRADER_UPDATE_EVENT:
		var eventMsg ProtoOATraderUpdatedEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA trader updated event",
				Err:         err,
			}
			return perr
		}

		events = append(events, &eventMsg)
	case PROTO_OA_EXECUTION_EVENT:
		var eventMsg ProtoOAExecutionEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA execution event",
				Err:         err,
			}
			return perr
		}

		events = append(events, &eventMsg)
	case PROTO_OA_ORDER_ERROR_EVENT:
		var eventMsg ProtoOAOrderErrorEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA order error event",
				Err:         err,
			}
			return perr
		}

		events = append(events, &eventMsg)
	case PROTO_OA_MARGIN_CHANGED_EVENT:
		var eventMsg ProtoOAMarginChangedEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA margin changed event",
				Err:         err,
			}
			return perr
		}

		events = append(events, &eventMsg)
	case PROTO_OA_ACCOUNTS_TOKEN_INVALIDATED_EVENT:
		var eventMsg ProtoOAAccountsTokenInvalidatedEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA accounts token invalidated event",
				Err:         err,
			}
			return perr
		}

		events = append(events, &eventMsg)
	case PROTO_OA_CLIENT_DISCONNECT_EVENT:
		if err := c.authenticateApp(); err != nil {
			c.fatalErrCh <- fmt.Errorf("failed to re-authenticate app after client disconnect: %w", err)
			return nil
		}

		var eventMsg ProtoOAClientDisconnectEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA client disconnect event",
				Err:         err,
			}
			return perr
		}

		events = append(events, &eventMsg)
	case PROTO_OA_ACCOUNT_DISCONNECT_EVENT:
		var eventMsg ProtoOAAccountDisconnectEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA account disconnect event",
				Err:         err,
			}
			return perr
		}

		ctid := CtraderAccountId(eventMsg.GetCtidTraderAccountId())
		c.onAccountDisconnect(ctid)

		events = append(events, &eventMsg)
	case PROTO_OA_MARGIN_CALL_UPDATE_EVENT:
		var eventMsg ProtoOAMarginCallUpdateEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA margin call update event",
				Err:         err,
			}
			return perr
		}

		events = append(events, &eventMsg)
	case PROTO_OA_MARGIN_CALL_TRIGGER_EVENT:
		var eventMsg ProtoOAMarginCallTriggerEvent

		if err := proto.Unmarshal(payloadBytes, &eventMsg); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto OA margin call trigger event",
				Err:         err,
			}
			return perr
		}

		events = append(events, &eventMsg)
	default:
		return &UnexpectedMessageTypeError{
			MsgType: msgType,
		}
	}

	if !hasEvent {
		// No listener for this event, so ignore it
		return nil
	}

	for _, event := range events {
		c.apiEventHandler.HandleEvent(eventType, event.(APIEvent))
	}

	return nil
}

// SpawnAPIEventHandler starts a simple goroutine that reads `APIEvent`
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
func SpawnAPIEventHandler[T APIEvent](ctx context.Context, onEventCh chan APIEvent, onEvent func(T)) error {
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

	go func(onEvent func(T), onEventCh chan APIEvent) {
		for event := range onEventCh {
			typedEvent, ok := CastToAPIEventType[T](event)
			if ok {
				onEvent(typedEvent)
			}
		}
	}(onEvent, onEventCh)
	return nil
}

// SpawnClientEventHandler starts a simple goroutine that reads `ClientEvent`
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
func SpawnClientEventHandler[T ClientEvent](ctx context.Context, onEventCh chan ClientEvent, onEvent func(T)) error {
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

	go func(onEvent func(T), onEventCh chan ClientEvent) {
		for event := range onEventCh {
			typedEvent, ok := CastToClientEventType[T](event)
			if ok {
				onEvent(typedEvent)
			}
		}
	}(onEvent, onEventCh)
	return nil
}
