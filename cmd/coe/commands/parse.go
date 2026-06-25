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
	"github.com/pterm/pterm"
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
		sp, _ := pterm.DefaultSpinner.Start("Parsing documents from " + inputDir)
		parseStart := time.Now()
		docs, err := parser.Parse(inputDir)
		if err != nil {
			sp.Fail("Failed to discover and parse files")
			return fmt.Errorf("failed to discover and parse files: %w", err)
		}
		sp.Success(fmt.Sprintf("Loaded %d documents [%s]", len(docs), time.Since(parseStart).Round(time.Millisecond)))

		// 2. Saving Documents Register
		docPath := filepath.Join(outputDir, "documents.json")
		sp, _ = pterm.DefaultSpinner.Start("Saving documents metadata to " + docPath)
		saveDocsStart := time.Now()
		docData, err := json.MarshalIndent(docs, "", "  ")
		if err != nil {
			sp.Fail("Failed to marshal documents json")
			return fmt.Errorf("failed to marshal documents json: %w", err)
		}
		if err := os.WriteFile(docPath, docData, 0644); err != nil {
			sp.Fail("Failed to write documents.json")
			return fmt.Errorf("failed to write %q: %w", docPath, err)
		}
		sp.Success(fmt.Sprintf("Saved documents registry [%s]", time.Since(saveDocsStart).Round(time.Millisecond)))

		// 3. Extracting and Parsing Chunks
		sp, _ = pterm.DefaultSpinner.Start("Extracting structural text chunks")
		chunkStart := time.Now()
		chunks := parser.Chunk(docs)
		sp.Success(fmt.Sprintf("Segmented into %d flat sequence chunks [%s]", len(chunks), time.Since(chunkStart).Round(time.Millisecond)))

		// 4. Saving Chunks Register
		chunkPath := filepath.Join(outputDir, "chunks.json")
		sp, _ = pterm.DefaultSpinner.Start("Saving chunks metadata to " + chunkPath)
		saveChunksStart := time.Now()
		chunkData, err := json.MarshalIndent(chunks, "", "  ")
		if err != nil {
			sp.Fail("Failed to marshal chunks json")
			return fmt.Errorf("failed to marshal chunks json: %w", err)
		}
		if err := os.WriteFile(chunkPath, chunkData, 0644); err != nil {
			sp.Fail("Failed to write chunks.json")
			return fmt.Errorf("failed to write %q: %w", chunkPath, err)
		}
		sp.Success(fmt.Sprintf("Saved chunks registry [%s]", time.Since(saveChunksStart).Round(time.Millisecond)))

		// 5. Inverted index indexing & saving
		idxPath := filepath.Join(outputDir, "index.bin")
		sp, _ = pterm.DefaultSpinner.Start("Indexing text into keyword proximity engine")
		iiStart := time.Now()
		ii := text.NewInvertedIndex()
		if err := ii.Index(chunks, nil); err != nil {
			sp.Fail("Failed during keyword indexing loop")
			return fmt.Errorf("failed during keyword indexing loop: %w", err)
		}
		idxBytes, err := ii.MarshalBinary()
		if err != nil {
			sp.Fail("Failed to serialize InvertedIndex")
			return fmt.Errorf("failed to serialize InvertedIndex: %w", err)
		}
		if err := os.WriteFile(idxPath, idxBytes, 0644); err != nil {
			sp.Fail("Failed to write index.bin")
			return fmt.Errorf("failed to write %q: %w", idxPath, err)
		}
		sp.Success(fmt.Sprintf("Serialized index.bin (%d bytes) [%s]", len(idxBytes), time.Since(iiStart).Round(time.Millisecond)))

		// 6. Vector/Semantic index indexing & saving (Optional behind runEmbed flag)
		if runEmbed {
			embPath := filepath.Join(outputDir, "embeddings.bin")
			sp, _ = pterm.DefaultSpinner.Start("Generating and quantizing dense semantic vectors")
			client := semantic.NewClient(nil) // Defaults to local llama-server on localhost:8080
			vi := semantic.NewVectorIndex(client)

			viStart := time.Now()
			if err := vi.Index(chunks, func(step string) { sp.UpdateText(step) }); err != nil {
				sp.Fail("Failed during semantic chunk embedding loop")
				return fmt.Errorf("failed during semantic chunk embedding loop: %w", err)
			}
			embBytes, err := vi.MarshalBinary()
			if err != nil {
				sp.Fail("Failed to serialize VectorIndex raw binary")
				return fmt.Errorf("failed to serialize VectorIndex raw binary: %w", err)
			}
			if err := os.WriteFile(embPath, embBytes, 0644); err != nil {
				sp.Fail("Failed to write embeddings.bin")
				return fmt.Errorf("failed to write %q: %w", embPath, err)
			}
			sp.Success(fmt.Sprintf("Serialized embeddings.bin (%d dims, %d bytes) [%s]", vi.Dim, len(embBytes), time.Since(viStart).Round(time.Millisecond)))
		} else {
			pterm.Warning.Println("Skipping semantic embeddings generation (use flag -e / --embed to activate neural parser).")
		}

		// Calculate stats & print report to screen
		printExecutionStats(chunks, totalStart)

		return nil
	},
}

func printExecutionStats(chunks []*models.Chunk, startTime time.Time) {
	if len(chunks) == 0 {
		pterm.Warning.Println("No chunks registered to output statistics.")
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

	pterm.DefaultSection.Println("Chronicles of Aethelgard Corpus Profile")

	_ = pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(pterm.TableData{
		{"Metric", "Value"},
		{"Total Chunk Elements", fmt.Sprintf("%d", len(chunks))},
		{"Average Chunk Size", fmt.Sprintf("%0.2f bytes", avgSize)},
		{"Smallest Chunk", fmt.Sprintf("%d bytes", len(smallest.Content))},
		{"  Heading", fmt.Sprintf("%q (Level %d)", smallest.Heading, smallest.HeadingLevel)},
		{"  Range", fmt.Sprintf("lines %d-%d", smallest.StartLine, smallest.EndLine)},
		{"  Source", smallest.Source},
		{"Largest Chunk", fmt.Sprintf("%d bytes", len(largest.Content))},
		{"  Heading", fmt.Sprintf("%q (Level %d)", largest.Heading, largest.HeadingLevel)},
		{"  Range", fmt.Sprintf("lines %d-%d", largest.StartLine, largest.EndLine)},
		{"  Source", largest.Source},
		{"Total Elapsed", time.Since(startTime).Round(time.Millisecond).String()},
	}).Render()
}
