package parser

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/phrozen/hybrid-agentic-rag/internal/models"
)

// validExtensions tracks valid markdown file extensions supported by this parser.
var validExtensions = []string{".md", ".markdown", ".mdown", ".mkd", ".mdx"}

// Parse recursively walks the root folder, parsing any discovered markdown files
// and loading them into memory. It returns a slice of documents.
func Parse(root string) ([]*models.Document, error) {
	var files []string

	// Walk the directory tree to collect all valid markdown files
	err := filepath.WalkDir(root, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if dir.IsDir() {
			return nil
		}

		// Save the path if its file extension is valid
		if slices.Contains(validExtensions, filepath.Ext(path)) {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Pre-allocate the documented slice to avoid reallocation penalties
	docs := make([]*models.Document, 0, len(files))

	// Load each document in turn
	for _, file := range files {
		source, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		d := &models.Document{
			Source:  file,
			Content: source,
			Size:    len(source),
			// Since bytes.Count counts delimiters, adding 1 gives us the line count.
			Lines: bytes.Count(source, []byte{'\n'}) + 1,
		}
		docs = append(docs, d)
	}

	return docs, nil
}

// Chunk processes all documents in the slice, splitting them into logical chunks at heading boundaries
// and flattening them into a single ordered sequence.
func Chunk(docs []*models.Document) []*models.Chunk {
	var allChunks []*models.Chunk

	for _, d := range docs {
		// Split raw byte content by LF
		lines := bytes.Split(d.Content, []byte{'\n'})

		// Track current heading hierarchy
		var h1, h2, h3 string
		var currentChunkStart int
		var currentChunkContent [][]byte // Use byte slices to defer string creation
		var currentHeading string
		var currentHeadingLevel int

		// Define prefixes as byte slices for zero-allocation checks
		prefixH1 := []byte("# ")
		prefixH2 := []byte("## ")
		prefixH3 := []byte("### ")

		// Helper to flush current chunk
		flushChunk := func(endLine int) {
			if len(currentChunkContent) == 0 {
				return
			}

			// Join and trim the gathered bytes
			joinedBytes := bytes.Join(currentChunkContent, []byte{'\n'})
			trimmedBytes := bytes.TrimSpace(joinedBytes)

			// Create chunk mapping only if the content is not empty
			if len(trimmedBytes) > 0 {
				breadcrumb := buildBreadcrumb(h1, h2, h3)
				allChunks = append(allChunks, &models.Chunk{
					Source:       d.Source,
					Breadcrumb:   breadcrumb,
					Heading:      currentHeading,
					HeadingLevel: currentHeadingLevel,
					Content:      string(trimmedBytes),
					StartLine:    currentChunkStart,
					EndLine:      endLine,
				})
			}

			currentChunkContent = [][]byte{}
		}

		// Process each line sequence
		for i, lineBytes := range lines {
			lineNum := i + 1 // 1-indexed

			// Strip carriage return at the end of the line (handles CRLF cleanly)
			lineBytes = bytes.TrimSuffix(lineBytes, []byte{'\r'})

			// Identify heading level by checking prefixes
			if bytes.HasPrefix(lineBytes, prefixH1) {
				// Found H1 heading: flush the currently accumulating chunk
				flushChunk(lineNum - 1)
				h1 = string(bytes.TrimSpace(bytes.TrimPrefix(lineBytes, prefixH1)))
				h2 = ""
				h3 = ""
				currentHeading = h1
				currentHeadingLevel = 1
				currentChunkStart = lineNum
				currentChunkContent = [][]byte{}
			} else if bytes.HasPrefix(lineBytes, prefixH2) {
				// Found H2 heading: flush the currently accumulating chunk
				flushChunk(lineNum - 1)
				h2 = string(bytes.TrimSpace(bytes.TrimPrefix(lineBytes, prefixH2)))
				h3 = ""
				currentHeading = h2
				currentHeadingLevel = 2
				currentChunkStart = lineNum
				currentChunkContent = [][]byte{}
			} else if bytes.HasPrefix(lineBytes, prefixH3) {
				// Found H3 heading: flush the currently accumulating chunk
				flushChunk(lineNum - 1)
				h3 = string(bytes.TrimSpace(bytes.TrimPrefix(lineBytes, prefixH3)))
				currentHeading = h3
				currentHeadingLevel = 3
				currentChunkStart = lineNum
				currentChunkContent = [][]byte{}
			} else {
				// Treat as standard content line
				// Initialize the start line for files that have content before any headings are encountered
				if currentChunkStart == 0 {
					currentChunkStart = lineNum
				}
				currentChunkContent = append(currentChunkContent, lineBytes)
			}
		}

		// Flush the final chunk at the end of the file
		flushChunk(len(lines))
	}

	return allChunks
}

// buildBreadcrumb creates a hierarchical breadcrumb string from the given h1, h2, and h3 headings.
// It is optimized to use strings.Builder to minimize heap allocations during string concatenation.
func buildBreadcrumb(h1, h2, h3 string) string {
	if h1 == "" && h2 == "" && h3 == "" {
		return ""
	}

	var sb strings.Builder
	// Pre-allocate estimated capacity
	totalLen := len(h1) + len(h2) + len(h3)
	if h1 != "" && h2 != "" {
		totalLen += 3 // for " > "
	}
	if (h1 != "" || h2 != "") && h3 != "" {
		totalLen += 3 // for " > "
	}
	sb.Grow(totalLen)

	if h1 != "" {
		sb.WriteString(h1)
	}
	if h2 != "" {
		if sb.Len() > 0 {
			sb.WriteString(" > ")
		}
		sb.WriteString(h2)
	}
	if h3 != "" {
		if sb.Len() > 0 {
			sb.WriteString(" > ")
		}
		sb.WriteString(h3)
	}

	return sb.String()
}
