package transport

import (
	"bytes"

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

// Append adds a chunk of data to the buffer.
func (r *ReadBuffer) Append(chunk []byte) error {
	_, err := r.buffer.Write(chunk)
	return err
}

// ReadMessage tries to read a complete JSON-RPC message from the buffer.
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
	line := make([]byte, index)
	copy(line, data[:index])

	// CRがあれば削除
	if len(line) > 0 && line[len(line)-1] == '\r' {
		line = line[:len(line)-1]
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

// Clear empties the buffer.
func (r *ReadBuffer) Clear() {
	r.buffer.Reset()
}
