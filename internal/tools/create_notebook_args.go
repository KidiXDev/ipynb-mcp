package tools

import (
	"fmt"
	"strings"

	"github.com/kidixdev/ipynb-mcp/internal/notebook"
	"github.com/mark3labs/mcp-go/mcp"
)

type createNotebookArgs struct {
	Cells []createNotebookCellArg `json:"cells"`
}

type createNotebookCellArg struct {
	CellType string `json:"cell_type"`
	Source   string `json:"source"`
}

func parseCreateNotebookInitialCells(req mcp.CallToolRequest) ([]notebook.InitialCell, error) {
	args := req.GetArguments()
	if _, ok := args["cells"]; !ok {
		return nil, nil
	}

	var parsed createNotebookArgs
	if err := req.BindArguments(&parsed); err != nil {
		return nil, err
	}

	if len(parsed.Cells) == 0 {
		return nil, fmt.Errorf("cells must not be empty when provided")
	}

	initialCells := make([]notebook.InitialCell, 0, len(parsed.Cells))
	for i, c := range parsed.Cells {
		cellType := strings.ToLower(strings.TrimSpace(c.CellType))
		if cellType != notebook.CellTypeMarkdown && cellType != notebook.CellTypeCode {
			return nil, fmt.Errorf("cells[%d].cell_type must be %q or %q", i, notebook.CellTypeMarkdown, notebook.CellTypeCode)
		}
		initialCells = append(initialCells, notebook.InitialCell{
			CellType: cellType,
			Source:   c.Source,
		})
	}

	return initialCells, nil
}
