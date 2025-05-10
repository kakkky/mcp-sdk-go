package schema

type Request interface {
	Method() string
	Params() any
}

type InitializeRequestSchema struct {
	MethodName string                  `json:"method"`
	ParamsData InitializeRequestParams `json:"params"`
}

type InitializeRequestParams struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ClientCapabilities `json:"capabilities"`
	ClientInfo      Implementation     `json:"clientInfo"`
}

func (r *InitializeRequestSchema) Method() string {
	return r.MethodName
}

func (r *InitializeRequestSchema) Params() any {
	return r.ParamsData
}
