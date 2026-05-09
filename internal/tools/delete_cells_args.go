package tools

import (
	"fmt"
	"sort"

	"github.com/mark3labs/mcp-go/mcp"
)

func parseDeleteCellIndices(req mcp.CallToolRequest) ([]int, error) {
	args := req.GetArguments()
	if _, hasIndices := args["indices"]; hasIndices {
		indices, err := req.RequireIntSlice("indices")
		if err != nil {
			return nil, err
		}
		if len(indices) == 0 {
			return nil, fmt.Errorf("indices must not be empty")
		}
		sort.Slice(indices, func(i, j int) bool {
			return indices[i] > indices[j]
		})
		for i := 1; i < len(indices); i++ {
			if indices[i] == indices[i-1] {
				return nil, fmt.Errorf("indices must not contain duplicates")
			}
		}
		return indices, nil
	}

	index, err := req.RequireInt("index")
	if err != nil {
		return nil, err
	}
	return []int{index}, nil
}
