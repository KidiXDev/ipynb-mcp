package notebook

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type InitialCell struct {
	CellType string
	Source   string
}

func ReadNotebook(path string) (*Notebook, error) {
	if err := ValidateNotebookPath(path); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read notebook: %w", err)
	}

	var nb Notebook
	if err := json.Unmarshal(data, &nb); err != nil {
		return nil, fmt.Errorf("parse notebook: %w", err)
	}

	if nb.Metadata == nil {
		nb.Metadata = map[string]any{}
	}
	if nb.Cells == nil {
		nb.Cells = []Cell{}
	}

	return &nb, nil
}

func CreateNotebook(path string) (*Notebook, error) {
	return CreateNotebookWithCells(path, nil)
}

func CreateNotebookWithCells(path string, initialCells []InitialCell) (*Notebook, error) {
	if err := ValidateNotebookPath(path); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create parent directories: %w", err)
	}

	nb := &Notebook{
		Cells:         []Cell{},
		Metadata:      map[string]any{},
		NBFormat:      4,
		NBFormatMinor: 5,
	}

	for i, cell := range initialCells {
		switch cell.CellType {
		case CellTypeMarkdown:
			nb.Cells = append(nb.Cells, NewMarkdownCell(cell.Source, nil))
		case CellTypeCode:
			nb.Cells = append(nb.Cells, NewCodeCell(cell.Source, nil))
		default:
			return nil, fmt.Errorf("initial cell %d has unsupported cell_type %q", i, cell.CellType)
		}
	}

	payload, err := marshalNotebookPayload(nb)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil, fmt.Errorf("notebook already exists: %s", path)
		}
		return nil, fmt.Errorf("create notebook file: %w", err)
	}

	if _, err := f.Write(payload); err != nil {
		_ = f.Close()
		_ = os.Remove(path)
		return nil, fmt.Errorf("write notebook file: %w", err)
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(path)
		return nil, fmt.Errorf("close notebook file: %w", err)
	}

	return nb, nil
}

func WriteNotebook(path string, nb *Notebook) error {
	if err := ValidateNotebookPath(path); err != nil {
		return err
	}
	if nb == nil {
		return fmt.Errorf("notebook is nil")
	}
	if nb.Metadata == nil {
		nb.Metadata = map[string]any{}
	}

	payload, err := marshalNotebookPayload(nb)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("ensure parent directories: %w", err)
	}

	return writeFileAtomically(path, payload)
}

func marshalNotebookPayload(nb *Notebook) ([]byte, error) {
	payload, err := json.MarshalIndent(nb, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal notebook: %w", err)
	}
	return append(payload, '\n'), nil
}

func ValidateNotebookPath(path string) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("path is required")
	}
	if !strings.EqualFold(filepath.Ext(path), ".ipynb") {
		return fmt.Errorf("invalid notebook extension %q: expected .ipynb", filepath.Ext(path))
	}
	return nil
}

func writeFileAtomically(path string, payload []byte) error {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	tmpFile, err := os.CreateTemp(dir, "."+base+".*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	cleanupTemp := func() {
		_ = os.Remove(tmpPath)
	}

	if _, err := tmpFile.Write(payload); err != nil {
		_ = tmpFile.Close()
		cleanupTemp()
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		cleanupTemp()
		return fmt.Errorf("close temp file: %w", err)
	}

	backupPath := path + ".bak"
	backupCreated := false

	if _, err := os.Stat(path); err == nil {
		_ = os.Remove(backupPath)
		if err := os.Rename(path, backupPath); err != nil {
			cleanupTemp()
			return fmt.Errorf("prepare backup: %w", err)
		}
		backupCreated = true
	} else if !errors.Is(err, os.ErrNotExist) {
		cleanupTemp()
		return fmt.Errorf("check existing notebook: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		if backupCreated {
			_ = os.Rename(backupPath, path)
		}
		cleanupTemp()
		return fmt.Errorf("replace notebook file: %w", err)
	}

	if backupCreated {
		_ = os.Remove(backupPath)
	}

	return nil
}
