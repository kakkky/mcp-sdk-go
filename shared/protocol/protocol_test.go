package protocol

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kakkky/mcp-sdk-go/shared/mcp_err"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
	gomock "go.uber.org/mock/gomock"
)

func TestProtocol_Connect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// モックの Transport を作成
	mockTransport := NewMocktransport(ctrl)

	// Protocol インスタンスを作成
	protocol := NewProtocol()

	// モックの期待値を設定
	mockTransport.EXPECT().setOnClose(gomock.Any())
	mockTransport.EXPECT().setOnError(gomock.Any())
	mockTransport.EXPECT().setOnReceiveMessage(gomock.Any())
	mockTransport.EXPECT().start()

	// Connect メソッドを呼び出し
	protocol.Connect(mockTransport)

	// Transport が正しく設定されているか確認
	if protocol.Transport() != mockTransport {
		t.Errorf("Expected transport to be set, but it was not")
	}
}

func TestProtocolRequest(t *testing.T) {
	tests := []struct {
		name            string
		request         schema.Request
		expectedResult  schema.Result
		isExpectedErr   bool
		transportMockFn func(ch chan schema.JsonRpcMessage, serverTransport, clientTransport *Mocktransport, server, client *Protocol)
	}{
		{
			name: "nomal case : send request and receive response successfully",
			request: schema.Request{
				Method: "test",
				Params: map[string]any{
					"param1": "value1",
					"param2": "value2",
				}},
			expectedResult: schema.Result{
				Result: map[string]string{
					"status": "success",
				},
			},
			isExpectedErr: false,
			transportMockFn: func(ch chan schema.JsonRpcMessage, serverTransport, clientTransport *Mocktransport, server, client *Protocol) {
				// Connect
				serverTransport.EXPECT().setOnClose(gomock.Any())
				serverTransport.EXPECT().setOnError(gomock.Any())
				serverTransport.EXPECT().setOnReceiveMessage(gomock.Any())
				serverTransport.EXPECT().start().Do(func() {
					go func() {
						for v := range ch {
							serverTransport.onReceiveMessage(v)
						}
					}()
				})
				clientTransport.EXPECT().setOnClose(gomock.Any())
				clientTransport.EXPECT().setOnError(gomock.Any())
				clientTransport.EXPECT().setOnReceiveMessage(gomock.Any())
				clientTransport.EXPECT().start().Do(func() {
					go func() {
						for v := range ch {
							clientTransport.onReceiveMessage(v)
						}
					}()
				})

				// クライアント→サーバーへのリクエスト送信
				clientTransport.EXPECT().sendMessage(gomock.Any()).Do(func(schema.JsonRpcMessage) {
					message := schema.JsonRpcRequest{
						BaseMessage: schema.BaseMessage{
							Jsonrpc: schema.JSON_RPC_VERSION,
							Id:      client.requestMessageId,
						},
						Request: schema.Request{
							Method: "test",
							Params: map[string]any{
								"param1": "value1",
								"param2": "value2",
							},
						},
					}
					ch <- message
				})

				// サーバーがリクエストを受信
				serverTransport.EXPECT().onReceiveMessage(gomock.Any()).Do(func(message schema.JsonRpcMessage) {
					server.onReceiveMessage(message)
				})

				// サーバー→クライアントへのレスポンス送信
				serverTransport.EXPECT().sendMessage(gomock.Any()).Do(func(schema.JsonRpcMessage) {
					message := schema.JsonRpcResponse{
						BaseMessage: schema.BaseMessage{
							Jsonrpc: schema.JSON_RPC_VERSION,
							Id:      1,
						},
						Result: schema.Result{
							Result: map[string]string{
								"status": "success",
							},
						},
					}
					ch <- message
				})

				// クライアントがレスポンスを受信
				clientTransport.EXPECT().onReceiveMessage(gomock.Any()).Do(func(message schema.JsonRpcMessage) {
					client.onReceiveMessage(message)
				})

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// トランスポートのモックを作成
			ServerTransport := NewMocktransport(ctrl)
			ClientTransport := NewMocktransport(ctrl)

			// プロトコルを作成
			server := NewProtocol()
			client := NewProtocol()

			// あらかじめサーバー側のリクエストハンドラを登録
			server.SetRequestHandler(tt.request, func(request schema.JsonRpcRequest) (schema.Result, *mcp_err.McpErr) {
				return schema.Result{
					Result: map[string]string{
						"result": "success",
					},
				}, nil
			})

			// 通信方法としてチャネルを使用
			ch := make(chan schema.JsonRpcMessage, 1)
			// トランスポート（モック）の振る舞いを設定
			tt.transportMockFn(ch, ServerTransport, ClientTransport, server, client)
			// プロトコルにトランスポートを接続
			server.Connect(ServerTransport)
			client.Connect(ClientTransport)

			// クライアント側でリクエストを送信し、レスポンスを受信する
			got, err := client.Request(tt.request, schema.Result{map[string]string{}})
			if (err != nil) != tt.isExpectedErr {
				t.Errorf("Request() error = %v, wantErr %v", err, tt.isExpectedErr)
				return
			}
			if diff := cmp.Diff(got, &tt.expectedResult); diff != "" {
				t.Errorf("Request() got = %v, want %v, diff %s", got, tt.expectedResult, diff)
			}
		})
	}
}
