package semantic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Config contains configuration options for initializing the embedding client.
// It supports both local servers (like llama-server or Ollama) and commercial cloud APIs
// (such as OpenAI, Voyage AI, or Cohere) by implementing standard OpenAI-compatible fields.
type Config struct {
	BaseURL string        // Base URL of the embedding server (defaults to "http://localhost:8080")
	APIKey  string        // Optional: API key for standard cloud provider bearers (e.g., Bearer auth keys)
	Model   string        // Optional: Model identifier required by some cloud platforms (e.g., "text-embedding-3-small")
	Timeout time.Duration // Optional: Timeout duration (defaults to 30s)
}

// Client represents a simple, high-performance, OpenAI-compatible HTTP client.
// It communicates with standard `/v1/embeddings` API structures, making it a drop-in replacement
// for migrating between local llama.cpp backends and remote commercial model endpoints.
type Client struct {
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewClient instantiates a new Client. If cfg is nil, default local parameters pointing to
// "http://localhost:8080" are loaded.
func NewClient(cfg *Config) *Client {
	// Load fallback default config
	if cfg == nil {
		cfg = &Config{}
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	return &Client{
		baseURL: baseURL,
		apiKey:  cfg.APIKey,
		model:   cfg.Model,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// EmbedRequest matches the standard format accepted by OpenAI-compatible /v1/embeddings endpoints.
type EmbedRequest struct {
	Input []string `json:"input"`
	Model string   `json:"model,omitempty"` // Model name parameter required by most cloud APIs
}

// EmbedResponseData holds the individual vector data inside the response array.
type EmbedResponseData struct {
	Index     int       `json:"index"`
	Embedding []float32 `json:"embedding"`
}

// EmbedUsage represents the token usage details returned by llama-server.
type EmbedUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// EmbedResponse matches the OpenAI-compatible JSON response from llama-server.
type EmbedResponse struct {
	Model  string              `json:"model"`
	Object string              `json:"object"`
	Data   []EmbedResponseData `json:"data"`
	Usage  EmbedUsage          `json:"usage"`
}

// Embed retrieves dense vectors for the given slice of texts.
func (c *Client) Embed(texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	reqBody := EmbedRequest{
		Input: texts,
		Model: c.model,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal embed request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/embeddings", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// If APIKey is configured, inject standard Bearer authentication header.
	if c.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embedding server returned non-ok status: %d", resp.StatusCode)
	}

	var response EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode embed response: %w", err)
	}

	// Reconstruct the embedding slice in the exact requested order
	embeddings := make([][]float32, len(texts))
	for _, item := range response.Data {
		if item.Index >= 0 && item.Index < len(embeddings) {
			embeddings[item.Index] = item.Embedding
		}
	}

	// Handle case where any embedding came back missing
	for i, emb := range embeddings {
		if len(emb) == 0 {
			return nil, fmt.Errorf("missing embedding for input index %d", i)
		}
	}

	return embeddings, nil
}
