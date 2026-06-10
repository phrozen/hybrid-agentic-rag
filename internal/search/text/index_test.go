package text

import (
	"testing"

	"github.com/phrozen/hybrid-agentic-rag/internal/models"
)

func TestInvertedIndex_IndexAndSearch(t *testing.T) {
	ii := NewInvertedIndex()

	chunks := []*models.Chunk{
		{Content: "The absolute coordinates of Silverhaven capitals"},
		{Content: "Dwarven miner blacksmiths working in Stonehold mountains"},
		{Content: "The Silvermoon Academy of Starry Magic"},
	}

	err := ii.Index(chunks)
	if err != nil {
		t.Fatalf("Failed to Index chunks: %v", err)
	}

	// Search for "Silverhaven"
	hits := ii.Search("Silverhaven", 2)
	if len(hits) == 0 {
		t.Fatalf("Expected results, got 0")
	}

	if hits[0].Index != 0 {
		t.Errorf("Expected leading match to be index 0, got %d", hits[0].Index)
	}

	// Search for "Starry Magic"
	hits2 := ii.Search("Starry Magic", 2)
	t.Logf("Hits in 'Starry Magic':")
	for i, hit := range hits2 {
		t.Logf("  #%d: Index=%d, Score=%f", i, hit.Index, hit.Score)
	}
	if len(hits2) == 0 {
		t.Fatalf("Expected results, got 0")
	}

	if hits2[0].Index != 2 {
		t.Errorf("Expected leading match to be index 2, got %d", hits2[0].Index)
	}
}

func TestInvertedIndex_Serialization(t *testing.T) {
	ii := NewInvertedIndex()

	chunks := []*models.Chunk{
		{Content: "Silverhaven wizard mages"},
		{Content: "Stonehold steel weapons"},
	}

	if err := ii.Index(chunks); err != nil {
		t.Fatalf("Failed to Index: %v", err)
	}

	data, err := ii.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to Marshal: %v", err)
	}

	iiDecoded := NewInvertedIndex()
	if err := iiDecoded.UnmarshalBinary(data); err != nil {
		t.Fatalf("Failed to Unmarshal: %v", err)
	}

	hits := iiDecoded.Search("Stonehold", 2)
	if len(hits) != 1 {
		t.Fatalf("Expected exactly 1 hit, got %d", len(hits))
	}

	if hits[0].Index != 1 {
		t.Errorf("Expected hit to map back to chunk 1, got %d", hits[0].Index)
	}
}
