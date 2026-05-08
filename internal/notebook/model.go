package notebook

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	CellTypeMarkdown = "markdown"
	CellTypeCode     = "code"
)

type Notebook struct {
	Cells         []Cell         `json:"cells"`
	Metadata      map[string]any `json:"metadata"`
	NBFormat      int            `json:"nbformat"`
	NBFormatMinor int            `json:"nbformat_minor"`
}

type Cell struct {
	CellType       string                     `json:"-"`
	Metadata       map[string]any             `json:"-"`
	Source         SourceLines                `json:"-"`
	ExecutionCount *int                       `json:"-"`
	Outputs        []json.RawMessage          `json:"-"`
	Attachments    map[string]any             `json:"-"`
	Extras         map[string]json.RawMessage `json:"-"`
}

type SourceLines []string

func NewSourceLines(source string) SourceLines {
	if source == "" {
		return SourceLines{}
	}
	parts := strings.SplitAfter(source, "\n")
	if len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return SourceLines(parts)
}

func (s SourceLines) String() string {
	return strings.Join(s, "")
}

func (s SourceLines) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(s))
}

func (s *SourceLines) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*s = SourceLines{}
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		*s = NewSourceLines(asString)
		return nil
	}

	var asSlice []string
	if err := json.Unmarshal(data, &asSlice); err == nil {
		*s = SourceLines(asSlice)
		return nil
	}

	return fmt.Errorf("invalid cell source, expected string or []string")
}

func (c *Cell) UnmarshalJSON(data []byte) error {
	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var cellType string
	if err := mustUnmarshal(raw, "cell_type", &cellType); err != nil {
		return err
	}
	c.CellType = cellType

	c.Metadata = map[string]any{}
	if _, ok := raw["metadata"]; ok {
		if err := json.Unmarshal(raw["metadata"], &c.Metadata); err != nil {
			return fmt.Errorf("invalid cell metadata: %w", err)
		}
	}

	if err := mustUnmarshal(raw, "source", &c.Source); err != nil {
		return err
	}

	if v, ok := raw["execution_count"]; ok {
		if string(v) != "null" {
			var count int
			if err := json.Unmarshal(v, &count); err != nil {
				return fmt.Errorf("invalid execution_count: %w", err)
			}
			c.ExecutionCount = &count
		}
	}

	if v, ok := raw["outputs"]; ok {
		if string(v) == "null" {
			c.Outputs = []json.RawMessage{}
		} else {
			if err := json.Unmarshal(v, &c.Outputs); err != nil {
				return fmt.Errorf("invalid outputs: %w", err)
			}
		}
	}

	if v, ok := raw["attachments"]; ok {
		if err := json.Unmarshal(v, &c.Attachments); err != nil {
			return fmt.Errorf("invalid attachments: %w", err)
		}
	}

	c.Extras = map[string]json.RawMessage{}
	for key, value := range raw {
		if key == "cell_type" || key == "metadata" || key == "source" || key == "execution_count" || key == "outputs" || key == "attachments" {
			continue
		}
		c.Extras[key] = value
	}

	return nil
}

func (c Cell) MarshalJSON() ([]byte, error) {
	raw := map[string]any{}
	for key, value := range c.Extras {
		raw[key] = value
	}

	raw["cell_type"] = c.CellType
	if c.Metadata != nil {
		raw["metadata"] = c.Metadata
	} else {
		raw["metadata"] = map[string]any{}
	}
	raw["source"] = []string(c.Source)

	if c.CellType == CellTypeCode {
		raw["execution_count"] = c.ExecutionCount
		if c.Outputs == nil {
			raw["outputs"] = []json.RawMessage{}
		} else {
			raw["outputs"] = c.Outputs
		}
	} else if c.Attachments != nil {
		raw["attachments"] = c.Attachments
	}

	return json.Marshal(raw)
}

func mustUnmarshal(raw map[string]json.RawMessage, key string, target any) error {
	value, ok := raw[key]
	if !ok {
		return fmt.Errorf("missing required field %q", key)
	}
	if err := json.Unmarshal(value, target); err != nil {
		return fmt.Errorf("invalid %s: %w", key, err)
	}
	return nil
}
