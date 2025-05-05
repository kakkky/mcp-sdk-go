package test

import "github.com/kakkky/mcp-sdk-go/shared/schema"

type MockChannelTransport struct {
	communicationCh  chan schema.JsonRpcMessage
	onReceiveMessage func(schema.JsonRpcMessage)
	onClose          func()
	onError          func(error)
}

func NewMockChannelTransport(communicateCh chan schema.JsonRpcMessage) *MockChannelTransport {
	return &MockChannelTransport{
		communicationCh: communicateCh,
	}
}

func (m *MockChannelTransport) Start() {
	localCh := m.communicationCh
	localOnReceive := m.onReceiveMessage

	if localOnReceive == nil {
		return
	}

	go func() {
		for v := range localCh {
			if localOnReceive != nil {
				localOnReceive(v)
			}
		}
	}()
}

func (m *MockChannelTransport) Close() {}

func (m *MockChannelTransport) SendMessage(message schema.JsonRpcMessage) error {
	m.communicationCh <- message
	return nil
}

func (m *MockChannelTransport) OnReceiveMessage(message schema.JsonRpcMessage) {
	m.onReceiveMessage(message)
}

func (m *MockChannelTransport) OnClose() {
	m.onClose()
}

func (m *MockChannelTransport) OnError(err error) {
	m.onError(err)
}

func (m *MockChannelTransport) SetOnClose(onClose func()) {
	m.onClose = onClose
}
func (m *MockChannelTransport) SetOnError(onError func(error)) {
	m.onError = onError
}

func (m *MockChannelTransport) SetOnReceiveMessage(onReceiveMessage func(schema.JsonRpcMessage)) {
	m.onReceiveMessage = onReceiveMessage
}
