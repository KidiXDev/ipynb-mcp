# ipynb-mcp

Go-based MCP server for safe `.ipynb` notebook editing through explicit tools.

The server exposes notebook operations without returning raw notebook JSON to the assistant. Notebook JSON parsing, validation, and writes are handled internally.

## Features

- Read notebook as human-readable preview (`read_notebook`)
- Token-efficient output previews with truncation controls
- Chunked output retrieval for large code outputs (`read_cell_output`)
- Create notebook with nbformat `4.5` (`create_notebook`)
- Insert/update markdown and code cells with explicit tools
- Delete cells by index
- Index and extension validation with clear errors
- Safe file writes (temp file + replace)
- Preserves notebook-level metadata and unrelated cell fields where possible

## Tools

1. `read_notebook`
2. `create_notebook`
3. `insert_markdown_cell`
4. `insert_code_cell`
5. `update_markdown_cell`
6. `update_code_cell`
7. `delete_cell`
8. `read_cell_output`

## Run

```bash
go run .
```

The server uses stdio transport.

## Install

Build a local binary:

```bash
go build -o ipynb-mcp .
```

Then configure your AI client to run the binary as an MCP stdio server.

Prebuilt binaries are also published in GitHub Releases for:

- Windows (`amd64`, `arm64`)
- Linux (`amd64`, `arm64`)
- macOS / Darwin (`amd64`, `arm64`)

See full per-client setup:

- `docs/AI_TOOLS_SETUP.md`

## Example MCP Config (Local)

```json
{
  "mcpServers": {
    "ipynb": {
      "command": "go",
      "args": ["run", "."],
      "cwd": "/path/to/ipynb-mcp"
    }
  }
}
```

You can also use a built binary:

```json
{
  "mcpServers": {
    "ipynb": {
      "command": "/path/to/ipynb-mcp/ipynb-mcp"
    }
  }
}
```

## Tests

```bash
go test ./...
```
