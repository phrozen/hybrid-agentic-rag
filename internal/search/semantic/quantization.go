package semantic

import (
	"math"
)

// quantize projects an L2-normalized float32 vector onto an int8 range [-127, 127]
// assuming a fixed input range of [-1.0, 1.0]. This eliminates the need
// to scan the vector for a dynamic absolute maximum, speeding up quantization,
// and removes the per-vector scale factor.
func quantize(v []float32) []int8 {
	if len(v) == 0 {
		return nil
	}

	vector := make([]int8, len(v))
	for i, val := range v {
		// Map [-1.0, 1.0] to [-127, 127] using a fixed scaling factor of 127.0
		qVal := math.Round(float64(val * 127.0))
		if qVal > 127 {
			qVal = 127
		} else if qVal < -127 {
			qVal = -127
		}
		vector[i] = int8(qVal)
	}

	return vector
}

// similarity computes an ultra-fast raw dot product between two quantized int8 vectors
// without safety bounds checks, maximizing throughput for internal vector database calculations and KNN.
// The result is mapped back to Cosine Similarity [-1.0, 1.0] by dividing by 16129.0 (127 * 127).
func similarity(a, b []int8) float32 {
	var dot int32
	for i := 0; i < len(a); i++ {
		dot += int32(a[i]) * int32(b[i])
	}

	// Direct scaling mapping from [-16129, 16129] back to Cosine Similarity [-1.0, 1.0]
	// 127 * 127 = 16129.0
	return float32(dot) / 16129.0
}

// similarityFlat computes an ultra-fast raw dot product against a contiguous, flat slice representation of vectors.
// It bypasses intermediate slice allocations and optimizes CPU cache architectures (Spatial Locality) for KNN.
func similarityFlat(vecs []int8, offset int, query []int8, dim int) float32 {
	subSlice := vecs[offset : offset+dim]
	_ = query[dim-1] // Single bounds check for query up front to eliminate checks inside the loop
	var dot int32
	for i := 0; i < dim; i++ {
		dot += int32(subSlice[i]) * int32(query[i])
	}

	// Direct scaling mapping from [-16129, 16129] back to Cosine Similarity [-1.0, 1.0]
	// 127 * 127 = 16129.0
	return float32(dot) / 16129.0
}

// l2NormalizeFloat32 normalizes a float32 vector to unit length (L2 norm = 1.0).
func l2NormalizeFloat32(v []float32) []float32 {
	var sum float64
	for _, val := range v {
		sum += float64(val * val)
	}
	norm := math.Sqrt(sum)
	if norm == 0 {
		return v
	}
	res := make([]float32, len(v))
	for i, val := range v {
		res[i] = float32(float64(val) / norm)
	}
	return res
}

// l2NormalizeInt8 scales an int8 vector so that its L2 norm is equal to 127.0.
// This is useful for constructing simulated mock vectors matching expected unit-length magnitudes.
func l2NormalizeInt8(v []int8) []int8 {
	var sum float64
	for _, val := range v {
		sum += float64(int32(val) * int32(val))
	}
	norm := math.Sqrt(sum)
	if norm == 0 {
		return v
	}
	res := make([]int8, len(v))
	for i, val := range v {
		qVal := math.Round(float64(val) * 127.0 / norm)
		if qVal > 127 {
			qVal = 127
		} else if qVal < -127 {
			qVal = -127
		}
		res[i] = int8(qVal)
	}
	return res
}
