package protocol

import "github.com/kakkky/mcp-sdk-go/shared/schema"

//go:generate  mockgen -source=transport_interface.go -destination=./transport_mock.go -package=protocol
type transport interface {
	start()
	close()
	sendMessage(message schema.JsonRpcMessage) error
	onReceiveMessage(message schema.JsonRpcMessage)
	onClose()
	onError(error)
	onMessage(message schema.JsonRpcMessage)

	// 通信の基本的なイベントはProtocolにより注入される
	setOnClose(onClose func())
	setOnError(onError func(error))
	setOnReceiveMessage(onReceiveMessage func(schema.JsonRpcMessage))
}
