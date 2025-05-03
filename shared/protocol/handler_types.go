package protocol

import "github.com/kakkky/mcp-sdk-go/shared/schema"

type requestHandler = func(request schema.JsonRpcRequest)

type notificationHandler = func(notification schema.JsonRpcNotification)

type responseHandler = func(response schema.JsonRpcResponse)
