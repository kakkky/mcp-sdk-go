package schema

type Notification struct {
	Method string `json:"method"`
	Params any    `json:"params"`
}
