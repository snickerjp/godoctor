package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"godoctor/internal/tools/docs"
)

func main() {
	// Create a server with a name
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "hello-mcp-server",
		Version: "1.0.0",
	}, nil)

	// Add the hello_world tool
	// The third argument is the handler function. 
	// Since we don't need any arguments for hello_world, we can use an empty struct or just ignore it.
	mcp.AddTool(server, &mcp.Tool{
		Name:        "hello_world",
		Description: "Returns a hello world message from the MCP server.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, any, error) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Hello, MCP world!",
				},
			},
		}, nil, nil
	})

	// Add the read_docs tool
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
				&mcp.TextContent{
					Text: doc,
				},
			},
		}, nil, nil
	})

	// Run the server on stdio transport
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
