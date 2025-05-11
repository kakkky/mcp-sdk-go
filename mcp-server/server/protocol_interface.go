package server

import (
	"github.com/kakkky/mcp-sdk-go/shared/protocol"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

//go:generate mockgen -source=./protocol_interface.go -destination=./mock/protocol_mock.go -package=mock
type Protocol interface {
	SetRequestHandler(schema schema.Request, handler func(schema.JsonRpcRequest) (schema.Result, error))
	SetNotificationHandler(schema schema.Notification, handler func(schema.JsonRpcNotification) error)

	SetValidateCapabilityForMethod(validator func(method string) error)
	SetValidateNotificationCapability(validator func(method string) error)

	Transport() protocol.Transport

	Request(request schema.Request, resultSchema any) (schema.Result, error)
}
