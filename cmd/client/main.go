package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	listTools := flag.Bool("tools-list", false, "List available tools on the server")
	callTool := flag.String("tool-call", "", "Call a specific tool on the server")
	serverPath := flag.String("server-path", "./bin/server", "Path to the server binary")
	flag.Parse()

	if !*listTools && *callTool == "" {
		flag.Usage()
		os.Exit(1)
	}

	ctx := context.Background()

	// Connect to the server using CommandTransport
	cmd := exec.Command(*serverPath)
	// We want to see stderr from the server for debugging
	cmd.Stderr = os.Stderr

	transport := &mcp.CommandTransport{Command: cmd}
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}, nil)

	cs, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer cs.Close()

	if *listTools {
		fmt.Println("Tools:")
		tools := cs.Tools(ctx, nil)
		for t, err := range tools {
			if err != nil {
				log.Fatalf("Error iterating tools: %v", err)
			}
			fmt.Printf("- %s: %s\n", t.Name, t.Description)
		}
	}

	if *callTool != "" {
		var toolArgs map[string]any
		if *callTool == "read_docs" {
			args := flag.Args()
			if len(args) == 0 {
				log.Fatal("read_docs requires at least a package name")
			}
			
			fullPath := args[0]
			// Try to find if there's a symbol (last part after a dot, if not in a URL-like path)
			// Actually, go doc handles it by itself if we just give the full path.
			// But the tool expects separate package and symbol.
			// Let's assume for now that if it has a slash, the last part after a dot is a symbol only if the dot is after the last slash.
			
			pkg := fullPath
			symbol := ""
			
			lastSlash := -1
			for i := len(fullPath) - 1; i >= 0; i-- {
				if fullPath[i] == '/' {
					lastSlash = i
					break
				}
			}
			
			lastDot := -1
			for i := len(fullPath) - 1; i > lastSlash; i-- {
				if fullPath[i] == '.' {
					lastDot = i
					break
				}
			}
			
			if lastDot != -1 {
				pkg = fullPath[:lastDot]
				symbol = fullPath[lastDot+1:]
			}
			
			toolArgs = map[string]any{
				"package": pkg,
			}
			if symbol != "" {
				toolArgs["symbol"] = symbol
			}
		}

		result, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      *callTool,
			Arguments: toolArgs,
		})
		if err != nil {
			log.Fatalf("Failed to call tool %s: %v", *callTool, err)
		}
		for _, content := range result.Content {
			if text, ok := content.(*mcp.TextContent); ok {
				fmt.Println(text.Text)
			}
		}
	}
}
