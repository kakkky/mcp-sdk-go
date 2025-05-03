package server

import (
	"github.com/kakkky/mcp-sdk-go/server/capability"
	"github.com/kakkky/mcp-sdk-go/server/protocol"
	"github.com/kakkky/mcp-sdk-go/server/tansport"
)

type McpServer struct {
	capabilities *capability.Capability
	transport    *tansport.Transport
	protocol     *protocol.Protocol
}
