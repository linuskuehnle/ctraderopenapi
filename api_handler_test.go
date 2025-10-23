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

func TestHandlerConnectDisconnect(t *testing.T) {
	h, err := createApiHandler(Environment_Demo)
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

func TestHandlerConnectDisconnectReconnect(t *testing.T) {
	h, err := createApiHandler(Environment_Demo)
	if err != nil {
		t.Fatal(err)
	}

	if err = h.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}

	if err = h.Disconnect(); err != nil {
		t.Fatalf("error disconnecting: %v", err)
	}

	if err = h.Reconnect(); err != nil {
		t.Fatalf("error reconnecting: %v", err)
	}
}

func TestHandlerConnectHeartbeatDisconnect(t *testing.T) {
	h, err := createApiHandler(Environment_Demo)
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

func TestHandlerReqConcurrency(t *testing.T) {
	h, err := createApiHandler(Environment_Demo)
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

			req := messages.ProtoOAVersionReq{
				PayloadType: PROTO_OA_VERSION_REQ.Enum(),
			}
			var res messages.ProtoOAVersionRes

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

	h, err := createApiHandler(env)
	if err != nil {
		t.Fatal(err)
	}

	if err = h.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}

	// Make unauthenticated request that will return a proto error response
	req := messages.ProtoOAGetTrendbarsReq{
		CtidTraderAccountId: proto.Int64(accountId),
		Period:              ProtoOATrendbarPeriod_D1.Enum(),
		SymbolId:            proto.Int64(1),
	}
	var res messages.ProtoOAGetTrendbarsRes

	reqCtx := context.Background()

	reqData := RequestData{
		Ctx:     reqCtx,
		ReqType: PROTO_OA_GET_TRENDBARS_REQ,
		Req:     &req,
		ResType: PROTO_OA_GET_TRENDBARS_RES,
		Res:     &res,
	}

	if err = h.SendRequest(reqData); err != nil {
		if _, ok := err.(*ResponseError); !ok {
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
	h, err := createApiHandler(Environment_Demo)
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
	h, err := createApiHandler(Environment_Demo)
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

	h, err := createApiHandler(env)
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
	h, err := createApiHandler(env)
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
		Req: &messages.ProtoOASymbolsListReq{
			CtidTraderAccountId:    proto.Int64(accountId),
			IncludeArchivedSymbols: proto.Bool(false),
		},
		ResType: PROTO_OA_SYMBOLS_LIST_RES,
		Res:     &messages.ProtoOASymbolsListRes{},
	}

	if err = h.SendRequest(reqData); err != nil {
		t.Fatalf("error sending request: %v", err)
	}

	if err := h.Disconnect(); err != nil {
		t.Fatalf("error disconnecting: %v", err)
	}
}
