package schema

import "github.com/kakkky/mcp-sdk-go/shared/mcp_err"

const JSON_RPC_VERSION = "2.0"

// Request , Notification, Response の抽象型。
// JsonRpcMessage()メソッド自体は意味をなさない。
type JsonRpcMessage interface {
	JsonRpcMessage()
}

type JsonRpcRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      string `json:"id"`
	Request
}

type JsonRpcNotification struct {
	Jsonrpc string `json:"jsonrpc"`
	Notification
}

type JsonRpcResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      string `json:"id"`
	Result
}

type JsonRpcError struct {
	Code    mcp_err.ErrCode `json:"code"`
	Message string          `json:"message"`
	Data    any             `json:"data,omitempty"`
}

func (r JsonRpcRequest) JsonRpcMessage()      {}
func (n JsonRpcNotification) JsonRpcMessage() {}
func (r JsonRpcResponse) JsonRpcMessage()     {}
func (e JsonRpcError) JsonRpcMessage()        {}
