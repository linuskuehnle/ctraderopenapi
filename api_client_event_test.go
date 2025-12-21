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
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// Subscribe + unsubscribe
func TestSubscribeUnsubscribeAPIEvent(t *testing.T) {
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

	var symbolId int64 = 1 // "EURUSD, Spot CFD"
	subData := SubscribableAPIEventData{
		EventType: APIEventType_Spots,
		SubcriptionData: &SubscriptionDataSpotEvent{
			CtraderAccountId: CtraderAccountId(accountId),
			SymbolIds:        []int64{symbolId},
		},
	}
	if err = c.SubscribeAPIEvent(subData); err != nil {
		t.Fatalf("error subscribing event: %v", err)
	}

	if err = c.UnsubscribeAPIEvent(subData); err != nil {
		t.Fatalf("error unsubscribing event: %v", err)
	}
}

// Subscribe, unsubscribe + listen to event and wait for immediate event message
func TestListenEvent(t *testing.T) {
	// Add listener and cancel context
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

	wg := sync.WaitGroup{}
	var once sync.Once

	wg.Add(1)
	eventCh := make(chan ListenableEvent)

	go func() {
		<-eventCh
		once.Do(func() {
			wg.Done()
		})
	}()

	ctx, cancelCtx := context.WithCancel(context.Background())

	if err = c.ListenToAPIEvent(ctx, APIEventType_Spots, eventCh); err != nil {
		t.Fatalf("error registering event listener: %v", err)
	}

	var symbolId int64 = 1 // "EURUSD, Spot CFD"

	subData := SubscribableAPIEventData{
		EventType: APIEventType_Spots,
		SubcriptionData: &SubscriptionDataSpotEvent{
			CtraderAccountId: CtraderAccountId(accountId),
			SymbolIds:        []int64{symbolId},
		},
	}
	if err = c.SubscribeAPIEvent(subData); err != nil {
		t.Fatalf("error subscribing event: %v", err)
	}

	// Wait for immediate event message
	wg.Wait()

	cancelCtx()

	if err = c.UnsubscribeAPIEvent(subData); err != nil {
		t.Fatalf("error unsubscribing event: %v", err)
	}
}

// Subscribe + Reconnect + Listen + cleanup (unsubscribe, deconnect)
func TestResubscribeOnReconnect(t *testing.T) {
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

	ctx, cancelCtx := context.WithCancel(context.Background())

	var symbolId int64 = 1 // "EURUSD, Spot CFD"
	subData := SubscribableAPIEventData{
		EventType: APIEventType_Spots,
		SubcriptionData: &SubscriptionDataSpotEvent{
			CtraderAccountId: CtraderAccountId(accountId),
			SymbolIds:        []int64{symbolId},
		},
	}
	if err = c.SubscribeAPIEvent(subData); err != nil {
		t.Fatalf("error subscribing event: %v", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	var onceSpots sync.Once
	spotsEventCh := make(chan ListenableEvent)
	go func() {
		<-spotsEventCh
		onceSpots.Do(func() {
			wg.Done()
		})
	}()

	errCh := make(chan error)

	var onceConnLoss sync.Once
	connLossCh := make(chan ListenableClientEvent)
	go func() {
		<-connLossCh
		onceConnLoss.Do(func() {
			if err = c.ListenToAPIEvent(ctx, APIEventType_Spots, spotsEventCh); err != nil {
				errCh <- fmt.Errorf("error registering event listener: %v", err)
			}
			close(errCh)
		})
	}()

	if err := c.ListenToClientEvent(ctx, ClientEventType_ConnectionLossEvent, connLossCh); err != nil {
		t.Fatalf("error listening to connection loss event: %v", err)
	}

	// Simulate connection loss
	c.tcpClient.JustCloseConn()

	// Ensure the initial event message from before JustCloseConn is flushed since no listener is installed yet.
	time.Sleep(time.Second * 2)

	err, ok := <-errCh
	if ok {
		t.Fatal(err)
	}

	// Wait for event message
	wg.Wait()
	cancelCtx()

	// Cleanup
	if err = c.UnsubscribeAPIEvent(subData); err != nil {
		t.Fatalf("error unsubscribing event: %v", err)
	}
}

// Subscribe, unsubscribe + listen to event and wait for immediate event message using SpawnAPIEventHandler
func TestListenEventWithSpawnAPIEventHandler(t *testing.T) {
	// Add listener and cancel context
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

	wg := sync.WaitGroup{}
	var once sync.Once

	wg.Add(1)
	eventCh := make(chan ListenableEvent)

	eventCallback := func(event *ProtoOASpotEvent) {
		once.Do(func() {
			wg.Done()
		})
	}

	if err := SpawnAPIEventHandler(context.Background(), eventCh, eventCallback); err != nil {
		t.Fatalf("error registering event handler: %v", err)
	}

	ctx, cancelCtx := context.WithCancel(context.Background())

	if err = c.ListenToAPIEvent(ctx, APIEventType_Spots, eventCh); err != nil {
		t.Fatalf("error registering event listener: %v", err)
	}

	var symbolId int64 = 1 // "EURUSD, Spot CFD"

	subData := SubscribableAPIEventData{
		EventType: APIEventType_Spots,
		SubcriptionData: &SubscriptionDataSpotEvent{
			CtraderAccountId: CtraderAccountId(accountId),
			SymbolIds:        []int64{symbolId},
		},
	}
	if err = c.SubscribeAPIEvent(subData); err != nil {
		t.Fatalf("error subscribing event: %v", err)
	}

	// Wait for immediate event message
	wg.Wait()

	cancelCtx()

	if err = c.UnsubscribeAPIEvent(subData); err != nil {
		t.Fatalf("error unsubscribing event: %v", err)
	}
}

// Test live trendbars event association with spots event
func TestLiveTrendbarEvent(t *testing.T) {
	// Add listener and cancel context
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

	wg := sync.WaitGroup{}
	var once sync.Once

	wg.Add(1)
	eventCh := make(chan ListenableEvent)

	go func() {
		for range eventCh {
			once.Do(func() {
				wg.Done()
			})
		}
	}()

	ctx, cancelCtx := context.WithCancel(context.Background())

	if err = c.ListenToAPIEvent(ctx, APIEventType_Spots, eventCh); err != nil {
		t.Fatalf("error registering spots event listener: %v", err)
	}

	var symbolId int64 = 1 // "EURUSD, Spot CFD"

	subDataLiveTrendbars := SubscribableAPIEventData{
		EventType: APIEventType_LiveTrendbars,
		SubcriptionData: &SubscriptionDataLiveTrendbarEvent{
			CtraderAccountId: CtraderAccountId(accountId),
			SymbolId:         symbolId,
			Period:           ProtoOATrendbarPeriod_D1,
		},
	}
	// Subscribe live trendbars before spots (expected to return an error)
	err = c.SubscribeAPIEvent(subDataLiveTrendbars)
	if err == nil {
		t.Fatal("error subscribing live trendbars event before spots: expected an error, returned nil")
	}
	if err.Error() != "INVALID_REQUEST. Impossible to get trendbars before the spot subscribing" {
		t.Fatalf("error subscribing live trendbars event before spots: expected "+
			"\"INVALID_REQUEST. Impossible to get trendbars before the spot subscribing\", got %v", err)
	}

	subDataSpots := SubscribableAPIEventData{
		EventType: APIEventType_Spots,
		SubcriptionData: &SubscriptionDataSpotEvent{
			CtraderAccountId: CtraderAccountId(accountId),
			SymbolIds:        []int64{symbolId},
		},
	}
	if err = c.SubscribeAPIEvent(subDataSpots); err != nil {
		t.Fatalf("error subscribing spots event: %v", err)
	}

	// Subscribe live trendbars after spots (expected to not return an error)
	if err = c.SubscribeAPIEvent(subDataLiveTrendbars); err != nil {
		t.Fatalf("error subscribing live trendbars event: %v", err)
	}

	// Wait for immediate event message
	wg.Wait()

	cancelCtx()

	if err = c.UnsubscribeAPIEvent(subDataSpots); err != nil {
		t.Fatalf("error unsubscribing spots event: %v", err)
	}
}

// Test client events
func TestClientEvents(t *testing.T) {
	c, err := createApiClient(Environment_Demo)
	if err != nil {
		t.Fatal(err)
	}

	connLossCh := make(chan ListenableClientEvent)
	if err := c.ListenToClientEvent(context.Background(), ClientEventType_ConnectionLossEvent, connLossCh); err != nil {
		t.Fatalf("error registering connection loss event listener: %v", err)
	}
	reconnectSuccessCh := make(chan ListenableClientEvent)
	if err := c.ListenToClientEvent(context.Background(), ClientEventType_ReconnectSuccessEvent, reconnectSuccessCh); err != nil {
		t.Fatalf("error registering connection loss event listener: %v", err)
	}

	if err = c.Connect(); err != nil {
		t.Fatalf("error connecting: %v", err)
	}
	defer c.Disconnect()

	c.tcpClient.JustCloseConn()

	// Wait for connection loss event
	<-connLossCh

	// Wait for reconnect success event
	<-reconnectSuccessCh

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
