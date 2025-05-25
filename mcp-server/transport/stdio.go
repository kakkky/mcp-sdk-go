package transport

import (
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

type stdioServerTransport struct {
	onReceiveMessage func(schema.JsonRpcMessage)
	onClose          func()
	onError          func(error)
}

func NewStdioServerTransport() *stdioServerTransport {
	return &stdioServerTransport{}
}
func (s *stdioServerTransport) Start() error {
	return nil
}

func (s *stdioServerTransport) Close() error {
	return nil
}

func (s *stdioServerTransport) SendMessage(message schema.JsonRpcMessage) error {
	return nil
}

func (s *stdioServerTransport) OnReceiveMessage(message schema.JsonRpcMessage) {
	if s.onReceiveMessage != nil {
		s.onReceiveMessage(message)
	}
}

func (s *stdioServerTransport) OnClose() {
	if s.onClose != nil {
		s.onClose()
	}
}

func (s *stdioServerTransport) OnError(err error) {
	if s.onError != nil {
		s.onError(err)
	}
}

func (s *stdioServerTransport) SetOnReceiveMessage(onReceiveMessage func(schema.JsonRpcMessage)) {
	s.onReceiveMessage = onReceiveMessage
}

func (s *stdioServerTransport) SetOnClose(onClose func()) {
	s.onClose = onClose
}

func (s *stdioServerTransport) SetOnError(onError func(error)) {
	s.onError = onError
}
