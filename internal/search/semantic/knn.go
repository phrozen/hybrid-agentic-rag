package semantic

import (
	"container/heap"

	"github.com/phrozen/hybrid-agentic-rag/internal/search"
)

type item struct {
	index   int // Positional index of the chunk in the global array
	score   float32
	heapIdx int // Internal index tracking for container/heap operations
}

type minHeap []*item

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].score < h[j].score }
func (h minHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].heapIdx = i
	h[j].heapIdx = j
}

// Push adds an item to the heap.
func (h *minHeap) Push(x any) {
	n := len(*h)
	it := x.(*item)
	it.heapIdx = n
	*h = append(*h, it)
}

// Pop removes the minimum item (the root) from the heap.
func (h *minHeap) Pop() any {
	old := *h
	n := len(old)
	it := old[n-1]
	old[n-1] = nil
	it.heapIdx = -1
	*h = old[0 : n-1]
	return it
}

// knn performs a 1-pass top-K similarity search using a min-heap over a contiguous, flat slice of vectors.
// It accepts 'vectors' as a single flattened []int8 array, where the vector for index 'i' begins
// at 'i * dim' and has length 'dim'. This layout minimizes pointer chasing and cache misses.
func knn(vectors []int8, query []int8, dim int, k int) []search.Hit {
	if k <= 0 || len(vectors) == 0 || len(query) == 0 || dim <= 0 {
		return nil
	}

	numVectors := len(vectors) / dim
	h := &minHeap{}
	heap.Init(h)

	for idx := 0; idx < numVectors; idx++ {
		offset := idx * dim
		// Slice bounds are check-free inside similarity flat operation
		score := similarityFlat(vectors, offset, query, dim)

		if h.Len() < k {
			heap.Push(h, &item{
				index: idx,
				score: score,
			})
		} else if score > (*h)[0].score {
			heap.Pop(h)
			heap.Push(h, &item{
				index: idx,
				score: score,
			})
		}
	}

	// Drain the heap and store results in descending order by similarity
	res := make([]search.Hit, h.Len())
	for h.Len() > 0 {
		it := heap.Pop(h).(*item)
		res[h.Len()] = search.Hit{
			Index: it.index,
			Score: it.score,
		}
	}

	return res
}
