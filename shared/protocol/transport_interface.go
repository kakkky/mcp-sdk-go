package protocol

import "github.com/kakkky/mcp-sdk-go/shared/schema"

//go:generate mockgen -source=./transport_interface.go -destination=./mock/transport_mock.go -package=mock
//go:generate mockgen -source=./transport_interface.go -destination=../../mcp-server/server/mock/transport_mock.go -package=mock
type Transport interface {
	Start() error
	Close() error
	SendMessage(message schema.JsonRpcMessage) error
	OnReceiveMessage(message schema.JsonRpcMessage)
	OnClose()
	OnError(error)

	// 通信の基本的なイベントはProtocolにより注入される
	SetOnClose(onClose func())
	SetOnError(onError func(error))
	SetOnReceiveMessage(onReceiveMessage func(schema.JsonRpcMessage))
}
