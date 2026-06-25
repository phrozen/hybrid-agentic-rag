package search

import (
	"sort"

	"github.com/phrozen/hybrid-agentic-rag/internal/models"
)

const (
	RFFK           = 60
	TextWeight     = 0.5
	SemanticWeight = 0.5
)

type SearchMode int

const (
	SearchRFF SearchMode = iota
	SearchText
	SearchSemantic
)

type Hit struct {
	Index int
	Score float32
}

type SearchResult struct {
	Index      int
	Score      float32
	Source     string
	Breadcrumb string
}

// ProgressFunc reports a human-readable progress step during indexing.
type ProgressFunc func(step string)

type Indexer interface {
	Index(chunks []*models.Chunk, progress ProgressFunc) error
	Search(query string, k int) []Hit
	MarshalBinary() ([]byte, error)
	UnmarshalBinary(data []byte) error
}

// RFF merges pre-sorted text and semantic search hits using Reciprocal Rank Fusion (RFF).
// It returns a unified slice of Hits ranked in descending order of their fused scores.
func RFF(text, semantic []Hit) []Hit {
	fusedScores := make(map[int]float32)

	// 1. Process Keyword search hits
	for rank, hit := range text {
		reciprocalRank := 1.0 / float32(RFFK+(rank+1))
		fusedScores[hit.Index] += TextWeight * reciprocalRank
	}

	// 2. Process Dense vector search hits
	for rank, hit := range semantic {
		reciprocalRank := 1.0 / float32(RFFK+(rank+1))
		fusedScores[hit.Index] += SemanticWeight * reciprocalRank
	}

	// 3. Convert map back into a flat slice of Hits
	results := make([]Hit, 0, len(fusedScores))
	for chunkIdx, score := range fusedScores {
		results = append(results, Hit{
			Index: chunkIdx,
			Score: score,
		})
	}

	// 4. Sort the finalized hits descending by RFF score, resolving ties deterministic by index
	sort.Slice(results, func(i, j int) bool {
		if results[i].Score == results[j].Score {
			return results[i].Index < results[j].Index
		}
		return results[i].Score > results[j].Score
	})

	return results
}
