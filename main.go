package main

import (
	"log"

	"github.com/kidixdev/ipynb-mcp/internal/tools"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"ipynb-mcp",
		"1.1.6",
		server.WithToolCapabilities(true),
		server.WithRecovery(),
	)

	tools.Register(s)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("mcp server error: %v", err)
	}
}
