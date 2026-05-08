package notebook

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var ansiEscapePattern = regexp.MustCompile(`\x1b\[[0-9;]*[A-Za-z]`)

const (
	DefaultMaxOutputCharsPerCell = 1200
	DefaultMaxTotalOutputChars   = 6000
)

type RenderOptions struct {
	IncludeOutputs        bool
	MaxOutputCharsPerCell int
	MaxTotalOutputChars   int
}

func DefaultRenderOptions() RenderOptions {
	return RenderOptions{
		IncludeOutputs:        true,
		MaxOutputCharsPerCell: DefaultMaxOutputCharsPerCell,
		MaxTotalOutputChars:   DefaultMaxTotalOutputChars,
	}
}

func RenderNotebook(path string, nb *Notebook) string {
	return RenderNotebookWithOptions(path, nb, DefaultRenderOptions())
}

func RenderNotebookWithOptions(path string, nb *Notebook, opts RenderOptions) string {
	if !opts.IncludeOutputs {
		opts.MaxOutputCharsPerCell = 0
		opts.MaxTotalOutputChars = 0
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("Notebook: %s\n", filepath.Base(path)))
	remainingTotal := opts.MaxTotalOutputChars
	if remainingTotal <= 0 {
		remainingTotal = -1
	}

	for i, cell := range nb.Cells {
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("Cell %d [%s]\n", i, cell.CellType))

		source := cell.Source.String()
		if source == "" {
			b.WriteString("(empty)\n")
			continue
		}

		b.WriteString(source)
		if !strings.HasSuffix(source, "\n") {
			b.WriteString("\n")
		}

		if opts.IncludeOutputs && cell.CellType == CellTypeCode && len(cell.Outputs) > 0 {
			rendered := RenderCellOutputs(cell.Outputs)
			rendered, wasTruncated := truncateOutputForPreview(rendered, opts.MaxOutputCharsPerCell, &remainingTotal)
			if rendered != "" {
				b.WriteString("\n")
				b.WriteString(rendered)
				if wasTruncated {
					b.WriteString("\n[output truncated for token efficiency; use read_cell_output to fetch more]")
				}
				if !strings.HasSuffix(rendered, "\n") {
					b.WriteString("\n")
				}
			}
		}
	}

	return strings.TrimRight(b.String(), "\n")
}

type notebookOutput struct {
	OutputType string         `json:"output_type"`
	Name       string         `json:"name"`
	Text       SourceLines    `json:"text"`
	Data       map[string]any `json:"data"`
	EName      string         `json:"ename"`
	EValue     string         `json:"evalue"`
	Traceback  []string       `json:"traceback"`
}

func RenderCellOutputs(rawOutputs []json.RawMessage) string {
	var sections []string
	for _, raw := range rawOutputs {
		text, isError := renderSingleOutput(raw)
		if text == "" {
			continue
		}
		if isError {
			sections = append(sections, "Error:\n"+text)
		} else {
			sections = append(sections, "Output:\n"+text)
		}
	}
	return strings.Join(sections, "\n\n")
}

func renderSingleOutput(raw json.RawMessage) (string, bool) {
	var out notebookOutput
	if err := json.Unmarshal(raw, &out); err != nil {
		return "", false
	}

	switch out.OutputType {
	case "stream":
		return normalizeOutputText(out.Text.String()), false
	case "execute_result", "display_data":
		return renderMimeData(out.Data), false
	case "error":
		return renderErrorOutput(out), true
	default:
		return "", false
	}
}

func renderMimeData(data map[string]any) string {
	if len(data) == 0 {
		return ""
	}

	keys := []string{
		"text/markdown",
		"text/plain",
		"text/html",
		"application/json",
	}
	for _, key := range keys {
		if value, ok := data[key]; ok {
			if text := valueToText(value); text != "" {
				return normalizeOutputText(text)
			}
		}
	}

	for _, value := range data {
		if text := valueToText(value); text != "" {
			return normalizeOutputText(text)
		}
	}

	return ""
}

func renderErrorOutput(out notebookOutput) string {
	lines := make([]string, 0, len(out.Traceback)+1)
	for _, trace := range out.Traceback {
		clean := strings.TrimSpace(ansiEscapePattern.ReplaceAllString(trace, ""))
		if clean != "" {
			lines = append(lines, clean)
		}
	}

	if len(lines) > 0 {
		return strings.Join(lines, "\n")
	}

	head := strings.TrimSpace(strings.Trim(out.EName+" "+out.EValue, " "))
	return head
}

func valueToText(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case []any:
		var b strings.Builder
		for _, item := range v {
			b.WriteString(valueToText(item))
		}
		return b.String()
	default:
		return ""
	}
}

func normalizeOutputText(text string) string {
	return strings.TrimRight(text, "\n")
}

func truncateOutputForPreview(text string, maxPerCell int, remainingTotal *int) (string, bool) {
	if text == "" {
		return "", false
	}

	originalLen := runeLen(text)
	visibleLen := originalLen
	truncated := false

	if maxPerCell > 0 && visibleLen > maxPerCell {
		visibleLen = maxPerCell
		truncated = true
	}

	if remainingTotal != nil && *remainingTotal >= 0 {
		if *remainingTotal == 0 {
			return "", true
		}
		if visibleLen > *remainingTotal {
			visibleLen = *remainingTotal
			truncated = true
		}
		*remainingTotal -= visibleLen
	}

	if visibleLen <= 0 {
		return "", true
	}

	if visibleLen < originalLen {
		return truncateRunes(text, visibleLen), true
	}
	return text, truncated
}

func runeLen(s string) int {
	return len([]rune(s))
}

func truncateRunes(s string, count int) string {
	if count <= 0 {
		return ""
	}
	runes := []rune(s)
	if count >= len(runes) {
		return s
	}
	return string(runes[:count])
}
