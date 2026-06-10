package store

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/phrozen/hybrid-agentic-rag/internal/models"
	"github.com/phrozen/hybrid-agentic-rag/internal/search"
)

type Store struct {
	documents map[string]*models.Document
	chunks    []*models.Chunk

	text   search.Indexer
	vector search.Indexer

	documentsTree string
}

func New(docs []*models.Document, chunks []*models.Chunk, text search.Indexer, vector search.Indexer) *Store {
	docmap := make(map[string]*models.Document)

	var tree strings.Builder
	fmt.Fprintf(&tree, "# %d documents\n\n", len(docs))

	var prevDir string
	for _, doc := range docs {
		docmap[doc.Source] = doc

		dir := filepath.Dir(doc.Source)
		if dir == "." {
			dir = ""
		}
		if dir != prevDir {
			if dir != "" {
				fmt.Fprintf(&tree, "%s/\n", dir)
			}
			prevDir = dir
		}
		fmt.Fprintf(&tree, "  %-40s (%dL, %dB)\n", filepath.Base(doc.Source), doc.Lines, doc.Size)
	}

	return &Store{
		documents:     docmap,
		chunks:        chunks,
		text:          text,
		vector:        vector,
		documentsTree: tree.String(),
	}
}

func (s *Store) GetDocument(source string) (*models.Document, bool) {
	doc, exists := s.documents[source]
	return doc, exists
}

func (s *Store) ListDocuments() string {
	return s.documentsTree
}

func (s *Store) GetContent(index int) (*models.Chunk, bool) {
	if index < 0 || index >= len(s.chunks) {
		return nil, false
	}
	return s.chunks[index], true
}

func (s *Store) Search(query string, k int, mode search.SearchMode) []search.SearchResult {
	var rawHits []search.Hit
	switch mode {
	case search.SearchText:
		rawHits = s.text.Search(query, k)
	case search.SearchSemantic:
		rawHits = s.vector.Search(query, k)
	default:
		textHits := s.text.Search(query, k)
		semanticHits := s.vector.Search(query, k)
		rawHits = search.RFF(textHits, semanticHits)
	}

	results := make([]search.SearchResult, 0, len(rawHits))
	for _, hit := range rawHits {
		if hit.Index < 0 || hit.Index >= len(s.chunks) {
			continue
		}
		if len(results) >= k {
			break
		}
		ch := s.chunks[hit.Index]
		results = append(results, search.SearchResult{
			Index:      hit.Index,
			Score:      hit.Score,
			Source:     ch.Source,
			Breadcrumb: ch.Breadcrumb,
		})
	}

	return results
}
