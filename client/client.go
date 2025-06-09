package client

import (
	"errors"
	"fmt"

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

func (c *Client) RegisterCapabilities(capabilities schema.ClientCapabilities) error {
	if c.Transport() == nil {
		return errors.New("cannot register capabilities after connecting to transport")
	}
	c.capabilities = protocol.MergeCapabilities(c.capabilities, capabilities)
	return nil
}

// 引数にとったCapabilityをServerが提供しているかを検証する。
func (c *Client) ValidateCapabilities(capability any, method string) error {
	switch capability.(type) {
	case *schema.Logging:
		if c.serverCapabilities.Logging == nil {
			return fmt.Errorf("%s requires logging capability which is not supported by the server", method)
		}
	case *schema.Completion:
		if c.serverCapabilities.Completion == nil {
			return fmt.Errorf("%s requires completion capability which is not supported by the server", method)
		}
	case *schema.Prompts:
		if c.serverCapabilities.Prompts == nil {
			return fmt.Errorf("%s requires prompts capability which is not supported by the server", method)
		}
	case *schema.Resources:
		if c.serverCapabilities.Resources == nil {
			return fmt.Errorf("%s requires resources capability which is not supported by the server", method)
		}
	case *schema.Tools:
		if c.serverCapabilities.Tools == nil {
			return fmt.Errorf("%s requires tools capability which is not supported by the server", method)
		}
	default:
		return fmt.Errorf("%s unknown capability type for method", method)
	}
	return nil
}
