# ipynb-mcp Installation Guide for AI Tools

Last updated: 2026-05-08

This guide shows how to install and configure this MCP server (`ipynb-mcp`) in common AI coding tools.

## 1. Build the server once

From this repository root:

```bash
go build -o ipynb-mcp .
```

On Windows, this produces `ipynb-mcp.exe`.

## 1b. Use prebuilt binaries from Releases

You can download prebuilt binaries from the GitHub Releases page instead of building locally.

Available targets:

- Windows: `amd64`, `arm64`
- Linux: `amd64`, `arm64`
- macOS (Darwin): `amd64`, `arm64`

## 2. Recommended server command

Use an absolute binary path in client configs.

Example:

- Windows: `C:\tools\ipynb-mcp\ipynb-mcp.exe`
- macOS/Linux: `/opt/ipynb-mcp/ipynb-mcp`

## 3. Per-client setup

### Codex (CLI + IDE extension)

Codex uses `~/.codex/config.toml` (or project `.codex/config.toml`).

```toml
[mcp_servers.ipynb]
command = "C:\\tools\\ipynb-mcp\\ipynb-mcp.exe"
```

Or via CLI:

```bash
codex mcp add ipynb -- C:\tools\ipynb-mcp\ipynb-mcp.exe
```

### Claude Code

Add a local stdio server:

```bash
claude mcp add --transport stdio ipynb -- C:\tools\ipynb-mcp\ipynb-mcp.exe
```

Check:

```bash
claude mcp list
```

### Cursor

Use either project config `.cursor/mcp.json` or global `~/.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "ipynb": {
      "command": "C:\\tools\\ipynb-mcp\\ipynb-mcp.exe"
    }
  }
}
```

### Windsurf

Edit `~/.codeium/windsurf/mcp_config.json`:

```json
{
  "mcpServers": {
    "ipynb": {
      "command": "C:\\tools\\ipynb-mcp\\ipynb-mcp.exe"
    }
  }
}
```

You can also add servers from the MCP UI in Cascade, then verify in raw config.

### Cline

Edit `cline_mcp_settings.json` from the MCP Servers UI:

```json
{
  "mcpServers": {
    "ipynb": {
      "command": "C:\\tools\\ipynb-mcp\\ipynb-mcp.exe",
      "disabled": false
    }
  }
}
```

For Cline CLI, the default config path is:

- `~/.cline/data/settings/cline_mcp_settings.json`

### OpenCode

Edit `opencode.json` or `opencode.jsonc`:

```json
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "ipynb": {
      "type": "local",
      "command": [
        "C:\\tools\\ipynb-mcp\\ipynb-mcp.exe"
      ],
      "enabled": true
    }
  }
}
```

### Antigravity (Google)

Note: this setup path is documented in Google Firebase's official MCP docs (Antigravity section).

From the Agent pane:

1. Open `MCP Servers`.
2. Open `Manage MCP Servers`.
3. Open `View raw config`.
4. Add this server entry to `mcp_config.json`:

```json
{
  "mcpServers": {
    "ipynb": {
      "command": "C:\\tools\\ipynb-mcp\\ipynb-mcp.exe"
    }
  }
}
```

### Optional: VS Code Copilot (Agent mode)

Workspace-level `.vscode/mcp.json`:

```json
{
  "servers": {
    "ipynb": {
      "type": "stdio",
      "command": "C:\\tools\\ipynb-mcp\\ipynb-mcp.exe"
    }
  }
}
```

## 4. Verify installation

In your AI tool, ask:

`Call read_notebook on <path>.ipynb`

If connected correctly, you should see tools like:

- `read_notebook`
- `create_notebook`
- `insert_markdown_cell`
- `insert_code_cell`
- `update_markdown_cell`
- `update_code_cell`
- `delete_cell`
- `read_cell_output`

## Notes

- Prefer binary command over `go run .` for faster startup.
- Use absolute paths in all configs.
- Restart the AI tool after config changes if tools do not appear.
- Release assets are generated from `.goreleaser.yaml` and include binaries for Windows/Linux/Darwin.

## References

- Codex MCP: https://developers.openai.com/codex/mcp
- Codex config reference: https://developers.openai.com/codex/config-reference
- Claude Code MCP: https://code.claude.com/docs/en/mcp
- Cursor MCP: https://docs.cursor.com/advanced/model-context-protocol
- Cursor MCP CLI: https://docs.cursor.com/cli/mcp
- Windsurf MCP: https://docs.windsurf.com/windsurf/cascade/mcp
- Cline MCP setup: https://docs.cline.bot/mcp/adding-and-configuring-servers
- Cline CLI config: https://docs.cline.bot/cline-cli/configuration
- OpenCode MCP servers: https://opencode.ai/docs/mcp-servers/
- OpenCode install: https://opencode.ai/docs/
- Firebase MCP (includes Antigravity/Cursor/Cline/Windsurf examples): https://firebase.google.com/docs/cli/mcp-server
- VS Code MCP config reference: https://code.visualstudio.com/docs/copilot/reference/mcp-configuration
