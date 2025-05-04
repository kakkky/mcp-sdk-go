package protocol

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/kakkky/mcp-sdk-go/shared/mcp_err"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

type Protocol struct {
	transport        transport
	handlers         handlers
	onClose          func()
	onError          func(error)
	requestMessageId int

	resultCh chan *schema.Result
	errCh    chan error
}

func NewProtocol() *Protocol {
	p := &Protocol{
		handlers: handlers{
			requestHandlers:      make(map[string]requestHandler),
			notificationHandlers: make(map[string]notificationHandler),
			responseHandlers:     make(map[int]responseHandler),
		},
		requestMessageId: 0,
		resultCh:         make(chan *schema.Result, 1),
		errCh:            make(chan error, 1),
	}
	p.onClose = func() {
		responseHandlers := p.handlers.responseHandlers
		for _, handler := range responseHandlers {
			handler(nil, mcp_err.NewMcpErr(mcp_err.CONNECTION_CLOSED, "connection closed", nil))
		}
		p.handlers.responseHandlers = make(map[int]responseHandler)
		p.transport = nil

	}
	p.onError = func(err error) {}

	return p
}

func (p *Protocol) SetOnClose(onClose func()) {
	p.onClose = func() {
		p.onClose()
		onClose()
	}
}

func (p *Protocol) SetOnError(onError func(error)) {
	p.onError = onError
}

func (p *Protocol) Connect(transport transport) {
	p.transport = transport
	p.transport.setOnClose(p.onClose)
	p.transport.setOnError(p.onError)
	p.transport.setOnMessage(
		func(message schema.JsonRpcMessage) {
			switch m := message.(type) {
			case schema.JsonRpcResponse:
				p.onResponse(m)
			case schema.JsonRpcError:
				p.onErrResponse(m)
			case schema.JsonRpcRequest:
				p.onRequest(m)
			case schema.JsonRpcNotification:
				p.onNotification(m)
			default:
				p.onError(errors.New("unknown message type"))
			}
		})
}

func (p *Protocol) Transport() transport {
	return p.transport
}

func (p *Protocol) Request(request schema.Request, resultSchema any) (*schema.Result, error) {
	if p.transport == nil {
		return nil, fmt.Errorf("not connected")
	}
	messageId := p.requestMessageId + 1
	jsonRpcRequest := schema.JsonRpcRequest{
		BaseMessage: schema.BaseMessage{
			Jsonrpc: schema.JSON_RPC_VERSION,
			Id:      messageId,
		},
		Request: request,
	}
	// リクエストに紐づくレスポンスハンドラを登録する
	p.SetResponseHandler(messageId, func(response *schema.JsonRpcResponse, mcpErr error) (*schema.Result, error) {
		// レスポンスの型をチェック
		result := response.Result
		resultT := reflect.TypeOf(result)
		schemaT := reflect.TypeOf(resultSchema)
		if resultT != schemaT {
			return nil, fmt.Errorf("result type mismatch: expected %s, got %s", schemaT, resultT)
		}
		return &result, nil
	})
	// リクエストの送信
	if err := p.transport.sendMessage(jsonRpcRequest); err != nil {
		return nil, err
	}
	// レスポンスハンドラーからの結果を待つ
	select {
	case result := <-p.resultCh:
		return result, nil
	case err := <-p.errCh:
		return nil, err
	}
}

func (p *Protocol) Notificate(notification schema.Notification) error {
	if p.transport == nil {
		return fmt.Errorf("not connected")
	}
	jsonRpcNotification := schema.JsonRpcNotification{
		Jsonrpc:      schema.JSON_RPC_VERSION,
		Notification: notification,
	}
	if err := p.transport.sendMessage(jsonRpcNotification); err != nil {
		return err
	}
	return nil
}
