package tools

import (
	"context"

	"github.com/kidixdev/ipynb-mcp/internal/notebook"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterCreateNotebook(s *server.MCPServer) {
	tool := mcp.NewTool(
		"create_notebook",
		mcp.WithDescription("Create a new nbformat 4.5 .ipynb notebook."),
		mcp.WithString(
			"path",
			mcp.Required(),
			mcp.Description("Path where the new .ipynb file should be created."),
		),
		mcp.WithString(
			"title",
			mcp.Description("Optional notebook title. When present, a first markdown cell '# {title}' is created."),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, err := req.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		title := req.GetString("title", "")

		nb, err := notebook.CreateNotebook(path, title)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(notebook.RenderNotebook(path, nb)), nil
	})
}
