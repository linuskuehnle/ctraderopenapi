package datatypes

import (
	"github.com/linuskuehnle/ctraderopenapi/internal/messages"
)

// Mapped responses to requests
var resTypeByReqType = map[messages.ProtoOAPayloadType]messages.ProtoOAPayloadType{
	messages.ProtoOAPayloadType_PROTO_OA_APPLICATION_AUTH_REQ:             messages.ProtoOAPayloadType_PROTO_OA_APPLICATION_AUTH_RES,
	messages.ProtoOAPayloadType_PROTO_OA_ACCOUNT_AUTH_REQ:                 messages.ProtoOAPayloadType_PROTO_OA_ACCOUNT_AUTH_RES,
	messages.ProtoOAPayloadType_PROTO_OA_VERSION_REQ:                      messages.ProtoOAPayloadType_PROTO_OA_VERSION_RES,
	messages.ProtoOAPayloadType_PROTO_OA_NEW_ORDER_REQ:                    0, // no response defined
	messages.ProtoOAPayloadType_PROTO_OA_CANCEL_ORDER_REQ:                 0, // no response defined
	messages.ProtoOAPayloadType_PROTO_OA_AMEND_ORDER_REQ:                  0, // no response defined
	messages.ProtoOAPayloadType_PROTO_OA_AMEND_POSITION_SLTP_REQ:          0, // no response defined
	messages.ProtoOAPayloadType_PROTO_OA_CLOSE_POSITION_REQ:               0, // no response defined
	messages.ProtoOAPayloadType_PROTO_OA_ASSET_LIST_REQ:                   messages.ProtoOAPayloadType_PROTO_OA_ASSET_LIST_RES,
	messages.ProtoOAPayloadType_PROTO_OA_SYMBOLS_LIST_REQ:                 messages.ProtoOAPayloadType_PROTO_OA_SYMBOLS_LIST_RES,
	messages.ProtoOAPayloadType_PROTO_OA_SYMBOL_BY_ID_REQ:                 messages.ProtoOAPayloadType_PROTO_OA_SYMBOL_BY_ID_RES,
	messages.ProtoOAPayloadType_PROTO_OA_SYMBOLS_FOR_CONVERSION_REQ:       messages.ProtoOAPayloadType_PROTO_OA_SYMBOLS_FOR_CONVERSION_RES,
	messages.ProtoOAPayloadType_PROTO_OA_TRADER_REQ:                       messages.ProtoOAPayloadType_PROTO_OA_TRADER_RES,
	messages.ProtoOAPayloadType_PROTO_OA_RECONCILE_REQ:                    messages.ProtoOAPayloadType_PROTO_OA_RECONCILE_RES,
	messages.ProtoOAPayloadType_PROTO_OA_SUBSCRIBE_SPOTS_REQ:              messages.ProtoOAPayloadType_PROTO_OA_SUBSCRIBE_SPOTS_RES,
	messages.ProtoOAPayloadType_PROTO_OA_UNSUBSCRIBE_SPOTS_REQ:            messages.ProtoOAPayloadType_PROTO_OA_UNSUBSCRIBE_SPOTS_RES,
	messages.ProtoOAPayloadType_PROTO_OA_DEAL_LIST_REQ:                    messages.ProtoOAPayloadType_PROTO_OA_DEAL_LIST_RES,
	messages.ProtoOAPayloadType_PROTO_OA_SUBSCRIBE_LIVE_TRENDBAR_REQ:      messages.ProtoOAPayloadType_PROTO_OA_SUBSCRIBE_LIVE_TRENDBAR_RES,
	messages.ProtoOAPayloadType_PROTO_OA_UNSUBSCRIBE_LIVE_TRENDBAR_REQ:    messages.ProtoOAPayloadType_PROTO_OA_UNSUBSCRIBE_LIVE_TRENDBAR_RES,
	messages.ProtoOAPayloadType_PROTO_OA_GET_TRENDBARS_REQ:                messages.ProtoOAPayloadType_PROTO_OA_GET_TRENDBARS_RES,
	messages.ProtoOAPayloadType_PROTO_OA_EXPECTED_MARGIN_REQ:              messages.ProtoOAPayloadType_PROTO_OA_EXPECTED_MARGIN_RES,
	messages.ProtoOAPayloadType_PROTO_OA_CASH_FLOW_HISTORY_LIST_REQ:       messages.ProtoOAPayloadType_PROTO_OA_CASH_FLOW_HISTORY_LIST_RES,
	messages.ProtoOAPayloadType_PROTO_OA_GET_TICKDATA_REQ:                 messages.ProtoOAPayloadType_PROTO_OA_GET_TICKDATA_RES,
	messages.ProtoOAPayloadType_PROTO_OA_GET_ACCOUNTS_BY_ACCESS_TOKEN_REQ: messages.ProtoOAPayloadType_PROTO_OA_GET_ACCOUNTS_BY_ACCESS_TOKEN_RES,
	messages.ProtoOAPayloadType_PROTO_OA_GET_CTID_PROFILE_BY_TOKEN_REQ:    messages.ProtoOAPayloadType_PROTO_OA_GET_CTID_PROFILE_BY_TOKEN_RES,
	messages.ProtoOAPayloadType_PROTO_OA_ASSET_CLASS_LIST_REQ:             messages.ProtoOAPayloadType_PROTO_OA_ASSET_CLASS_LIST_RES,
	messages.ProtoOAPayloadType_PROTO_OA_SUBSCRIBE_DEPTH_QUOTES_REQ:       messages.ProtoOAPayloadType_PROTO_OA_SUBSCRIBE_DEPTH_QUOTES_RES,
	messages.ProtoOAPayloadType_PROTO_OA_UNSUBSCRIBE_DEPTH_QUOTES_REQ:     messages.ProtoOAPayloadType_PROTO_OA_UNSUBSCRIBE_DEPTH_QUOTES_RES,
	messages.ProtoOAPayloadType_PROTO_OA_SYMBOL_CATEGORY_REQ:              messages.ProtoOAPayloadType_PROTO_OA_SYMBOL_CATEGORY_RES,
	messages.ProtoOAPayloadType_PROTO_OA_ACCOUNT_LOGOUT_REQ:               messages.ProtoOAPayloadType_PROTO_OA_ACCOUNT_LOGOUT_RES,
	messages.ProtoOAPayloadType_PROTO_OA_MARGIN_CALL_LIST_REQ:             messages.ProtoOAPayloadType_PROTO_OA_MARGIN_CALL_LIST_RES,
	messages.ProtoOAPayloadType_PROTO_OA_MARGIN_CALL_UPDATE_REQ:           messages.ProtoOAPayloadType_PROTO_OA_MARGIN_CALL_UPDATE_RES,
	messages.ProtoOAPayloadType_PROTO_OA_REFRESH_TOKEN_REQ:                messages.ProtoOAPayloadType_PROTO_OA_REFRESH_TOKEN_RES,
	messages.ProtoOAPayloadType_PROTO_OA_ORDER_LIST_REQ:                   messages.ProtoOAPayloadType_PROTO_OA_ORDER_LIST_RES,
	messages.ProtoOAPayloadType_PROTO_OA_GET_DYNAMIC_LEVERAGE_REQ:         messages.ProtoOAPayloadType_PROTO_OA_GET_DYNAMIC_LEVERAGE_RES,
	messages.ProtoOAPayloadType_PROTO_OA_DEAL_LIST_BY_POSITION_ID_REQ:     messages.ProtoOAPayloadType_PROTO_OA_DEAL_LIST_BY_POSITION_ID_RES,
	messages.ProtoOAPayloadType_PROTO_OA_ORDER_DETAILS_REQ:                messages.ProtoOAPayloadType_PROTO_OA_ORDER_DETAILS_RES,
	messages.ProtoOAPayloadType_PROTO_OA_ORDER_LIST_BY_POSITION_ID_REQ:    messages.ProtoOAPayloadType_PROTO_OA_ORDER_LIST_BY_POSITION_ID_RES,
	messages.ProtoOAPayloadType_PROTO_OA_DEAL_OFFSET_LIST_REQ:             messages.ProtoOAPayloadType_PROTO_OA_DEAL_OFFSET_LIST_RES,
	messages.ProtoOAPayloadType_PROTO_OA_GET_POSITION_UNREALIZED_PNL_REQ:  messages.ProtoOAPayloadType_PROTO_OA_GET_POSITION_UNREALIZED_PNL_RES,
}
