package protocol

import (
	"fmt"

	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

type handlers struct {
	requestHandlers             map[string]requestHandler
	notificationHandlers        map[string]notificationHandler
	responseHandlers            map[int]responseHandler
	fallbackNotificationHandler func()
	fallbackRequestHandler      func()
}

type requestHandler func(request schema.JsonRpcRequest) (schema.Result, error)

type notificationHandler func(notification schema.JsonRpcNotification) error

type responseHandler func(response *schema.JsonRpcResponse, mcpErr error) (schema.Result, error)

func (p *Protocol) SetRequestHandler(requestSchema schema.Request, handler func(request schema.JsonRpcRequest) (schema.Result, error)) {
	method := requestSchema.Method
	if p.capabilityValidators.validateRequestHandlerCapability != nil {
		if err := p.capabilityValidators.validateRequestHandlerCapability(method()); err != nil {
			fmt.Println("Capability validation failed:", err)
			return
		}
	}
	// TODO: ここで、指定されたmethodをすでに登録していないか確認
	p.handlers.requestHandlers[method()] = handler
}

func (p *Protocol) SetNotificationHandler(notificationSchema schema.Notification, handler func(notification schema.JsonRpcNotification) error) {
	method := notificationSchema.Method
	p.handlers.notificationHandlers[method()] = handler
}

// リクエスト送信の際に、対応するレスポンスハンドラを登録する
func (p *Protocol) SetResponseHandler(messageId int, handler func(response *schema.JsonRpcResponse, mcpErr error) (schema.Result, error)) {
	p.handlers.responseHandlers[messageId] = handler
}

func (p *Protocol) SetFallbackNotificationHandler(handler func(), notification schema.Notification) {
	p.handlers.fallbackNotificationHandler = handler
}
func (p *Protocol) SetFallbackRequestHandler(handler func(), request schema.Request) {
	p.handlers.fallbackRequestHandler = handler
}

func (p *Protocol) ValidateCanSetRequestHandler(method string) error {
	if p.handlers.requestHandlers[method] != nil {
		return fmt.Errorf("request handler for method %s already exists , which would be overridden", method)
	}
	return nil
}
