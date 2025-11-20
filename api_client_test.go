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
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"
)

func TestClientConnectDisconnect(t *testing.T) {
	h, err := createApiClient(Environment_Demo)
	if err != nil {
		t.Fatal(err)
	}

	if err = h.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}

	if err = h.Disconnect(); err != nil {
		t.Fatalf("error disconnecting: %v", err)
	}
}

func TestClientConnectHeartbeatDisconnect(t *testing.T) {
	h, err := createApiClient(Environment_Demo)
	if err != nil {
		t.Fatal(err)
	}

	if err = h.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}

	if err := h.emitHeartbeat(); err != nil {
		t.Fatalf("error emitting heartbeat: %v", err)
	}

	if err = h.Disconnect(); err != nil {
		t.Fatalf("error disconnecting: %v", err)
	}
}

func TestClientReqConcurrency(t *testing.T) {
	h, err := createApiClient(Environment_Demo)
	if err != nil {
		t.Fatal(err)
	}

	if err = h.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}

	// Create requests and send them concurrently to check if each request makes a proper round trip
	numRequests := 20
	wg := sync.WaitGroup{}
	wg.Add(numRequests)

	var reqErrs []error = make([]error, numRequests)

	for i := range numRequests {
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

			if err := h.SendRequest(reqData); err != nil {
				reqErrs[i] = err
				return
			}
		}()
	}

	// Wait for all requests to finish
	wg.Wait()

	if err = h.Disconnect(); err != nil {
		t.Fatalf("error disconnecting: %v", err)
	}

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

	h, err := createApiClient(env)
	if err != nil {
		t.Fatal(err)
	}

	if err = h.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}

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

	if err = h.SendRequest(reqData); err != nil {
		var respErr *ResponseError
		if !errors.As(err, &respErr) {
			t.Fatalf("unexpected error of type RequestError. got: %v", err)
		}
	} else {
		t.Fatal("expected RequestError error return value")
	}

	if err := h.Disconnect(); err != nil {
		t.Fatalf("error disconnecting: %v", err)
	}
}

func TestFatalError(t *testing.T) {
	h, err := createApiClient(Environment_Demo)
	if err != nil {
		t.Fatal(err)
	}

	errCh, err := h.MakeFatalErrChan()
	if err != nil {
		t.Fatalf("error creating fatal error channel: %v", err)
	}

	if err = h.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}

	// Manually trigger a fatal error
	h.onFatalError(errors.New("test fatal error"))

	select {
	case fe := <-errCh:
		if fe.Error() != "test fatal error" {
			t.Fatalf("expected fatal error 'test fatal error', got: %v", fe)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for fatal error")
	}
}

func TestHeartbeatEmission(t *testing.T) {
	h, err := createApiClient(Environment_Demo)
	if err != nil {
		t.Fatal(err)
	}

	if err = h.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}

	// Test if any panics occur after emitting the heartbeat event
	time.Sleep(time.Second * (Heartbeat_Timeout_Seconds + 2))

	if err = h.Disconnect(); err != nil {
		t.Fatalf("error disconnecting: %v", err)
	}
}

func TestAccountAuth(t *testing.T) {
	accountId, accessToken, env, err := loadTestAccountCredentials()
	if err != nil {
		t.Fatalf("error loading test account credentials: %v", err)
	}

	h, err := createApiClient(env)
	if err != nil {
		t.Fatal(err)
	}

	if err = h.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}

	res, err := autheticateAccount(h, accountId, accessToken)
	if err != nil {
		t.Fatal(err)
	}

	if err := h.Disconnect(); err != nil {
		t.Fatalf("error disconnecting: %v", err)
	}

	if res.GetCtidTraderAccountId() != accountId {
		t.Fatalf("expected account ID %d, got %d", accountId, res.GetCtidTraderAccountId())
	}
}

func TestGetSymbols(t *testing.T) {
	accountId, accessToken, env, err := loadTestAccountCredentials()
	if err != nil {
		t.Fatalf("error loading test account credentials: %v", err)
	}
	h, err := createApiClient(env)
	if err != nil {
		t.Fatal(err)
	}

	if err = h.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}

	if _, err = autheticateAccount(h, accountId, accessToken); err != nil {
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

	if err = h.SendRequest(reqData); err != nil {
		t.Fatalf("error sending request: %v", err)
	}

	if err := h.Disconnect(); err != nil {
		t.Fatalf("error disconnecting: %v", err)
	}
}

func TestReconnect(t *testing.T) {
	h, err := createApiClient(Environment_Demo)
	if err != nil {
		t.Fatal(err)
	}

	var once1 sync.Once
	var once2 sync.Once
	wg := sync.WaitGroup{}
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())

	onConnLossCh := make(chan ListenableClientEvent)
	onReconnectSuccessCh := make(chan ListenableClientEvent)
	onReconnectFailCh := make(chan ListenableClientEvent)

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

	if err := h.ListenToClientEvent(ClientEventType_ConnectionLossEvent, onConnLossCh, ctx); err != nil {
		t.Fatalf("error listening to connection loss events: %v", err)
	}
	if err := h.ListenToClientEvent(ClientEventType_ReconnectSuccessEvent, onReconnectSuccessCh, ctx); err != nil {
		t.Fatalf("error listening to reconnect success events: %v", err)
	}
	if err := h.ListenToClientEvent(ClientEventType_ReconnectFailEvent, onReconnectFailCh, ctx); err != nil {
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

	if err = h.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}

	h.tcpClient.JustCloseConn()

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

	if err := h.SendRequest(reqData); err != nil {
		t.Fatalf("error getting proxy version: %v", err)
		return
	}

	if err = h.Disconnect(); err != nil {
		t.Fatalf("error disconnecting: %v", err)
	}
}
