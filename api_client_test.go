// Copyright 2025 Linus Kühnle
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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
	"sync"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"
)

func TestClientConnectDisconnect(t *testing.T) {
	c, err := createApiClient(Environment_Demo)
	if err != nil {
		t.Fatal(err)
	}

	if err = c.Connect(); err != nil {
		t.Fatalf("error connecting (1): %v", err)
	}
	c.Disconnect()

	if err = c.Connect(); err != nil {
		t.Fatalf("error connecting (2): %v", err)
	}
	if err = c.Connect(); err == nil {
		t.Fatalf("no error connecting after already connected. expected error")
	}
	c.Disconnect()
}

func TestClientConnectHeartbeatDisconnect(t *testing.T) {
	c, err := createApiClient(Environment_Demo)
	if err != nil {
		t.Fatal(err)
	}

	if err = c.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}
	defer c.Disconnect()

	if err := c.emitHeartbeat(); err != nil {
		t.Fatalf("error emitting heartbeat: %v", err)
	}
}

func TestClientReqConcurrency(t *testing.T) {
	c, err := createApiClient(Environment_Demo)
	if err != nil {
		t.Fatal(err)
	}

	if err = c.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}

	// Create requests and send them concurrently to check if each request makes a proper round trip
	numRequests := 20
	wg := sync.WaitGroup{}
	wg.Add(numRequests)

	var reqErrs []error = make([]error, numRequests)

	for i := range numRequests {
		go func(i int) {
			defer wg.Done()

			req := ProtoOAVersionReq{
				PayloadType: PROTO_OA_VERSION_REQ.Enum(),
			}
			var res ProtoOAVersionRes

			reqCtx := context.Background()

			reqData := RequestData{
				Ctx:     reqCtx,
				ReqType: PROTO_OA_VERSION_REQ,
				Req:     &req,
				ResType: PROTO_OA_VERSION_RES,
				Res:     &res,
			}

			if err := c.SendRequest(reqData); err != nil {
				reqErrs[i] = err
				return
			}
		}(i)
	}

	// Wait for all requests to finish
	wg.Wait()

	c.Disconnect()

	// Check for request errors
	for _, e := range reqErrs {
		if e != nil {
			t.Fatalf("request error: %v", e)
		}
	}
}

func TestProtoOAErrorResponse(t *testing.T) {
	accountId, _, env, err := loadTestAccountCredentials()
	if err != nil {
		t.Fatalf("error loading test account credentials: %v", err)
	}

	c, err := createApiClient(env)
	if err != nil {
		t.Fatal(err)
	}

	if err = c.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}
	defer c.Disconnect()

	// Make unauthenticated request that will return a proto error response
	req := ProtoOAGetTrendbarsReq{
		CtidTraderAccountId: proto.Int64(accountId),
		Period:              ProtoOATrendbarPeriod_D1.Enum(),
		SymbolId:            proto.Int64(1),
	}
	var res ProtoOAGetTrendbarsRes

	reqCtx := context.Background()

	reqData := RequestData{
		Ctx:     reqCtx,
		ReqType: PROTO_OA_GET_TRENDBARS_REQ,
		Req:     &req,
		ResType: PROTO_OA_GET_TRENDBARS_RES,
		Res:     &res,
	}

	if err = c.SendRequest(reqData); err != nil {
		var respErr *ResponseError
		if !errors.As(err, &respErr) {
			t.Fatalf("unexpected error of type RequestError. got: %v", err)
		}
		if respErr.ErrorCode != messages.ProtoErrorCode_name[int32(messages.ProtoErrorCode_INVALID_REQUEST)] {
			t.Fatalf("expected error code INVALID_REQUEST, got: %v", respErr.ErrorCode)
		}
		if respErr.Description != "Trading account is not authorized" {
			t.Fatalf("expected error description 'Trading account is not authorized', got: %v", respErr.Description)
		}
	} else {
		t.Fatal("expected RequestError error return value")
	}
}

func TestListenToFatalError(t *testing.T) {
	c, err := createApiClient(Environment_Demo)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	var onceFatalErr sync.Once

	errCh := make(chan error)

	fatalErrCh := make(chan ClientEvent)
	onFatalErr := func(e *FatalErrorEvent) {
		fatalErr := e.Err
		if fatalErr.Error() != "test fatal error" {
			t.Fatalf("expected fatal error 'test fatal error', got: %v", fatalErr)
		}
		onceFatalErr.Do(func() {
			errCh <- e.Err
			close(errCh)
		})
	}
	if err := c.ListenToClientEvent(ctx, ClientEventListenData{ClientEventType_FatalErrorEvent, fatalErrCh}); err != nil {
		t.Fatalf("error listening to fatal error client event: %v", err)
	}
	if err := SpawnClientEventHandler(ctx, fatalErrCh, onFatalErr); err != nil {
		t.Fatalf("error spawning fatal error client event handler: %v", err)
	}

	var onceRecSucc sync.Once
	wg := sync.WaitGroup{}
	wg.Add(1)

	recSuccCh := make(chan ClientEvent)
	onRecSucc := func(e *ReconnectSuccessEvent) {
		onceRecSucc.Do(func() {
			wg.Done()
		})
	}
	if err := c.ListenToClientEvent(ctx, ClientEventListenData{ClientEventType_ReconnectSuccessEvent, recSuccCh}); err != nil {
		t.Fatalf("error listening to reconnect success client event: %v", err)
	}
	if err := SpawnClientEventHandler(ctx, recSuccCh, onRecSucc); err != nil {
		t.Fatalf("error spawning reconnect success client event handler: %v", err)
	}

	if err = c.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}
	defer c.Disconnect()

	emitErrText := "test fatal error"

	// Manually trigger a fatal error
	c.fatalErrCh <- errors.New(emitErrText)

	err, ok := <-errCh
	if !ok {
		t.Fatal("expected error, received no error before errCh close")
	}
	if receivedErrText := err.Error(); receivedErrText != emitErrText {
		t.Fatalf("expected error text\"%s\", received %s", emitErrText, receivedErrText)
	}

	wg.Wait()
	cancel()
}

func TestHeartbeatEmission(t *testing.T) {
	c, err := createApiClient(Environment_Demo)
	if err != nil {
		t.Fatal(err)
	}

	if err = c.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}
	defer c.Disconnect()

	// Test if any panics occur after emitting the heartbeat event
	time.Sleep(time.Second * (Heartbeat_Timeout_Seconds + 2))
}

func TestAccountAuth(t *testing.T) {
	accountId, accessToken, env, err := loadTestAccountCredentials()
	if err != nil {
		t.Fatalf("error loading test account credentials: %v", err)
	}

	c, err := createApiClient(env)
	if err != nil {
		t.Fatal(err)
	}

	if err = c.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}
	defer c.Disconnect()

	res, err := c.AuthenticateAccount(CtraderAccountId(accountId), AccessToken(accessToken))
	if err != nil {
		t.Fatal(err)
	}

	if res.GetCtidTraderAccountId() != accountId {
		t.Fatalf("expected account ID %d, got %d", accountId, res.GetCtidTraderAccountId())
	}

	// Test account disconnect confirm
	if _, err = c.LogoutAccount(CtraderAccountId(accountId), true); err != nil {
		t.Fatalf("error logging out out of account: %v", err)
	}
}

func TestGetSymbols(t *testing.T) {
	accountId, accessToken, env, err := loadTestAccountCredentials()
	if err != nil {
		t.Fatalf("error loading test account credentials: %v", err)
	}
	c, err := createApiClient(env)
	if err != nil {
		t.Fatal(err)
	}

	if err = c.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}
	defer c.Disconnect()

	if _, err = c.AuthenticateAccount(CtraderAccountId(accountId), AccessToken(accessToken)); err != nil {
		t.Fatal(err)
	}

	reqData := RequestData{
		ReqType: PROTO_OA_SYMBOLS_LIST_REQ,
		Req: &ProtoOASymbolsListReq{
			CtidTraderAccountId:    proto.Int64(accountId),
			IncludeArchivedSymbols: proto.Bool(false),
		},
		ResType: PROTO_OA_SYMBOLS_LIST_RES,
		Res:     &ProtoOASymbolsListRes{},
	}

	if err = c.SendRequest(reqData); err != nil {
		t.Fatalf("error sending request: %v", err)
	}
}

func TestReconnect(t *testing.T) {
	c, err := createApiClient(Environment_Demo)
	if err != nil {
		t.Fatal(err)
	}

	var once1 sync.Once
	var once2 sync.Once
	wg := sync.WaitGroup{}
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())

	onConnLossCh := make(chan ClientEvent)
	onReconnectSuccessCh := make(chan ClientEvent)
	onReconnectFailCh := make(chan ClientEvent)

	onConnectionLoss := func(e *ConnectionLossEvent) {
		once1.Do(func() {
			t.Log("closed tcp conn. waiting for reconnect...")
		})
	}
	onReconnectSuccess := func(e *ReconnectSuccessEvent) {
		// On reconnect success
		t.Log("reconnect successful")

		once2.Do(func() {
			wg.Done()
		})
	}
	onReconnectFail := func(e *ReconnectFailEvent) {
		// On reconnect failure
		t.Logf("reconnect failed: %v", e.Err)
	}

	if err := c.ListenToClientEvent(ctx, ClientEventListenData{ClientEventType_ConnectionLossEvent, onConnLossCh}); err != nil {
		t.Fatalf("error listening to connection loss events: %v", err)
	}
	if err := c.ListenToClientEvent(ctx, ClientEventListenData{ClientEventType_ReconnectSuccessEvent, onReconnectSuccessCh}); err != nil {
		t.Fatalf("error listening to reconnect success events: %v", err)
	}
	if err := c.ListenToClientEvent(ctx, ClientEventListenData{ClientEventType_ReconnectFailEvent, onReconnectFailCh}); err != nil {
		t.Fatalf("error listening to reconnect failure events: %v", err)
	}

	if err := SpawnClientEventHandler(ctx, onConnLossCh, onConnectionLoss); err != nil {
		t.Fatalf("error spawning connection loss event handler: %v", err)
	}
	if err := SpawnClientEventHandler(ctx, onReconnectSuccessCh, onReconnectSuccess); err != nil {
		t.Fatalf("error spawning reconnect success event handler: %v", err)
	}
	if err := SpawnClientEventHandler(ctx, onReconnectFailCh, onReconnectFail); err != nil {
		t.Fatalf("error spawning reconnect failure event handler: %v", err)
	}

	if err = c.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}
	defer c.Disconnect()

	c.tcpClient.JustCloseConn()

	// Wait for reconnect to complete
	wg.Wait()
	cancel()

	// Make a request to verify connection is working

	reqData := RequestData{
		ReqType: PROTO_OA_VERSION_REQ,
		Req: &ProtoOAVersionReq{
			PayloadType: PROTO_OA_VERSION_REQ.Enum(),
		},
		ResType: PROTO_OA_VERSION_RES,
		Res:     &ProtoOAVersionRes{},
	}

	if err := c.SendRequest(reqData); err != nil {
		t.Fatalf("error getting proxy version: %v", err)
		return
	}
}

func TestClientRateLimiter(t *testing.T) {
	c, err := createApiClient(Environment_Demo)
	if err != nil {
		t.Fatal(err)
	}

	// Create requests and send them concurrently to test if the rate limiter is able to handle
	// the requests without causing the server to rate limit
	numRequests := 500
	// Apply a big enough request timeout to ensure the requests do not fail due to timeout
	c.WithRequestTimeout(
		( /* Client side rate limiter constraint */ time.Second/49)*time.Duration(numRequests) +
			/* Round trip time */ time.Millisecond*10*time.Duration(numRequests),
	)

	if err = c.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}
	wg := sync.WaitGroup{}
	wg.Add(numRequests)

	var reqErrs []error = make([]error, numRequests)

	for i := range numRequests {
		go func(i int) {
			defer wg.Done()

			req := ProtoOAVersionReq{
				PayloadType: PROTO_OA_VERSION_REQ.Enum(),
			}
			var res ProtoOAVersionRes

			reqCtx := context.Background()

			reqData := RequestData{
				Ctx:     reqCtx,
				ReqType: PROTO_OA_VERSION_REQ,
				Req:     &req,
				ResType: PROTO_OA_VERSION_RES,
				Res:     &res,
			}

			if err := c.SendRequest(reqData); err != nil {
				reqErrs[i] = err
				return
			}
		}(i)
	}

	// Wait for all requests to finish
	wg.Wait()

	c.Disconnect()

	// Check for request errors
	for _, e := range reqErrs {
		if e != nil {
			t.Fatalf("request error: %v", e)
		}
	}
}

func TestBatchRequestSendOnConnLoss(t *testing.T) {
	c, err := createApiClient(Environment_Demo)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	// Create requests and send them concurrently to check if each request makes a proper round trip
	numRequests := 20
	wg := sync.WaitGroup{}
	wg.Add(numRequests)
	var reqErrs []error = make([]error, numRequests)

	connLossCh := make(chan ClientEvent)
	onConnLoss := func(e *ConnectionLossEvent) {
		for i := range numRequests {
			go func(i int) {
				defer wg.Done()

				req := ProtoOAVersionReq{
					PayloadType: PROTO_OA_VERSION_REQ.Enum(),
				}
				var res ProtoOAVersionRes

				reqCtx := context.Background()

				reqData := RequestData{
					Ctx:     reqCtx,
					ReqType: PROTO_OA_VERSION_REQ,
					Req:     &req,
					ResType: PROTO_OA_VERSION_RES,
					Res:     &res,
				}

				if err := c.SendRequest(reqData); err != nil {
					reqErrs[i] = err
					return
				}
			}(i)
		}
	}

	if err := SpawnClientEventHandler(ctx, connLossCh, onConnLoss); err != nil {
		t.Fatalf("error starting connection loss client event handler: %v", err)
	}

	if err := c.ListenToClientEvent(ctx, ClientEventListenData{ClientEventType_ConnectionLossEvent, connLossCh}); err != nil {
		t.Fatalf("error listening to connection loss client event: %v", err)
	}

	if err = c.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}

	// Simulate connection loss right before sending requests
	c.tcpClient.JustCloseConn()

	// Wait for all requests to finish
	wg.Wait()

	c.Disconnect()

	// Check for request errors
	for _, e := range reqErrs {
		if e != nil {
			t.Fatalf("request error: %v", e)
		}
	}
}

func TestPreElectReqCtxCancel(t *testing.T) {
	c, err := createApiClient(Environment_Demo)
	if err != nil {
		t.Fatal(err)
	}

	/*
		Scenario 1:
		Simulate connection loss state where BlockUntilReconnect handles the request ctx.
		Cancel context before SendRequest call to force this behaviour.
	*/
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	// Stuff the request queue with 10 requests which is 200ms worth of
	// request emission time (200ms due to rate limiter)
	errCh := make(chan error)

	connLossCh := make(chan ClientEvent)
	onConnLoss := func(e *ConnectionLossEvent) {
		req := ProtoOAVersionReq{
			PayloadType: PROTO_OA_VERSION_REQ.Enum(),
		}
		var res ProtoOAVersionRes

		ctx, cancelCtx := context.WithCancel(context.Background())
		cancelCtx()

		reqData := RequestData{
			Ctx:     ctx,
			ReqType: PROTO_OA_VERSION_REQ,
			Req:     &req,
			ResType: PROTO_OA_VERSION_RES,
			Res:     &res,
		}

		if err := c.SendRequest(reqData); err != nil {
			errCh <- err
		}
		close(errCh)
	}

	if err := SpawnClientEventHandler(ctx, connLossCh, onConnLoss); err != nil {
		t.Fatalf("error starting connection loss client event handler: %v", err)
	}

	if err := c.ListenToClientEvent(ctx, ClientEventListenData{ClientEventType_ConnectionLossEvent, connLossCh}); err != nil {
		t.Fatalf("error listening to connection loss client event: %v", err)
	}

	if err = c.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}
	defer c.Disconnect()

	// Simulate connection loss right before sending requests
	c.tcpClient.JustCloseConn()

	// Read from error channel
	err, ok := <-errCh
	if !ok {
		t.Fatal("expected request context expired error, got no error")
	}
	var reqCtxExpErr *RequestContextExpiredError
	if !errors.As(err, &reqCtxExpErr) {
		t.Fatalf("expected request context expired error, got %v", err)
	}
}

func TestPreElectReqCtxCancelBackpressure(t *testing.T) {
	/*
		Apply backpressure to the request queue to force the request heap to cancel in maximum
		duration DefaultRequestHeapIterationTimeout (50ms) * 2
	*/

	c, err := createApiClient(Environment_Demo)
	if err != nil {
		t.Fatal(err)
	}

	if err = c.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}
	defer c.Disconnect()

	numRequests := 250 // 5 seconds worth of backpressure due to rate limiter
	wg := sync.WaitGroup{}
	wg.Add(numRequests)
	for range numRequests {
		go func() {
			defer wg.Done()
			req := ProtoOAVersionReq{
				PayloadType: PROTO_OA_VERSION_REQ.Enum(),
			}
			var res ProtoOAVersionRes

			reqCtx := context.Background()

			reqData := RequestData{
				Ctx:     reqCtx,
				ReqType: PROTO_OA_VERSION_REQ,
				Req:     &req,
				ResType: PROTO_OA_VERSION_RES,
				Res:     &res,
			}

			c.SendRequest(reqData)
		}()
	}

	time.Sleep(time.Second * 2) // Wait two seconds to ensure backpressure applies

	// Send the request with immediately cancelled context
	req := ProtoOAVersionReq{
		PayloadType: PROTO_OA_VERSION_REQ.Enum(),
	}
	var res ProtoOAVersionRes

	ctx, cancelCtx := context.WithCancel(context.Background())
	cancelCtx()

	reqData := RequestData{
		Ctx:     ctx,
		ReqType: PROTO_OA_VERSION_REQ,
		Req:     &req,
		ResType: PROTO_OA_VERSION_RES,
		Res:     &res,
	}

	err = c.SendRequest(reqData)
	if err == nil {
		t.Fatal("expected request context expired error, got no error")
	}

	var reqCtxExpErr *RequestContextExpiredError
	if !errors.As(err, &reqCtxExpErr) {
		t.Fatalf("expected request context expired error, got %v", err)
	}

	wg.Wait()
}

func TestDisabledDefaultRateLimiter(t *testing.T) {
	/*
		TestGetSymbols mock
	*/
	accountId, accessToken, env, err := loadTestAccountCredentials()
	if err != nil {
		t.Fatalf("error loading test account credentials: %v", err)
	}
	c, err := createApiClient(env)
	if err != nil {
		t.Fatal(err)
	}
	c.DisableDefaultRateLimiter()

	if err = c.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}
	defer c.Disconnect()

	if _, err = c.AuthenticateAccount(CtraderAccountId(accountId), AccessToken(accessToken)); err != nil {
		t.Fatal(err)
	}

	reqData := RequestData{
		ReqType: PROTO_OA_SYMBOLS_LIST_REQ,
		Req: &ProtoOASymbolsListReq{
			CtidTraderAccountId:    proto.Int64(accountId),
			IncludeArchivedSymbols: proto.Bool(false),
		},
		ResType: PROTO_OA_SYMBOLS_LIST_RES,
		Res:     &ProtoOASymbolsListRes{},
	}

	if err = c.SendRequest(reqData); err != nil {
		t.Fatalf("error sending request: %v", err)
	}

	/*
		TestReconnect mock
	*/
	var once1 sync.Once
	var once2 sync.Once
	wg := sync.WaitGroup{}
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())

	onConnLossCh := make(chan ClientEvent)
	onReconnectSuccessCh := make(chan ClientEvent)
	onReconnectFailCh := make(chan ClientEvent)

	onConnectionLoss := func(e *ConnectionLossEvent) {
		once1.Do(func() {
			t.Log("closed tcp conn. waiting for reconnect...")
		})
	}
	onReconnectSuccess := func(e *ReconnectSuccessEvent) {
		// On reconnect success
		t.Log("reconnect successful")

		once2.Do(func() {
			wg.Done()
		})
	}
	onReconnectFail := func(e *ReconnectFailEvent) {
		// On reconnect failure
		t.Logf("reconnect failed: %v", e.Err)
	}

	if err := c.ListenToClientEvent(ctx, ClientEventListenData{ClientEventType_ConnectionLossEvent, onConnLossCh}); err != nil {
		t.Fatalf("error listening to connection loss events: %v", err)
	}
	if err := c.ListenToClientEvent(ctx, ClientEventListenData{ClientEventType_ReconnectSuccessEvent, onReconnectSuccessCh}); err != nil {
		t.Fatalf("error listening to reconnect success events: %v", err)
	}
	if err := c.ListenToClientEvent(ctx, ClientEventListenData{ClientEventType_ReconnectFailEvent, onReconnectFailCh}); err != nil {
		t.Fatalf("error listening to reconnect failure events: %v", err)
	}

	if err := SpawnClientEventHandler(ctx, onConnLossCh, onConnectionLoss); err != nil {
		t.Fatalf("error spawning connection loss event handler: %v", err)
	}
	if err := SpawnClientEventHandler(ctx, onReconnectSuccessCh, onReconnectSuccess); err != nil {
		t.Fatalf("error spawning reconnect success event handler: %v", err)
	}
	if err := SpawnClientEventHandler(ctx, onReconnectFailCh, onReconnectFail); err != nil {
		t.Fatalf("error spawning reconnect failure event handler: %v", err)
	}

	c.tcpClient.JustCloseConn()

	// Wait for reconnect to complete
	wg.Wait()
	cancel()

	// Make a request to verify connection is working

	reqData = RequestData{
		ReqType: PROTO_OA_VERSION_REQ,
		Req: &ProtoOAVersionReq{
			PayloadType: PROTO_OA_VERSION_REQ.Enum(),
		},
		ResType: PROTO_OA_VERSION_RES,
		Res:     &ProtoOAVersionRes{},
	}

	if err := c.SendRequest(reqData); err != nil {
		t.Fatalf("error getting proxy version: %v", err)
		return
	}
}
