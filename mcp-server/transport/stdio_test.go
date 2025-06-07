package transport

import (
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func TestStdioServeTransport(t *testing.T) {
	tests := []struct {
		name string
		msg  string
	}{
		{
			name: "normal: send ping request",
			msg:  `{"jsonrpc":"2.0","id":1,"method":"ping"}`,
		},
		{
			name: "normal: send initialize request",
			msg:  `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"0.1.0"}}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 標準入出力のバックアップ
			originStdout := os.Stdout
			originStdin := os.Stdin

			stdinR, stdinW, _ := os.Pipe()
			stdoutR, stdoutW, _ := os.Pipe()

			os.Stdout = stdoutW // transportからのメッセージ受信をキャプチャするため
			os.Stdin = stdinR   // transportへのメッセージ送信をキャプチャするため

			// テスト終了時に元に戻す
			defer func() {
				os.Stdout = originStdout
				os.Stdin = originStdin
				stdinR.Close()
				stdoutR.Close()
				stdinW.Close()
				stdoutW.Close()
			}()

			sut := NewStdioServerTransport()

			// メッセージ送信時に通知する
			sendMsgDoneChan := make(chan struct{})
			// メッセージ受信時のコールバックを設定
			// ここでは、受信したメッセージをそのまま送信する
			sut.SetOnReceiveMessage(func(jrm schema.JsonRpcMessage) {
				sut.SendMessage(jrm)
				sendMsgDoneChan <- struct{}{}
			})

			// 非同期でStartを実行
			go func() {
				err := sut.Start()
				if err != nil {
					t.Errorf("Start() error = %v", err)
				}
			}()

			// テストデータを標準入力に書き込む
			// 注意: 実際のJSONメッセージを送信する必要がある
			_, err := stdinW.WriteString(tt.msg + "\n")
			if err != nil {
				t.Fatalf("Failed to write to stdin: %v", err)
			}

			// Transportがメッセージを返すまで待機
			<-sendMsgDoneChan
			stdoutW.Close() // 標準出力の書き込みを終了

			// 標準出力を読み取る
			stdoutResp, err := io.ReadAll(stdoutR)
			if err != nil {
				t.Fatalf("Failed to read from stdout: %v", err)
			}

			if diff := cmp.Diff(tt.msg+"\n", string(stdoutResp)); diff != "" {
				t.Errorf("Output mismatch (-want +got):\n%s", diff)
			}

			// シグナルを送信してStart()を終了させる
			proc, _ := os.FindProcess(os.Getpid())
			proc.Signal(os.Interrupt)
		})
	}
}
