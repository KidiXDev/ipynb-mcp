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
		mcp.WithDescription("Delete one or more cells by index."),
		mcp.WithString("path", mcp.Required(), mcp.Description("Path to the .ipynb file.")),
		mcp.WithNumber("index", mcp.Description("Cell index to delete (single delete mode).")),
		mcp.WithArray("indices", mcp.Description("Cell indices to delete (batch mode)."), mcp.WithNumberItems()),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, err := req.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		indices, err := parseDeleteCellIndices(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return withNotebookMutation(path, func(nb *notebook.Notebook) error {
			for i := range indices {
				if err := notebook.DeleteCell(nb, indices[i]); err != nil {
					return err
				}
			}
			return nil
		})
	})
}
