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

package ctraderopenapi

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func LoadAppCredentialsFromEnv() (ApplicationCredentials, error) {
	godotenv.Load()

	cred := ApplicationCredentials{
		ClientId:     os.Getenv("_TEST_OPENAPI_CLIENT_ID"),
		ClientSecret: os.Getenv("_TEST_OPENAPI_CLIENT_SECRET"),
	}
	return cred, cred.CheckError()
}

func loadTestAccountCredentials() (int64, string, Environment, error) {
	godotenv.Load()

	strAccountId := os.Getenv("_TEST_OPENAPI_ACCOUNT_ID")
	strAccessToken := os.Getenv("_TEST_OPENAPI_ACCOUNT_ACCESS_TOKEN")

	if strAccountId == "" {
		return 0, "", 0, fmt.Errorf("account ID env var _TEST_OPENAPI_ACCOUNT_ID not set")
	}
	if strAccessToken == "" {
		return 0, "", 0, fmt.Errorf("access token env var _TEST_OPENAPI_ACCOUNT_ACCESS_TOKEN not set")
	}

	accountId, err := strconv.ParseInt(strAccountId, 10, 64)
	if err != nil {
		return 0, "", 0, fmt.Errorf("error parsing account ID: %v", err)
	}
	envStr := os.Getenv("_TEST_OPENAPI_ACCOUNT_ENVIRONMENT")
	var env Environment
	switch strings.ToUpper(envStr) {
	case "DEMO":
		env = Environment_Demo
	case "LIVE":
		env = Environment_Live
	default:
		return 0, "", 0, fmt.Errorf("invalid account environment in env var _TEST_OPENAPI_ACCOUNT_ENVIRONMENT: %s", envStr)
	}

	return accountId, strAccessToken, env, nil
}

func createApiClient(env Environment) (*apiClient, error) {
	cred, err := LoadAppCredentialsFromEnv()
	if err != nil {
		return nil, fmt.Errorf("invalid credentials: %v", err)
	}

	h, err := newApiClient(cred, env)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %v", err)
	}

	return h, nil
}
