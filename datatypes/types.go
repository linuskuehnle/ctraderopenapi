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
	"github.com/google/uuid"
	"github.com/linuskuehnle/ctraderopenapi/messages"

	"context"
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
)

type RequestId string

type RequestData struct {
	Ctx     context.Context
	ReqType messages.ProtoOAPayloadType
	Req     proto.Message
	ResType messages.ProtoOAPayloadType
	Res     proto.Message
}

type ResponseData struct {
	ProtoMsg    *messages.ProtoMessage
	PayloadType messages.ProtoOAPayloadType
}

type RequestMetaData struct {
	*RequestData

	Id RequestId

	ErrCh     chan error
	ResDataCh chan *ResponseData
}

func NewRequestMetaData(reqData *RequestData, errCh chan error, resDataCh chan *ResponseData) (*RequestMetaData, error) {
	if reqData == nil {
		return nil, &FunctionInvalidArgError{
			FunctionName: "NewRequestMetaData",
			Err:          errors.New("reqData mustn't be nil"),
		}
	}
	if errCh == nil {
		return nil, &FunctionInvalidArgError{
			FunctionName: "NewRequestMetaData",
			Err:          errors.New("errCh mustn't be nil"),
		}
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate request ID: %w", err)
	}

	r := RequestMetaData{
		RequestData: reqData,
		Id:          RequestId(id.String()),
		ErrCh:       errCh,
		ResDataCh:   resDataCh,
	}

	return &r, nil
}
