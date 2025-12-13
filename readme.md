# ctrader — Go client for cTrader OpenAPI

This repository provides a compact, idiomatic Go client for the cTrader
OpenAPI proxy. It focuses on a small, testable public surface for:

- sending RPC-like requests (request/response mapped to protobuf messages),
- subscribing to and listening for streaming server events (push-style),
- a concurrency-safe event dispatcher and small adapters to wire typed
	handlers without reflection.
- an client-side request rate limiter to prevent server-side rate limiting.

Contents
- Quick highlights
- Installation (Go modules)
- Exported types & API summary
- Examples (connect/request, subscribe + listen, client events)
- Running tests
- License

## Quick highlights

- Typed, protobuf-backed request/response helpers.
- Event handler for both API and client events providing easy abstraction via
	subscribe/unsubscribe and listen functions.
- Client-side rate limiter to avoid server-side rate limiting.
- Small typed adapters (`SpawnAPIEventHandler`, `SpawnClientEventHandler`) that
	adapt generic channels to typed callbacks (they perform runtime type
	assertions and will drop non-matching events).
- Connection loss handling; recovering authenticated accounts and active
	subscriptions on reconnect.

## Installation

This library uses Go modules. From your project directory you can add the
dependency with:

```powershell
go get github.com/linuskuehnle/ctraderopenapi@latest
```

Then import the package in your code:

```go
import "github.com/linuskuehnle/ctraderopenapi"
```

To install a binary (if provided by this repository) or to run tooling:

```powershell
go install github.com/linuskuehnle/ctraderopenapi@latest
```

## Exported types & API summary

Key public types and functions (see `types.go` and `api_client.go` for
full comments and examples):

- `NewAPIClient(cred ApplicationCredentials, env Environment) (APIClient, error)`
	— create a new client instance. Call `WithConfig` before `Connect` to
	adjust runtime buffers/timeouts.

- `APIClient` — main interface.
	Important methods:
	- `Connect() error` / `Disconnect()`
	- `SendRequest(RequestData) error` — sends a protobuf-typed request and
		unmarshals the response into the provided response object.
	- `SubscribeAPIEvent(SubscribableAPIEventData)` / `UnsubscribeAPIEvent(...)` —
		subscribe/unsubscribe for server-side subscription-based events.
	- `SubscribeClientEvent(SubscribableClientEventData)` / `UnsubscribeClientEvent(...)` —
		subscribe/unsubscribe for client-side events.
	- `ListenToAPIEvent(ctx, eventType, chan ListenableAPIEvent)` — register a
		long-running listener channel for push events.
	- `ListenToClientEvent(ctx, clientEventType, chan ListenableClientEvent)`
		— listen for client-side events
	
	Client-side events are e.g. fatal client errors, connection loss, reconnect success
	and reconnect fail.
	
	Fatal (non-recoverable) client errors will be emitted as client event of type
	FatalErrorEvent. If there is a listener channel installed, the error will be sent to
	the channel and the API client recovers by dropping the tcp connection and running the
	reconnect sequence.
	If there is no listener channel installed for the fatal error event, a fatal error will
	be raised as a panic instead.

	The API Client takes care of connection losses. You can listen to reconnect events
	(ConnectionLossEvent, ReconnectSuccessEvent, ReconnectFailEvent) using ListenToClientEvent.
	API Events subscribed to will automatically be resubscribed before ReconnectSuccessEvent is
	emitted.

	Runtime configuration:
	- `WithQueueBufferSize(int)` updates the number of queued requests that may be buffered
		by the internal request queue before backpressure applies.
	- `WithTCPMessageBufferSize(int)` updates the size of the channel used to receive inbound
		TCP messages from the network reader.
	- `WithRequestHeapIterationTimeout(time.Duration)` updates interval used by the request heap
		to periodically check for expired request contexts.
	- `DisableDefaultRateLimiter()` to disable the client-side rate limiter.

	The rate limiter ensures that the 50 live requests per second and 5 historical requests per
	second are not exceeded. In case there still occurs a server-side rate limit due to network
	divergence (very unlikely) it enqueues the request again, so the caller does not have to worry
	about it.

- `ApplicationCredentials{ ClientId, ClientSecret }` — credentials used by the application to
	authenticate to the OpenAPI. Validate with `CheckError()`.

- Event adapters & helpers:
	- `ListenableEvent` and `ListenableClientEvent` are marker interfaces.
	- `CastToEventType[T]` and `CastToClientEventType[T]` helpers to cast the generic event type
		to the concrete event types.
	- `SpawnAPIEventHandler` and `SpawnClientEventHandler` start small goroutines that forward
		typed events to your handler.

## Examples

1) Basic connect and simple request

```go
package main

import (
	"context"
	"fmt"
	"github.com/linuskuehnle/ctraderopenapi"
)

func main() {
	cred := ctraderopenapi.ApplicationCredentials{ClientId: "id", ClientSecret: "secret"}
	client, err := ctraderopenapi.NewAPIClient(cred, ctraderopenapi.Environment_Demo)
	if err != nil {
		panic(err)
	}

	if err := client.Connect(); err != nil {
		panic(err)
	}
	defer client.Disconnect()

	req := ctraderopenapi.ProtoOAVersionReq{}
	var res ctraderopenapi.ProtoOAVersionRes
	reqData := ctraderopenapi.RequestData{
		Ctx:     context.Background(),
		ReqType: ctraderopenapi.PROTO_OA_VERSION_REQ,
		Req:     &req,
		ResType: ctraderopenapi.PROTO_OA_VERSION_RES,
		Res:     &res,
	}
	if err := client.SendRequest(reqData); err != nil {
		fmt.Println("request failed:", err)
		return
	}

	fmt.Println("version:", res.GetVersion())
}
```

2) Subscribe to spot events and handle them with the typed adapter

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

onEventCh := make(chan ctraderopenapi.ListenableEvent)
if err := client.ListenToAPIEvent(ctx, ctraderopenapi.APIEventType_Spots, onEventCh); err != nil {
	panic(err)
}

_ = ctraderopenapi.SpawnAPIEventHandler(ctx, onEventCh, func(e *ctraderopenapi.ProtoOASpotEvent) {
	fmt.Println("spot event:", e)
})

sub := ctraderopenapi.SubscribableAPIEventData{
	EventType: ctraderopenapi.EventType_Spots,
	SubcriptionData: &ctraderopenapi.SubscriptionDataSpotEvent{
		CtraderAccountId: ctraderopenapi.CtraderAccountId(123456),
		SymbolIds:        []int64{1, 2, 3},
	},
}
if err := client.SubscribeAPIEvent(sub); err != nil {
	panic(err)
}
```

3) Listen for client events; e.g. ReconnectSuccessEvent:

```go

clientCh := make(chan ctraderopenapi.ListenableClientEvent)
if err := client.ListenToClientEvent(context.Background(), ctraderopenapi.ClientEventType_ReconnectSuccessEvent, clientCh); err != nil {
	panic(err)
}
ctraderopenapi.SpawnClientEventHandler(context.Background(), clientCh, func(e *ctraderopenapi.ReconnectSuccessEvent) {
	fmt.Println("reconnected")
})
```

I suggest registering all event listener channels before calling apiClient.Connect().
apiClient.Disconnect() unregisters all event listener channels, hence if you want to
connect again and use the same listeners you explicitly need to register them again.

## Running tests

Some tests exercise live interactions and expect credentials via environment
variables or a `.env` file (see `stubs_test.go`). To run unit/integration
tests locally:

```powershell
go test ./...
```

To run a focused test and vet the code:

```powershell
go vet ./...
go test ./... -run TestClientConnectDisconnect
```

## License

This project is licensed under the Apache License 2.0 — see the
[LICENSE](./LICENSE) file for details.

