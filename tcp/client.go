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
	"github.com/linuskuehnle/ctraderopenapi/datatypes"

	"bufio"
	"context"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type TCPClient interface {
	WithTimeout(timeout time.Duration) TCPClient
	WithTLS(config *tls.Config) TCPClient
	WithMessageCallbackChan(onMessageCh chan []byte) TCPClient
	WithReconnectTimeout(timeout time.Duration) TCPClient
	WithMaxReconnectAttempts(maxAttempts int) TCPClient
	WithInfiniteReconnectAttempts() TCPClient
	WithConnEventHooks(onConnectionLoss func(), onReconnectSuccess func(), onReconnectFail func(error)) TCPClient
	WithRetryBackoff(backoff datatypes.RetryBackoff) TCPClient

	GetTimeout() time.Duration

	HasConn() bool
	OpenConn() error
	CloseConn() error
	CleanupConn()
	JustCloseConn()

	Send(data []byte) error
	Read() ([]byte, error)
}

type messageHandling struct {
	stoppedCh   chan struct{}
	onMessageCh chan []byte
}

type tcpClient struct {
	mu sync.RWMutex

	address string

	conn   net.Conn
	reader *bufio.Reader

	connCtx       context.Context
	connCtxCancel context.CancelFunc

	useTLS    bool
	tlsConfig *tls.Config

	timeout time.Duration

	reconnectBackoff datatypes.RetryBackoff
	reconnectTimeout time.Duration
	maxReconnects    int

	onConnectionLoss   func()
	onReconnectSuccess func()
	onReconnectFail    func(error)

	messageHandling *messageHandling

	onFatalErr func(error)
}

func NewTCPClient(address string, onFatalErr func(error)) TCPClient {
	return newTCPClient(address, onFatalErr)
}

func newTCPClient(address string, onFatalErr func(error)) *tcpClient {
	if onFatalErr == nil {
		onFatalErr = func(err error) {}
	}

	return &tcpClient{
		mu:               sync.RWMutex{},
		address:          address,
		timeout:          DefaultTimeout,
		reconnectTimeout: DefaultReconnectTimeout,
		maxReconnects:    DefaultMaxReconnects,
		onFatalErr:       onFatalErr,
	}
}

func (c *tcpClient) WithReconnectTimeout(timeout time.Duration) TCPClient {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return c
	}

	c.reconnectTimeout = timeout
	return c
}

func (c *tcpClient) WithMaxReconnectAttempts(maxAttempts int) TCPClient {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil || maxAttempts <= 0 {
		return c
	}

	c.maxReconnects = maxAttempts
	return c
}

func (c *tcpClient) WithInfiniteReconnectAttempts() TCPClient {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return c
	}

	c.maxReconnects = 0
	return c
}

func (c *tcpClient) WithConnEventHooks(onConnectionLoss func(), onReconnectSuccess func(), onReconnectFail func(error)) TCPClient {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return c
	}

	c.onConnectionLoss = onConnectionLoss
	c.onReconnectSuccess = onReconnectSuccess
	c.onReconnectFail = onReconnectFail
	return c
}

func (c *tcpClient) WithRetryBackoff(backoff datatypes.RetryBackoff) TCPClient {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return c
	}

	c.reconnectBackoff = backoff
	return c
}

func (c *tcpClient) WithTimeout(timeout time.Duration) TCPClient {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil { // Disallow changing timeout while connection is open
		return c
	}

	c.timeout = timeout
	return c
}

func (c *tcpClient) WithTLS(config *tls.Config) TCPClient {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil { // Disallow changing timeout while connection is open
		return c
	}

	c.useTLS = true
	c.tlsConfig = config
	return c
}

func (c *tcpClient) WithMessageCallbackChan(onMessageCh chan []byte) TCPClient {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil { // Disallow changing timeout while connection is open
		return c
	}
	if onMessageCh == nil {
		return c
	}

	c.messageHandling = &messageHandling{
		stoppedCh:   make(chan struct{}),
		onMessageCh: onMessageCh,
	}
	return c
}

func (c *tcpClient) GetTimeout() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()

	timeout := c.timeout
	return timeout
}

func (c *tcpClient) HasConn() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	hasConn := c.conn != nil
	return hasConn
}

func (c *tcpClient) OpenConn() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return &OpenConnectionError{
			ErrorText: "connection is already open",
		}
	}
	host, _, err := splitAddress(c.address)
	if err != nil {
		return &datatypes.FunctionInvalidArgError{
			FunctionName: "OpenConn",
			Err:          fmt.Errorf("invalid address '%s': %w", c.address, err),
		}
	}

	if err := c.establishConn(host); err != nil {
		return err
	}

	c.connCtx, c.connCtxCancel = context.WithCancel(context.Background())

	if c.messageHandling != nil {
		// Start a goroutine to handle incoming messages
		go c.handleInputStream(c.connCtx)
	}

	return nil
}

func (c *tcpClient) establishConn(host string) error {
	if host == "" {
		host, _, _ = splitAddress(c.address)
	}

	// Check if TLS is enabled and config is provided
	if c.useTLS && c.tlsConfig != nil {
		// Update tls config sever name to match address
		if net.ParseIP(host) == nil {
			// Only update server name if it is a hostname
			c.tlsConfig.ServerName = host
		}
	}

	// TCP dial
	rawConn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}

	if tcpConn, ok := rawConn.(*net.TCPConn); ok {
		// Enable TCP keepalive to detect dead connections early without overhead.
		// The OS-level keepalive probes will detect lost connections immediately.
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(time.Second * 3) // Check periodically
	}

	if !c.useTLS {
		c.conn = rawConn
		c.reader = bufio.NewReader(c.conn)

		return nil
	}

	// Wrap with TLS
	tlsConn := tls.Client(rawConn, c.tlsConfig)
	if err = tlsConn.Handshake(); err != nil {
		rawConn.Close()
		return err
	}

	c.conn = tlsConn
	c.reader = bufio.NewReader(c.conn)

	return nil
}

func (c *tcpClient) CloseConn() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return &CloseConnectionError{
			ErrorText: "connection is not open",
		}
	}

	c.cleanupMessageHandling()

	if err := c.conn.Close(); err != nil {
		return &CloseConnectionError{
			ErrorText: err.Error(),
		}
	}

	c.clear()

	return nil
}

func (c *tcpClient) CleanupConn() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cleanupMessageHandling()

	if c.conn != nil {
		c.conn.Close()
	}

	c.clear()
}

func (c *tcpClient) JustCloseConn() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *tcpClient) cleanupMessageHandling() {
	if c.connCtx != nil {
		c.connCtxCancel()
	}

	if c.messageHandling != nil {
		// Wait for message handling to stop
		<-c.messageHandling.stoppedCh
	}

	if c.messageHandling == nil {
		return
	}

	c.messageHandling = nil
}

func (c *tcpClient) clear() {
	c.conn = nil
	c.reader = nil

	c.connCtx = nil
	c.connCtxCancel = nil
}

func (c *tcpClient) Send(bytes []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if bytes == nil {
		return &datatypes.FunctionInvalidArgError{
			FunctionName: "Send",
			Err:          errors.New("bytes mustn't be nil"),
		}
	}

	if c.conn == nil {
		return &NoConnectionError{
			CallContext: "Send",
		}
	}

	// create header including 4 byte length
	header := make([]byte, 4)
	bytesLen := len(bytes)
	binary.BigEndian.PutUint32(header, uint32(bytesLen))

	// Concatenate header with bytes
	payloadBytes := append(header, bytes...)

	if err := c.conn.SetWriteDeadline(time.Now().Add(c.timeout)); err != nil {
		return err
	}

	n, err := c.conn.Write(payloadBytes)
	if err != nil {
		return err
	}

	if n != len(payloadBytes) {
		return fmt.Errorf("incomplete write: expected %d, wrote %d", len(bytes), n)
	}

	return nil
}

func (c *tcpClient) Read() ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.messageHandling != nil {
		// If a message handling is set, Read mustn't be used directly
		// This is to prevent confusion between explicit reads and handler-based reads
		return nil, &OperationBlockedError{
			CallContext: "Read",
			ErrorText:   "message handler function is set, cannot read explicitly",
		}
	}

	return c.read()
}

func (c *tcpClient) read() ([]byte, error) {
	if c.conn == nil {
		return nil, &NoConnectionError{
			CallContext: "Read",
		}
	}

	deadline := time.Now().Add(c.timeout)
	if err := c.conn.SetReadDeadline(deadline); err != nil {
		return nil, err
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), c.timeout)
	defer cancelCtx()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("waiting for message timed out")
		default:
			// Peek header
			if _, err := c.reader.Peek(4); err != nil {
				if err == io.EOF {
					continue
				}
				return nil, fmt.Errorf("error peeking header: %w", err)
			}
		}

		return c.readPayload()
	}
}

func (c *tcpClient) waitForHeader(ctx context.Context) (bool, error) {
	errCh := make(chan error, 1)

	go func(errCh chan error) {
		defer close(errCh)

		_, err := c.reader.Peek(4)
		errCh <- err
	}(errCh)

	select {
	case <-ctx.Done():
		return false, nil
	case err := <-errCh:
		return err == nil, err
	}
}

func (c *tcpClient) readPayload() ([]byte, error) {
	// Read header
	headerLen := 4
	header := make([]byte, headerLen)
	if _, err := io.ReadFull(c.reader, header); err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Extract payload length and type
	payloadLen := int(binary.BigEndian.Uint32(header))
	if payloadLen < 0 {
		return nil, fmt.Errorf("invalid payload length: %d", payloadLen)
	}

	// Read payload
	payloadBytes := make([]byte, payloadLen)
	if _, err := io.ReadFull(c.reader, payloadBytes); err != nil {
		return nil, fmt.Errorf("failed to read payload: %w", err)
	}
	return payloadBytes, nil
}

func (c *tcpClient) readNextDataBlock(ctx context.Context) ([]byte, error) {
	// Wait for the message header bytes
	if hasHeader, err := c.waitForHeader(ctx); err != nil || !hasHeader {
		return nil, err
	}

	// Set read deadline
	deadline := time.Now().Add(c.timeout)
	if err := c.conn.SetReadDeadline(deadline); err != nil {
		return nil, err
	}

	payload, err := c.readPayload()
	if err != nil {
		return nil, err
	}

	// Reset read deadline
	if err := c.conn.SetReadDeadline(time.Time{}); err != nil {
		return nil, err
	}

	return payload, nil
}

func (c *tcpClient) handleInputStream(ctx context.Context) {
	defer func() {
		c.messageHandling.stoppedCh <- struct{}{}
		close(c.messageHandling.stoppedCh)
	}()

	for {
		select {
		case <-ctx.Done():
			// Context canceled — connection is being closed gracefully
			return
		default:
			payload, err := c.readNextDataBlock(ctx)
			if err != nil {
				if !isNetworkFatal(err) {
					c.onFatalErr(err)
				}

				// Fatal network connection error — connection is lost
				c.reconnectBackoff.Backoff()

				reconnectCh := make(chan struct{})
				go c.execReconnectLoop(ctx, reconnectCh)

				select {
				case <-ctx.Done():
					// Context canceled — connection is being closed gracefully
					return
				case _, ok := <-reconnectCh:
					if !ok {
						// Reconnect failed or context has been cancelled. No reconnect signal has been emitted.
						return
					}
				}

				// Reconnect successful
				c.handleReconnectSuccess()
			}

			if payload != nil && c.messageHandling != nil {
				c.messageHandling.onMessageCh <- payload
			}
		}
	}
}

func (c *tcpClient) execReconnectLoop(ctx context.Context, reconnectedCh chan struct{}) {
	defer close(reconnectedCh) // Ensure channel is always closed on all exit paths

	permitted := c.reconnectBackoff.WaitForPermit(ctx)
	if !permitted {
		// Context canceled during backoff wait - clean shutdown
		return
	}

	c.mu.Lock()

	c.conn = nil
	c.reader = nil

	reconnectAttempt := 0
	maxAttempts := c.maxReconnects
	timeout := c.reconnectTimeout
	c.mu.Unlock()

	if c.onConnectionLoss != nil {
		c.onConnectionLoss()
	}

	timer := time.NewTimer(0) // Immediate first attempt
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			// Context canceled during reconnect - clean shutdown
			return
		case <-timer.C:
			if maxAttempts > 0 {
				if reconnectAttempt >= maxAttempts {
					go c.onFatalErr(&MaxReconnectAttemptsReachedError{
						MaxAttempts: maxAttempts,
					})
					return
				}

				reconnectAttempt++
			}

			c.mu.Lock()
			err := c.establishConn("")
			if err == nil {
				c.mu.Unlock()

				reconnectedCh <- struct{}{} // Signal successful reconnect before closing
				return
			}
			c.mu.Unlock()

			c.handleReconnectFail(err)

			timer.Reset(timeout)
		}
	}
}

func (c *tcpClient) handleReconnectFail(err error) {
	if c.onReconnectFail != nil {
		c.onReconnectFail(err)
	} else {
		go c.onFatalErr(err)
	}
}

func (c *tcpClient) handleReconnectSuccess() {
	// Reset backoff on successful reconnect
	if c.reconnectBackoff != nil {
		c.reconnectBackoff.Reset()
	}

	if c.onReconnectSuccess != nil {
		c.onReconnectSuccess()
	}
}
