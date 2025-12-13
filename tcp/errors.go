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

package tcp

import (
	"fmt"
)

/*
	client.go
*/

type OpenConnectionError struct {
	ErrorText string
}

func (e *OpenConnectionError) Error() string {
	return fmt.Sprintf("error opening connection: %s", e.ErrorText)
}

type CloseConnectionError struct {
	ErrorText string
}

func (e *CloseConnectionError) Error() string {
	return fmt.Sprintf("error closing connection: %s", e.ErrorText)
}

type NoConnectionError struct {
	CallContext string
}

func (e *NoConnectionError) Error() string {
	return fmt.Sprintf("%s: no open connection available", e.CallContext)
}

type OperationBlockedError struct {
	CallContext string
	ErrorText   string
}

func (e *OperationBlockedError) Error() string {
	return fmt.Sprintf("%s operation blocked: %s", e.CallContext, e.ErrorText)
}

type MaxReconnectAttemptsReachedError struct {
	MaxAttempts int
}

func (e *MaxReconnectAttemptsReachedError) Error() string {
	return fmt.Sprintf("maximum reconnect attempts (%d) reached", e.MaxAttempts)
}

/*
	helpers.go
*/

type InvalidAddressError struct {
	Address   string
	ErrorText string
}

func (e *InvalidAddressError) Error() string {
	return fmt.Sprintf("invalid address '%s': %s", e.Address, e.ErrorText)
}
