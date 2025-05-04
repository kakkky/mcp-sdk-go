package protocol

import (
	"github.com/kakkky/mcp-sdk-go/shared/mcp_err"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func (p *Protocol) onResponse(message schema.JsonRpcMessage) {}

func (p *Protocol) onRequest(request schema.JsonRpcRequest) {
	handler := p.handlers.requestHandlers[request.Method]
	if handler == nil && p.handlers.fallbackRequestHandler != nil {
		p.handlers.fallbackRequestHandler()
		return
	}
	if handler == nil {
		err := p.transport.sendMessage(
			schema.JsonRpcError{
				Jsonrpc: schema.JSON_RPC_VERSION,
				Id:      request.Id,
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
		code := err.Code
		err := p.transport.sendMessage(
			schema.JsonRpcError{
				Jsonrpc: schema.JSON_RPC_VERSION,
				Id:      request.Id,
				Error: schema.Error{
					Code:    code,
					Message: err.Error(),
				},
			},
		)
		if err != nil {
			p.onError(err)
		}
	}
	if err := p.transport.sendMessage(schema.JsonRpcResponse{
		Jsonrpc: schema.JSON_RPC_VERSION,
		Id:      request.Id,
		Result:  result,
	}); err != nil {
		p.onError(err)
	}
}

func (p *Protocol) onNotification(notification schema.JsonRpcNotification) {
	handler := p.handlers.notificationHandlers[notification.Method]
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
