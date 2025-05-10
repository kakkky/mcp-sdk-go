package server

import (
	"testing"

	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func TestAssertCapabilityForMethod(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		capabilities  schema.ClientCapabilities
		expectedError bool
	}{
		{
			name:   "normal : client supports sampling",
			method: "sampling/createMessage",
			capabilities: schema.ClientCapabilities{
				Sampling: &schema.Sampling{},
			},
			expectedError: false,
		},
		{
			name:          "semi normal : client does not support sampling",
			method:        "sampling/createMessage",
			capabilities:  schema.ClientCapabilities{},
			expectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := NewServer(schema.Implementation{}, nil)
			sut.clientCapabilities = tt.capabilities

			err := sut.assertCapabilityForMethod(tt.method)

			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}
		})
	}
}

func TestAssertNotificationCapability(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		capabilities  schema.ServerCapabilities
		expectedError bool
	}{
		{
			name:   "normal : server supports logging",
			method: "notifications/message",
			capabilities: schema.ServerCapabilities{
				Logging: &schema.Logging{},
			},
			expectedError: false,
		},
		{
			name:          "semi normal : server does not support logging",
			method:        "notifications/message",
			capabilities:  schema.ServerCapabilities{},
			expectedError: true,
		},
		{
			name:          "semi normal : server does not support resources",
			method:        "notifications/resources/updated",
			capabilities:  schema.ServerCapabilities{},
			expectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := NewServer(schema.Implementation{}, nil)
			sut.capabilities = tt.capabilities

			err := sut.assertNotificationCapability(tt.method)

			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}
		})
	}
}
