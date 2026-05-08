package tools

import (
	"context"

	"github.com/kidixdev/ipynb-mcp/internal/notebook"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterInsertCodeCell(s *server.MCPServer) {
	tool := mcp.NewTool(
		"insert_code_cell",
		mcp.WithDescription("Insert a code cell at a given index."),
		mcp.WithString("path", mcp.Required(), mcp.Description("Path to the .ipynb file.")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("Cell index to insert at.")),
		mcp.WithString("source", mcp.Required(), mcp.Description("Code source for the new cell.")),
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
			return notebook.InsertCodeCell(nb, index, source)
		})
	})
}
