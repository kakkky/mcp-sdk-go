package jsonrpc

import (
	"encoding/json"

	mcperr "github.com/kakkky/mcp-sdk-go/shared/mcp-err"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

type Message struct {
	Jsonrpc string          `json:"jsonrpc"`
	Id      int             `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *struct {
		Code    mcperr.ErrCode  `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data,omitempty"`
	} `json:"error,omitempty"`
}

func Unmarshal(jsonData []byte) (schema.JsonRpcMessage, error) {
	message := &Message{}
	// Jsonrpc,Id,Method,ErrorのCodeとMessageをUnmarshalする
	if err := json.Unmarshal(jsonData, message); err != nil {
		return nil, err
	}
	// メッセージの種類を判定し、Unmarshalする
	switch {
	// Request
	case message.Method != "" && message.Id != 0:
		request, err := unmarshalRequest(message)
		if err != nil {
			return nil, err
		}
		return schema.JsonRpcRequest{
			BaseMessage: schema.BaseMessage{
				Jsonrpc: message.Jsonrpc,
				Id:      message.Id,
			},
			Request: request,
		}, nil

	// Notification
	case message.Method != "" && message.Id == 0:
		notification, err := unmarshalNotification(message)
		if err != nil {
			return nil, err
		}
		return schema.JsonRpcNotification{
			Jsonrpc:      message.Jsonrpc,
			Notification: notification,
		}, nil
	// Error
	case message.Error != nil:
		errorData, err := unmarshalError(message)
		if err != nil {
			return nil, err
		}
		return schema.JsonRpcError{
			BaseMessage: schema.BaseMessage{
				Jsonrpc: message.Jsonrpc,
				Id:      message.Id,
			},
			Error: *errorData,
		}, nil
	// Response
	default:
		result, err := unmarshalResult(message)
		if err != nil {
			return nil, err
		}
		return schema.JsonRpcResponse{
			BaseMessage: schema.BaseMessage{
				Jsonrpc: message.Jsonrpc,
				Id:      message.Id,
			},
			Result: result,
		}, nil
	}
}
