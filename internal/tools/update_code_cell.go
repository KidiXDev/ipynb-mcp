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
		mcp.WithDescription("Replace one or more existing cells with code cells at given indices."),
		mcp.WithString("path", mcp.Required(), mcp.Description("Path to the .ipynb file.")),
		mcp.WithNumber("index", mcp.Description("Cell index to replace (single update mode).")),
		mcp.WithString("source", mcp.Description("New code source (single update mode).")),
		mcp.WithArray("indices", mcp.Description("Cell indices to replace (batch mode)."), mcp.WithNumberItems()),
		mcp.WithArray("sources", mcp.Description("New code sources matching the indices order (batch mode)."), mcp.WithStringItems()),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, err := req.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		indices, sources, err := parseUpdateCellsArgs(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return withNotebookMutation(path, func(nb *notebook.Notebook) error {
			for i := range indices {
				if err := notebook.UpdateCodeCell(nb, indices[i], sources[i]); err != nil {
					return err
				}
			}
			return nil
		})
	})
}
