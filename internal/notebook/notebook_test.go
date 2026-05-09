package notebook

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadNotebookParsesStringAndSliceSources(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "sample.ipynb")

	raw := `{
  "cells": [
    {
      "cell_type": "markdown",
      "metadata": {"role":"intro"},
      "source": ["# Title\n", "Paragraph\n"]
    },
    {
      "cell_type": "code",
      "metadata": {"collapsed": false},
      "source": "print('hello')\n",
      "execution_count": 3,
      "outputs": []
    }
  ],
  "metadata": {"kernelspec": {"name":"python3"}},
  "nbformat": 4,
  "nbformat_minor": 5
}`

	if err := os.WriteFile(path, []byte(raw), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	nb, err := ReadNotebook(path)
	if err != nil {
		t.Fatalf("read notebook: %v", err)
	}

	if got, want := len(nb.Cells), 2; got != want {
		t.Fatalf("cells length: got %d want %d", got, want)
	}

	if got, want := nb.Cells[0].Source.String(), "# Title\nParagraph\n"; got != want {
		t.Fatalf("markdown source: got %q want %q", got, want)
	}

	if got, want := nb.Cells[1].Source.String(), "print('hello')\n"; got != want {
		t.Fatalf("code source: got %q want %q", got, want)
	}
}

func TestCreateNotebookEmpty(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "new", "empty.ipynb")

	nb, err := CreateNotebook(path)
	if err != nil {
		t.Fatalf("create notebook: %v", err)
	}

	if got, want := nb.NBFormat, 4; got != want {
		t.Fatalf("nbformat: got %d want %d", got, want)
	}
	if got, want := nb.NBFormatMinor, 5; got != want {
		t.Fatalf("nbformat_minor: got %d want %d", got, want)
	}
	if got, want := len(nb.Cells), 0; got != want {
		t.Fatalf("cells length: got %d want %d", got, want)
	}

	reloaded, err := ReadNotebook(path)
	if err != nil {
		t.Fatalf("reload notebook: %v", err)
	}
	if got, want := len(reloaded.Cells), 0; got != want {
		t.Fatalf("reloaded cells length: got %d want %d", got, want)
	}
}

func TestCreateNotebookWithInitialCells(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "new", "with-cells.ipynb")

	nb, err := CreateNotebookWithCells(path, []InitialCell{
		{CellType: CellTypeMarkdown, Source: "## Section\n"},
		{CellType: CellTypeCode, Source: "x = 1\n"},
	})
	if err != nil {
		t.Fatalf("create notebook: %v", err)
	}

	if got, want := len(nb.Cells), 2; got != want {
		t.Fatalf("cells length: got %d want %d", got, want)
	}
	if got, want := nb.Cells[0].CellType, CellTypeMarkdown; got != want {
		t.Fatalf("first cell type: got %q want %q", got, want)
	}
	if got, want := nb.Cells[1].CellType, CellTypeCode; got != want {
		t.Fatalf("second cell type: got %q want %q", got, want)
	}
}

func TestCreateNotebookWithInitialCellsRejectsInvalidCellType(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "new", "invalid-cell-type.ipynb")

	_, err := CreateNotebookWithCells(path, []InitialCell{
		{CellType: "raw", Source: "x"},
	})
	if err == nil {
		t.Fatalf("expected error for unsupported initial cell type")
	}
}

func TestInsertUpdateDeleteCells(t *testing.T) {
	t.Parallel()

	nb := &Notebook{
		Cells:    []Cell{},
		Metadata: map[string]any{"project": "x"},
	}

	if err := InsertMarkdownCell(nb, 0, "# Intro\n"); err != nil {
		t.Fatalf("insert markdown: %v", err)
	}
	if err := InsertCodeCell(nb, 1, "a = 1\n"); err != nil {
		t.Fatalf("insert code: %v", err)
	}

	if got, want := len(nb.Cells), 2; got != want {
		t.Fatalf("cells length after insert: got %d want %d", got, want)
	}
	if got, want := nb.Cells[1].CellType, CellTypeCode; got != want {
		t.Fatalf("inserted code type: got %q want %q", got, want)
	}
	if nb.Cells[1].Outputs == nil {
		t.Fatalf("inserted code outputs should be initialized")
	}
	if got, want := len(nb.Cells[1].Outputs), 0; got != want {
		t.Fatalf("inserted code outputs length: got %d want %d", got, want)
	}

	nb.Cells[0].Metadata = map[string]any{"tag": "keep"}
	if err := UpdateMarkdownCell(nb, 0, "## Updated\n"); err != nil {
		t.Fatalf("update markdown: %v", err)
	}
	if got, want := nb.Cells[0].Metadata["tag"], "keep"; got != want {
		t.Fatalf("markdown metadata not preserved: got %v want %v", got, want)
	}

	nb.Cells[1].Metadata = map[string]any{"custom": true}
	if err := UpdateCodeCell(nb, 1, "a = 2\n"); err != nil {
		t.Fatalf("update code: %v", err)
	}
	if got, want := nb.Cells[1].Metadata["custom"], true; got != want {
		t.Fatalf("code metadata not preserved: got %v want %v", got, want)
	}
	if nb.Cells[1].ExecutionCount != nil {
		t.Fatalf("execution_count should be reset to nil")
	}
	if nb.Cells[1].Outputs == nil || len(nb.Cells[1].Outputs) != 0 {
		t.Fatalf("outputs should be reset to an empty list")
	}

	if err := DeleteCell(nb, 0); err != nil {
		t.Fatalf("delete cell: %v", err)
	}
	if got, want := len(nb.Cells), 1; got != want {
		t.Fatalf("cells length after delete: got %d want %d", got, want)
	}
}

func TestWriteAndEditPreserveNotebookMetadata(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "meta.ipynb")

	nb := &Notebook{
		Cells: []Cell{
			NewMarkdownCell("# A", nil),
		},
		Metadata: map[string]any{
			"kernelspec": map[string]any{"name": "python3"},
			"custom":     map[string]any{"team": "ai"},
		},
		NBFormat:      4,
		NBFormatMinor: 5,
	}

	if err := WriteNotebook(path, nb); err != nil {
		t.Fatalf("write notebook: %v", err)
	}

	loaded, err := ReadNotebook(path)
	if err != nil {
		t.Fatalf("read notebook: %v", err)
	}

	if err := InsertCodeCell(loaded, 1, "print('ok')\n"); err != nil {
		t.Fatalf("insert code cell: %v", err)
	}
	if err := WriteNotebook(path, loaded); err != nil {
		t.Fatalf("rewrite notebook: %v", err)
	}

	reloaded, err := ReadNotebook(path)
	if err != nil {
		t.Fatalf("re-read notebook: %v", err)
	}

	custom, ok := reloaded.Metadata["custom"].(map[string]any)
	if !ok {
		t.Fatalf("custom metadata missing or wrong type: %#v", reloaded.Metadata["custom"])
	}
	if got, want := custom["team"], "ai"; got != want {
		t.Fatalf("custom metadata changed: got %v want %v", got, want)
	}
}

func TestValidationErrors(t *testing.T) {
	t.Parallel()

	if err := ValidateNotebookPath("bad.txt"); err == nil {
		t.Fatalf("expected invalid extension error")
	}

	nb := &Notebook{}
	err := DeleteCell(nb, 0)
	if err == nil {
		t.Fatalf("expected out-of-range error")
	}
	if !strings.Contains(err.Error(), "out of range") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInsertIntoEmptyNotebook(t *testing.T) {
	t.Parallel()

	nb := &Notebook{Cells: []Cell{}}

	if err := InsertMarkdownCell(nb, 0, "# Intro\n"); err != nil {
		t.Fatalf("insert markdown into empty notebook: %v", err)
	}
	if err := InsertCodeCell(nb, 1, "a = 1\n"); err != nil {
		t.Fatalf("insert code at end after first insert: %v", err)
	}

	if got, want := len(nb.Cells), 2; got != want {
		t.Fatalf("cells length after inserts: got %d want %d", got, want)
	}
}

func TestInsertIntoEmptyNotebookInvalidIndex(t *testing.T) {
	t.Parallel()

	nb := &Notebook{Cells: []Cell{}}

	err := InsertMarkdownCell(nb, 1, "# Intro\n")
	if err == nil {
		t.Fatalf("expected out-of-range error")
	}
	if !strings.Contains(err.Error(), "out of range") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderNotebookIncludesCodeOutputs(t *testing.T) {
	t.Parallel()

	nb := &Notebook{
		Cells: []Cell{
			NewCodeCell("print(\"hello\")\n", nil),
			NewCodeCell("df.head()\n", nil),
		},
	}

	nb.Cells[0].Outputs = []json.RawMessage{
		json.RawMessage(`{"output_type":"stream","name":"stdout","text":"hello\n"}`),
	}
	nb.Cells[1].Outputs = []json.RawMessage{
		json.RawMessage(`{"output_type":"execute_result","data":{"text/plain":"| name | age |\n|------|-----|\n| John | 20  |\n"},"metadata":{},"execution_count":1}`),
	}

	rendered := RenderNotebook("example.ipynb", nb)

	if !strings.Contains(rendered, "Cell 0 [code]\nprint(\"hello\")\n\nOutput:\nhello") {
		t.Fatalf("missing stream output rendering:\n%s", rendered)
	}
	if !strings.Contains(rendered, "Cell 1 [code]\ndf.head()\n\nOutput:\n| name | age |\n|------|-----|\n| John | 20  |") {
		t.Fatalf("missing execute_result output rendering:\n%s", rendered)
	}
}

func TestRenderNotebookIncludesErrorOutputs(t *testing.T) {
	t.Parallel()

	nb := &Notebook{
		Cells: []Cell{
			NewCodeCell("1 / 0\n", nil),
		},
	}
	nb.Cells[0].Outputs = []json.RawMessage{
		json.RawMessage(`{
			"output_type":"error",
			"ename":"ZeroDivisionError",
			"evalue":"division by zero",
			"traceback":[
				"\u001b[31mZeroDivisionError: division by zero\u001b[0m"
			]
		}`),
	}

	rendered := RenderNotebook("example.ipynb", nb)

	if !strings.Contains(rendered, "Error:\nZeroDivisionError: division by zero") {
		t.Fatalf("missing error output rendering:\n%s", rendered)
	}
}

func TestRenderNotebookTruncatesLongOutputs(t *testing.T) {
	t.Parallel()

	long := strings.Repeat("x", 50)
	nb := &Notebook{
		Cells: []Cell{
			NewCodeCell("print('long')\n", nil),
		},
	}
	nb.Cells[0].Outputs = []json.RawMessage{
		json.RawMessage(`{"output_type":"stream","name":"stdout","text":"` + long + `"}`),
	}

	rendered := RenderNotebookWithOptions("example.ipynb", nb, RenderOptions{
		IncludeOutputs:        true,
		MaxOutputCharsPerCell: 20,
		MaxTotalOutputChars:   100,
	})

	if !strings.Contains(rendered, "Output:\nxxxxxxxxxxxx") {
		t.Fatalf("expected truncated output preview:\n%s", rendered)
	}
	if !strings.Contains(rendered, "[output truncated for token efficiency; use read_cell_output to fetch more]") {
		t.Fatalf("expected truncation hint:\n%s", rendered)
	}
}

func TestRenderNotebookMaxTotalOutputChars(t *testing.T) {
	t.Parallel()

	nb := &Notebook{
		Cells: []Cell{
			NewCodeCell("print('a')\n", nil),
			NewCodeCell("print('b')\n", nil),
		},
	}
	nb.Cells[0].Outputs = []json.RawMessage{
		json.RawMessage(`{"output_type":"stream","name":"stdout","text":"1234567890"}`),
	}
	nb.Cells[1].Outputs = []json.RawMessage{
		json.RawMessage(`{"output_type":"stream","name":"stdout","text":"abcdefghij"}`),
	}

	rendered := RenderNotebookWithOptions("example.ipynb", nb, RenderOptions{
		IncludeOutputs:        true,
		MaxOutputCharsPerCell: 100,
		MaxTotalOutputChars:   15,
	})

	if !strings.Contains(rendered, "Output:\n1234567") {
		t.Fatalf("missing first output:\n%s", rendered)
	}
	if strings.Contains(rendered, "Cell 1 [code]\nprint('b')\n\nOutput:") {
		t.Fatalf("expected total output budget to be exhausted before second output:\n%s", rendered)
	}
}
