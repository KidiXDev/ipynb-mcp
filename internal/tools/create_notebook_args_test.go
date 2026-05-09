package tools

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestParseCreateNotebookInitialCellsMissing(t *testing.T) {
	t.Parallel()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"path": "x.ipynb",
			},
		},
	}

	cells, err := parseCreateNotebookInitialCells(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cells) != 0 {
		t.Fatalf("expected empty cells, got %#v", cells)
	}
}

func TestParseCreateNotebookInitialCellsValid(t *testing.T) {
	t.Parallel()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"path": "x.ipynb",
				"cells": []any{
					map[string]any{"cell_type": "markdown", "source": "# Intro\n"},
					map[string]any{"cell_type": "code", "source": "a = 1\n"},
				},
			},
		},
	}

	cells, err := parseCreateNotebookInitialCells(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cells) != 2 {
		t.Fatalf("expected 2 cells, got %d", len(cells))
	}
	if cells[0].CellType != "markdown" || cells[1].CellType != "code" {
		t.Fatalf("unexpected cell types: %#v", cells)
	}
}

func TestParseCreateNotebookInitialCellsNormalized(t *testing.T) {
	t.Parallel()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"path": "x.ipynb",
				"cells": []any{
					map[string]any{"cell_type": " Markdown ", "source": "# Intro\n"},
					map[string]any{"cell_type": " CODE ", "source": "a = 1\n"},
				},
			},
		},
	}

	cells, err := parseCreateNotebookInitialCells(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cells) != 2 {
		t.Fatalf("expected 2 cells, got %d", len(cells))
	}
	if cells[0].CellType != "markdown" || cells[1].CellType != "code" {
		t.Fatalf("unexpected normalized cell types: %#v", cells)
	}
}

func TestParseCreateNotebookInitialCellsInvalidType(t *testing.T) {
	t.Parallel()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"path": "x.ipynb",
				"cells": []any{
					map[string]any{"cell_type": "raw", "source": "x"},
				},
			},
		},
	}

	_, err := parseCreateNotebookInitialCells(req)
	if err == nil {
		t.Fatalf("expected error for unsupported cell type")
	}
}

func TestParseCreateNotebookInitialCellsEmptyProvided(t *testing.T) {
	t.Parallel()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"path":  "x.ipynb",
				"cells": []any{},
			},
		},
	}

	_, err := parseCreateNotebookInitialCells(req)
	if err == nil {
		t.Fatalf("expected error when empty cells is provided")
	}
}
