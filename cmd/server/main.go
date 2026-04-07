package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"godoctor/internal/tools/code"
	"godoctor/internal/tools/docs"
	"godoctor/internal/tools/explain"
	"godoctor/internal/tools/gendocs"
	"godoctor/internal/tools/gotest"
	"godoctor/internal/tools/govet"
	"godoctor/internal/tools/sbom"
)

func main() {
	project := flag.String("project", os.Getenv("GOOGLE_CLOUD_PROJECT"), "Google Cloud Project ID")
	location := flag.String("location", os.Getenv("GOOGLE_CLOUD_LOCATION"), "Google Cloud Location")
	useHTTP := flag.Bool("http", false, "Use Streamable HTTP transport")
	listen := flag.String("listen", ":8080", "Address to listen on (used with -http)")
	flag.Parse()

	newServer := func() *mcp.Server {
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

		mcp.AddTool(server, &mcp.Tool{
			Name:        "sbom_generate",
			Description: "Generates a Software Bill of Materials from go.mod content.",
		}, func(ctx context.Context, req *mcp.CallToolRequest, args sbom.Args) (*mcp.CallToolResult, any, error) {
			result, err := sbom.Generate(args)
			if err != nil {
				return nil, nil, err
			}
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: result},
				},
			}, nil, nil
		})

		mcp.AddTool(server, &mcp.Tool{
			Name:        "go_test",
			Description: "Runs go test on a specified package and returns the results.",
		}, func(ctx context.Context, req *mcp.CallToolRequest, args gotest.Args) (*mcp.CallToolResult, any, error) {
			result, err := gotest.Run(args)
			if err != nil {
				return &mcp.CallToolResult{
					Content: []mcp.Content{
						&mcp.TextContent{Text: result},
					},
					IsError: true,
				}, nil, nil
			}
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: result},
				},
			}, nil, nil
		})

		mcp.AddTool(server, &mcp.Tool{
			Name:        "go_vet",
			Description: "Runs go vet on a specified package and reports code issues.",
		}, func(ctx context.Context, req *mcp.CallToolRequest, args govet.Args) (*mcp.CallToolResult, any, error) {
			result, err := govet.Run(args)
			if err != nil {
				return &mcp.CallToolResult{
					Content: []mcp.Content{
						&mcp.TextContent{Text: result},
					},
					IsError: true,
				}, nil, nil
			}
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: result},
				},
			}, nil, nil
		})

		mcp.AddTool(server, &mcp.Tool{
			Name:        "explain_code",
			Description: "Explains Go code in English or Japanese using Gemini.",
		}, func(ctx context.Context, req *mcp.CallToolRequest, args explain.Args) (*mcp.CallToolResult, any, error) {
			result, err := explain.Run(ctx, *project, *location, args)
			if err != nil {
				return nil, nil, err
			}
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: result},
				},
			}, nil, nil
		})

		mcp.AddTool(server, &mcp.Tool{
			Name:        "generate_docs",
			Description: "Generates GoDoc comments for exported identifiers in Go code.",
		}, func(ctx context.Context, req *mcp.CallToolRequest, args gendocs.Args) (*mcp.CallToolResult, any, error) {
			result, err := gendocs.Run(ctx, *project, *location, args)
			if err != nil {
				return nil, nil, err
			}
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: result},
				},
			}, nil, nil
		})

		return server
	}

	if *useHTTP {
		handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
			return newServer()
		}, nil)
		http.Handle("/mcp", handler)
		log.Printf("Listening on %s", *listen)
		log.Fatal(http.ListenAndServe(*listen, nil))
	} else {
		if err := newServer().Run(context.Background(), &mcp.StdioTransport{}); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}
}
