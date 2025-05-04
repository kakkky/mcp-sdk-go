package mcp_err

import (
	"fmt"
)

type McpErr struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
	Data    any     `json:"data,omitempty"`
}

func NewMcpErr(code ErrCode, message string, data any) *McpErr {
	return &McpErr{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func (e *McpErr) Error() string {
	return fmt.Sprintf("MCP Error: %d (%s)", e.Code, e.Message)
}

type ErrCode int

const (
	// Custom error codes
	CONNECTION_CLOSED = -32000
	REQUEST_TIMEOUT   = -32001

	// Standard JSON-RPC error codes
	PARSE_ERROR      = -32700
	INVALID_REQUEST  = -32600
	METHOD_NOT_FOUND = -32601
	INVALID_PARAMS   = -32602
	INTERNAL_ERROR   = -32603
)
