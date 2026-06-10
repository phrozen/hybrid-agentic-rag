package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/phrozen/hybrid-agentic-rag/internal/models"
	"github.com/phrozen/hybrid-agentic-rag/internal/search/semantic"
	"github.com/phrozen/hybrid-agentic-rag/internal/search/text"
	"github.com/phrozen/hybrid-agentic-rag/internal/store"
	"github.com/phrozen/hybrid-agentic-rag/internal/tools"
	"github.com/spf13/cobra"
)

var (
	dataDir string
	port    int
)

func init() {
	serveCmd.Flags().StringVarP(&dataDir, "data", "d", "", "Directory path containing parsed index artifacts")
	serveCmd.Flags().IntVarP(&port, "port", "p", 3001, "Port to listen on")
	_ = serveCmd.MarkFlagRequired("data")

	RootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start search engine query server",
	Long:  `Launches an MCP server exposing hybrid search tools over HTTP.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		docPath := filepath.Join(dataDir, "documents.json")
		fmt.Printf("Loading documents from %s...\n", docPath)
		docBytes, err := os.ReadFile(docPath)
		if err != nil {
			return fmt.Errorf("failed to read documents file: %w", err)
		}
		var docs []*models.Document
		if err := json.Unmarshal(docBytes, &docs); err != nil {
			return fmt.Errorf("failed to unmarshal documents: %w", err)
		}
		fmt.Printf("  └─ Loaded %d documents\n", len(docs))

		chunkPath := filepath.Join(dataDir, "chunks.json")
		fmt.Printf("Loading chunks from %s...\n", chunkPath)
		chunkBytes, err := os.ReadFile(chunkPath)
		if err != nil {
			return fmt.Errorf("failed to read chunks file: %w", err)
		}
		var chunks []*models.Chunk
		if err := json.Unmarshal(chunkBytes, &chunks); err != nil {
			return fmt.Errorf("failed to unmarshal chunks: %w", err)
		}
		fmt.Printf("  └─ Loaded %d chunks\n", len(chunks))

		idxPath := filepath.Join(dataDir, "index.bin")
		fmt.Printf("Loading keyword inverted index from %s...\n", idxPath)
		idxBytes, err := os.ReadFile(idxPath)
		if err != nil {
			return fmt.Errorf("failed to read inverted index binary file: %w", err)
		}
		ii := text.NewInvertedIndex()
		if err := ii.UnmarshalBinary(idxBytes); err != nil {
			return fmt.Errorf("failed to unmarshal inverted index binary: %w", err)
		}
		fmt.Printf("  └─ Loaded keyword search index successfully\n")

		embPath := filepath.Join(dataDir, "embeddings.bin")
		var vi *semantic.VectorIndex
		if _, err := os.Stat(embPath); os.IsNotExist(err) {
			fmt.Printf("Notice: Semantic embeddings file %s does not exist. Initializing empty semantic index.\n", embPath)
			vi = semantic.NewVectorIndex(semantic.NewClient(nil))
		} else {
			fmt.Printf("Loading dense semantic vector index from %s...\n", embPath)
			embBytes, err := os.ReadFile(embPath)
			if err != nil {
				return fmt.Errorf("failed to read embeddings binary file: %w", err)
			}
			vi = semantic.NewVectorIndex(semantic.NewClient(nil))
			if err := vi.UnmarshalBinary(embBytes); err != nil {
				return fmt.Errorf("failed to unmarshal semantic vector index binary: %w", err)
			}
			fmt.Printf("  └─ Loaded dense vector search index successfully (dimensions: %d)\n", vi.Dim)
		}

		fmt.Println("Initializing unified search Store...")
		unifiedStore := store.New(docs, chunks, ii, vi)
		fmt.Println("  └─ Success: Search engine Store fully initialized and primed in memory!")

		mcpServer := mcp.NewServer(&mcp.Implementation{Name: "Chronicles of Aethelgard Hybrid RAG", Version: "0.1.0"}, &mcp.ServerOptions{
			Instructions: `You are searching the Chronicles of Aethelgard, a fantasy world corpus with 200+ interconnected documents covering races, characters, locations, events, factions, and lore.

- Use search for both exact entity lookup (mode=text, e.g. "Throndor the Wise") and conceptual queries (mode=semantic, e.g. "ancient magical traditions"). Default mode=rff fuses both.
- After search, use get_content with the returned index to read full chunk text.
- Use get_document to read a complete entry by its source path.
- Use list_documents for navigation or to browse what is available.

Always cite your sources — include the source path and breadcrumb from search results in your answers.`,
		})
		tools.Register(mcpServer, unifiedStore)

		handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
			return mcpServer
		}, nil)

		addr := fmt.Sprintf(":%d", port)
		fmt.Printf("\nMCP server listening on http://localhost%s/mcp\n", addr)
		http.Handle("/mcp", handler)
		return http.ListenAndServe(addr, nil)
	},
}
