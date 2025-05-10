package schema

type Result interface {
	Result() any
}

type InitializeResultSchema struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      Implementation     `json:"serverInfo"`
	Instructions    *string            `json:"instructions,omitempty"`
}

func (r *InitializeResultSchema) Result() any {
	return r
}

type EmptyResultSchema struct{}

func (r *EmptyResultSchema) Result() any {
	return r
}
