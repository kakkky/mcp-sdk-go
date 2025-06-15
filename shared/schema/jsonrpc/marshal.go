package jsonrpc

import (
	"encoding/json"
	"fmt"

	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func Marshal(message schema.JsonRpcMessage) ([]byte, error) {
	switch m := message.(type) {
	case schema.JsonRpcRequest:
		return marshalRequest(m)
	case schema.JsonRpcResponse:
		return json.Marshal(m)
	case schema.JsonRpcNotification:
		return marshalNotification(m)
	case schema.JsonRpcError:
		return json.Marshal(m)
	default:
		return nil, fmt.Errorf("unsupported message type: %T", message)
	}
}

func marshalRequest(req schema.JsonRpcRequest) ([]byte, error) {
	type requestJSON struct {
		Jsonrpc string      `json:"jsonrpc"`
		Id      int         `json:"id"`
		Method  string      `json:"method"`
		Params  interface{} `json:"params,omitempty"`
	}

	jsonObj := requestJSON{
		Jsonrpc: req.Jsonrpc,
		Id:      req.Id,
		Method:  req.Method(),
	}

	if req.Params() != nil {
		jsonObj.Params = req.Params()
	}

	return json.Marshal(jsonObj)
}

func marshalNotification(notif schema.JsonRpcNotification) ([]byte, error) {
	type notificationJSON struct {
		Jsonrpc string      `json:"jsonrpc"`
		Method  string      `json:"method"`
		Params  interface{} `json:"params,omitempty"`
	}

	jsonObj := notificationJSON{
		Jsonrpc: notif.Jsonrpc,
		Method:  notif.Method(),
	}

	if notif.Params() != nil {
		jsonObj.Params = notif.Params()
	}

	return json.Marshal(jsonObj)
}
