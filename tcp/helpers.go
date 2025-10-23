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
	"net"
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
