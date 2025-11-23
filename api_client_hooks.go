package ctraderopenapi

import (
	"errors"

	"github.com/linuskuehnle/ctraderopenapi/datatypes"
	"github.com/linuskuehnle/ctraderopenapi/messages"
	"google.golang.org/protobuf/proto"
)

func (h *apiClient) onQueueData() {
	reqMetaData, _ := h.requestQueue.Dequeue()
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
		if err := h.requestHeap.AddNode(reqMetaData); err != nil {
			reqErrCh <- err
			close(reqErrCh)
			return
		}
	} else {
		// If no response is expected, close the error channel immediately after dispatching the request
		defer close(reqErrCh)
	}

	if err := h.handleSendPayload(reqMetaData); err != nil {
		// Check if the queue data callback has been called when the tcp client connection has just been closed
		var noConnErr *NoConnectionError
		if !errors.As(err, &noConnErr) {
			if reqErrCh != nil {
				reqErrCh <- err

				// This check is required to avoid closing a closed channel when no response is expected
				if expectRes {
					close(reqErrCh)
				}
			}
		}
		return
	}
}

func (h *apiClient) onTCPMessage(msgBytes []byte) {
	if h.lifecycleData.IsClientConnected() {
		// Only lock if life cycle is running already, on Connect the mutex is already locked
		h.mu.RLock()
		defer h.mu.RUnlock()
	}

	var protoMsg messages.ProtoMessage
	if err := proto.Unmarshal(msgBytes, &protoMsg); err != nil {
		perr := &ProtoUnmarshalError{
			CallContext: "proto message",
			Err:         err,
		}
		h.fatalErrCh <- perr
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
			h.fatalErrCh <- perr
			return
		}

		genericResErr := &GenericResponseError{
			ErrorCode:               protoErrorRes.GetErrorCode(),
			Description:             protoErrorRes.GetDescription(),
			MaintenanceEndTimestamp: protoErrorRes.GetMaintenanceEndTimestamp(),
		}

		// Since this response cannot be mapped to any request, make it a fatal error
		h.fatalErrCh <- genericResErr
		return
	case messages.ProtoPayloadType_HEARTBEAT_EVENT:
		// Ignore heartbeat events
		return
	default:
		msgOAPayloadType := ProtoOAPayloadType(msgPayloadType)
		if isListenableEvent[msgOAPayloadType] {
			if err := h.handleListenableEvent(msgOAPayloadType, &protoMsg); err != nil {
				h.fatalErrCh <- err
			}
			return
		}

		if protoMsg.ClientMsgId == nil {
			h.fatalErrCh <- errors.New("invalid proto message on response: field ClientMsgId type is nil")
			return
		}
		reqId := datatypes.RequestId(protoMsg.GetClientMsgId())

		// Remove the request meta data from the request heap
		reqMetaData, err := h.requestHeap.RemoveNode(reqId)
		if err != nil {
			var reqHeapErr *RequestHeapNodeNotIncludedError
			if errors.As(err, &reqHeapErr) {
				return
			}
			h.fatalErrCh <- err
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
func (h *apiClient) onFatalError(err error) {
	if h.lifecycleData.IsClientConnected() {
		// Only lock if life cycle is running already, on Connect the mutex is already locked
		h.mu.Lock()
		defer h.mu.Unlock()
	}

	event := &FatalErrorEvent{
		Err: err,
	}
	eventHandled := h.clientEventHandler.HandleEvent(datatypes.EventId(ClientEventType_FatalErrorEvent), event)
	if !eventHandled {
		panic(err)
	}

	h.disconnect()
}

func (h *apiClient) onReconnectAttempt(errCh chan error) {
	defer close(errCh)

	if err := h.authenticateApp(); err != nil {
		errCh <- err
		return
	}

	errCh <- nil
}

func (h *apiClient) onConnectionLoss() {
	event := &ConnectionLossEvent{}
	h.clientEventHandler.HandleEvent(datatypes.EventId(ClientEventType_ConnectionLossEvent), event)
}

func (h *apiClient) onReconnectSuccess() {
	event := &ReconnectSuccessEvent{}
	h.clientEventHandler.HandleEvent(datatypes.EventId(ClientEventType_ReconnectSuccessEvent), event)
}

func (h *apiClient) onReconnectFail(err error) {
	event := &ReconnectFailEvent{
		Err: err,
	}
	h.clientEventHandler.HandleEvent(datatypes.EventId(ClientEventType_ReconnectFailEvent), event)
}
