package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/phrozen/hybrid-agentic-rag/internal/store"
)

type ContentInput struct {
	Index int `json:"index" jsonschema:"content index from search results"`
}

type ContentOutput struct {
	Index      int    `json:"index" jsonschema:"content index"`
	Source     string `json:"source" jsonschema:"source document path"`
	Breadcrumb string `json:"breadcrumb" jsonschema:"heading hierarchy path"`
	Heading    string `json:"heading" jsonschema:"section heading"`
	Content    string `json:"content" jsonschema:"full markdown content"`
}

func Content(s *store.Store) func(context.Context, *mcp.CallToolRequest, ContentInput) (*mcp.CallToolResult, ContentOutput, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, input ContentInput) (*mcp.CallToolResult, ContentOutput, error) {
		chunk, ok := s.GetContent(input.Index)
		if !ok {
			return nil, ContentOutput{}, fmt.Errorf("chunk %d not found", input.Index)
		}
		out := ContentOutput{
			Index:      input.Index,
			Source:     chunk.Source,
			Breadcrumb: chunk.Breadcrumb,
			Heading:    chunk.Heading,
			Content:    chunk.Content,
		}
		return prettyJSONText(out), ContentOutput{}, nil
	}
}
