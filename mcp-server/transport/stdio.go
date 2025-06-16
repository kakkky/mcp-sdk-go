package transport

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/kakkky/mcp-sdk-go/shared/schema"
	"github.com/kakkky/mcp-sdk-go/shared/schema/jsonrpc"
	"github.com/kakkky/mcp-sdk-go/shared/transport"
)

type stdioServerTransport struct {
	onReceiveMessage func(schema.JsonRpcMessage)
	onClose          func()
	onError          func(error)

	readBuffer *transport.ReadBuffer
	stdScanner *bufio.Scanner
	isStarted  bool
}

func NewStdioServerTransport() *stdioServerTransport {
	return &stdioServerTransport{
		readBuffer: transport.NewReadBuffer(),
		isStarted:  false,
		stdScanner: bufio.NewScanner(os.Stdin),
	}
}
func (s *stdioServerTransport) Start() error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	if s.isStarted {
		return errors.New("stdio server transport is already started. If using Server class, note that connect() calls start() automatically")
	}
	s.isStarted = true
	go func() {
		s.stdinOnData()
		s.stdinOnError()
	}()
	<-sig
	return nil
}

func (s *stdioServerTransport) Close() error {
	if !s.isStarted {
		return errors.New("stdio server transport is not started")
	}
	s.readBuffer.Clear()
	s.onClose()
	return nil
}

func (s *stdioServerTransport) SendMessage(message schema.JsonRpcMessage) error {
	data, err := jsonrpc.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	// 改行を追加して標準出力に書き込む
	_, err = os.Stdout.Write(append(data, '\n'))
	if err != nil {
		return fmt.Errorf("failed to write message to stdout: %w", err)
	}
	return nil
}

func (s *stdioServerTransport) OnClose() {
	if s.onClose != nil {
		s.onClose()
	}
}

func (s *stdioServerTransport) OnError(err error) {
	if s.onError != nil {
		s.onError(err)
	}
}

func (s *stdioServerTransport) SetOnReceiveMessage(onReceiveMessage func(schema.JsonRpcMessage)) {
	s.onReceiveMessage = onReceiveMessage
}

func (s *stdioServerTransport) SetOnClose(onClose func()) {
	s.onClose = onClose
}

func (s *stdioServerTransport) SetOnError(onError func(error)) {
	s.onError = onError
}

// 標準入力でデータを\nごとに受け取り、onDataコールバックを呼び出す
func (s *stdioServerTransport) stdinOnData() {
	// 標準入力のスキャナーを使用して、データを読み取る
	for s.stdScanner.Scan() {
		data := s.stdScanner.Text()
		// Scannerは改行を含まないので、改行を追加して
		if err := s.onData([]byte(data + "\n")); err != nil {
			s.OnError(fmt.Errorf("failed to read data from stdin: %w", err))
			return
		}
	}
}
func (s *stdioServerTransport) stdinOnError() {
	if s.stdScanner.Err() != nil && s.onError != nil {
		s.onError(s.stdScanner.Err())
	}
}

func (s *stdioServerTransport) onData(chunk []byte) error {
	// バッファにメッセージを書き込む
	if err := s.readBuffer.Append(chunk); err != nil {
		return err
	}
	s.processReadBuffer()
	return nil
}

// バッファからメッセージを読み取り、onReceiveMessageコールバックを呼び出す
func (s *stdioServerTransport) processReadBuffer() {
	msg, err := s.readBuffer.ReadMessage()
	if err != nil {
		s.OnError(err)
	}
	if msg == nil {
		return
	}
	s.onReceiveMessage(msg)
}
