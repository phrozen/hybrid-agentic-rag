package tools

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/phrozen/hybrid-agentic-rag/internal/logger"
)

// LoggingMiddleware returns MCP receiving middleware that emits one structured
// record per tool call: the tool name, its input arguments (grouped), a
// stats-only summary of the result (grouped), and the call latency.
//
// It logs through the package-level logger.Default, which is an independent
// sink from log/slog's discarded default logger (see internal/logger).
func LoggingMiddleware() mcp.Middleware {
	return func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
			params, ok := req.GetParams().(*mcp.CallToolParamsRaw)
			if !ok {
				// Not a tools/call (e.g. initialize, tools/list) — pass through.
				return next(ctx, method, req)
			}

			start := time.Now()
			result, err := next(ctx, method, req)
			elapsed := time.Since(start)

			input := group("input", decodeArgs(params.Arguments))

			if err != nil {
				logger.Error(params.Name,
					input,
					"elapsed", elapsed,
					"err", err,
				)
				return result, err
			}

			logger.Info(params.Name,
				input,
				group("output", summarizeResult(result)),
				"elapsed", elapsed,
			)
			return result, nil
		}
	}
}

// group builds an slog group attr from a map. Order follows Go's map iteration
// (unsorted); pterm renders the group value as-is. String values containing
// spaces are wrapped in single quotes so multi-word values (e.g. a search
// query) stay visually distinct within the bracketed group.
func group(name string, m map[string]any) slog.Attr {
	kv := make([]any, 0, len(m)*2)
	for k, v := range m {
		if s, ok := v.(string); ok && strings.ContainsRune(s, ' ') {
			v = "'" + s + "'"
		}
		kv = append(kv, k, v)
	}
	return slog.Group(name, kv...)
}

// decodeArgs unmarshals raw tool arguments into a map for structured logging.
// On failure it preserves the raw payload under a "raw" key rather than dropping it.
func decodeArgs(raw json.RawMessage) map[string]any {
	if len(raw) == 0 {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return map[string]any{"raw": string(raw)}
	}
	return m
}

// summarizeResult reduces a CallToolResult to stats only (no content preview):
// call status, the number of content blocks, and total text size in bytes.
func summarizeResult(result mcp.Result) map[string]any {
	ctr, ok := result.(*mcp.CallToolResult)
	if !ok || ctr == nil {
		return map[string]any{"status": "no-result"}
	}

	status := "ok"
	if ctr.IsError {
		status = "error"
	}

	var textBytes int
	for _, c := range ctr.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			textBytes += len(tc.Text)
		}
	}

	return map[string]any{
		"status": status,
		"blocks": len(ctr.Content),
		"bytes":  textBytes,
	}
}
