package mock

import (
	"context"
	"fmt"

	"github.com/kakkky/mcp-sdk-go/shared/schema"
	"github.com/kakkky/mcp-sdk-go/shared/schema/jsonrpc"
)

type MockChannelClientTransport struct {
	clientToServerCh chan []byte
	serverToClientCh chan []byte
	cancel           context.CancelFunc // Close時にgoroutineを終了するためのcancel関数
	onReceiveMessage func(schema.JsonRpcMessage)
	onClose          func()
	onError          func(error)
}

func NewMockChannelClientTransport(clientToServerCh chan []byte, serverToClientChan chan []byte) *MockChannelClientTransport {
	return &MockChannelClientTransport{
		clientToServerCh: clientToServerCh,
		serverToClientCh: serverToClientChan,
	}
}

func (m *MockChannelClientTransport) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case message, ok := <-m.serverToClientCh:
				if !ok {
					return
				}
				msg, err := jsonrpc.Unmarshal(message)
				if err != nil {
					m.onError(err)
					continue
				}
				m.onReceiveMessage(msg)
			}
		}
	}()
	return nil
}

func (m *MockChannelClientTransport) Close() error {
	if m.onClose != nil {
		m.onClose()
	}
	if m.cancel == nil {
		return fmt.Errorf("cancel is nil")
	}
	m.cancel()
	return nil
}

func (m *MockChannelClientTransport) SendMessage(message schema.JsonRpcMessage) error {
	msg, err := jsonrpc.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling message:", err)
		return err
	}
	m.clientToServerCh <- msg
	return nil
}

func (m *MockChannelClientTransport) OnReceiveMessage(message schema.JsonRpcMessage) {
	m.onReceiveMessage(message)
}

func (m *MockChannelClientTransport) OnClose() {
	m.onClose()
}

func (m *MockChannelClientTransport) OnError(err error) {
	m.onError(err)
}

func (m *MockChannelClientTransport) SetOnClose(onClose func()) {
	m.onClose = onClose
}
func (m *MockChannelClientTransport) SetOnError(onError func(error)) {
	m.onError = onError
}

func (m *MockChannelClientTransport) SetOnReceiveMessage(onReceiveMessage func(schema.JsonRpcMessage)) {
	m.onReceiveMessage = onReceiveMessage
}
