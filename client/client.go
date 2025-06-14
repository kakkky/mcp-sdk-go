package client

import (
	"errors"
	"fmt"
	"log"

	"github.com/kakkky/mcp-sdk-go/shared"
	"github.com/kakkky/mcp-sdk-go/shared/protocol"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

type ClientOptions struct {
	Capabilities schema.ClientCapabilities
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
		c.capabilities = options.Capabilities
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

// Transport側で接続が確立されたことを通知するためのチャネル
var TransportStartedNotify = make(chan struct{}, 1)

// Initialization phaseが完了し、Operation phaseを開始するための通知チャネル
var OpetationPhaseStartNotify = make(chan struct{}, 1)

func (c *Client) Connect(transport protocol.Transport) error {
	if transport == nil {
		return errors.New("transport is required")
	}
	// transportに接続
	go func() {
		if err := c.Protocol.Connect(transport); err != nil {
			if err := c.Close(); err != nil {
				log.Fatalln("Failed to close protocol after connection error:", err)
			}
			log.Fatalf("failed to connect to transport: %v", err)
		}
	}()
	// transportへの接続が確立後に後続のinitialiation phaseを開始する
	<-TransportStartedNotify
	// initializeリクエスト
	result, err := c.Request(&schema.InitializeRequestSchema{
		MethodName: "initialize",
		ParamsData: schema.InitializeRequestParams{
			ProtocolVersion: schema.LATEST_PROTOCOL_VERSION,
			Capabilities:    c.capabilities,
			ClientInfo:      c.clientInfo,
		},
	}, &schema.InitializeResultSchema{})
	if err != nil {
		if err := c.Close(); err != nil {
			fmt.Println("Failed to close protocol after connection error:", err)
		}
		return fmt.Errorf("failed to initialize: %w", err)
	}
	if result == nil {
		if err := c.Close(); err != nil {
			fmt.Println("Failed to close protocol after connection error:", err)
		}
		return fmt.Errorf("server sent invalid initialize result")
	}
	initializeResult, ok := result.(*schema.InitializeResultSchema)
	if !ok {
		if err := c.Close(); err != nil {
			fmt.Println("Failed to close protocol after connection error:", err)
		}
		return fmt.Errorf("server sent invalid initialize result: %T", result)
	}

	// サーバーのプロトコルバージョンがサポートされているかを確認
	protocolVersion := initializeResult.ProtocolVersion
	for i := 0; i < len(schema.SUPPORTED_PROTOCOL_VERSIONS); i++ {
		if protocolVersion == schema.SUPPORTED_PROTOCOL_VERSIONS[i] {
			break
		}
		if i == len(schema.SUPPORTED_PROTOCOL_VERSIONS)-1 {
			if err := c.Close(); err != nil {
				fmt.Println("Failed to close protocol after connection error:", err)
			}
			fmt.Println("Server's protocol version is not supported:", protocolVersion)
			return fmt.Errorf("server's protocol version is not supported: %s", protocolVersion)
		}
	}
	c.serverCapabilities = initializeResult.Capabilities
	c.serverVersion = initializeResult.ServerInfo
	c.instruction = initializeResult.Instructions
	if err := c.Notificate(&schema.InitializeNotificationSchema{
		MethodName: "notifications/initialized",
	}); err != nil {
		if err := c.Close(); err != nil {
			fmt.Println("Failed to close protocol after connection error:", err)
		}
		return fmt.Errorf("failed to send initialized notification: %w", err)
	}
	OpetationPhaseStartNotify <- struct{}{}
	return nil
}
