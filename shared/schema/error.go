package schema

import mcperr "github.com/kakkky/mcp-sdk-go/shared/mcp-err"

type Error struct {
	Code    mcperr.ErrCode `json:"code"`
	Message string         `json:"message"`
	Data    any            `json:"data,omitempty"`
}
