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
	case isListRootsResult(rawResult):
		var result schema.ListRootsResultSchema
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

	case isCompleteResult(rawResult):
		var result schema.CompleteResultSchema
		if err := json.Unmarshal(message.Result, &result); err != nil {
			return nil, err
		}
		return &result, nil
	case isListToolsResult(rawResult):
		var result schema.ListToolsResultSchema
		if err := json.Unmarshal(message.Result, &result); err != nil {
			return nil, err
		}
		return &result, nil
	case isCallToolResult(rawResult):
		rawContents := rawResult["content"].([]any)
		contents := make([]schema.ToolContentSchema, 0, len(rawContents))
		for _, content := range rawContents {
			content := content.(map[string]any)
			switch content["type"].(string) {
			case "text":
				contents = append(contents, &schema.TextContentSchema{
					Type: content["type"].(string),
					Text: content["text"].(string),
				})
			case "image":
				contents = append(contents, &schema.ImageContentSchema{
					Type:     content["type"].(string),
					Data:     content["data"].(string),
					MimeType: content["mimeType"].(string),
				})
			case "audio":
				contents = append(contents, &schema.AudioContentSchema{
					Type:     content["type"].(string),
					Data:     content["data"].(string),
					MimeType: content["mimeType"].(string),
				})
			}
		}
		isError, found := rawResult["isError"]
		if !found {
			isError = false // デフォルトはエラーではない
		}
		return &schema.CallToolResultSchema{
			Content: contents,
			IsError: isError.(bool),
		}, nil
	case isListPromptsResult(rawResult):
		var result schema.ListPromptsResultSchema
		if err := json.Unmarshal(message.Result, &result); err != nil {
			return nil, err
		}
		return &result, nil
	case isGetPromptResult(rawResult):
		var result schema.GetPromptResultSchema
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
func isListRootsResult(data map[string]any) bool {
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
func isListToolsResult(data map[string]any) bool {
	return hasResultFields(data, "tools")
}
func isCallToolResult(data map[string]any) bool {
	return hasResultFields(data, "content")
}
func isListPromptsResult(data map[string]any) bool {
	return hasResultFields(data, "prompts")
}
func isGetPromptResult(data map[string]any) bool {
	return hasResultFields(data, "description", "messages")
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
