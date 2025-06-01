package transport

import (
	"bytes"
	"strings"

	"github.com/kakkky/mcp-sdk-go/shared/schema"
	"github.com/kakkky/mcp-sdk-go/shared/schema/jsonrpc"
)

type ReadBuffer struct {
	buffer *bytes.Buffer
}

func NewReadBuffer() *ReadBuffer {
	return &ReadBuffer{
		buffer: bytes.NewBuffer(
			[]byte{},
		),
	}
}

// バッファにチャンクを追加する
func (r *ReadBuffer) Append(chunk []byte) error {
	_, err := r.buffer.Write(chunk)
	return err
}

// バッファから一行分メッセージを読み取る
func (r *ReadBuffer) ReadMessage() (schema.JsonRpcMessage, error) {
	// バッファが空の場合
	if r.buffer.Len() == 0 {
		return nil, nil
	}

	// バッファ内容を取得（コピーせず）
	data := r.buffer.Bytes()

	// 改行を探す
	index := bytes.IndexByte(data, '\n')
	if index == -1 {
		return nil, nil
	}

	// 改行までの部分を取り出す
	// JSONRPCメッセージは改行ごとに区切られている
	// 参照：(https://modelcontextprotocol.io/specification/draft/basic/transports#stdio)
	line := make([]byte, index)
	copy(line, data[:index])

	// CRLFを考慮（Windowsの場合）
	if strings.HasSuffix(string(line), "\r") {
		line = line[:len(line)-1] // CRを削除
	}

	// バッファを更新（読み取った部分を削除）
	r.buffer.Next(index + 1)

	// JSONメッセージにデシリアライズ
	message, err := jsonrpc.Unmarshal(line)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// バッファをクリアする
func (r *ReadBuffer) Clear() {
	r.buffer.Reset()
}
