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
	Method  string `json:"method"`
	Params  any    `json:"params"`
}

func (r JsonRpcRequest) JsoRpcnMessage() {}

type JsonRpcNotification struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
}

func (n JsonRpcNotification) JsonRpcMessage() {}

type JsonRpcResponse struct {
	Jsonrpc string        `json:"jsonrpc"`
	Id      string        `json:"id"`
	Result  any           `json:"result"`
	Error   *JsonRpcError `json:"error,omitempty"`
}

func (r JsonRpcResponse) JsonRpcMessage() {}

type JsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
