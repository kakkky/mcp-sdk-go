package protocol

import "github.com/kakkky/mcp-sdk-go/shared/schema"

func (p *Protocol) onResponse(message schema.JsonRpcMessage) {}

func (p *Protocol) onRequest(message schema.JsonRpcMessage) {}

func (p *Protocol) onNotification(message schema.JsonRpcMessage) {}
