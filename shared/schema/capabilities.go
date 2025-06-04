package schema

// Server
type ServerCapabilities struct {
	Experimental any `json:"experimental,omitempty"`
	*Logging     `json:"logging,omitempty"`
	*Prompts     `json:"prompts,omitempty"`
	*Resources   `json:"resources,omitempty"`
	*Tools       `json:"tools,omitempty"`
}

type Logging struct{}

type Prompts struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type Resources struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	ListChanged bool `json:"listChanged,omitempty"`
}
type Tools struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// Client
type ClientCapabilities struct {
	Experimental any `json:"experimental,omitempty"`
	*Sampling    `json:"sampling,omitempty"`
	*Roots       `json:"roots,omitempty"`
}

type Sampling struct{}

type Roots struct {
	ListChanged bool `json:"listChanged,omitempty"`
}
