package protocol

import "github.com/kakkky/mcp-sdk-go/shared/schema"

type transport interface {
	start()
	close()
	sendMessage(message schema.JsonRpcMessage)
	onReceiveMessage(message schema.JsonRpcMessage)
	onClose()
	onError(error)
	onMessage(message schema.JsonRpcMessage)

	// 通信の基本的なイベントはProtocolにより注入される
	setOnClose(onClose func())
	setOnError(onError func(error))
	setOnMessage(onMessage func(schema.JsonRpcMessage))
}
