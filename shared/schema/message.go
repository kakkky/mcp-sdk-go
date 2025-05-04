package schema

const JSON_RPC_VERSION = "2.0"

// Request , Notification, Response の抽象型。
// JsonRpcMessage()メソッド自体は意味をなさない。
type JsonRpcMessage interface {
	JsonRpcMessage()
}

type BaseMessage struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
}

type JsonRpcRequest struct {
	BaseMessage
	Request
}

type JsonRpcNotification struct {
	Jsonrpc string `json:"jsonrpc"`
	Notification
}

type JsonRpcResponse struct {
	BaseMessage
	Result
}

type JsonRpcError struct {
	BaseMessage
	Error Error `json:"error"`
}

func (r JsonRpcRequest) JsonRpcMessage()      {}
func (n JsonRpcNotification) JsonRpcMessage() {}
func (r JsonRpcResponse) JsonRpcMessage()     {}
func (e JsonRpcError) JsonRpcMessage()        {}
