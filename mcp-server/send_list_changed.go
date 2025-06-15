package mcpserver

func (m *McpServer) sendResourceListChanged() {
	if m.isConnected() {
		_ = m.Server.SendResourceListChanged()
	}
}
