package tools

import "github.com/mark3labs/mcp-go/server"

func Register(s *server.MCPServer) {
	RegisterReadNotebook(s)
	RegisterReadCellOutput(s)
	RegisterCreateNotebook(s)
	RegisterInsertMarkdownCell(s)
	RegisterInsertCodeCell(s)
	RegisterUpdateMarkdownCell(s)
	RegisterUpdateCodeCell(s)
	RegisterDeleteCell(s)
}
