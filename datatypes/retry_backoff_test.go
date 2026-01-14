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
	"testing"
	"time"
)

func TestRetryBackoffInvalidArgs(t *testing.T) {
	_, err := NewRetryBackoff([]time.Duration{}, time.Second)
	if err == nil {
		t.Fatalf("expected error for empty ladder")
	}

	_, err = NewRetryBackoff([]time.Duration{time.Second * 2, time.Second}, time.Second*10)
	if err == nil {
		t.Fatalf("expected error for non-monotonic ladder")
	}
}

func TestRetryBackoffSequences(t *testing.T) {
	ladder := []time.Duration{
		time.Millisecond * 10, // Step 1
		time.Millisecond * 20, // Step 2
		time.Millisecond * 40, // Step 3
		time.Millisecond * 80, // Step 4
	}
	stepDownAfter := time.Millisecond * 100

	delta := time.Millisecond * 5

	b := newRetryBackoff(ladder, stepDownAfter)

	ctx := context.Background()

	// Sequence 1: Immediate pass
	now := time.Now()
	b.WaitForPermit(ctx)
	waitTime := time.Since(now)
	if waitTime > delta {
		t.Fatalf("should not wait. waited for %dms", waitTime.Milliseconds())
	}

	// Sequence 2: Ladder step 1
	b.Backoff()
	now = time.Now()
	b.WaitForPermit(ctx)
	waitTime = time.Since(now)
	if waitTime < ladder[0] {
		t.Fatalf("did not wait. should wait for ladder step 1: %dms", ladder[0].Milliseconds())
	}

	// Sequence 2: Ladder step 3
	b.Backoff()
	b.Backoff()
	now = time.Now()
	b.WaitForPermit(ctx)
	waitTime = time.Since(now)
	if waitTime < ladder[2] {
		t.Fatalf("did not wait. should wait for ladder step 3: %dms", ladder[2].Milliseconds())
	}
	if waitTime > ladder[2]+delta {
		t.Fatalf("waited for too long. should wait for ladder step 3: %dms. waited for %dms", ladder[2].Milliseconds(), waitTime.Milliseconds())
	}

	// Sequence 3: Wait for two steps down. Ladder step 1
	time.Sleep(stepDownAfter * 2)
	now = time.Now()
	b.WaitForPermit(ctx)
	waitTime = time.Since(now)
	if waitTime < ladder[0] {
		t.Fatalf("did not wait. should wait for ladder step 1: %dms", ladder[0].Milliseconds())
	}
	if waitTime > ladder[0]+delta {
		t.Fatalf("waited for too long. should wait for ladder step 1: %dms. waited for %dms", ladder[0].Milliseconds(), waitTime.Milliseconds())
	}

	// Sequence 4: Wait for ladder reset. Ladder step 0
	time.Sleep(stepDownAfter)
	now = time.Now()
	b.WaitForPermit(ctx)
	waitTime = time.Since(now)
	if waitTime > delta {
		t.Fatalf("should not wait. waited for %dms", waitTime.Milliseconds())
	}

	// Sequence 5: Ladder step 4 after 5 backoffs
	b.Backoff()
	b.Backoff()
	b.Backoff()
	b.Backoff()
	b.Backoff()
	now = time.Now()
	b.WaitForPermit(ctx)
	waitTime = time.Since(now)
	if waitTime < ladder[3] {
		t.Fatalf("did not wait. should wait for ladder step 4: %dms", ladder[3].Milliseconds())
	}
	if waitTime > ladder[3]+delta {
		t.Fatalf("waited for too long. should wait for ladder step 4: %dms. waited for %dms", ladder[3].Milliseconds(), waitTime.Milliseconds())
	}

	// Sequence 6: Wait for ladder reset. Ladder step 0
	time.Sleep(stepDownAfter * 4)
	now = time.Now()
	b.WaitForPermit(ctx)
	waitTime = time.Since(now)
	if waitTime > delta {
		t.Fatalf("should not wait. waited for %dms", waitTime.Milliseconds())
	}

	// Sequence 7: Ladder step 4 with ongoing backoffs
	b.Backoff()
	b.Backoff()
	b.Backoff()
	b.Backoff()

	go func() {
		for range 4 {
			time.Sleep(ladder[3] / 2)
			b.Backoff()
		}
	}()

	now = time.Now()
	b.WaitForPermit(ctx)
	waitTime = time.Since(now)
	if waitTime < ladder[3]*3 {
		t.Fatalf("did not wait. should wait for ladder step 4 and 4 concurrent backoffs: %dms", (ladder[3] * 3).Milliseconds())
	}
	if waitTime > ladder[3]*3+delta {
		t.Fatalf("waited for too long. should wait for ladder step 4 and 4 concurrent backoffs: %dms. waited for %dms", (ladder[3] * 3).Milliseconds(), waitTime.Milliseconds())
	}

	// Sequence 8: Wait for ladder reset. Ladder step 0
	time.Sleep(stepDownAfter * 4)
	now = time.Now()
	b.WaitForPermit(ctx)
	waitTime = time.Since(now)
	if waitTime > delta {
		t.Fatalf("should not wait. waited for %dms", waitTime.Milliseconds())
	}
}
