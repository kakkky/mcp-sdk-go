package client

import "fmt"

// リクエストを送信する際に、メソッドにサーバーが対応しているのかを検証する
func (s *Client) validateCapabilityForMethod(method string) error {
	switch method {
	case "logging/setLevel":
		if s.serverCapabilities.Logging == nil {
			return fmt.Errorf("client does not support logging (required for %s)", method)
		}
	case "prompts/get",
		"prompts/list":
		if s.serverCapabilities.Prompts == nil {
			return fmt.Errorf("client does not support prompts (required for %s)", method)
		}
	case "resources/list",
		"resources/templates/list",
		"resources/read",
		"resources/subscribe",
		"resources/unsubscribe":
		if s.serverCapabilities.Resources == nil {
			return fmt.Errorf("client does not support resources (required for %s)", method)
		}
		resources := s.serverCapabilities.Resources
		if (method == "resources/subscribe") && !resources.Subscribe {
			return fmt.Errorf("client does not support subscribing to resources (required for %s)", method)
		}
	case "tools/call",
		"tools/list":
		if s.serverCapabilities.Tools == nil {
			return fmt.Errorf("client does not support tools (required for %s)", method)
		}
	case "completion/complete":
		if s.serverCapabilities.Completion == nil {
			return fmt.Errorf("client does not support completion (required for %s)", method)
		}
	case "initialize":
		break
	case "ping":
		break
	}
	return nil
}

// 通知を送る前に、クライアントがメソッドに対応しているのかを検証する
func (s *Client) validateNotificationCapability(method string) error {
	switch method {
	case "notifications/roots/list_changed":
		roots := s.capabilities.Roots
		if !roots.ListChanged {
			return fmt.Errorf("Client does not support notifying about roots (required for %s)", method)
		}
	case "notifications/initialized":
		break
	case "notifications/cancelled":
		break
	case "notifications/progress":
		break
	}
	return nil
}

// リクエストハンドラを登録する前に、クライアントがメソッドに対応しているのかを検証する
func (s *Client) validateRequestHandlerCapability(method string) error {
	switch method {
	case "sampling/createMessage":
		if s.capabilities.Sampling == nil {
			return fmt.Errorf("Client does not support sampling (required for %s)", method)
		}
	case "roots/list":
		if s.capabilities.Roots == nil {
			return fmt.Errorf("Client does not support roots (required for %s)", method)
		}
	case "ping":
		break
	}
	return nil
}
