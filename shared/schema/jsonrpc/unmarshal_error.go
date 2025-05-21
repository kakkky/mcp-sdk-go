package jsonrpc

import (
	"encoding/json"

	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func unmarshalError(message *Message) (*schema.Error, error) {
	errorData := schema.Error{
		Code:    message.Error.Code,
		Message: message.Error.Message,
	}

	if message.Error.Data != nil {
		errorData.Data = message.Error.Data
	}
	if message.Error.Data != nil {
		err := json.Unmarshal(message.Error.Data, &errorData.Data)
		if err != nil {
			return nil, err
		}
	}
	return &errorData, nil
}
