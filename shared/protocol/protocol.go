package protocol

import (
	"errors"

	"github.com/kakkky/mcp-sdk-go/shared/mcp_err"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

type Protocol struct {
	transport transport
	handlers  handlers
	onClose   func()
	onError   func(error)
}

func NewProtocol() *Protocol {
	p := &Protocol{
		handlers: handlers{
			requestHandlers:      make(map[string]requestHandler),
			notificationHandlers: make(map[string]notificationHandler),
			responseHandlers:     make(map[int]responseHandler),
		},
	}
	p.onClose = func() {
		responseHandlers := p.handlers.responseHandlers
		err := mcp_err.NewMcpErr(mcp_err.CONNECTION_CLOSED, "connection closed")
		for _, handler := range responseHandlers {
			handler(nil, err)
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
	p.transport.setOnMessage(func(message schema.JsonRpcMessage) {
		switch m := message.(type) {
		case schema.JsonRpcResponse, schema.JsonRpcError:
			p.onResponse(m)
		case schema.JsonRpcRequest:
			p.onRequest(m)
		case schema.JsonRpcNotification:
			p.onNotification(m)
		default:
			p.onError(errors.New("unknown message type"))
		}
	})
}
