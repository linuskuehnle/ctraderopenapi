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
	"time"

	"github.com/linuskuehnle/ctraderopenapi/datatypes"
)

const (
	ConfigDefault_RequestTimeout              = time.Second * 5
	ConfigDefault_QueueBufferSize             = 10
	ConfigDefault_TCPMessageBufferSize        = 10
	ConfigDefault_RequestHeapIterationTimeout = datatypes.DefaultRequestHeapIterationTimeout
)

const (
	rateLimitN_Live       = 50
	rateLimitN_Historical = 5

	rateLimitInterval   = time.Second
	rateLimitHitTimeout = time.Millisecond * 5
)

const (
	rateLimitType_Live rateLimitType = iota
	rateLimitType_Historical
)

var rateLimitTypeByReqType = map[ProtoOAPayloadType]rateLimitType{
	PROTO_OA_APPLICATION_AUTH_REQ:             rateLimitType_Live,
	PROTO_OA_ACCOUNT_AUTH_REQ:                 rateLimitType_Live,
	PROTO_OA_VERSION_REQ:                      rateLimitType_Live,
	PROTO_OA_NEW_ORDER_REQ:                    rateLimitType_Live,
	PROTO_OA_CANCEL_ORDER_REQ:                 rateLimitType_Live,
	PROTO_OA_AMEND_ORDER_REQ:                  rateLimitType_Live,
	PROTO_OA_AMEND_POSITION_SLTP_REQ:          rateLimitType_Live,
	PROTO_OA_CLOSE_POSITION_REQ:               rateLimitType_Live,
	PROTO_OA_ASSET_LIST_REQ:                   rateLimitType_Live,
	PROTO_OA_SYMBOLS_LIST_REQ:                 rateLimitType_Live,
	PROTO_OA_SYMBOL_BY_ID_REQ:                 rateLimitType_Live,
	PROTO_OA_SYMBOLS_FOR_CONVERSION_REQ:       rateLimitType_Live,
	PROTO_OA_TRADER_REQ:                       rateLimitType_Live,
	PROTO_OA_RECONCILE_REQ:                    rateLimitType_Live,
	PROTO_OA_SUBSCRIBE_SPOTS_REQ:              rateLimitType_Live,
	PROTO_OA_UNSUBSCRIBE_SPOTS_REQ:            rateLimitType_Live,
	PROTO_OA_DEAL_LIST_REQ:                    rateLimitType_Historical,
	PROTO_OA_SUBSCRIBE_LIVE_TRENDBAR_REQ:      rateLimitType_Live,
	PROTO_OA_UNSUBSCRIBE_LIVE_TRENDBAR_REQ:    rateLimitType_Live,
	PROTO_OA_GET_TRENDBARS_REQ:                rateLimitType_Historical,
	PROTO_OA_EXPECTED_MARGIN_REQ:              rateLimitType_Live,
	PROTO_OA_CASH_FLOW_HISTORY_LIST_REQ:       rateLimitType_Historical,
	PROTO_OA_GET_TICKDATA_REQ:                 rateLimitType_Historical,
	PROTO_OA_GET_ACCOUNTS_BY_ACCESS_TOKEN_REQ: rateLimitType_Live,
	PROTO_OA_GET_CTID_PROFILE_BY_TOKEN_REQ:    rateLimitType_Live,
	PROTO_OA_ASSET_CLASS_LIST_REQ:             rateLimitType_Live,
	PROTO_OA_SUBSCRIBE_DEPTH_QUOTES_REQ:       rateLimitType_Live,
	PROTO_OA_UNSUBSCRIBE_DEPTH_QUOTES_REQ:     rateLimitType_Live,
	PROTO_OA_SYMBOL_CATEGORY_REQ:              rateLimitType_Live,
	PROTO_OA_ACCOUNT_LOGOUT_REQ:               rateLimitType_Live,
	PROTO_OA_MARGIN_CALL_LIST_REQ:             rateLimitType_Live,
	PROTO_OA_MARGIN_CALL_UPDATE_REQ:           rateLimitType_Live,
	PROTO_OA_REFRESH_TOKEN_REQ:                rateLimitType_Live,
	PROTO_OA_ORDER_LIST_REQ:                   rateLimitType_Live,
	PROTO_OA_GET_DYNAMIC_LEVERAGE_REQ:         rateLimitType_Live,
	PROTO_OA_DEAL_LIST_BY_POSITION_ID_REQ:     rateLimitType_Live,
	PROTO_OA_ORDER_DETAILS_REQ:                rateLimitType_Live,
	PROTO_OA_ORDER_LIST_BY_POSITION_ID_REQ:    rateLimitType_Live,
	PROTO_OA_DEAL_OFFSET_LIST_REQ:             rateLimitType_Live,
	PROTO_OA_GET_POSITION_UNREALIZED_PNL_REQ:  rateLimitType_Live,
}

const (
	// EndpointAddress represents the cTrader OpenAPI server address.

	// Demo environment endpoint address
	EndpointAddress_Demo endpointAddress = "demo.ctraderapi.com:5035"
	// Live environment endpoint address
	EndpointAddress_Live endpointAddress = "live.ctraderapi.com:5035"
)

const (
	// Demo environment is used for cTrader Demo accounts
	Environment_Demo Environment = iota
	// Live environment is used for cTrader Live accounts
	Environment_Live Environment = iota
)

const (
	// Heartbeat_Interval_Seconds is the interval at which heartbeat messages are sent to the server.
	Heartbeat_Timeout_Seconds = 9
)

const (
	resErrorCode_appAlreadyAuthenticated   = "ALREADY_LOGGED_IN"
	resErrorCode_apiEventAlreadySubscribed = "ALREADY_SUBSCRIBED"
	resErrorCode_serverSideRateLimitHit    = "REQUEST_FREQUENCY_EXCEEDED"
)

// Mapped responses to requests
var resTypeByReqType = map[ProtoOAPayloadType]ProtoOAPayloadType{
	PROTO_OA_APPLICATION_AUTH_REQ:             PROTO_OA_APPLICATION_AUTH_RES,
	PROTO_OA_ACCOUNT_AUTH_REQ:                 PROTO_OA_ACCOUNT_AUTH_RES,
	PROTO_OA_VERSION_REQ:                      PROTO_OA_VERSION_RES,
	PROTO_OA_NEW_ORDER_REQ:                    0, // no response defined
	PROTO_OA_CANCEL_ORDER_REQ:                 0, // no response defined
	PROTO_OA_AMEND_ORDER_REQ:                  0, // no response defined
	PROTO_OA_AMEND_POSITION_SLTP_REQ:          0, // no response defined
	PROTO_OA_CLOSE_POSITION_REQ:               0, // no response defined
	PROTO_OA_ASSET_LIST_REQ:                   PROTO_OA_ASSET_LIST_RES,
	PROTO_OA_SYMBOLS_LIST_REQ:                 PROTO_OA_SYMBOLS_LIST_RES,
	PROTO_OA_SYMBOL_BY_ID_REQ:                 PROTO_OA_SYMBOL_BY_ID_RES,
	PROTO_OA_SYMBOLS_FOR_CONVERSION_REQ:       PROTO_OA_SYMBOLS_FOR_CONVERSION_RES,
	PROTO_OA_TRADER_REQ:                       PROTO_OA_TRADER_RES,
	PROTO_OA_RECONCILE_REQ:                    PROTO_OA_RECONCILE_RES,
	PROTO_OA_SUBSCRIBE_SPOTS_REQ:              PROTO_OA_SUBSCRIBE_SPOTS_RES,
	PROTO_OA_UNSUBSCRIBE_SPOTS_REQ:            PROTO_OA_UNSUBSCRIBE_SPOTS_RES,
	PROTO_OA_DEAL_LIST_REQ:                    PROTO_OA_DEAL_LIST_RES,
	PROTO_OA_SUBSCRIBE_LIVE_TRENDBAR_REQ:      PROTO_OA_SUBSCRIBE_LIVE_TRENDBAR_RES,
	PROTO_OA_UNSUBSCRIBE_LIVE_TRENDBAR_REQ:    PROTO_OA_UNSUBSCRIBE_LIVE_TRENDBAR_RES,
	PROTO_OA_GET_TRENDBARS_REQ:                PROTO_OA_GET_TRENDBARS_RES,
	PROTO_OA_EXPECTED_MARGIN_REQ:              PROTO_OA_EXPECTED_MARGIN_RES,
	PROTO_OA_CASH_FLOW_HISTORY_LIST_REQ:       PROTO_OA_CASH_FLOW_HISTORY_LIST_RES,
	PROTO_OA_GET_TICKDATA_REQ:                 PROTO_OA_GET_TICKDATA_RES,
	PROTO_OA_GET_ACCOUNTS_BY_ACCESS_TOKEN_REQ: PROTO_OA_GET_ACCOUNTS_BY_ACCESS_TOKEN_RES,
	PROTO_OA_GET_CTID_PROFILE_BY_TOKEN_REQ:    PROTO_OA_GET_CTID_PROFILE_BY_TOKEN_RES,
	PROTO_OA_ASSET_CLASS_LIST_REQ:             PROTO_OA_ASSET_CLASS_LIST_RES,
	PROTO_OA_SUBSCRIBE_DEPTH_QUOTES_REQ:       PROTO_OA_SUBSCRIBE_DEPTH_QUOTES_RES,
	PROTO_OA_UNSUBSCRIBE_DEPTH_QUOTES_REQ:     PROTO_OA_UNSUBSCRIBE_DEPTH_QUOTES_RES,
	PROTO_OA_SYMBOL_CATEGORY_REQ:              PROTO_OA_SYMBOL_CATEGORY_RES,
	PROTO_OA_ACCOUNT_LOGOUT_REQ:               PROTO_OA_ACCOUNT_LOGOUT_RES,
	PROTO_OA_MARGIN_CALL_LIST_REQ:             PROTO_OA_MARGIN_CALL_LIST_RES,
	PROTO_OA_MARGIN_CALL_UPDATE_REQ:           PROTO_OA_MARGIN_CALL_UPDATE_RES,
	PROTO_OA_REFRESH_TOKEN_REQ:                PROTO_OA_REFRESH_TOKEN_RES,
	PROTO_OA_ORDER_LIST_REQ:                   PROTO_OA_ORDER_LIST_RES,
	PROTO_OA_GET_DYNAMIC_LEVERAGE_REQ:         PROTO_OA_GET_DYNAMIC_LEVERAGE_RES,
	PROTO_OA_DEAL_LIST_BY_POSITION_ID_REQ:     PROTO_OA_DEAL_LIST_BY_POSITION_ID_RES,
	PROTO_OA_ORDER_DETAILS_REQ:                PROTO_OA_ORDER_DETAILS_RES,
	PROTO_OA_ORDER_LIST_BY_POSITION_ID_REQ:    PROTO_OA_ORDER_LIST_BY_POSITION_ID_RES,
	PROTO_OA_DEAL_OFFSET_LIST_REQ:             PROTO_OA_DEAL_OFFSET_LIST_RES,
	PROTO_OA_GET_POSITION_UNREALIZED_PNL_REQ:  PROTO_OA_GET_POSITION_UNREALIZED_PNL_RES,
}

var isAPIEvent = map[ProtoOAPayloadType]bool{
	PROTO_OA_SPOT_EVENT:                       true,
	PROTO_OA_DEPTH_EVENT:                      true,
	PROTO_OA_TRAILING_SL_CHANGED_EVENT:        true,
	PROTO_OA_SYMBOL_CHANGED_EVENT:             true,
	PROTO_OA_TRADER_UPDATE_EVENT:              true,
	PROTO_OA_EXECUTION_EVENT:                  true,
	PROTO_OA_ORDER_ERROR_EVENT:                true,
	PROTO_OA_MARGIN_CHANGED_EVENT:             true,
	PROTO_OA_ACCOUNTS_TOKEN_INVALIDATED_EVENT: true,
	PROTO_OA_CLIENT_DISCONNECT_EVENT:          true,
	PROTO_OA_ACCOUNT_DISCONNECT_EVENT:         true,
	PROTO_OA_MARGIN_CALL_UPDATE_EVENT:         true,
	PROTO_OA_MARGIN_CALL_TRIGGER_EVENT:        true,
}

var hasHookForAPIEvent = map[ProtoOAPayloadType]bool{
	PROTO_OA_ACCOUNT_DISCONNECT_EVENT: true,
}

const (
	// Subscribable events
	APIEventType_Spots apiEventType = iota
	APIEventType_LiveTrendbars
	APIEventType_DepthQuotes

	// Listenable events
	APIEventType_TrailingSLChanged
	APIEventType_SymbolChanged
	APIEventType_TraderUpdated
	APIEventType_Execution
	APIEventType_OrderError
	APIEventType_MarginChanged
	APIEventType_AccountsTokenInvalidated
	APIEventType_ClientDisconnect
	APIEventType_AccountDisconnect
	APIEventType_MarginCallUpdate
	APIEventType_MarginCallTrigger
)

const (
	// API Client events
	ClientEventType_FatalErrorEvent clientEventType = iota
	ClientEventType_ConnectionLossEvent
	ClientEventType_ReconnectSuccessEvent
	ClientEventType_ReconnectFailEvent
)
