package ctraderopenapi

import (
	"context"
	"errors"
	"sync"

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

	// Check if a response is expected
	expectRes := reqMetaData.ResDataCh != nil
	if expectRes {
		// Add the request meta data to the request heap
		if err := c.requestHeap.AddNode(reqMetaData); err != nil {
			reqErrCh <- err
			close(reqErrCh)
			return
		}
	} else {
		// If no response is expected, close the error channel immediately after dispatching the request
		defer close(reqErrCh)
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
	if c.lifecycleData.IsClientConnected() {
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
			if err := c.handleListenableEvent(msgOAPayloadType, &protoMsg); err != nil {
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
	if c.lifecycleData.IsClientConnected() {
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

	c.disconnect()
}

func (c *apiClient) onConnectionLoss() {
	event := &ConnectionLossEvent{}
	c.clientEventHandler.HandleEvent(datatypes.EventId(ClientEventType_ConnectionLossEvent), event)
}

func (c *apiClient) onReconnectSuccess() {
	authErrCh := make(chan error)
	resubErrCh := make(chan error)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		if err := c.authenticateApp(); err != nil {
			authErrCh <- err
			close(authErrCh)
			return
		}
		close(authErrCh)

		// Resubscribe current sessions active subscriptions
		wg := sync.WaitGroup{}

		c.mu.RLock()
		defer c.mu.RUnlock()
		for t, s := range c.activeSubscriptions {
			wg.Add(1)

			go func() {
				defer wg.Done()

				select {
				case <-ctx.Done():
					return
				default:
					eventData := SubscribableEventData{
						EventType:       t,
						SubcriptionData: s,
					}
					if err := c.SubscribeEvent(eventData); err != nil {
						resubErrCh <- err
					}
				}
			}()
		}

		wg.Wait()
		close(resubErrCh)
	}()

	go func() {
		authErr := <-authErrCh
		if authErr != nil {
			c.fatalErrCh <- authErr
			return
		}

		resubErr, ok := <-resubErrCh
		cancelCtx()
		if ok {
			c.fatalErrCh <- resubErr
			return
		}

		event := &ReconnectSuccessEvent{}
		c.clientEventHandler.HandleEvent(datatypes.EventId(ClientEventType_ReconnectSuccessEvent), event)
	}()
}

func (c *apiClient) onReconnectFail(err error) {
	event := &ReconnectFailEvent{
		Err: err,
	}
	c.clientEventHandler.HandleEvent(datatypes.EventId(ClientEventType_ReconnectFailEvent), event)
}
