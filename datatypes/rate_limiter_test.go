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
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	n := 50
	interval := time.Second
	hitTimeout := time.Millisecond * 5

	_, err := NewRateLimiter(n, interval, hitTimeout)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRateLimiterConsistency(t *testing.T) {
	n := 50
	interval := time.Second
	hitTimeout := time.Millisecond * 5

	// Choose a safe delta for all operations
	delta := time.Millisecond * 100

	r := newRateLimiter(n, interval, hitTimeout)

	tNow := time.Now()

	for range n {
		r.WaitForPermit()
	}

	if time.Since(tNow) > delta {
		t.Fatalf("waiting for permit on first %d (n) calls did not pass below time delta", n)
	}

	done := make(chan struct{})
	go func() {
		r.WaitForPermit()
		done <- struct{}{}
	}()

	waitTime := interval - delta
	select {
	case <-time.After(waitTime):
	case <-done:
		t.Fatalf("did not wait at least %dms before giving permit", waitTime.Milliseconds())
	}
}

func TestRateLimiterTimeout(t *testing.T) {
	n := 50
	interval := time.Second
	hitTimeout := time.Millisecond * 5

	// Choose a safe delta for all operations
	delta := time.Millisecond * 100

	r := newRateLimiter(n, interval, hitTimeout)

	r.SetPenalty(interval)

	done := make(chan struct{})
	go func() {
		r.WaitForPermit()
		done <- struct{}{}
	}()

	waitTime := interval - delta
	select {
	case <-time.After(waitTime):
	case <-done:
		t.Fatalf("did not wait at least %dms before giving permit", waitTime.Milliseconds())
	}

	waitTime = delta * 2
	select {
	case <-time.After(waitTime):
		t.Fatalf("expected permit before %dms passed", waitTime.Milliseconds())
	case <-done:
	}
}
