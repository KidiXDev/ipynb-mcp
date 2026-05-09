package tools

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestParseInsertCellsArgsSingle(t *testing.T) {
	t.Parallel()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"index":  2.0,
				"source": "x = 1\n",
			},
		},
	}

	indices, sources, err := parseInsertCellsArgs(req)
	if err != nil {
		t.Fatalf("parseInsertCellsArgs returned error: %v", err)
	}
	if len(indices) != 1 || indices[0] != 2 {
		t.Fatalf("unexpected indices: %#v", indices)
	}
	if len(sources) != 1 || sources[0] != "x = 1\n" {
		t.Fatalf("unexpected sources: %#v", sources)
	}
}

func TestParseInsertCellsArgsBatch(t *testing.T) {
	t.Parallel()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"indices": []any{0.0, 2.0},
				"sources": []any{"# A\n", "# B\n"},
			},
		},
	}

	indices, sources, err := parseInsertCellsArgs(req)
	if err != nil {
		t.Fatalf("parseInsertCellsArgs returned error: %v", err)
	}
	if len(indices) != 2 || indices[0] != 2 || indices[1] != 0 {
		t.Fatalf("unexpected indices: %#v", indices)
	}
	if len(sources) != 2 || sources[0] != "# B\n" || sources[1] != "# A\n" {
		t.Fatalf("unexpected sources: %#v", sources)
	}
}

func TestParseInsertCellsArgsBatchLengthMismatch(t *testing.T) {
	t.Parallel()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"indices": []any{0.0, 1.0},
				"sources": []any{"# A\n"},
			},
		},
	}

	_, _, err := parseInsertCellsArgs(req)
	if err == nil {
		t.Fatalf("expected an error for mismatched lengths")
	}
}

func TestParseInsertCellsArgsBatchEmpty(t *testing.T) {
	t.Parallel()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"indices": []any{},
				"sources": []any{},
			},
		},
	}

	_, _, err := parseInsertCellsArgs(req)
	if err == nil {
		t.Fatalf("expected an error for empty indices")
	}
}

func TestParseInsertCellsArgsOnlyIndices(t *testing.T) {
	t.Parallel()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"indices": []any{0.0, 1.0},
			},
		},
	}

	_, _, err := parseInsertCellsArgs(req)
	if err == nil {
		t.Fatalf("expected an error when only indices are provided")
	}
}

func TestParseInsertCellsArgsOnlySources(t *testing.T) {
	t.Parallel()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"sources": []any{"# A\n", "# B\n"},
			},
		},
	}

	_, _, err := parseInsertCellsArgs(req)
	if err == nil {
		t.Fatalf("expected an error when only sources are provided")
	}
}
