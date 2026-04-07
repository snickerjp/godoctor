// Package gotest runs go test on a specified package.
package gotest

import (
	"bytes"
	"fmt"
	"os/exec"
)

// Args defines the arguments for the go_test tool.
type Args struct {
	Package string `json:"package" jsonschema:"the Go package to test, e.g. ./internal/tools/..."`
	Verbose bool   `json:"verbose,omitempty" jsonschema:"run tests with -v flag"`
}

// Run executes go test on the specified package and returns the output.
func Run(args Args) (string, error) {
	cmdArgs := []string{"test"}
	if args.Verbose {
		cmdArgs = append(cmdArgs, "-v")
	}
	cmdArgs = append(cmdArgs, args.Package)

	cmd := exec.Command("go", cmdArgs...)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	result := out.String() + stderr.String()
	if err != nil {
		return result, fmt.Errorf("go test failed: %w\n%s", err, result)
	}
	return result, nil
}
