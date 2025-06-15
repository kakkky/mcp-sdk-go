package transport

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/kakkky/mcp-sdk-go/client"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
	"github.com/kakkky/mcp-sdk-go/shared/schema/jsonrpc"
	"github.com/kakkky/mcp-sdk-go/shared/transport"
)

type IOType string

const (
	INHERIT    IOType = "inherit"    // サーバープロセス内のstderrを親プロセスに継承（クライアント側でそのまま出力される）
	PIPE       IOType = "pipe"       // パイプを使用し、プログラムから読み取れるようにする
	OVERLAPPED IOType = "overlapped" // Windowsのオーバーラップモードでパイプを使用
	IGNORE     IOType = "ignore"     // 無視
)

type StdioServerParameters struct {
	Command string
	Args    []string
	Env     []string // "key=value"形式の環境変数のリスト
	Stderr  IOType   // サーバープロセスの標準エラー出力の取り扱い
	// サーバーのプロセスが実行されるディレクトリ
	// 指定しなかった場合は、curent working directoryが使用される
	Cwd string
}

// サーバープロセスに継承されるデフォルトの環境変数のリスト
func defaultInheritedEnvVars() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{
			"APPDATA",
			"HOMEDRIVE",
			"HOMEPATH",
			"LOCALAPPDATA",
			"PATH",
			"PROCESSOR_ARCHITECTURE",
			"SYSTEMDRIVE",
			"SYSTEMROOT",
			"TEMP",
			"USERNAME",
			"USERPROFILE",
		}
	default:
		return []string{
			"HOME", "LOGNAME", "PATH", "SHELL", "TERM", "USER",
		}
	}
}

// 起動するサーバープログラムに渡すための環境変数を取得する
func getDefaultEnvironment() []string {
	env := []string{}
	for _, key := range defaultInheritedEnvVars() {
		value, ok := os.LookupEnv(key)
		if !ok || strings.HasPrefix(value, "()") {
			continue // 環境変数が存在しないか、値が不正な場合はスキップ
		}
		env = append(env, key+"="+value)
	}
	return env
}

type StdioClientTransport struct {
	process      *exec.Cmd
	readBuffer   *transport.ReadBuffer
	serverParams StdioServerParameters
	stderrChan   chan error     // PIPEの場合にサーバープロセスの標準エラー出力を受け取るチャネル
	stdinPipe    io.WriteCloser // 標準入力のパイプ（サーバープロセスにメッセージを送信するため）
	stdoutPipe   io.ReadCloser  // 標準出力のパイプ（サーバープロセスからのメッセージを受信するため）

	onReceiveMessage func(schema.JsonRpcMessage)
	onClose          func()
	onError          func(error)
}

func NewStdioClientTransport(server StdioServerParameters) *StdioClientTransport {
	s := &StdioClientTransport{
		serverParams: server,
		readBuffer:   transport.NewReadBuffer(),
	}
	if (s.serverParams.Stderr == PIPE) || (s.serverParams.Stderr == OVERLAPPED) {
		s.stderrChan = make(chan error, 1) // 標準エラー出力を受け取るチャネルを初期化
	}
	return s
}

func (s *StdioClientTransport) Start() error {
	s.process = exec.Command(s.serverParams.Command, s.serverParams.Args...)
	if len(s.serverParams.Env) > 0 {
		s.process.Env = s.serverParams.Env
	} else {
		s.process.Env = getDefaultEnvironment()
	}
	stdinPipe, err := s.process.StdinPipe()
	if err != nil {
		return err
	}
	s.stdinPipe = stdinPipe // 標準入力のパイプ
	stdoutPipe, err := s.process.StdoutPipe()
	if err != nil {
		return err
	}
	s.stdoutPipe = stdoutPipe // 標準出力のパイプ
	switch s.serverParams.Stderr {
	case INHERIT:
		s.process.Stderr = os.Stderr // 標準エラー出力を親プロセスに継承
	case PIPE, OVERLAPPED:
		// チャネルで標準エラー出力を受け取る
		stderrPipe, err := s.process.StderrPipe()
		if err != nil {
			return err
		}
		go func() {
			defer stderrPipe.Close()
			for {
				data := make([]byte, 1024)
				n, err := stderrPipe.Read(data)
				if err != nil {
					if err.Error() != io.EOF.Error() {
						s.OnError(err) // エラーを通知
					}
					break
				}
				if n > 0 {
					s.stderrChan <- errors.New(string(data[:n])) // チャネルにデータを送信
				}
			}
		}()
	case IGNORE:
		s.process.Stderr = nil // 標準エラー出力を無視
	default:
		// defaultはINHERITと同じ扱い
		s.process.Stderr = os.Stderr // 標準エラー出力を親プロセスに継承
	}
	s.process.Dir = s.serverParams.Cwd
	if err := s.process.Start(); err != nil {
		return err
	}
	go func() {
		client.TransportStartedNotify <- struct{}{}
		s.stdinOnData() // 標準入力からのデータを読み取る
	}()
	return nil
}

func (s *StdioClientTransport) SendMessage(message schema.JsonRpcMessage) error {
	data, err := jsonrpc.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	// 改行を追加して標準入力に書き込む
	fmt.Println("Client:", string(data))
	_, err = s.stdinPipe.Write(append(data, '\n'))
	if err != nil {
		return fmt.Errorf("client failed to write message to stdout: %w", err)
	}
	return nil
}

func (s *StdioClientTransport) Close() error {
	if s.process.Process == nil {
		return fmt.Errorf("process is not running")
	}
	if err := s.process.Process.Kill(); err != nil {
		return fmt.Errorf("failed to kill process: %w", err)
	}
	s.readBuffer.Clear()
	return nil
}

func (s *StdioClientTransport) Stderr() <-chan error {
	if s.stderrChan == nil {
		return nil // 標準エラー出力を受け取るチャネルがない場合はnilを返す
	}
	return s.stderrChan

}

func (s *StdioClientTransport) OnReceiveMessage(message schema.JsonRpcMessage) {
	if s.onReceiveMessage != nil {
		s.onReceiveMessage(message)
	}
}

func (s *StdioClientTransport) OnClose() {
	if s.onClose != nil {
		s.onClose()
	}
}

func (s *StdioClientTransport) OnError(err error) {
	if s.onError != nil {
		s.onError(err)
	}
}

func (s *StdioClientTransport) SetOnReceiveMessage(onReceiveMessage func(schema.JsonRpcMessage)) {
	s.onReceiveMessage = onReceiveMessage
}

func (s *StdioClientTransport) SetOnClose(onClose func()) {
	s.onClose = onClose
}

func (s *StdioClientTransport) SetOnError(onError func(error)) {
	s.onError = onError
}

// 標準入力でデータを\nごとに受け取り、onDataコールバックを呼び出す
func (s *StdioClientTransport) stdinOnData() {
	scanner := bufio.NewScanner(s.stdoutPipe)
	// 標準入力のスキャナーを使用して、データを読み取る
	for scanner.Scan() {
		data := scanner.Text()
		// Scannerは改行を含まないので、改行を追加して
		if err := s.onData([]byte(data + "\n")); err != nil {
			s.OnError(fmt.Errorf("failed to read data from stdin: %w", err))
			return
		}
	}
}

func (s *StdioClientTransport) onData(chunk []byte) error {
	// バッファにメッセージを書き込む
	if err := s.readBuffer.Append(chunk); err != nil {
		return err
	}
	fmt.Println("Server:", string(chunk))
	s.processReadBuffer()
	return nil
}

// バッファからメッセージを読み取り、onReceiveMessageコールバックを呼び出す
func (s *StdioClientTransport) processReadBuffer() {
	msg, err := s.readBuffer.ReadMessage()
	if err != nil {
		s.OnError(err)
	}
	if msg == nil {
		return
	}
	s.onReceiveMessage(msg)
}
