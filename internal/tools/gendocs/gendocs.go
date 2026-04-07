// Package gendocs provides the generate_docs tool for generating GoDoc comments.
package gendocs

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

// Args defines the arguments for the generate_docs tool.
type Args struct {
	Code string `json:"code" jsonschema:"the Go source code to generate documentation for"`
}

// Run generates GoDoc comments for Go code using Gemini on Vertex AI.
func Run(ctx context.Context, project, location string, args Args) (string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  project,
		Location: location,
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		return "", fmt.Errorf("creating genai client: %w", err)
	}

	prompt := "You are an expert Go developer. Add GoDoc comments to all exported identifiers " +
		"(functions, types, methods, constants, variables) in the following Go code. " +
		"Follow the Go documentation conventions: comments should start with the identifier name, " +
		"be complete sentences, and explain the why not just the what. " +
		"Return the complete code with the added comments.\n\n```go\n" + args.Code + "\n```"

	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-pro", genai.Text(prompt), nil)
	if err != nil {
		return "", fmt.Errorf("generating content: %w", err)
	}
	return result.Text(), nil
}
