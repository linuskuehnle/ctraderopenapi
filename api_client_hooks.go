package ctraderopenapi

import (
	"context"
	"errors"

	"github.com/linuskuehnle/ctraderopenapi/datatypes"
	"github.com/linuskuehnle/ctraderopenapi/messages"
	"google.golang.org/protobuf/proto"
)

func (c *apiClient) onQueueData() {
	reqMetaData, _ := c.requestQueue.Dequeue()
	// Don't check the error of Dequeue since by the time this
	// callback goroutine is called, all previously enqueued requests might
	// have been context cancelled. Hence we check if theres non-nil returned reqMetaData
	if reqMetaData == nil {
		return
	}

	reqErrCh := reqMetaData.ErrCh

	// Check if the request context has expired already. In that case the request has
	// been cleaned up already and should not be further processed.
	if reqMetaData.Ctx.Err() != nil {
		close(reqErrCh) // For completeness only
		return
	}

	// Check if a response is expected
	expectRes := reqMetaData.ResDataCh != nil
	if !expectRes {
		// If no response is expected, close the channel immediately after sending payload
		defer close(reqErrCh)
	}

	// Check if internal rate limiter is being used
	if c.rateLimiters != nil {
		// Get the rate limit related type
		rateLimitType, ok := rateLimitTypeByReqType[reqMetaData.ReqType]
		if ok {
			// Schedule the request via rate limiter
			c.rateLimiters[rateLimitType].WaitForPermit()
		}
	}

	if err := c.handleSendPayload(reqMetaData); err != nil {
		if reqErrCh != nil {
			reqErrCh <- err

			// This check is required to avoid closing a closed channel when no response is expected
			if expectRes {
				close(reqErrCh)
			}
		}
	}
}

func (c *apiClient) onTCPMessage(msgBytes []byte) {
	if c.lifecycleData.IsClientInitialized() {
		// Only lock if life cycle is running already, on Connect the mutex is already locked
		c.mu.RLock()
		defer c.mu.RUnlock()
	}

	var protoMsg messages.ProtoMessage
	if err := proto.Unmarshal(msgBytes, &protoMsg); err != nil {
		perr := &ProtoUnmarshalError{
			CallContext: "proto message",
			Err:         err,
		}
		c.fatalErrCh <- perr
		return
	}

	msgPayloadType := messages.ProtoPayloadType(protoMsg.GetPayloadType())

	switch msgPayloadType {
	case messages.ProtoPayloadType_ERROR_RES:
		var protoErrorRes messages.ProtoErrorRes
		if err := proto.Unmarshal(msgBytes, &protoErrorRes); err != nil {
			perr := &ProtoUnmarshalError{
				CallContext: "proto error response",
				Err:         err,
			}
			c.fatalErrCh <- perr
			return
		}

		genericResErr := &GenericResponseError{
			ErrorCode:               protoErrorRes.GetErrorCode(),
			Description:             protoErrorRes.GetDescription(),
			MaintenanceEndTimestamp: protoErrorRes.GetMaintenanceEndTimestamp(),
		}

		// Since this response cannot be mapped to any request, make it a fatal error
		c.fatalErrCh <- genericResErr
		return
	case messages.ProtoPayloadType_HEARTBEAT_EVENT:
		// Ignore heartbeat events
		return
	default:
		msgOAPayloadType := ProtoOAPayloadType(msgPayloadType)
		if isListenableEvent[msgOAPayloadType] {
			if err := c.handleListenableAPIEvent(msgOAPayloadType, &protoMsg); err != nil {
				c.fatalErrCh <- err
			}
			return
		}

		if protoMsg.ClientMsgId == nil {
			c.fatalErrCh <- errors.New("invalid proto message on response: field ClientMsgId type is nil")
			return
		}
		reqId := datatypes.RequestId(protoMsg.GetClientMsgId())

		// Remove the request meta data from the request heap
		reqMetaData, err := c.requestHeap.RemoveNode(reqId)
		if err != nil {
			var reqHeapErr *RequestHeapNodeNotIncludedError
			if errors.As(err, &reqHeapErr) {
				// Request context has been cancelled before response has been received
				return
			}
			c.fatalErrCh <- err
			return
		}

		close(reqMetaData.ErrCh)

		reqMetaData.ResDataCh <- &datatypes.ResponseData{
			ProtoMsg:    &protoMsg,
			PayloadType: msgOAPayloadType,
		}
		return
	}
}

// Note: do not call this function, pass errors to apiClient.fatalErrCh
func (c *apiClient) onFatalError(err error) {
	if c.lifecycleData.IsClientInitialized() {
		// Only lock if life cycle is running already, on Connect the mutex is already locked
		c.mu.Lock()
		defer c.mu.Unlock()
	}

	event := &FatalErrorEvent{
		Err: err,
	}
	eventHandled := c.clientEventHandler.HandleEvent(datatypes.EventId(ClientEventType_FatalErrorEvent), event)
	if !eventHandled {
		panic(err)
	}

	// Drop TCP connection to trigger reconnect routine to attempt recovery from the fatal error
	c.tcpClient.JustCloseConn()

	c.lifecycleData.SetClientDisconnected()
}

func (c *apiClient) onConnectionLoss() {
	c.lifecycleData.SetClientDisconnected()

	event := &ConnectionLossEvent{}
	c.clientEventHandler.HandleEvent(datatypes.EventId(ClientEventType_ConnectionLossEvent), event)
}

func (c *apiClient) onReconnectSuccess() {
	jobErrCh := make(chan error, 3)

	go func() {
		if err := c.authenticateApp(); err != nil {
			jobErrCh <- err
			close(jobErrCh)
			return
		}

		ctx, cancelCtx := context.WithCancel(context.Background())
		defer cancelCtx()

		c.accManager.LockModification(ctx)

		if err := c.reloginActiveAccounts(); err != nil {
			jobErrCh <- err
			close(jobErrCh)
			return
		}

		if err := c.resubscribeActiveSubs(); err != nil {
			jobErrCh <- err
			close(jobErrCh)
			return
		}

		close(jobErrCh)
	}()

	go func() {
		jobErr, ok := <-jobErrCh
		if ok {
			// Error occured in authenticateApp, reloginActiveAccounts or resubscribeActiveSubs task.

			// Pass error to reconnect fail routine.
			c.onReconnectFail(jobErr)

			// Force TCP disconnect to start the reconnect process again.
			c.tcpClient.JustCloseConn()

			return
		}

		event := &ReconnectSuccessEvent{}
		c.clientEventHandler.HandleEvent(datatypes.EventId(ClientEventType_ReconnectSuccessEvent), event)

		c.lifecycleData.SetClientConnected()
	}()
}

func (c *apiClient) onReconnectFail(err error) {
	event := &ReconnectFailEvent{
		Err: err,
	}
	c.clientEventHandler.HandleEvent(datatypes.EventId(ClientEventType_ReconnectFailEvent), event)
}
