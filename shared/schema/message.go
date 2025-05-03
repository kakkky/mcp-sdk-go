package schema

const JSON_RPC_VERSION = "2.0"

type JsonRpcRequeststruct[T any] struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      string `json:"id"`
	Method  string `json:"method"`
	Params  T      `json:"params"`
}

type JsonRpcNotification struct {
	Jsonrc string `json:"jsonrpc"`
	Method string `json:"method"`
}

type JsonRpcResponse[T, U any] struct {
	Jsonrpc string           `json:"jsonrpc"`
	Id      string           `json:"id"`
	Result  T                `json:"result"`
	Error   *JsonRpcError[U] `json:"error,omitempty"`
}
type JsonRpcError[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}
