package tools

import (
	"fmt"
	"sort"

	"github.com/mark3labs/mcp-go/mcp"
)

func parseInsertCellsArgs(req mcp.CallToolRequest) ([]int, []string, error) {
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

		type insertPair struct {
			index  int
			source string
		}

		pairs := make([]insertPair, len(indices))
		for i := range indices {
			pairs[i] = insertPair{index: indices[i], source: sources[i]}
		}

		sort.SliceStable(pairs, func(i, j int) bool {
			return pairs[i].index > pairs[j].index
		})

		for i := range pairs {
			indices[i] = pairs[i].index
			sources[i] = pairs[i].source
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
