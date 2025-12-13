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
)

type ProtoPayloadType = messages.ProtoPayloadType

type ProtoOAPayloadType = messages.ProtoOAPayloadType

const (
	PROTO_OA_APPLICATION_AUTH_REQ             ProtoOAPayloadType = 2100
	PROTO_OA_APPLICATION_AUTH_RES             ProtoOAPayloadType = 2101
	PROTO_OA_ACCOUNT_AUTH_REQ                 ProtoOAPayloadType = 2102
	PROTO_OA_ACCOUNT_AUTH_RES                 ProtoOAPayloadType = 2103
	PROTO_OA_VERSION_REQ                      ProtoOAPayloadType = 2104
	PROTO_OA_VERSION_RES                      ProtoOAPayloadType = 2105
	PROTO_OA_NEW_ORDER_REQ                    ProtoOAPayloadType = 2106
	PROTO_OA_TRAILING_SL_CHANGED_EVENT        ProtoOAPayloadType = 2107
	PROTO_OA_CANCEL_ORDER_REQ                 ProtoOAPayloadType = 2108
	PROTO_OA_AMEND_ORDER_REQ                  ProtoOAPayloadType = 2109
	PROTO_OA_AMEND_POSITION_SLTP_REQ          ProtoOAPayloadType = 2110
	PROTO_OA_CLOSE_POSITION_REQ               ProtoOAPayloadType = 2111
	PROTO_OA_ASSET_LIST_REQ                   ProtoOAPayloadType = 2112
	PROTO_OA_ASSET_LIST_RES                   ProtoOAPayloadType = 2113
	PROTO_OA_SYMBOLS_LIST_REQ                 ProtoOAPayloadType = 2114
	PROTO_OA_SYMBOLS_LIST_RES                 ProtoOAPayloadType = 2115
	PROTO_OA_SYMBOL_BY_ID_REQ                 ProtoOAPayloadType = 2116
	PROTO_OA_SYMBOL_BY_ID_RES                 ProtoOAPayloadType = 2117
	PROTO_OA_SYMBOLS_FOR_CONVERSION_REQ       ProtoOAPayloadType = 2118
	PROTO_OA_SYMBOLS_FOR_CONVERSION_RES       ProtoOAPayloadType = 2119
	PROTO_OA_SYMBOL_CHANGED_EVENT             ProtoOAPayloadType = 2120
	PROTO_OA_TRADER_REQ                       ProtoOAPayloadType = 2121
	PROTO_OA_TRADER_RES                       ProtoOAPayloadType = 2122
	PROTO_OA_TRADER_UPDATE_EVENT              ProtoOAPayloadType = 2123
	PROTO_OA_RECONCILE_REQ                    ProtoOAPayloadType = 2124
	PROTO_OA_RECONCILE_RES                    ProtoOAPayloadType = 2125
	PROTO_OA_EXECUTION_EVENT                  ProtoOAPayloadType = 2126
	PROTO_OA_SUBSCRIBE_SPOTS_REQ              ProtoOAPayloadType = 2127
	PROTO_OA_SUBSCRIBE_SPOTS_RES              ProtoOAPayloadType = 2128
	PROTO_OA_UNSUBSCRIBE_SPOTS_REQ            ProtoOAPayloadType = 2129
	PROTO_OA_UNSUBSCRIBE_SPOTS_RES            ProtoOAPayloadType = 2130
	PROTO_OA_SPOT_EVENT                       ProtoOAPayloadType = 2131
	PROTO_OA_ORDER_ERROR_EVENT                ProtoOAPayloadType = 2132
	PROTO_OA_DEAL_LIST_REQ                    ProtoOAPayloadType = 2133
	PROTO_OA_DEAL_LIST_RES                    ProtoOAPayloadType = 2134
	PROTO_OA_SUBSCRIBE_LIVE_TRENDBAR_REQ      ProtoOAPayloadType = 2135
	PROTO_OA_UNSUBSCRIBE_LIVE_TRENDBAR_REQ    ProtoOAPayloadType = 2136
	PROTO_OA_GET_TRENDBARS_REQ                ProtoOAPayloadType = 2137
	PROTO_OA_GET_TRENDBARS_RES                ProtoOAPayloadType = 2138
	PROTO_OA_EXPECTED_MARGIN_REQ              ProtoOAPayloadType = 2139
	PROTO_OA_EXPECTED_MARGIN_RES              ProtoOAPayloadType = 2140
	PROTO_OA_MARGIN_CHANGED_EVENT             ProtoOAPayloadType = 2141
	PROTO_OA_ERROR_RES                        ProtoOAPayloadType = 2142
	PROTO_OA_CASH_FLOW_HISTORY_LIST_REQ       ProtoOAPayloadType = 2143
	PROTO_OA_CASH_FLOW_HISTORY_LIST_RES       ProtoOAPayloadType = 2144
	PROTO_OA_GET_TICKDATA_REQ                 ProtoOAPayloadType = 2145
	PROTO_OA_GET_TICKDATA_RES                 ProtoOAPayloadType = 2146
	PROTO_OA_ACCOUNTS_TOKEN_INVALIDATED_EVENT ProtoOAPayloadType = 2147
	PROTO_OA_CLIENT_DISCONNECT_EVENT          ProtoOAPayloadType = 2148
	PROTO_OA_GET_ACCOUNTS_BY_ACCESS_TOKEN_REQ ProtoOAPayloadType = 2149
	PROTO_OA_GET_ACCOUNTS_BY_ACCESS_TOKEN_RES ProtoOAPayloadType = 2150
	PROTO_OA_GET_CTID_PROFILE_BY_TOKEN_REQ    ProtoOAPayloadType = 2151
	PROTO_OA_GET_CTID_PROFILE_BY_TOKEN_RES    ProtoOAPayloadType = 2152
	PROTO_OA_ASSET_CLASS_LIST_REQ             ProtoOAPayloadType = 2153
	PROTO_OA_ASSET_CLASS_LIST_RES             ProtoOAPayloadType = 2154
	PROTO_OA_DEPTH_EVENT                      ProtoOAPayloadType = 2155
	PROTO_OA_SUBSCRIBE_DEPTH_QUOTES_REQ       ProtoOAPayloadType = 2156
	PROTO_OA_SUBSCRIBE_DEPTH_QUOTES_RES       ProtoOAPayloadType = 2157
	PROTO_OA_UNSUBSCRIBE_DEPTH_QUOTES_REQ     ProtoOAPayloadType = 2158
	PROTO_OA_UNSUBSCRIBE_DEPTH_QUOTES_RES     ProtoOAPayloadType = 2159
	PROTO_OA_SYMBOL_CATEGORY_REQ              ProtoOAPayloadType = 2160
	PROTO_OA_SYMBOL_CATEGORY_RES              ProtoOAPayloadType = 2161
	PROTO_OA_ACCOUNT_LOGOUT_REQ               ProtoOAPayloadType = 2162
	PROTO_OA_ACCOUNT_LOGOUT_RES               ProtoOAPayloadType = 2163
	PROTO_OA_ACCOUNT_DISCONNECT_EVENT         ProtoOAPayloadType = 2164
	PROTO_OA_SUBSCRIBE_LIVE_TRENDBAR_RES      ProtoOAPayloadType = 2165
	PROTO_OA_UNSUBSCRIBE_LIVE_TRENDBAR_RES    ProtoOAPayloadType = 2166
	PROTO_OA_MARGIN_CALL_LIST_REQ             ProtoOAPayloadType = 2167
	PROTO_OA_MARGIN_CALL_LIST_RES             ProtoOAPayloadType = 2168
	PROTO_OA_MARGIN_CALL_UPDATE_REQ           ProtoOAPayloadType = 2169
	PROTO_OA_MARGIN_CALL_UPDATE_RES           ProtoOAPayloadType = 2170
	PROTO_OA_MARGIN_CALL_UPDATE_EVENT         ProtoOAPayloadType = 2171
	PROTO_OA_MARGIN_CALL_TRIGGER_EVENT        ProtoOAPayloadType = 2172
	PROTO_OA_REFRESH_TOKEN_REQ                ProtoOAPayloadType = 2173
	PROTO_OA_REFRESH_TOKEN_RES                ProtoOAPayloadType = 2174
	PROTO_OA_ORDER_LIST_REQ                   ProtoOAPayloadType = 2175
	PROTO_OA_ORDER_LIST_RES                   ProtoOAPayloadType = 2176
	PROTO_OA_GET_DYNAMIC_LEVERAGE_REQ         ProtoOAPayloadType = 2177
	PROTO_OA_GET_DYNAMIC_LEVERAGE_RES         ProtoOAPayloadType = 2178
	PROTO_OA_DEAL_LIST_BY_POSITION_ID_REQ     ProtoOAPayloadType = 2179
	PROTO_OA_DEAL_LIST_BY_POSITION_ID_RES     ProtoOAPayloadType = 2180
	PROTO_OA_ORDER_DETAILS_REQ                ProtoOAPayloadType = 2181
	PROTO_OA_ORDER_DETAILS_RES                ProtoOAPayloadType = 2182
	PROTO_OA_ORDER_LIST_BY_POSITION_ID_REQ    ProtoOAPayloadType = 2183
	PROTO_OA_ORDER_LIST_BY_POSITION_ID_RES    ProtoOAPayloadType = 2184
	PROTO_OA_DEAL_OFFSET_LIST_REQ             ProtoOAPayloadType = 2185
	PROTO_OA_DEAL_OFFSET_LIST_RES             ProtoOAPayloadType = 2186
	PROTO_OA_GET_POSITION_UNREALIZED_PNL_REQ  ProtoOAPayloadType = 2187
	PROTO_OA_GET_POSITION_UNREALIZED_PNL_RES  ProtoOAPayloadType = 2188
)

/*
	Proto message derived types
*/
/**/

type ProtoOAOrderType = messages.ProtoOAOrderType

const (
	ProtoOAOrderType_MARKET                ProtoOAOrderType = 1
	ProtoOAOrderType_LIMIT                 ProtoOAOrderType = 2
	ProtoOAOrderType_STOP                  ProtoOAOrderType = 3
	ProtoOAOrderType_STOP_LOSS_TAKE_PROFIT ProtoOAOrderType = 4
	ProtoOAOrderType_MARKET_RANGE          ProtoOAOrderType = 5
	ProtoOAOrderType_STOP_LIMIT            ProtoOAOrderType = 6
)

type ProtoOATradeSide = messages.ProtoOATradeSide

const (
	ProtoOATradeSide_BUY  ProtoOATradeSide = 1
	ProtoOATradeSide_SELL ProtoOATradeSide = 2
)

type ProtoOATimeInForce = messages.ProtoOATimeInForce

const (
	ProtoOATimeInForce_GOOD_TILL_DATE      ProtoOATimeInForce = 1
	ProtoOATimeInForce_GOOD_TILL_CANCEL    ProtoOATimeInForce = 2
	ProtoOATimeInForce_IMMEDIATE_OR_CANCEL ProtoOATimeInForce = 3
	ProtoOATimeInForce_FILL_OR_KILL        ProtoOATimeInForce = 4
	ProtoOATimeInForce_MARKET_ON_OPEN      ProtoOATimeInForce = 5
)

type ProtoOAOrderTriggerMethod = messages.ProtoOAOrderTriggerMethod

const (
	ProtoOAOrderTriggerMethod_TRADE           ProtoOAOrderTriggerMethod = 1 // Stop Order: buy is triggered by ask, sell by bid; Stop Loss Order: for buy position is triggered by bid and for sell position by ask.
	ProtoOAOrderTriggerMethod_OPPOSITE        ProtoOAOrderTriggerMethod = 2 // Stop Order: buy is triggered by bid, sell by ask; Stop Loss Order: for buy position is triggered by ask and for sell position by bid.
	ProtoOAOrderTriggerMethod_DOUBLE_TRADE    ProtoOAOrderTriggerMethod = 3 // The same as TRADE, but trigger is checked after the second consecutive tick.
	ProtoOAOrderTriggerMethod_DOUBLE_OPPOSITE ProtoOAOrderTriggerMethod = 4 // The same as OPPOSITE, but trigger is checked after the second consecutive tick.
)

type ProtoOAAsset = messages.ProtoOAAsset

type ProtoOALightSymbol = messages.ProtoOALightSymbol

type ProtoOAArchivedSymbol = messages.ProtoOAArchivedSymbol

type ProtoOAAssetClass = messages.ProtoOAAssetClass

type ProtoOATrader = messages.ProtoOATrader

type ProtoOAAccessRights = messages.ProtoOAAccessRights

const (
	ProtoOAAccessRights_FULL_ACCESS ProtoOAAccessRights = 0 // Enable all trading.
	ProtoOAAccessRights_CLOSE_ONLY  ProtoOAAccessRights = 1 // Only closing trading request are enabled.
	ProtoOAAccessRights_NO_TRADING  ProtoOAAccessRights = 2 // View only access.
	ProtoOAAccessRights_NO_LOGIN    ProtoOAAccessRights = 3 // No access.
)

type ProtoOATotalMarginCalculationType = messages.ProtoOATotalMarginCalculationType

const (
	ProtoOATotalMarginCalculationType_MAX ProtoOATotalMarginCalculationType = 0
	ProtoOATotalMarginCalculationType_SUM ProtoOATotalMarginCalculationType = 1
	ProtoOATotalMarginCalculationType_NET ProtoOATotalMarginCalculationType = 2
)

type ProtoOAAccountType = messages.ProtoOAAccountType

const (
	ProtoOAAccountType_HEDGED         ProtoOAAccountType = 0 // Allows multiple positions on a trading account for a symbol.
	ProtoOAAccountType_NETTED         ProtoOAAccountType = 1 // Only one position per symbol is allowed on a trading account.
	ProtoOAAccountType_SPREAD_BETTING ProtoOAAccountType = 2 // Spread betting type account.
)

type ProtoOALimitedRiskMarginCalculationStrategy = messages.ProtoOALimitedRiskMarginCalculationStrategy

const (
	ProtoOALimitedRiskMarginCalculationStrategy_ACCORDING_TO_LEVERAGE         ProtoOALimitedRiskMarginCalculationStrategy = 0
	ProtoOALimitedRiskMarginCalculationStrategy_ACCORDING_TO_GSL              ProtoOALimitedRiskMarginCalculationStrategy = 1
	ProtoOALimitedRiskMarginCalculationStrategy_ACCORDING_TO_GSL_AND_LEVERAGE ProtoOALimitedRiskMarginCalculationStrategy = 2
)

type ProtoOAStopOutStrategy = messages.ProtoOAStopOutStrategy

const (
	ProtoOAStopOutStrategy_MOST_MARGIN_USED_FIRST ProtoOAStopOutStrategy = 0 //A Stop Out strategy that closes a Position with the largest Used Margin
	ProtoOAStopOutStrategy_MOST_LOSING_FIRST      ProtoOAStopOutStrategy = 1 //A Stop Out strategy that closes a Position with the least PnL
)

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

type ProtoOATrendbar = messages.ProtoOATrendbar

type ProtoOAQuoteType = messages.ProtoOAQuoteType

const (
	ProtoOAQuoteType_BID ProtoOAQuoteType = 1
	ProtoOAQuoteType_ASK ProtoOAQuoteType = 2
)

type ProtoOATickData = messages.ProtoOATickData

type ProtoOASymbolCategory = messages.ProtoOASymbolCategory

type ProtoOAMarginCall = messages.ProtoOAMarginCall

type ProtoOANotificationType = messages.ProtoOANotificationType

const (
	ProtoOANotificationType_MARGIN_LEVEL_THRESHOLD_1 ProtoOANotificationType = 61 // one of three margin calls, they are all similar.
	ProtoOANotificationType_MARGIN_LEVEL_THRESHOLD_2 ProtoOANotificationType = 62 // one of three margin calls, they are all similar.
	ProtoOANotificationType_MARGIN_LEVEL_THRESHOLD_3 ProtoOANotificationType = 63 // one of three margin calls, they are all similar.
)

type ProtoOADynamicLeverage = messages.ProtoOADynamicLeverage

type ProtoOADynamicLeverageTier = messages.ProtoOADynamicLeverageTier

type ProtoOADeal = messages.ProtoOADeal

type ProtoOADealStatus = messages.ProtoOADealStatus

const (
	ProtoOADealStatus_FILLED              ProtoOADealStatus = 2 // Deal filled.
	ProtoOADealStatus_PARTIALLY_FILLED    ProtoOADealStatus = 3 // Deal is partially filled.
	ProtoOADealStatus_REJECTED            ProtoOADealStatus = 4 // Deal is correct but was rejected by liquidity provider (e.g. no liquidity).
	ProtoOADealStatus_INTERNALLY_REJECTED ProtoOADealStatus = 5 // Deal rejected by server (e.g. no price quotes).
	ProtoOADealStatus_ERROR               ProtoOADealStatus = 6 // Deal is rejected by LP due to error (e.g. symbol is unknown).
	ProtoOADealStatus_MISSED              ProtoOADealStatus = 7 // Liquidity provider did not sent response on the deal during specified execution time period.
)

type ProtoOAClosePositionDetail = messages.ProtoOAClosePositionDetail

type ProtoOAOrder = messages.ProtoOAOrder

type ProtoOATradeData = messages.ProtoOATradeData

type ProtoOAOrderStatus = messages.ProtoOAOrderStatus

const (
	ProtoOAOrderStatus_ORDER_STATUS_ACCEPTED  ProtoOAOrderStatus = 1 // Order request validated and accepted for execution.
	ProtoOAOrderStatus_ORDER_STATUS_FILLED    ProtoOAOrderStatus = 2 // Order is fully filled.
	ProtoOAOrderStatus_ORDER_STATUS_REJECTED  ProtoOAOrderStatus = 3 // Order is rejected due to validation.
	ProtoOAOrderStatus_ORDER_STATUS_EXPIRED   ProtoOAOrderStatus = 4 // Order expired. Might be valid for orders with partially filled volume that were expired on LP.
	ProtoOAOrderStatus_ORDER_STATUS_CANCELLED ProtoOAOrderStatus = 5 // Order is cancelled. Might be valid for orders with partially filled volume that were cancelled by LP.
)

type ProtoOADealOffset = messages.ProtoOADealOffset

type ProtoOAPositionUnrealizedPnL = messages.ProtoOAPositionUnrealizedPnL

type ProtoOADepthQuote = messages.ProtoOADepthQuote

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

type ProtoOAPosition = messages.ProtoOAPosition

type ProtoOAPositionStatus = messages.ProtoOAPositionStatus

const (
	ProtoOAPositionStatus_POSITION_STATUS_OPEN    ProtoOAPositionStatus = 1
	ProtoOAPositionStatus_POSITION_STATUS_CLOSED  ProtoOAPositionStatus = 2
	ProtoOAPositionStatus_POSITION_STATUS_CREATED ProtoOAPositionStatus = 3 // Empty position is created for pending order.
	ProtoOAPositionStatus_POSITION_STATUS_ERROR   ProtoOAPositionStatus = 4
)

type ProtoOABonusDepositWithdraw = messages.ProtoOABonusDepositWithdraw

type ProtoOAChangeBonusType = messages.ProtoOAChangeBonusType

const (
	ProtoOAChangeBonusType_BONUS_DEPOSIT  ProtoOAChangeBonusType = 0
	ProtoOAChangeBonusType_BONUS_WITHDRAW ProtoOAChangeBonusType = 1
)

type ProtoOADepositWithdraw = messages.ProtoOADepositWithdraw

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

/*
Accessory functions
*/

var trendbarPeriods []ProtoOATrendbarPeriod = []ProtoOATrendbarPeriod{
	ProtoOATrendbarPeriod_M1,
	ProtoOATrendbarPeriod_M2,
	ProtoOATrendbarPeriod_M3,
	ProtoOATrendbarPeriod_M4,
	ProtoOATrendbarPeriod_M5,
	ProtoOATrendbarPeriod_M10,
	ProtoOATrendbarPeriod_M15,
	ProtoOATrendbarPeriod_M30,
	ProtoOATrendbarPeriod_H1,
	ProtoOATrendbarPeriod_H4,
	ProtoOATrendbarPeriod_H12,
	ProtoOATrendbarPeriod_D1,
	ProtoOATrendbarPeriod_W1,
	ProtoOATrendbarPeriod_MN1,
}

func GetTrendbarPeriods() []ProtoOATrendbarPeriod {
	cpy := make([]ProtoOATrendbarPeriod, len(trendbarPeriods))
	copy(cpy, trendbarPeriods)
	return cpy
}

/*
	Proto Message types
*/
/**/

// * Response to the ProtoOAApplicationAuthRes request.
type ProtoOAAccountAuthRes = messages.ProtoOAAccountAuthRes

// * Request for getting the proxy version. Can be used to check the current version of the Open API scheme.
type ProtoOAVersionReq = messages.ProtoOAVersionReq

// * Response to the ProtoOAVersionReq request.
type ProtoOAVersionRes = messages.ProtoOAVersionRes

// * Request for sending a new trading order. Allowed only if the accessToken has the "trade" permissions for the trading account.
type ProtoOANewOrderReq = messages.ProtoOANewOrderReq

// * Request for cancelling existing pending order. Allowed only if the accessToken has "trade" permissions for the trading account.
type ProtoOACancelOrderReq = messages.ProtoOACancelOrderReq

// * Request for amending the existing pending order. Allowed only if the Access Token has "trade" permissions for the trading account.
type ProtoOAAmendOrderReq = messages.ProtoOAAmendOrderReq

// * Request for amending StopLoss and TakeProfit of existing position. Allowed only if the accessToken has "trade" permissions for the trading account.
type ProtoOAAmendPositionSLTPReq = messages.ProtoOAAmendPositionSLTPReq

// * Request for closing or partially closing of an existing position. Allowed only if the accessToken has "trade" permissions for the trading account.
type ProtoOAClosePositionReq = messages.ProtoOAClosePositionReq

// * Request for the list of assets available for a trader's account.
type ProtoOAAssetListReq = messages.ProtoOAAssetListReq

// * Response to the ProtoOAAssetListReq request.
type ProtoOAAssetListRes = messages.ProtoOAAssetListRes

// * Request for a list of symbols available for a trading account. Symbol entries are returned with the limited set of fields.
type ProtoOASymbolsListReq = messages.ProtoOASymbolsListReq

// * Response to the ProtoOASymbolsListReq request.
type ProtoOASymbolsListRes = messages.ProtoOASymbolsListRes

// * Request for getting a full symbol entity.
type ProtoOASymbolByIdReq = messages.ProtoOASymbolByIdReq

// * Response to the ProtoOASymbolByIdReq request.
type ProtoOASymbolByIdRes = messages.ProtoOASymbolByIdRes

// * Request for getting a conversion chain between two assets that consists of several symbols. Use when no direct quote is available.
type ProtoOASymbolsForConversionReq = messages.ProtoOASymbolsForConversionReq

// * Response to the ProtoOASymbolsForConversionReq request.
type ProtoOASymbolsForConversionRes = messages.ProtoOASymbolsForConversionRes

// * Request for a list of asset classes available for the trader's account.
type ProtoOAAssetClassListReq = messages.ProtoOAAssetClassListReq

// * Response to the ProtoOAAssetListReq request.
type ProtoOAAssetClassListRes = messages.ProtoOAAssetClassListRes

// * Request for getting data of Trader's Account.
type ProtoOATraderReq = messages.ProtoOATraderReq

// * Response to the ProtoOATraderReq request.
type ProtoOATraderRes = messages.ProtoOATraderRes

// * Request for getting Trader's current open positions and pending orders data.
type ProtoOAReconcileReq = messages.ProtoOAReconcileReq

// * The response to the ProtoOAReconcileReq request.
type ProtoOAReconcileRes = messages.ProtoOAReconcileRes

// * Request for getting Trader's deals historical data (execution details).
type ProtoOADealListReq = messages.ProtoOADealListReq

// * The response to the ProtoOADealListRes request.
type ProtoOADealListRes = messages.ProtoOADealListRes

// * Request for getting Trader's orders filtered by timestamp.
type ProtoOAOrderListReq = messages.ProtoOAOrderListReq

// * The response to the ProtoOAOrderListReq request.
type ProtoOAOrderListRes = messages.ProtoOAOrderListRes

// * Request for getting the margin estimate according to leverage profiles. Can be used before sending a new order request.
// This doesn't consider ACCORDING_TO_GSL margin calculation type, as this calculation is trivial: usedMargin = (VWAP price
// of the position - GSL price) * volume * Quote2Deposit.
type ProtoOAExpectedMarginReq = messages.ProtoOAExpectedMarginReq

// * The response to the ProtoOAExpectedMarginReq request.
type ProtoOAExpectedMarginRes = messages.ProtoOAExpectedMarginRes

// * Request for getting Trader's historical data of deposits and withdrawals.
type ProtoOACashFlowHistoryListReq = messages.ProtoOACashFlowHistoryListReq

// * Response to the ProtoOACashFlowHistoryListReq request.
type ProtoOACashFlowHistoryListRes = messages.ProtoOACashFlowHistoryListRes

// * Request for getting the list of granted trader's account for the access token.
type ProtoOAGetAccountListByAccessTokenReq = messages.ProtoOAGetAccountListByAccessTokenReq

// * Response to the ProtoOAGetAccountListByAccessTokenReq request.
type ProtoOAGetAccountListByAccessTokenRes = messages.ProtoOAGetAccountListByAccessTokenRes

// * Response to the ProtoOARefreshTokenReq request.
type ProtoOARefreshTokenRes = messages.ProtoOARefreshTokenRes

// * Request for getting historical trend bars for the symbol.
type ProtoOAGetTrendbarsReq = messages.ProtoOAGetTrendbarsReq

// * Response to the ProtoOAGetTrendbarsReq request.
type ProtoOAGetTrendbarsRes = messages.ProtoOAGetTrendbarsRes

// * Request for getting historical tick data for the symbol.
type ProtoOAGetTickDataReq = messages.ProtoOAGetTickDataReq

// * Response to the ProtoOAGetTickDataReq request.
type ProtoOAGetTickDataRes = messages.ProtoOAGetTickDataRes

/*
	End quotes section
*/
/**/

// * Request for getting details of Trader's profile. Limited due to GDRP requirements.
type ProtoOAGetCtidProfileByTokenReq = messages.ProtoOAGetCtidProfileByTokenReq

// * Response to the ProtoOAGetCtidProfileByTokenReq request.
type ProtoOAGetCtidProfileByTokenRes = messages.ProtoOAGetCtidProfileByTokenRes

// * Request for a list of symbol categories available for a trading account.
type ProtoOASymbolCategoryListReq = messages.ProtoOASymbolCategoryListReq

// * Response to the ProtoSymbolCategoryListReq request.
type ProtoOASymbolCategoryListRes = messages.ProtoOASymbolCategoryListRes

// * Response to the ProtoOATraderLogoutReq request. Actual logout of trading account will be completed on ProtoOAAccountDisconnectEvent.
type ProtoOAAccountLogoutRes = messages.ProtoOAAccountLogoutRes

// * Request for a list of existing margin call thresholds configured for a user.
type ProtoOAMarginCallListReq = messages.ProtoOAMarginCallListReq

// * Response with a list of existing user Margin Calls, usually contains 3 items.
type ProtoOAMarginCallListRes = messages.ProtoOAMarginCallListRes

// * Request to modify marginLevelThreshold of specified marginCallType for ctidTraderAccountId.
type ProtoOAMarginCallUpdateReq = messages.ProtoOAMarginCallUpdateReq

// * If this response received, it means that margin call was successfully updated.
type ProtoOAMarginCallUpdateRes = messages.ProtoOAMarginCallUpdateRes

// * Request for getting a dynamic leverage entity referenced in ProtoOASymbol.leverageId.
type ProtoOAGetDynamicLeverageByIDReq = messages.ProtoOAGetDynamicLeverageByIDReq

// * Response to the ProtoOAGetDynamicLeverageByIDReq request.
type ProtoOAGetDynamicLeverageByIDRes = messages.ProtoOAGetDynamicLeverageByIDRes

// * Request for retrieving the deals related to a position.
type ProtoOADealListByPositionIdReq = messages.ProtoOADealListByPositionIdReq

// * Response to the ProtoOADealListByPositionIdReq request.
type ProtoOADealListByPositionIdRes = messages.ProtoOADealListByPositionIdRes

// * Request for getting Order and its related Deals.
type ProtoOAOrderDetailsReq = messages.ProtoOAOrderDetailsReq

// * Response to the ProtoOAOrderDetailsReq request.
type ProtoOAOrderDetailsRes = messages.ProtoOAOrderDetailsRes

// * Request for retrieving Orders related to a Position by using Position ID. Filtered by utcLastUpdateTimestamp.
type ProtoOAOrderListByPositionIdReq = messages.ProtoOAOrderListByPositionIdReq

// * Response to ProtoOAOrderListByPositionIdReq request.
type ProtoOAOrderListByPositionIdRes = messages.ProtoOAOrderListByPositionIdRes

// * Request for getting sets of Deals that were offset by a specific Deal and that are offsetting the Deal.
type ProtoOADealOffsetListReq = messages.ProtoOADealOffsetListReq

// * Response for ProtoOADealOffsetListReq.
type ProtoOADealOffsetListRes = messages.ProtoOADealOffsetListRes

// * Request for getting trader's positions' unrealized PnLs.
type ProtoOAGetPositionUnrealizedPnLReq = messages.ProtoOAGetPositionUnrealizedPnLReq

// * Response to ProtoOAGetPositionUnrealizedPnLReq request.
type ProtoOAGetPositionUnrealizedPnLRes = messages.ProtoOAGetPositionUnrealizedPnLRes

/*
	Subscription based event quotes section
*/
/**/

type ProtoOASpotEvent = messages.ProtoOASpotEvent

type ProtoOADepthEvent = messages.ProtoOADepthEvent

/*
	Listenable event quotes section
*/
/**/

// * Event that is sent when a new spot event is generated on the server side. Requires subscription on the spot events, see
// ProtoOASubscribeSpotsReq. First event, received after subscription will contain latest spot prices even if market is closed.
type ProtoOATrailingSLChangedEvent = messages.ProtoOATrailingSLChangedEvent

// * Event that is sent when the symbol is changed on the Server side.
type ProtoOASymbolChangedEvent = messages.ProtoOASymbolChangedEvent

// * Event that is sent when a Trader is updated on Server side.
type ProtoOATraderUpdatedEvent = messages.ProtoOATraderUpdatedEvent

// * Event that is sent following the successful order acceptance or execution by the server. Acts as response to the
// ProtoOANewOrderReq, ProtoOACancelOrderReq, ProtoOAAmendOrderReq, ProtoOAAmendPositionSLTPReq, ProtoOAClosePositionReq requests.
// Also, the event is sent when a Deposit/Withdrawal took place.
type ProtoOAExecutionEvent = messages.ProtoOAExecutionEvent

// * Event that is sent when errors occur during the order requests.
type ProtoOAOrderErrorEvent = messages.ProtoOAOrderErrorEvent

// * Event that is sent when the margin allocated to a specific position is changed.
type ProtoOAMarginChangedEvent = messages.ProtoOAMarginChangedEvent

// * Event that is sent when a session to a specific trader's account is terminated by the server but the existing connections
// with the other trader's accounts are maintained. Reasons to trigger: account was deleted, cTID was deleted, token was refreshed,
// token was revoked.
type ProtoOAAccountsTokenInvalidatedEvent = messages.ProtoOAAccountsTokenInvalidatedEvent

// * Event that is sent when the connection with the client application is cancelled by the server. All the sessions for the traders'
// accounts will be terminated.
type ProtoOAClientDisconnectEvent = messages.ProtoOAClientDisconnectEvent

// * Event that is sent when the established session for an account is dropped on the server side. A new session must be authorized
// for the account.
type ProtoOAAccountDisconnectEvent = messages.ProtoOAAccountDisconnectEvent

// * Event that is sent when a Margin Call threshold configuration is updated.
type ProtoOAMarginCallUpdateEvent = messages.ProtoOAMarginCallUpdateEvent

// * Event that is sent when account margin level reaches target marginLevelThreshold. Event is sent no more than once every 10 minutes
// to avoid spamming.
type ProtoOAMarginCallTriggerEvent = messages.ProtoOAMarginCallTriggerEvent
