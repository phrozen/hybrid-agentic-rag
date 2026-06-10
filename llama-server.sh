llama-server \
  -m ~/models/jina-embeddings-v5-small-retrieval-Q4_K_M.gguf \
  --embedding \
  --pooling last \
  -c 2048 \
  -b 2048 \
  -ub 2048 \
  --no-cache-prompt \
  --parallel 1