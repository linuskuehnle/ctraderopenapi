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
	"crypto/tls"
	"crypto/x509"
	"fmt"
)

func NewSystemCertTLSConfig() (*tls.Config, error) {
	// Load system CA certificates
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to load system CA certificates: %w", err)
	}

	return &tls.Config{
		RootCAs:    rootCAs,
		MinVersion: tls.VersionTLS12,
	}, nil
}
