package messages

import (
	"fmt"
)

type APIEventKey string

// BaseEvent is the interface that all events must implement to be distributed.
type BaseEvent interface {
	ToDistributableEvent() DistributableEvent
}

// APIEvent marks protobuf API event message types that can be listened to
// (push-style events).
//
// These event types are delivered by the server independently from any
// single client's request; users can register listener channels of type
// `APIEvent` using `ListenToAPIEvent`.
type APIEvent interface {
	BaseEvent
	IsListenableAPIEvent()
}

// DistributableEvent marks protobuf API event message types that can be
// listened to (push-style events) AND can be used in an EventDistributor
// instance.
//
// These event types are delivered by the server independently from any
// single client's request; users can register channels using `AddListener`
// on EventDistributor.
type DistributableEvent interface {
	APIEvent
	BuildKeys() []APIEventKey
}

// ClientEvent marks events emitted by the API client itself
// (for example connection loss, reconnect success, etc).
//
// These event types are emitted internally by the client based on the
// client behaviour and its connection status; users can register listener
// channels of type `ClientEvent` using `ListenToClientEvent`.
type ClientEvent interface {
	BaseEvent
	IsListenableClientEvent()
}

/*
Subscribable event assertions
*/
/**/

func (*ProtoOASpotEvent) IsSubscribableEvent()  {}
func (*ProtoOASpotEvent) IsListenableAPIEvent() {}
func (e *ProtoOASpotEvent) BuildKeys() []APIEventKey {
	return []APIEventKey{
		APIEventKey(fmt.Sprintf("SpotEvent__Ctid:%d__SymbolId:%d__TrendbarPeriod:%d",
			e.GetCtidTraderAccountId(),
			e.GetSymbolId(),
			e.GetTrendbarPeriod(),
		)),
		APIEventKey(fmt.Sprintf("SpotEvent__Ctid:%d__SymbolId:%d__TrendbarPeriod:",
			e.GetCtidTraderAccountId(),
			e.GetSymbolId(),
		)),
		APIEventKey(fmt.Sprintf("SpotEvent__Ctid:%d__SymbolId:__TrendbarPeriod:%d",
			e.GetCtidTraderAccountId(),
			e.GetTrendbarPeriod(),
		)),
		APIEventKey(fmt.Sprintf("SpotEvent__Ctid:%d__SymbolId:__TrendbarPeriod:",
			e.GetCtidTraderAccountId(),
		)),
		APIEventKey(fmt.Sprintf("SpotEvent__Ctid:__SymbolId:%d__TrendbarPeriod:%d",
			e.GetSymbolId(),
			e.GetTrendbarPeriod(),
		)),
		APIEventKey(fmt.Sprintf("SpotEvent__Ctid:__SymbolId:%d__TrendbarPeriod:",
			e.GetSymbolId(),
		)),
		APIEventKey(fmt.Sprintf("SpotEvent__Ctid:__SymbolId:__TrendbarPeriod:%d",
			e.GetTrendbarPeriod(),
		)),
	}
}
func (e *ProtoOASpotEvent) ToDistributableEvent() DistributableEvent {
	return e
}

func (*ProtoOADepthEvent) IsSubscribableEvent()  {}
func (*ProtoOADepthEvent) IsListenableAPIEvent() {}
func (e *ProtoOADepthEvent) BuildKeys() []APIEventKey {
	return []APIEventKey{
		APIEventKey(fmt.Sprintf("DepthEvent__Ctid:%d__SymbolId:%d",
			e.GetCtidTraderAccountId(),
			e.GetSymbolId(),
		)),
		APIEventKey(fmt.Sprintf("DepthEvent__Ctid:%d__SymbolId:",
			e.GetCtidTraderAccountId(),
		)),
		APIEventKey(fmt.Sprintf("DepthEvent__Ctid:__SymbolId:%d",
			e.GetSymbolId(),
		)),
	}
}
func (e *ProtoOADepthEvent) ToDistributableEvent() DistributableEvent {
	return e
}

/*
Listenable event assertions
*/
/**/

func (*ProtoOATrailingSLChangedEvent) IsListenableAPIEvent() {}
func (e *ProtoOATrailingSLChangedEvent) BuildKeys() []APIEventKey {
	return []APIEventKey{
		APIEventKey(fmt.Sprintf("TrailingSLChangedEvent__Ctid:%d_PositionId:%d__OrderId:%d",
			e.GetCtidTraderAccountId(),
			e.GetPositionId(),
			e.GetOrderId(),
		)),
		APIEventKey(fmt.Sprintf("TrailingSLChangedEvent__Ctid:%d_PositionId:%d__OrderId:",
			e.GetCtidTraderAccountId(),
			e.GetPositionId(),
		)),
		APIEventKey(fmt.Sprintf("TrailingSLChangedEvent__Ctid:%d_PositionId:__OrderId:%d",
			e.GetCtidTraderAccountId(),
			e.GetOrderId(),
		)),
		APIEventKey(fmt.Sprintf("TrailingSLChangedEvent__Ctid:%d_PositionId:__OrderId:",
			e.GetCtidTraderAccountId(),
		)),
		APIEventKey(fmt.Sprintf("TrailingSLChangedEvent__Ctid:_PositionId:%d__OrderId:%d",
			e.GetPositionId(),
			e.GetOrderId(),
		)),
		APIEventKey(fmt.Sprintf("TrailingSLChangedEvent__Ctid:_PositionId:%d__OrderId:",
			e.GetPositionId(),
		)),
		APIEventKey(fmt.Sprintf("TrailingSLChangedEvent__Ctid:_PositionId:__OrderId:%d",
			e.GetOrderId(),
		)),
	}
}
func (e *ProtoOATrailingSLChangedEvent) ToDistributableEvent() DistributableEvent {
	return e
}

func (*ProtoOASymbolChangedEvent) IsListenableAPIEvent() {}
func (e *ProtoOASymbolChangedEvent) BuildKeys() []APIEventKey {
	return []APIEventKey{
		APIEventKey(fmt.Sprintf("SymbolChangedEvent__Ctid:%d",
			e.GetCtidTraderAccountId(),
		)),
	}
}
func (e *ProtoOASymbolChangedEvent) ToDistributableEvent() DistributableEvent {
	return e
}

func (*ProtoOATraderUpdatedEvent) IsListenableAPIEvent() {}
func (e *ProtoOATraderUpdatedEvent) BuildKeys() []APIEventKey {
	return []APIEventKey{
		APIEventKey(fmt.Sprintf("TraderUpdatedEvent__Ctid:%d",
			e.GetCtidTraderAccountId(),
		)),
	}
}
func (e *ProtoOATraderUpdatedEvent) ToDistributableEvent() DistributableEvent {
	return e
}

func (*ProtoOAExecutionEvent) IsListenableAPIEvent() {}
func (e *ProtoOAExecutionEvent) BuildKeys() []APIEventKey {
	return []APIEventKey{
		APIEventKey(fmt.Sprintf("ExecutionEvent__Ctid:%d__PositionId:%d__OrderId:%d__DealId:%d",
			e.GetCtidTraderAccountId(),
			e.GetPosition().GetPositionId(),
			e.GetOrder().GetOrderId(),
			e.GetDeal().GetDealId(),
		)),
		APIEventKey(fmt.Sprintf("ExecutionEvent__Ctid:%d__PositionId:%d__OrderId:%d__DealId:",
			e.GetCtidTraderAccountId(),
			e.GetPosition().GetPositionId(),
			e.GetOrder().GetOrderId(),
		)),
		APIEventKey(fmt.Sprintf("ExecutionEvent__Ctid:%d__PositionId:%d__OrderId:__DealId:%d",
			e.GetCtidTraderAccountId(),
			e.GetPosition().GetPositionId(),
			e.GetDeal().GetDealId(),
		)),
		APIEventKey(fmt.Sprintf("ExecutionEvent__Ctid:%d__PositionId:%d__OrderId:__DealId:",
			e.GetCtidTraderAccountId(),
			e.GetPosition().GetPositionId(),
		)),
		APIEventKey(fmt.Sprintf("ExecutionEvent__Ctid:%d__PositionId:__OrderId:%d__DealId:%d",
			e.GetCtidTraderAccountId(),
			e.GetOrder().GetOrderId(),
			e.GetDeal().GetDealId(),
		)),
		APIEventKey(fmt.Sprintf("ExecutionEvent__Ctid:%d__PositionId:__OrderId:%d__DealId:",
			e.GetCtidTraderAccountId(),
			e.GetOrder().GetOrderId(),
		)),
		APIEventKey(fmt.Sprintf("ExecutionEvent__Ctid:%d__PositionId:__OrderId:__DealId:%d",
			e.GetCtidTraderAccountId(),
			e.GetDeal().GetDealId(),
		)),
		APIEventKey(fmt.Sprintf("ExecutionEvent__Ctid:%d__PositionId:__OrderId:__DealId:",
			e.GetCtidTraderAccountId(),
		)),
		APIEventKey(fmt.Sprintf("ExecutionEvent__Ctid:__PositionId:%d__OrderId:%d__DealId:%d",
			e.GetPosition().GetPositionId(),
			e.GetOrder().GetOrderId(),
			e.GetDeal().GetDealId(),
		)),
		APIEventKey(fmt.Sprintf("ExecutionEvent__Ctid:__PositionId:%d__OrderId:%d__DealId:",
			e.GetPosition().GetPositionId(),
			e.GetOrder().GetOrderId(),
		)),
		APIEventKey(fmt.Sprintf("ExecutionEvent__Ctid:__PositionId:%d__OrderId:__DealId:%d",
			e.GetPosition().GetPositionId(),
			e.GetDeal().GetDealId(),
		)),
		APIEventKey(fmt.Sprintf("ExecutionEvent__Ctid:__PositionId:%d__OrderId:__DealId:",
			e.GetPosition().GetPositionId(),
		)),
		APIEventKey(fmt.Sprintf("ExecutionEvent__Ctid:__PositionId:__OrderId:%d__DealId:%d",
			e.GetOrder().GetOrderId(),
			e.GetDeal().GetDealId(),
		)),
		APIEventKey(fmt.Sprintf("ExecutionEvent__Ctid:__PositionId:__OrderId:%d__DealId:",
			e.GetOrder().GetOrderId(),
		)),
		APIEventKey(fmt.Sprintf("ExecutionEvent__Ctid:__PositionId:__OrderId:__DealId:%d",
			e.GetDeal().GetDealId(),
		)),
	}
}
func (e *ProtoOAExecutionEvent) ToDistributableEvent() DistributableEvent {
	return e
}

func (*ProtoOAOrderErrorEvent) IsListenableAPIEvent() {}
func (e *ProtoOAOrderErrorEvent) BuildKeys() []APIEventKey {
	return []APIEventKey{
		APIEventKey(fmt.Sprintf("OrderErrorEvent__Ctid:%d__PositionId:%d__OrderId:%d",
			e.GetCtidTraderAccountId(),
			e.GetPositionId(),
			e.GetOrderId(),
		)),
		APIEventKey(fmt.Sprintf("OrderErrorEvent__Ctid:%d__PositionId:%d__OrderId:",
			e.GetCtidTraderAccountId(),
			e.GetPositionId(),
		)),
		APIEventKey(fmt.Sprintf("OrderErrorEvent__Ctid:%d__PositionId:__OrderId:%d",
			e.GetCtidTraderAccountId(),
			e.GetOrderId(),
		)),
		APIEventKey(fmt.Sprintf("OrderErrorEvent__Ctid:%d__PositionId:__OrderId:",
			e.GetCtidTraderAccountId(),
		)),
		APIEventKey(fmt.Sprintf("OrderErrorEvent__Ctid:__PositionId:%d__OrderId:%d",
			e.GetPositionId(),
			e.GetOrderId(),
		)),
		APIEventKey(fmt.Sprintf("OrderErrorEvent__Ctid:__PositionId:%d__OrderId:",
			e.GetPositionId(),
		)),
		APIEventKey(fmt.Sprintf("OrderErrorEvent__Ctid:__PositionId:__OrderId:%d",
			e.GetOrderId(),
		)),
	}
}
func (e *ProtoOAOrderErrorEvent) ToDistributableEvent() DistributableEvent {
	return e
}

func (*ProtoOAMarginChangedEvent) IsListenableAPIEvent() {}
func (e *ProtoOAMarginChangedEvent) BuildKeys() []APIEventKey {
	return []APIEventKey{
		APIEventKey(fmt.Sprintf("MarginChangedEvent__Ctid:%d__PositionId:%d",
			e.GetCtidTraderAccountId(),
			e.GetPositionId(),
		)),
		APIEventKey(fmt.Sprintf("MarginChangedEvent__Ctid:%d__PositionId:",
			e.GetCtidTraderAccountId(),
		)),
		APIEventKey(fmt.Sprintf("MarginChangedEvent__Ctid:__PositionId:%d",
			e.GetPositionId(),
		)),
	}
}
func (e *ProtoOAMarginChangedEvent) ToDistributableEvent() DistributableEvent {
	return e
}

func (*ProtoOAAccountsTokenInvalidatedEvent) IsListenableAPIEvent() {}
func (*ProtoOAAccountsTokenInvalidatedEvent) ToDistributableEvent() DistributableEvent {
	return nil
}

func (*ProtoOAClientDisconnectEvent) IsListenableAPIEvent()  {}
func (*ProtoOAAccountDisconnectEvent) IsListenableAPIEvent() {}
func (e *ProtoOAAccountDisconnectEvent) BuildKeys() []APIEventKey {
	return []APIEventKey{
		APIEventKey(fmt.Sprintf("AccountDisconnectEvent__Ctid:%d",
			e.GetCtidTraderAccountId(),
		)),
	}
}
func (e *ProtoOAAccountDisconnectEvent) ToDistributableEvent() DistributableEvent {
	return e
}

func (*ProtoOAMarginCallUpdateEvent) IsListenableAPIEvent() {}
func (e *ProtoOAMarginCallUpdateEvent) BuildKeys() []APIEventKey {
	return []APIEventKey{
		APIEventKey(fmt.Sprintf("MarginCallUpdateEvent__Ctid:%d",
			e.GetCtidTraderAccountId(),
		)),
	}
}
func (e *ProtoOAMarginCallUpdateEvent) ToDistributableEvent() DistributableEvent {
	return e
}

func (*ProtoOAMarginCallTriggerEvent) IsListenableAPIEvent() {}
func (e *ProtoOAMarginCallTriggerEvent) BuildKeys() []APIEventKey {
	return []APIEventKey{
		APIEventKey(fmt.Sprintf("MarginCallTriggerEvent__Ctid:%d",
			e.GetCtidTraderAccountId(),
		)),
	}
}

func (e *ProtoOAMarginCallTriggerEvent) ToDistributableEvent() DistributableEvent {
	return e
}
