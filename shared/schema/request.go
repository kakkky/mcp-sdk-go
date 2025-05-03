package schema

type Request struct {
	Method string `json:"method"`
	Params any    `json:"params"`
}
