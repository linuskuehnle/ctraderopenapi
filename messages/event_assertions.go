package messages

/*
Subscribable event assertions
*/
/**/

func (*ProtoOASpotEvent) IsSubscribableEvent() {}
func (*ProtoOASpotEvent) IsListenableEvent()   {}

func (*ProtoOADepthEvent) IsSubscribableEvent() {}
func (*ProtoOADepthEvent) IsListenableEvent()   {}

/*
Listenable event assertions
*/
/**/

func (*ProtoOATrailingSLChangedEvent) IsListenableEvent() {}

func (*ProtoOASymbolChangedEvent) IsListenableEvent() {}

func (*ProtoOATraderUpdatedEvent) IsListenableEvent() {}

func (*ProtoOAExecutionEvent) IsListenableEvent() {}

func (*ProtoOAOrderErrorEvent) IsListenableEvent() {}

func (*ProtoOAMarginChangedEvent) IsListenableEvent() {}

func (*ProtoOAAccountsTokenInvalidatedEvent) IsListenableEvent() {}

func (*ProtoOAClientDisconnectEvent) IsListenableEvent() {}

func (*ProtoOAAccountDisconnectEvent) IsListenableEvent() {}

func (*ProtoOAMarginCallUpdateEvent) IsListenableEvent() {}

func (*ProtoOAMarginCallTriggerEvent) IsListenableEvent() {}
