package mcpserver

func (m *McpServer) sendResourceListChanged() {
	if m.isConnected() {
		_ = m.Server.SendResourceListChanged()
	}
}

func (m *McpServer) sendToolListChanged() {
	if m.isConnected() {
		_ = m.Server.SendToolListChanged()
	}
}
