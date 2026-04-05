package docs

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

// ReadDocsArgs defines the arguments for the read_docs tool.
type ReadDocsArgs struct {
	Package string `json:"package" jsonschema:"the package to read documentation for"`
	Symbol  string `json:"symbol,omitempty" jsonschema:"the optional symbol to read documentation for"`
}

// ReadDocs invokes 'go doc' to retrieve documentation.
func ReadDocs(args ReadDocsArgs) (string, error) {
	docPath := args.Package
	if args.Symbol != "" {
		docPath = fmt.Sprintf("%s.%s", args.Package, args.Symbol)
	}

	// Try local 'go doc' first
	cmd := exec.Command("go", "doc", docPath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		return out.String(), nil
	}

	// Fallback to External API (pkg.go.dev via r.jina.ai for clean markdown)
	url := fmt.Sprintf("https://r.jina.ai/https://pkg.go.dev/%s?tab=doc", args.Package)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch from external API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("external API returned status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	markdown := string(body)

	// If a symbol is requested, try to find it in the markdown
	if args.Symbol != "" {
		// Simple heuristic: search for the symbol name in the markdown
		// In Jina markdown, headers or code blocks usually contain the symbol
		lines := strings.Split(markdown, "\n")
		var result []string
		found := false
		for _, line := range lines {
			if strings.Contains(line, "### "+args.Symbol) || strings.Contains(line, "## "+args.Symbol) || strings.Contains(line, "func "+args.Symbol) || strings.Contains(line, "type "+args.Symbol) {
				found = true
			}
			if found {
				result = append(result, line)
				// Stop at next major header if we've found something
				if len(result) > 1 && (strings.HasPrefix(line, "## ") || strings.HasPrefix(line, "### ")) {
					// Check if this header is a different symbol
					if !strings.Contains(line, args.Symbol) {
						break
					}
				}
			}
		}
		if found {
			return strings.Join(result, "\n"), nil
		}
		return markdown, nil // Return full markdown if symbol search fails
	}

	return markdown, nil
}
