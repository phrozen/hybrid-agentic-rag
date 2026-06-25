package semantic

import (
	"encoding/binary"
	"errors"
	"fmt"
	"unsafe"

	"github.com/phrozen/hybrid-agentic-rag/internal/models"
	"github.com/phrozen/hybrid-agentic-rag/internal/search"
)

// compile time assertion to ensure VectorIndex implements the Indexer interface
var _ search.Indexer = (*VectorIndex)(nil)

var (
	// ErrInvalidData is returned when binary deserialization fails due to corrupted inputs.
	ErrInvalidData = errors.New("invalid or corrupted binary data")
)

// Embedder abstracts the embedding generation pipeline.
// This allows VectorIndex to generate contextual embeddings internally,
// preserving a clean string-in/hits-out surface-level search API.
type Embedder interface {
	Embed(texts []string) ([][]float32, error)
}

// VectorIndex holds the flat contiguous slice of quantized vectors.
// It is designed to be fully driven by the unified Store using a raw string-in/hits-out API.
type VectorIndex struct {
	Vectors []int8   // Flat sequence of quantized elements of size N * Dim
	Dim     uint32   // Dimension size of vectors (e.g. 1024)
	client  Embedder // Dependency-injected generator abstraction
}

// NewVectorIndex instantiates a VectorIndex with its embedding dependencies.
func NewVectorIndex(client Embedder) *VectorIndex {
	return &VectorIndex{
		Vectors: nil,
		Dim:     0,
		client:  client,
	}
}

func (vi *VectorIndex) Index(chunks []*models.Chunk, progress search.ProgressFunc) error {
	if len(chunks) == 0 {
		return nil
	}

	const batchSize = 32
	var allFloatVectors [][]float32
	totalBatches := (len(chunks) + batchSize - 1) / batchSize

	for i := 0; i < len(chunks); i += batchSize {
		batchIdx := (i / batchSize) + 1
		end := i + batchSize
		if end > len(chunks) {
			end = len(chunks)
		}

		batch := chunks[i:end]
		texts := make([]string, len(batch))
		for j, ch := range batch {
			// =================================================================
			// USER FINE-TUNING AREA (EMBEDDING INPUT PREPARATION)
			// You can adjust exactly how the chunk's content is formatted
			// before sending it to Jina (e.g. prefixing, breadcrumbs, etc.)
			// =================================================================

			texts[j] = fmt.Sprintf("Document: %s\nContext: %s\n\n%s", ch.Source, ch.Breadcrumb, ch.Content)

			// =================================================================
		}

		// Report pre-batch status through the progress callback (if provided)
		if progress != nil {
			progress(fmt.Sprintf("Embedding batch %d/%d (chunks %d-%d)", batchIdx, totalBatches, i, end-1))
		}

		emb, err := vi.client.Embed(texts)
		if err != nil {
			return err
		}
		allFloatVectors = append(allFloatVectors, emb...)
	}

	// 1. Establish index dimension size
	vi.Dim = uint32(len(allFloatVectors[0]))

	// 2. Pre-allocate flat contiguous int8 array
	vi.Vectors = make([]int8, 0, len(chunks)*int(vi.Dim))

	// 3. Project down and append
	for _, fVec := range allFloatVectors {
		qVec := quantize(fVec)
		vi.Vectors = append(vi.Vectors, qVec...)
	}

	return nil
}

// Search performs a highly-optimized KNN search over the quantized vectors for a query string.
// It generates the search coordinates, quantizes them, and filters them internally.
func (vi *VectorIndex) Search(query string, k int) []search.Hit {
	var hits []search.Hit
	if vi.client == nil || len(vi.Vectors) == 0 || vi.Dim <= 0 {
		return hits
	}

	// 1. Generate standard embedding using Jina Asymmetric query prefixes
	outputVec, err := vi.client.Embed([]string{"Query: " + query})
	if err != nil {
		return hits
	}

	// 2. Keep Store oblivious to quantization by projecting onto int8 internally
	qVecQuantized := quantize(outputVec[0])

	// 3. Complete KNN search
	return knn(vi.Vectors, qVecQuantized, int(vi.Dim), k)
}

// MarshalBinary encodes the VectorIndex into a high-performance raw byte stream.
// It serializes as: [Dim (4 bytes)] [VectorsLength (4 bytes)] [Vectors (N * Dim bytes)].
// This implements encoding.BinaryMarshaler.
func (vi *VectorIndex) MarshalBinary() ([]byte, error) {
	if len(vi.Vectors) == 0 || vi.Dim <= 0 {
		return nil, nil
	}

	buf := make([]byte, 8+len(vi.Vectors))
	binary.LittleEndian.PutUint32(buf[0:4], vi.Dim)
	binary.LittleEndian.PutUint32(buf[4:8], uint32(len(vi.Vectors)))

	// Cast []int8 to []byte without extra heap allocation using modern unsafe.Slice
	src := unsafe.Slice((*byte)(unsafe.Pointer(&vi.Vectors[0])), len(vi.Vectors))
	copy(buf[8:], src)

	return buf, nil
}

// UnmarshalBinary decodes a raw byte stream into the VectorIndex with clean memory boundaries.
// This implements encoding.BinaryUnmarshaler.
func (vi *VectorIndex) UnmarshalBinary(data []byte) error {
	if len(data) < 8 {
		return ErrInvalidData
	}

	// Read header metadata
	vi.Dim = binary.LittleEndian.Uint32(data[0:4])
	vecDataLen := binary.LittleEndian.Uint32(data[4:8])

	// Strict bounds checking
	if vi.Dim == 0 || vecDataLen == 0 {
		vi.Vectors = nil
		return nil
	}

	if len(data) < 8+int(vecDataLen) {
		return ErrInvalidData
	}

	// Extract exact vector segment (accounting for any potential trailing bytes)
	exactBytes := data[8 : 8+vecDataLen]

	// Allocate fresh memory to guarantee ownership safety and prevent side-channel updates
	vi.Vectors = make([]int8, vecDataLen)
	src := unsafe.Slice((*int8)(unsafe.Pointer(&exactBytes[0])), len(exactBytes))
	copy(vi.Vectors, src)

	return nil
}
