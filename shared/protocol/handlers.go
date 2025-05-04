package protocol

import (
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

type handlers struct {
	requestHandlers      map[string]requestHandler
	notificationHandlers map[string]notificationHandler
	responseHandlers     map[int]responseHandler

	fallbackNotificationHandler func()
}

type requestHandler = func(request *schema.JsonRpcRequest) error

type notificationHandler = func(notification *schema.JsonRpcNotification) error

type responseHandler = func(response *schema.JsonRpcResponse, mcpErr error) error

func (p *Protocol) SetRequestHandler(requestSchema schema.Request, handler requestHandler) {
	method := requestSchema.Method
	// TODO: ここで、指定されたmethodをすでに登録していないか確認
	p.handlers.requestHandlers[method] = handler
}

func (p *Protocol) SetNotificationHandler(notificationSchema schema.Notification, handler notificationHandler) {
	method := notificationSchema.Method
	p.handlers.notificationHandlers[method] = handler
}

// リクエスト送信の際に、対応するレスポンスハンドラを登録する
func (p *Protocol) SetResponseHandler(messageId int, handler responseHandler) {
	p.handlers.responseHandlers[messageId] = handler
}

func (p *Protocol) SetFallbackNotificationHandler(handler func(), notification schema.Notification) {
	p.handlers.fallbackNotificationHandler = handler
}
