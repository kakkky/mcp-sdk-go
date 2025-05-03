package capability

type Capability struct {
	tool     []tool
	prompt   []prompt
	resource []resource
}

type tool struct{}

type prompt struct{}

type resource struct{}
