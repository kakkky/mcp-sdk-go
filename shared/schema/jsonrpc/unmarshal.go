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

	// Request
	if message.Method != "" && message.Id != 0 {
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
	}

	// Notification
	if message.Method != "" && message.Id == 0 {

	}
	// Error
	if message.Error != nil {

	}
	// Response
	return nil, nil
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
