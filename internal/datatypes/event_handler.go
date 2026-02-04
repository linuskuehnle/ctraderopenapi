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

// EventHandler is the type used for distributing API events to
// specific channels based on subscription keys.
type EventHandler[T BaseEvent] interface {
	HasListenerSource(eventType EventType) bool

	SetDefaultListener(eventType EventType, eventCh chan T) error
	UnsetDefaultListener(eventType EventType) error
	SetListener(eventType EventType, keyData APIEventKeyData, eventCh chan T) error
	UnsetListener(eventType EventType, keyData APIEventKeyData) error

	HandleEvent(eventType EventType, event T) error
}

type IEventHandler[T BaseEvent] interface {
	EventHandler[T]
	Clear()
}

type eventHandler[T BaseEvent] struct {
	mu sync.RWMutex

	listenerSource   map[EventType]chan T
	defaultListeners map[EventType]chan T
	listeners        map[APIEventKey]chan T
}

func NewEventHandler[T BaseEvent]() EventHandler[T] {
	return newEventHandler[T]()
}

func NewIEventHandler[T BaseEvent]() IEventHandler[T] {
	return newEventHandler[T]()
}

func newEventHandler[T BaseEvent]() *eventHandler[T] {
	d := &eventHandler[T]{
		mu: sync.RWMutex{},

		listenerSource:   make(map[EventType]chan T),
		defaultListeners: make(map[EventType]chan T),
		listeners:        make(map[APIEventKey]chan T),
	}

	return d
}

func (d *eventHandler[T]) Clear() {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Close all listener source channels
	for _, ch := range d.listenerSource {
		close(ch)
	}

	d.listenerSource = make(map[EventType]chan T)
	d.defaultListeners = make(map[EventType]chan T)
	d.listeners = make(map[APIEventKey]chan T)
}

func (d *eventHandler[T]) setListenerSource(eventType EventType) {
	sourceCh := make(chan T)
	d.listenerSource[eventType] = sourceCh

	go d.executeEventListen(eventType, sourceCh)
}

func (d *eventHandler[T]) unsetListenerSource(eventType EventType) {
	sourceCh := d.listenerSource[eventType]

	close(sourceCh)
	delete(d.listenerSource, eventType)
}

func (d *eventHandler[T]) HasListenerSource(eventType EventType) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.hasListenerSource(eventType)
}

func (d *eventHandler[T]) hasListenerSource(eventType EventType) bool {
	_, exists := d.listenerSource[eventType]
	return exists
}

func (d *eventHandler[T]) SetDefaultListener(eventType EventType, eventCh chan T) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.defaultListeners[eventType]; exists {
		return &ListenerAlreadySetError{callContext: "default listener", EventType: eventType}
	}

	if !d.hasListenerSource(eventType) {
		d.setListenerSource(eventType)
	}

	d.defaultListeners[eventType] = eventCh
	return nil
}

func (d *eventHandler[T]) UnsetDefaultListener(eventType EventType) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	listenerCh, exists := d.defaultListeners[eventType]
	if !exists {
		return &ListenerNotSetError{callContext: "default listener", EventType: eventType}
	}

	close(listenerCh)
	delete(d.defaultListeners, eventType)
	return nil
}

func (d *eventHandler[T]) SetListener(eventType EventType, keyData APIEventKeyData, eventCh chan T) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	key := keyData.BuildKey()

	if _, exists := d.listeners[key]; exists {
		return &ListenerAlreadySetError{callContext: "listener", EventType: eventType}
	}

	if !d.hasListenerSource(eventType) {
		d.setListenerSource(eventType)
	}

	d.listeners[key] = eventCh
	return nil
}

func (d *eventHandler[T]) UnsetListener(eventType EventType, keyData APIEventKeyData) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	key := keyData.BuildKey()

	listenerCh, exists := d.listeners[key]
	if !exists {
		return &ListenerNotSetError{callContext: "listener", EventType: eventType}
	}

	close(listenerCh)
	delete(d.listeners, key)
	return nil
}

func (d *eventHandler[T]) HandleEvent(eventType EventType, event T) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	sourceCh, exists := d.listenerSource[eventType]
	if !exists {
		return &SourceChannelNotExistingError{EventType: eventType}
	}

	sourceCh <- event
	return nil
}

func (d *eventHandler[T]) executeEventListen(eventType EventType, sourceCh chan T) {
	for event := range sourceCh {
		distrEvent := event.ToDistributableEvent()
		isDistr := distrEvent != nil

		d.mu.RLock()

		var sentSpecific bool
		if isDistr {
			// Event is distributable
			keys := distrEvent.BuildKeys()

			for _, key := range keys {
				listenerCh, hasListener := d.listeners[key]
				if hasListener {
					sentSpecific = true
					listenerCh <- event
					break
				}
			}
		}

		// Send to default listener if event was not sent to specific listener
		if !sentSpecific {
			defaultCh, hasDefaultListener := d.defaultListeners[eventType]
			if hasDefaultListener {
				defaultCh <- event
			}
		}

		d.mu.RUnlock()
	}
}
