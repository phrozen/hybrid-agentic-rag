package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/phrozen/hybrid-agentic-rag/internal/models"
	"github.com/phrozen/hybrid-agentic-rag/internal/search/semantic"
	"github.com/phrozen/hybrid-agentic-rag/internal/search/text"
	"github.com/phrozen/hybrid-agentic-rag/internal/store"
	"github.com/phrozen/hybrid-agentic-rag/internal/tools"
	"github.com/pterm/pterm"
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
		// 1. Documents
		docPath := filepath.Join(dataDir, "documents.json")
		sp, _ := pterm.DefaultSpinner.Start("Loading documents from " + docPath)
		docBytes, err := os.ReadFile(docPath)
		if err != nil {
			sp.Fail("Failed to read documents file")
			return fmt.Errorf("failed to read documents file: %w", err)
		}
		var docs []*models.Document
		if err := json.Unmarshal(docBytes, &docs); err != nil {
			sp.Fail("Failed to unmarshal documents")
			return fmt.Errorf("failed to unmarshal documents: %w", err)
		}
		sp.Success(fmt.Sprintf("Loaded %d documents", len(docs)))

		// 2. Chunks
		chunkPath := filepath.Join(dataDir, "chunks.json")
		sp, _ = pterm.DefaultSpinner.Start("Loading chunks from " + chunkPath)
		chunkBytes, err := os.ReadFile(chunkPath)
		if err != nil {
			sp.Fail("Failed to read chunks file")
			return fmt.Errorf("failed to read chunks file: %w", err)
		}
		var chunks []*models.Chunk
		if err := json.Unmarshal(chunkBytes, &chunks); err != nil {
			sp.Fail("Failed to unmarshal chunks")
			return fmt.Errorf("failed to unmarshal chunks: %w", err)
		}
		sp.Success(fmt.Sprintf("Loaded %d chunks", len(chunks)))

		// 3. Keyword inverted index
		idxPath := filepath.Join(dataDir, "index.bin")
		sp, _ = pterm.DefaultSpinner.Start("Loading keyword inverted index from " + idxPath)
		idxBytes, err := os.ReadFile(idxPath)
		if err != nil {
			sp.Fail("Failed to read inverted index binary file")
			return fmt.Errorf("failed to read inverted index binary file: %w", err)
		}
		ii := text.NewInvertedIndex()
		if err := ii.UnmarshalBinary(idxBytes); err != nil {
			sp.Fail("Failed to unmarshal inverted index binary")
			return fmt.Errorf("failed to unmarshal inverted index binary: %w", err)
		}
		sp.Success("Loaded keyword search index")

		// 4. Semantic vector index (optional)
		embPath := filepath.Join(dataDir, "embeddings.bin")
		embeddingsSummary := "disabled"
		var vi *semantic.VectorIndex
		if _, err := os.Stat(embPath); os.IsNotExist(err) {
			pterm.Warning.Printf("Semantic embeddings file %s does not exist. Initializing empty semantic index.\n", embPath)
			vi = semantic.NewVectorIndex(semantic.NewClient(nil))
		} else {
			sp, _ = pterm.DefaultSpinner.Start("Loading dense semantic vector index from " + embPath)
			embBytes, err := os.ReadFile(embPath)
			if err != nil {
				sp.Fail("Failed to read embeddings binary file")
				return fmt.Errorf("failed to read embeddings binary file: %w", err)
			}
			vi = semantic.NewVectorIndex(semantic.NewClient(nil))
			if err := vi.UnmarshalBinary(embBytes); err != nil {
				sp.Fail("Failed to unmarshal semantic vector index binary")
				return fmt.Errorf("failed to unmarshal semantic vector index binary: %w", err)
			}
			sp.Success(fmt.Sprintf("Loaded dense vector search index (%d dimensions)", vi.Dim))
			embeddingsSummary = fmt.Sprintf("%d dims", vi.Dim)
		}

		// 5. Unified store
		sp, _ = pterm.DefaultSpinner.Start("Initializing unified search Store")
		unifiedStore := store.New(docs, chunks, ii, vi)
		sp.Success("Search engine Store primed in memory")

		mcpServer := mcp.NewServer(&mcp.Implementation{Name: "Chronicles of Aethelgard Hybrid RAG", Version: "0.1.0"}, &mcp.ServerOptions{
			Instructions: `You are searching the Chronicles of Aethelgard, a fantasy world corpus with 200+ interconnected documents covering races, characters, locations, events, factions, and lore.

- Use search for both exact entity lookup (mode=text, e.g. "Throndor the Wise") and conceptual queries (mode=semantic, e.g. "ancient magical traditions"). Default mode=rff fuses both.
- After search, use get_content with the returned index to read full chunk text.
- Use get_document to read a complete entry by its source path.
- Use list_documents for navigation or to browse what is available.

Always cite your sources — include the source path and breadcrumb from search results in your answers.`,
		})
		mcpServer.AddReceivingMiddleware(tools.LoggingMiddleware())
		tools.Register(mcpServer, unifiedStore)

		handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
			return mcpServer
		}, nil)

		addr := fmt.Sprintf(":%d", port)
		endpoint := fmt.Sprintf("http://localhost%s/mcp", addr)

		_ = pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(pterm.TableData{
			{"Resource", "Value"},
			{"Documents", strconv.Itoa(len(docs))},
			{"Chunks", strconv.Itoa(len(chunks))},
			{"Embeddings", embeddingsSummary},
			{"Endpoint", endpoint},
		}).Render()

		pterm.Info.Printf("MCP server listening on %s\n", endpoint)
		http.Handle("/mcp", handler)
		return http.ListenAndServe(addr, nil)
	},
}
