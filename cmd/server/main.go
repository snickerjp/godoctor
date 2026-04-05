package main

import (
	"context"
	"flag"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"godoctor/internal/tools/code"
	"godoctor/internal/tools/docs"
)

func main() {
	project := flag.String("project", "", "Google Cloud Project ID")
	location := flag.String("location", "", "Google Cloud Location")
	flag.Parse()

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "hello-mcp-server",
		Version: "1.0.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "hello_world",
		Description: "Returns a hello world message from the MCP server.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, any, error) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Hello, MCP world!"},
			},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "read_docs",
		Description: "Returns documentation for a Go package or symbol.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args docs.ReadDocsArgs) (*mcp.CallToolResult, any, error) {
		doc, err := docs.ReadDocs(args)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: doc},
			},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "code_review",
		Description: "Analyzes Go code using Gemini and returns improvements in Markdown.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args code.ReviewArgs) (*mcp.CallToolResult, any, error) {
		result, err := code.Review(ctx, *project, *location, args)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: result},
			},
		}, nil, nil
	})

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
