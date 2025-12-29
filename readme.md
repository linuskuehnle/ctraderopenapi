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
	
	**Connection management:**
	- `Connect() error` / `Disconnect()` — establish or close the connection
	
	**Account authentication:**
	- `AuthenticateAccount(CtraderAccountId, AccessToken) (*ProtoOAAccountAuthRes, error)`
		— authenticate with a specific cTrader account using an access token.
		Must be called before making account-specific requests or subscribing to account events.
	- `LogoutAccount(CtraderAccountId, bool) (*ProtoOAAccountLogoutRes, error)` — logout from
		a cTrader account. The boolean parameter controls whether to wait for the server's
		`ProtoOAAccountDisconnectEvent` confirmation before returning.
	- `RefreshAccessToken(AccessToken, RefreshToken) (*ProtoOARefreshTokenRes, error)`
		— refresh an expired access token using a refresh token, returns a new access token.
	
	**Request/Response:**
	- `SendRequest(RequestData) error` — sends a protobuf-typed request and
		unmarshals the response into the provided response object.
	
	**Event subscriptions:**
	- `SubscribeAPIEvent(APIEventData)` / `UnsubscribeAPIEvent(...)` —
		subscribe/unsubscribe for server-side subscription-based events.
	- `SubscribeClientEvent(SubscribableClientEventData)` / `UnsubscribeClientEvent(...)` —
		subscribe/unsubscribe for client-side events.
	
	**Event listening:**
	- `ListenToAPIEvent(ctx, eventType, chan APIEvent)` — register a
		long-running listener channel for push events.
	- `ListenToClientEvent(ctx, clientEventType, chan ClientEvent)`
		— listen for client-side events (connection loss, reconnect events, fatal errors).
	
	**Client event types:** Fatal client errors, connection loss, reconnect success, and reconnect fail.
	
	**Error handling:** Fatal (non-recoverable) client errors are emitted as `FatalErrorEvent`.
	If a listener channel is registered, the error is sent to the channel and the client
	recovers by reconnecting. Without a listener, a fatal error raises a panic.

	**Reconnection:** The API Client automatically handles connection losses. Register listeners
	for `ConnectionLossEvent`, `ReconnectSuccessEvent`, and `ReconnectFailEvent` using
	`ListenToClientEvent`. Previously subscribed API Events are automatically resubscribed
	before `ReconnectSuccessEvent` is emitted.

	Runtime configuration:
	- `WithQueueBufferSize(int)` updates the number of queued requests that may be buffered
		by the internal request queue before backpressure applies.
	- `WithTCPMessageBufferSize(int)` updates the size of the channel used to receive inbound
		TCP messages from the network reader.
	- `WithRequestHeapIterationTimeout(time.Duration)` updates interval used by the request heap
		to periodically check for expired request contexts.
	- `WithRequestTimeout(time.Duration)` updates the duration until a request roundtrip is aborted
		no matter if already sent or not.
	- `DisableDefaultRateLimiter()` to disable the client-side rate limiter.

	**Rate limiting:** The client enforces rate limits to prevent server-side throttling:
	- 50 live requests per second (ProtoOALiveView subscriptions)
	- 5 historical requests per second (ProtoOAHistoricalData requests)
	
	If the server-side rate limit is still exceeded due to network divergence, the request
	is automatically re-enqueued, so callers do not need to handle rate limit errors.

- `ApplicationCredentials{ ClientId, ClientSecret }` — credentials used by the application to
	authenticate with the OpenAPI. Validate with `CheckError()`.

- `CtraderAccountId` — thin typed alias over int64 for explicit account ID references.

- `AccessToken` — thin typed alias over string for access tokens obtained after authenticating an account.

- `RefreshToken` — thin typed alias over string for refresh tokens used to obtain new access tokens.

- Event helpers and adapters:
	- `APIEvent` and `ClientEvent` — marker interfaces for event types.
	- `CastToEventType[T]` and `CastToClientEventType[T]` — helpers to cast generic event types
		to concrete event types.
	- `SpawnAPIEventHandler` and `SpawnClientEventHandler` — start small goroutines that forward
		typed events to your handler, eliminating the need for manual type assertions.

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

2) Account authentication and logout

```go
package main

import (
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

	// Authenticate with a cTrader account
	accountId := ctraderopenapi.CtraderAccountId(123456)
	accessToken := ctraderopenapi.AccessToken("your-access-token")
	
	authRes, err := client.AuthenticateAccount(accountId, accessToken)
	if err != nil {
		fmt.Println("authentication failed:", err)
		return
	}

	fmt.Println("authenticated successfully:", authRes)

	// Use the authenticated account to make requests or subscribe to events...

	// Logout from the account (waitForConfirm=true waits for server confirmation)
	if _, err := client.LogoutAccount(accountId, true); err != nil {
		fmt.Println("logout failed:", err)
		return
	}

	fmt.Println("logged out successfully")
}
```

3) Token refresh

```go
package main

import (
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

	// Refresh an expired access token
	expiredToken := ctraderopenapi.AccessToken("expired-token")
	refreshToken := ctraderopenapi.RefreshToken("refresh-token")
	
	refreshRes, err := client.RefreshAccessToken(expiredToken, refreshToken)
	if err != nil {
		fmt.Println("token refresh failed:", err)
		return
	}

	// Use the new access token for subsequent requests
	newAccessToken := ctraderopenapi.AccessToken(refreshRes.GetAccessToken())
	fmt.Println("new access token:", newAccessToken)
}
```

4) Subscribe to spot events and handle them with the typed adapter

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

onEventCh := make(chan ctraderopenapi.APIEvent)
if err := client.ListenToAPIEvent(ctx, ctraderopenapi.APIEventType_Spots, onEventCh); err != nil {
	panic(err)
}

_ = ctraderopenapi.SpawnAPIEventHandler(ctx, onEventCh, func(e *ctraderopenapi.ProtoOASpotEvent) {
	fmt.Println("spot event:", e)
})

sub := ctraderopenapi.APIEventData{
	EventType: ctraderopenapi.EventType_Spots,
	SubcriptionData: &ctraderopenapi.SpotEventData{
		CtraderAccountId: ctraderopenapi.CtraderAccountId(123456),
		SymbolIds:        []int64{1, 2, 3},
	},
}
if err := client.SubscribeAPIEvent(sub); err != nil {
	panic(err)
}
```

5) Listen for client events; e.g. ReconnectSuccessEvent:

```go

clientCh := make(chan ctraderopenapi.ClientEvent)
if err := client.ListenToClientEvent(context.Background(), ctraderopenapi.ClientEventType_ReconnectSuccessEvent, clientCh); err != nil {
	panic(err)
}
ctraderopenapi.SpawnClientEventHandler(context.Background(), clientCh, func(e *ctraderopenapi.ReconnectSuccessEvent) {
	fmt.Println("reconnected")
})
```

I suggest registering all event listener channels before calling `apiClient.Connect()`.
`apiClient.Disconnect()` unregisters all event listener channels; if you want to
reconnect and use the same listeners, explicitly register them again before the next `Connect()` call.

**Authentication workflow:**
1. Create and connect the API client with application credentials
2. Authenticate with a specific account using `AuthenticateAccount(accountId, accessToken)`
3. Make account-specific requests or subscribe to events
4. Optionally refresh access tokens with `RefreshAccessToken()` if tokens expire
	(see APIEvent `ProtoOAAccountsTokenInvalidatedEvent`)
5. Logout with `LogoutAccount()` when finished with the account

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

