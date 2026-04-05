// Package code provides the code_review tool for analyzing Go code.
package code

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

// ReviewArgs defines the arguments for the code_review tool.
type ReviewArgs struct {
	Code string `json:"code" jsonschema:"the Go source code to review"`
	Hint string `json:"hint,omitempty" jsonschema:"optional guidance for the reviewer, e.g. focus on security"`
}

// Review analyzes Go code using Gemini on Vertex AI and returns improvements in Markdown.
func Review(ctx context.Context, project, location string, args ReviewArgs) (string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  project,
		Location: location,
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		return "", fmt.Errorf("creating genai client: %w", err)
	}

	prompt := "You are an expert Go code reviewer. " +
		"Analyze the following Go code and provide a list of improvements " +
		"according to Go community best practices (Effective Go, Go Code Review Comments, Google Go Style Guide). " +
		"Format your response in Markdown.\n\n"

	if args.Hint != "" {
		prompt += fmt.Sprintf("Focus area: %s\n\n", args.Hint)
	}

	prompt += "```go\n" + args.Code + "\n```"

	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-pro", genai.Text(prompt), nil)
	if err != nil {
		return "", fmt.Errorf("generating content: %w", err)
	}

	return result.Text(), nil
}
