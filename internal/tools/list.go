package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/phrozen/hybrid-agentic-rag/internal/store"
)

type ListInput struct{}

func List(s *store.Store) func(context.Context, *mcp.CallToolRequest, ListInput) (*mcp.CallToolResult, struct{}, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, _ ListInput) (*mcp.CallToolResult, struct{}, error) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: s.ListDocuments()}},
		}, struct{}{}, nil
	}
}
