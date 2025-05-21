package jsonrpc

import (
	"encoding/json"
	"fmt"

	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func unmarshalResult(message *Message) (schema.Result, error) {
	var rawResult map[string]any
	if err := json.Unmarshal(message.Result, &rawResult); err != nil {
		return nil, err
	}
	switch {
	case isInitializeResult(rawResult):
		var result schema.InitializeResultSchema
		if err := json.Unmarshal(message.Result, &result); err != nil {
			return nil, err
		}
		return &result, nil
	case isEmptyResult(rawResult):
		return &schema.EmptyResultSchema{}, nil
	case isCreateMessageResult(rawResult):
		result := struct {
			Model      string         `json:"model"`
			Role       string         `json:"role"`
			StopReason string         `json:"stopReason,omitempty"`
			Content    map[string]any `json:"content"`
		}{}
		if err := json.Unmarshal(message.Result, &result); err != nil {
			return nil, err
		}
		contentType := result.Content["type"].(string)
		// contentTypeに応じて、適切な構造体に変換
		switch contentType {
		case "text":
			return &schema.CreateMessageResultSchema[schema.TextContentSchema]{
				Model:      result.Model,
				Role:       result.Role,
				StopReason: result.StopReason,
				Content: schema.TextContentSchema{
					Type: result.Content["type"].(string),
					Text: result.Content["text"].(string),
				},
			}, nil
		case "image":
			return &schema.CreateMessageResultSchema[schema.ImageContentSchema]{
				Model:      result.Model,
				Role:       result.Role,
				StopReason: result.StopReason,
				Content: schema.ImageContentSchema{
					Type:     result.Content["type"].(string),
					Data:     result.Content["data"].(string),
					MimeType: result.Content["mimeType"].(string),
				},
			}, nil
		case "audio":
			return &schema.CreateMessageResultSchema[schema.AudioContentSchema]{
				Model:      result.Model,
				Role:       result.Role,
				StopReason: result.StopReason,
				Content: schema.AudioContentSchema{
					Type:     result.Content["type"].(string),
					Data:     result.Content["data"].(string),
					MimeType: result.Content["mimeType"].(string),
				},
			}, nil
		default:
			return nil, fmt.Errorf("unknown content type: %s", contentType)
		}
	case isReadResourceResult(rawResult):
		rawContents := rawResult["contents"].([]any)
		contents := make([]schema.ResourceContentSchema, 0, len(rawContents))
		for _, content := range rawContents {
			content := content.(map[string]any)
			if _, ok := content["text"]; ok {
				contents = append(contents, &schema.TextResourceContentsSchema{
					UriData:      content["uri"].(string),
					MimeTypeData: content["mimeType"].(string),
					ContentData:  content["text"].(string),
				})
			}
			if _, ok := content["blob"]; ok {
				contents = append(contents, &schema.BlobResourceContentsSchema{
					UriData:      content["uri"].(string),
					MimeTypeData: content["mimeType"].(string),
					ContentData:  content["blob"].(string),
				})
			}
		}
		return &schema.ReadResourceResultSchema{
			Contents: contents,
		}, nil
	case isListRootResult(rawResult):
		var result schema.ListRootResultSchema
		if err := json.Unmarshal(message.Result, &result); err != nil {
			return nil, err
		}
		return &result, nil
	case isListResourcesResult(rawResult):
		var result schema.ListResourcesResultSchema
		if err := json.Unmarshal(message.Result, &result); err != nil {
			return nil, err
		}
		return &result, nil

	case isListResourceTemplatesResult(rawResult):
		var result schema.ListResourceTemplatesResultSchema
		if err := json.Unmarshal(message.Result, &result); err != nil {
			return nil, err
		}
		return &result, nil

	case isListRootResult(rawResult):
		var result schema.ListRootResultSchema
		if err := json.Unmarshal(message.Result, &result); err != nil {
			return nil, err
		}
		return &result, nil

	case isCompleteResult(rawResult):
		var result schema.CompleteResultSchema
		if err := json.Unmarshal(message.Result, &result); err != nil {
			return nil, err
		}
		return &result, nil
	}
	return nil, fmt.Errorf("unknown result: %v", rawResult)
}

// 各種Resultの判定関数
func isInitializeResult(data map[string]any) bool {
	return hasResultFields(data, "protocolVersion", "serverInfo", "capabilities")
}
func isEmptyResult(data map[string]any) bool {
	return len(data) == 0
}
func isCreateMessageResult(data map[string]any) bool {
	return hasResultFields(data, "model", "role", "content")
}
func isListRootResult(data map[string]any) bool {
	return hasResultFields(data, "roots")
}
func isReadResourceResult(data map[string]any) bool {
	return hasResultFields(data, "contents")
}
func isListResourcesResult(data map[string]any) bool {
	return hasResultFields(data, "resources")
}
func isListResourceTemplatesResult(data map[string]any) bool {
	return hasResultFields(data, "resourceTemplates")
}
func isCompleteResult(data map[string]any) bool {
	return hasResultFields(data, "completion")
}

// 指定されたフィールドがすべて存在するか確認
// 上の関数の共通処理
func hasResultFields(data map[string]interface{}, fields ...string) bool {
	for _, field := range fields {
		if _, ok := data[field]; !ok {
			return false
		}
	}
	return true
}
