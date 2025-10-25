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
	"errors"
	"io"
	"net"
	"strings"
	"syscall"
)

func splitAddress(address string) (string, string, error) {
	if address == "" {
		return "", "", &InvalidAddressError{
			Address:   address,
			ErrorText: "address cannot be empty",
		}
	}

	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return "", "", &InvalidAddressError{
			Address:   address,
			ErrorText: "invalid format, expected 'host:port'",
		}
	}

	if net.ParseIP(host) == nil && host == "" {
		return "", "", &InvalidAddressError{
			Address:   address,
			ErrorText: "host part of the address cannot be empty",
		}
	}

	if port == "" {
		return "", "", &InvalidAddressError{
			Address:   address,
			ErrorText: "port part of the address cannot be empty",
		}
	}

	return host, port, nil
}

// isNetworkFatal determines whether an error means the connection is lost
// and should trigger a reconnection attempt.
func isNetworkFatal(err error) bool {
	if err == nil {
		return false
	}

	// Graceful remote close or unexpected EOF
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return true
	}

	// Underlying connection reset, broken pipe, or other low-level errors
	if errors.Is(err, syscall.ECONNRESET) ||
		errors.Is(err, syscall.EPIPE) ||
		strings.Contains(err.Error(), "connection reset by peer") ||
		strings.Contains(err.Error(), "broken pipe") {
		return true
	}

	// Network errors may be temporary/retriable
	var netErr net.Error
	if errors.As(err, &netErr) {
		// Consider temporary errors and timeouts as potentially recoverable
		return !netErr.Timeout()
	}

	// Any unrecognized error is not considered network-fatal
	return false
}
