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

package datatypes

import (
	"context"
	"sync"
)

type LifeCycleData interface {
	Start() error
	Stop() error
	IsRunning() bool

	GetContext() (context.Context, error)
	CancelContext() error
	SetClientConnected(bool) error
	IsClientConnected() bool
	SignalMessageSend()
	GetOnMessageSendCh() (chan struct{}, error)
}

type lifeCycleData struct {
	mu                sync.RWMutex
	isRunning         bool
	onMessageSendCh   chan struct{}
	ctx               context.Context
	cancelCtx         context.CancelFunc
	isClientConnected bool
}

func NewLifeCycleData() LifeCycleData {
	return newLifeCycleData()
}

func newLifeCycleData() *lifeCycleData {
	return &lifeCycleData{
		mu: sync.RWMutex{},
	}
}

func (d *lifeCycleData) Start() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.isRunning {
		return &LifeCycleAlreadyRunningError{
			CallContext: "error starting",
		}
	}

	d.init()

	d.isRunning = true

	return nil
}

func (d *lifeCycleData) Stop() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.isRunning {
		return &LifeCycleNotRunningError{
			CallContext: "error stopping",
		}
	}

	close(d.onMessageSendCh)
	d.cancelCtx()

	d.clear()

	d.isRunning = false

	return nil
}

func (d *lifeCycleData) IsRunning() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	isRunning := d.isRunning
	return isRunning
}

func (d *lifeCycleData) init() {
	d.onMessageSendCh = make(chan struct{}, 1)
	d.ctx, d.cancelCtx = context.WithCancel(context.Background())
	d.isClientConnected = false
}

func (d *lifeCycleData) clear() {
	d.onMessageSendCh = nil
	d.ctx = nil
	d.cancelCtx = nil
	d.isClientConnected = false
}

func (d *lifeCycleData) checkLifeCycleRunning(callContext string) error {
	if d.isRunning {
		return nil
	}

	return &LifeCycleNotRunningError{
		CallContext: callContext,
	}
}

func (d *lifeCycleData) GetContext() (context.Context, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if err := d.checkLifeCycleRunning("error on GetContext"); err != nil {
		return nil, err
	}

	ctx := d.ctx
	return ctx, nil
}

func (d *lifeCycleData) CancelContext() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.checkLifeCycleRunning("error on CancelContext"); err != nil {
		return err
	}

	d.cancelCtx()
	return nil
}

func (d *lifeCycleData) SetClientConnected(isConnected bool) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.checkLifeCycleRunning("error on SetClientConnected"); err != nil {
		return err
	}

	d.isClientConnected = isConnected
	return nil
}

func (d *lifeCycleData) IsClientConnected() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if err := d.checkLifeCycleRunning("error on IsClientConnected"); err != nil {
		return false
	}

	isConnected := d.isClientConnected
	return isConnected
}

func (d *lifeCycleData) SignalMessageSend() {
	d.mu.RLock()
	defer d.mu.RUnlock()

	d.onMessageSendCh <- struct{}{}
}

func (d *lifeCycleData) GetOnMessageSendCh() (chan struct{}, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if err := d.checkLifeCycleRunning("error on IsClientConnected"); err != nil {
		return nil, err
	}

	ch := d.onMessageSendCh
	return ch, nil
}
