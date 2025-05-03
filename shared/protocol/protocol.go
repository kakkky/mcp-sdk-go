package protocol

import (
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

type Protocol struct {
	transport Transport
	handlers  handlers
}

func NewProtocol() *Protocol {
	return &Protocol{
		handlers: handlers{
			requestHandlers:      make(map[string]requestHandler),
			notificationHandlers: make(map[string]notificationHandler),
			responseHandlers:     make(map[int]responseHandler),
		},
	}
}

func (p *Protocol) SetRequestHandler(requestSchema schema.Request, handler func(request schema.JsonRpcRequest)) {
	method := requestSchema.Method
	// TODO: ここで、指定されたmethodをすでに登録していないか確認
	p.handlers.requestHandlers[method] = handler
}

func (p *Protocol) SetNotificationHandler(notificationSchema schema.Notification, handler func(notification schema.JsonRpcNotification)) {
	method := notificationSchema.Method
	p.handlers.notificationHandlers[method] = handler
}

// リクエスト送信の際に、対応するレスポンスハンドラを登録する
func (p *Protocol) SetResponseHandler(messageId int, handler func(response schema.JsonRpcResponse)) {
	p.handlers.responseHandlers[messageId] = handler
}

func (p *Protocol) Connect(transport Transport) {

}
