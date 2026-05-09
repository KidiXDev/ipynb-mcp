package tools

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestParseUpdateCellsArgsSingle(t *testing.T) {
	t.Parallel()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"index":  2.0,
				"source": "x = 1\n",
			},
		},
	}

	indices, sources, err := parseUpdateCellsArgs(req)
	if err != nil {
		t.Fatalf("parseUpdateCellsArgs returned error: %v", err)
	}
	if len(indices) != 1 || indices[0] != 2 {
		t.Fatalf("unexpected indices: %#v", indices)
	}
	if len(sources) != 1 || sources[0] != "x = 1\n" {
		t.Fatalf("unexpected sources: %#v", sources)
	}
}

func TestParseUpdateCellsArgsBatch(t *testing.T) {
	t.Parallel()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"indices": []any{0.0, 2.0},
				"sources": []any{"# A\n", "# B\n"},
			},
		},
	}

	indices, sources, err := parseUpdateCellsArgs(req)
	if err != nil {
		t.Fatalf("parseUpdateCellsArgs returned error: %v", err)
	}
	if len(indices) != 2 || indices[0] != 0 || indices[1] != 2 {
		t.Fatalf("unexpected indices: %#v", indices)
	}
	if len(sources) != 2 || sources[0] != "# A\n" || sources[1] != "# B\n" {
		t.Fatalf("unexpected sources: %#v", sources)
	}
}

func TestParseUpdateCellsArgsBatchLengthMismatch(t *testing.T) {
	t.Parallel()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"indices": []any{0.0, 1.0},
				"sources": []any{"# A\n"},
			},
		},
	}

	_, _, err := parseUpdateCellsArgs(req)
	if err == nil {
		t.Fatalf("expected an error for mismatched lengths")
	}
}

func TestParseUpdateCellsArgsBatchEmpty(t *testing.T) {
	t.Parallel()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"indices": []any{},
				"sources": []any{},
			},
		},
	}

	_, _, err := parseUpdateCellsArgs(req)
	if err == nil {
		t.Fatalf("expected an error for empty indices")
	}
}
