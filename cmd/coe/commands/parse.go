package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/phrozen/hybrid-agentic-rag/internal/models"
	"github.com/phrozen/hybrid-agentic-rag/internal/parser"
	"github.com/phrozen/hybrid-agentic-rag/internal/search/semantic"
	"github.com/phrozen/hybrid-agentic-rag/internal/search/text"
	"github.com/spf13/cobra"
)

var (
	inputDir  string
	outputDir string
	runEmbed  bool
)

func init() {
	parseCmd.Flags().StringVarP(&inputDir, "input", "i", "", "Directory path containing input chronicles markdown sources")
	parseCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Directory path where output index artifacts will be saved")
	parseCmd.Flags().BoolVarP(&runEmbed, "embed", "e", false, "Generate and quantize semantic neural embeddings (optional)")

	// Mark required flags
	_ = parseCmd.MarkFlagRequired("input")
	_ = parseCmd.MarkFlagRequired("output")

	RootCmd.AddCommand(parseCmd)
}

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Ingest and build semantic and keyword lookup files",
	Long:  `Parses fantasy chronicles into flat metadata chunks, runs batch embeddings over localhost:8080 (optional), and serializes indexes for serve operations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		totalStart := time.Now()

		// 0. Ensure target output folder exists
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory %q: %w", outputDir, err)
		}

		// 1. Parsing Documents
		fmt.Printf("Parsing documents from path: %s...\n", inputDir)
		parseStart := time.Now()
		docs, err := parser.Parse(inputDir)
		if err != nil {
			return fmt.Errorf("failed to discover and parse files: %w", err)
		}
		parseElapsed := time.Since(parseStart)
		fmt.Printf("  └─ Success: Loaded %d documents [%s]\n\n", len(docs), parseElapsed.Round(time.Millisecond))

		// 2. Saving Documents Register
		docPath := filepath.Join(outputDir, "documents.json")
		fmt.Printf("Saving documents metadata to %s...\n", docPath)
		saveDocsStart := time.Now()
		docData, err := json.MarshalIndent(docs, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal documents json: %w", err)
		}
		if err := os.WriteFile(docPath, docData, 0644); err != nil {
			return fmt.Errorf("failed to write %q: %w", docPath, err)
		}
		saveDocsElapsed := time.Since(saveDocsStart)
		fmt.Printf("  └─ Success: Saved documents registry [%s]\n\n", saveDocsElapsed.Round(time.Millisecond))

		// 3. Extracting and Parsing Chunks
		fmt.Println("Extracting structural text chunks...")
		chunkStart := time.Now()
		chunks := parser.Chunk(docs)
		chunkElapsed := time.Since(chunkStart)
		fmt.Printf("  └─ Success: Segmented into %d flat sequence chunks [%s]\n\n", len(chunks), chunkElapsed.Round(time.Millisecond))

		// 4. Saving Chunks Register
		chunkPath := filepath.Join(outputDir, "chunks.json")
		fmt.Printf("Saving chunks metadata to %s...\n", chunkPath)
		saveChunksStart := time.Now()
		chunkData, err := json.MarshalIndent(chunks, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal chunks json: %w", err)
		}
		if err := os.WriteFile(chunkPath, chunkData, 0644); err != nil {
			return fmt.Errorf("failed to write %q: %w", chunkPath, err)
		}
		saveChunksElapsed := time.Since(saveChunksStart)
		fmt.Printf("  └─ Success: Saved chunks registry [%s]\n\n", saveChunksElapsed.Round(time.Millisecond))

		// 5. Inverted index indexing & saving
		idxPath := filepath.Join(outputDir, "index.bin")
		fmt.Println("Indexing text into keyword proximity engine...")
		iiStart := time.Now()
		ii := text.NewInvertedIndex()
		if err := ii.Index(chunks); err != nil {
			return fmt.Errorf("failed during keyword indexing loop: %w", err)
		}
		iiIndexElapsed := time.Since(iiStart)
		fmt.Printf("  ├─ Indexed keyword coordinates [%s]\n", iiIndexElapsed.Round(time.Millisecond))

		iiSaveStart := time.Now()
		idxBytes, err := ii.MarshalBinary()
		if err != nil {
			return fmt.Errorf("failed to serialize InvertedIndex: %w", err)
		}
		if err := os.WriteFile(idxPath, idxBytes, 0644); err != nil {
			return fmt.Errorf("failed to write %q: %w", idxPath, err)
		}
		iiSaveElapsed := time.Since(iiSaveStart)
		fmt.Printf("  └─ Success: Serialized index.bin (%d bytes) [%s]\n\n", len(idxBytes), iiSaveElapsed.Round(time.Millisecond))

		// 6. Vector/Semantic index indexing & saving (Optional behind runEmbed flag)
		if runEmbed {
			embPath := filepath.Join(outputDir, "embeddings.bin")
			fmt.Println("Generating and quantizing dense semantic vectors...")
			client := semantic.NewClient(nil) // Defaults to local llama-server on localhost:8080
			vi := semantic.NewVectorIndex(client)

			viStart := time.Now()
			if err := vi.Index(chunks); err != nil {
				return fmt.Errorf("failed during semantic chunk embedding loop: %w", err)
			}
			viIndexElapsed := time.Since(viStart)
			fmt.Printf("\n  ├─ Compacted neural vector arrays [%s]\n", viIndexElapsed.Round(time.Millisecond))

			viSaveStart := time.Now()
			embBytes, err := vi.MarshalBinary()
			if err != nil {
				return fmt.Errorf("failed to serialize VectorIndex raw binary: %w", err)
			}
			if err := os.WriteFile(embPath, embBytes, 0644); err != nil {
				return fmt.Errorf("failed to write %q: %w", embPath, err)
			}
			viSaveElapsed := time.Since(viSaveStart)
			fmt.Printf("  └─ Success: Serialized embeddings.bin (%d dims, %d bytes) [%s]\n\n", vi.Dim, len(embBytes), viSaveElapsed.Round(time.Millisecond))
		} else {
			fmt.Println("Skipping semantic embeddings generation (use flag -e / --embed to activate neural parser).")
			fmt.Println()
		}

		// Calculate stats & print report to screen
		printExecutionStats(chunks, totalStart)

		return nil
	},
}

func printExecutionStats(chunks []*models.Chunk, startTime time.Time) {
	if len(chunks) == 0 {
		fmt.Println("No chunks registered to output statistics.")
		return
	}

	smallest := chunks[0]
	largest := chunks[0]
	var totalBytes int64

	for _, ch := range chunks {
		chLen := int64(len(ch.Content))
		totalBytes += chLen

		if chLen < int64(len(smallest.Content)) {
			smallest = ch
		}
		if chLen > int64(len(largest.Content)) {
			largest = ch
		}
	}

	avgSize := float64(totalBytes) / float64(len(chunks))

	fmt.Println("================================================================================")
	fmt.Println("                  CHRONICLES OF AETHELGARD CORPUS PROFILE                       ")
	fmt.Println("================================================================================")
	fmt.Printf("  • Total Chunk Elements  : %d\n", len(chunks))
	fmt.Printf("  • Average Chunk Size    : %0.2f bytes\n", avgSize)
	fmt.Println()
	fmt.Printf("  • Smallest Chunk Profile: %d bytes\n", len(smallest.Content))
	fmt.Printf("    ├─ Heading            : %q (Level %d)\n", smallest.Heading, smallest.HeadingLevel)
	fmt.Printf("    ├─ Range              : lines %d-%d\n", smallest.StartLine, smallest.EndLine)
	fmt.Printf("    └─ Source Location    : %s\n", smallest.Source)
	fmt.Println()
	fmt.Printf("  • Largest Chunk Profile : %d bytes\n", len(largest.Content))
	fmt.Printf("    ├─ Heading            : %q (Level %d)\n", largest.Heading, largest.HeadingLevel)
	fmt.Printf("    ├─ Range              : lines %d-%d\n", largest.StartLine, largest.EndLine)
	fmt.Printf("    └─ Source Location    : %s\n", largest.Source)
	fmt.Println("================================================================================")
	fmt.Printf("Total Elapsed Execution Time: %s\n", time.Since(startTime).Round(time.Millisecond))
	fmt.Println("================================================================================")
	fmt.Println()
}
