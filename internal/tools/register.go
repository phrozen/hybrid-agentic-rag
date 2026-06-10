package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/phrozen/hybrid-agentic-rag/internal/store"
)

func Register(server *mcp.Server, s *store.Store) {
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "list_documents",
			Description: "List all documents in the corpus grouped by category with line and byte counts",
		},
		List(s),
	)

	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "search",
			Description: "Search the corpus using hybrid (RFF), keyword (text), or semantic mode and return ranked chunk matches",
		},
		Search(s),
	)

	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "get_content",
			Description: "Retrieve a specific chunk's full markdown content by its index from search results",
		},
		Content(s),
	)

	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "get_document",
			Description: "Retrieve a full document's markdown content by its source path",
		},
		Document(s),
	)
}
