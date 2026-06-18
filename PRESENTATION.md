# PRESENTATION.md — Hybrid Agentic RAG

## Talk Info

- **Title:** Hybrid Agentic RAG: Where Search Meets AI
- **Duration:** 30 minutes (20 min talk + 10 min demo & QA)
- **Format:** 7 content slides + live demo + Q&A
- **Tool:** Google Slides (Gemini-generated from this outline)

---

## Slide 1: Title / Who I Am

### Visual

- Title slide, bold text
- Subtitle: "Giving AI agents a search engine that actually works"
- Your name, title, company/handle

### Content

- Your name and role
- One-liner on what you build
- "Today: how we combined two search paradigms into one tool that AI agents actually want to use"

### Presenter Notes

> Don't spend more than 60 seconds here. State who you are, what the talk is about, and move on. The audience came for the tech, not the bio. Skip any slide animations — we have 20 minutes.

---

## Slide 2: What is RAG?

### Visual

- Simple diagram, center of slide:
  ```
  LLM + Context → Better Answer
        ↑
      Retrieved Data
  ```
- Below: two paths leading to that retrieved data:
  - Left path (Classic): harness retrieves → stuffs into prompt → LLM generates once
  - Right path (Agentic): LLM calls search tool → gets results → decides if more needed → generates

### Content

- LLMs are powerful but frozen at training time — they don't know your data
- RAG = Retrieval-Augmented Generation: augment the LLM's context with external data
- The core idea is simple: look it up, don't memorize it
- Two architectures for getting that data into context:
  - **Classic**: a harness retrieves data before the LLM turn, stuffs it into the prompt
  - **Agentic**: the LLM decides when and what to retrieve using tools
- "Same goal, different control flow. Let's see how each works."

### Presenter Notes

> Keep this to 2 minutes. The key distinction to plant: RAG is about augmenting context — that's it. The *how* is the architectural split we'll unpack next. Don't go deep into either architecture yet; just name them and move on. The diagram is the anchor — point to it twice: once for "what RAG is" (the center), once for "how it's done" (the two paths).

---

## Slide 3: Classic RAG vs Agentic RAG

### Visual

- Two side-by-side diagrams:

**Classic RAG (left):**
```
User Query
    ↓
[Retrieve] → [Stuff Context] → [Generate] → Answer
    ↑              ↑
  Pipeline     One shot
  fixed         fixed context
```

**Agentic RAG (right):**
```
User Query
    ↓
LLM → search("volcanic landmark") → results
  ↓                                  ↓
  LLM evaluates → need more?
  ↓           ↓ yes              ↓ no
  get_content(2274)            generate → Answer
  ↓
  get_document(everflame-peak)
  ↓
  generate → Answer
```

### Content

- **Classic RAG**: retrieve once, stuff into prompt, generate once
  - Fixed pipeline — the harness decides what to retrieve
  - Context window constraint — you can only stuff so much
  - No follow-up — if retrieval misses, the answer is wrong

- **Agentic RAG**: the LLM has tools, it decides when and how to search
  - `search(query, k, mode)` — find relevant chunks
  - `get_content(index)` — read a specific chunk
  - `get_document(source)` — pull the full document
  - The agent can iterate — search, evaluate, search again, go deeper
  - Heading-based breadcrumbs tell the agent *where* it landed

- "Classic is a pipeline. Agentic is a loop. The agent has agency."
- Classic RAG pipelines often add query rewriting (LLM rephrases the query before retrieval) and reranking (cross-encoder re-scores results after retrieval). In agentic RAG, the LLM naturally rewrites its own queries, and RFF fusion makes reranking less necessary.

### Presenter Notes

> This is the key architectural distinction. Take 3 minutes here. Walk through the left side fast — "retrieve, stuff, generate, done. Simple, but if retrieval misses, you're stuck." Then spend most time on the right side — "the LLM calls search, gets results, decides if it needs more, calls get_content or get_document, then generates. It's a loop, not a pipeline." Emphasize the breadcrumb: "Siege of Silverhaven > Participants > Defenders" tells the agent *where* it is in the document hierarchy. This is why heading-based chunking matters. This slide sets up the entire architecture — the rest is about *how* the search tool works.

---

## Slide 4: The Old Way — Keyword Search

### Visual

- Left side: diagram of an inverted index (term → [doc1, doc3, doc7])
- Right side: two query results:
  - ✅ Query: `"Throndor the Wise"` → Found! (exact term match)
  - ❌ Query: `"volcanic landmark with eternal flame"` → No results (vocabulary gap)
- Bottom: callout box — "BM25: decades of search, still the foundation"

### Content

- Inverted index: map every word to every document that contains it
- BM25 scoring: term frequency, inverse document frequency, document length normalization
- Strengths: fast, interpretable, exact matching on names and rare terms
- The vocabulary gap: you must use the same words the document uses
- "If you search for 'volcanic forge' but the document says 'Everflame Peak', BM25 finds nothing"

### Presenter Notes

> This slide sets up the problem. BM25 is not broken — it's excellent at what it does. Emphasize: keyword search is still the right tool for exact entity lookups. The problem is that humans describe things in ways documents don't. Use the "Everflame Peak" example — it's from the corpus and it'll recur in the demo, creating a callback moment. 2.5 minutes max.

---

## Slide 5: The New Way — Semantic Search

### Visual

- Diagram: text → embedding model → vector (show a few numbers: [0.12, -0.87, 0.34, ...]) → point in space
- Two clusters in "vector space" with a query arrow landing near the right cluster
- Callout: "Similar meaning → similar vector → nearby in vector space"

### Content

- Embedding models convert text into high-dimensional vectors (768, 1024, 1536 dims)
- Training: learns that "volcanic landmark" and "Everflame Peak" are related
- Cosine similarity: measures angle between vectors (1.0 = identical, 0.0 = unrelated)
- Int8 quantization: compress 1024-dim float32 vectors to int8 — 4x smaller, <1% quality loss
- The payoff: "volcanic landmark with eternal flame" ≈ "Everflame Peak" ✓

### Visual: Quantization Table

| | float32 | int8 |
|---|---------|------|
| "Elven Magic" ↔ "Elven Weavings" (related) | 0.5376 | 0.5403 |
| "Elven Magic" ↔ "Dwarven Mining" (unrelated) | 0.3667 | 0.3694 |

### Sidebar: Common Pipeline Enhancements

Two techniques you'll see in production RAG pipelines:

- **Query rewriting**: an LLM rephrases the user's query before retrieval. Useful when user input is vague or poorly worded. In agentic RAG the LLM already chooses its own search terms, so this is effectively built in.
- **Reranking**: a cross-encoder or LLM re-scores the top-k results after retrieval. Adds latency but improves precision. With RFF fusing two already-ranked lists and k=10, we found reranking unnecessary for this corpus.

### Presenter Notes

> This is the conceptual heavy lift. Keep it visual — the cluster diagram does more work than words. Walk through: (1) text goes into embedding model, (2) comes out as a vector, (3) similar concepts cluster together, (4) cosine similarity finds nearest neighbors. Then hit the quantization table fast — "we compress these vectors 4x with under 1% relative error." Don't explain MRL or quantization math — the table speaks for itself. 3.5 minutes. If running long, skip the quantization detail and just say "we compress 4x with minimal quality loss."
>
> The sidebar on query rewriting and reranking is a 30-second mention, not a deep dive. Say: "You'll see query rewriting and reranking in a lot of RAG pipelines. Query rewriting is effectively built into agentic RAG — the LLM picks its own search terms. Reranking adds latency for marginal gains when you've already got two ranked lists fusing with RFF. Worth knowing about, not worth adding here." Move on fast.

---

## Slide 6: Hybrid Agentic RAG — This Is What We Built

### Visual

- Architecture diagram showing the full pipeline:

```
                     ┌──────────────┐
                     │   Corpus     │
                     │  200 docs    │
                     │  3,536 chunks│
                     └──────┬───────┘
                            │ parse + index
                    ┌───────┴────────┐
                    │                │
              ┌─────┴──────┐   ┌─────┴──────┐
              │  BM25      │   │  Vectors   │
              │  index.bin │   │  embed.bin │
              └─────┬──────┘   └─────┬──────┘
                    │                │
                    └───────┬────────┘
                            │ RFF fusion
                     ┌──────┴───────┐
                     │  MCP Server  │
                     │  HTTP /mcp   │
                     └──────┬───────┘
                            │ tool calls
                     ┌──────┴───────┐
                     │  LLM Agent   │
                     │  (opencode)  │
                     └──────────────┘
```

- Three search mode callouts:
  - `text` → BM25 only
  - `semantic` → vectors only
  - `rff` (default) → both, fused

- Key numbers: RFF k=60, weights 0.5/0.5, int8 quantization, ~2.4MB vectors

### Content

- BM25 catches exact names and rare terms that vectors miss
- Vector search catches concepts and paraphrases that keywords miss
- RFF merges both with equal weight — simple, no training, no tuning
- Three search modes in one tool call: `search(query, k, mode)`
- Heading-based chunking → each chunk is a focused retrieval unit with a breadcrumb
- The agentic loop: `search` → `get_content` → `get_document` — the agent decides how deep to go
- "This is the system. Let me show you it working."

### Presenter Notes

> This is the payoff slide. Everything before was building to this. 3 minutes. Walk the architecture diagram bottom to top: "The agent calls search, which hits both indexes, fuses results with RFF, and returns ranked chunks with breadcrumbs. If the agent needs more, it calls get_content for a specific chunk or get_document for the full file." The diagram IS the slide — don't read every box, trace the two paths (BM25 and vectors) converging into RFF. Then transition hard: "Let me show you it working."

---

## Slide 7: Demo

### Visual

- Title: "Live Demo" with the three search modes
- The architecture diagram from slide 6 stays visible as a reference (small, bottom corner)
- No content bullets — this is live. Have a terminal window ready.

### Demo Script (7 minutes)

**Query 1 — Text mode (exact entity) — 1.5 min:**
```
search("Throndor the Wise", mode="text", k=3)
```
Point out: exact hit, top result, show the breadcrumb. "BM25 nails this — it found 'Throndor the Wise > Overview' instantly."

**Query 2 — Semantic mode (conceptual, no name) — 2 min:**
```
search("volcanic landmark with an eternal flame", mode="semantic", k=3)
```
Point out: "I never said 'Everflame Peak'. BM25 would find nothing here. But the embedding model knows that 'volcanic' + 'eternal flame' is semantically close to this document." Show the score (~0.73). "That's cosine similarity on int8-quantized vectors — 0.73 is a strong match."

**Query 3 — RFF mode (multi-hop, the hero moment) — 2 min:**
```
search("characters on both sides of the Siege of Silverhaven", mode="rff", k=5)
```
Point out: "This requires understanding 'both sides' and 'Siege of Silverhaven' and connecting characters to the event. RFF merges keyword hits for 'Siege of Silverhaven' with semantic matches for 'characters on both sides'." Show multiple results from the same document at different heading levels.

**Follow-up — the agentic pattern — 1.5 min:**
```
get_content(index=532)   # Show just the chunk
get_document(source="chronicles-of-aethelgard/events/siege-of-silverhaven.md")
```
"Search gives you the chunk. get_content gives you the detail. get_document gives you the full context. The agent decides how deep to go — that's agentic RAG."

### Presenter Notes

> The demo is the moment. Practice it. Have the MCP server running before the talk — don't make the audience watch you start it. If something breaks, have a screenshot backup of each query result. The three-mode comparison is the hero moment — it proves the thesis in 60 seconds. Don't narrate every field in the output — highlight the breadcrumb, the score, and the relevance. Keep demo to exactly 7 minutes.

---

## Q&A (3 minutes)

### Anticipated Questions

- **"Why k=60 for RFF?"** — It's the canonical value from the Cormack et al. paper. It controls how quickly scores decay across ranks. Lower k overweights rank 1; higher k flattens everything. 60 is the sweet spot for top-10 retrieval.

- **"Why equal weights 0.5/0.5?"** — Because both modes are equally strong on this corpus. If your keyword index is noisy or your embedding model is weak, you'd adjust. But for a well-structured corpus with a good embedding model, 50/50 works.

- **"Why int8 quantization instead of float16?"** — 4x compression vs 2x. The relative cosine error is under 1% for int8, and the vectors fit in ~2.4MB for 3,536 chunks. Float16 would be ~4.6MB. For a demo corpus this doesn't matter, but at 1M+ chunks it's the difference between fitting in memory or not.

- **"Why heading-based chunking instead of fixed-size?"** — Fixed-size chunks split mid-sentence and lose semantic coherence. Heading-based chunks align with document structure, so each chunk is a focused topic unit with a breadcrumb telling the agent where it landed.

- **"Can I use this corpus for my own RAG demos?"** — Yes, it's all in the repo. MIT licensed. Clone it, run `go run ./cmd/coe parse -i ./chronicles-of-aethelgard -o ./data -e`, then `go run ./cmd/coe serve -d ./data`.

- **"What embedding model did you use?"** — Jina embeddings v5 small (1024 dimensions), but the pipeline works with any OpenAI-compatible /v1/embeddings endpoint. Swap the model, reindex, done.

### Presenter Notes

> If Q&A runs long, cut the quantization question — it's the most likely to eat time. If no one asks questions, have the "why heading-based chunking" answer ready as a closing remark. It's the most impactful architectural decision and worth emphasizing even without a prompt.

---

## Slide Design Notes

- **Font:** Monospace for code/queries, sans-serif for body. Keep it readable from the back row.
- **Colors:** Dark background, light text. One accent color for highlights (e.g., green for ✅, red for ❌, blue for search terms).
- **Diagrams:** Use consistent iconography — magnifying glass for search, brain/vector for semantic, merge icon for RFF.
- **Code blocks:** Show actual tool calls and responses, not pseudocode. The audience wants to see the real interface.
- **Animations:** Minimal. The demo is the animation. Slides should be static, clear, and fast.
- **Timing:** 2-3 minutes per slide, 7 minutes for demo, 3 minutes for QA. If running long, cut quantization detail from slide 5 and the RFF formula from slide 6.

---

## Timing Breakdown

| Slide | Topic | Time |
|-------|-------|------|
| 1 | Title / Who I Am | 1 min |
| 2 | What is RAG? | 2 min |
| 3 | Classic vs Agentic RAG | 3 min |
| 4 | Keyword Search (BM25) | 2.5 min |
| 5 | Semantic Search (vectors) | 3.5 min |
| 6 | Hybrid Agentic RAG (architecture) | 3 min |
| 7 | Live Demo | 7 min |
| - | Q&A | 3 min |
| **Total** | | **25 min** |

Buffer: 5 minutes for transitions, audience reactions, and overages.