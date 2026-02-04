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
	"errors"
	"fmt"
	"sync"
	"time"
)

type RetryBackoff interface {
	Reset()

	Backoff()
	WaitForPermit(context.Context) bool
}

type retryBackoff struct {
	mu sync.RWMutex

	ladder        []time.Duration
	stepDownAfter time.Duration

	ladderStep int
	gateTime   time.Time
}

func NewRetryBackoff(ladder []time.Duration, stepDownAfter time.Duration) (RetryBackoff, error) {
	if len(ladder) == 0 {
		return nil, &FunctionInvalidArgError{
			FunctionName: "NewRetryBackoff",
			Err:          errors.New("ladder must not be empty"),
		}
	}

	prev := ladder[0]
	for i, d := range ladder {
		if d < prev {
			return nil, &FunctionInvalidArgError{
				FunctionName: "NewRetryBackoff",
				Err:          fmt.Errorf("duration at ladder index %d must be higher than predecessor duration", i),
			}
		}
		prev = d
	}

	return newRetryBackoff(ladder, stepDownAfter), nil
}

func newRetryBackoff(ladder []time.Duration, stepDownAfter time.Duration) *retryBackoff {
	cpyLadder := make([]time.Duration, len(ladder))
	copy(cpyLadder, ladder)

	return &retryBackoff{
		mu: sync.RWMutex{},

		ladder:        cpyLadder,
		stepDownAfter: stepDownAfter,

		gateTime: time.Now(),
	}
}

func (d *retryBackoff) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.gateTime = time.Now()
	d.ladderStep = 0
}

func (b *retryBackoff) WaitForPermit(ctx context.Context) bool {
	b.mu.Lock()

	b.updateLadderStep()
	b.updateGateTime()

	gateTimeDiff := time.Until(b.gateTime)
	b.mu.Unlock()

	for gateTimeDiff > 0 {
		select {
		case <-ctx.Done():
			return false
		case <-time.After(gateTimeDiff):
			b.mu.RLock()
			gateTimeDiff = time.Until(b.gateTime)
			b.mu.RUnlock()
		}
	}

	return true
}

func (b *retryBackoff) Backoff() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.updateLadderStep()

	if b.ladderStep < len(b.ladder) {
		b.ladderStep++
	}

	b.updateGateTime()
}

func (b *retryBackoff) updateLadderStep() {
	diff := time.Since(b.gateTime)
	if diff < b.stepDownAfter {
		return
	}

	numStepsDown := 0
	for diff > b.stepDownAfter {
		diff -= b.stepDownAfter
		numStepsDown++
	}

	b.ladderStep -= numStepsDown
	if b.ladderStep < 0 {
		b.ladderStep = 0
	}
}

func (b *retryBackoff) updateGateTime() {
	if b.ladderStep == 0 {
		// Pull gate time along for efficiency
		b.gateTime = time.Now()
		return
	}

	d := b.ladder[b.ladderStep-1]
	b.gateTime = time.Now().Add(d)
}
