package search_test

import (
	"sort"
	"testing"

	"github.com/phrozen/hybrid-agentic-rag/internal/models"
	"github.com/phrozen/hybrid-agentic-rag/internal/search"
	"github.com/phrozen/hybrid-agentic-rag/internal/search/semantic"
	"github.com/phrozen/hybrid-agentic-rag/internal/search/text"
)

// MockEmbedder serves as a mock client for Jina / v1 / embeddings responses.
type MockEmbedder struct {
	embedFunc func(texts []string) ([][]float32, error)
}

func (m *MockEmbedder) Embed(texts []string) ([][]float32, error) {
	if m.embedFunc != nil {
		return m.embedFunc(texts)
	}
	// Return arbitrary zero embeddings of dimension 4
	out := make([][]float32, len(texts))
	for i := range out {
		out[i] = []float32{0.1, -0.2, 0.5, 0.9}
	}
	return out, nil
}

// TestRFFBasic verifies merging logic under balanced weights.
func TestRFFBasic(t *testing.T) {
	// Sorted text inputs
	textHits := []search.Hit{
		{Index: 2, Score: 12.1},
		{Index: 1, Score: 2.5},
	}

	// Sorted semantic inputs
	semanticHits := []search.Hit{
		{Index: 1, Score: 0.9},
		{Index: 3, Score: -0.1},
	}

	fused := search.RFF(textHits, semanticHits)

	// Verify outputs are ranked correctly.
	// Index 1 was rank 2 (text, score 2.5) and rank 1 (semantic, score 0.9)
	// Index 2 was rank 1 (text, score 12.1)
	// Index 3 was rank 2 (semantic, score -0.1)
	if len(fused) != 3 {
		t.Fatalf("expected 3 fused results, got %d", len(fused))
	}

	if fused[0].Index != 1 {
		t.Errorf("expected top hit to be index 1, got index %d", fused[0].Index)
	}

	// Double check score ordering
	for i := 0; i < len(fused)-1; i++ {
		if fused[i].Score < fused[i+1].Score {
			t.Errorf("fused results not sorted in descending order: index %d score (%f) < index %d score (%f)",
				fused[i].Index, fused[i].Score, fused[i+1].Index, fused[i+1].Score)
		}
	}
}

// TestTextSearchSorting verifies that text.InvertedIndex.Search returns results sorted descending by score.
func TestTextSearchSorting(t *testing.T) {
	idx := text.NewInvertedIndex()
	chunks := []*models.Chunk{
		{Content: "apple bananna orange core"},
		{Content: "apple apple and more apple things"},
		{Content: "completely unrelated text query words"},
	}
	if err := idx.Index(chunks); err != nil {
		t.Fatalf("failed to index chunks: %v", err)
	}

	hits := idx.Search("apple", 10)
	if len(hits) == 0 {
		t.Fatal("expected hits back but got none")
	}

	// Assert that hits are strictly sorted descending by Score
	isSorted := sort.SliceIsSorted(hits, func(i, j int) bool {
		return hits[i].Score > hits[j].Score
	})
	if !isSorted {
		t.Errorf("text index Search results not descendingly sorted by score: %+v", hits)
	}
}

// TestSemanticSearchSorting verifies that semantic.VectorIndex.Search returns results sorted descending by score.
func TestSemanticSearchSorting(t *testing.T) {
	// Set mock response where index 1 gets a high similarity match and index 0 gets lower
	mockClient := &MockEmbedder{
		embedFunc: func(texts []string) ([][]float32, error) {
			if len(texts) == 1 && texts[0] == "Query: target" {
				return [][]float32{{0.5, 0.5, 0.5, 0.5}}, nil
			}
			// When indexing elements:
			// chunk 0: lower similarity
			// chunk 1: higher similarity
			res := make([][]float32, len(texts))
			for i, text := range texts {
				if i == 1 || text == "Document: \nContext: \n\nbanana" {
					res[i] = []float32{0.4, 0.4, 0.4, 0.4}
				} else {
					res[i] = []float32{0.1, 0.1, 0.1, 0.1}
				}
			}
			return res, nil
		},
	}

	vi := semantic.NewVectorIndex(mockClient)
	chunks := []*models.Chunk{
		{Content: "apple"},
		{Content: "banana"},
	}

	if err := vi.Index(chunks); err != nil {
		t.Fatalf("indexing failed: %v", err)
	}

	hits := vi.Search("target", 2)
	if len(hits) == 0 {
		t.Fatal("expected hits back but got none")
	}

	// Assert that hits are strictly sorted descending by Score
	isSorted := sort.SliceIsSorted(hits, func(i, j int) bool {
		return hits[i].Score > hits[j].Score
	})
	if !isSorted {
		t.Errorf("semantic index Search results not descendingly sorted by score: %+v", hits)
	}
}
