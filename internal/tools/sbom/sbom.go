// Package sbom generates a Software Bill of Materials from go.mod content.
package sbom

import (
	"fmt"
	"strings"
	"time"
)

// Args defines the arguments for the sbom_generate tool.
type Args struct {
	GoMod string `json:"gomod" jsonschema:"the content of a go.mod file"`
}

// Generate parses go.mod content and returns an SBOM in Markdown format.
func Generate(args Args) (string, error) {
	lines := strings.Split(args.GoMod, "\n")

	var moduleName, goVersion string
	var deps []struct{ path, version string }
	inRequire := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			moduleName = strings.TrimPrefix(line, "module ")
		} else if strings.HasPrefix(line, "go ") {
			goVersion = strings.TrimPrefix(line, "go ")
		} else if line == "require (" {
			inRequire = true
		} else if line == ")" {
			inRequire = false
		} else if inRequire && line != "" && !strings.HasPrefix(line, "//") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				indirect := strings.Contains(line, "// indirect")
				tag := ""
				if indirect {
					tag = " (indirect)"
				}
				deps = append(deps, struct{ path, version string }{parts[0] + tag, parts[1]})
			}
		}
	}

	if moduleName == "" {
		return "", fmt.Errorf("could not parse module name from go.mod")
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# SBOM — %s\n\n", moduleName)
	fmt.Fprintf(&b, "- **Generated:** %s\n", time.Now().UTC().Format(time.RFC3339))
	fmt.Fprintf(&b, "- **Go version:** %s\n", goVersion)
	fmt.Fprintf(&b, "- **Dependencies:** %d\n\n", len(deps))
	fmt.Fprintln(&b, "| Module | Version |")
	fmt.Fprintln(&b, "|--------|---------|")
	for _, d := range deps {
		fmt.Fprintf(&b, "| %s | %s |\n", d.path, d.version)
	}
	return b.String(), nil
}
