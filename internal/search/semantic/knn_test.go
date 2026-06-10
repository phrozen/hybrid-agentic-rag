package semantic

import (
	"testing"
)

func TestKNN_Search(t *testing.T) {
	// Let's mock 5 candidate vectors of dimension 8, flattened into a single contiguous []int8 array.
	dim := 8
	rawVectors := [][]int8{
		l2NormalizeInt8([]int8{-1, -1, -1, -1, 1, 1, 1, 1}),     // unrelated_1 (idx = 0)
		l2NormalizeInt8([]int8{-1, -1, -1, -1, -1, -1, -1, -1}), // unrelated_2 (idx = 1)
		l2NormalizeInt8([]int8{1, -1, 1, -1, 1, 1, 1, 1}),       // match_third (idx = 2)
		l2NormalizeInt8([]int8{1, 1, 1, 1, 1, 1, 1, 1}),         // match_first (idx = 3)
		l2NormalizeInt8([]int8{1, 0, 1, 1, 1, 1, 1, 1}),         // match_second (idx = 4)
	}

	// Flatten vectors sequentially into a single slice to avoid pointer chasing
	flatVectors := make([]int8, 0, len(rawVectors)*dim)
	for _, vec := range rawVectors {
		flatVectors = append(flatVectors, vec...)
	}

	query := l2NormalizeInt8([]int8{1, 1, 1, 1, 1, 1, 1, 1}) // Query points in perfect positive direction

	// Search for top 3 matches
	k := 3
	results := knn(flatVectors, query, dim, k)

	if len(results) != k {
		t.Fatalf("Expected exactly %d results, got %d", k, len(results))
	}

	// Verify that index 3 ("match_first") is top 1, with highest similarity
	if results[0].Index != 3 {
		t.Errorf("Expected leading match to be index 3, got %d", results[0].Index)
	}

	// Verify that index 4 ("match_second") is second
	if results[1].Index != 4 {
		t.Errorf("Expected second match to be index 4, got %d", results[1].Index)
	}

	// Verify that index 2 ("match_third") is third
	if results[2].Index != 2 {
		t.Errorf("Expected third match to be index 2, got %d", results[2].Index)
	}

	// Output candidate list
	t.Logf("Rankings returned successfully:")
	for i, res := range results {
		t.Logf("  #%d: Index=%d, Similarity=%0.6f", i+1, res.Index, res.Score)
	}
}
