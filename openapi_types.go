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
	"github.com/linuskuehnle/ctraderopenapi/internal/messages"
)

/*
-------------------
CommonMessages.pb.go
-------------------
*/

// ProtoMessage left out because it is only used internally

// ProtoErrorRes left out because it is only used internally

// ProtoHeartbeatEvent left out because it is only used internally

/*
-------------------
CommonModelMessages.pb.go
-------------------
*/

// ProtoPayloadType left out because it is only used internally

type ProtoErrorCode = messages.ProtoErrorCode

const (
	ProtoErrorCode_UNKNOWN_ERROR           ProtoErrorCode = 1  // Generic error.
	ProtoErrorCode_UNSUPPORTED_MESSAGE     ProtoErrorCode = 2  // Message is not supported. Wrong message.
	ProtoErrorCode_INVALID_REQUEST         ProtoErrorCode = 3  // Generic error.  Usually used when input value is not correct.
	ProtoErrorCode_TIMEOUT_ERROR           ProtoErrorCode = 5  // Deal execution is reached timeout and rejected.
	ProtoErrorCode_ENTITY_NOT_FOUND        ProtoErrorCode = 6  // Generic error for requests by id.
	ProtoErrorCode_CANT_ROUTE_REQUEST      ProtoErrorCode = 7  // Connection to Server is lost or not supported.
	ProtoErrorCode_FRAME_TOO_LONG          ProtoErrorCode = 8  // Message is too large.
	ProtoErrorCode_MARKET_CLOSED           ProtoErrorCode = 9  // Market is closed.
	ProtoErrorCode_CONCURRENT_MODIFICATION ProtoErrorCode = 10 // Order is blocked (e.g. under execution) and change cannot be applied.
	ProtoErrorCode_BLOCKED_PAYLOAD_TYPE    ProtoErrorCode = 11 // Message is blocked by server or rate limit is reached.
)

var (
	ProtoErrorCodeToString = map[ProtoErrorCode]string{
		1:  "UNKNOWN_ERROR",
		2:  "UNSUPPORTED_MESSAGE",
		3:  "INVALID_REQUEST",
		5:  "TIMEOUT_ERROR",
		6:  "ENTITY_NOT_FOUND",
		7:  "CANT_ROUTE_REQUEST",
		8:  "FRAME_TOO_LONG",
		9:  "MARKET_CLOSED",
		10: "CONCURRENT_MODIFICATION",
		11: "BLOCKED_PAYLOAD_TYPE",
	}
	ProtoErrorCodeByString = map[string]ProtoErrorCode{
		"UNKNOWN_ERROR":           1,
		"UNSUPPORTED_MESSAGE":     2,
		"INVALID_REQUEST":         3,
		"TIMEOUT_ERROR":           5,
		"ENTITY_NOT_FOUND":        6,
		"CANT_ROUTE_REQUEST":      7,
		"FRAME_TOO_LONG":          8,
		"MARKET_CLOSED":           9,
		"CONCURRENT_MODIFICATION": 10,
		"BLOCKED_PAYLOAD_TYPE":    11,
	}
)

/*
-------------------
ModelMessages.pb.go
-------------------
*/

/*
Enum for payload type.
*/
type protoOAPayloadType = messages.ProtoOAPayloadType

const (
	proto_OA_APPLICATION_AUTH_REQ protoOAPayloadType = 2100
	// proto_OA_APPLICATION_AUTH_RES             protoOAPayloadType = 2101
	proto_OA_ACCOUNT_AUTH_REQ protoOAPayloadType = 2102
	// proto_OA_ACCOUNT_AUTH_RES                 protoOAPayloadType = 2103
	proto_OA_VERSION_REQ protoOAPayloadType = 2104
	// proto_OA_VERSION_RES                      protoOAPayloadType = 2105
	proto_OA_NEW_ORDER_REQ             protoOAPayloadType = 2106
	proto_OA_TRAILING_SL_CHANGED_EVENT protoOAPayloadType = 2107
	proto_OA_CANCEL_ORDER_REQ          protoOAPayloadType = 2108
	proto_OA_AMEND_ORDER_REQ           protoOAPayloadType = 2109
	proto_OA_AMEND_POSITION_SLTP_REQ   protoOAPayloadType = 2110
	proto_OA_CLOSE_POSITION_REQ        protoOAPayloadType = 2111
	proto_OA_ASSET_LIST_REQ            protoOAPayloadType = 2112
	// proto_OA_ASSET_LIST_RES                   protoOAPayloadType = 2113
	proto_OA_SYMBOLS_LIST_REQ protoOAPayloadType = 2114
	// proto_OA_SYMBOLS_LIST_RES                 protoOAPayloadType = 2115
	proto_OA_SYMBOL_BY_ID_REQ protoOAPayloadType = 2116
	// proto_OA_SYMBOL_BY_ID_RES                 protoOAPayloadType = 2117
	proto_OA_SYMBOLS_FOR_CONVERSION_REQ protoOAPayloadType = 2118
	// proto_OA_SYMBOLS_FOR_CONVERSION_RES       protoOAPayloadType = 2119
	proto_OA_SYMBOL_CHANGED_EVENT protoOAPayloadType = 2120
	proto_OA_TRADER_REQ           protoOAPayloadType = 2121
	// proto_OA_TRADER_RES                       protoOAPayloadType = 2122
	proto_OA_TRADER_UPDATE_EVENT protoOAPayloadType = 2123
	proto_OA_RECONCILE_REQ       protoOAPayloadType = 2124
	// proto_OA_RECONCILE_RES                    protoOAPayloadType = 2125
	proto_OA_EXECUTION_EVENT     protoOAPayloadType = 2126
	proto_OA_SUBSCRIBE_SPOTS_REQ protoOAPayloadType = 2127
	// proto_OA_SUBSCRIBE_SPOTS_RES              protoOAPayloadType = 2128
	proto_OA_UNSUBSCRIBE_SPOTS_REQ protoOAPayloadType = 2129
	// proto_OA_UNSUBSCRIBE_SPOTS_RES            protoOAPayloadType = 2130
	proto_OA_SPOT_EVENT        protoOAPayloadType = 2131
	proto_OA_ORDER_ERROR_EVENT protoOAPayloadType = 2132
	proto_OA_DEAL_LIST_REQ     protoOAPayloadType = 2133
	// proto_OA_DEAL_LIST_RES                    protoOAPayloadType = 2134
	proto_OA_SUBSCRIBE_LIVE_TRENDBAR_REQ   protoOAPayloadType = 2135
	proto_OA_UNSUBSCRIBE_LIVE_TRENDBAR_REQ protoOAPayloadType = 2136
	proto_OA_GET_TRENDBARS_REQ             protoOAPayloadType = 2137
	// proto_OA_GET_TRENDBARS_RES                protoOAPayloadType = 2138
	proto_OA_EXPECTED_MARGIN_REQ protoOAPayloadType = 2139
	// proto_OA_EXPECTED_MARGIN_RES              protoOAPayloadType = 2140
	proto_OA_MARGIN_CHANGED_EVENT       protoOAPayloadType = 2141
	proto_OA_ERROR_RES                  protoOAPayloadType = 2142
	proto_OA_CASH_FLOW_HISTORY_LIST_REQ protoOAPayloadType = 2143
	// proto_OA_CASH_FLOW_HISTORY_LIST_RES       protoOAPayloadType = 2144
	proto_OA_GET_TICKDATA_REQ protoOAPayloadType = 2145
	// proto_OA_GET_TICKDATA_RES                 protoOAPayloadType = 2146
	proto_OA_ACCOUNTS_TOKEN_INVALIDATED_EVENT protoOAPayloadType = 2147
	proto_OA_CLIENT_DISCONNECT_EVENT          protoOAPayloadType = 2148
	proto_OA_GET_ACCOUNTS_BY_ACCESS_TOKEN_REQ protoOAPayloadType = 2149
	// proto_OA_GET_ACCOUNTS_BY_ACCESS_TOKEN_RES protoOAPayloadType = 2150
	proto_OA_GET_CTID_PROFILE_BY_TOKEN_REQ protoOAPayloadType = 2151
	// proto_OA_GET_CTID_PROFILE_BY_TOKEN_RES    protoOAPayloadType = 2152
	proto_OA_ASSET_CLASS_LIST_REQ protoOAPayloadType = 2153
	// proto_OA_ASSET_CLASS_LIST_RES             protoOAPayloadType = 2154
	proto_OA_DEPTH_EVENT                protoOAPayloadType = 2155
	proto_OA_SUBSCRIBE_DEPTH_QUOTES_REQ protoOAPayloadType = 2156
	// proto_OA_SUBSCRIBE_DEPTH_QUOTES_RES       protoOAPayloadType = 2157
	proto_OA_UNSUBSCRIBE_DEPTH_QUOTES_REQ protoOAPayloadType = 2158
	// proto_OA_UNSUBSCRIBE_DEPTH_QUOTES_RES     protoOAPayloadType = 2159
	proto_OA_SYMBOL_CATEGORY_REQ protoOAPayloadType = 2160
	// proto_OA_SYMBOL_CATEGORY_RES              protoOAPayloadType = 2161
	proto_OA_ACCOUNT_LOGOUT_REQ protoOAPayloadType = 2162
	// proto_OA_ACCOUNT_LOGOUT_RES               protoOAPayloadType = 2163
	proto_OA_ACCOUNT_DISCONNECT_EVENT protoOAPayloadType = 2164
	// proto_OA_SUBSCRIBE_LIVE_TRENDBAR_RES      protoOAPayloadType = 2165
	// proto_OA_UNSUBSCRIBE_LIVE_TRENDBAR_RES    protoOAPayloadType = 2166
	proto_OA_MARGIN_CALL_LIST_REQ protoOAPayloadType = 2167
	// proto_OA_MARGIN_CALL_LIST_RES             protoOAPayloadType = 2168
	proto_OA_MARGIN_CALL_UPDATE_REQ protoOAPayloadType = 2169
	// proto_OA_MARGIN_CALL_UPDATE_RES           protoOAPayloadType = 2170
	proto_OA_MARGIN_CALL_UPDATE_EVENT  protoOAPayloadType = 2171
	proto_OA_MARGIN_CALL_TRIGGER_EVENT protoOAPayloadType = 2172
	proto_OA_REFRESH_TOKEN_REQ         protoOAPayloadType = 2173
	// proto_OA_REFRESH_TOKEN_RES                protoOAPayloadType = 2174
	proto_OA_ORDER_LIST_REQ protoOAPayloadType = 2175
	// proto_OA_ORDER_LIST_RES                   protoOAPayloadType = 2176
	proto_OA_GET_DYNAMIC_LEVERAGE_REQ protoOAPayloadType = 2177
	// proto_OA_GET_DYNAMIC_LEVERAGE_RES         protoOAPayloadType = 2178
	proto_OA_DEAL_LIST_BY_POSITION_ID_REQ protoOAPayloadType = 2179
	// proto_OA_DEAL_LIST_BY_POSITION_ID_RES     protoOAPayloadType = 2180
	proto_OA_ORDER_DETAILS_REQ protoOAPayloadType = 2181
	// proto_OA_ORDER_DETAILS_RES                protoOAPayloadType = 2182
	proto_OA_ORDER_LIST_BY_POSITION_ID_REQ protoOAPayloadType = 2183
	// proto_OA_ORDER_LIST_BY_POSITION_ID_RES    protoOAPayloadType = 2184
	proto_OA_DEAL_OFFSET_LIST_REQ protoOAPayloadType = 2185
	// proto_OA_DEAL_OFFSET_LIST_RES             protoOAPayloadType = 2186
	proto_OA_GET_POSITION_UNREALIZED_PNL_REQ protoOAPayloadType = 2187
	// proto_OA_GET_POSITION_UNREALIZED_PNL_RES  protoOAPayloadType = 2188
)

/*
Enum for day of week.
*/
type ProtoOADayOfWeek = messages.ProtoOADayOfWeek

const (
	ProtoOADayOfWeek_NONE      ProtoOADayOfWeek = 0
	ProtoOADayOfWeek_MONDAY    ProtoOADayOfWeek = 1
	ProtoOADayOfWeek_TUESDAY   ProtoOADayOfWeek = 2
	ProtoOADayOfWeek_WEDNESDAY ProtoOADayOfWeek = 3
	ProtoOADayOfWeek_THURSDAY  ProtoOADayOfWeek = 4
	ProtoOADayOfWeek_FRIDAY    ProtoOADayOfWeek = 5
	ProtoOADayOfWeek_SATURDAY  ProtoOADayOfWeek = 6
	ProtoOADayOfWeek_SUNDAY    ProtoOADayOfWeek = 7
)

var (
	ProtoOADayOfWeekToString = map[ProtoOADayOfWeek]string{
		0: "NONE",
		1: "MONDAY",
		2: "TUESDAY",
		3: "WEDNESDAY",
		4: "THURSDAY",
		5: "FRIDAY",
		6: "SATURDAY",
		7: "SUNDAY",
	}
	ProtoOADayOfWeekByString = map[string]ProtoOADayOfWeek{
		"NONE":      0,
		"MONDAY":    1,
		"TUESDAY":   2,
		"WEDNESDAY": 3,
		"THURSDAY":  4,
		"FRIDAY":    5,
		"SATURDAY":  6,
		"SUNDAY":    7,
	}
)

/*
Enum for specifying type of trading commission.
*/
type ProtoOACommissionType = messages.ProtoOACommissionType

const (
	ProtoOACommissionType_USD_PER_MILLION_USD   ProtoOACommissionType = 1 // USD per million USD volume - usually used for FX. Example: 50 USD for 1 mil USD of trading volume.
	ProtoOACommissionType_USD_PER_LOT           ProtoOACommissionType = 2 // USD per 1 lot - usually used for CFDs and futures for commodities, and indices. Example: 15 USD for 1 contract.
	ProtoOACommissionType_PERCENTAGE_OFByString ProtoOACommissionType = 3 // Percentage of trading volume - usually used for Equities. Example: 0.005% of notional trading volume. Multiplied by 100,000.
	ProtoOACommissionType_QUOTE_CCY_PER_LOT     ProtoOACommissionType = 4 // Quote ccy of Symbol per 1 lot - will be used for CFDs and futures for commodities, and indices. Example: 15 EUR for 1 contract of DAX.
)

var (
	ProtoOACommissionTypeToString = map[ProtoOACommissionType]string{
		1: "USD_PER_MILLION_USD",
		2: "USD_PER_LOT",
		3: "PERCENTAGE_OFByString",
		4: "QUOTE_CCY_PER_LOT",
	}
	ProtoOACommissionTypeByString = map[string]ProtoOACommissionType{
		"USD_PER_MILLION_USD":   1,
		"USD_PER_LOT":           2,
		"PERCENTAGE_OFByString": 3,
		"QUOTE_CCY_PER_LOT":     4,
	}
)

/*
Enum for specifying stop loss and take profit distances.
*/
type ProtoOASymbolDistanceType = messages.ProtoOASymbolDistanceType

const (
	ProtoOASymbolDistanceType_SYMBOL_DISTANCE_IN_POINTS     ProtoOASymbolDistanceType = 1
	ProtoOASymbolDistanceType_SYMBOL_DISTANCE_IN_PERCENTAGE ProtoOASymbolDistanceType = 2
)

var (
	ProtoOASymbolDistanceTypeToString = map[ProtoOASymbolDistanceType]string{
		1: "SYMBOL_DISTANCE_IN_POINTS",
		2: "SYMBOL_DISTANCE_IN_PERCENTAGE",
	}
	ProtoOASymbolDistanceTypeByString = map[string]ProtoOASymbolDistanceType{
		"SYMBOL_DISTANCE_IN_POINTS":     1,
		"SYMBOL_DISTANCE_IN_PERCENTAGE": 2,
	}
)

/*
Enum for specifying type of minimum trading commission.
*/
type ProtoOAMinCommissionType = messages.ProtoOAMinCommissionType

const (
	ProtoOAMinCommissionType_CURRENCY       ProtoOAMinCommissionType = 1
	ProtoOAMinCommissionType_QUOTE_CURRENCY ProtoOAMinCommissionType = 2
)

var (
	ProtoOAMinCommissionTypeToString = map[ProtoOAMinCommissionType]string{
		1: "CURRENCY",
		2: "QUOTE_CURRENCY",
	}
	ProtoOAMinCommissionTypeByString = map[string]ProtoOAMinCommissionType{
		"CURRENCY":       1,
		"QUOTE_CURRENCY": 2,
	}
)

/*
Enum for specifying symbol trading mode.
*/
type ProtoOATradingMode = messages.ProtoOATradingMode

const (
	ProtoOATradingMode_ENABLED                             ProtoOATradingMode = 0
	ProtoOATradingMode_DISABLED_WITHOUT_PENDINGS_EXECUTION ProtoOATradingMode = 1
	ProtoOATradingMode_DISABLED_WITH_PENDINGS_EXECUTION    ProtoOATradingMode = 2
	ProtoOATradingMode_CLOSE_ONLY_MODE                     ProtoOATradingMode = 3
)

var (
	ProtoOATradingModeToString = map[ProtoOATradingMode]string{
		0: "ENABLED",
		1: "DISABLED_WITHOUT_PENDINGS_EXECUTION",
		2: "DISABLED_WITH_PENDINGS_EXECUTION",
		3: "CLOSE_ONLY_MODE",
	}
	ProtoOATradingModeByString = map[string]ProtoOATradingMode{
		"ENABLED":                             0,
		"DISABLED_WITHOUT_PENDINGS_EXECUTION": 1,
		"DISABLED_WITH_PENDINGS_EXECUTION":    2,
		"CLOSE_ONLY_MODE":                     3,
	}
)

/*
Enum for specifying SWAP calculation type for symbol.
*/
type ProtoOASwapCalculationType = messages.ProtoOASwapCalculationType

const (
	ProtoOASwapCalculationType_PIPS       ProtoOASwapCalculationType = 0 // Specifies type of SWAP computation as PIPS (0).
	ProtoOASwapCalculationType_PERCENTAGE ProtoOASwapCalculationType = 1 // Specifies type of SWAP computation as PERCENTAGE (1, annual, in percent).
	ProtoOASwapCalculationType_POINTS     ProtoOASwapCalculationType = 2 // Specifies type of SWAP computation as POINTS (2).
)

var (
	ProtoOASwapCalculationTypeToString = map[ProtoOASwapCalculationType]string{
		0: "PIPS",
		1: "PERCENTAGE",
		2: "POINTS",
	}
	ProtoOASwapCalculationTypeByString = map[string]ProtoOASwapCalculationType{
		"PIPS":       0,
		"PERCENTAGE": 1,
		"POINTS":     2,
	}
)

/*
Enum for specifying access right for a trader.
*/
type ProtoOAAccessRights = messages.ProtoOAAccessRights

const (
	ProtoOAAccessRights_FULL_ACCESS ProtoOAAccessRights = 0 // Enable all trading.
	ProtoOAAccessRights_CLOSE_ONLY  ProtoOAAccessRights = 1 // Only closing trading request are enabled.
	ProtoOAAccessRights_NO_TRADING  ProtoOAAccessRights = 2 // View only access.
	ProtoOAAccessRights_NO_LOGIN    ProtoOAAccessRights = 3 // No access.
)

var (
	ProtoOAAccessRightsToString = map[ProtoOAAccessRights]string{
		0: "FULL_ACCESS",
		1: "CLOSE_ONLY",
		2: "NO_TRADING",
		3: "NO_LOGIN",
	}
	ProtoOAAccessRightsByString = map[string]ProtoOAAccessRights{
		"FULL_ACCESS": 0,
		"CLOSE_ONLY":  1,
		"NO_TRADING":  2,
		"NO_LOGIN":    3,
	}
)

/*
Enum for specifying margin calculation type for an account.
*/
type ProtoOATotalMarginCalculationType = messages.ProtoOATotalMarginCalculationType

const (
	ProtoOATotalMarginCalculationType_MAX ProtoOATotalMarginCalculationType = 0
	ProtoOATotalMarginCalculationType_SUM ProtoOATotalMarginCalculationType = 1
	ProtoOATotalMarginCalculationType_NET ProtoOATotalMarginCalculationType = 2
)

var (
	ProtoOATotalMarginCalculationTypeToString = map[ProtoOATotalMarginCalculationType]string{
		0: "MAX",
		1: "SUM",
		2: "NET",
	}
	ProtoOATotalMarginCalculationTypeByString = map[string]ProtoOATotalMarginCalculationType{
		"MAX": 0,
		"SUM": 1,
		"NET": 2,
	}
)

/*
Enum for specifying type of an account.
*/
type ProtoOAAccountType = messages.ProtoOAAccountType

const (
	ProtoOAAccountType_HEDGED         ProtoOAAccountType = 0 // Allows multiple positions on a trading account for a symbol.
	ProtoOAAccountType_NETTED         ProtoOAAccountType = 1 // Only one position per symbol is allowed on a trading account.
	ProtoOAAccountType_SPREAD_BETTING ProtoOAAccountType = 2 // Spread betting type account.
)

var (
	ProtoOAAccountTypeToString = map[ProtoOAAccountType]string{
		0: "HEDGED",
		1: "NETTED",
		2: "SPREAD_BETTING",
	}
	ProtoOAAccountTypeByString = map[string]ProtoOAAccountType{
		"HEDGED":         0,
		"NETTED":         1,
		"SPREAD_BETTING": 2,
	}
)

/*
Enum for position status.
*/
type ProtoOAPositionStatus = messages.ProtoOAPositionStatus

const (
	ProtoOAPositionStatus_POSITION_STATUS_OPEN    ProtoOAPositionStatus = 1
	ProtoOAPositionStatus_POSITION_STATUS_CLOSED  ProtoOAPositionStatus = 2
	ProtoOAPositionStatus_POSITION_STATUS_CREATED ProtoOAPositionStatus = 3 // Empty position is created for pending order.
	ProtoOAPositionStatus_POSITION_STATUS_ERROR   ProtoOAPositionStatus = 4
)

var (
	ProtoOAPositionStatusToString = map[ProtoOAPositionStatus]string{
		1: "POSITION_STATUS_OPEN",
		2: "POSITION_STATUS_CLOSED",
		3: "POSITION_STATUS_CREATED",
		4: "POSITION_STATUS_ERROR",
	}
	ProtoOAPositionStatusByString = map[string]ProtoOAPositionStatus{
		"POSITION_STATUS_OPEN":    1,
		"POSITION_STATUS_CLOSED":  2,
		"POSITION_STATUS_CREATED": 3,
		"POSITION_STATUS_ERROR":   4,
	}
)

/*
Enum for trade side. Used for order, position, deal.
*/
type ProtoOATradeSide = messages.ProtoOATradeSide

const (
	ProtoOATradeSide_BUY  ProtoOATradeSide = 1
	ProtoOATradeSide_SELL ProtoOATradeSide = 2
)

var (
	ProtoOATradeSideToString = map[ProtoOATradeSide]string{
		1: "BUY",
		2: "SELL",
	}
	ProtoOATradeSideByString = map[string]ProtoOATradeSide{
		"BUY":  1,
		"SELL": 2,
	}
)

/*
Enum for order type.
*/
type ProtoOAOrderType = messages.ProtoOAOrderType

const (
	ProtoOAOrderType_MARKET                ProtoOAOrderType = 1
	ProtoOAOrderType_LIMIT                 ProtoOAOrderType = 2
	ProtoOAOrderType_STOP                  ProtoOAOrderType = 3
	ProtoOAOrderType_STOP_LOSS_TAKE_PROFIT ProtoOAOrderType = 4
	ProtoOAOrderType_MARKET_RANGE          ProtoOAOrderType = 5
	ProtoOAOrderType_STOP_LIMIT            ProtoOAOrderType = 6
)

var (
	ProtoOAOrderTypeToString = map[ProtoOAOrderType]string{
		1: "MARKET",
		2: "LIMIT",
		3: "STOP",
		4: "STOP_LOSS_TAKE_PROFIT",
		5: "MARKET_RANGE",
		6: "STOP_LIMIT",
	}
	ProtoOAOrderTypeByString = map[string]ProtoOAOrderType{
		"MARKET":                1,
		"LIMIT":                 2,
		"STOP":                  3,
		"STOP_LOSS_TAKE_PROFIT": 4,
		"MARKET_RANGE":          5,
		"STOP_LIMIT":            6,
	}
)

/*
Enum for order's time in force.
*/
type ProtoOATimeInForce = messages.ProtoOATimeInForce

const (
	ProtoOATimeInForce_GOOD_TILL_DATE      ProtoOATimeInForce = 1
	ProtoOATimeInForce_GOOD_TILL_CANCEL    ProtoOATimeInForce = 2
	ProtoOATimeInForce_IMMEDIATE_OR_CANCEL ProtoOATimeInForce = 3
	ProtoOATimeInForce_FILL_OR_KILL        ProtoOATimeInForce = 4
	ProtoOATimeInForce_MARKET_ON_OPEN      ProtoOATimeInForce = 5
)

var (
	ProtoOATimeInForceToString = map[ProtoOATimeInForce]string{
		1: "GOOD_TILL_DATE",
		2: "GOOD_TILL_CANCEL",
		3: "IMMEDIATE_OR_CANCEL",
		4: "FILL_OR_KILL",
		5: "MARKET_ON_OPEN",
	}
	ProtoOATimeInForceByString = map[string]ProtoOATimeInForce{
		"GOOD_TILL_DATE":      1,
		"GOOD_TILL_CANCEL":    2,
		"IMMEDIATE_OR_CANCEL": 3,
		"FILL_OR_KILL":        4,
		"MARKET_ON_OPEN":      5,
	}
)

/*
Enum for order status.
*/
type ProtoOAOrderStatus = messages.ProtoOAOrderStatus

const (
	ProtoOAOrderStatus_ORDER_STATUS_ACCEPTED  ProtoOAOrderStatus = 1 // Order request validated and accepted for execution.
	ProtoOAOrderStatus_ORDER_STATUS_FILLED    ProtoOAOrderStatus = 2 // Order is fully filled.
	ProtoOAOrderStatus_ORDER_STATUS_REJECTED  ProtoOAOrderStatus = 3 // Order is rejected due to validation.
	ProtoOAOrderStatus_ORDER_STATUS_EXPIRED   ProtoOAOrderStatus = 4 // Order expired. Might be valid for orders with partially filled volume that were expired on LP.
	ProtoOAOrderStatus_ORDER_STATUS_CANCELLED ProtoOAOrderStatus = 5 // Order is cancelled. Might be valid for orders with partially filled volume that were cancelled by LP.
)

var (
	ProtoOAOrderStatusToString = map[ProtoOAOrderStatus]string{
		1: "ORDER_STATUS_ACCEPTED",
		2: "ORDER_STATUS_FILLED",
		3: "ORDER_STATUS_REJECTED",
		4: "ORDER_STATUS_EXPIRED",
		5: "ORDER_STATUS_CANCELLED",
	}
	ProtoOAOrderStatusByString = map[string]ProtoOAOrderStatus{
		"ORDER_STATUS_ACCEPTED":  1,
		"ORDER_STATUS_FILLED":    2,
		"ORDER_STATUS_REJECTED":  3,
		"ORDER_STATUS_EXPIRED":   4,
		"ORDER_STATUS_CANCELLED": 5,
	}
)

/*
Enum for stop order and stop loss triggering method.
*/
type ProtoOAOrderTriggerMethod = messages.ProtoOAOrderTriggerMethod

const (
	ProtoOAOrderTriggerMethod_TRADE           ProtoOAOrderTriggerMethod = 1 // Stop Order: buy is triggered by ask, sell by bid; Stop Loss Order: for buy position is triggered by bid and for sell position by ask.
	ProtoOAOrderTriggerMethod_OPPOSITE        ProtoOAOrderTriggerMethod = 2 // Stop Order: buy is triggered by bid, sell by ask; Stop Loss Order: for buy position is triggered by ask and for sell position by bid.
	ProtoOAOrderTriggerMethod_DOUBLE_TRADE    ProtoOAOrderTriggerMethod = 3 // The same as TRADE, but trigger is checked after the second consecutive tick.
	ProtoOAOrderTriggerMethod_DOUBLE_OPPOSITE ProtoOAOrderTriggerMethod = 4 // The same as OPPOSITE, but trigger is checked after the second consecutive tick.
)

var (
	ProtoOAOrderTriggerMethodToString = map[ProtoOAOrderTriggerMethod]string{
		1: "TRADE",
		2: "OPPOSITE",
		3: "DOUBLE_TRADE",
		4: "DOUBLE_OPPOSITE",
	}
	ProtoOAOrderTriggerMethodByString = map[string]ProtoOAOrderTriggerMethod{
		"TRADE":           1,
		"OPPOSITE":        2,
		"DOUBLE_TRADE":    3,
		"DOUBLE_OPPOSITE": 4,
	}
)

/*
Enum for execution event type.
*/
type ProtoOAExecutionType = messages.ProtoOAExecutionType

const (
	ProtoOAExecutionType_ORDER_ACCEPTED         ProtoOAExecutionType = 2  // Order passed validation.
	ProtoOAExecutionType_ORDER_FILLED           ProtoOAExecutionType = 3  // Order filled.
	ProtoOAExecutionType_ORDER_REPLACED         ProtoOAExecutionType = 4  // Pending order is changed with a new one.
	ProtoOAExecutionType_ORDER_CANCELLED        ProtoOAExecutionType = 5  // Order cancelled.
	ProtoOAExecutionType_ORDER_EXPIRED          ProtoOAExecutionType = 6  // Order with GTD time in force is expired.
	ProtoOAExecutionType_ORDER_REJECTED         ProtoOAExecutionType = 7  // Order is rejected due to validations.
	ProtoOAExecutionType_ORDER_CANCEL_REJECTED  ProtoOAExecutionType = 8  // Cancel order request is rejected.
	ProtoOAExecutionType_SWAP                   ProtoOAExecutionType = 9  // Type related to SWAP execution events.
	ProtoOAExecutionType_DEPOSIT_WITHDRAW       ProtoOAExecutionType = 10 // Type related to event of deposit or withdrawal cash flow operation.
	ProtoOAExecutionType_ORDER_PARTIAL_FILL     ProtoOAExecutionType = 11 // Order is partially filled.
	ProtoOAExecutionType_BONUS_DEPOSIT_WITHDRAW ProtoOAExecutionType = 12 // Type related to event of bonus deposit or bonus withdrawal.
)

var (
	ProtoOAExecutionTypeToString = map[ProtoOAExecutionType]string{
		2:  "ORDER_ACCEPTED",
		3:  "ORDER_FILLED",
		4:  "ORDER_REPLACED",
		5:  "ORDER_CANCELLED",
		6:  "ORDER_EXPIRED",
		7:  "ORDER_REJECTED",
		8:  "ORDER_CANCEL_REJECTED",
		9:  "SWAP",
		10: "DEPOSIT_WITHDRAW",
		11: "ORDER_PARTIAL_FILL",
		12: "BONUS_DEPOSIT_WITHDRAW",
	}
	ProtoOAExecutionTypeByString = map[string]ProtoOAExecutionType{
		"ORDER_ACCEPTED":         2,
		"ORDER_FILLED":           3,
		"ORDER_REPLACED":         4,
		"ORDER_CANCELLED":        5,
		"ORDER_EXPIRED":          6,
		"ORDER_REJECTED":         7,
		"ORDER_CANCEL_REJECTED":  8,
		"SWAP":                   9,
		"DEPOSIT_WITHDRAW":       10,
		"ORDER_PARTIAL_FILL":     11,
		"BONUS_DEPOSIT_WITHDRAW": 12,
	}
)

/*
Enum for bonus operation type.
*/
type ProtoOAChangeBonusType = messages.ProtoOAChangeBonusType

const (
	ProtoOAChangeBonusType_BONUS_DEPOSIT  ProtoOAChangeBonusType = 0
	ProtoOAChangeBonusType_BONUS_WITHDRAW ProtoOAChangeBonusType = 1
)

var (
	ProtoOAChangeBonusTypeToString = map[ProtoOAChangeBonusType]string{
		0: "BONUS_DEPOSIT",
		1: "BONUS_WITHDRAW",
	}
	ProtoOAChangeBonusTypeByString = map[string]ProtoOAChangeBonusType{
		"BONUS_DEPOSIT":  0,
		"BONUS_WITHDRAW": 1,
	}
)

/*
Balance operation entity. Covers all cash movement operations related to account, trading, IB operations, mirroring, etc.
*/
type ProtoOAChangeBalanceType = messages.ProtoOAChangeBalanceType

const (
	ProtoOAChangeBalanceType_BALANCE_DEPOSIT                                  ProtoOAChangeBalanceType = 0  // Cash deposit.
	ProtoOAChangeBalanceType_BALANCE_WITHDRAW                                 ProtoOAChangeBalanceType = 1  // Cash withdrawal.
	ProtoOAChangeBalanceType_BALANCE_DEPOSIT_STRATEGY_COMMISSION_INNER        ProtoOAChangeBalanceType = 3  // Received mirroring commission.
	ProtoOAChangeBalanceType_BALANCE_WITHDRAW_STRATEGY_COMMISSION_INNER       ProtoOAChangeBalanceType = 4  // Paid mirroring commission.
	ProtoOAChangeBalanceType_BALANCE_DEPOSIT_IB_COMMISSIONS                   ProtoOAChangeBalanceType = 5  // For IB account. Commissions paid by trader.
	ProtoOAChangeBalanceType_BALANCE_WITHDRAW_IB_SHARED_PERCENTAGE            ProtoOAChangeBalanceType = 6  // For IB account. Withdrawal of commissions shared with broker.
	ProtoOAChangeBalanceType_BALANCE_DEPOSIT_IB_SHARED_PERCENTAGE_FROM_SUB_IB ProtoOAChangeBalanceType = 7  // For IB account. Commissions paid by sub-ibs.
	ProtoOAChangeBalanceType_BALANCE_DEPOSIT_IB_SHARED_PERCENTAGE_FROM_BROKER ProtoOAChangeBalanceType = 8  // For IB account. Commissions paid by broker.
	ProtoOAChangeBalanceType_BALANCE_DEPOSIT_REBATE                           ProtoOAChangeBalanceType = 9  // Deposit rebate for trading volume for period.
	ProtoOAChangeBalanceType_BALANCE_WITHDRAW_REBATE                          ProtoOAChangeBalanceType = 10 // Withdrawal of rebate.
	ProtoOAChangeBalanceType_BALANCE_DEPOSIT_STRATEGY_COMMISSION_OUTER        ProtoOAChangeBalanceType = 11 // Mirroring commission.
	ProtoOAChangeBalanceType_BALANCE_WITHDRAW_STRATEGY_COMMISSION_OUTER       ProtoOAChangeBalanceType = 12 // Mirroring commission.
	ProtoOAChangeBalanceType_BALANCE_WITHDRAW_BONUS_COMPENSATION              ProtoOAChangeBalanceType = 13 // For IB account. Share commission with the Broker.
	ProtoOAChangeBalanceType_BALANCE_WITHDRAW_IB_SHARED_PERCENTAGE_TO_BROKER  ProtoOAChangeBalanceType = 14 // IB commissions.
	ProtoOAChangeBalanceType_BALANCE_DEPOSIT_DIVIDENDS                        ProtoOAChangeBalanceType = 15 // Deposit dividends payments.
	ProtoOAChangeBalanceType_BALANCE_WITHDRAW_DIVIDENDS                       ProtoOAChangeBalanceType = 16 // Negative dividend charge for short position.
	ProtoOAChangeBalanceType_BALANCE_WITHDRAW_GSL_CHARGE                      ProtoOAChangeBalanceType = 17 // Charge for guaranteedStopLoss.
	ProtoOAChangeBalanceType_BALANCE_WITHDRAW_ROLLOVER                        ProtoOAChangeBalanceType = 18 // Charge of rollover fee for Shariah compliant accounts.
	ProtoOAChangeBalanceType_BALANCE_DEPOSIT_NONWITHDRAWABLE_BONUS            ProtoOAChangeBalanceType = 19 // Broker's operation to deposit bonus.
	ProtoOAChangeBalanceType_BALANCE_WITHDRAW_NONWITHDRAWABLE_BONUS           ProtoOAChangeBalanceType = 20 // Broker's operation to withdrawal bonus.
	ProtoOAChangeBalanceType_BALANCE_DEPOSIT_SWAP                             ProtoOAChangeBalanceType = 21 // Deposits of negative SWAP.
	ProtoOAChangeBalanceType_BALANCE_WITHDRAW_SWAP                            ProtoOAChangeBalanceType = 22 // SWAP charges.
	ProtoOAChangeBalanceType_BALANCE_DEPOSIT_MANAGEMENT_FEE                   ProtoOAChangeBalanceType = 27 // Mirroring commission.
	ProtoOAChangeBalanceType_BALANCE_WITHDRAW_MANAGEMENT_FEE                  ProtoOAChangeBalanceType = 28 // Mirroring commission. Deprecated since 7.1 in favor of BALANCE_WITHDRAW_COPY_FEE (34).
	ProtoOAChangeBalanceType_BALANCE_DEPOSIT_PERFORMANCE_FEE                  ProtoOAChangeBalanceType = 29 // Mirroring commission.
	ProtoOAChangeBalanceType_BALANCE_WITHDRAW_FOR_SUBACCOUNT                  ProtoOAChangeBalanceType = 30 // Withdraw for subaccount creation (cTrader Copy).
	ProtoOAChangeBalanceType_BALANCE_DEPOSIT_TO_SUBACCOUNT                    ProtoOAChangeBalanceType = 31 // Deposit to subaccount on creation (cTrader Copy).
	ProtoOAChangeBalanceType_BALANCE_WITHDRAW_FROM_SUBACCOUNT                 ProtoOAChangeBalanceType = 32 // Manual user's withdraw from subaccount (cTrader Copy), to parent account.
	ProtoOAChangeBalanceType_BALANCE_DEPOSIT_FROM_SUBACCOUNT                  ProtoOAChangeBalanceType = 33 // Manual user's deposit to subaccount (cTrader Copy), from parent account.
	ProtoOAChangeBalanceType_BALANCE_WITHDRAW_COPY_FEE                        ProtoOAChangeBalanceType = 34 // Withdrawal fees to Strategy Provider.
	ProtoOAChangeBalanceType_BALANCE_WITHDRAW_INACTIVITY_FEE                  ProtoOAChangeBalanceType = 35 // Withdraw of inactivity fee from the balance.
	ProtoOAChangeBalanceType_BALANCE_DEPOSIT_TRANSFER                         ProtoOAChangeBalanceType = 36 // Deposit within the same server (from another account).
	ProtoOAChangeBalanceType_BALANCE_WITHDRAW_TRANSFER                        ProtoOAChangeBalanceType = 37 // Withdraw within the same server (to another account).
	ProtoOAChangeBalanceType_BALANCE_DEPOSIT_CONVERTED_BONUS                  ProtoOAChangeBalanceType = 38 // Bonus being converted from virtual bonus to real deposit.
	ProtoOAChangeBalanceType_BALANCE_DEPOSIT_NEGATIVE_BALANCE_PROTECTION      ProtoOAChangeBalanceType = 39 // Applies if negative balance protection is configured by broker, should make balance = 0.
)

var (
	ProtoOAChangeBalanceTypeToString = map[ProtoOAChangeBalanceType]string{
		0:  "BALANCE_DEPOSIT",
		1:  "BALANCE_WITHDRAW",
		3:  "BALANCE_DEPOSIT_STRATEGY_COMMISSION_INNER",
		4:  "BALANCE_WITHDRAW_STRATEGY_COMMISSION_INNER",
		5:  "BALANCE_DEPOSIT_IB_COMMISSIONS",
		6:  "BALANCE_WITHDRAW_IB_SHARED_PERCENTAGE",
		7:  "BALANCE_DEPOSIT_IB_SHARED_PERCENTAGE_FROM_SUB_IB",
		8:  "BALANCE_DEPOSIT_IB_SHARED_PERCENTAGE_FROM_BROKER",
		9:  "BALANCE_DEPOSIT_REBATE",
		10: "BALANCE_WITHDRAW_REBATE",
		11: "BALANCE_DEPOSIT_STRATEGY_COMMISSION_OUTER",
		12: "BALANCE_WITHDRAW_STRATEGY_COMMISSION_OUTER",
		13: "BALANCE_WITHDRAW_BONUS_COMPENSATION",
		14: "BALANCE_WITHDRAW_IB_SHARED_PERCENTAGE_TO_BROKER",
		15: "BALANCE_DEPOSIT_DIVIDENDS",
		16: "BALANCE_WITHDRAW_DIVIDENDS",
		17: "BALANCE_WITHDRAW_GSL_CHARGE",
		18: "BALANCE_WITHDRAW_ROLLOVER",
		19: "BALANCE_DEPOSIT_NONWITHDRAWABLE_BONUS",
		20: "BALANCE_WITHDRAW_NONWITHDRAWABLE_BONUS",
		21: "BALANCE_DEPOSIT_SWAP",
		22: "BALANCE_WITHDRAW_SWAP",
		27: "BALANCE_DEPOSIT_MANAGEMENT_FEE",
		28: "BALANCE_WITHDRAW_MANAGEMENT_FEE",
		29: "BALANCE_DEPOSIT_PERFORMANCE_FEE",
		30: "BALANCE_WITHDRAW_FOR_SUBACCOUNT",
		31: "BALANCE_DEPOSIT_TO_SUBACCOUNT",
		32: "BALANCE_WITHDRAW_FROM_SUBACCOUNT",
		33: "BALANCE_DEPOSIT_FROM_SUBACCOUNT",
		34: "BALANCE_WITHDRAW_COPY_FEE",
		35: "BALANCE_WITHDRAW_INACTIVITY_FEE",
		36: "BALANCE_DEPOSIT_TRANSFER",
		37: "BALANCE_WITHDRAW_TRANSFER",
		38: "BALANCE_DEPOSIT_CONVERTED_BONUS",
		39: "BALANCE_DEPOSIT_NEGATIVE_BALANCE_PROTECTION",
	}
	ProtoOAChangeBalanceTypeByString = map[string]ProtoOAChangeBalanceType{
		"BALANCE_DEPOSIT":                                  0,
		"BALANCE_WITHDRAW":                                 1,
		"BALANCE_DEPOSIT_STRATEGY_COMMISSION_INNER":        3,
		"BALANCE_WITHDRAW_STRATEGY_COMMISSION_INNER":       4,
		"BALANCE_DEPOSIT_IB_COMMISSIONS":                   5,
		"BALANCE_WITHDRAW_IB_SHARED_PERCENTAGE":            6,
		"BALANCE_DEPOSIT_IB_SHARED_PERCENTAGE_FROM_SUB_IB": 7,
		"BALANCE_DEPOSIT_IB_SHARED_PERCENTAGE_FROM_BROKER": 8,
		"BALANCE_DEPOSIT_REBATE":                           9,
		"BALANCE_WITHDRAW_REBATE":                          10,
		"BALANCE_DEPOSIT_STRATEGY_COMMISSION_OUTER":        11,
		"BALANCE_WITHDRAW_STRATEGY_COMMISSION_OUTER":       12,
		"BALANCE_WITHDRAW_BONUS_COMPENSATION":              13,
		"BALANCE_WITHDRAW_IB_SHARED_PERCENTAGE_TO_BROKER":  14,
		"BALANCE_DEPOSIT_DIVIDENDS":                        15,
		"BALANCE_WITHDRAW_DIVIDENDS":                       16,
		"BALANCE_WITHDRAW_GSL_CHARGE":                      17,
		"BALANCE_WITHDRAW_ROLLOVER":                        18,
		"BALANCE_DEPOSIT_NONWITHDRAWABLE_BONUS":            19,
		"BALANCE_WITHDRAW_NONWITHDRAWABLE_BONUS":           20,
		"BALANCE_DEPOSIT_SWAP":                             21,
		"BALANCE_WITHDRAW_SWAP":                            22,
		"BALANCE_DEPOSIT_MANAGEMENT_FEE":                   27,
		"BALANCE_WITHDRAW_MANAGEMENT_FEE":                  28,
		"BALANCE_DEPOSIT_PERFORMANCE_FEE":                  29,
		"BALANCE_WITHDRAW_FOR_SUBACCOUNT":                  30,
		"BALANCE_DEPOSIT_TO_SUBACCOUNT":                    31,
		"BALANCE_WITHDRAW_FROM_SUBACCOUNT":                 32,
		"BALANCE_DEPOSIT_FROM_SUBACCOUNT":                  33,
		"BALANCE_WITHDRAW_COPY_FEE":                        34,
		"BALANCE_WITHDRAW_INACTIVITY_FEE":                  35,
		"BALANCE_DEPOSIT_TRANSFER":                         36,
		"BALANCE_WITHDRAW_TRANSFER":                        37,
		"BALANCE_DEPOSIT_CONVERTED_BONUS":                  38,
		"BALANCE_DEPOSIT_NEGATIVE_BALANCE_PROTECTION":      39,
	}
)

/*
Enum for deal status.
*/
type ProtoOADealStatus = messages.ProtoOADealStatus

const (
	ProtoOADealStatus_FILLED              ProtoOADealStatus = 2 // Deal filled.
	ProtoOADealStatus_PARTIALLY_FILLED    ProtoOADealStatus = 3 // Deal is partially filled.
	ProtoOADealStatus_REJECTED            ProtoOADealStatus = 4 // Deal is correct but was rejected by liquidity provider (e.g. no liquidity).
	ProtoOADealStatus_INTERNALLY_REJECTED ProtoOADealStatus = 5 // Deal rejected by server (e.g. no price quotes).
	ProtoOADealStatus_ERROR               ProtoOADealStatus = 6 // Deal is rejected by LP due to error (e.g. symbol is unknown).
	ProtoOADealStatus_MISSED              ProtoOADealStatus = 7 // Liquidity provider did not sent response on the deal during specified execution time period.
)

var (
	ProtoOADealStatusToString = map[ProtoOADealStatus]string{
		2: "FILLED",
		3: "PARTIALLY_FILLED",
		4: "REJECTED",
		5: "INTERNALLY_REJECTED",
		6: "ERROR",
		7: "MISSED",
	}
	ProtoOADealStatusByString = map[string]ProtoOADealStatus{
		"FILLED":              2,
		"PARTIALLY_FILLED":    3,
		"REJECTED":            4,
		"INTERNALLY_REJECTED": 5,
		"ERROR":               6,
		"MISSED":              7,
	}
)

/*
Enum for trendbar period.
*/
type ProtoOATrendbarPeriod = messages.ProtoOATrendbarPeriod

const (
	ProtoOATrendbarPeriod_M1  ProtoOATrendbarPeriod = 1
	ProtoOATrendbarPeriod_M2  ProtoOATrendbarPeriod = 2
	ProtoOATrendbarPeriod_M3  ProtoOATrendbarPeriod = 3
	ProtoOATrendbarPeriod_M4  ProtoOATrendbarPeriod = 4
	ProtoOATrendbarPeriod_M5  ProtoOATrendbarPeriod = 5
	ProtoOATrendbarPeriod_M10 ProtoOATrendbarPeriod = 6
	ProtoOATrendbarPeriod_M15 ProtoOATrendbarPeriod = 7
	ProtoOATrendbarPeriod_M30 ProtoOATrendbarPeriod = 8
	ProtoOATrendbarPeriod_H1  ProtoOATrendbarPeriod = 9
	ProtoOATrendbarPeriod_H4  ProtoOATrendbarPeriod = 10
	ProtoOATrendbarPeriod_H12 ProtoOATrendbarPeriod = 11
	ProtoOATrendbarPeriod_D1  ProtoOATrendbarPeriod = 12
	ProtoOATrendbarPeriod_W1  ProtoOATrendbarPeriod = 13
	ProtoOATrendbarPeriod_MN1 ProtoOATrendbarPeriod = 14
)

var (
	ProtoOATrendbarPeriodToString = map[ProtoOATrendbarPeriod]string{
		1:  "M1",
		2:  "M2",
		3:  "M3",
		4:  "M4",
		5:  "M5",
		6:  "M10",
		7:  "M15",
		8:  "M30",
		9:  "H1",
		10: "H4",
		11: "H12",
		12: "D1",
		13: "W1",
		14: "MN1",
	}
	ProtoOATrendbarPeriodByString = map[string]ProtoOATrendbarPeriod{
		"M1":  1,
		"M2":  2,
		"M3":  3,
		"M4":  4,
		"M5":  5,
		"M10": 6,
		"M15": 7,
		"M30": 8,
		"H1":  9,
		"H4":  10,
		"H12": 11,
		"D1":  12,
		"W1":  13,
		"MN1": 14,
	}
)

/*
Price quote type.
*/
type ProtoOAQuoteType = messages.ProtoOAQuoteType

const (
	ProtoOAQuoteType_BID ProtoOAQuoteType = 1
	ProtoOAQuoteType_ASK ProtoOAQuoteType = 2
)

var (
	ProtoOAQuoteTypeToString = map[ProtoOAQuoteType]string{
		1: "BID",
		2: "ASK",
	}
	ProtoOAQuoteTypeByString = map[string]ProtoOAQuoteType{
		"BID": 1,
		"ASK": 2,
	}
)

/*
Enum for Open API application permission in regards to token.
*/
type ProtoOAClientPermissionScope = messages.ProtoOAClientPermissionScope

const (
	ProtoOAClientPermissionScope_SCOPE_VIEW  ProtoOAClientPermissionScope = 0 // Allows to use only view commends. Trade is prohibited.
	ProtoOAClientPermissionScope_SCOPE_TRADE ProtoOAClientPermissionScope = 1 // Allows to use all commands.
)

var (
	ProtoOAClientPermissionScopeToString = map[ProtoOAClientPermissionScope]string{
		0: "SCOPE_VIEW",
		1: "SCOPE_TRADE",
	}
	ProtoOAClientPermissionScopeByString = map[string]ProtoOAClientPermissionScope{
		"SCOPE_VIEW":  0,
		"SCOPE_TRADE": 1,
	}
)

/*
Type of notification, currently only 3 instances of marginCall are supported.
*/
type ProtoOANotificationType = messages.ProtoOANotificationType

const (
	ProtoOANotificationType_MARGIN_LEVEL_THRESHOLD_1 ProtoOANotificationType = 61 // One of three margin calls, they are all similar.
	ProtoOANotificationType_MARGIN_LEVEL_THRESHOLD_2 ProtoOANotificationType = 62 // One of three margin calls, they are all similar.
	ProtoOANotificationType_MARGIN_LEVEL_THRESHOLD_3 ProtoOANotificationType = 63 // One of three margin calls, they are all similar.
)

var (
	ProtoOANotificationTypeToString = map[ProtoOANotificationType]string{
		61: "MARGIN_LEVEL_THRESHOLD_1",
		62: "MARGIN_LEVEL_THRESHOLD_2",
		63: "MARGIN_LEVEL_THRESHOLD_3",
	}
	ProtoOANotificationTypeByString = map[string]ProtoOANotificationType{
		"MARGIN_LEVEL_THRESHOLD_1": 61,
		"MARGIN_LEVEL_THRESHOLD_2": 62,
		"MARGIN_LEVEL_THRESHOLD_3": 63,
	}
)

/*
Enum for error code.
*/
type ProtoOAErrorCode = messages.ProtoOAErrorCode

const (
	// Authorization
	ProtoOAErrorCode_OA_AUTH_TOKEN_EXPIRED            ProtoOAErrorCode = 1   // When token used for account authorization is expired.
	ProtoOAErrorCode_ACCOUNT_NOT_AUTHORIZED           ProtoOAErrorCode = 2   // When account is not authorized.
	ProtoOAErrorCode_RET_NO_SUCH_LOGIN                ProtoOAErrorCode = 12  // When such account no longer exists.
	ProtoOAErrorCode_ALREADY_LOGGED_IN                ProtoOAErrorCode = 14  // When client tries to authorize after it was already authorized.
	ProtoOAErrorCode_RET_ACCOUNT_DISABLED             ProtoOAErrorCode = 64  // When account is disabled.
	ProtoOAErrorCode_CH_CLIENT_AUTH_FAILURE           ProtoOAErrorCode = 101 // Open API client is not activated or wrong client credentials.
	ProtoOAErrorCode_CH_CLIENT_NOT_AUTHENTICATED      ProtoOAErrorCode = 102 // When a command is sent for not authorized Open API client.
	ProtoOAErrorCode_CH_CLIENT_ALREADY_AUTHENTICATED  ProtoOAErrorCode = 103 // Client is trying to authenticate twice.
	ProtoOAErrorCode_CH_ACCESS_TOKEN_INVALID          ProtoOAErrorCode = 104 // Access token is invalid.
	ProtoOAErrorCode_CH_SERVER_NOT_REACHABLE          ProtoOAErrorCode = 105 // Trading service is not available.
	ProtoOAErrorCode_CH_CTID_TRADER_ACCOUNT_NOT_FOUND ProtoOAErrorCode = 106 // Trading account is not found.
	ProtoOAErrorCode_CH_OA_CLIENT_NOT_FOUND           ProtoOAErrorCode = 107 // Could not find this client id.
	// General
	ProtoOAErrorCode_REQUEST_FREQUENCY_EXCEEDED  ProtoOAErrorCode = 108 // Request frequency is reached.
	ProtoOAErrorCode_SERVER_IS_UNDER_MAINTENANCE ProtoOAErrorCode = 109 // Server is under maintenance.
	ProtoOAErrorCode_CHANNEL_IS_BLOCKED          ProtoOAErrorCode = 110 // Operations are not allowed for this account.
	ProtoOAErrorCode_CONNECTIONS_LIMIT_EXCEEDED  ProtoOAErrorCode = 67  // Limit of connections is reached for this Open API client.
	ProtoOAErrorCode_WORSE_GSL_NOT_ALLOWED       ProtoOAErrorCode = 68  // Not allowed to increase risk for Positions with Guaranteed Stop Loss.
	ProtoOAErrorCode_SYMBOL_HAS_HOLIDAY          ProtoOAErrorCode = 69  // Trading disabled because symbol has holiday.
	// Pricing
	ProtoOAErrorCode_NOT_SUBSCRIBED_TO_SPOTS ProtoOAErrorCode = 112 // When trying to subscribe to depth, trendbars, etc. without spot subscription.
	ProtoOAErrorCode_ALREADY_SUBSCRIBED      ProtoOAErrorCode = 113 // When subscription is requested for an active.
	ProtoOAErrorCode_SYMBOL_NOT_FOUND        ProtoOAErrorCode = 114 // Symbol not found.
	ProtoOAErrorCode_UNKNOWN_SYMBOL          ProtoOAErrorCode = 115 // Note: to be merged with SYMBOL_NOT_FOUND.
	ProtoOAErrorCode_INCORRECT_BOUNDARIES    ProtoOAErrorCode = 35  // When requested period (from,to) is too large or invalid values are set to from/to.
	// Trading
	ProtoOAErrorCode_NO_QUOTES                         ProtoOAErrorCode = 117 // Trading cannot be done as not quotes are available. Applicable for Book B.
	ProtoOAErrorCode_NOT_ENOUGH_MONEY                  ProtoOAErrorCode = 118 // Not enough funds to allocate margin.
	ProtoOAErrorCode_MAX_EXPOSURE_REACHED              ProtoOAErrorCode = 119 // Max exposure limit is reached for a {trader, symbol, side}.
	ProtoOAErrorCode_POSITION_NOT_FOUND                ProtoOAErrorCode = 120 // Position not found.
	ProtoOAErrorCode_ORDER_NOT_FOUND                   ProtoOAErrorCode = 121 // Order not found.
	ProtoOAErrorCode_POSITION_NOT_OPEN                 ProtoOAErrorCode = 122 // When trying to close a position that it is not open.
	ProtoOAErrorCode_POSITION_LOCKED                   ProtoOAErrorCode = 123 // Position in the state that does not allow to perform an operation.
	ProtoOAErrorCode_TOO_MANY_POSITIONS                ProtoOAErrorCode = 124 // Trading account reached its limit for max number of open positions and orders.
	ProtoOAErrorCode_TRADING_BAD_VOLUME                ProtoOAErrorCode = 125 // Invalid volume.
	ProtoOAErrorCode_TRADING_BAD_STOPS                 ProtoOAErrorCode = 126 // Invalid stop price.
	ProtoOAErrorCode_TRADING_BAD_PRICES                ProtoOAErrorCode = 127 // Invalid price (e.g. negative).
	ProtoOAErrorCode_TRADING_BAD_STAKE                 ProtoOAErrorCode = 128 // Invalid stake volume (e.g. negative).
	ProtoOAErrorCode_PROTECTION_IS_TOO_CLOSE_TO_MARKET ProtoOAErrorCode = 129 // Invalid protection prices.
	ProtoOAErrorCode_TRADING_BAD_EXPIRATION_DATE       ProtoOAErrorCode = 130 // Invalid expiration.
	ProtoOAErrorCode_PENDING_EXECUTION                 ProtoOAErrorCode = 131 // Unable to apply changes as position has an order under execution.
	ProtoOAErrorCode_TRADING_DISABLED                  ProtoOAErrorCode = 132 // Trading is blocked for the symbol.
	ProtoOAErrorCode_TRADING_NOT_ALLOWED               ProtoOAErrorCode = 133 // Trading account is in read only mode.
	ProtoOAErrorCode_UNABLE_TO_CANCEL_ORDER            ProtoOAErrorCode = 134 // Unable to cancel order.
	ProtoOAErrorCode_UNABLE_TO_AMEND_ORDER             ProtoOAErrorCode = 135 // Unable to amend order.
	ProtoOAErrorCode_SHORT_SELLING_NOT_ALLOWED         ProtoOAErrorCode = 136 // Short selling is not allowed.
)

var (
	ProtoOAErrorCodeToString = map[ProtoOAErrorCode]string{
		1:   "OA_AUTH_TOKEN_EXPIRED",
		2:   "ACCOUNT_NOT_AUTHORIZED",
		12:  "RET_NO_SUCH_LOGIN",
		14:  "ALREADY_LOGGED_IN",
		64:  "RET_ACCOUNT_DISABLED",
		101: "CH_CLIENT_AUTH_FAILURE",
		102: "CH_CLIENT_NOT_AUTHENTICATED",
		103: "CH_CLIENT_ALREADY_AUTHENTICATED",
		104: "CH_ACCESS_TOKEN_INVALID",
		105: "CH_SERVER_NOT_REACHABLE",
		106: "CH_CTID_TRADER_ACCOUNT_NOT_FOUND",
		107: "CH_OA_CLIENT_NOT_FOUND",
		108: "REQUEST_FREQUENCY_EXCEEDED",
		109: "SERVER_IS_UNDER_MAINTENANCE",
		110: "CHANNEL_IS_BLOCKED",
		67:  "CONNECTIONS_LIMIT_EXCEEDED",
		68:  "WORSE_GSL_NOT_ALLOWED",
		69:  "SYMBOL_HAS_HOLIDAY",
		112: "NOT_SUBSCRIBED_TO_SPOTS",
		113: "ALREADY_SUBSCRIBED",
		114: "SYMBOL_NOT_FOUND",
		115: "UNKNOWN_SYMBOL",
		35:  "INCORRECT_BOUNDARIES",
		117: "NO_QUOTES",
		118: "NOT_ENOUGH_MONEY",
		119: "MAX_EXPOSURE_REACHED",
		120: "POSITION_NOT_FOUND",
		121: "ORDER_NOT_FOUND",
		122: "POSITION_NOT_OPEN",
		123: "POSITION_LOCKED",
		124: "TOO_MANY_POSITIONS",
		125: "TRADING_BAD_VOLUME",
		126: "TRADING_BAD_STOPS",
		127: "TRADING_BAD_PRICES",
		128: "TRADING_BAD_STAKE",
		129: "PROTECTION_IS_TOO_CLOSE_TO_MARKET",
		130: "TRADING_BAD_EXPIRATION_DATE",
		131: "PENDING_EXECUTION",
		132: "TRADING_DISABLED",
		133: "TRADING_NOT_ALLOWED",
		134: "UNABLE_TO_CANCEL_ORDER",
		135: "UNABLE_TO_AMEND_ORDER",
		136: "SHORT_SELLING_NOT_ALLOWED",
	}
	ProtoOAErrorCodeByString = map[string]ProtoOAErrorCode{
		"OA_AUTH_TOKEN_EXPIRED":             1,
		"ACCOUNT_NOT_AUTHORIZED":            2,
		"RET_NO_SUCH_LOGIN":                 12,
		"ALREADY_LOGGED_IN":                 14,
		"RET_ACCOUNT_DISABLED":              64,
		"CH_CLIENT_AUTH_FAILURE":            101,
		"CH_CLIENT_NOT_AUTHENTICATED":       102,
		"CH_CLIENT_ALREADY_AUTHENTICATED":   103,
		"CH_ACCESS_TOKEN_INVALID":           104,
		"CH_SERVER_NOT_REACHABLE":           105,
		"CH_CTID_TRADER_ACCOUNT_NOT_FOUND":  106,
		"CH_OA_CLIENT_NOT_FOUND":            107,
		"REQUEST_FREQUENCY_EXCEEDED":        108,
		"SERVER_IS_UNDER_MAINTENANCE":       109,
		"CHANNEL_IS_BLOCKED":                110,
		"CONNECTIONS_LIMIT_EXCEEDED":        67,
		"WORSE_GSL_NOT_ALLOWED":             68,
		"SYMBOL_HAS_HOLIDAY":                69,
		"NOT_SUBSCRIBED_TO_SPOTS":           112,
		"ALREADY_SUBSCRIBED":                113,
		"SYMBOL_NOT_FOUND":                  114,
		"UNKNOWN_SYMBOL":                    115,
		"INCORRECT_BOUNDARIES":              35,
		"NO_QUOTES":                         117,
		"NOT_ENOUGH_MONEY":                  118,
		"MAX_EXPOSURE_REACHED":              119,
		"POSITION_NOT_FOUND":                120,
		"ORDER_NOT_FOUND":                   121,
		"POSITION_NOT_OPEN":                 122,
		"POSITION_LOCKED":                   123,
		"TOO_MANY_POSITIONS":                124,
		"TRADING_BAD_VOLUME":                125,
		"TRADING_BAD_STOPS":                 126,
		"TRADING_BAD_PRICES":                127,
		"TRADING_BAD_STAKE":                 128,
		"PROTECTION_IS_TOO_CLOSE_TO_MARKET": 129,
		"TRADING_BAD_EXPIRATION_DATE":       130,
		"PENDING_EXECUTION":                 131,
		"TRADING_DISABLED":                  132,
		"TRADING_NOT_ALLOWED":               133,
		"UNABLE_TO_CANCEL_ORDER":            134,
		"UNABLE_TO_AMEND_ORDER":             135,
		"SHORT_SELLING_NOT_ALLOWED":         136,
	}
)

/*
Enum for limited risk margin calculation strategy.
*/
type ProtoOALimitedRiskMarginCalculationStrategy = messages.ProtoOALimitedRiskMarginCalculationStrategy

const (
	ProtoOALimitedRiskMarginCalculationStrategy_ACCORDING_TO_LEVERAGE         ProtoOALimitedRiskMarginCalculationStrategy = 0
	ProtoOALimitedRiskMarginCalculationStrategy_ACCORDING_TO_GSL              ProtoOALimitedRiskMarginCalculationStrategy = 1
	ProtoOALimitedRiskMarginCalculationStrategy_ACCORDING_TO_GSL_AND_LEVERAGE ProtoOALimitedRiskMarginCalculationStrategy = 2
)

var (
	ProtoOALimitedRiskMarginCalculationStrategyToString = map[ProtoOALimitedRiskMarginCalculationStrategy]string{
		0: "ACCORDING_TO_LEVERAGE",
		1: "ACCORDING_TO_GSL",
		2: "ACCORDING_TO_GSL_AND_LEVERAGE",
	}
	ProtoOALimitedRiskMarginCalculationStrategyByString = map[string]ProtoOALimitedRiskMarginCalculationStrategy{
		"ACCORDING_TO_LEVERAGE":         0,
		"ACCORDING_TO_GSL":              1,
		"ACCORDING_TO_GSL_AND_LEVERAGE": 2,
	}
)

/*
Enum for stop out strategy.
*/
type ProtoOAStopOutStrategy = messages.ProtoOAStopOutStrategy

const (
	ProtoOAStopOutStrategy_MOST_MARGIN_USED_FIRST ProtoOAStopOutStrategy = 0 // A Stop Out strategy that closes a Position with the largest Used Margin.
	ProtoOAStopOutStrategy_MOST_LOSING_FIRST      ProtoOAStopOutStrategy = 1 // A Stop Out strategy that closes a Position with the least PnL.
)

var (
	ProtoOAStopOutStrategyToString = map[ProtoOAStopOutStrategy]string{
		0: "MOST_MARGIN_USED_FIRST",
		1: "MOST_LOSING_FIRST",
	}
	ProtoOAStopOutStrategyByString = map[string]ProtoOAStopOutStrategy{
		"MOST_MARGIN_USED_FIRST": 0,
		"MOST_LOSING_FIRST":      1,
	}
)

/*
Asset entity.
*/
type ProtoOAAsset = messages.ProtoOAAsset

/*
Trading symbol entity.
*/
type ProtoOASymbol = messages.ProtoOASymbol

/*
ProtoOALightSymbol entity.
*/
type ProtoOALightSymbol = messages.ProtoOALightSymbol

/*
Archived symbol entity.
*/
type ProtoOAArchivedSymbol = messages.ProtoOAArchivedSymbol

/*
Symbol category entity.
*/
type ProtoOASymbolCategory = messages.ProtoOASymbolCategory

/*
Symbol trading session entity.
*/
type ProtoOAInterval = messages.ProtoOAInterval

/*
Trading account entity.
*/
type ProtoOATrader = messages.ProtoOATrader

/*
Trade position entity.
*/
type ProtoOAPosition = messages.ProtoOAPosition

/*
Position/order trading details entity.
*/
type ProtoOATradeData = messages.ProtoOATradeData

/*
Trade order entity.
*/
type ProtoOAOrder = messages.ProtoOAOrder

/*
Bonus deposit/withdrawal entity.
*/
type ProtoOABonusDepositWithdraw = messages.ProtoOABonusDepositWithdraw

/*
Account deposit/withdrawal operation entity.
*/
type ProtoOADepositWithdraw = messages.ProtoOADepositWithdraw

/*
Execution entity.
*/
type ProtoOADeal = messages.ProtoOADeal

/*
Deal details for ProtoOADealOffsetListReq.
*/
type ProtoOADealOffset = messages.ProtoOADealOffset

/*
Trading details for closing deal.
*/
type ProtoOAClosePositionDetail = messages.ProtoOAClosePositionDetail

/*
Historical Trendbar entity.
*/
type ProtoOATrendbar = messages.ProtoOATrendbar

/*
Expected margin computation entity.
*/
type ProtoOAExpectedMargin = messages.ProtoOAExpectedMargin

/*
Historical tick data type.
*/
type ProtoOATickData = messages.ProtoOATickData

/*
Trader profile entity. Empty due to GDPR.
*/
type ProtoOACtidProfile = messages.ProtoOACtidProfile

/*
Trader account entity.
*/
type ProtoOACtidTraderAccount = messages.ProtoOACtidTraderAccount

/*
Asset class entity.
*/
type ProtoOAAssetClass = messages.ProtoOAAssetClass

/*
Depth of market entity.
*/
type ProtoOADepthQuote = messages.ProtoOADepthQuote

/*
Margin call entity, specifies threshold for exact margin call type. Only 3 instances of margin calls are supported, identified by marginCallType. See ProtoOANotificationType for details.
*/
type ProtoOAMarginCall = messages.ProtoOAMarginCall

/*
Holiday entity.
*/
type ProtoOAHoliday = messages.ProtoOAHoliday

/*
Dynamic leverage entity.
*/
type ProtoOADynamicLeverage = messages.ProtoOADynamicLeverage

/*
Dynamic leverage tier entity.
*/
type ProtoOADynamicLeverageTier = messages.ProtoOADynamicLeverageTier

/*
Position unrealized profit and loss entity.
*/
type ProtoOAPositionUnrealizedPnL = messages.ProtoOAPositionUnrealizedPnL

/*
--------------
Messages.pb.go
--------------
*/

// ProtoOAApplicationAuthReq left out due to explicit handling

// ProtoOAApplicationAuthRes left out due to explicit handling

// ProtoOAAccountAuthReq left out due to explicit handling

/*
Response to the ProtoOAApplicationAuthReq request.
*/
type ProtoOAAccountAuthRes = messages.ProtoOAAccountAuthRes

/*
Generic response when an ERROR occurred.
*/
type ProtoOAErrorRes = messages.ProtoOAErrorRes

/*
Event that is sent when the connection with the client application is cancelled by the server. All the sessions for the traders' accounts will be terminated.

IMPORTANT: Do not try to recover the connection. The client takes care of reconnecting and re-establishing the sessions internally.
*/
type ProtoOAClientDisconnectEvent = messages.ProtoOAClientDisconnectEvent

/*
Event that is sent when a session to a specific trader's account is terminated by the server but the existing connections with the other trader's accounts are maintained. Reasons to trigger: account was deleted, cTID was deleted, token was refreshed, token was revoked.
*/
type ProtoOAAccountsTokenInvalidatedEvent = messages.ProtoOAAccountsTokenInvalidatedEvent

/*
Request for getting the proxy version. Can be used to check the current version of the Open API scheme.
*/
type ProtoOAVersionReq = messages.ProtoOAVersionReq

/*
Response to the ProtoOAVersionReq request.
*/
type ProtoOAVersionRes = messages.ProtoOAVersionRes

/*
Request for sending a new trading order. Allowed only if the accessToken has the "trade" permissions for the trading account.
*/
type ProtoOANewOrderReq = messages.ProtoOANewOrderReq

/*
Event that is sent following the successful order acceptance or execution by the server. Acts as response to the ProtoOANewOrderReq, ProtoOACancelOrderReq, ProtoOAAmendOrderReq, ProtoOAAmendPositionSLTPReq, ProtoOAClosePositionReq requests. Also, the event is sent when a Deposit/Withdrawal took place.
*/
type ProtoOAExecutionEvent = messages.ProtoOAExecutionEvent

/*
Request for cancelling existing pending order. Allowed only if the accessToken has "trade" permissions for the trading account.
*/
type ProtoOACancelOrderReq = messages.ProtoOACancelOrderReq

/*
Request for amending the existing pending order. Allowed only if the Access Token has "trade" permissions for the trading account.
*/
type ProtoOAAmendOrderReq = messages.ProtoOAAmendOrderReq

/*
Request for amending StopLoss and TakeProfit of existing position. Allowed only if the accessToken has "trade" permissions for the trading account.
*/
type ProtoOAAmendPositionSLTPReq = messages.ProtoOAAmendPositionSLTPReq

/*
Request for closing or partially closing of an existing position. Allowed only if the accessToken has "trade" permissions for the trading account.
*/
type ProtoOAClosePositionReq = messages.ProtoOAClosePositionReq

/*
Event that is sent when the level of the Trailing Stop Loss is changed due to the price level changes.
*/
type ProtoOATrailingSLChangedEvent = messages.ProtoOATrailingSLChangedEvent

/*
Request for the list of assets available for a trader's account.
*/
type ProtoOAAssetListReq = messages.ProtoOAAssetListReq

/*
Response to the ProtoOAAssetListReq request.
*/
type ProtoOAAssetListRes = messages.ProtoOAAssetListRes

/*
Request for a list of symbols available for a trading account. Symbol entries are returned with the limited set of fields.
*/
type ProtoOASymbolsListReq = messages.ProtoOASymbolsListReq

/*
Response to the ProtoOASymbolsListReq request.
*/
type ProtoOASymbolsListRes = messages.ProtoOASymbolsListRes

/*
Request for getting a full symbol entity.
*/
type ProtoOASymbolByIdReq = messages.ProtoOASymbolByIdReq

/*
Response to the ProtoOASymbolByIdReq request.
*/
type ProtoOASymbolByIdRes = messages.ProtoOASymbolByIdRes

/*
Request for getting a conversion chain between two assets that consists of several symbols. Use when no direct quote is available.
*/
type ProtoOASymbolsForConversionReq = messages.ProtoOASymbolsForConversionReq

/*
Response to the ProtoOASymbolsForConversionReq request.
*/
type ProtoOASymbolsForConversionRes = messages.ProtoOASymbolsForConversionRes

/*
Event that is sent when the symbol is changed on the Server side.
*/
type ProtoOASymbolChangedEvent = messages.ProtoOASymbolChangedEvent

/*
Request for a list of asset classes available for the trader's account.
*/
type ProtoOAAssetClassListReq = messages.ProtoOAAssetClassListReq

/*
Response to the ProtoOAAssetListReq request.
*/
type ProtoOAAssetClassListRes = messages.ProtoOAAssetClassListRes

/*
Request for getting data of Trader's Account.
*/
type ProtoOATraderReq = messages.ProtoOATraderReq

/*
Response to the ProtoOATraderReq request.
*/
type ProtoOATraderRes = messages.ProtoOATraderRes

/*
Event that is sent when a Trader is updated on Server side.
*/
type ProtoOATraderUpdatedEvent = messages.ProtoOATraderUpdatedEvent

/*
Request for getting Trader's current open positions and pending orders data.
*/
type ProtoOAReconcileReq = messages.ProtoOAReconcileReq

/*
Response to the ProtoOAReconcileReq request.
*/
type ProtoOAReconcileRes = messages.ProtoOAReconcileRes

/*
Event that is sent when errors occur during the order requests.
*/
type ProtoOAOrderErrorEvent = messages.ProtoOAOrderErrorEvent

/*
Request for getting Trader's deals historical data (execution details).
*/
type ProtoOADealListReq = messages.ProtoOADealListReq

/*
Response to the ProtoOADealListReq request.
*/
type ProtoOADealListRes = messages.ProtoOADealListRes

/*
Request for getting Trader's orders filtered by timestamp.
*/
type ProtoOAOrderListReq = messages.ProtoOAOrderListReq

/*
Response to the ProtoOAOrderListReq request.
*/
type ProtoOAOrderListRes = messages.ProtoOAOrderListRes

/*
Request for getting the margin estimate according to leverage profiles. Can be used before sending a new order request. This doesn't consider ACCORDING_TO_GSL margin calculation type, as this calculation is trivial: usedMargin = (VWAP price of the position - GSL price) * volume * Quote2Deposit.
*/
type ProtoOAExpectedMarginReq = messages.ProtoOAExpectedMarginReq

/*
Response to the ProtoOAExpectedMarginReq request.
*/
type ProtoOAExpectedMarginRes = messages.ProtoOAExpectedMarginRes

/*
Event that is sent when the margin allocated to a specific position is changed.
*/
type ProtoOAMarginChangedEvent = messages.ProtoOAMarginChangedEvent

/*
Request for getting Trader's historical data of deposits and withdrawals.
*/
type ProtoOACashFlowHistoryListReq = messages.ProtoOACashFlowHistoryListReq

/*
Response to the ProtoOACashFlowHistoryListReq request.
*/
type ProtoOACashFlowHistoryListRes = messages.ProtoOACashFlowHistoryListRes

/*
Request for getting the list of granted trader's account for the access token.
*/
type ProtoOAGetAccountListByAccessTokenReq = messages.ProtoOAGetAccountListByAccessTokenReq

/*
Response to the ProtoOAGetAccountListByAccessTokenReq request.
*/
type ProtoOAGetAccountListByAccessTokenRes = messages.ProtoOAGetAccountListByAccessTokenRes

// ProtoOARefreshTokenReq left out due to explicit handling

/*
Response to the ProtoOARefreshTokenReq request.
*/
type ProtoOARefreshTokenRes = messages.ProtoOARefreshTokenRes

// ProtoOASubscribeSpotsReq left out due to explicit handling

// ProtoOASubscribeSpotsRes left out due to explicit handling

// ProtoOAUnsubscribeSpotsReq left out due to explicit handling

// ProtoOAUnsubscribeSpotsRes left out due to explicit handling

/*
Event that is sent when a new spot event is generated on the server side. Requires subscription on the spot events, see ProtoOASubscribeSpotsReq. First event, received after subscription will contain latest spot prices even if market is closed.
*/
type ProtoOASpotEvent = messages.ProtoOASpotEvent

// ProtoOASubscribeLiveTrendbarReq left out due to explicit handling

// ProtoOASubscribeLiveTrendbarRes left out due to explicit handling

// ProtoOAUnsubscribeLiveTrendbarReq left out due to explicit handling

// ProtoOAUnsubscribeLiveTrendbarRes left out due to explicit handling

/*
Request for getting historical trend bars for the symbol.
*/
type ProtoOAGetTrendbarsReq = messages.ProtoOAGetTrendbarsReq

/*
Response to the ProtoOAGetTrendbarsReq request.
*/
type ProtoOAGetTrendbarsRes = messages.ProtoOAGetTrendbarsRes

/*
Request for getting historical tick data for the symbol.
*/
type ProtoOAGetTickDataReq = messages.ProtoOAGetTickDataReq

/*
Response to the ProtoOAGetTickDataReq request.
*/
type ProtoOAGetTickDataRes = messages.ProtoOAGetTickDataRes

/*
Request for getting details of Trader's profile. Limited due to GDRP requirements.
*/
type ProtoOAGetCtidProfileByTokenReq = messages.ProtoOAGetCtidProfileByTokenReq

/*
Response to the ProtoOAGetCtidProfileByTokenReq request.
*/
type ProtoOAGetCtidProfileByTokenRes = messages.ProtoOAGetCtidProfileByTokenRes

/*
Event that is sent when the structure of depth of market is changed. Requires subscription on the depth of markets for the symbol, see ProtoOASubscribeDepthQuotesReq.
*/
type ProtoOADepthEvent = messages.ProtoOADepthEvent

// ProtoOASubscribeDepthQuotesReq left out due to explicit handling

// ProtoOASubscribeDepthQuotesRes left out due to explicit handling

// ProtoOAUnsubscribeDepthQuotesReq left out due to explicit handling

// ProtoOAUnsubscribeDepthQuotesRes left out due to explicit handling

/*
Request for a list of symbol categories available for a trading account.
*/
type ProtoOASymbolCategoryListReq = messages.ProtoOASymbolCategoryListReq

/*
Response to the ProtoSymbolCategoryListReq request.
*/
type ProtoOASymbolCategoryListRes = messages.ProtoOASymbolCategoryListRes

// ProtoOAAccountLogoutReq left out due to explicit handling

/*
Response to the ProtoOAAccountLogoutReq request. Actual logout of trading account will be completed on ProtoOAAccountDisconnectEvent.
*/
type ProtoOAAccountLogoutRes = messages.ProtoOAAccountLogoutRes

/*
Event that is sent when the established session for an account is dropped on the server side. A new session must be authorized for the account.
*/
type ProtoOAAccountDisconnectEvent = messages.ProtoOAAccountDisconnectEvent

/*
Request for a list of existing margin call thresholds configured for a user.
*/
type ProtoOAMarginCallListReq = messages.ProtoOAMarginCallListReq

/*
Response to the ProtoOATraderLogoutReq request. Has a list of existing user Margin Calls, usually contains 3 items.
*/
type ProtoOAMarginCallListRes = messages.ProtoOAMarginCallListRes

/*
Request to modify marginLevelThreshold of specified marginCallType for ctidTraderAccountId.
*/
type ProtoOAMarginCallUpdateReq = messages.ProtoOAMarginCallUpdateReq

/*
Response to the ProtoOAMarginCallUpdateReq request. If this response received, it means that margin call was successfully updated.
*/
type ProtoOAMarginCallUpdateRes = messages.ProtoOAMarginCallUpdateRes

/*
Event that is sent when a Margin Call threshold configuration is updated.
*/
type ProtoOAMarginCallUpdateEvent = messages.ProtoOAMarginCallUpdateEvent

/*
Event that is sent when account margin level reaches target marginLevelThreshold. Event is sent no more than once every 10 minutes to avoid spamming.
*/
type ProtoOAMarginCallTriggerEvent = messages.ProtoOAMarginCallTriggerEvent

/*
Request for getting a dynamic leverage entity referenced in ProtoOASymbol.leverageId.
*/
type ProtoOAGetDynamicLeverageByIDReq = messages.ProtoOAGetDynamicLeverageByIDReq

/*
Response to the ProtoOAGetDynamicLeverageByIDReq request.
*/
type ProtoOAGetDynamicLeverageByIDRes = messages.ProtoOAGetDynamicLeverageByIDRes

/*
Request for retrieving the deals related to a position.
*/
type ProtoOADealListByPositionIdReq = messages.ProtoOADealListByPositionIdReq

/*
Response to the ProtoOADealListByPositionIdReq request.
*/
type ProtoOADealListByPositionIdRes = messages.ProtoOADealListByPositionIdRes

/*
Request for getting Order and its related Deals.
*/
type ProtoOAOrderDetailsReq = messages.ProtoOAOrderDetailsReq

/*
Response to the ProtoOAOrderDetailsReq request.
*/
type ProtoOAOrderDetailsRes = messages.ProtoOAOrderDetailsRes

/*
Request for retrieving Orders related to a Position by using Position ID. Filtered by utcLastUpdateTimestamp.
*/
type ProtoOAOrderListByPositionIdReq = messages.ProtoOAOrderListByPositionIdReq

/*
Response to the ProtoOAOrderListByPositionIdReq request.
*/
type ProtoOAOrderListByPositionIdRes = messages.ProtoOAOrderListByPositionIdRes

/*
Request for getting sets of Deals that were offset by a specific Deal and that are offsetting the Deal.
*/
type ProtoOADealOffsetListReq = messages.ProtoOADealOffsetListReq

/*
Response to the ProtoOADealOffsetListReq request.
*/
type ProtoOADealOffsetListRes = messages.ProtoOADealOffsetListRes

/*
Request for getting trader's positions' unrealized PnLs.
*/
type ProtoOAGetPositionUnrealizedPnLReq = messages.ProtoOAGetPositionUnrealizedPnLReq

/*
Response to the ProtoOAGetPositionUnrealizedPnLReq request.
*/
type ProtoOAGetPositionUnrealizedPnLRes = messages.ProtoOAGetPositionUnrealizedPnLRes
