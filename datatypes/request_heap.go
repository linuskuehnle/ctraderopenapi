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
	"time"
)

const DefaultRequestHeapIterationTimeout = 50 * time.Millisecond
const MinRequestHeapIterationTimeout = 10 * time.Millisecond

type RequestHeap interface {
	Start() error
	Stop() error
	SetIterationInterval(time.Duration) error
	AddNode(reqMetaData *RequestMetaData) error
	RemoveNode(id RequestId) (*RequestMetaData, error)
}

type requestHeap struct {
	mu               sync.RWMutex
	iterationTimeout time.Duration
	heap             []*RequestMetaData
	expiredNodes     []*RequestMetaData
	ctx              context.Context
	cancelCtx        context.CancelFunc
}

func NewRequestHeap() RequestHeap {
	return newRequestHeap()
}

func newRequestHeap() *requestHeap {
	return &requestHeap{
		mu:               sync.RWMutex{},
		iterationTimeout: DefaultRequestHeapIterationTimeout,
		heap:             []*RequestMetaData{},
		expiredNodes:     []*RequestMetaData{},
	}
}

func (h *requestHeap) Start() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.ctx != nil {
		return &RequestHeapAlreadyRunningError{
			CallContext: "Start",
		}
	}

	h.ctx, h.cancelCtx = context.WithCancel(context.Background())

	go func(ctx context.Context, timeout time.Duration) {
		t := time.NewTicker(timeout)

		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				h.checkRequestContexts()
			}
		}
	}(h.ctx, h.iterationTimeout)

	return nil
}

func (h *requestHeap) Stop() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.ctx == nil {
		return &RequestHeapNotRunningError{
			CallContext: "Stop",
		}
	}

	// Stop goroutine by cancelling runtime context
	h.cancelCtx()

	h.clear()

	return nil
}

func (h *requestHeap) SetIterationInterval(interval time.Duration) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.ctx != nil {
		return &RequestHeapAlreadyRunningError{
			CallContext: "SetIterationInterval",
		}
	}

	h.iterationTimeout = interval
	return nil
}

func (h *requestHeap) clear() {
	h.heap = []*RequestMetaData{}
	h.ctx = nil
	h.cancelCtx = nil
}

func (h *requestHeap) checkRequestContexts() {
	// We gather expired nodes while holding the lock, remove them from the internal slice,
	// then notify their ErrCh channels without holding the lock to avoid blocking other
	// operations on the heap.
	h.mu.Lock()

	var expired []*RequestMetaData

	// iterate and remove expired nodes in-place
	for i := 0; i < len(h.heap); {
		node := h.heap[i]
		// consider nil context as expired (defensive)
		if node.Ctx == nil || node.Ctx.Err() != nil {
			expired = append(expired, node)

			// remove element i by swapping with last and shrinking slice
			last := len(h.heap) - 1
			h.heap[i] = h.heap[last]
			h.heap = h.heap[:last]
			// do not increment i, re-evaluate swapped-in element
			continue
		}
		i++
	}

	h.mu.Unlock()

	// Notify expired nodes outside lock. For each expired node send a RequestContextExpiredError
	// and then close the channel. We perform the send inside a goroutine to avoid blocking the
	// requestHeap ticker if a receiver isn't actively reading.
	for _, n := range expired {
		// capture channel and context error
		ch := n.HeapErrCh
		ctxErr := error(nil)
		if n.Ctx != nil {
			ctxErr = n.Ctx.Err()
		}

		go func(ch chan error, err error) {
			// Send the RequestContextExpiredError if possible and then close channel to signal
			// completion. If err is nil, still close the channel to indicate removal.
			if err != nil {
				// best-effort send; block until receiver takes it so they reliably observe the error
				ch <- &RequestContextExpiredError{Err: err}
			}
			close(ch)
		}(ch, ctxErr)
	}
}

func (h *requestHeap) AddNode(node *RequestMetaData) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.ctx == nil {
		return &RequestHeapNotRunningError{
			CallContext: "AddNode",
		}
	}

	h.heap = append(h.heap, node)

	return nil
}

func (h *requestHeap) RemoveNode(id RequestId) (*RequestMetaData, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.ctx == nil {
		return nil, &RequestHeapNotRunningError{
			CallContext: "RenoveNode",
		}
	}

	for i := 0; i < len(h.heap); i++ {
		if h.heap[i].Id == id {
			node := h.heap[i]

			// remove element i by swapping with last and shrinking
			last := len(h.heap) - 1
			h.heap[i] = h.heap[last]
			h.heap = h.heap[:last]

			// close the heap error channel with no error
			close(node.HeapErrCh)

			return node, nil
		}
	}

	return nil, &RequestHeapNodeNotIncludedError{
		CallContext: "RemoveNode",
		Id:          id,
	}
}
