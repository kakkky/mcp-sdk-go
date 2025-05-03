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
