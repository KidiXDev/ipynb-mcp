package tools

import (
	"context"

	"github.com/kidixdev/ipynb-mcp/internal/notebook"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterUpdateCodeCell(s *server.MCPServer) {
	tool := mcp.NewTool(
		"update_code_cell",
		mcp.WithDescription("Replace an existing cell with a code cell at a given index."),
		mcp.WithString("path", mcp.Required(), mcp.Description("Path to the .ipynb file.")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("Cell index to replace.")),
		mcp.WithString("source", mcp.Required(), mcp.Description("New code source.")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, err := req.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		index, err := req.RequireInt("index")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		source, err := req.RequireString("source")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return withNotebookMutation(path, func(nb *notebook.Notebook) error {
			return notebook.UpdateCodeCell(nb, index, source)
		})
	})
}
