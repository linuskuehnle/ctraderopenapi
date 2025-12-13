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
)

type node[T comparable] struct {
	value *T
	next  *node[T]
	prev  *node[T]
}

func (n *node[T]) isHead() bool {
	return n.prev == nil
}

func (n *node[T]) isTail() bool {
	return n.next == nil
}

type LinkedList[T comparable] interface {
	WithAllowDuplicates() LinkedList[T]

	Clear()
	Has(t T) bool
	IsEmpty() bool
	Length() int
	Append(t T) error
	Prepend(t T) error
	PopHead() (T, error)
	PopTail() (T, error)
	Remove(t T) error
	PeekHead() (T, error)
	PeekTail() (T, error)
}

// No mutex locks here as its used only in request queue
// which is synchronized by mutexes
type linkedList[T comparable] struct {
	head   *node[T]
	tail   *node[T]
	length int

	allowDuplicates bool
}

func NewLinkedList[T comparable]() LinkedList[T] {
	return newLinkedList[T]()
}

func newLinkedList[T comparable]() *linkedList[T] {
	return &linkedList[T]{}
}

func (l *linkedList[T]) WithAllowDuplicates() LinkedList[T] {
	if !l.IsEmpty() {
		return l
	}

	l.allowDuplicates = true
	return l
}

func (l *linkedList[T]) Clear() {
	l.length = 0
	l.head = nil
	l.tail = nil

	l.allowDuplicates = false
}

func (l *linkedList[T]) Has(t T) bool {
	return l.find(t) != nil
}

func (l *linkedList[T]) IsEmpty() bool {
	return l.length == 0
}

func (l *linkedList[T]) Length() int {
	return l.length
}

func (l *linkedList[T]) find(t T) *node[T] {
	i := l.head
	for i != nil {
		if *i.value == t {
			return i
		}
		i = i.next
	}
	return nil
}

func (l *linkedList[T]) Append(t T) error {
	if !l.allowDuplicates && l.Has(t) {
		return fmt.Errorf("duplicate: cannot append: %v", t)
	}

	l.length++

	n := &node[T]{
		value: &t,
	}

	if l.tail == nil {
		l.tail = n
		l.head = n
		return nil
	}

	l.tail.next = n
	n.prev = l.tail

	l.tail = n
	return nil
}

func (l *linkedList[T]) Prepend(t T) error {
	if !l.allowDuplicates && l.Has(t) {
		return fmt.Errorf("duplicate: cannot prepend: %v", t)
	}

	l.length++

	n := &node[T]{
		value: &t,
	}

	if l.head == nil {
		l.head = n
		l.tail = n
		return nil
	}

	l.head.prev = n
	n.next = l.head

	l.head = n
	return nil
}

func (l *linkedList[T]) PopHead() (T, error) {
	if l.length == 0 {
		var zero T
		return zero, fmt.Errorf("cannot pop head from empty list")
	}

	t := *l.head.value

	l.popHead()

	return t, nil
}

func (l *linkedList[T]) PopTail() (T, error) {
	if l.length == 0 {
		var zero T
		return zero, fmt.Errorf("cannot pop tail from empty list")
	}

	t := *l.tail.value

	l.popTail()

	return t, nil
}

func (l *linkedList[T]) popHead() {
	l.length--

	head := l.head

	l.head = l.head.next
	if l.head != nil /* equivalent to l.length != 0 */ {
		l.head.prev = nil
	} else /* equivalent to l.length == 0 */ {
		l.tail = nil
	}

	// Garbage collector speed-up
	head.value = nil
	head.next = nil
}

func (l *linkedList[T]) popTail() {
	l.length--

	tail := l.tail

	l.tail = l.tail.prev
	if l.tail != nil /* equivalent to l.length != 0 */ {
		l.tail.next = nil
	} else /* equivalent to l.length == 0 */ {
		l.head = nil
	}

	// Garbage collector speed-up
	tail.value = nil
	tail.prev = nil
}

func (l *linkedList[T]) Remove(t T) error {
	if l.allowDuplicates {
		return fmt.Errorf("remove: operation not usable when allowing duplicates")
	}

	n := l.find(t)
	if n == nil {
		return fmt.Errorf("remove: not in list: %v", t)
	}

	if n.isHead() {
		l.popHead()
		return nil
	}
	if n.isTail() {
		l.popTail()
		return nil
	}

	l.length--
	n.next.prev = n.prev
	n.prev.next = n.next

	// Garbage collector speed-up
	n.value = nil
	n.next = nil
	n.prev = nil

	return nil
}

func (l *linkedList[T]) PeekHead() (T, error) {
	if l.length == 0 {
		var zero T
		return zero, fmt.Errorf("cannot peek head on empty list")
	}

	t := *l.head.value
	return t, nil
}

func (l *linkedList[T]) PeekTail() (T, error) {
	if l.length == 0 {
		var zero T
		return zero, fmt.Errorf("cannot peek tail on empty list")
	}

	t := *l.tail.value
	return t, nil
}
