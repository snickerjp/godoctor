// Package explain provides the explain_code tool for explaining Go code.
package explain

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

// Args defines the arguments for the explain_code tool.
type Args struct {
	Code string `json:"code" jsonschema:"the Go source code to explain"`
	Lang string `json:"lang,omitempty" jsonschema:"output language: en or ja (default: en)"`
}

// Run explains Go code using Gemini on Vertex AI.
func Run(ctx context.Context, project, location string, args Args) (string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  project,
		Location: location,
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		return "", fmt.Errorf("creating genai client: %w", err)
	}

	lang := "English"
	if args.Lang == "ja" {
		lang = "Japanese"
	}

	prompt := fmt.Sprintf(
		"You are an expert Go developer. Explain the following Go code in %s. "+
			"Break down what each part does clearly for someone new to the codebase. "+
			"Format your response in Markdown.\n\n```go\n%s\n```",
		lang, args.Code)

	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-pro", genai.Text(prompt), nil)
	if err != nil {
		return "", fmt.Errorf("generating content: %w", err)
	}
	return result.Text(), nil
}
