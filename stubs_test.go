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
	"github.com/linuskuehnle/ctraderopenapi/messages"

	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"google.golang.org/protobuf/proto"
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

func createApiHandler(env Environment) (*apiHandler, error) {
	cred, err := LoadAppCredentialsFromEnv()
	if err != nil {
		return nil, fmt.Errorf("invalid credentials: %v", err)
	}

	h, err := newApiHandler(cred, env)
	if err != nil {
		return nil, fmt.Errorf("error creating handler: %v", err)
	}

	return h, nil
}

func autheticateAccount(h *apiHandler, ctid int64, accessToken string) (*ProtoOAAccountAuthRes, error) {
	req := messages.ProtoOAAccountAuthReq{
		CtidTraderAccountId: proto.Int64(ctid),
		AccessToken:         proto.String(accessToken),
	}
	var res messages.ProtoOAAccountAuthRes

	reqCtx := context.Background()

	reqData := RequestData{
		Ctx:     reqCtx,
		ReqType: PROTO_OA_ACCOUNT_AUTH_REQ,
		Req:     &req,
		ResType: PROTO_OA_ACCOUNT_AUTH_RES,
		Res:     &res,
	}

	if err := h.SendRequest(reqData); err != nil {
		return nil, fmt.Errorf("error sending account auth request: %v", err)
	}
	return &res, nil
}
