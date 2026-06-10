package tools

import (
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func prettyJSONText(v any) *mcp.CallToolResult {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(bytes)}},
	}
}

func markdownText(source string, metadata map[string]string, body string) *mcp.CallToolResult {
	var header string
	header += fmt.Sprintf("# %s\n\n", source)
	for k, v := range metadata {
		header += fmt.Sprintf("**%s:** %s\n", k, v)
	}
	text := header + "\n" + body
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}
}
