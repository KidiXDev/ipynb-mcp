package tools

import (
	"context"

	"github.com/kidixdev/ipynb-mcp/internal/notebook"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterInsertMarkdownCell(s *server.MCPServer) {
	tool := mcp.NewTool(
		"insert_markdown_cell",
		mcp.WithDescription("Insert one or more markdown cells at given indices."),
		mcp.WithString("path", mcp.Required(), mcp.Description("Path to the .ipynb file.")),
		mcp.WithNumber("index", mcp.Description("Cell index to insert at (single insert mode).")),
		mcp.WithString("source", mcp.Description("Markdown source for the new cell (single insert mode).")),
		mcp.WithArray("indices", mcp.Description("Cell indices to insert at (batch mode)."), mcp.WithNumberItems()),
		mcp.WithArray("sources", mcp.Description("Markdown sources matching the indices order (batch mode)."), mcp.WithStringItems()),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, err := req.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		indices, sources, err := parseInsertCellsArgs(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return withNotebookMutation(path, func(nb *notebook.Notebook) error {
			for i := range indices {
				if err := notebook.InsertMarkdownCell(nb, indices[i], sources[i]); err != nil {
					return err
				}
			}
			return nil
		})
	})
}
