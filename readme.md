# ctrader — Go client for cTrader OpenAPI

This repository provides a compact, idiomatic Go client for the cTrader
OpenAPI proxy. It focuses on a small, testable public surface for:

- sending RPC-like requests (request/response mapped to protobuf messages),
- subscribing to and listening for streaming server events (push-style),
- a concurrency-safe event dispatcher and small adapters to wire typed
	handlers without reflection.

The library deliberately keeps runtime behaviour explicit and conservative
— networking resilience and reconnect logic are provided but intentionally
simple so they are easy to reason about and test.

Contents
- Quick highlights
- Installation (Go modules)
- Exported types & API summary
- Examples (connect/request, subscribe + listen, client events)
- Running tests
- License

## Quick highlights

- Typed, protobuf-backed request/response helpers (package `messages`).
- Streaming event subscription and listening API: `SubscribeEvent`,
	`UnsubscribeEvent`, `ListenToEvent` on the `APIClient` implementation.
- Concurrency-safe event dispatch implemented by `datatypes.EventHandler`.
- Small typed adapters (`SpawnEventHandler`, `SpawnClientEventHandler`) that
	adapt generic channels to typed callbacks (they perform runtime type
	assertions and will drop non-matching events).

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

- `APIClient` — main interface. Important methods:
	- `Connect() error` / `Disconnect() error`
	- `SendRequest(RequestData) error` — sends a protobuf-typed request and
		unmarshals the response into the provided response object.
	- `SubscribeEvent(SubscribableEventData)` / `UnsubscribeEvent(...)` —
		subscribe/unsubscribe for server-side subscription-based events.
	- `ListenToEvent(eventType, chan ListenableEvent, ctx)` — register a
		long-running listener channel for push events.
	- `ListenToClientEvent(clientEventType, chan ListenableClientEvent, ctx)`
		— listen for client lifecycle events (connection loss, reconnect success
		and reconnect fail).

- `ApplicationCredentials{ClientId, ClientSecret}` — credentials used by
	the application to authenticate to the OpenAPI. Validate with
	`CheckError()`.

- `APIClientConfig` — runtime configuration (queue buffer sizes and the
	request heap iteration timeout). Use `DefaultAPIClientConfig()` to obtain
	sensible defaults and call `WithConfig` on the client before `Connect`.

- Subscription helpers (used with `SubscribeEvent` / `UnsubscribeEvent`):
	- `SubscribableEventData{EventType, SubcriptionData}`
	- `SubscriptionDataSpotEvent{ CtraderAccountId, SymbolIds []int64 }`
	- `SubscriptionDataLiveTrendbarEvent{ CtraderAccountId, SymbolId, Period }`
	- `SubscriptionDataDepthQuoteEvent{ CtraderAccountId, SymbolIds }`

- Event adapters & helpers:
	- `ListenableEvent` and `ListenableClientEvent` are marker interfaces.
	- `CastToEventType[T]` and `CastToClientEventType[T]` helpers to cast
		generic events to typed values.
	- `SpawnEventHandler` and `SpawnClientEventHandler` start small goroutines
		that forward typed events to your handler.

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

	client = client.WithConfig(ctraderopenapi.DefaultAPIClientConfig())
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
	fmt.Println("version response:", res)
}
```

2) Subscribe to spot events and handle them with the typed adapter

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

onEventCh := make(chan ctraderopenapi.ListenableEvent)
if err := client.ListenToEvent(ctraderopenapi.EventType_Spots, onEventCh, ctx); err != nil {
	panic(err)
}

_ = ctraderopenapi.SpawnEventHandler(ctx, onEventCh, func(e *ctraderopenapi.ProtoOASpotEvent) {
	fmt.Println("spot event:", e)
})

sub := ctraderopenapi.SubscribableEventData{
	EventType: ctraderopenapi.EventType_Spots,
	SubcriptionData: &ctraderopenapi.SubscriptionDataSpotEvent{
		CtraderAccountId: ctraderopenapi.CtraderAccountId(123456),
		SymbolIds:        []int64{1, 2, 3},
	},
}
if err := client.SubscribeEvent(sub); err != nil {
	panic(err)
}
```

3) Listen for client lifecycle events

```go

clientCh := make(chan ctraderopenapi.ListenableClientEvent)
if err := client.ListenToClientEvent(ctraderopenapi.ClientEventType_ReconnectSuccessEvent, clientCh, context.Background()); err != nil {
	panic(err)
}
ctraderopenapi.SpawnClientEventHandler(context.Background(), clientCh, func(e *ctraderopenapi.ReconnectSuccessEvent) {
	fmt.Println("reconnected")
})
```

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

