package store

import (
	"testing"

	"github.com/phrozen/hybrid-agentic-rag/internal/models"
	"github.com/phrozen/hybrid-agentic-rag/internal/search"
)

type mockIndexer struct {
	hits []search.Hit
}

func (m *mockIndexer) Index(chunks []*models.Chunk) error      { return nil }
func (m *mockIndexer) MarshalBinary() ([]byte, error)          { return nil, nil }
func (m *mockIndexer) UnmarshalBinary(data []byte) error       { return nil }
func (m *mockIndexer) Search(query string, k int) []search.Hit { return m.hits }

func TestSearch_ModeDispatch(t *testing.T) {
	chunks := []*models.Chunk{
		{Source: "factions/mages-guild.md", Breadcrumb: "Mages' Guild > Overview"},
		{Source: "characters/throndor.md", Breadcrumb: "Throndor the Wise > Abilities"},
		{Source: "events/war.md", Breadcrumb: "War of Shadows > Aftermath"},
	}

	textIdx := &mockIndexer{
		hits: []search.Hit{{Index: 0, Score: 10.0}},
	}
	semanticIdx := &mockIndexer{
		hits: []search.Hit{{Index: 1, Score: 0.95}},
	}
	docs := []*models.Document{
		{Source: "factions/mages-guild.md"},
		{Source: "characters/throndor.md"},
		{Source: "events/war.md"},
	}

	s := New(docs, chunks, textIdx, semanticIdx)

	t.Run("text mode uses text indexer only", func(t *testing.T) {
		results := s.Search("mage", 5, search.SearchText)
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if results[0].Index != 0 {
			t.Errorf("expected index 0, got %d", results[0].Index)
		}
		if results[0].Source != "factions/mages-guild.md" {
			t.Errorf("expected source factions/mages-guild.md, got %s", results[0].Source)
		}
	})

	t.Run("semantic mode uses semantic indexer only", func(t *testing.T) {
		results := s.Search("ancient powers", 5, search.SearchSemantic)
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if results[0].Index != 1 {
			t.Errorf("expected index 1, got %d", results[0].Index)
		}
		if results[0].Breadcrumb != "Throndor the Wise > Abilities" {
			t.Errorf("unexpected breadcrumb: %s", results[0].Breadcrumb)
		}
	})

	t.Run("rff mode fuses both indexers", func(t *testing.T) {
		results := s.Search("mage", 5, search.SearchRFF)
		if len(results) != 2 {
			t.Fatalf("expected 2 results from RFF fusion, got %d", len(results))
		}
	})
}

func TestSearch_RFFTruncation(t *testing.T) {
	chunks := make([]*models.Chunk, 5)
	for i := range chunks {
		chunks[i] = &models.Chunk{Source: "doc.md", Breadcrumb: "Doc > Section"}
	}
	docs := []*models.Document{{Source: "doc.md"}}

	textIdx := &mockIndexer{
		hits: []search.Hit{
			{Index: 0, Score: 10.0},
			{Index: 1, Score: 9.0},
			{Index: 2, Score: 8.0},
		},
	}
	semanticIdx := &mockIndexer{
		hits: []search.Hit{
			{Index: 3, Score: 0.95},
			{Index: 4, Score: 0.90},
		},
	}

	s := New(docs, chunks, textIdx, semanticIdx)
	results := s.Search("test", 3, search.SearchRFF)

	if len(results) != 3 {
		t.Fatalf("expected 3 truncated results, got %d", len(results))
	}
}

func TestSearch_SkipsOutOfBoundsHits(t *testing.T) {
	chunks := []*models.Chunk{
		{Source: "a.md", Breadcrumb: "A > Overview"},
	}
	docs := []*models.Document{
		{Source: "a.md"},
	}

	textIdx := &mockIndexer{
		hits: []search.Hit{
			{Index: 0, Score: 5.0},
			{Index: 99, Score: 3.0},
			{Index: -1, Score: 1.0},
		},
	}
	semanticIdx := &mockIndexer{}

	s := New(docs, chunks, textIdx, semanticIdx)
	results := s.Search("test", 5, search.SearchText)

	if len(results) != 1 {
		t.Fatalf("expected 1 valid result, got %d", len(results))
	}
	if results[0].Index != 0 {
		t.Errorf("expected index 0, got %d", results[0].Index)
	}
}

func TestGetContent(t *testing.T) {
	chunks := []*models.Chunk{
		{Source: "a.md", Heading: "Overview", Content: "hello world"},
	}
	docs := []*models.Document{
		{Source: "a.md"},
	}

	s := New(docs, chunks, &mockIndexer{}, &mockIndexer{})

	t.Run("valid index", func(t *testing.T) {
		chunk, ok := s.GetContent(0)
		if !ok {
			t.Fatal("expected chunk to be found")
		}
		if chunk.Content != "hello world" {
			t.Errorf("unexpected content: %s", chunk.Content)
		}
	})

	t.Run("negative index", func(t *testing.T) {
		_, ok := s.GetContent(-1)
		if ok {
			t.Error("expected false for negative index")
		}
	})

	t.Run("out of bounds index", func(t *testing.T) {
		_, ok := s.GetContent(1)
		if ok {
			t.Error("expected false for out of bounds index")
		}
	})
}
