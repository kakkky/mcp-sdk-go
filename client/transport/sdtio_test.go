package transport

import (
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kakkky/mcp-sdk-go/client"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func TestStdioClientTransport(t *testing.T) {
	tests := []struct {
		name    string
		message schema.JsonRpcMessage
	}{
		{
			name: "normal : ping request",
			message: schema.JsonRpcRequest{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: "2.0",
					Id:      1,
				},
				Request: &schema.PingRequestSchema{
					MethodName: "ping",
				},
			},
		},
		{
			name: "normal : initialize request",
			message: schema.JsonRpcRequest{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: "2.0",
					Id:      1,
				},
				Request: &schema.InitializeRequestSchema{
					MethodName: "initialize",
					ParamsData: schema.InitializeRequestParams{
						ProtocolVersion: "1.0",
						Capabilities:    schema.ClientCapabilities{},
						ClientInfo: schema.Implementation{
							Name:    "test-client",
							Version: "1.0.0",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got schema.JsonRpcMessage

			// catコマンドの存在確認（エコーサーバーとして使用）
			catPath, err := exec.LookPath("cat")
			if err != nil {
				t.Skip("cannot find 'cat' command, skipping test")
			}

			// トランスポート初期化
			transport := NewStdioClientTransport(StdioServerParameters{
				Command: catPath, // catコマンドはエコーサーバーとして機能
			})

			var receiveMsgNotify = make(chan struct{})
			// メッセージ受信ハンドラを設定
			transport.SetOnReceiveMessage(func(msg schema.JsonRpcMessage) {
				// 擬似的に受信したメッセージをチャネルに送信することとする
				receiveMsgNotify <- struct{}{}
				got = msg
			})

			// エラーハンドラを設定
			transport.SetOnError(func(err error) {
				t.Errorf("transport error: %v", err)
			})

			// トランスポート開始
			go func() {
				if err := transport.Start(); err != nil {
					t.Errorf("failed to start transport: %v", err)
				}
			}()
			<-client.TransportStartedNotify
			// メッセージ送信
			if err := transport.SendMessage(tt.message); err != nil {
				t.Errorf("failed to send message: %v", err)
			}

			<-receiveMsgNotify
			// サーバープログラムはcatコマンドなので、送信したメッセージがそのまま返ってくることを確認する
			if diff := cmp.Diff(tt.message, got); diff != "" {
				t.Errorf("received message mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
