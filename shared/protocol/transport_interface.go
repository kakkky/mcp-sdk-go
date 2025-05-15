package protocol

import "github.com/kakkky/mcp-sdk-go/shared/schema"

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
