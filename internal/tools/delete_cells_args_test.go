package tools

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestParseDeleteCellIndicesSingle(t *testing.T) {
	t.Parallel()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"index": 2.0,
			},
		},
	}

	indices, err := parseDeleteCellIndices(req)
	if err != nil {
		t.Fatalf("parseDeleteCellIndices returned error: %v", err)
	}
	if len(indices) != 1 || indices[0] != 2 {
		t.Fatalf("unexpected indices: %#v", indices)
	}
}

func TestParseDeleteCellIndicesBatchSortedDescending(t *testing.T) {
	t.Parallel()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"indices": []any{1.0, 3.0, 2.0},
			},
		},
	}

	indices, err := parseDeleteCellIndices(req)
	if err != nil {
		t.Fatalf("parseDeleteCellIndices returned error: %v", err)
	}
	if len(indices) != 3 || indices[0] != 3 || indices[1] != 2 || indices[2] != 1 {
		t.Fatalf("unexpected sorted indices: %#v", indices)
	}
}

func TestParseDeleteCellIndicesBatchEmpty(t *testing.T) {
	t.Parallel()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"indices": []any{},
			},
		},
	}

	_, err := parseDeleteCellIndices(req)
	if err == nil {
		t.Fatalf("expected error for empty indices")
	}
}

func TestParseDeleteCellIndicesBatchDuplicate(t *testing.T) {
	t.Parallel()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"indices": []any{1.0, 2.0, 1.0},
			},
		},
	}

	_, err := parseDeleteCellIndices(req)
	if err == nil {
		t.Fatalf("expected error for duplicate indices")
	}
}
