package tools

import (
	"context"
	"fmt"

	"github.com/kidixdev/ipynb-mcp/internal/notebook"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	defaultOutputChunkLimit = 4000
	maxOutputChunkLimit     = 50000
)

func RegisterReadCellOutput(s *server.MCPServer) {
	tool := mcp.NewTool(
		"read_cell_output",
		mcp.WithDescription("Read a code cell's rendered output in chunks for token-efficient access."),
		mcp.WithString("path", mcp.Required(), mcp.Description("Path to the .ipynb file.")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("Code cell index to read output from.")),
		mcp.WithNumber("offset", mcp.Description("Character offset into rendered output. Default: 0.")),
		mcp.WithNumber("limit", mcp.Description("Maximum characters to return. Default: 4000, max: 50000.")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, err := req.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		index, err := req.RequireInt("index")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		offset := req.GetInt("offset", 0)
		limit := req.GetInt("limit", defaultOutputChunkLimit)
		if offset < 0 {
			return mcp.NewToolResultError("offset must be >= 0"), nil
		}
		if limit <= 0 || limit > maxOutputChunkLimit {
			return mcp.NewToolResultError(fmt.Sprintf("limit must be between 1 and %d", maxOutputChunkLimit)), nil
		}

		nb, err := notebook.ReadNotebook(path)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if len(nb.Cells) == 0 {
			return mcp.NewToolResultError("notebook has no cells"), nil
		}
		if index < 0 || index >= len(nb.Cells) {
			return mcp.NewToolResultError(fmt.Sprintf("cell index %d out of range: valid range is 0..%d", index, len(nb.Cells)-1)), nil
		}

		cell := nb.Cells[index]
		if cell.CellType != notebook.CellTypeCode {
			return mcp.NewToolResultError(fmt.Sprintf("cell %d is %q, expected %q", index, cell.CellType, notebook.CellTypeCode)), nil
		}

		rendered := notebook.RenderCellOutputs(cell.Outputs)
		if rendered == "" {
			return mcp.NewToolResultText(fmt.Sprintf("Cell %d has no rendered outputs.", index)), nil
		}

		chunk, start, end, total, hasMore, err := paginateText(rendered, offset, limit)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		nextOffset := end
		return mcp.NewToolResultText(
			fmt.Sprintf(
				"Cell %d output chunk (%d-%d of %d chars)\n\n%s\n\nhas_more: %t\nnext_offset: %d",
				index,
				start,
				end,
				total,
				chunk,
				hasMore,
				nextOffset,
			),
		), nil
	})
}

func paginateText(text string, offset int, limit int) (chunk string, start int, end int, total int, hasMore bool, err error) {
	runes := []rune(text)
	total = len(runes)
	if offset > total {
		return "", 0, 0, total, false, fmt.Errorf("offset %d out of range for output size %d", offset, total)
	}
	start = offset
	end = offset + limit
	if end > total {
		end = total
	}
	chunk = string(runes[start:end])
	hasMore = end < total
	return chunk, start, end, total, hasMore, nil
}
