package protocol

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	mcperr "github.com/kakkky/mcp-sdk-go/shared/mcp-err"
	"github.com/kakkky/mcp-sdk-go/shared/protocol/mock"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func TestProtocol_Connect(t *testing.T) {
	// モックの Transport を作成
	transport := mock.NewMockChannelServerTransport(
		make(chan []byte, 1),
		make(chan []byte, 1),
	)

	// Protocol インスタンスを作成
	protocol := NewProtocol(nil)

	// Connect メソッドを呼び出し
	err := protocol.Connect(transport)
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	// Transport が正しく設定されているか確認
	if protocol.Transport() != transport {
		t.Errorf("Expected transport to be set, but it was not")
	}
}

func TestProtocol_Request(t *testing.T) {
	tests := []struct {
		name             string
		setHandler       func(p *Protocol)
		request          schema.Request
		resultSchema     schema.Result
		expectedResult   schema.Result
		expectedErrCode  int
		isExpectedMcpErr bool
	}{
		{
			name: "nomal case :client send request and receive response successfully",
			setHandler: func(p *Protocol) {
				p.SetRequestHandler(&schema.InitializeRequestSchema{MethodName: "initialize"}, func(request schema.JsonRpcRequest) (schema.Result, error) {
					return &schema.InitializeResultSchema{
						ServerInfo: schema.Implementation{
							Name:    "test-server",
							Version: "1.0.0",
						},
						Capabilities: schema.ServerCapabilities{
							Tools: &schema.Tools{
								ListChanged: true,
							},
						},
					}, nil
				})
			},
			request: &schema.InitializeRequestSchema{
				MethodName: "initialize",
				ParamsData: schema.InitializeRequestParams{
					ProtocolVersion: "2025-01-01",
					Capabilities: schema.ClientCapabilities{
						Sampling: &schema.Sampling{},
					},
					ClientInfo: schema.Implementation{
						Name:    "test-client",
						Version: "1.0.0",
					},
				},
			},
			resultSchema: &schema.InitializeResultSchema{},
			expectedResult: &schema.InitializeResultSchema{
				ServerInfo: schema.Implementation{
					Name:    "test-server",
					Version: "1.0.0",
				},
				Capabilities: schema.ServerCapabilities{
					Tools: &schema.Tools{
						ListChanged: true,
					},
				},
			},
			isExpectedMcpErr: false,
		},
		{
			name:             "semi normal case :client send unknown request and receive 'method not found' error",
			request:          &schema.ListResourceRequestSchema{MethodName: "resources/list"},
			expectedErrCode:  mcperr.METHOD_NOT_FOUND,
			isExpectedMcpErr: true,
		},
		{
			name:    "semi normal case :client send request and receive something error (not mcpErr)",
			request: &schema.PingRequestSchema{MethodName: "ping"},
			setHandler: func(p *Protocol) {
				p.SetRequestHandler(&schema.PingRequestSchema{MethodName: "ping"}, func(request schema.JsonRpcRequest) (schema.Result, error) {
					return nil, errors.New("some error")
				})
			},
			isExpectedMcpErr: true,
			expectedErrCode:  mcperr.INTERNAL_ERROR,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// プロトコルのインスタンスを作成
			server := NewProtocol(nil)
			client := NewProtocol(nil)

			// トランスポートのモックを作成
			serverToClientCh := make(chan []byte, 1)
			clientToServerCh := make(chan []byte, 1)

			// トランスポートを初期化
			serverTransport := mock.NewMockChannelServerTransport(clientToServerCh, serverToClientCh)
			clientTransport := mock.NewMockChannelClientTransport(clientToServerCh, serverToClientCh)

			// Close時コールバックを登録
			server.SetOnClose(func() {
				close(serverToClientCh)
			})
			client.SetOnClose(func() {
				close(clientToServerCh)
			})

			// サーバー側でリクエストハンドラを登録
			if tt.setHandler != nil {
				tt.setHandler(server)
			}

			// 通信を開始
			if err := server.Connect(serverTransport); err != nil {
				t.Errorf("Connect() error = %v", err)
				return
			}
			if err := client.Connect(clientTransport); err != nil {
				t.Errorf("Connect() error = %v", err)
				return
			}
			// クリーンアップ
			defer func() {
				if err := server.Close(); err != nil {
					t.Errorf("Close() error = %v", err)
					return
				}
				if err := client.Close(); err != nil {
					t.Errorf("Close() error = %v", err)
					return
				}
			}()
			// リクエストを受け取ったら、レスポンスを返すことを確認する
			got, err := client.Request(tt.request, tt.resultSchema)
			// テストケースがMCPエラーを期待する場合、エラーが期待通りか確認
			if tt.isExpectedMcpErr {
				if err == nil {
					t.Errorf("Request() got error = %v, want %v", err, tt.isExpectedMcpErr)
					return
				}
				e, ok := err.(*mcperr.McpErr)
				if !ok {
					t.Errorf("Request() got error = %v, want %v", err, tt.isExpectedMcpErr)
					return
				}
				if e.Code != mcperr.ErrCode(tt.expectedErrCode) {
					t.Errorf("Request() got error code = %v, want %v", e.Code, tt.expectedErrCode)
					return
				}
				return
			}
			if err != nil {
				t.Errorf("Request() got error = %v", err)
				return
			}
			if diff := cmp.Diff(got, tt.expectedResult); diff != "" {
				t.Errorf("Request() got = %v, want %v, diff %s", got, tt.expectedResult, diff)
			}
		})
	}
}

func TestProtocol_Notificate(t *testing.T) {
	tests := []struct {
		name         string
		notification schema.Notification
	}{
		{
			name: "normal case :client send notification successfully",
			notification: &schema.InitializeNotificationSchema{
				MethodName: "notifications/initialize",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// プロトコルのインスタンスを作成
			server := NewProtocol(nil)
			client := NewProtocol(nil)

			// トランスポートのモックを作成
			serverToClientCh := make(chan []byte, 1)
			clientToServerCh := make(chan []byte, 1)

			// トランスポートを初期化
			serverTransport := mock.NewMockChannelServerTransport(clientToServerCh, serverToClientCh)
			clientTransport := mock.NewMockChannelClientTransport(clientToServerCh, serverToClientCh)

			// Close時コールバックを登録
			server.SetOnClose(func() {
				close(serverToClientCh)
			})
			client.SetOnClose(func() {
				close(clientToServerCh)
			})

			// 通信を開始
			if err := server.Connect(serverTransport); err != nil {
				t.Errorf("サーバー接続に失敗: %v", err)
			}
			if err := client.Connect(clientTransport); err != nil {
				t.Errorf("クライアント接続に失敗: %v", err)
				return
			}
			// クリーンアップ
			defer func() {
				if err := server.Close(); err != nil {
					t.Errorf("サーバークローズに失敗: %v", err)
				}
				if err := client.Close(); err != nil {
					t.Errorf("クライアントクローズに失敗: %v", err)
				}
			}()

			// 通知を送信
			err := client.Notificate(tt.notification)
			if err != nil {
				t.Errorf("Notify() got error = %v", err)
				return
			}

		})
	}
}

func TestProtocol_Close(t *testing.T) {
	tests := []struct {
		name     string
		onClose  func()
		expected bool
	}{
		{
			name: "normal case :set onClose callback",
			onClose: func() {
				// do nothing
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			protocol := NewProtocol(nil)
			serverTransport := mock.NewMockChannelServerTransport(
				make(chan []byte, 1),
				make(chan []byte, 1),
			)
			protocol.SetOnClose(tt.onClose)
			if err := protocol.Connect(serverTransport); err != nil {
				t.Errorf("接続に失敗: %v", err)
			}

			if err := protocol.Close(); err != nil {
				t.Errorf("クローズに失敗: %v", err)
			}
			if protocol.Transport() != nil {
				t.Errorf("Expected transport to be nil after close, but it was not")
			}

		})
	}
}
