package protocol

import (
	"fmt"
	"reflect"

	"github.com/kakkky/mcp-sdk-go/shared/mcp_err"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

type Protocol struct {
	transport            transport
	handlers             *handlers
	requestMessageId     int
	onClose              func()
	onError              func(error)
	options              *ProtocolOptions
	capabilityValidators *capabilityValidators

	respCh    chan schema.Result
	errRespCh chan error
}

func NewProtocol(options *ProtocolOptions) *Protocol {
	p := &Protocol{
		handlers: &handlers{
			requestHandlers:      make(map[string]requestHandler),
			notificationHandlers: make(map[string]notificationHandler),
			responseHandlers:     make(map[int]responseHandler),
		},
		requestMessageId: 0,
		options:          options,
		capabilityValidators: &capabilityValidators{
			validateCapabilityForMethod:      nil,
			validateNotificationCapability:   nil,
			validateRequestHandlerCapability: nil,
		},
		respCh:    make(chan schema.Result, 1),
		errRespCh: make(chan error, 1),
	}
	p.onError = func(err error) {}
	p.onClose = func() {
		responseHandlers := p.handlers.responseHandlers
		for _, handler := range responseHandlers {
			handler(nil, mcp_err.NewMcpErr(mcp_err.CONNECTION_CLOSED, "connection closed", nil))
		}
		p.handlers.responseHandlers = make(map[int]responseHandler)
		p.transport = nil
	}

	return p
}

// クローズ時のコールバック処理を追加する
func (p *Protocol) SetOnClose(onClose func()) {
	defaultOnClose := p.onClose
	p.onClose = func() {
		defaultOnClose()
		onClose()
	}
}

func (p *Protocol) SetOnError(onError func(error)) {
	p.onError = onError
}

func (p *Protocol) Connect(transport transport) {
	p.transport = transport
	p.transport.SetOnClose(p.onClose)

	p.transport.SetOnError(p.onError)
	p.transport.SetOnReceiveMessage(p.onReceiveMessage)
	p.transport.Start()
}

func (p *Protocol) Close() {
	if p.transport == nil {
		return
	}
	p.transport.Close()
}

func (p *Protocol) Transport() transport {
	return p.transport
}

func (p *Protocol) Request(request schema.Request, resultSchema any) (schema.Result, error) {
	if p.transport == nil {
		return nil, fmt.Errorf("not connected")
	}

	if p.options != nil && p.options.enforceStrictCapabilities && p.capabilityValidators.validateCapabilityForMethod != nil {
		if err := p.capabilityValidators.validateCapabilityForMethod(request.Method()); err != nil {
			return nil, err
		}

	}
	p.requestMessageId += 1
	messageId := p.requestMessageId
	jsonRpcRequest := schema.JsonRpcRequest{
		BaseMessage: schema.BaseMessage{
			Jsonrpc: schema.JSON_RPC_VERSION,
			Id:      messageId,
		},
		Request: request,
	}
	// リクエストに紐づくレスポンスハンドラを登録する
	p.SetResponseHandler(messageId, func(response *schema.JsonRpcResponse, mcpErr error) (schema.Result, error) {
		// レスポンスの型をチェック
		result := response.Result
		resultT := reflect.TypeOf(result)
		schemaT := reflect.TypeOf(resultSchema)
		if resultT != schemaT {
			return nil, fmt.Errorf("result type mismatch: expected %s, got %s", schemaT, resultT)
		}
		return result, nil
	})
	// リクエストの送信
	if err := p.transport.SendMessage(jsonRpcRequest); err != nil {
		return nil, err
	}
	// 登録したレスポンスハンドラーからの結果を待つ
	select {
	case result := <-p.respCh:
		return result, nil
	case err := <-p.errRespCh:
		return nil, err
	}
}

func (p *Protocol) Notificate(notification schema.Notification) error {
	if p.transport == nil {
		return fmt.Errorf("not connected")
	}
	if p.capabilityValidators.validateNotificationCapability != nil {
		if err := p.capabilityValidators.validateNotificationCapability(notification.Method()); err != nil {
			return err
		}
	}
	jsonRpcNotification := schema.JsonRpcNotification{
		Jsonrpc:      schema.JSON_RPC_VERSION,
		Notification: notification,
	}
	if err := p.transport.SendMessage(jsonRpcNotification); err != nil {
		return err
	}
	return nil
}
