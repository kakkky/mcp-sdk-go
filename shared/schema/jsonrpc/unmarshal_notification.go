package jsonrpc

import (
	"encoding/json"
	"fmt"

	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func unmarshalNotification(message *Message) (schema.Notification, error) {
	// メソッド名から判断して、適切な通知型を選択
	switch message.Method {
	case "notifications/initialized":
		return &schema.InitializeNotificationSchema{
			MethodName: message.Method,
		}, nil

	case "notifications/message":
		params := schema.LoggingMessageNotificationParams{}
		if err := json.Unmarshal(message.Params, &params); err != nil {
			return nil, err
		}
		return &schema.LoggingMessageNotificationSchema{
			MethodName: message.Method,
			ParamsData: params,
		}, nil

	case "notifications/resources/updated":
		params := schema.ResourceUpdatedNotificationParams{}
		if err := json.Unmarshal(message.Params, &params); err != nil {
			return nil, err
		}
		return &schema.ResourceUpdatedNotificationSchema{
			MethodName: message.Method,
			ParamsData: params,
		}, nil

	case "notifications/resources/list_changed":
		return &schema.ResourceListChangedNotificationSchema{
			MethodName: message.Method,
		}, nil

	case "notifications/tools/list_changed":
		return &schema.ToolListChangedNotificationSchema{
			MethodName: message.Method,
		}, nil

	case "notifications/prompts/list_changed":
		return &schema.PromptListChangedNotificationSchema{
			MethodName: message.Method,
		}, nil

	// その他の通知タイプはここに追加

	default:
		return nil, fmt.Errorf("unknown notification method: %s", message.Method)
	}
}
