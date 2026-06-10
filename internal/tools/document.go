package tools

import (
	"context"
	"fmt"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/phrozen/hybrid-agentic-rag/internal/store"
)

type DocumentInput struct {
	Source string `json:"source" jsonschema:"document path (e.g. characters/throndor-the-wise.md)"`
}

type DocumentOutput struct {
	Source  string `json:"source" jsonschema:"document path"`
	Lines   int    `json:"lines" jsonschema:"total lines"`
	Size    int    `json:"size" jsonschema:"size in bytes"`
	Content string `json:"content" jsonschema:"full document markdown"`
}

func Document(s *store.Store) func(context.Context, *mcp.CallToolRequest, DocumentInput) (*mcp.CallToolResult, DocumentOutput, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, input DocumentInput) (*mcp.CallToolResult, DocumentOutput, error) {
		doc, ok := s.GetDocument(input.Source)
		if !ok {
			return nil, DocumentOutput{}, fmt.Errorf("document %q not found", input.Source)
		}
		return markdownText(doc.Source, map[string]string{
			"Lines": strconv.Itoa(doc.Lines),
			"Size":  strconv.Itoa(doc.Size) + " bytes",
		}, string(doc.Content)), DocumentOutput{}, nil
	}
}
