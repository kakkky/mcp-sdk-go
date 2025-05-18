package server

import (
	"errors"
	"fmt"

	mcperr "github.com/kakkky/mcp-sdk-go/shared/mcp-err"
	"github.com/kakkky/mcp-sdk-go/shared/protocol"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

type ServerOptions struct {
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
	Protocol

	onInitialized func() error
}

func NewServer(serverInfo schema.Implementation, options *ServerOptions) *Server {
	s := &Server{
		serverInfo: serverInfo,
	}
	if options == nil {
		s.capabilities = schema.ServerCapabilities{}
		s.instructions = ""
		s.Protocol = protocol.NewProtocol(nil)
	} else {
		s.capabilities = options.Capabilities
		s.instructions = options.Instructions
		s.Protocol = protocol.NewProtocol(&options.ProtocolOptions)
	}
	// 初期化時のやり取りを行うためのハンドラをセット
	s.SetRequestHandler(&schema.InitializeRequestSchema{MethodName: "initialize"}, func(request schema.JsonRpcRequest) (schema.Result, error) {
		return s.onInitialize(request)
	})
	s.SetNotificationHandler(&schema.InitializeNotificationSchema{MethodName: "notifications/initialized"}, func(notification schema.JsonRpcNotification) error {
		if s.onInitialized != nil {
			return s.onInitialized()
		}
		return nil
	})

	s.SetValidateCapabilityForMethod(s.validateCapabilityForMethod)
	s.SetValidateNotificationCapability(s.validateNotificationCapability)

	return s
}

func (s *Server) RegisterCapabilities(capabilities schema.ServerCapabilities) error {
	if s.Transport() == nil {
		return errors.New("cannot register capabilities after connecting to transport")
	}
	s.capabilities = protocol.MergeCapabilities(s.capabilities, capabilities)
	return nil
}

func (s *Server) onInitialize(request schema.JsonRpcRequest) (*schema.InitializeResultSchema, error) {
	requestData, ok := request.Request.(*schema.InitializeRequestSchema)
	if !ok {
		return nil, mcperr.NewMcpErr(mcperr.INVALID_REQUEST, "invalid request", nil)
	}

	requestVersion := requestData.ParamsData.ProtocolVersion

	s.clientCapabilities = requestData.ParamsData.Capabilities
	s.clientVersion = requestData.ParamsData.ClientInfo

	return &schema.InitializeResultSchema{
		ProtocolVersion: requestVersion,
		Capabilities:    s.capabilities,
		ServerInfo:      s.serverInfo,
		Instructions:    s.instructions,
	}, nil

}

// 基本的な通信メソッド
func (s *Server) Ping() (schema.Result, error) {
	return s.Request(&schema.PingRequestSchema{
		MethodName: "ping",
	}, &schema.EmptyResultSchema{})
}

func (s *Server) CreateMessage(params any, contentType string) (schema.Result, error) {
	switch contentType {
	case "text":
		typedParams, ok := params.(schema.CreateMessageRequestParams[schema.TextContentSchema])
		if !ok {
			return nil, fmt.Errorf("invalid params type: %T", params)
		}
		return s.Request(&schema.CreateMessageRequestSchema[schema.TextContentSchema]{
			MethodName: "sampling/createMessage",
			ParamsData: typedParams,
		}, &schema.CreateMessageResultSchema[schema.TextContentSchema]{})
	case "image":
		typedParams, ok := params.(schema.CreateMessageRequestParams[schema.ImageContentSchema])
		if !ok {
			return nil, fmt.Errorf("invalid params type: %T", params)
		}
		return s.Request(&schema.CreateMessageRequestSchema[schema.ImageContentSchema]{
			MethodName: "sampling/createMessage",
			ParamsData: typedParams,
		}, &schema.CreateMessageResultSchema[schema.ImageContentSchema]{})
	case "audio":
		typedParams, ok := params.(schema.CreateMessageRequestParams[schema.AudioContentSchema])
		if !ok {
			return nil, fmt.Errorf("invalid params type: %T", params)
		}
		return s.Request(&schema.CreateMessageRequestSchema[schema.AudioContentSchema]{
			MethodName: "sampling/createMessage",
			ParamsData: typedParams,
		}, &schema.CreateMessageResultSchema[schema.AudioContentSchema]{})
	}
	return nil, fmt.Errorf("invalid content type: %s", contentType)
}

func (s *Server) ListRoots() (schema.Result, error) {
	return s.Request(&schema.ListRootsRequestSchema{
		MethodName: "roots/list",
	}, &schema.ListRootResultSchema{})
}

func (s *Server) SendLogginMessage(params schema.LoggingMessageNotificationParams) error {
	return s.Notificate(&schema.LoggingMessageNotificationSchema{
		MethodName: "notifications/message",
		ParamsData: params,
	})
}

func (s *Server) SendResourceUpdated(params schema.ResourceUpdatedNotificationParams) error {
	return s.Notificate(&schema.ResourceUpdatedNotificationSchema{
		MethodName: "notifications/resources/updated",
		ParamsData: params,
	})
}

func (s *Server) SendResourceListChanged() error {
	return s.Notificate(&schema.ResourceListChangedNotificationSchema{
		MethodName: "notifications/resources/list_changed",
	})
}

func (s *Server) SendToolListChanged() error {
	return s.Notificate(&schema.ToolListChangedNotificationSchema{
		MethodName: "notifications/tools/list_changed",
	})
}

func (s *Server) SendPromptListChanged() error {
	return s.Notificate(&schema.PromptListChangedNotificationSchema{
		MethodName: "notifications/prompts/list_changed",
	})
}
