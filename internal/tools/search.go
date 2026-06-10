package tools

import (
	"context"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/phrozen/hybrid-agentic-rag/internal/search"
	"github.com/phrozen/hybrid-agentic-rag/internal/store"
)

type SearchInput struct {
	Query string `json:"query" jsonschema:"search query string"`
	K     int    `json:"k,omitempty" jsonschema:"maximum number of results (default 10)"`
	Mode  string `json:"mode,omitempty" jsonschema:"search mode: rff, text, or semantic (default rff)"`
}

type SearchOutput struct {
	Results []SearchResultItem `json:"results" jsonschema:"ranked search results"`
}

type SearchResultItem struct {
	Index      int     `json:"index" jsonschema:"content index"`
	Score      float32 `json:"score" jsonschema:"relevance score"`
	Source     string  `json:"source" jsonschema:"source document path"`
	Breadcrumb string  `json:"breadcrumb" jsonschema:"heading hierarchy path"`
}

func Search(s *store.Store) func(context.Context, *mcp.CallToolRequest, SearchInput) (*mcp.CallToolResult, SearchOutput, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, input SearchInput) (*mcp.CallToolResult, SearchOutput, error) {
		mode := search.SearchRFF
		switch strings.ToLower(input.Mode) {
		case "text":
			mode = search.SearchText
		case "semantic":
			mode = search.SearchSemantic
		}

		k := input.K
		if k <= 0 {
			k = 10
		}

		hits := s.Search(input.Query, k, mode)
		items := make([]SearchResultItem, len(hits))
		for i, h := range hits {
			items[i] = SearchResultItem{
				Index:      h.Index,
				Score:      h.Score,
				Source:     h.Source,
				Breadcrumb: h.Breadcrumb,
			}
		}

		return prettyJSONText(SearchOutput{Results: items}), SearchOutput{}, nil
	}
}
