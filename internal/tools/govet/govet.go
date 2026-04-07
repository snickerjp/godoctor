// Package govet runs go vet on a specified package.
package govet

import (
	"bytes"
	"fmt"
	"os/exec"
)

// Args defines the arguments for the go_vet tool.
type Args struct {
	Package string `json:"package" jsonschema:"the Go package to vet, e.g. ./internal/tools/..."`
}

// Run executes go vet on the specified package and returns the output.
func Run(args Args) (string, error) {
	cmd := exec.Command("go", "vet", args.Package)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	result := out.String() + stderr.String()
	if err != nil {
		return result, fmt.Errorf("go vet found issues: %w\n%s", err, result)
	}
	if result == "" {
		return "go vet: no issues found", nil
	}
	return result, nil
}
