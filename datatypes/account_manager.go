package datatypes

import (
	"context"
	"maps"
	"slices"
	"sync"
)

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

type AccountManager[EventT comparable, SubDataT any] interface {
	LockModification(context.Context)

	HasAccessToken(AccessToken) bool
	AddAccessToken(AccessToken) error
	RemoveAccessToken(AccessToken) error
	UpdateAccessToken(old, new AccessToken) error

	AddAccountId(AccessToken, CtraderAccountId) error
	RemoveAccountId(CtraderAccountId) error

	AddEventSubscription(CtraderAccountId, EventT, SubDataT) error
	RemoveEventSubscription(CtraderAccountId, EventT) error

	GetAccessTokenByAccountId() map[CtraderAccountId]AccessToken
	GetAllAccountIds() []CtraderAccountId
	GetAccountIdsOfAccessToken(AccessToken) ([]CtraderAccountId, error)
	GetEventSubscriptionsOfAccountId(CtraderAccountId) (map[EventT]SubDataT, error)
}

type accountManager[EventT comparable, SubDataT any] struct {
	mu sync.RWMutex

	// Context based modification locking mechanism
	modMu      sync.Mutex
	modLockCtx context.Context

	// Accessory
	accessTokenByAccountId  map[CtraderAccountId]AccessToken
	accountIdsByAccessToken map[AccessToken][]CtraderAccountId

	// Event subscription referencing
	subscriptionsByAccount map[CtraderAccountId]map[EventT]SubDataT
}

func NewAccountManager[EventT comparable, SubDataT any]() AccountManager[EventT, SubDataT] {
	return newAccountManager[EventT, SubDataT]()
}

func newAccountManager[EventT comparable, SubDataT any]() *accountManager[EventT, SubDataT] {
	return &accountManager[EventT, SubDataT]{
		mu: sync.RWMutex{},

		modMu:      sync.Mutex{},
		modLockCtx: nil,

		accessTokenByAccountId:  make(map[CtraderAccountId]AccessToken),
		accountIdsByAccessToken: make(map[AccessToken][]CtraderAccountId),

		subscriptionsByAccount: make(map[CtraderAccountId]map[EventT]SubDataT),
	}
}

func (m *accountManager[EventT, SubDataT]) LockModification(ctx context.Context) {
	if ctx == nil {
		return
	}

	m.modMu.Lock()
	defer m.modMu.Unlock()

	m.modLockCtx = ctx
}

func (m *accountManager[EventT, SubDataT]) waitForModificationUnlock() {
	m.modMu.Lock()
	defer m.modMu.Unlock()

	if m.modLockCtx != nil {
		<-m.modLockCtx.Done()
		m.modLockCtx = nil
	}

	m.mu.Lock()
}

func (m *accountManager[EventT, SubDataT]) HasAccessToken(token AccessToken) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, hasAccessToken := m.accountIdsByAccessToken[token]
	return hasAccessToken
}

func (m *accountManager[EventT, SubDataT]) AddAccessToken(token AccessToken) error {
	// Wait for any active modification lock before acquiring the write lock
	m.waitForModificationUnlock()
	defer m.mu.Unlock()

	// Check if token already exists
	if _, exists := m.accountIdsByAccessToken[token]; exists {
		return &AccessTokenAlreadyExistsError{AccessToken: token}
	}

	m.accountIdsByAccessToken[token] = make([]CtraderAccountId, 0)
	return nil
}

func (m *accountManager[EventT, SubDataT]) RemoveAccessToken(token AccessToken) error {
	// Wait for any active modification lock before acquiring the write lock
	m.waitForModificationUnlock()
	defer m.mu.Unlock()

	accountIds, exists := m.accountIdsByAccessToken[token]
	if !exists {
		return &AccessTokenDoesNotExist{AccessToken: token}
	}

	// Remove all accounts associated with this token
	for _, accountId := range accountIds {
		delete(m.accessTokenByAccountId, accountId)
		delete(m.subscriptionsByAccount, accountId)
	}

	delete(m.accountIdsByAccessToken, token)
	return nil
}

func (m *accountManager[EventT, SubDataT]) UpdateAccessToken(old, new AccessToken) error {
	// Wait for any active modification lock before acquiring the write lock
	m.waitForModificationUnlock()
	defer m.mu.Unlock()

	// Check if old token exists
	accountIds, exists := m.accountIdsByAccessToken[old]
	if !exists {
		return &AccessTokenDoesNotExist{AccessToken: old}
	}

	// Check if new token already exists
	if _, exists := m.accountIdsByAccessToken[new]; exists {
		return &AccessTokenAlreadyExistsError{AccessToken: new}
	}

	// Move all accountIds from old token to new token
	m.accountIdsByAccessToken[new] = accountIds

	// Update the access token for each account
	for _, accountId := range accountIds {
		m.accessTokenByAccountId[accountId] = new
	}

	delete(m.accountIdsByAccessToken, old)
	return nil
}

func (m *accountManager[EventT, SubDataT]) AddAccountId(token AccessToken, accountId CtraderAccountId) error {
	// Wait for any active modification lock before acquiring the write lock
	m.waitForModificationUnlock()
	defer m.mu.Unlock()

	// Check if token exists
	accountIds, exists := m.accountIdsByAccessToken[token]
	if !exists {
		return &AccessTokenDoesNotExist{AccessToken: token}
	}

	// Check if account already exists
	if slices.Contains(accountIds, accountId) {
		return &AccountIdAlreadyExistsError{AccountId: accountId}
	}

	m.accountIdsByAccessToken[token] = append(accountIds, accountId)
	m.accessTokenByAccountId[accountId] = token
	m.subscriptionsByAccount[accountId] = make(map[EventT]SubDataT)

	return nil
}

func (m *accountManager[EventT, SubDataT]) RemoveAccountId(accountId CtraderAccountId) error {
	// Wait for any active modification lock before acquiring the write lock
	m.waitForModificationUnlock()
	defer m.mu.Unlock()

	token, exists := m.accessTokenByAccountId[accountId]
	if !exists {
		return &AccountIdDoesNotExist{AccountId: accountId}
	}

	// Check if token exists
	accountIds, exists := m.accountIdsByAccessToken[token]
	if !exists {
		return &AccessTokenDoesNotExist{AccessToken: token}
	}

	// Find and remove the account
	foundIndex := -1
	for i, id := range accountIds {
		if id == accountId {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		// Should never happen actually
		return &AccountIdDoesNotExistOnToken{AccountId: accountId, AccessToken: token}
	}

	// Remove from accountIds slice
	m.accountIdsByAccessToken[token] = append(accountIds[:foundIndex], accountIds[foundIndex+1:]...)

	// Remove from other maps
	delete(m.accessTokenByAccountId, accountId)
	delete(m.subscriptionsByAccount, accountId)

	return nil
}

func (m *accountManager[EventT, SubDataT]) AddEventSubscription(accountId CtraderAccountId, event EventT, subData SubDataT) error {
	// Wait for any active modification lock before acquiring the write lock
	m.waitForModificationUnlock()
	defer m.mu.Unlock()

	// Check if account exists
	subs, exists := m.subscriptionsByAccount[accountId]
	if !exists {
		return &AccountIdDoesNotExist{AccountId: accountId}
	}

	// Check if subscription already exists
	if _, exists := subs[event]; exists {
		return &EventSubscriptionAlreadyExistsError[EventT]{EventType: event, AccountId: accountId}
	}

	// Add or update the subscription
	subs[event] = subData
	return nil
}

func (m *accountManager[EventT, SubDataT]) RemoveEventSubscription(accountId CtraderAccountId, event EventT) error {
	// Wait for any active modification lock before acquiring the write lock
	m.waitForModificationUnlock()
	defer m.mu.Unlock()

	// Check if account exists
	subs, exists := m.subscriptionsByAccount[accountId]
	if !exists {
		return &AccountIdDoesNotExist{AccountId: accountId}
	}

	// Check if subscription exists
	if _, exists := subs[event]; !exists {
		return &EventSubscriptionNotExistingError[EventT]{EventType: event, AccountId: accountId}
	}

	delete(subs, event)
	return nil
}

func (m *accountManager[EventT, SubDataT]) GetAccessTokenByAccountId() map[CtraderAccountId]AccessToken {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy of the subscriptions map
	result := make(map[CtraderAccountId]AccessToken)
	maps.Copy(result, m.accessTokenByAccountId)
	return result
}

func (m *accountManager[EventT, SubDataT]) GetAllAccountIds() []CtraderAccountId {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ctids := []CtraderAccountId{}
	for ctid := range m.accessTokenByAccountId {
		ctids = append(ctids, ctid)
	}

	return ctids
}

func (m *accountManager[EventT, SubDataT]) GetAccountIdsOfAccessToken(token AccessToken) ([]CtraderAccountId, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Check if token exists
	accountIds, exists := m.accountIdsByAccessToken[token]
	if !exists {
		return nil, &AccessTokenDoesNotExist{AccessToken: token}
	}

	// Return a copy of the accountIds slice
	result := make([]CtraderAccountId, len(accountIds))
	copy(result, accountIds)
	return result, nil
}

func (m *accountManager[EventT, SubDataT]) GetEventSubscriptionsOfAccountId(accountId CtraderAccountId) (map[EventT]SubDataT, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Check if account exists
	subs, exists := m.subscriptionsByAccount[accountId]
	if !exists {
		return nil, &AccountIdDoesNotExist{AccountId: accountId}
	}

	// Return a copy of the subscriptions map
	result := make(map[EventT]SubDataT)
	maps.Copy(result, subs)
	return result, nil
}
