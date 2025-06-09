package client

import (
	"github.com/kakkky/mcp-sdk-go/shared"
	"github.com/kakkky/mcp-sdk-go/shared/protocol"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

type ClientOptions struct {
	capabilities schema.ClientCapabilities
	protocol.ProtocolOptions
}

// プラグイン可能なトランスポートの上に構築されたMCPクライアント。
// connect() が呼び出されると、クライアントはサーバーとの初期化フローを自動的に開始する。
type Client struct {
	serverCapabilities schema.ServerCapabilities
	serverVersion      schema.Implementation
	capabilities       schema.ClientCapabilities
	instruction        string
	clientInfo         schema.Implementation
	shared.Protocol
}

func NewClient(clientInfo schema.Implementation, options *ClientOptions) *Client {
	c := &Client{
		clientInfo: clientInfo,
	}
	if options == nil {
		c.capabilities = schema.ClientCapabilities{}
		c.Protocol = protocol.NewProtocol(nil)
	} else {
		c.capabilities = options.capabilities
		c.Protocol = protocol.NewProtocol(&options.ProtocolOptions)
	}

	return c

}
