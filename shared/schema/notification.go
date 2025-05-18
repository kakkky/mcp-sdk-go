package schema

import "fmt"

type Notification interface {
	Method() string
	Params() any
}

type InitializeNotificationSchema struct {
	MethodName string `json:"method"`
}

func (n *InitializeNotificationSchema) Method() string {
	if n.MethodName != "notifications/initialized" {
		fmt.Println("Invalid method name")
	}
	return n.MethodName
}

func (n *InitializeNotificationSchema) Params() any {
	return nil
}

type LoggingMessageNotificationSchema struct {
	MethodName string                           `json:"method"`
	ParamsData LoggingMessageNotificationParams `json:"params"`
}

type LoggingMessageNotificationParams struct {
	Level  LoggingLevelSchema `json:"level"`
	Logger string             `json:"logger,omitempty"`
	Data   any                `json:"data"`
}

func (n *LoggingMessageNotificationSchema) Method() string {
	return n.MethodName
}

func (n *LoggingMessageNotificationSchema) Params() any {
	return n.ParamsData
}

type ResourceUpdatedNotificationSchema struct {
	MethodName string                            `json:"method"`
	ParamsData ResourceUpdatedNotificationParams `json:"params"`
}

type ResourceUpdatedNotificationParams struct {
	Uri string `json:"uri"`
}

func (n *ResourceUpdatedNotificationSchema) Method() string {
	return n.MethodName
}

func (n *ResourceUpdatedNotificationSchema) Params() any {
	return n.ParamsData
}

type ResourceListChangedNotificationSchema struct {
	MethodName string `json:"method"`
}

func (n *ResourceListChangedNotificationSchema) Method() string {
	return n.MethodName
}
func (n *ResourceListChangedNotificationSchema) Params() any {
	return nil
}

type ToolListChangedNotificationSchema struct {
	MethodName string `json:"method"`
}

func (n *ToolListChangedNotificationSchema) Method() string {
	return n.MethodName
}

func (n *ToolListChangedNotificationSchema) Params() any {
	return nil
}

type PromptListChangedNotificationSchema struct {
	MethodName string `json:"method"`
}

func (n *PromptListChangedNotificationSchema) Method() string {
	return n.MethodName
}

func (n *PromptListChangedNotificationSchema) Params() any {
	return nil
}
