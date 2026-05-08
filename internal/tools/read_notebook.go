package tools

import (
	"context"

	"github.com/kidixdev/ipynb-mcp/internal/notebook"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterReadNotebook(s *server.MCPServer) {
	tool := mcp.NewTool(
		"read_notebook",
		mcp.WithDescription("Read a .ipynb notebook and return a human-readable notebook preview."),
		mcp.WithString(
			"path",
			mcp.Required(),
			mcp.Description("Path to the .ipynb file."),
		),
		mcp.WithBoolean(
			"include_outputs",
			mcp.Description("Whether to include code cell outputs in the preview. Default: true."),
		),
		mcp.WithNumber(
			"max_output_chars_per_cell",
			mcp.Description("Maximum output characters per code cell in preview. Default: 1200."),
		),
		mcp.WithNumber(
			"max_total_output_chars",
			mcp.Description("Maximum total output characters across notebook preview. Default: 6000."),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, err := req.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		nb, err := notebook.ReadNotebook(path)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		includeOutputs := req.GetBool("include_outputs", true)
		maxPerCell := req.GetInt("max_output_chars_per_cell", notebook.DefaultMaxOutputCharsPerCell)
		maxTotal := req.GetInt("max_total_output_chars", notebook.DefaultMaxTotalOutputChars)
		if maxPerCell < 0 {
			return mcp.NewToolResultError("max_output_chars_per_cell must be >= 0"), nil
		}
		if maxTotal < 0 {
			return mcp.NewToolResultError("max_total_output_chars must be >= 0"), nil
		}

		opts := notebook.RenderOptions{
			IncludeOutputs:        includeOutputs,
			MaxOutputCharsPerCell: maxPerCell,
			MaxTotalOutputChars:   maxTotal,
		}
		return mcp.NewToolResultText(notebook.RenderNotebookWithOptions(path, nb, opts)), nil
	})
}
