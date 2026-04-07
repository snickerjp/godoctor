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
	hint := flag.String("hint", "", "Optional hint for code_review (e.g. 'focus on security')")
	addr := flag.String("addr", "", "Server HTTP endpoint (e.g. http://localhost:8080/mcp)")
	flag.Parse()

	if !*listTools && *callTool == "" {
		flag.Usage()
		os.Exit(1)
	}

	ctx := context.Background()

	var transport mcp.Transport
	if *addr != "" {
		transport = &mcp.StreamableClientTransport{Endpoint: *addr, DisableStandaloneSSE: true}
	} else {
		cmd := exec.Command(*serverPath)
		cmd.Stderr = os.Stderr
		transport = &mcp.CommandTransport{Command: cmd}
	}

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
		toolArgs, err := buildToolArgs(*callTool, *hint, flag.Args())
		if err != nil {
			log.Fatal(err)
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

func buildToolArgs(tool, hint string, args []string) (map[string]any, error) {
	switch tool {
	case "read_docs":
		if len(args) == 0 {
			return nil, fmt.Errorf("read_docs requires at least a package name")
		}
		fullPath := args[0]
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
		toolArgs := map[string]any{"package": pkg}
		if symbol != "" {
			toolArgs["symbol"] = symbol
		}
		return toolArgs, nil

	case "code_review", "explain_code", "generate_docs":
		if len(args) == 0 {
			return nil, fmt.Errorf("%s requires a file path", tool)
		}
		content, err := os.ReadFile(args[0])
		if err != nil {
			return nil, fmt.Errorf("reading file: %w", err)
		}
		toolArgs := map[string]any{"code": string(content)}
		if tool == "code_review" && hint != "" {
			toolArgs["hint"] = hint
		}
		if tool == "explain_code" && hint != "" {
			toolArgs["lang"] = hint
		}
		return toolArgs, nil

	case "sbom_generate":
		if len(args) == 0 {
			return nil, fmt.Errorf("sbom_generate requires a go.mod file path")
		}
		content, err := os.ReadFile(args[0])
		if err != nil {
			return nil, fmt.Errorf("reading file: %w", err)
		}
		return map[string]any{"gomod": string(content)}, nil

	case "go_test":
		if len(args) == 0 {
			return nil, fmt.Errorf("go_test requires a package path")
		}
		return map[string]any{"package": args[0]}, nil

	case "go_vet":
		if len(args) == 0 {
			return nil, fmt.Errorf("go_vet requires a package path")
		}
		return map[string]any{"package": args[0]}, nil

	default:
		return nil, nil
	}
}
