# ctrader — Go client for cTrader OpenAPI

This repository contains a compact, idiomatic Go client for the cTrader OpenAPI. It focuses on providing a clear, testable API surface for sending requests and handling streaming events from the cTrader OpenAPI proxy.

The client intentionally keeps runtime behaviour explicit and minimal: it does not try to hide or magically repair protocol errors. Instead it provides small, well-documented building blocks (an API handler, request/response mapping, and a concurrency-safe event dispatcher) that you can compose in production systems.

> Note: this README describes the features implemented in the codebase; behaviour such as network resilience, reconnect strategies or production hardening are provided as utilities in the code but are intentionally conservative so they remain easy to reason about and test.

## Quick highlights

- Typed, protobuf-backed request/response helpers (package `messages`).
- Streaming event subscription and listening API: `SubscribeEvent`, `UnsubscribeEvent`, `ListenToEvent` on the `APIHandler` implementation.
- Concurrency-safe event dispatch implemented by `datatypes.EventHandler` (uses protobuf messages as the event payload type).
- Small adapters to let callers register typed handlers without reflection; these adapters perform a runtime type assertion and will panic if the handler type doesn't match the event type (caller responsibility).
- Test helpers and small integration-like tests that read credentials from environment variables (see "Running tests").

## Package layout (important files)

- `api_handler.go` — main public API (`APIHandler`) and implementation (`apiHandler`). Connect/disconnect, SendRequest, ListenToEvent, SubscribeEvent, UnsubscribeEvent.
- `messages/` — generated protobuf types and a small `event_assert.go` file that adds marker methods so only approved message types are usable as events.
- `datatypes/` — small reusable data structures: `EventHandler`, `RequestQueue`, `RequestMapper`, a generic linked list, and typed errors used internally.
- `tcp/` — low-level TCP/TLS client code driving the connection to the OpenAPI proxy.

## Installation

Use `go get` (module-aware):

```powershell
go get github.com/linuskuehnle/ctrader@latest
```

Then import in your code:

```go
import "github.com/linuskuehnle/ctraderopenapi"
```

## Basic usage

Below are minimal usage snippets that show the primary flows. See the tests in the repository for more realistic examples.

Create the handler and connect:

```go
cred := ctraderopenapi.ApplicationCredentials{ClientId: "id", ClientSecret: "secret"}
h, err := ctraderopenapi.NewApiHandler(cred, ctraderopenapi.Environment_Demo)
if err != nil {
    // ...
}

// Optionally adjust runtime configuration before connecting
// cfg := ctraderopenapi.APIHandlerConfig{
//     /* ... */
// }
// h = h.WithConfig(cfg)

if err := h.Connect(); err != nil {
    // ...
}
defer h.Disconnect()
```

Sending a request (RPC-style):

```go
req := messages.ProtoOAVersionReq{ /* fill fields if needed */ }
var res messages.ProtoOAVersionRes
reqData := ctraderopenapi.RequestData{
    Ctx:     context.Background(),
    ReqType: PROTO_OA_VERSION_REQ,
    Req:     &req,
    ResType: PROTO_OA_VERSION_RES,
    Res:     &res,
}
if err := h.SendRequest(reqData); err != nil {
    // ...
}
// use res
```

Listening to streaming events

The API exposes two related concepts for events:

- Subscribable events: events that require a server-side subscription (for example market spots or depth events). Use `SubscribeEvent`/`UnsubscribeEvent` to request or cancel server subscriptions.
- Listenable events: events that are delivered by the server and can be listened to using `ListenToEvent`.

`ListenToEvent` accepts a generic callback signature that receives a `ListenableEvent` (a small marker interface implemented by the generated protobuf pointer types). To make callers' life easier you can use the provided typed adapter so you pass a typed function like `func(*messages.ProtoOASpotEvent)` instead of writing the wrapper yourself.

Example using the provided channel-based listener adapter:

```go
wg := sync.WaitGroup{}
wg.Add(1)

onSpot := func(e *messages.ProtoOASpotEvent) {
    defer wg.Done()
    // handle spot event
}

// Create a channel that will receive generic ListenableEvent values.
onEventCh := make(chan datatypes.ListenableEvent)

// Register the listener — ListenToEvent will close `onEventCh` when the
// listener's context is canceled or when the handler disconnects.
if err := h.ListenToEvent(EventType_Spots, onEventCh, ctx); err != nil {
    // handle error
}

// SpawnEventHandler starts a small goroutine that converts generic events
// received on `onEventCh` to the concrete typed event and calls `onSpot`.
if err := SpawnEventHandler(ctx, onEventCh, onSpot); err != nil {
    // handle error
}

// subscribe so the server will send immediate event messages
sub := SubscribableEventData{
    EventType: EventType_Spots,
    SubcriptionData: &SubscriptionDataSpotEvent{
        CtraderAccountId: CtraderAccountId(accountId),
        SymbolIds:        []int64{symbolId},
    },
}
if err := h.SubscribeEvent(sub); err != nil {
    // ...
}

wg.Wait()

// cleanup
cancel()
```

Important notes about adapters and type-safety

- The adapter uses a runtime type assertion (e.(T)). That means if you pass a typed handler for the wrong `eventType` the assertion will panic. The library intentionally avoids using reflection to keep the runtime surface small and explicit; callers are responsible for matching the handler type to the `eventType` they listen for.
- The generated protobuf pointer types in `messages/` implement small marker methods (see `messages/event_assert.go`) so only intended message types satisfy the event interfaces. This prevents accidental registration of unrelated message types at compile-time for simple cases.

Running the repository tests

Some tests in this repo are integration-like and expect real credentials. The test helpers call `godotenv.Load()` which will load a `.env` file from the current working directory if present; alternatively you can set the required environment variables directly in your shell. The `.env` file is optional and not required for normal usage of the library — it is only convenient for running the package's internal tests locally.

Create a `.env` file in the repository root (or set these variables in your shell) with the following keys if you want to run the tests that require a live account:

- `_TEST_OPENAPI_CLIENT_ID` — application client id used by the tests
- `_TEST_OPENAPI_CLIENT_SECRET` — application client secret used by the tests
- `_TEST_OPENAPI_ACCOUNT_ID` — the numeric cTrader account id used for tests
- `_TEST_OPENAPI_ACCOUNT_ACCESS_TOKEN` — a valid account access token
- `_TEST_OPENAPI_ACCOUNT_ENVIRONMENT` — `DEMO` or `LIVE`

Example on Windows PowerShell (alternative to a `.env` file):

```powershell
$env:_TEST_OPENAPI_ACCOUNT_ID = "123456"
$env:_TEST_OPENAPI_ACCOUNT_ACCESS_TOKEN = "<token>"
$env:_TEST_OPENAPI_ACCOUNT_ENVIRONMENT = "DEMO"
go test ./...
```

Notes and caveats

- The library uses generated protobuf types; do not edit the generated `.pb.go` files by hand.
- Error types returned from server responses are modelled and wrapped into client errors (see `errors.go`). If you see garbled error text it usually means the library attempted to print raw payload bytes instead of unmarshalling them into the expected protobuf error message — this repository already contains code paths that unmarshal server error payloads into `messages.ProtoErrorRes`.
- The project includes small helpers to adapt typed listeners to the generic listener signature — there are two equivalently functioning adapters in the codebase; consider unifying their names if you prefer a single canonical API.

Contributing

Contributions are welcome. Good first issues include:

- Add an `examples/` directory with a minimal bot that authenticates and subscribes to spots.
- Add unit tests that mock `tcp` and exercise the request mapper and event handler without network access.
- Consolidate adapter helper names and improve the documentation examples.

## License

This project is licensed under the Apache License 2.0 — see the [LICENSE](./LICENSE) file for details.

---

If you'd like, I can also add a small `examples/` directory and a focused unit test that demonstrates `ListenToEvent` + `SubscribeEvent` without needing a networked cTrader account — say the word and I'll scaffold it.