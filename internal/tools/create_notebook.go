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
		mcp.WithArray(
			"cells",
			mcp.Description("Optional initial cells. Each item requires cell_type ('markdown' or 'code') and source."),
			mcp.Items(map[string]any{
				"type": "object",
				"properties": map[string]any{
					"cell_type": map[string]any{
						"type": "string",
						"enum": []string{notebook.CellTypeMarkdown, notebook.CellTypeCode},
					},
					"source": map[string]any{
						"type": "string",
					},
				},
				"required": []string{"cell_type", "source"},
			}),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, err := req.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if _, hasTitle := req.GetArguments()["title"]; hasTitle {
			return mcp.NewToolResultError("title is not supported; provide all content using cells"), nil
		}
		initialCells, err := parseCreateNotebookInitialCells(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		nb, err := notebook.CreateNotebookWithCells(path, initialCells)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(notebook.RenderNotebook(path, nb)), nil
	})
}
