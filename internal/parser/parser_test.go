package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseAndChunk(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "test-chronicles-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write a mock markdown file
	mockContent := `# Test Race

## Overview
This is a test race overview.

## Biological Characteristics
Some biological content here.
`
	testFilePath := filepath.Join(tmpDir, "test-race.md")
	if err := os.WriteFile(testFilePath, []byte(mockContent), 0644); err != nil {
		t.Fatalf("Failed to write mock md file: %v", err)
	}

	// Run Parse
	docs, err := Parse(tmpDir)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(docs) != 1 {
		t.Fatalf("Expected exactly 1 document, got %d", len(docs))
	}

	doc := docs[0]
	if doc.Source != testFilePath {
		t.Errorf("Expected source path %q, got %q", testFilePath, doc.Source)
	}

	// Run Chunk
	chunks := Chunk(docs)

	// We expect 2 chunks, one for Overview (H2) and one for Biological Characteristics (H2)
	// (Note: H1 heading doesn't contain non-empty body before the next heading in this mockup so it gets skipped or flushed empty)
	if len(chunks) != 2 {
		t.Fatalf("Expected exactly 2 chunks, got %d", len(chunks))
	}

	c1 := chunks[0]
	if c1.Heading != "Overview" {
		t.Errorf("Expected first heading to be Overview, got %q", c1.Heading)
	}
	if c1.Content != "This is a test race overview." {
		t.Errorf("Expected chunk content mismatch, got %q", c1.Content)
	}
	if c1.Breadcrumb != "Test Race > Overview" {
		t.Errorf("Expected breadcrumb 'Test Race > Overview', got %q", c1.Breadcrumb)
	}

	c2 := chunks[1]
	if c2.Heading != "Biological Characteristics" {
		t.Errorf("Expected second heading to be Biological Characteristics, got %q", c2.Heading)
	}
	if c2.Content != "Some biological content here." {
		t.Errorf("Expected chunk content mismatch, got %q", c2.Content)
	}
}
