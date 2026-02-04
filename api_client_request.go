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
	"github.com/linuskuehnle/ctraderopenapi/internal/datatypes"
	"github.com/linuskuehnle/ctraderopenapi/internal/messages"

	"context"
	"errors"
	"fmt"
	"sync"

	"google.golang.org/protobuf/proto"
)

func (c *apiClient) AuthenticateAccount(ctid CtraderAccountId, accessToken AccessToken) (*ProtoOAAccountAuthRes, error) {
	req := messages.ProtoOAAccountAuthReq{
		CtidTraderAccountId: proto.Int64(int64(ctid)),
		AccessToken:         proto.String(string(accessToken)),
	}
	var res ProtoOAAccountAuthRes

	reqCtx := context.Background()

	reqData := RequestData{
		Ctx: reqCtx,
		Req: &req,
		Res: &res,
	}

	if err := c.SendRequest(reqData); err != nil {
		return nil, err
	}

	if !c.accManager.HasAccessToken(accessToken) {
		c.accManager.AddAccessToken(accessToken)
	}
	c.accManager.AddAccountId(accessToken, ctid)

	return &res, nil
}

func (c *apiClient) LogoutAccount(ctid CtraderAccountId, waitForConfirm bool) (*ProtoOAAccountLogoutRes, error) {
	req := messages.ProtoOAAccountLogoutReq{
		CtidTraderAccountId: proto.Int64(int64(ctid)),
	}
	var res ProtoOAAccountLogoutRes

	reqCtx := context.Background()

	reqData := RequestData{
		Ctx: reqCtx,
		Req: &req,
		Res: &res,
	}

	wg := sync.WaitGroup{}

	if waitForConfirm {
		wg.Go(func() {
			// Wait for the disconnect event confirming the logout
			c.accManager.WaitForAccDisconnectConfirm(ctid)
		})
	}

	if err := c.SendRequest(reqData); err != nil {
		return nil, err
	}

	c.accManager.RemoveAccountId(ctid)

	wg.Wait()
	return &res, nil
}

func (c *apiClient) RefreshAccessToken(expiredToken AccessToken, refreshToken RefreshToken) (*ProtoOARefreshTokenRes, error) {
	req := messages.ProtoOARefreshTokenReq{
		RefreshToken: proto.String(string(refreshToken)),
	}
	var res ProtoOARefreshTokenRes

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	c.accManager.LockModification(ctx)

	if !c.accManager.HasAccessToken(expiredToken) {
		return nil, &AccessTokenDoesNotExistError{
			AccessToken: expiredToken,
		}
	}

	reqData := RequestData{
		Ctx: context.Background(),
		Req: &req,
		Res: &res,
	}

	if err := c.SendRequest(reqData); err != nil {
		return nil, err
	}

	c.accManager.UpdateAccessToken(expiredToken, AccessToken(res.GetAccessToken()))
	return &res, nil
}

func (c *apiClient) SendRequest(reqData RequestData) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if blockPassed := c.lifecycleData.BlockUntilReconnected(reqData.Ctx); !blockPassed {
		return &RequestContextExpiredError{
			Err: reqData.Ctx.Err(),
		}
	}

	return c.sendRequest(reqData)
}

func (c *apiClient) sendRequest(reqData RequestData) error {
	if err := reqData.CheckError(); err != nil {
		return &FunctionInvalidArgError{
			FunctionName: "SendRequest",
			Err:          err,
		}
	}

	if !c.lifecycleData.IsRunning() {
		return &LifeCycleNotRunningError{
			CallContext: "SendRequest",
		}
	}

	if reqData.Ctx == nil {
		reqData.Ctx = context.Background()
	}
	ctx, cancelCtx := context.WithTimeout(reqData.Ctx, c.cfg.requestTimeout)
	reqData.Ctx = ctx
	defer cancelCtx()

	errCh := make(chan error)
	heapErrCh := make(chan error)
	resDataCh := make(chan *datatypes.ResponseData)

	metaData, err := datatypes.NewRequestMetaData(&reqData, errCh, heapErrCh, resDataCh)
	if err != nil {
		return err
	}

	// Add the request meta data to the request heap
	if err := c.requestHeap.AddNode(metaData); err != nil {
		return err
	}

	// Enqueue the request
	if err := c.enqueueRequest(metaData); err != nil {
		return err
	}

	err = nil
	var heapErrChClosed, errChClosed bool
	for !heapErrChClosed || !errChClosed {
		select {
		case reqCtxExpiredErr, ok := <-heapErrCh:
			if !ok {
				heapErrChClosed = true
				continue
			}
			err = reqCtxExpiredErr
		case reqErr, ok := <-errCh:
			if !ok {
				errChClosed = true
				continue
			}
			err = reqErr
		}

		close(resDataCh)
		return err
	}

	resData := <-resDataCh

	payload := resData.ProtoMsg.GetPayload()

	if err := checkResponseForError(payload, resData.PayloadType); err != nil {
		var resErr *ResponseError
		/*
			Retry requests on either of those conditions:
			- The error is not a ResponseError, meaning it is a transport level error
			- The error is a ResponseError with error code "server side rate limit hit"
		*/
		if !errors.As(err, &resErr) || resErr.ErrorCode == resErrorCode_serverSideRateLimitHit {
			// Set rate limiter penalty
			rateLimitType := rateLimitTypeByReqType[metaData.Req.GetOAType()]
			c.rateLimiters[rateLimitType].SetPenalty(rateLimitInterval)

			// Execute request again
			return c.sendRequest(reqData)
		}
		return err
	}

	if resData.PayloadType != reqData.Res.GetOAType() {
		return fmt.Errorf("unexpected response payload type: got %d, expected %d",
			resData.PayloadType, reqData.Res.GetOAType(),
		)
	}

	// Unmarshal payload into provided response struct
	if err := proto.Unmarshal(payload, reqData.Res); err != nil {
		return &ProtoUnmarshalError{
			CallContext: fmt.Sprintf("proto response [%d]", reqData.Res.GetOAType()),
			Err:         err,
		}
	}

	return nil
}

func (c *apiClient) authenticateApp() error {
	reqData := RequestData{
		Ctx: context.Background(),
		Req: &messages.ProtoOAApplicationAuthReq{
			ClientId:     proto.String(c.cred.ClientId),
			ClientSecret: proto.String(c.cred.ClientSecret),
		},
		Res: &messages.ProtoOAApplicationAuthRes{},
	}

	err := c.sendRequest(reqData)
	if err != nil {
		var resErr *ResponseError
		if !errors.As(err, &resErr) || resErr.ErrorCode != resErrorCode_appAlreadyAuthenticated {
			// err either is not a ResponseError or it is, but error code does not indicate tolerable error
			// is already authenticated. So we need to return it.
			return err
		}
	}

	// Either err is nil or err indicates app has already been authenticated.
	return nil
}

func (c *apiClient) emitHeartbeat() error {
	errCh := make(chan error)

	reqData := RequestData{
		Ctx: context.Background(),
		Req: &messages.ProtoHeartbeatEvent{},
	}

	metaData, err := datatypes.NewRequestMetaData(&reqData, errCh, nil, nil)
	if err != nil {
		return err
	}

	if err := c.enqueueRequest(metaData); err != nil {
		return err
	}

	if err, ok := <-errCh; ok {
		return err
	}
	return nil
}

func (c *apiClient) reloginActiveAccounts() error {
	wg := sync.WaitGroup{}

	tokenByCtid := c.accManager.GetAccessTokenByAccountId()
	errCh := make(chan error, len(tokenByCtid))

	for ctid, accessToken := range tokenByCtid {
		wg.Add(1)

		go func(ctid CtraderAccountId, accessToken AccessToken) {
			defer wg.Done()

			req := messages.ProtoOAAccountAuthReq{
				CtidTraderAccountId: proto.Int64(int64(ctid)),
				AccessToken:         proto.String(string(accessToken)),
			}
			var res ProtoOAAccountAuthRes

			reqCtx := context.Background()

			reqData := RequestData{
				Ctx: reqCtx,
				Req: &req,
				Res: &res,
			}

			if err := c.sendRequest(reqData); err != nil {
				errCh <- err
				return
			}
		}(ctid, accessToken)
	}

	wg.Wait()
	close(errCh)

	if err, ok := <-errCh; ok {
		return err
	}
	return nil
}

func (c *apiClient) resubscribeActiveSubs() error {
	ctids := c.accManager.GetAllAccountIds()

	wg := sync.WaitGroup{}
	errCh := make(chan error, len(ctids))

	for _, ctid := range ctids {
		wg.Add(1)

		go func(ctid CtraderAccountId) {
			defer wg.Done()

			if err := c.resubscribeAccountSubs(ctid); err != nil {
				errCh <- err
				return
			}
		}(ctid)
	}

	wg.Wait()
	close(errCh)

	if err, ok := <-errCh; ok {
		return err
	}
	return nil
}

func (c *apiClient) resubscribeAccountSubs(ctid CtraderAccountId) error {
	// Filter out live trendbar events as they have to be subscribed after
	// successful spot event subscription.
	liveTrendbarEventData := []APIEventSubData{}
	otherEventData := []APIEventSubData{}

	currentSubs, err := c.accManager.GetEventSubscriptionsOfAccountId(ctid)
	if err != nil {
		return err
	}

	for t, s := range currentSubs {
		d := APIEventSubData{
			EventType:       t,
			SubcriptionData: s,
		}

		if t == APIEventType_LiveTrendbars {
			liveTrendbarEventData = append(liveTrendbarEventData, d)
		} else {
			otherEventData = append(otherEventData, d)
		}
	}

	resubscribeBatch := func(eventDataBatch []APIEventSubData) error {
		// Resubscribe current sessions active subscriptions
		wg := sync.WaitGroup{}
		errCh := make(chan error, len(eventDataBatch))

		for _, e := range eventDataBatch {
			wg.Add(1)
			go func(e APIEventSubData) {
				defer wg.Done()

				if _, err := c.subscribeAPIEvent(e, true); err != nil {
					var resErr *ResponseError
					if !errors.As(err, &resErr) || resErr.ErrorCode != resErrorCode_apiEventAlreadySubscribed {
						// Either err is not a ResponseError or it is, but error code does not indicate tolerable error
						// is already subscribed. So we need to return it.
						errCh <- err
						return
					}
				}
				// Successfully resubscribed
			}(e)
		}

		wg.Wait()
		close(errCh)

		if err, ok := <-errCh; ok {
			return err
		}
		return nil
	}

	if err := resubscribeBatch(otherEventData); err != nil {
		return err
	}
	if err := resubscribeBatch(liveTrendbarEventData); err != nil {
		return err
	}
	return nil
}

func checkResponseForError(payloadBytes []byte, payloadType protoOAPayloadType) error {
	if payloadType != proto_OA_ERROR_RES {
		return nil
	}

	var errorMsg messages.ProtoOAErrorRes
	if err := proto.Unmarshal(payloadBytes, &errorMsg); err != nil {
		return &ProtoUnmarshalError{
			CallContext: "proto OA error response",
			Err:         err,
		}
	}

	return &ResponseError{
		ErrorCode:               errorMsg.GetErrorCode(),
		Description:             errorMsg.GetDescription(),
		MaintenanceEndTimestamp: errorMsg.GetMaintenanceEndTimestamp(),
		RetryAfter:              errorMsg.GetRetryAfter(),
	}
}
