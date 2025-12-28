package messages

/*
Subscribable event assertions
*/
/**/

func (*ProtoOASpotEvent) IsSubscribableEvent()  {}
func (*ProtoOASpotEvent) IsListenableAPIEvent() {}

func (*ProtoOADepthEvent) IsSubscribableEvent()  {}
func (*ProtoOADepthEvent) IsListenableAPIEvent() {}

/*
Listenable event assertions
*/
/**/

func (*ProtoOATrailingSLChangedEvent) IsListenableAPIEvent() {}

func (*ProtoOASymbolChangedEvent) IsListenableAPIEvent() {}

func (*ProtoOATraderUpdatedEvent) IsListenableAPIEvent() {}

func (*ProtoOAExecutionEvent) IsListenableAPIEvent() {}

func (*ProtoOAOrderErrorEvent) IsListenableAPIEvent() {}

func (*ProtoOAMarginChangedEvent) IsListenableAPIEvent() {}

func (*ProtoOAAccountsTokenInvalidatedEvent) IsListenableAPIEvent() {}

func (*ProtoOAClientDisconnectEvent) IsListenableAPIEvent() {}

func (*ProtoOAAccountDisconnectEvent) IsListenableAPIEvent() {}

func (*ProtoOAMarginCallUpdateEvent) IsListenableAPIEvent() {}

func (*ProtoOAMarginCallTriggerEvent) IsListenableAPIEvent() {}
