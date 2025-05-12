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
	if n.MethodName != "notification/initialize" {
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
	Logger *string            `json:"logger,omitempty"`
	Data   any                `json:"data"`
}

func (n *LoggingMessageNotificationSchema) Method() string {
	return n.MethodName
}

func (n *LoggingMessageNotificationSchema) Params() any {
	return n.ParamsData
}
