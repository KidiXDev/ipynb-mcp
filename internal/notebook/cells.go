package notebook

import (
	"encoding/json"
	"fmt"
)

func InsertMarkdownCell(nb *Notebook, index int, source string) error {
	if err := validateInsertIndex(nb, index); err != nil {
		return err
	}
	cell := NewMarkdownCell(source, nil)
	nb.Cells = append(nb.Cells[:index], append([]Cell{cell}, nb.Cells[index:]...)...)
	return nil
}

func InsertCodeCell(nb *Notebook, index int, source string) error {
	if err := validateInsertIndex(nb, index); err != nil {
		return err
	}
	cell := NewCodeCell(source, nil)
	nb.Cells = append(nb.Cells[:index], append([]Cell{cell}, nb.Cells[index:]...)...)
	return nil
}

func UpdateMarkdownCell(nb *Notebook, index int, source string) error {
	if err := validateCellIndex(nb, index); err != nil {
		return err
	}
	existing := nb.Cells[index]
	nb.Cells[index] = NewMarkdownCell(source, &existing)
	return nil
}

func UpdateCodeCell(nb *Notebook, index int, source string) error {
	if err := validateCellIndex(nb, index); err != nil {
		return err
	}
	existing := nb.Cells[index]
	nb.Cells[index] = NewCodeCell(source, &existing)
	return nil
}

func DeleteCell(nb *Notebook, index int) error {
	if err := validateCellIndex(nb, index); err != nil {
		return err
	}
	nb.Cells = append(nb.Cells[:index], nb.Cells[index+1:]...)
	return nil
}

func NewMarkdownCell(source string, existing *Cell) Cell {
	cell := Cell{
		CellType:    CellTypeMarkdown,
		Metadata:    map[string]any{},
		Source:      NewSourceLines(source),
		Attachments: nil,
		Extras:      map[string]json.RawMessage{},
	}

	if existing == nil {
		return cell
	}

	cell.Metadata = cloneMap(existing.Metadata)
	cell.Extras = cloneRawMap(existing.Extras)
	if existing.CellType == CellTypeMarkdown && existing.Attachments != nil {
		cell.Attachments = cloneMap(existing.Attachments)
	}
	return cell
}

func NewCodeCell(source string, existing *Cell) Cell {
	cell := Cell{
		CellType:       CellTypeCode,
		Metadata:       map[string]any{},
		Source:         NewSourceLines(source),
		ExecutionCount: nil,
		Outputs:        []json.RawMessage{},
		Extras:         map[string]json.RawMessage{},
	}

	if existing == nil {
		return cell
	}

	cell.Metadata = cloneMap(existing.Metadata)
	cell.Extras = cloneRawMap(existing.Extras)
	return cell
}

func validateInsertIndex(nb *Notebook, index int) error {
	if nb == nil {
		return fmt.Errorf("notebook is nil")
	}
	if index < 0 || index > len(nb.Cells) {
		return fmt.Errorf("insert index %d out of range: valid range is 0..%d", index, len(nb.Cells))
	}
	return nil
}

func validateCellIndex(nb *Notebook, index int) error {
	if nb == nil {
		return fmt.Errorf("notebook is nil")
	}
	max := len(nb.Cells) - 1
	if index < 0 || index > max {
		return fmt.Errorf("cell index %d out of range: valid range is 0..%d", index, max)
	}
	return nil
}

func cloneMap(src map[string]any) map[string]any {
	if src == nil {
		return map[string]any{}
	}
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func cloneRawMap(src map[string]json.RawMessage) map[string]json.RawMessage {
	if src == nil {
		return map[string]json.RawMessage{}
	}
	dst := make(map[string]json.RawMessage, len(src))
	for k, v := range src {
		copied := make([]byte, len(v))
		copy(copied, v)
		dst[k] = copied
	}
	return dst
}
