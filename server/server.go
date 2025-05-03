package server

import (
	"github.com/kakkky/mcp-sdk-go/server/capability"
	"github.com/kakkky/mcp-sdk-go/server/tansport"
	"github.com/kakkky/mcp-sdk-go/shared/protocol"
)

type McpServer struct {
	capabilities *capability.Capability
	transport    *tansport.Transport
	protocol     *protocol.Protocol
}
