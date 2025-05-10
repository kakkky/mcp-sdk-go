package server

import (
	"errors"
	"fmt"

	"github.com/kakkky/mcp-sdk-go/shared/mcp_err"
	"github.com/kakkky/mcp-sdk-go/shared/protocol"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

type ServerOptopns struct {
	Capabilities schema.ServerCapabilities
	Instructions string
	protocol.ProtocolOptions
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
		Protocol:     protocol.NewProtocol(&options.ProtocolOptions),
	}
	// 初期化時のやり取りを行うためのハンドラをセット
	s.SetRequestHandler(&schema.InitializeRequestSchema{}, func(request schema.JsonRpcRequest) (schema.Result, error) {
		return s.onInitialize(request)
	})
	s.SetNotificationHandler(&schema.InitializeNotificationSchema{}, func(notification schema.JsonRpcNotification) error {
		return s.onInitialized()
	})

	s.SetOnAssertCapabilityForMethod(s.assertCapabilityForMethod)

	return s
}

func (s *Server) registerCapabilities(capabilities schema.ServerCapabilities) error {
	if s.Transport() == nil {
		return errors.New("cannot register capabilities after connecting to transport")
	}
	s.capabilities = protocol.MergeCapabilities(s.capabilities, capabilities)
	return nil
}

func (s *Server) assertCapabilityForMethod(method string) error {
	switch method {
	case "sampling/createMessage":
		if s.clientCapabilities.Sampling == nil {
			return fmt.Errorf("client does not support sampling (required for %s)", method)
		}
	case "roots/list":
		if s.clientCapabilities.Roots == nil {
			return fmt.Errorf("client does not support roots (required for %s)", method)
		}
	case "ping":
		break
	}
	return nil
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
