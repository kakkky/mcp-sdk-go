package server

import (
	"github.com/kakkky/mcp-sdk-go/shared/mcp_err"
	"github.com/kakkky/mcp-sdk-go/shared/protocol"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

type ServerOptopns struct {
	Capabilities schema.ServerCapabilities
	Instructions string
}

type Server struct {
	clientCapabilities schema.ClientCapabilities
	clientVersion      schema.Implementation
	capabilities       schema.ServerCapabilities
	instructions       string
	serverInfo         schema.Implementation
	*protocol.Protocol

	onInitialized func() error
}

func NewServer(serverInfo schema.Implementation, options ServerOptopns) *Server {
	s := &Server{
		capabilities: options.Capabilities,
		instructions: options.Instructions,
		serverInfo:   serverInfo,
		Protocol:     protocol.NewProtocol(),
	}
	// 初期化時のやり取りを行うためのハンドラをセット
	s.SetRequestHandler(&schema.InitializeRequestSchema{}, func(request schema.JsonRpcRequest) (schema.Result, error) {
		return s.onInitialize(request)
	})
	s.SetNotificationHandler(&schema.InitializeNotificationSchema{}, func(notification schema.JsonRpcNotification) error {
		return s.onInitialized()
	})

	return s
}

func (s *Server) onInitialize(request schema.JsonRpcRequest) (*schema.InitializeResultSchema, error) {
	requestData, ok := request.Request.(*schema.InitializeRequestSchema)
	if !ok {
		return nil, mcp_err.NewMcpErr(mcp_err.INVALID_REQUEST, "invalid request", nil)
	}

	requestVersion := requestData.ParamsData.ProtocolVersion

	s.clientCapabilities = requestData.ParamsData.Capabilities
	s.clientVersion = requestData.ParamsData.ClientInfo

	return &schema.InitializeResultSchema{
		ProtocolVersion: requestVersion,
		Capabilities:    s.capabilities,
		ServerInfo:      s.serverInfo,
		Instructions:    &s.instructions,
	}, nil

}
