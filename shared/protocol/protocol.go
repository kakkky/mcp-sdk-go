package protocol

import (
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

type Protocol struct {
	requestHandlers      map[string]requestHandler
	notificationHandlers map[string]notificationHandler
	responseHandlers     map[int]responseHandler
}

func (p *Protocol) SetRequestHandler(requestSchema schema.Request, handler func(request schema.JsonRpcRequest)) {
	method := requestSchema.Method
	// TODO: ここで、指定されたmethodをすでに登録していないか確認
	p.requestHandlers[method] = handler
}

func (p *Protocol) SetNotificationHandler(notificationSchema schema.Notification, handler func(notification schema.JsonRpcNotification)) {
	method := notificationSchema.Method
	p.notificationHandlers[method] = handler
}

// リクエスト送信の際に、対応するレスポンスハンドラを登録する
func (p *Protocol) SetResponseHandler(messageId int, handler func(response schema.JsonRpcResponse)) {
	p.responseHandlers[messageId] = handler
}
