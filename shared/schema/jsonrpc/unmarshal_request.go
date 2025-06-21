package jsonrpc

import (
	"encoding/json"
	"fmt"

	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

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
	case "roots/list":
		return &schema.ListRootsRequestSchema{
			MethodName: message.Method,
		}, nil
	case "logging/setLevel":
		params := &schema.SetLoggingLevelRequestParams{}
		if err := json.Unmarshal(message.Params, &params); err != nil {
			return nil, err
		}
		return &schema.SetLevelRequestSchema{
			MethodName: message.Method,
			ParamsData: *params,
		}, nil
	case "tools/list":
		return &schema.ListToolsRequestSchema{
			MethodName: message.Method,
		}, nil
	case "tools/call":
		params := &schema.CallToolRequestParams{}
		if err := json.Unmarshal(message.Params, &params); err != nil {
			return nil, err
		}
		return &schema.CallToolRequestSchema{
			MethodName: message.Method,
			ParamsData: *params,
		}, nil
	}
	return nil, fmt.Errorf("unknown method: %s", message.Method)
}
