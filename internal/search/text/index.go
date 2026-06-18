package text

import (
	"fmt"
	"sort"

	"github.com/phrozen/hybrid-agentic-rag/internal/models"
	"github.com/phrozen/hybrid-agentic-rag/internal/search"
	"github.com/wizenheimer/blaze"
)

// Ensure InvertedIndex implements search.Indexer at compile time.
var _ search.Indexer = (*InvertedIndex)(nil)

// InvertedIndex wraps the high-performance blaze.InvertedIndex implementation
// to fit seamlessly into our unified search pipeline.
type InvertedIndex struct {
	index *blaze.InvertedIndex
}

// NewInvertedIndex instantiates a clean keyword search index.
func NewInvertedIndex() *InvertedIndex {
	return &InvertedIndex{
		index: blaze.NewInvertedIndex(),
	}
}

// Index populates the inverted index with chunks.
func (ii *InvertedIndex) Index(chunks []*models.Chunk) error {
	// Re-initialize index on bulk-load execution
	ii.index = blaze.NewInvertedIndex()

	for i, ch := range chunks {
		// =================================================================
		// USER FINE-TUNING AREA (KEYWORD INPUT PREPARATION)
		// You can fine-tune what fields are indexed by the keyword engine
		// (e.g. adding headings to help direct matching, filtering stop-words,
		// or indexing raw body context without heritance pollution etc.)
		// =================================================================

		preparedText := fmt.Sprintf("%s\n%s", ch.Heading, ch.Content)

		// =================================================================

		// Index the document using the flat chunk index (i) as the unique DocID
		ii.index.Index(i, preparedText)
	}

	return nil
}

// Search queries the inverted index using BM25 relevance ranking
// (IDF weighting + document-length normalization).
func (ii *InvertedIndex) Search(query string, k int) []search.Hit {
	if ii.index == nil {
		return nil
	}

	matches := ii.index.RankBM25(query, k)
	hits := make([]search.Hit, len(matches))
	for i, m := range matches {
		hits[i] = search.Hit{
			Index: m.DocID,
			Score: float32(m.Score),
		}
	}

	// Always sort hits descendingly by score to guarantee sorted output for upstream callers/RFF.
	sort.Slice(hits, func(i, j int) bool {
		if hits[i].Score == hits[j].Score {
			return hits[i].Index < hits[j].Index
		}
		return hits[i].Score > hits[j].Score
	})

	return hits
}

// MarshalBinary serializes the inverted index into a binary payload using blaze's protocol.
func (ii *InvertedIndex) MarshalBinary() ([]byte, error) {
	if ii.index == nil {
		return nil, nil
	}
	return ii.index.Encode()
}

// UnmarshalBinary restores the inverted index from a raw binary payload.
func (ii *InvertedIndex) UnmarshalBinary(data []byte) error {
	ii.index = blaze.NewInvertedIndex()
	return ii.index.Decode(data)
}
