package protocol

import (
	"errors"
	"fmt"

	"github.com/kakkky/mcp-sdk-go/shared/mcp_err"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func (p *Protocol) onReceiveMessage(message schema.JsonRpcMessage) {
	switch m := message.(type) {
	case schema.JsonRpcResponse:
		p.onResponse(m)
	case schema.JsonRpcError:
		p.onErrResponse(m)
	case schema.JsonRpcRequest:
		p.onRequest(m)
	case schema.JsonRpcNotification:
		p.onNotification(m)
	default:
		p.onError(errors.New("unknown message type"))
	}
}

func (p *Protocol) onRequest(request schema.JsonRpcRequest) {
	handler := p.handlers.requestHandlers[request.Method()]
	if handler == nil && p.handlers.fallbackRequestHandler != nil {
		p.handlers.fallbackRequestHandler()
		return
	}
	if handler == nil {
		err := p.transport.SendMessage(
			schema.JsonRpcError{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: schema.JSON_RPC_VERSION,
					Id:      request.Id,
				},
				Error: schema.Error{
					Code:    mcp_err.METHOD_NOT_FOUND,
					Message: "method not found",
				},
			},
		)
		if err != nil {
			p.onError(err)
		}
		return
	}
	result, err := handler(request)
	if err != nil {
		// MCPエラー
		if mcpErr, ok := err.(*mcp_err.McpErr); ok {
			code := mcpErr.Code
			if err := p.transport.SendMessage(
				schema.JsonRpcError{
					BaseMessage: schema.BaseMessage{
						Jsonrpc: schema.JSON_RPC_VERSION,
						Id:      request.Id,
					},
					Error: schema.Error{
						Code:    code,
						Message: err.Error(),
					},
				},
			); err != nil {
				p.onError(err)
			}
			return
		}
		// MCPエラーではないエラー
		if err := p.transport.SendMessage(
			schema.JsonRpcError{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: schema.JSON_RPC_VERSION,
					Id:      request.Id,
				},
				Error: schema.Error{
					Code:    mcp_err.INTERNAL_ERROR,
					Message: err.Error(),
				},
			},
		); err != nil {
			p.onError(err)
		}
		return
	}
	if err := p.transport.SendMessage(schema.JsonRpcResponse{
		BaseMessage: schema.BaseMessage{
			Jsonrpc: schema.JSON_RPC_VERSION,
			Id:      request.Id,
		},
		Result: result,
	}); err != nil {
		p.onError(err)
	}
}

func (p *Protocol) onResponse(response schema.JsonRpcResponse) {
	messageId := response.Id
	handler := p.handlers.responseHandlers[messageId]
	if handler == nil {
		err := fmt.Errorf("received a response for an unknown message ID: %d", messageId)
		p.onError(err)
		return
	}
	defer delete(p.handlers.responseHandlers, messageId)
	result, err := handler(&response, nil)
	if err != nil {
		p.onError(err)
		return
	}

	p.respCh <- result

}

func (p *Protocol) onErrResponse(errResponse schema.JsonRpcError) {
	messageId := errResponse.Id
	handler := p.handlers.responseHandlers[messageId]
	if handler == nil {
		p.onError(fmt.Errorf("received a response for an unknown message ID: %d", messageId))
		return
	}
	defer delete(p.handlers.responseHandlers, messageId)
	err := mcp_err.NewMcpErr(errResponse.Error.Code, errResponse.Error.Message, errResponse.Error.Data)

	p.errRespCh <- err

}

func (p *Protocol) onNotification(notification schema.JsonRpcNotification) {
	handler := p.handlers.notificationHandlers[notification.Method()]
	if handler == nil && p.handlers.fallbackNotificationHandler != nil {
		p.handlers.fallbackNotificationHandler()
		return
	}
	if handler == nil {
		return
	}
	if err := handler(notification); err != nil {
		p.onError(err)
	}
}
