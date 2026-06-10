# AGENTS.md — hybrid-agentic-rag

## Quick commands

```bash
go run ./cmd/coe parse -i ./chronicles-of-aethelgard -o ./data         # keyword index only
go run ./cmd/coe parse -i ./chronicles-of-aethelgard -o ./data -e      # keyword + semantic
go run ./cmd/coe serve -d ./data                                        # start MCP server
go build -o coe ./cmd/coe/
```

## Tests

```bash
# Unit tests
go test ./internal/parser/ ./internal/search/text/ ./internal/search/semantic/ ./internal/store/ -run 'TestParse|TestInvertedIndex|TestKNN|TestVectorIndex_Serial|TestVectorIndex_Unmarshal|TestRFF|TestSearch|TestGetContent'

# Integration tests — require llama-server on :8080 with Jina v5 model
go test ./internal/search/semantic/ -run 'TestClient_Embed|TestQuantization'
```

## Architecture

Single Go module, single binary `coe`. Three-layer pipeline:

1. **Parse** — walks markdown corpus, chunks on H1/H2/H3 heading boundaries
2. **Index** — BM25 keyword index (blaze → `index.bin`) + int8-quantized vector index (OpenAI-compatible `/v1/embeddings` → `embeddings.bin`)
3. **Serve** — loads artifacts into `Store`, registers MCP tools, serves over HTTP

```
internal/
├── parser/         file discovery + heading-based chunking
├── search/
│   ├── search.go   Indexer interface + RFF fusion + SearchResult/SearchMode types
│   ├── text/       BM25 wrapper for github.com/wizenheimer/blaze
│   └── semantic/   embedding client + int8 quantization + flat-vector KNN
├── store/
│   ├── store.go    unified Store: Search, GetContent, GetDocument, ListDocuments
│   └── store_test.go
├── tools/          MCP tool handlers: search, get_content, get_document, list_documents
└── models/         Document and Chunk structs
```

## Data artifacts

Pre-built in `data/`: `documents.json`, `chunks.json`, `index.bin`, `embeddings.bin`.
The `parse` command regenerates all four; `serve` reads them.

## Dependencies

| Package | Purpose |
|---------|---------|
| `spf13/cobra` | CLI subcommands |
| `wizenheimer/blaze` | BM25/proximity keyword index |
| `modelcontextprotocol/go-sdk` | MCP tool protocol |

## Gotchas

- `cmd/coe/commands/root.go` silences all `log/slog` globally
- `parse -e` needs an OpenAI-compatible `/v1/embeddings` endpoint on `:8080`. Pre-built artifacts were built with `jina-embeddings-v5-small-retrieval` (1024 dimensions, int8 quantized)
- Semantic index uses Jina asymmetric retrieval prefixes (`Document: ...` / `Query: ...`)
- Embeddings binary format: 4 bytes dim (LE) + 4 bytes data len (LE) + raw byte array

## Corpus conventions

The fantasy corpus in `chronicles-of-aethelgard/` (~200 markdown files) follows strict rules:
- **kebab-case filenames** matching H1 (e.g. `throndor-the-wise.md` → `# Throndor the Wise`)
- **Heading hierarchy H1/H2/H3** is critical — chunking splits on these
- Natural cross-references in prose, 3–5 per entry
- Full guidelines in `AETHELGARD.md`
