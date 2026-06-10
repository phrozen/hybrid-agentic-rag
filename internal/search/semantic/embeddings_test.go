package semantic

import (
	"math"
	"testing"
)

func TestClient_Embed(t *testing.T) {
	// Initialize client using nil configuration (defaults to localhost:8080)
	client := NewClient(nil)

	// Test inputs
	inputs := []string{
		"Hello Chronicles of Aethelgard",
		"Magic systems and ancient runes",
	}

	embeddings, err := client.Embed(inputs)
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}

	if len(embeddings) != len(inputs) {
		t.Errorf("Expected %d embeddings, got %d", len(inputs), len(embeddings))
	}

	for i, emb := range embeddings {
		if len(emb) == 0 {
			t.Errorf("Embedding at index %d is empty", i)
		} else {
			// Calculate L2 Norm: sqrt(sum of squares)
			var sumSquares float64
			for _, val := range emb {
				sumSquares += float64(val * val)
			}
			l2Norm := math.Sqrt(sumSquares)

			t.Logf("Embedding at index %d has dimension %d, L2 Norm = %f", i, len(emb), l2Norm)
		}
	}
}

func TestClient_NormalizationEquivalence(t *testing.T) {
	client := NewClient(nil)

	inputs := []string{"High Elves study runemagic within the Arcane College of Thaldor"}
	embeddings, err := client.Embed(inputs)
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}

	if len(embeddings) == 0 {
		t.Fatalf("Expected at least one embedding response")
	}

	original := embeddings[0]

	// Run our helper l2NormalizeFloat32 on the returned vector
	normalized := l2NormalizeFloat32(original)

	// Since Jina v5 embedding response is already L2-normalized, the original vector
	// and the normalized output from our helper should be practically identical
	// (difference in every coordinate should be virtually zero).
	if len(original) != len(normalized) {
		t.Fatalf("Vector length mismatch: original=%d, normalized=%d", len(original), len(normalized))
	}

	var maxDiff float64
	for i := 0; i < len(original); i++ {
		diff := math.Abs(float64(original[i] - normalized[i]))
		if diff > maxDiff {
			maxDiff = diff
		}
	}

	const threshold = 1e-6
	if maxDiff > threshold {
		t.Errorf("Embedding and normalized vector differ significantly: max difference = %e", maxDiff)
	} else {
		t.Logf("Normalization equivalence verified! Max coordinate difference = %e", maxDiff)
	}
}
