package tools

import (
	"context"

	"github.com/kidixdev/ipynb-mcp/internal/notebook"
	"github.com/mark3labs/mcp-go/mcp"
)

func withNotebookMutation(
	path string,
	mutate func(nb *notebook.Notebook) error,
) (*mcp.CallToolResult, error) {
	nb, err := notebook.ReadNotebook(path)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if err := mutate(nb); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if err := notebook.WriteNotebook(path, nb); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(notebook.RenderNotebook(path, nb)), nil
}

type toolHandler func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)
