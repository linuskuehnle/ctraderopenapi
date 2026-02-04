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
	"errors"
	"sync"
	"time"
)

type RateLimiter interface {
	SetPenalty(duration time.Duration)
	WaitForPermit()
}

type rateLimiter struct {
	mu sync.RWMutex

	n          int           // Number of requests allowed per interval
	interval   time.Duration // Interval in which n requests are allowed
	hitTimeout time.Duration // Duration to pause whenever rate limit is being hit until next check iteration

	reqTimeList LinkedList[time.Time]

	noPermitUntil time.Time
}

func NewRateLimiter(n int, interval, hitTimeout time.Duration) (RateLimiter, error) {
	if n <= 0 {
		return nil, &FunctionInvalidArgError{
			FunctionName: "NewRateLimiter",
			Err:          errors.New("n must be at least 1"),
		}
	}
	if interval <= 0 {
		return nil, &FunctionInvalidArgError{
			FunctionName: "NewRateLimiter",
			Err:          errors.New("interval must be greater than 0"),
		}
	}
	if hitTimeout <= 0 {
		return nil, &FunctionInvalidArgError{
			FunctionName: "NewRateLimiter",
			Err:          errors.New("hitTimeout must be greater than 0"),
		}
	}

	return newRateLimiter(n, interval, hitTimeout), nil
}

func newRateLimiter(n int, interval, hitTimeout time.Duration) *rateLimiter {
	return &rateLimiter{
		mu: sync.RWMutex{},

		n:          n,
		interval:   interval,
		hitTimeout: hitTimeout,

		reqTimeList: NewLinkedList[time.Time]().
			WithAllowDuplicates(),
	}
}

func (d *rateLimiter) SetPenalty(duration time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.noPermitUntil = time.Now().Add(duration)
}

func (d *rateLimiter) WaitForPermit() {
	d.mu.Lock()
	defer d.mu.Unlock()

	defer d.addRequestTime(time.Now()) // Add request permit time to request time list on return

	if d.reqTimeList.Length() < d.n && d.noPermitUntil.IsZero() {
		return
	}

	limitT := d.getLimitTime()
	for time.Now().Before(limitT.Add(d.interval)) {
		// Rate limit is being hit. wait for duration d.hitTimeout until next check interation
		time.Sleep(d.hitTimeout)
	}
}

func (d *rateLimiter) getLimitTime() time.Time {
	limitT, _ := d.reqTimeList.PeekHead()

	if limitT.Before(d.noPermitUntil) {
		limitT = d.noPermitUntil.Add(-d.interval)
		d.noPermitUntil = time.Time{}
	}

	return limitT
}

func (d *rateLimiter) addRequestTime(t time.Time) {
	d.reqTimeList.Append(t)
	if d.reqTimeList.Length() > d.n {
		d.reqTimeList.PopHead()
	}
}
