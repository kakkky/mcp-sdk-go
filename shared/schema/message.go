package schema

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
	Code    errCode `json:"code"`
	Message string  `json:"message"`
	Data    any     `json:"data,omitempty"`
}

type errCode int

const (
	// Custom error codes
	CONNECTION_CLOSED = -32000
	REQUEST_TIMEOUT   = -32001

	// Standard JSON-RPC error codes
	PARSE_ERROR      = -32700
	INVALID_REQUEST  = -32600
	METHOD_NOT_FOUND = -32601
	INVALID_PARAMS   = -32602
	INTERNAL_ERROR   = -32603
)

func (r JsonRpcRequest) JsonRpcMessage()      {}
func (n JsonRpcNotification) JsonRpcMessage() {}
func (r JsonRpcResponse) JsonRpcMessage()     {}
func (e JsonRpcError) JsonRpcMessage()        {}
