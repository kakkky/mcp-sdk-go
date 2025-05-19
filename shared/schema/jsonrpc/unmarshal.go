package json

import (
	"encoding/json"
	"fmt"

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

	// Error
	case message.Error != nil:

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
	return nil, fmt.Errorf("unknown message: %v", message)
}

func unmarshalRequest(message *Message) (schema.Request, error) {
	// Methodから判断して、RequestフィールドをUnmarshalする
	switch message.Method {
	case "initialize":
		params := &schema.InitializeRequestParams{}
		if err := json.Unmarshal(message.Params, &params); err != nil {
			return nil, err
		}
		return &schema.InitializeRequestSchema{
			MethodName: message.Method,
			ParamsData: *params,
		}, nil
	case "ping":
		return &schema.PingRequestSchema{
			MethodName: message.Method,
		}, nil
	case "resources/list":
		return &schema.ListResourceRequestSchema{
			MethodName: message.Method,
		}, nil
	case "resources/read":
		params := &schema.ReadResourceRequestParams{}
		if err := json.Unmarshal(message.Params, &params); err != nil {
			return nil, err
		}
		return &schema.ReadResourceRequestSchema{
			MethodName: message.Method,
			ParamsData: *params,
		}, nil
	case "resources/templates/list":
		return &schema.ListResourceTemplatesRequestSchema{
			MethodName: message.Method,
		}, nil
	case "completion/complete":
		params := struct {
			Ref      json.RawMessage                 `json:"ref"`
			Argument schema.CompleteRequestParamsArg `json:"argument"`
		}{}
		if err := json.Unmarshal(message.Params, &params); err != nil {
			return nil, err
		}
		ref := struct {
			Type string `json:"type"`
			Uri  string `json:"uri"`
			Name string `json:"name"`
		}{}
		if err := json.Unmarshal(params.Ref, &ref); err != nil {
			return nil, err
		}
		var request *schema.CompleteRequestSchema
		switch ref.Type {
		case "ref/resource":
			request = &schema.CompleteRequestSchema{
				MethodName: message.Method,
				ParamsData: schema.CompleteRequestParams{
					Ref: &schema.ResourceReferenceSchema{
						TypeData: ref.Type,
						UriData:  ref.Uri,
					},
					Argument: params.Argument,
				},
			}
		case "ref/prompt":
			request = &schema.CompleteRequestSchema{
				MethodName: message.Method,
				ParamsData: schema.CompleteRequestParams{
					Ref: &schema.PromptReferenceSchema{
						TypeData: ref.Type,
						NameData: ref.Name,
					},
					Argument: params.Argument,
				},
			}
		default:
			return nil, fmt.Errorf("unknown type: %s", ref.Type)
		}
		return request, nil
	}
	return nil, fmt.Errorf("unknown method: %s", message.Method)
}

func unmarshalResult(message *Message) (schema.Result, error) {
	var rawResult map[string]any
	if err := json.Unmarshal(message.Result, &rawResult); err != nil {
		return nil, err
	}
	switch {
	case hasResultFields(rawResult, "protocolVersion", "serverInfo", "capabilities"):
		var result schema.InitializeResultSchema
		if err := json.Unmarshal(message.Result, &result); err != nil {
			return nil, err
		}
		return &result, nil
	}
	return nil, fmt.Errorf("unknown result: %v", rawResult)
}

// 指定されたフィールドがすべて存在するか確認
func hasResultFields(data map[string]interface{}, fields ...string) bool {
	for _, field := range fields {
		if _, ok := data[field]; !ok {
			return false
		}
	}
	return true
}
