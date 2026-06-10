package models

// Document represents a parsed markdown file with its content and chunks
type Document struct {
	// Source is the file path of the markdown document
	Source string `json:"source"`

	// Lines is the total number of lines in the document
	Lines int `json:"lines"`

	// Size is the size of the document content in bytes
	Size int `json:"size"`

	// Content is the raw content of the markdown file
	Content []byte `json:"content"`
}

// Chunk represents a single chunk of markdown content with its hierarchical context
type Chunk struct {
	// Source is the file path of the original document (for traceability)
	Source string `json:"source"`
	// Breadcrumb is the hierarchical path from H1 -> H2 -> H3
	// e.g., "The Silvermoon Academy of Stars > Overview"
	Breadcrumb string `json:"breadcrumb"`

	// Heading is the H3 heading text (or H2/H1 if no H3 exists)
	Heading string `json:"heading"`

	// HeadingLevel is the level of the heading that starts this chunk (1, 2, or 3)
	HeadingLevel int `json:"heading_level"`

	// Content is the markdown content for this chunk (excluding the heading itself)
	Content string `json:"content"`

	// StartLine is the line number where this chunk begins (1-indexed)
	StartLine int `json:"start_line"`

	// EndLine is the line number where this chunk ends (1-indexed)
	EndLine int `json:"end_line"`
}
