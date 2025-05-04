package schema

import "github.com/kakkky/mcp-sdk-go/shared/mcp_err"

type Error struct {
	Code    mcp_err.ErrCode `json:"code"`
	Message string          `json:"message"`
	Data    any             `json:"data,omitempty"`
}
