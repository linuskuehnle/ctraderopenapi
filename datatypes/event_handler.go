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
	"sync"
)

type EventId int32

type EventHandler[T any] interface {
	WithIgnoreIdsNotIncluded() EventHandler[T]

	Clear()

	HasEvent(eventId EventId) bool
	HandleEvent(eventId EventId, event T) error
	AddEvent(eventId EventId, onEvent func(T)) error
	RemoveEvent(eventId EventId) error
}

type eventHandler[T any] struct {
	mu                   sync.RWMutex
	eventsMap            map[EventId]func(T)
	ignoreIdsNotIncluded bool
}

func NewEventHandler[T any]() EventHandler[T] {
	return newEventHandler[T]()
}

func newEventHandler[T any]() *eventHandler[T] {
	return &eventHandler[T]{
		mu:                   sync.RWMutex{},
		eventsMap:            make(map[EventId]func(T)),
		ignoreIdsNotIncluded: false,
	}
}

func (h *eventHandler[T]) WithIgnoreIdsNotIncluded() EventHandler[T] {
	h.ignoreIdsNotIncluded = true
	return h
}

func (h *eventHandler[T]) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.eventsMap = make(map[EventId]func(T))
}

func (h *eventHandler[T]) HasEvent(eventId EventId) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	_, exists := h.eventsMap[eventId]
	return exists
}

func (h *eventHandler[T]) HandleEvent(eventId EventId, event T) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	onEvent, exists := h.eventsMap[eventId]
	if !exists {
		if h.ignoreIdsNotIncluded {
			return nil
		}

		return &IdNotIncludedError{
			Id: eventId,
		}
	}

	if onEvent != nil {
		onEvent(event)
	}

	return nil
}

// Event: Live Trendbars
func (h *eventHandler[T]) AddEvent(eventId EventId, onEvent func(T)) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Check if already added
	if _, exists := h.eventsMap[eventId]; exists {
		return &IdAlreadyIncludedError{
			Id: eventId,
		}
	}

	// Register onEvent callback
	h.eventsMap[eventId] = onEvent

	return nil
}
func (h *eventHandler[T]) RemoveEvent(eventId EventId) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Check if not subscribed
	if _, exists := h.eventsMap[eventId]; !exists {
		return &IdNotIncludedError{
			Id: eventId,
		}
	}

	// Unregister onEvent callback
	delete(h.eventsMap, eventId)

	return nil
}
