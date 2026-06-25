package semantic

import (
	"errors"
	"testing"

	"github.com/phrozen/hybrid-agentic-rag/internal/models"
)

// mockEmbedder acts as a test-isolation double to verify VectorIndex Search behavior off-line
type mockEmbedder struct {
	embedFunc func(texts []string) ([][]float32, error)
}

func (m *mockEmbedder) Embed(texts []string) ([][]float32, error) {
	if m.embedFunc != nil {
		return m.embedFunc(texts)
	}
	return nil, errors.New("unimplemented mock embed method")
}

func TestVectorIndex_Serialization_RoundTrip(t *testing.T) {
	// Initialize a realistic Index structure
	vi := NewVectorIndex(nil)
	vi.Dim = 4
	vi.Vectors = []int8{
		10, -20, 30, -40,
		50, 60, -70, 80,
		-90, 100, -110, 120,
	}

	data, err := vi.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal VectorIndex: %v", err)
	}

	// Unmarshal into a clean VectorIndex structure
	viDecoded := NewVectorIndex(nil)
	err = viDecoded.UnmarshalBinary(data)
	if err != nil {
		t.Fatalf("Failed to unmarshal VectorIndex: %v", err)
	}

	if viDecoded.Dim != vi.Dim {
		t.Errorf("Expected Dim %d, got %d", vi.Dim, viDecoded.Dim)
	}

	if len(viDecoded.Vectors) != len(vi.Vectors) {
		t.Fatalf("Expected length %d, got %d", len(vi.Vectors), len(viDecoded.Vectors))
	}

	for i, val := range vi.Vectors {
		if viDecoded.Vectors[i] != val {
			t.Errorf("Element mismatch at index %d: expected %d, got %d", i, val, viDecoded.Vectors[i])
		}
	}
}

func TestVectorIndex_Serialization_Empty(t *testing.T) {
	vi := NewVectorIndex(nil)
	data, err := vi.MarshalBinary()
	if err != nil {
		t.Fatalf("Marshal target empty failed: %v", err)
	}
	if len(data) != 0 {
		t.Errorf("Expected empty byte slice and no error, got len = %d", len(data))
	}

	vi.Dim = 1024
	// Vectors is nil/empty
	data, err = vi.MarshalBinary()
	if err != nil {
		t.Fatalf("Marshal empty vectors failed: %v", err)
	}
	if len(data) != 0 {
		t.Errorf("Expected empty byte slice, got len = %d", len(data))
	}
}

func TestVectorIndex_Unmarshal_CorruptedOrShort(t *testing.T) {
	vi := NewVectorIndex(nil)

	// 1. Data completely empty or shorter than header length (8 bytes)
	err := vi.UnmarshalBinary([]byte{1, 2, 3, 4})
	if !errors.Is(err, ErrInvalidData) {
		t.Errorf("Expected ErrInvalidData when binary data too short, got: %v", err)
	}

	// 2. Data has header, but indicates body length of 100 on only 2 extra bytes
	// Header: Dim = 4 (little-endian: 4, 0, 0, 0), VecDataLen = 100 (little-endian: 100, 0, 0, 0)
	corruptPayload := []byte{4, 0, 0, 0, 100, 0, 0, 0, 99, 99}
	err = vi.UnmarshalBinary(corruptPayload)
	if !errors.Is(err, ErrInvalidData) {
		t.Errorf("Expected ErrInvalidData when payload smaller than indicated size, got: %v", err)
	}
}

func TestVectorIndex_Unmarshal_TrailBytes(t *testing.T) {
	vi := NewVectorIndex(nil)
	vi.Dim = 2
	vi.Vectors = []int8{10, 20}

	data, err := vi.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Append some trailing payload bytes (e.g. from a compound file stream)
	corruptData := append(data, []byte{99, 99, 99, 99}...)

	viDecoded := NewVectorIndex(nil)
	err = viDecoded.UnmarshalBinary(corruptData)
	if err != nil {
		t.Fatalf("Unmarshal failed with trailing data: %v", err)
	}

	// Decoded length must adhere strictly to header specifications (i.e. length of 2)
	// and NOT include the trailing garbage bytes.
	if len(viDecoded.Vectors) != 2 {
		t.Fatalf("Expected exactly 2 vectors, got %d", len(viDecoded.Vectors))
	}
	if viDecoded.Vectors[0] != 10 || viDecoded.Vectors[1] != 20 {
		t.Errorf("Expected [10, 20], got %v", viDecoded.Vectors)
	}
}

func TestVectorIndex_Unmarshal_MemoryIsolation(t *testing.T) {
	vi := NewVectorIndex(nil)
	vi.Dim = 2
	vi.Vectors = []int8{10, 20}

	data, err := vi.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	viDecoded := NewVectorIndex(nil)
	err = viDecoded.UnmarshalBinary(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Mutate index 1 byte in the source serialized byte slice
	// Index 8 is where the body starts
	data[8] = 99

	// If properly memory-isolated, modifying the source slice must not impact viDecoded.Vectors
	if viDecoded.Vectors[0] != 10 {
		t.Errorf("Memory leakage observed! Source modification bypassed boundary safety (got %d instead of 10)", viDecoded.Vectors[0])
	}
}

func TestVectorIndex_Search(t *testing.T) {
	// Instantiate mock embedder that outputs float32 unit-vectors matching dimension size of 3
	mockCl := &mockEmbedder{
		embedFunc: func(texts []string) ([][]float32, error) {
			// Mock high-value positive vector matching Index 2 coordinate direction
			// Since we mock floating ranges, [1.0, 1.0, 1.0] normalized works beautifully
			return [][]float32{{1.0 / 1.73205, 1.0 / 1.73205, 1.0 / 1.73205}}, nil
		},
	}

	vi := NewVectorIndex(mockCl)
	vi.Dim = 3
	// Mock three vectors of dimension 3, flat-aligned.
	// Vector 0 (unrelated): [-100, -100, -100] scale mapped
	// Vector 1 (match 2nd): [50, 50, 50] scale mapped
	// Vector 2 (match 1st): [120, 120, 120] scale mapped
	vi.Vectors = []int8{
		-100, -100, -100,
		50, 50, 50,
		120, 120, 120,
	}

	hits := vi.Search("test-query-text", 2)
	if len(hits) != 2 {
		t.Fatalf("Expected exactly 2 hits, got %d", len(hits))
	}

	// Highest similarity score first (index 2)
	if hits[0].Index != 2 {
		t.Errorf("Expected top hit to be index 2, got %d", hits[0].Index)
	}

	// Second query match (index 1)
	if hits[1].Index != 1 {
		t.Errorf("Expected second hit to be index 1, got %d", hits[1].Index)
	}
}

func TestVectorIndex_Index(t *testing.T) {
	mockCl := &mockEmbedder{
		embedFunc: func(texts []string) ([][]float32, error) {
			result := make([][]float32, len(texts))
			for i := range texts {
				result[i] = []float32{0.1, -0.2, 0.3}
			}
			return result, nil
		},
	}

	vi := NewVectorIndex(mockCl)
	chunks := []*models.Chunk{
		{Content: "chunk 1"},
		{Content: "chunk 2"},
	}

	err := vi.Index(chunks, nil)
	if err != nil {
		t.Fatalf("Index failed: %v", err)
	}

	if vi.Dim != 3 {
		t.Errorf("Expected Dim 3, got %d", vi.Dim)
	}

	expectedLen := len(chunks) * 3
	if len(vi.Vectors) != expectedLen {
		t.Errorf("Expected Vector length %d, got %d", expectedLen, len(vi.Vectors))
	}
}
