package tools

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

func parseUpdateCellsArgs(req mcp.CallToolRequest) ([]int, []string, error) {
	args := req.GetArguments()
	_, hasIndices := args["indices"]
	_, hasSources := args["sources"]

	if hasIndices || hasSources {
		indices, err := req.RequireIntSlice("indices")
		if err != nil {
			return nil, nil, err
		}
		sources, err := req.RequireStringSlice("sources")
		if err != nil {
			return nil, nil, err
		}
		if len(indices) == 0 {
			return nil, nil, fmt.Errorf("indices must not be empty")
		}
		if len(indices) != len(sources) {
			return nil, nil, fmt.Errorf("indices and sources must have the same length")
		}
		return indices, sources, nil
	}

	index, err := req.RequireInt("index")
	if err != nil {
		return nil, nil, err
	}
	source, err := req.RequireString("source")
	if err != nil {
		return nil, nil, err
	}

	return []int{index}, []string{source}, nil
}
