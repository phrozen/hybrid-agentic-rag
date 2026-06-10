package semantic

import (
	"fmt"
	"testing"
)

func TestQuantization_SimilarityLoss(t *testing.T) {
	// Initialize embedding client using nil configuration (defaults to localhost:8080)
	client := NewClient(nil)

	// Since we are using Jina-embeddings-v5 (with symmetric bi-directional training),
	// we prefix texts with "Query:" for queries and "Document:" for document chunks to trigger correct pooling.
	rawQueries := []string{
		"High Elves and magical runic spells in Silvermoon",
		"Elven weavers practicing textile arts in Silvermoon",
		"Dwarven mining mechanics and blacksmith guilds in Stonehold",
	}

	// Format inputs according to Jina v5 retrieval prompts
	// For asymmetric retrieval optimization: index 0 acts as the "Query", others are candidate "Documents".
	inputs := []string{
		"Query: " + rawQueries[0],
		"Document: " + rawQueries[1],
		"Document: " + rawQueries[2],
	}

	embeddings, err := client.Embed(inputs)
	if err != nil {
		t.Fatalf("Failed to fetch test embeddings: %v", err)
	}

	// We'll perform asymmetric similarity comparisons:
	// Pair 0 (Highly Related): 0 vs 1 (Elven magic query vs Elven weavers document)
	// Pair 1 (Unrelated): 0 vs 2 (Elven magic query vs Dwarven mining document)
	pairs := []struct {
		name string
		idx1 int
		idx2 int
	}{
		{"Query: Elven Magic vs Doc: Elven Weavings (Related)", 0, 1},
		{"Query: Elven Magic vs Doc: Dwarven Mining (Unrelated)", 0, 2},
	}

	totalDims := len(embeddings[0])

	fmt.Printf("\n--- JINA RETRIEVAL QUANTIZATION COMPARISON (Dimensions: %d) ---\n\n", totalDims)

	for _, p := range pairs {
		v1 := embeddings[p.idx1]
		v2 := embeddings[p.idx2]

		// 1. Reference similarity (Float32 Dot Product embedded directly)
		var floatSim float32
		for i := 0; i < len(v1); i++ {
			floatSim += v1[i] * v2[i]
		}

		// 2. Int8 Quantization
		i8_1 := quantize(v1)
		i8_2 := quantize(v2)
		int8Sim := similarity(i8_1, i8_2)
		int8Loss := floatSim - int8Sim

		fmt.Printf("Query Pair: %s\n", p.name)
		fmt.Printf("  • Original Float32 Cosine Similarity : %0.6f\n", floatSim)
		fmt.Printf("  • Int8 Quantized Cosine Similarity   : %0.6f (Loss: %0.6f, Rel Error: %0.2f%%)\n",
			int8Sim, int8Loss, float64(mathAbs(int8Loss)/floatSim)*100.0)
		fmt.Println()
	}
}

func mathAbs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}
