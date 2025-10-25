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
	"fmt"
	"sync"
)

type RequestQueue interface {
	WithDataCallbackChan(chan struct{}) RequestQueue
	Clear()
	Enqueue(*RequestMetaData) error
	Dequeue() (*RequestMetaData, error)
	RemoveFromQueue(RequestId) error
}

type requestQueue struct {
	mu       sync.RWMutex
	list     LinkedList[RequestId]
	reqMap   map[RequestId]*RequestMetaData
	onDataCh chan struct{}
}

func NewRequestQueue() RequestQueue {
	return newRequestQueue()
}

func newRequestQueue() *requestQueue {
	return &requestQueue{
		mu:     sync.RWMutex{},
		list:   NewLinkedList[RequestId](),
		reqMap: make(map[RequestId]*RequestMetaData),
	}
}

func (q *requestQueue) WithDataCallbackChan(onDataCh chan struct{}) RequestQueue {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.onDataCh = onDataCh
	return q
}

func (q *requestQueue) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.list.Clear()

	for reqId := range q.reqMap {
		delete(q.reqMap, reqId)
	}

	if q.onDataCh == nil {
		return
	}

	close(q.onDataCh)
	q.onDataCh = nil
}

func (q *requestQueue) Enqueue(r *RequestMetaData) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if err := q.list.Append(r.Id); err != nil {
		return err
	}

	q.reqMap[r.Id] = r

	if q.list.Length() == 1 && q.onDataCh != nil {
		q.onDataCh <- struct{}{}
	}

	return nil
}

func (q *requestQueue) Dequeue() (*RequestMetaData, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.dequeueNoLock()
}

func (q *requestQueue) dequeueNoLock() (*RequestMetaData, error) {
	if q.list.IsEmpty() {
		return nil, fmt.Errorf("error dequeueing: queue is empty")
	}
	id, _ := q.list.Pop()

	reqData := q.reqMap[id]
	delete(q.reqMap, id)

	if q.list.Length() > 0 && q.onDataCh != nil {
		q.onDataCh <- struct{}{}
	}

	return reqData, nil
}

func (q *requestQueue) RemoveFromQueue(id RequestId) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if err := q.list.Remove(id); err != nil {
		return err
	}

	delete(q.reqMap, id)

	// RemoveFromQueue should not signal data since its apart from the usual enqueue/dequeue dynamic

	return nil
}
