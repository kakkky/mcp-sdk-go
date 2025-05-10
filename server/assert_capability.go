package server

import "fmt"

// リクエストを送信する際に、メソッドにクライアントが対応しているのかを検証する
func (s *Server) assertCapabilityForMethod(method string) error {
	switch method {
	case "sampling/createMessage":
		if s.clientCapabilities.Sampling == nil {
			return fmt.Errorf("client does not support sampling (required for %s)", method)
		}
	case "roots/list":
		if s.clientCapabilities.Roots == nil {
			return fmt.Errorf("client does not support roots (required for %s)", method)
		}
	case "ping":
		break
	}
	return nil
}

// 通知を送る前に、サーバーがメソッドに対応しているのかを検証する
func (s *Server) assertNotificationCapability(method string) error {
	switch method {
	case "notifications/message":
		if s.capabilities.Logging == nil {
			return fmt.Errorf("server does not support logging (required for %s)", method)
		}
	case "notifications/resources/updated",
		"notifications/resources/list_changed":
		if s.capabilities.Resources == nil {
			return fmt.Errorf("server does not support notifying about resources (required for %s)", method)
		}
	case "notifications/tools/list_changed":
		if s.capabilities.Tools == nil {
			return fmt.Errorf("server does not support notifying about tools (required for %s)", method)
		}
	case "notifications/prompts/list_changed":
		if s.capabilities.Prompts == nil {
			return fmt.Errorf("server does not support notifying about prompts (required for %s)", method)
		}
	case "notifications/cancelled":
		break
	case "notifications/progress":
		break
	}
	return nil
}
