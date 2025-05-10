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
