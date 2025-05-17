package mcpserver

func (m *McpServer) sendResourceListChanged() {
	if m.isConnected() {
		m.server.SendResourceListChanged()
	}
}
