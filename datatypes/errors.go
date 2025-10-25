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
	"fmt"
)

type FunctionInvalidArgError struct {
	FunctionName string
	Err          error
}

func (e *FunctionInvalidArgError) Error() string {
	return fmt.Sprintf("invalid argument on function %s: %v", e.FunctionName, e.Err)
}

type RequestContextExpiredError struct {
	Err error
}

func (e *RequestContextExpiredError) Error() string {
	return fmt.Sprintf("request context expired: %v", e.Err)
}

type IdAlreadyIncludedError struct {
	Id        EventId
	EventName string
}

func (e *IdAlreadyIncludedError) Error() string {
	return fmt.Sprintf("id %d is already included in event handler %s", e.Id, e.EventName)
}

type IdNotIncludedError struct {
	Id EventId
}

func (e *IdNotIncludedError) Error() string {
	return fmt.Sprintf("id %d is not included in event handler", e.Id)
}

type LifeCycleAlreadyRunningError struct {
	CallContext string
}

func (e *LifeCycleAlreadyRunningError) Error() string {
	return fmt.Sprintf("%s: life cycle already running", e.CallContext)
}

type LifeCycleNotRunningError struct {
	CallContext string
}

func (e *LifeCycleNotRunningError) Error() string {
	return fmt.Sprintf("%s: life cycle not running", e.CallContext)
}

type RequestHeapAlreadyRunningError struct {
	CallContext string
}

func (e *RequestHeapAlreadyRunningError) Error() string {
	return fmt.Sprintf("%s: request heap already running", e.CallContext)
}

type RequestHeapNotRunningError struct {
	CallContext string
}

func (e *RequestHeapNotRunningError) Error() string {
	return fmt.Sprintf("%s: request heap not running", e.CallContext)
}

type RequestHeapNodeNotIncludedError struct {
	CallContext string
	Id          RequestId
}

func (e *RequestHeapNodeNotIncludedError) Error() string {
	return fmt.Sprintf("%s: heap node with request id %s is not included", e.CallContext, e.Id)
}
