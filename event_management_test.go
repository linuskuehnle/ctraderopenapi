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
	"sync"
	"testing"
)

// Subscribe + unsubscribe
func TestSubscribeUnsubscribeEvent(t *testing.T) {
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

	var symbolId int64 = 315 // "German 40 Index, Spot CFD"
	subData := SubscribableEventData{
		EventType: EventType_Spots,
		SubcriptionData: &SubscriptionDataSpotEvent{
			CtraderAccountId: CtraderAccountId(accountId),
			SymbolIds:        []int64{symbolId},
		},
	}
	if err = h.SubscribeEvent(subData); err != nil {
		t.Fatalf("error subscribing event: %v", err)
	}

	if err = h.UnsubscribeEvent(subData); err != nil {
		t.Fatalf("error unsubscribing event: %v", err)
	}

	if err := h.Disconnect(); err != nil {
		t.Fatalf("error disconnecting: %v", err)
	}
}

// Subscribe, unsubscribe + listen to event and wait for immediate event message
func TestListenEvent(t *testing.T) {
	// Add listener and cancel context
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

	// Listen to event. Use the adapter so callers can provide a typed
	// callback without having to wrap it manually.
	if err = h.ListenToEvent(EventType_Spots, eventCh, ctx); err != nil {
		t.Fatalf("error registering event listener: %v", err)
	}

	var symbolId int64 = 315 // "German 40 Index, Spot CFD"

	subData := SubscribableEventData{
		EventType: EventType_Spots,
		SubcriptionData: &SubscriptionDataSpotEvent{
			CtraderAccountId: CtraderAccountId(accountId),
			SymbolIds:        []int64{symbolId},
		},
	}
	if err = h.SubscribeEvent(subData); err != nil {
		t.Fatalf("error subscribing event: %v", err)
	}

	// Wait for immediate event message
	wg.Wait()

	cancelCtx()

	if err = h.UnsubscribeEvent(subData); err != nil {
		t.Fatalf("error unsubscribing event: %v", err)
	}

	if err := h.Disconnect(); err != nil {
		t.Fatalf("error disconnecting: %v", err)
	}
}

// Subscribe, unsubscribe + listen to event and wait for immediate event message using SpawnEventHandler
func TestListenEventWithSpawnEventHandler(t *testing.T) {
	// Add listener and cancel context
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

	wg := sync.WaitGroup{}
	var once sync.Once

	wg.Add(1)
	eventCh := make(chan ListenableEvent)

	eventCallback := func(event *ProtoOASpotEvent) {
		once.Do(func() {
			wg.Done()
		})
	}

	if err := SpawnEventHandler(context.Background(), eventCh, eventCallback); err != nil {
		t.Fatalf("error registering event handler: %v", err)
	}

	ctx, cancelCtx := context.WithCancel(context.Background())

	// Listen to event. Use the adapter so callers can provide a typed
	// callback without having to wrap it manually.
	if err = h.ListenToEvent(EventType_Spots, eventCh, ctx); err != nil {
		t.Fatalf("error registering event listener: %v", err)
	}

	var symbolId int64 = 315 // "German 40 Index, Spot CFD"

	subData := SubscribableEventData{
		EventType: EventType_Spots,
		SubcriptionData: &SubscriptionDataSpotEvent{
			CtraderAccountId: CtraderAccountId(accountId),
			SymbolIds:        []int64{symbolId},
		},
	}
	if err = h.SubscribeEvent(subData); err != nil {
		t.Fatalf("error subscribing event: %v", err)
	}

	// Wait for immediate event message
	wg.Wait()

	cancelCtx()

	if err = h.UnsubscribeEvent(subData); err != nil {
		t.Fatalf("error unsubscribing event: %v", err)
	}

	if err := h.Disconnect(); err != nil {
		t.Fatalf("error disconnecting: %v", err)
	}
}

// Test live trendbars event association with spots event
func TestLiveTrendbarEvent(t *testing.T) {
	// Add listener and cancel context
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

	// Listen to event. Use the adapter so callers can provide a typed
	// callback without having to wrap it manually.
	if err = h.ListenToEvent(EventType_Spots, eventCh, ctx); err != nil {
		t.Fatalf("error registering spots event listener: %v", err)
	}

	var symbolId int64 = 315 // "German 40 Index, Spot CFD"

	subDataLiveTrendbars := SubscribableEventData{
		EventType: EventType_LiveTrendbars,
		SubcriptionData: &SubscriptionDataLiveTrendbarEvent{
			CtraderAccountId: CtraderAccountId(accountId),
			SymbolId:         symbolId,
			Period:           ProtoOATrendbarPeriod_D1,
		},
	}
	// Subscribe live trendbars before spots (expected to return an error)
	err = h.SubscribeEvent(subDataLiveTrendbars)
	if err == nil {
		t.Fatal("error subscribing live trendbars event before spots: expected an error, returned nil")
	}
	if err.Error() != "INVALID_REQUEST. Impossible to get trendbars before the spot subscribing" {
		t.Fatalf("error subscribing live trendbars event before spots: expected "+
			"\"INVALID_REQUEST. Impossible to get trendbars before the spot subscribing\", got %v", err)
	}

	subDataSpots := SubscribableEventData{
		EventType: EventType_Spots,
		SubcriptionData: &SubscriptionDataSpotEvent{
			CtraderAccountId: CtraderAccountId(accountId),
			SymbolIds:        []int64{symbolId},
		},
	}
	if err = h.SubscribeEvent(subDataSpots); err != nil {
		t.Fatalf("error subscribing spots event: %v", err)
	}

	// Subscribe live trendbars after spots (expected to not return an error)
	if err = h.SubscribeEvent(subDataLiveTrendbars); err != nil {
		t.Fatalf("error subscribing live trendbars event: %v", err)
	}

	// Wait for immediate event message
	wg.Wait()

	cancelCtx()

	if err = h.UnsubscribeEvent(subDataSpots); err != nil {
		t.Fatalf("error unsubscribing spots event: %v", err)
	}

	if err := h.Disconnect(); err != nil {
		t.Fatalf("error disconnecting: %v", err)
	}
}
