package protocol

import "github.com/kakkky/mcp-sdk-go/shared/schema"

func (p *Protocol) onResponse(message schema.JsonRpcMessage) {}

func (p *Protocol) onRequest(message schema.JsonRpcMessage) {}

func (p *Protocol) onNotification(notification *schema.JsonRpcNotification) {
	handler := p.handlers.notificationHandlers[notification.Method]
	if handler == nil && p.handlers.fallbackNotificationHandler != nil {
		p.handlers.fallbackNotificationHandler()
		return
	}
	if err := handler(notification); err != nil {
		p.onError(err)
	}
}
