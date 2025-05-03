package protocol

import "github.com/kakkky/mcp-sdk-go/shared/schema"

type handlers struct {
	requestHandlers      map[string]requestHandler
	notificationHandlers map[string]notificationHandler
	responseHandlers     map[int]responseHandler
}

type requestHandler = func(request schema.JsonRpcRequest)

type notificationHandler = func(notification schema.JsonRpcNotification)

type responseHandler = func(response schema.JsonRpcResponse)
