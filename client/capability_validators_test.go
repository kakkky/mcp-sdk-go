package client

import (
	"strings"
	"testing"

	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func TestValidateCapabilityForMethod(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		capabilities  schema.ServerCapabilities
		expectedError bool
		errorContains string
	}{
		{
			name:   "normal: server supports logging",
			method: "logging/setLevel",
			capabilities: schema.ServerCapabilities{
				Logging: &schema.Logging{},
			},
			expectedError: false,
		},
		{
			name:          "semi normal: server does not support logging",
			method:        "logging/setLevel",
			capabilities:  schema.ServerCapabilities{},
			expectedError: true,
			errorContains: "client does not support logging",
		},
		{
			name:   "normal: server supports prompts",
			method: "prompts/list",
			capabilities: schema.ServerCapabilities{
				Prompts: &schema.Prompts{},
			},
			expectedError: false,
		},
		{
			name:          "semi normal: server does not support prompts",
			method:        "prompts/list",
			capabilities:  schema.ServerCapabilities{},
			expectedError: true,
			errorContains: "client does not support prompts",
		},
		{
			name:   "normal: server supports prompts for get",
			method: "prompts/get",
			capabilities: schema.ServerCapabilities{
				Prompts: &schema.Prompts{},
			},
			expectedError: false,
		},
		{
			name:          "semi normal: server does not support prompts for get",
			method:        "prompts/get",
			capabilities:  schema.ServerCapabilities{},
			expectedError: true,
			errorContains: "client does not support prompts",
		},
		{
			name:   "normal: server supports resources",
			method: "resources/list",
			capabilities: schema.ServerCapabilities{
				Resources: &schema.Resources{},
			},
			expectedError: false,
		},
		{
			name:          "semi normal: server does not support resources",
			method:        "resources/list",
			capabilities:  schema.ServerCapabilities{},
			expectedError: true,
			errorContains: "client does not support resources",
		},
		{
			name:   "normal: server supports resources/read",
			method: "resources/read",
			capabilities: schema.ServerCapabilities{
				Resources: &schema.Resources{},
			},
			expectedError: false,
		},
		{
			name:          "semi normal: server does not support resources/read",
			method:        "resources/read",
			capabilities:  schema.ServerCapabilities{},
			expectedError: true,
			errorContains: "client does not support resources",
		},
		{
			name:   "normal: server supports resources/templates/list",
			method: "resources/templates/list",
			capabilities: schema.ServerCapabilities{
				Resources: &schema.Resources{},
			},
			expectedError: false,
		},
		{
			name:          "semi normal: server does not support resources/templates/list",
			method:        "resources/templates/list",
			capabilities:  schema.ServerCapabilities{},
			expectedError: true,
			errorContains: "client does not support resources",
		},
		{
			name:   "normal: server supports resources/subscribe",
			method: "resources/subscribe",
			capabilities: schema.ServerCapabilities{
				Resources: &schema.Resources{
					Subscribe: true,
				},
			},
			expectedError: false,
		},
		{
			name:   "semi normal: server supports resources but not subscribe",
			method: "resources/subscribe",
			capabilities: schema.ServerCapabilities{
				Resources: &schema.Resources{
					Subscribe: false,
				},
			},
			expectedError: true,
			errorContains: "client does not support subscribing to resources",
		},
		{
			name:          "semi normal: server does not support resources/subscribe",
			method:        "resources/subscribe",
			capabilities:  schema.ServerCapabilities{},
			expectedError: true,
			errorContains: "client does not support resources",
		},
		{
			name:   "normal: server supports resources/unsubscribe",
			method: "resources/unsubscribe",
			capabilities: schema.ServerCapabilities{
				Resources: &schema.Resources{},
			},
			expectedError: false,
		},
		{
			name:          "semi normal: server does not support resources/unsubscribe",
			method:        "resources/unsubscribe",
			capabilities:  schema.ServerCapabilities{},
			expectedError: true,
			errorContains: "client does not support resources",
		},
		{
			name:   "normal: server supports tools",
			method: "tools/list",
			capabilities: schema.ServerCapabilities{
				Tools: &schema.Tools{},
			},
			expectedError: false,
		},
		{
			name:          "semi normal: server does not support tools",
			method:        "tools/list",
			capabilities:  schema.ServerCapabilities{},
			expectedError: true,
			errorContains: "client does not support tools",
		},
		{
			name:   "normal: server supports tools/call",
			method: "tools/call",
			capabilities: schema.ServerCapabilities{
				Tools: &schema.Tools{},
			},
			expectedError: false,
		},
		{
			name:          "semi normal: server does not support tools/call",
			method:        "tools/call",
			capabilities:  schema.ServerCapabilities{},
			expectedError: true,
			errorContains: "client does not support tools",
		},
		{
			name:   "normal: server supports completion",
			method: "completion/complete",
			capabilities: schema.ServerCapabilities{
				Completion: &schema.Completion{},
			},
			expectedError: false,
		},
		{
			name:          "semi normal: server does not support completion",
			method:        "completion/complete",
			capabilities:  schema.ServerCapabilities{},
			expectedError: true,
			errorContains: "client does not support completion",
		},
		{
			name:          "normal: initialize is always supported",
			method:        "initialize",
			capabilities:  schema.ServerCapabilities{},
			expectedError: false,
		},
		{
			name:          "normal: ping is always supported",
			method:        "ping",
			capabilities:  schema.ServerCapabilities{},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := NewClient(schema.Implementation{}, nil)
			sut.serverCapabilities = tt.capabilities

			err := sut.validateCapabilityForMethod(tt.method)

			if (err != nil) != tt.expectedError {
				t.Errorf("validateCapabilityForMethod() error = %v, expectedError %v", err, tt.expectedError)
			}

			if err != nil && tt.errorContains != "" {
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("validateCapabilityForMethod() error = %v, should contain %v", err, tt.errorContains)
				}
			}
		})
	}
}

func TestValidateNotificationCapability(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		capabilities  schema.ClientCapabilities
		expectedError bool
		errorContains string
	}{
		{
			name:   "normal: client supports roots list_changed notification",
			method: "notifications/roots/list_changed",
			capabilities: schema.ClientCapabilities{
				Roots: &schema.Roots{
					ListChanged: true,
				},
			},
			expectedError: false,
		},
		{
			name:   "semi normal: client does not support roots list_changed notification",
			method: "notifications/roots/list_changed",
			capabilities: schema.ClientCapabilities{
				Roots: &schema.Roots{
					ListChanged: false,
				},
			},
			expectedError: true,
			errorContains: "Client does not support notifying about roots",
		},
		{
			name:          "normal: initialized notification is always supported",
			method:        "notifications/initialized",
			capabilities:  schema.ClientCapabilities{},
			expectedError: false,
		},
		{
			name:          "normal: cancelled notification is always supported",
			method:        "notifications/cancelled",
			capabilities:  schema.ClientCapabilities{},
			expectedError: false,
		},
		{
			name:          "normal: progress notification is always supported",
			method:        "notifications/progress",
			capabilities:  schema.ClientCapabilities{},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := NewClient(schema.Implementation{}, &ClientOptions{
				Capabilities: tt.capabilities,
			})

			err := sut.validateNotificationCapability(tt.method)

			if (err != nil) != tt.expectedError {
				t.Errorf("validateNotificationCapability() error = %v, expectedError %v", err, tt.expectedError)
			}

			if err != nil && tt.errorContains != "" {
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("validateNotificationCapability() error = %v, should contain %v", err, tt.errorContains)
				}
			}
		})
	}
}

func TestValidateRequestHandlerCapability(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		capabilities  schema.ClientCapabilities
		expectedError bool
		errorContains string
	}{
		{
			name:   "normal: client supports sampling",
			method: "sampling/createMessage",
			capabilities: schema.ClientCapabilities{
				Sampling: &schema.Sampling{},
			},
			expectedError: false,
		},
		{
			name:          "semi normal: client does not support sampling",
			method:        "sampling/createMessage",
			capabilities:  schema.ClientCapabilities{},
			expectedError: true,
			errorContains: "Client does not support sampling",
		},
		{
			name:   "normal: client supports roots",
			method: "roots/list",
			capabilities: schema.ClientCapabilities{
				Roots: &schema.Roots{},
			},
			expectedError: false,
		},
		{
			name:          "semi normal: client does not support roots",
			method:        "roots/list",
			capabilities:  schema.ClientCapabilities{},
			expectedError: true,
			errorContains: "Client does not support roots",
		},
		{
			name:          "normal: ping is always supported",
			method:        "ping",
			capabilities:  schema.ClientCapabilities{},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := NewClient(schema.Implementation{}, &ClientOptions{
				Capabilities: tt.capabilities,
			})

			err := sut.validateRequestHandlerCapability(tt.method)

			if (err != nil) != tt.expectedError {
				t.Errorf("validateRequestHandlerCapability() error = %v, expectedError %v", err, tt.expectedError)
			}

			if err != nil && tt.errorContains != "" {
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("validateRequestHandlerCapability() error = %v, should contain %v", err, tt.errorContains)
				}
			}
		})
	}
}
