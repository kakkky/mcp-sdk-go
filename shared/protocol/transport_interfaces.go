package protocol

import "github.com/kakkky/mcp-sdk-go/shared/schema"

type Transport interface {
	Start()
	Close()
	SendMessage(message schema.JsonRpcMessage)
	OnReceiveMessage(message schema.JsonRpcMessage)
}
