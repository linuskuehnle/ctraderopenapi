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

	"sync"

	"google.golang.org/protobuf/proto"
)

type Event = proto.Message

type EventHandler interface {
	WithIgnoreIdsNotIncluded() EventHandler

	Clear()

	HasEvent(eventId messages.ProtoOAPayloadType) bool
	HandleEvent(eventId messages.ProtoOAPayloadType, event Event) error
	AddEvent(eventId messages.ProtoOAPayloadType, onEvent func(proto.Message)) error
	RemoveEvent(eventId messages.ProtoOAPayloadType) error
}

type eventHandler struct {
	mu                   sync.RWMutex
	eventsMap            map[messages.ProtoOAPayloadType]func(Event)
	ignoreIdsNotIncluded bool
}

func NewEventHandler() EventHandler {
	return newEventHandler()
}

func newEventHandler() *eventHandler {
	return &eventHandler{
		mu:                   sync.RWMutex{},
		eventsMap:            make(map[messages.ProtoOAPayloadType]func(Event)),
		ignoreIdsNotIncluded: false,
	}
}

func (h *eventHandler) WithIgnoreIdsNotIncluded() EventHandler {
	h.ignoreIdsNotIncluded = true
	return h
}

func (h *eventHandler) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.eventsMap = make(map[messages.ProtoOAPayloadType]func(Event))
}

func (h *eventHandler) HasEvent(eventId messages.ProtoOAPayloadType) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	_, exists := h.eventsMap[eventId]
	return exists
}

func (h *eventHandler) HandleEvent(eventId messages.ProtoOAPayloadType, event Event) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	onEvent, exists := h.eventsMap[eventId]
	if !exists {
		if h.ignoreIdsNotIncluded {
			return nil
		}

		return &IdNotIncludedError{
			Id: string(eventId),
		}
	}

	if onEvent != nil {
		onEvent(event)
	}

	return nil
}

// Event: Live Trendbars
func (h *eventHandler) AddEvent(eventId messages.ProtoOAPayloadType, onEvent func(Event)) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Check if already added
	if _, exists := h.eventsMap[eventId]; exists {
		return &IdAlreadyIncludedError{
			Id: string(eventId),
		}
	}

	// Register onEvent callback
	h.eventsMap[eventId] = onEvent

	return nil
}
func (h *eventHandler) RemoveEvent(eventId messages.ProtoOAPayloadType) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Check if not subscribed
	if _, exists := h.eventsMap[eventId]; !exists {
		return &IdNotIncludedError{
			Id: string(eventId),
		}
	}

	// Unregister onEvent callback
	delete(h.eventsMap, eventId)

	return nil
}
