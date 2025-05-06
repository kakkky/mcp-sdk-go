package test

import (
	"context"

	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

type MockChannelServerTransport struct {
	clientToServerCh chan schema.JsonRpcMessage
	serverToClientCh chan schema.JsonRpcMessage
	cancel           context.CancelFunc // Close時にgoroutineを終了するためのcancel関数
	onReceiveMessage func(schema.JsonRpcMessage)
	onClose          func()
	onError          func(error)
}

func NewMockChannelServerTransport(clientToServerCh chan schema.JsonRpcMessage, serverToClientCh chan schema.JsonRpcMessage) *MockChannelServerTransport {
	return &MockChannelServerTransport{
		clientToServerCh: clientToServerCh,
		serverToClientCh: serverToClientCh,
	}
}

func (m *MockChannelServerTransport) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case message, ok := <-m.clientToServerCh:
				if !ok {
					return
				}
				m.onReceiveMessage(message)
			}
		}
	}()
}

func (m *MockChannelServerTransport) Close() {
	if m.onClose != nil {
		m.onClose()
	}
	m.cancel()
}

func (m *MockChannelServerTransport) SendMessage(message schema.JsonRpcMessage) error {
	m.serverToClientCh <- message
	return nil
}

func (m *MockChannelServerTransport) OnReceiveMessage(message schema.JsonRpcMessage) {
	m.onReceiveMessage(message)
}

func (m *MockChannelServerTransport) OnClose() {
	m.onClose()
}

func (m *MockChannelServerTransport) OnError(err error) {
	m.onError(err)
}

func (m *MockChannelServerTransport) SetOnClose(onClose func()) {
	m.onClose = onClose
}
func (m *MockChannelServerTransport) SetOnError(onError func(error)) {
	m.onError = onError
}

func (m *MockChannelServerTransport) SetOnReceiveMessage(onReceiveMessage func(schema.JsonRpcMessage)) {
	m.onReceiveMessage = onReceiveMessage
}
