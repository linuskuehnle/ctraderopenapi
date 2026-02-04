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
	"testing"
	"time"
)

func TestTcpClient(t *testing.T) {
	errCh := make(chan error)
	onFatalErr := func(err error) {
		errCh <- err
	}

	client := newTCPClient("demo.ctraderapi.com:5035", onFatalErr)

	// Test setting timeout
	timeout := 5 * time.Second
	client.WithTimeout(timeout)
	if client.timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, client.timeout)
	}

	// Test setting TLS config
	tlsConfig, err := NewSystemCertTLSConfig()
	if err != nil {
		t.Fatalf("Failed to create TLS config: %v", err)
	}
	client.WithTLS(tlsConfig)

	if !client.useTLS || client.tlsConfig == nil {
		t.Error("Expected TLS to be enabled with a valid config")
	}

	close(errCh)
	if err := <-errCh; err != nil {
		t.Errorf("Unexpected fatal error: %v", err)
	}
}

func TestTcpClientConn(t *testing.T) {
	errCh := make(chan error)
	onFatalErr := func(err error) {
		errCh <- err
	}

	client := newTCPClient("demo.ctraderapi.com:5035", onFatalErr)

	// Test setting timeout
	timeout := 5 * time.Second
	client.WithTimeout(timeout)
	if client.timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, client.timeout)
	}

	// Test setting TLS config
	tlsConfig, err := NewSystemCertTLSConfig()
	if err != nil {
		t.Fatalf("Failed to create TLS config: %v", err)
	}
	client.WithTLS(tlsConfig)

	if !client.useTLS || client.tlsConfig == nil {
		t.Error("Expected TLS to be enabled with a valid config")
	}

	// Test opening connection
	err = client.OpenConn()
	if err != nil {
		t.Errorf("Error opening connection: %v", err)
	}

	if client.conn == nil {
		t.Error("Expected connection to be established, but it is nil")
	}

	err = client.CloseConn()
	if err != nil {
		t.Errorf("Error closing connection: %v", err)
	}

	close(errCh)
	if err := <-errCh; err != nil {
		t.Errorf("Unexpected fatal error: %v", err)
	}
}
