package protocol

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kakkky/mcp-sdk-go/shared/mcp_err"
	"github.com/kakkky/mcp-sdk-go/shared/protocol/test"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func TestProtocol_Connect(t *testing.T) {
	// モックの Transport を作成
	transport := test.NewMockChannelServerTransport(
		make(chan schema.JsonRpcMessage, 1),
		make(chan schema.JsonRpcMessage, 1),
	)

	// Protocol インスタンスを作成
	protocol := NewProtocol()

	// Connect メソッドを呼び出し
	protocol.Connect(transport)

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
				p.SetRequestHandler(schema.Request{Method: "test"}, func(request schema.JsonRpcRequest) (*schema.Result, error) {
					return &schema.Result{
						Result: map[string]string{
							"status": "success",
						}}, nil
				})
			},
			request: schema.Request{
				Method: "test",
			},
			resultSchema: schema.Result{
				Result: map[string]string{},
			},
			expectedResult: schema.Result{
				Result: map[string]string{
					"status": "success",
				},
			},
			isExpectedMcpErr: false,
		},
		{
			name: "sminormal case :client send unknown request and receive 'method not found' error",
			request: schema.Request{
				Method: "unknown",
			},
			expectedErrCode:  mcp_err.METHOD_NOT_FOUND,
			isExpectedMcpErr: true,
		},
		{
			name: "sminormal case :client send unknown request and receive something error (not mcpErr)",
			request: schema.Request{
				Method: "error",
			},
			setHandler: func(p *Protocol) {
				p.SetRequestHandler(schema.Request{Method: "error"}, func(request schema.JsonRpcRequest) (*schema.Result, error) {
					return nil, errors.New("some error")
				})
			},
			isExpectedMcpErr: true,
			expectedErrCode:  mcp_err.INTERNAL_ERROR,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// プロトコルのインスタンスを作成
			server := NewProtocol()
			client := NewProtocol()

			// トランスポートのモックを作成
			serverToClientCh := make(chan schema.JsonRpcMessage, 1)
			clientToServerCh := make(chan schema.JsonRpcMessage, 1)

			// トランスポートを初期化
			serverTransport := test.NewMockChannelServerTransport(clientToServerCh, serverToClientCh)
			clientTransport := test.NewMockChannelClientTransport(clientToServerCh, serverToClientCh)

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
			server.Connect(serverTransport)
			client.Connect(clientTransport)
			// クリーンアップ
			defer func() {
				server.Close()
				client.Close()
			}()
			// リクエストを受け取ったら、レスポンスを返すことを確認する
			got, err := client.Request(tt.request, tt.resultSchema)
			// テストケースがMCPエラーを期待する場合、エラーが期待通りか確認
			if (err != nil) != tt.isExpectedMcpErr {
				t.Errorf("Request() got error = %v, want %v", err, tt.isExpectedMcpErr)
				return
			}
			if tt.isExpectedMcpErr {
				e, ok := err.(*mcp_err.McpErr)
				if !ok {
					t.Errorf("Request() got error = %v, want %v", err, tt.isExpectedMcpErr)
					return
				}
				if e.Code != mcp_err.ErrCode(tt.expectedErrCode) {
					t.Errorf("Request() got error code = %v, want %v", e.Code, tt.expectedErrCode)
					return
				}
				return
			}
			if diff := cmp.Diff(got, &tt.expectedResult); diff != "" {
				t.Errorf("Request() got = %v, want %v, diff %s", got, tt.expectedResult, diff)
			}
		})
	}
}
