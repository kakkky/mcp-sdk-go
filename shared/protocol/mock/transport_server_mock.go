package mock

import (
	"context"
	"fmt"

	"github.com/kakkky/mcp-sdk-go/shared/schema"
	"github.com/kakkky/mcp-sdk-go/shared/schema/jsonrpc"
)

type MockChannelServerTransport struct {
	clientToServerCh chan []byte
	serverToClientCh chan []byte
	cancel           context.CancelFunc // Close時にgoroutineを終了するためのcancel関数
	onReceiveMessage func(schema.JsonRpcMessage)
	onClose          func()
	onError          func(error)
}

func NewMockChannelServerTransport(clientToServerCh chan []byte, serverToClientCh chan []byte) *MockChannelServerTransport {
	return &MockChannelServerTransport{
		clientToServerCh: clientToServerCh,
		serverToClientCh: serverToClientCh,
	}
}

func (m *MockChannelServerTransport) Start() error {
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

func (m *MockChannelServerTransport) Close() error {
	if m.onClose != nil {
		m.onClose()
	}
	if m.cancel == nil {
		return fmt.Errorf("cancel is nil")
	}
	m.cancel()
	return nil
}

func (m *MockChannelServerTransport) SendMessage(message schema.JsonRpcMessage) error {
	msg, err := jsonrpc.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling message:", err)
		return err
	}
	m.serverToClientCh <- msg
	return nil
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
