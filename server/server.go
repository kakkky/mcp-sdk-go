package server

import (
	"github.com/kakkky/mcp-sdk-go/server/capability"
	"github.com/kakkky/mcp-sdk-go/shared/protocol"
)

type Server struct {
	capabilities *capability.Capability
	protocol     *protocol.Protocol
}
