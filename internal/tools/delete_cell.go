package tools

import (
	"context"

	"github.com/kidixdev/ipynb-mcp/internal/notebook"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterDeleteCell(s *server.MCPServer) {
	tool := mcp.NewTool(
		"delete_cell",
		mcp.WithDescription("Delete a cell at a given index."),
		mcp.WithString("path", mcp.Required(), mcp.Description("Path to the .ipynb file.")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("Cell index to delete.")),
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

		return withNotebookMutation(path, func(nb *notebook.Notebook) error {
			return notebook.DeleteCell(nb, index)
		})
	})
}
