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
	"github.com/linuskuehnle/ctraderopenapi/messages"

	"context"
	"fmt"
	"sync"
)

type ResponseCallbackInfo struct {
	mCtx        context.Context
	mCtxCancel  context.CancelFunc
	Ctx         context.Context
	ResCallback func(*messages.ProtoMessage, messages.ProtoOAPayloadType)
}

type RequestMapper interface {
	Clear()
	Has(id RequestId) bool
	Add(id RequestId, callbackInfo *ResponseCallbackInfo) error
	Remove(id RequestId) error
	Call(id RequestId, protoMsg *messages.ProtoMessage, payloadType messages.ProtoOAPayloadType) error
}

// Map request ids to response callback functions
type requestMapper struct {
	mu              sync.RWMutex
	callbackInfoMap map[RequestId]*ResponseCallbackInfo
}

func NewRequestMapper() RequestMapper {
	return newRequestMapper()
}

func newRequestMapper() *requestMapper {
	return &requestMapper{
		mu:              sync.RWMutex{},
		callbackInfoMap: make(map[RequestId]*ResponseCallbackInfo),
	}
}

func (m *requestMapper) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for reqId := range m.callbackInfoMap {
		delete(m.callbackInfoMap, reqId)
	}
}

func (m *requestMapper) Has(id RequestId) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.has(id)
}

func (m *requestMapper) has(id RequestId) bool {
	_, exists := m.callbackInfoMap[id]
	return exists
}

func (m *requestMapper) Add(id RequestId, callbackInfo *ResponseCallbackInfo) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if callbackInfo == nil {
		return &FunctionInvalidArgError{
			FunctionName: "Add",
			Err:          fmt.Errorf("callbackInfo may not be nil"),
		}
	}
	if callbackInfo.Ctx == nil {
		return &FunctionInvalidArgError{
			FunctionName: "Add",
			Err:          fmt.Errorf("field Ctx may not be nil"),
		}
	}
	if callbackInfo.ResCallback == nil {
		return &FunctionInvalidArgError{
			FunctionName: "Add",
			Err:          fmt.Errorf("field ResCallback may not be nil"),
		}
	}

	if m.has(id) {
		return fmt.Errorf("response callback entry already exists for request id %v", id)
	}

	m.add(id, callbackInfo)
	return nil
}

func (m *requestMapper) add(id RequestId, callbackInfo *ResponseCallbackInfo) error {
	m.callbackInfoMap[id] = callbackInfo

	callbackInfo.mCtx, callbackInfo.mCtxCancel = context.WithCancel(context.Background())

	go func(i *ResponseCallbackInfo) {
		select {
		case <-i.Ctx.Done():
			// Request context has been cancelled early by callee
			m.Remove(id)
			return
		case <-i.mCtx.Done():
			// Internal context has been cancelled so res callback is called
			return
		}
	}(callbackInfo)

	return nil
}

func (m *requestMapper) Remove(id RequestId) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.has(id) {
		return fmt.Errorf("response callback entry does not exist for request id %v", id)
	}

	m.remove(id)
	return nil
}

func (m *requestMapper) remove(id RequestId) {
	delete(m.callbackInfoMap, id)
}

func (m *requestMapper) Call(id RequestId, protoMsg *messages.ProtoMessage, payloadType messages.ProtoOAPayloadType) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.has(id) {
		return fmt.Errorf("response callback entry does not exist for request id %v", id)
	}

	callbackInfo := m.callbackInfoMap[id]
	resCallback := callbackInfo.ResCallback

	// Cancel internal context to signal response callback is being called
	callbackInfo.mCtxCancel()
	// Remove callback info entry
	m.remove(id)

	go resCallback(protoMsg, payloadType)

	return nil
}
