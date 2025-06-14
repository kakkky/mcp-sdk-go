package shared

import (
	"github.com/kakkky/mcp-sdk-go/shared/protocol"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

//go:generate mockgen -source=./protocol_interface.go -destination=../client/mock/protocol_mock.go -package=mock
//go:generate mockgen -source=./protocol_interface.go -destination=../mcp-server/server/mock/protocol_mock.go -package=mock
type Protocol interface {
	SetRequestHandler(schema schema.Request, handler func(schema.JsonRpcRequest) (schema.Result, error))
	SetNotificationHandler(schema schema.Notification, handler func(schema.JsonRpcNotification) error)
	ValidateCanSetRequestHandler(method string) error

	SetValidateCapabilityForMethod(validator func(method string) error)
	SetValidateNotificationCapability(validator func(method string) error)
	SetValidateRequestHandlerCapability(validatror func(method string) error)
	Transport() protocol.Transport

	Connect(transport protocol.Transport) error
	Close() error

	Request(request schema.Request, resultSchema any) (schema.Result, error)
	Notificate(notification schema.Notification) error
}
